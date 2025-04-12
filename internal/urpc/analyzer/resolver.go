package analyzer

import (
	"fmt"
	"strings"

	"github.com/uforg/uforpc/internal/urpc/ast"
	"github.com/uforg/uforpc/internal/urpc/parser"
	"github.com/uforg/uforpc/internal/util/filepathutil"
)

// resolver is in charge of recursively resolve all imports from a given
// schema entry point and combine them into a single CombinedSchema.
//
// It verifies circular imports and ensures that all imported schemas exist
// and are parsable.
type resolver struct {
	fileProvider FileProvider
}

// resolverContext tracks the state of import resolution to detect circular imports
type resolverContext struct {
	// visitedFiles tracks files that have been processed to avoid duplicates
	visitedFiles map[string]*ast.Schema
	// importChain tracks the current import chain to detect circular imports
	importChain []string
	// diagnostics collects all diagnostics during resolution
	diagnostics []Diagnostic
}

// newResolver creates a new resolver. See resolver for more details.
func newResolver(fileProvider FileProvider) *resolver {
	return &resolver{fileProvider: fileProvider}
}

// resolve is the main entry point for resolving imports and combining schemas
//
// Returns:
//   - The combined schema.
//   - A list of diagnostics that occurred during the analysis.
//   - The first diagnostic converted to Error interface if any.
func (r *resolver) resolve(entryPointFilePath string) (CombinedSchema, []Diagnostic, error) {
	// Initialize the import context
	ctx := &resolverContext{
		visitedFiles: make(map[string]*ast.Schema),
		importChain:  []string{},
		diagnostics:  []Diagnostic{},
	}

	// Normalize the entry point file path
	normFilePath, err := filepathutil.Normalize("", entryPointFilePath)
	if err != nil {
		// Create a diagnostic for the normalization error
		diag := Diagnostic{
			Positions: Positions{
				Pos:    ast.Position{Filename: entryPointFilePath, Line: 1, Column: 1, Offset: 0},
				EndPos: ast.Position{Filename: entryPointFilePath, Line: 1, Column: 1, Offset: 0},
			},
			Message: fmt.Sprintf("error normalizing entry point file path: %v", err),
		}
		ctx.diagnostics = append(ctx.diagnostics, diag)
		return CombinedSchema{}, ctx.diagnostics, diag
	}

	// Resolve the entry point and all its imports
	combinedSchema := r.resolveFile(normFilePath, ctx)

	// Create an empty schema if resolution failed completely
	if combinedSchema == nil {
		combinedSchema = &ast.Schema{}
	}

	// Collect all definitions from the combined schema
	ruleDefs := make(map[string]Positions)
	typeDefs := make(map[string]Positions)
	procDefs := make(map[string]Positions)

	// Collect rule definitions
	for _, rule := range combinedSchema.GetRules() {
		ruleDefs[rule.Name] = Positions{
			Pos:    rule.Pos,
			EndPos: rule.EndPos,
		}
	}

	// Collect type definitions
	for _, typeDecl := range combinedSchema.GetTypes() {
		typeDefs[typeDecl.Name] = Positions{
			Pos:    typeDecl.Pos,
			EndPos: typeDecl.EndPos,
		}
	}

	// Collect procedure definitions
	for _, proc := range combinedSchema.GetProcs() {
		procDefs[proc.Name] = Positions{
			Pos:    proc.Pos,
			EndPos: proc.EndPos,
		}
	}

	// Create the final combined schema
	result := CombinedSchema{
		Schema:   combinedSchema,
		RuleDefs: ruleDefs,
		TypeDefs: typeDefs,
		ProcDefs: procDefs,
	}

	// Return the first diagnostic as error if any
	if len(ctx.diagnostics) > 0 {
		return result, ctx.diagnostics, ctx.diagnostics[0]
	}
	return result, nil, nil
}

