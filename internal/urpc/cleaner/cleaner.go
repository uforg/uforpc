// Package cleaner provides functionality to remove unused rules and custom types from URPC schemas.
package cleaner

import (
	"github.com/uforg/uforpc/internal/urpc/ast"
)

// Clean removes unused rules and custom types from a URPC schema.
// A rule or custom type is considered unused if it's not referenced by any field
// in another rule, type, or procedure.
//
// This function modifies the input schema directly and returns it.
func Clean(schema *ast.Schema) *ast.Schema {
	// Find referenced types and remove unused ones (first pass)
	//
	// This removes the types that are not used in procedures
	usedTypes := findReferencedTypes(schema)
	removeUnusedTypes(schema, usedTypes)

	// Find referenced types and remove unused ones (second pass)
	//
	// This removes the types that was used in the removed types
	// in the first pass
	usedTypes = findReferencedTypes(schema)
	removeUnusedTypes(schema, usedTypes)

	// Find referenced rules and remove unused ones
	usedRules := findReferencedRules(schema)
	removeUnusedRules(schema, usedRules)

	return schema
}

// findReferencedTypes identifies all custom types that are being used in the schema.
func findReferencedTypes(schema *ast.Schema) map[string]bool {
	usedTypes := make(map[string]bool)

	// Find types referenced in procedures
	for _, proc := range schema.GetProcs() {
		for _, child := range proc.Children {
			// Check input fields
			if child.Input != nil {
				inputFields := extractFields(child.Input.Children)
				for _, field := range inputFields {
					if field.Type.Base.Named != nil {
						typeName := *field.Type.Base.Named
						usedTypes[typeName] = true
					}
				}
			}

			// Check output fields
			if child.Output != nil {
				outputFields := extractFields(child.Output.Children)
				for _, field := range outputFields {
					if field.Type.Base.Named != nil {
						typeName := *field.Type.Base.Named
						usedTypes[typeName] = true
					}
				}
			}
		}
	}

	// Find types referenced in other types
	for _, typeDecl := range schema.GetTypes() {
		// Mark extended types as used
		if len(typeDecl.Extends) > 0 {
			for _, extendTypeName := range typeDecl.Extends {
				usedTypes[extendTypeName] = true
			}
		}

		// Check fields in types
		typeFields := extractFields(typeDecl.Children)
		for _, field := range typeFields {
			if field.Type.Base.Named != nil {
				typeName := *field.Type.Base.Named
				usedTypes[typeName] = true
			}
		}
	}

	return usedTypes
}

// removeUnusedTypes removes unused types from the schema.
func removeUnusedTypes(schema *ast.Schema, usedTypes map[string]bool) {
	var newChildren []*ast.SchemaChild

	for _, child := range schema.Children {
		if child.Kind() == ast.SchemaChildKindType {
			if usedTypes[child.Type.Name] {
				newChildren = append(newChildren, child)
			}
		} else {
			newChildren = append(newChildren, child)
		}
	}

	schema.Children = newChildren
}

// findReferencedRules identifies all rules that are being used in the schema.
func findReferencedRules(schema *ast.Schema) map[string]bool {
	usedRules := make(map[string]bool)

	// Find rules used in types
	for _, typeDecl := range schema.GetTypes() {
		typeFields := extractFields(typeDecl.Children)
		for _, field := range typeFields {
			// Rules applied directly to fields
			for _, child := range field.Children {
				if child.Rule != nil {
					usedRules[child.Rule.Name] = true
				}
			}
		}
	}

	// Find rules used in procedures
	for _, proc := range schema.GetProcs() {
		for _, child := range proc.Children {
			// Input fields
			if child.Input != nil {
				inputFields := extractFields(child.Input.Children)
				for _, field := range inputFields {
					for _, fieldChild := range field.Children {
						if fieldChild.Rule != nil {
							usedRules[fieldChild.Rule.Name] = true
						}
					}
				}
			}

			// Output fields
			if child.Output != nil {
				outputFields := extractFields(child.Output.Children)
				for _, field := range outputFields {
					for _, fieldChild := range field.Children {
						if fieldChild.Rule != nil {
							usedRules[fieldChild.Rule.Name] = true
						}
					}
				}
			}
		}
	}

	return usedRules
}

// removeUnusedRules removes unused rules from the schema.
func removeUnusedRules(schema *ast.Schema, usedRules map[string]bool) {
	var newChildren []*ast.SchemaChild

	for _, child := range schema.Children {
		if child.Kind() == ast.SchemaChildKindRule {
			if usedRules[child.Rule.Name] {
				newChildren = append(newChildren, child)
			}
		} else {
			newChildren = append(newChildren, child)
		}
	}

	schema.Children = newChildren
}

// extractFields extracts all fields from a slice of FieldOrComment
// it recursively checks inline objects and returns a flattened list
// of all named fields.
func extractFields(fieldOrComments []*ast.FieldOrComment) []*ast.Field {
	var fields []*ast.Field

	for _, foc := range fieldOrComments {
		if foc.Field == nil {
			continue
		}

		if foc.Field.Type.Base.Named != nil {
			fields = append(fields, foc.Field)
		}

		if foc.Field.Type.Base.Object != nil {
			inlineFields := extractFields(foc.Field.Type.Base.Object.Children)
			fields = append(fields, inlineFields...)
		}
	}

	return fields
}
