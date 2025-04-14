// Package typeflattener provides functionality to flatten extended types in URPC schemas.
// It resolves type extensions by copying fields from extended types to the base types,
// so that representations that don't support inheritance can work with the schema.
package typeflattener

import (
	"maps"

	"github.com/uforg/uforpc/internal/urpc/ast"
)

// Flatten flattens all extended types in the schema.
//
// For each type that extends other types, it copies all fields from the extended types
// to the base type. Fields are copied in the order they appear in the extends list,
// and then the type's own fields are added.
//
// This function modifies the input schema directly.
//
// Returns:
//   - The schema with flattened types.
func Flatten(schema *ast.Schema) *ast.Schema {
	tf := newTypeFlattener(schema)
	return tf.flatten()
}

// typeFlattener is responsible for flattening extended types in a URPC schema.
type typeFlattener struct {
	schema  *ast.Schema
	typeMap map[string]*ast.TypeDecl
}

// newTypeFlattener creates a new typeFlattener instance.
func newTypeFlattener(schema *ast.Schema) *typeFlattener {
	// Build a map of type names to type declarations for quick lookup
	typeMap := make(map[string]*ast.TypeDecl)
	for _, typeDecl := range schema.GetTypes() {
		typeMap[typeDecl.Name] = typeDecl
	}

	return &typeFlattener{
		schema:  schema,
		typeMap: typeMap,
	}
}

// flatten flattens all extended types in the schema.
func (f *typeFlattener) flatten() *ast.Schema {
	// Create a map to store the flattened fields for each type
	flattenedFields := make(map[string][]*ast.FieldOrComment)

	// For each type in the schema
	for _, typeDecl := range f.schema.GetTypes() {
		// If the type doesn't extend other types, there's nothing to flatten
		if len(typeDecl.Extends) == 0 {
			continue
		}

		// Collect all fields from extended types
		allFields := f.collectFieldsFromExtendedTypes(typeDecl)

		// Store the flattened fields for this type
		flattenedFields[typeDecl.Name] = allFields
	}

	// Apply the flattened fields to each type in the schema
	for typeName, fields := range flattenedFields {
		typeDecl := f.typeMap[typeName]
		// Add the type's own fields at the end
		fields = append(fields, typeDecl.Children...)
		// Update the type's fields
		typeDecl.Children = fields
		// Remove the extends property since it's been flattened
		typeDecl.Extends = nil
	}

	return f.schema
}

// collectFieldsFromExtendedTypes collects all fields from the types extended
// by a given type, in the order they appear in the extends list.
func (f *typeFlattener) collectFieldsFromExtendedTypes(typeDecl *ast.TypeDecl) []*ast.FieldOrComment {
	visited := make(map[string]bool)
	return f.collectFieldsRecursive(typeDecl, visited)
}

// hasCircularDependency checks if there is a circular dependency between types.
// It returns true if the target type is found in the extends chain of the current type.
func (f *typeFlattener) hasCircularDependency(currentType *ast.TypeDecl, targetName string, visited map[string]bool) bool {
	visited[currentType.Name] = true

	for _, extendTypeName := range currentType.Extends {
		// If we found the target, there is a circular dependency
		if extendTypeName == targetName {
			return true
		}

		// Skip if already visited
		if visited[extendTypeName] {
			continue
		}

		// Get the extended type
		extendedType, exists := f.typeMap[extendTypeName]
		if !exists {
			continue
		}

		// Recursively check the extended type
		if f.hasCircularDependency(extendedType, targetName, make(map[string]bool)) {
			return true
		}
	}

	return false
}

// collectFieldsRecursive is a helper function for collectFieldsFromExtendedTypes
// that avoids cycles in recursion.
func (f *typeFlattener) collectFieldsRecursive(typeDecl *ast.TypeDecl, visited map[string]bool) []*ast.FieldOrComment {
	var allFields []*ast.FieldOrComment

	// Mark this type as visited to avoid cycles
	visited[typeDecl.Name] = true

	// For each extended type
	for _, extendTypeName := range typeDecl.Extends {
		// Avoid cycles
		if visited[extendTypeName] {
			continue
		}

		// Get the extended type
		extendedType, exists := f.typeMap[extendTypeName]
		if !exists {
			continue
		}

		// Check for circular dependencies
		if f.hasCircularDependency(extendedType, typeDecl.Name, make(map[string]bool)) {
			continue
		}

		// If the extended type also extends other types, collect its fields first
		if len(extendedType.Extends) > 0 {
			branchVisited := make(map[string]bool)
			maps.Copy(branchVisited, visited)
			extendedFields := f.collectFieldsRecursive(extendedType, branchVisited)
			allFields = append(allFields, extendedFields...)
		}

		// Add the fields of the extended type, avoiding duplicates
		for _, child := range extendedType.Children {
			// Handle comments
			if child.Comment != nil {
				newChild := &ast.FieldOrComment{
					Comment: &ast.Comment{},
				}
				if child.Comment.Simple != nil {
					newChild.Comment.Simple = new(string)
					*newChild.Comment.Simple = *child.Comment.Simple
				}
				if child.Comment.Block != nil {
					newChild.Comment.Block = new(string)
					*newChild.Comment.Block = *child.Comment.Block
				}
				allFields = append(allFields, newChild)
				continue
			}

			// Handle fields
			if child.Field != nil {
				// Skip if this is a field that already exists
				duplicate := false
				for _, existingChild := range allFields {
					if existingChild.Field != nil && existingChild.Field.Name == child.Field.Name {
						duplicate = true
						break
					}
				}
				if duplicate {
					continue
				}

				// Create a copy of the field
				newChild := &ast.FieldOrComment{
					Field: &ast.Field{},
				}
				*newChild.Field = *child.Field

				// Copy field children if they exist
				if len(child.Field.Children) > 0 {
					newChild.Field.Children = make([]*ast.FieldChild, len(child.Field.Children))
					for i, fieldChild := range child.Field.Children {
						newFieldChild := &ast.FieldChild{}
						*newFieldChild = *fieldChild
						newChild.Field.Children[i] = newFieldChild
					}
				}

				// Add the field to the list
				allFields = append(allFields, newChild)
			}
		}
	}

	return allFields
}