// resolveFile resolves a single file and all its imports
func (r *resolver) resolveFile(filePath string, ctx *resolverContext) *ast.Schema {
	// Check if we've already processed this file
	if schema, exists := ctx.visitedFiles[filePath]; exists {
		// Check for circular imports
		if r.detectCircularImport(filePath, ctx) {
			// Return the already processed schema to continue processing
			return schema
		}

		// If we've already processed this file but it's not in our import chain,
		// just return the processed schema
		return schema
	}

	// Read and parse the file
	content, _, err := r.fileProvider.GetFileAndHash("", filePath)
	if err != nil {
		// Create a diagnostic for the file reading error
		diag := Diagnostic{
			Positions: Positions{
				Pos:    ast.Position{Filename: filePath, Line: 1, Column: 1, Offset: 0},
				EndPos: ast.Position{Filename: filePath, Line: 1, Column: 1, Offset: 0},
			},
			Message: fmt.Sprintf("error reading file: %v", err),
		}
		ctx.diagnostics = append(ctx.diagnostics, diag)
		return nil
	}

	schema, err := parser.ParserInstance.ParseString(filePath, content)
	if err != nil {
		// Assert parser error if possible
		if parserErr, ok := err.(parser.Error); ok {
			ctx.diagnostics = append(ctx.diagnostics, Diagnostic{
				Positions: Positions{
					Pos:    parserErr.Position(),
					EndPos: parserErr.Position(),
				},
				Message: fmt.Sprintf("error parsing file: %v", parserErr),
			})
			return nil
		}

		// Create generic diagnostic for other errors
		diag := Diagnostic{
			Positions: Positions{
				Pos:    ast.Position{Filename: filePath, Line: 1, Column: 1, Offset: 0},
				EndPos: ast.Position{Filename: filePath, Line: 1, Column: 1, Offset: 0},
			},
			Message: fmt.Sprintf("error parsing file: %v", err),
		}
		ctx.diagnostics = append(ctx.diagnostics, diag)
		return nil
	}

	// Validate the file version, returns false if version is invalid
	// and all diagnostics are added to the context by the function
	if !r.validateFileVersion(schema, ctx) {
		return nil
	}

	// Add this file to the visited files and import chain
	ctx.visitedFiles[filePath] = schema
	ctx.importChain = append(ctx.importChain, filePath)

	// Create a new schema with the same positions as the original
	combinedSchema := &ast.Schema{}
	combinedSchema.Positions = schema.Positions

	// Process all children of the schema in order, both imports and non-imports
	for _, child := range schema.Children {
		// Skip version because it's already validated
		if child.Kind() == ast.SchemaChildKindVersion {
			continue
		}

		if child.Kind() == ast.SchemaChildKindImport {
			importPath, err := filepathutil.Normalize(filePath, child.Import.Path)
			if err != nil {
				// Create a diagnostic for the import path normalization error
				ctx.diagnostics = append(ctx.diagnostics, Diagnostic{
					Positions: Positions{
						Pos:    child.Import.Pos,
						EndPos: child.Import.EndPos,
					},
					Message: fmt.Sprintf("error resolving import path: %v", err),
				})
				continue
			}

			// Recursively resolve the imported file
			importedSchema := r.resolveFile(importPath, ctx)
			if importedSchema == nil {
				// Skip this import if resolution failed (diagnostics already added)
				continue
			}

			// Add all non-import children from the imported schema to the combined schema
			for _, child := range importedSchema.Children {
				if child.Kind() != ast.SchemaChildKindImport {
					combinedSchema.Children = append(combinedSchema.Children, child)
				}
			}
		}

		if child.Kind() != ast.SchemaChildKindImport {
			combinedSchema.Children = append(combinedSchema.Children, child)
		}
	}

	// Remove this file from the import chain (backtracking)
	ctx.importChain = ctx.importChain[:len(ctx.importChain)-1]

	// Update the visited files map with the combined schema
	ctx.visitedFiles[filePath] = combinedSchema

	return combinedSchema
}

func (r *resolver) validateFileVersion(schema *ast.Schema, ctx *resolverContext) bool {
	if len(schema.Children) == 0 {
		return true
	}

	firstChild := schema.Children[0]
	if firstChild.Kind() != ast.SchemaChildKindVersion {
		ctx.diagnostics = append(ctx.diagnostics, Diagnostic{
			Positions: Positions{
				Pos:    firstChild.Pos,
				EndPos: firstChild.EndPos,
			},
			Message: "the first statement must be a version statement",
		})
		return false
	}
	if firstChild.Version.Number != 1 {
		ctx.diagnostics = append(ctx.diagnostics, Diagnostic{
			Positions: Positions{
				Pos:    firstChild.Version.Pos,
				EndPos: firstChild.Version.EndPos,
			},
			Message: "at the moment, the only supported version is 1",
		})
		return false
	}

	versions := schema.GetVersions()
	for i, version := range versions {
		if i == 0 {
			continue
		}

		ctx.diagnostics = append(ctx.diagnostics, Diagnostic{
			Positions: Positions{
				Pos:    version.Pos,
				EndPos: version.EndPos,
			},
			Message: "version statement already defined for this schema",
		})
	}

	return len(versions) == 1
}

// detectCircularImport checks if the given file path creates a circular import in the current import chain.
// If a circular import is detected, it adds a diagnostic and returns true.
func (r *resolver) detectCircularImport(filePath string, ctx *resolverContext) bool {
	for i, path := range ctx.importChain {
		if path == filePath {
			// Found a circular import - extract the circular chain
			circularChain := append(ctx.importChain[i:], filePath)

			// Get the first file in the circular chain
			firstFileInCircle := ctx.importChain[i]

			// Find the import statement in the first file that starts the circle
			if firstSchema, exists := ctx.visitedFiles[firstFileInCircle]; exists && firstSchema != nil {
				// Find the import that points to the next file in the circle
				nextFileInCircle := ctx.importChain[i+1]

				for _, importNode := range firstSchema.GetImports() {
					importPath, err := filepathutil.Normalize(firstFileInCircle, importNode.Path)
					if err == nil && importPath == nextFileInCircle {
						// Create a diagnostic with the position of the import statement in the first file
						ctx.diagnostics = append(ctx.diagnostics, Diagnostic{
							Positions: Positions{
								Pos:    importNode.Pos,
								EndPos: importNode.EndPos,
							},
							Message: fmt.Sprintf("circular import detected: %s", strings.Join(circularChain, " -> ")),
						})
						break
					}
				}
			}
			return true
		}
	}
	return false
}
