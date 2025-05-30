package analyzer

import (
	"errors"
	"fmt"
	"os"

	"github.com/uforg/uforpc/urpc/internal/urpc/ast"
	"github.com/uforg/uforpc/urpc/internal/urpc/parser"
	"github.com/uforg/uforpc/urpc/internal/util/filepathutil"
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
	// and correctly handle circular imports
	visitedFiles map[string]*ast.Schema
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
	combinedSchema := r.resolveFile(nil, normFilePath, ctx)

	// Create an empty schema if resolution failed completely
	if combinedSchema == nil {
		combinedSchema = &ast.Schema{}
	}

	// Resolve external docstrings, if any error happens, or if the
	// referenced markdown file is not found, the function will
	// add a diagnostic to the context.
	r.resolveExternalDocstrings(combinedSchema, ctx)

	// Collect all declarations from the combined schema
	ruleDecls := make(map[string]*ast.RuleDecl)
	typeDecls := make(map[string]*ast.TypeDecl)
	procDecls := make(map[string]*ast.ProcDecl)

	// Collect rule declarations
	for _, rule := range combinedSchema.GetRules() {
		ruleDecls[rule.Name] = rule
	}

	// Collect type declarations
	for _, typeDecl := range combinedSchema.GetTypes() {
		typeDecls[typeDecl.Name] = typeDecl
	}

	// Collect procedure declarations
	for _, proc := range combinedSchema.GetProcs() {
		procDecls[proc.Name] = proc
	}

	// Create the final combined schema
	result := CombinedSchema{
		Schema:    combinedSchema,
		RuleDecls: ruleDecls,
		TypeDecls: typeDecls,
		ProcDecls: procDecls,
	}

	// Return the first diagnostic as error if any
	if len(ctx.diagnostics) > 0 {
		return result, ctx.diagnostics, ctx.diagnostics[0]
	}
	return result, nil, nil
}

// resolveFile resolves a single file and all its imports
func (r *resolver) resolveFile(parentImport *ast.Import, filePath string, ctx *resolverContext) *ast.Schema {
	// Check if we've already processed this file
	if schema, exists := ctx.visitedFiles[filePath]; exists {
		return schema
	}

	// Read and parse the file
	content, _, err := r.fileProvider.GetFileAndHash("", filePath)
	if err != nil {
		pos := ast.Position{Filename: filePath, Line: 1, Column: 1, Offset: 0}
		endPos := ast.Position{Filename: filePath, Line: 1, Column: 1, Offset: 0}

		if parentImport != nil {
			pos = parentImport.Pos
			endPos = parentImport.EndPos
		}

		// Create a diagnostic for the file reading error
		diag := Diagnostic{
			Positions: Positions{
				Pos:    pos,
				EndPos: endPos,
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
				Message: parserErr.Message(),
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

	// Add this file to the visited files
	ctx.visitedFiles[filePath] = schema

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

			// If the import is already resolved, skip the resolution
			if _, exists := ctx.visitedFiles[importPath]; exists {
				continue
			}

			// Recursively resolve the imported file
			importedSchema := r.resolveFile(child.Import, importPath, ctx)
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

			continue
		}

		combinedSchema.Children = append(combinedSchema.Children, child)
	}

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
			Message: "version statement already declared for this schema",
		})
	}

	return len(versions) == 1
}

// resolveExternalDocstrings resolves external docstrings in the combined schema
// by reading the content of the referenced Markdown files and updating the docstring values.
func (r *resolver) resolveExternalDocstrings(combinedSchema *ast.Schema, ctx *resolverContext) {
	for _, docstring := range combinedSchema.GetDocstrings() {
		r.resolveExternalDocstring(docstring, ctx)
	}

	for _, rule := range combinedSchema.GetRules() {
		if rule.Docstring != nil {
			r.resolveExternalDocstring(rule.Docstring, ctx)
		}
	}

	for _, typeDecl := range combinedSchema.GetTypes() {
		if typeDecl.Docstring != nil {
			r.resolveExternalDocstring(typeDecl.Docstring, ctx)
		}
	}

	for _, proc := range combinedSchema.GetProcs() {
		if proc.Docstring != nil {
			r.resolveExternalDocstring(proc.Docstring, ctx)
		}
	}
}

// resolveExternalDocstring is the logic to resolve a single external docstring
// behind resolveExternalDocstrings.
func (r *resolver) resolveExternalDocstring(docstring *ast.Docstring, ctx *resolverContext) {
	externalPath, isExternal := docstring.GetExternal()
	if !isExternal {
		return
	}

	content, _, err := r.fileProvider.GetFileAndHash(docstring.Pos.Filename, externalPath)
	if errors.Is(err, os.ErrNotExist) {
		ctx.diagnostics = append(ctx.diagnostics, Diagnostic{
			Positions: Positions{
				Pos:    docstring.Pos,
				EndPos: docstring.EndPos,
			},
			Message: fmt.Sprintf("external markdown file not found: %s", externalPath),
		})
		return
	}
	if err != nil {
		ctx.diagnostics = append(ctx.diagnostics, Diagnostic{
			Positions: Positions{
				Pos:    docstring.Pos,
				EndPos: docstring.EndPos,
			},
			Message: fmt.Sprintf("error reading external markdown file: %v", err),
		})
		return
	}

	docstring.Value = content
}
