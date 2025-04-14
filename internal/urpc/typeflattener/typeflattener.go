// Package typeflattener provides functionality to flatten extended types in URPC schemas.
// It resolves type extensions by copying fields from extended types to the base types,
// so that representations that don't support inheritance can work with the schema.
package typeflattener

import (
	"encoding/json"
	"maps"

	"github.com/uforg/uforpc/internal/urpc/ast"
)

// TypeFlattener is responsible for flattening extended types in a URPC schema.
type TypeFlattener struct {
	schema *ast.Schema
	// Map of type name to type declaration for quick lookup
	typeMap map[string]*ast.TypeDecl
}

// NewTypeFlattener creates a new TypeFlattener instance.
func NewTypeFlattener(schema *ast.Schema) *TypeFlattener {
	// Build a map of type names to type declarations for quick lookup
	typeMap := make(map[string]*ast.TypeDecl)
	for _, typeDecl := range schema.GetTypes() {
		typeMap[typeDecl.Name] = typeDecl
	}

	return &TypeFlattener{
		schema:  schema,
		typeMap: typeMap,
	}
}

// Flatten flattens all extended types in the schema.
//
// For each type that extends other types, it copies all fields from the extended types
// to the base type. Fields are copied in the order they appear in the extends list,
// and then the type's own fields are added.
//
// This function creates a deep copy of the original schema to avoid side effects.
//
// Returns:
//   - A new schema with flattened types.
func (f *TypeFlattener) Flatten() *ast.Schema {
	// Create a deep copy of the original schema using JSON marshaling/unmarshaling
	copiedSchema := f.deepCopySchema()

	// Create a new type map for the copied schema
	copiedTypeMap := make(map[string]*ast.TypeDecl)
	for _, typeDecl := range copiedSchema.GetTypes() {
		copiedTypeMap[typeDecl.Name] = typeDecl
	}

	// Create a map to store the flattened fields for each type
	flattenedFields := make(map[string][]*ast.FieldOrComment)

	// For each type in the copied schema
	for _, typeDecl := range copiedSchema.GetTypes() {
		// If the type doesn't extend other types, there's nothing to flatten
		if len(typeDecl.Extends) == 0 {
			continue
		}

		// Collect all fields from extended types
		allFields := f.collectFieldsFromExtendedTypes(typeDecl)

		// Store the flattened fields for this type
		flattenedFields[typeDecl.Name] = allFields
	}

	// Apply the flattened fields to each type in the copied schema
	for typeName, fields := range flattenedFields {
		typeDecl := copiedTypeMap[typeName]
		// Add the type's own fields at the end
		fields = append(fields, typeDecl.Children...)
		// Update the type's fields
		typeDecl.Children = fields
		// Remove the extends property since it's been flattened
		typeDecl.Extends = nil
	}

	return copiedSchema
}

// deepCopySchema creates a deep copy of the original schema.
//
// Returns:
//   - A new schema that is a deep copy of the original.
func (f *TypeFlattener) deepCopySchema() *ast.Schema {
	// Marshal the original schema to JSON
	data, err := json.Marshal(f.schema)
	if err != nil {
		// If marshaling fails, return a new empty schema
		return &ast.Schema{}
	}

	// Unmarshal the JSON data into a new schema
	var copiedSchema ast.Schema
	err = json.Unmarshal(data, &copiedSchema)
	if err != nil {
		// If unmarshaling fails, return a new empty schema
		return &ast.Schema{}
	}

	return &copiedSchema
}

// collectFieldsFromExtendedTypes collects all fields from the types extended
// by a given type, in the order they appear in the extends list.
//
// Returns:
//   - A list of FieldOrComment containing all fields from extended types.
func (f *TypeFlattener) collectFieldsFromExtendedTypes(typeDecl *ast.TypeDecl) []*ast.FieldOrComment {
	// Use a map to avoid cycles in recursion
	visited := make(map[string]bool)
	return f.collectFieldsRecursive(typeDecl, visited)
}

// hasCircularDependency checks if there is a circular dependency between types.
// It returns true if the target type is found in the extends chain of the current type.
//
// Parameters:
//   - currentType: The current type to check
//   - targetName: The name of the target type to find in the extends chain
//   - visited: A map of visited types to avoid infinite recursion
//
// Returns:
//   - true if there is a circular dependency, false otherwise
func (f *TypeFlattener) hasCircularDependency(currentType *ast.TypeDecl, targetName string, visited map[string]bool) bool {
	// Mark this type as visited
	visited[currentType.Name] = true

	// Check if any of the extended types is the target
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
		if f.hasCircularDependency(extendedType, targetName, visited) {
			return true
		}
	}

	return false
}

// collectFieldsRecursive is a helper function for collectFieldsFromExtendedTypes
// that avoids cycles in recursion.
//
// Returns:
//   - A list of FieldOrComment containing all fields from extended types.
func (f *TypeFlattener) collectFieldsRecursive(typeDecl *ast.TypeDecl, visited map[string]bool) []*ast.FieldOrComment {
	var allFields []*ast.FieldOrComment

	// Mark this type as visited to avoid cycles
	visited[typeDecl.Name] = true

	// For each extended type
	for _, extendTypeName := range typeDecl.Extends {
		// Avoid cycles
		if visited[extendTypeName] {
			// Skip this type to avoid infinite recursion
			continue
		}

		// Get the extended type
		extendedType, exists := f.typeMap[extendTypeName]
		if !exists {
			// If the extended type doesn't exist, ignore it (this should have been
			// detected in the semantic analysis phase)
			continue
		}

		// Create a copy of the visited map for this branch of recursion
		branchVisited := make(map[string]bool)
		maps.Copy(branchVisited, visited)

		// Check for circular dependencies
		if f.hasCircularDependency(extendedType, typeDecl.Name, make(map[string]bool)) {
			// Skip this type to avoid infinite recursion
			continue
		}

		// If the extended type also extends other types, collect its fields first
		if len(extendedType.Extends) > 0 {
			extendedFields := f.collectFieldsRecursive(extendedType, branchVisited)
			allFields = append(allFields, extendedFields...)
		}

		// Add the fields of the extended type, avoiding duplicates
		for _, child := range extendedType.Children {
			// Create a deep copy of the child
			copiedChild := &ast.FieldOrComment{}

			// Copy comment if it exists
			if child.Comment != nil {
				copiedChild.Comment = &ast.Comment{}
				*copiedChild.Comment = *child.Comment
				allFields = append(allFields, copiedChild)
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

				// Create a deep copy of the field
				copiedChild.Field = &ast.Field{}
				*copiedChild.Field = *child.Field

				// Copy field children if they exist
				if len(child.Field.Children) > 0 {
					copiedChild.Field.Children = make([]*ast.FieldChild, len(child.Field.Children))
					for i, fieldChild := range child.Field.Children {
						copiedFieldChild := &ast.FieldChild{}
						*copiedFieldChild = *fieldChild
						copiedChild.Field.Children[i] = copiedFieldChild
					}
				}

				// Add the field to the list
				allFields = append(allFields, copiedChild)
			}
		}
	}

	return allFields
}
