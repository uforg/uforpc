package pieces

import (
	"encoding/json"
	"fmt"
	"strings"
)

/** START FROM HERE **/

// -----------------------------------------------------------------------------
// Required field validator
// -----------------------------------------------------------------------------

// ValidateJSONPaths validates JSON data against a set of required paths.
// It supports recursive type references with the syntax field->Type.
func ValidateJSONPaths(data []byte, templates map[string][]string, rootType string) error {
	var jsonData any
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}

	// Track visited paths to prevent infinite recursion
	visited := make(map[string]bool)

	return validateType(jsonData, templates, rootType, "", "", visited)
}

// validateType validates an object against a template type
func validateType(data any, templates map[string][]string, typeName, basePath, currentPath string, visited map[string]bool) error {
	// Get the template for this type
	paths, ok := templates[typeName]
	if !ok {
		return fmt.Errorf("no template defined for type: %s", typeName)
	}

	// Validate each path in the template
	for _, path := range paths {
		fullPath := joinPath(basePath, path)

		// Skip this path if we've already validated it (prevents infinite recursion)
		visitKey := typeName + ":" + fullPath
		if visited[visitKey] {
			continue
		}
		visited[visitKey] = true

		// Check if this path has a type reference
		if strings.Contains(path, "->") {
			parts := strings.Split(path, "->")
			if len(parts) != 2 {
				return fmt.Errorf("invalid type reference: %s", path)
			}

			fieldPath := parts[0]
			refType := parts[1]

			// Check if this is a wildcard path
			if strings.Contains(fieldPath, "[*]") {
				if err := validateWildcardTypeRef(data, fieldPath, refType, templates, currentPath, visited); err != nil {
					return err
				}
			} else {
				// Regular path
				fieldData, err := resolvePath(data, fieldPath)
				if err != nil {
					return fmt.Errorf("%s: %v", joinPath(currentPath, fieldPath), err)
				}

				// New base path for recursive validation
				newBasePath := joinPath(currentPath, fieldPath)

				// Validate the field as the referenced type
				if err := validateType(fieldData, templates, refType, newBasePath, newBasePath, visited); err != nil {
					return err
				}
			}
		} else {
			// Normal path validation
			if err := validatePath(data, path, currentPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateWildcardTypeRef handles type references with wildcards like "posts[*]->Post"
func validateWildcardTypeRef(data any, fieldPath, refType string, templates map[string][]string,
	basePath string, visited map[string]bool) error {
	// Extract the array field name (before the [*])
	parts := strings.Split(fieldPath, "[*]")
	arrayField := parts[0]

	// Get the array
	arrayData, err := resolvePath(data, arrayField)
	if err != nil {
		return fmt.Errorf("%s: %v", joinPath(basePath, arrayField), err)
	}

	// Check if it's an array
	array, ok := arrayData.([]any)
	if !ok {
		return fmt.Errorf("%s: expected array, got %T", joinPath(basePath, arrayField), arrayData)
	}

	// Apply the type validation to each element in the array
	for i, elem := range array {
		elemPath := fmt.Sprintf("%s[%d]", joinPath(basePath, arrayField), i)

		// Apply remaining path parts if any
		elemData := elem
		if len(parts) > 1 && parts[1] != "" {
			// Handle any remaining path after the wildcard
			subPath := parts[1]
			subPath = strings.TrimPrefix(subPath, ".")

			// Resolve the remaining path
			if subPath != "" {
				elemData, err = resolvePath(elem, subPath)
				if err != nil {
					return fmt.Errorf("%s.%s: %v", elemPath, subPath, err)
				}
				elemPath = joinPath(elemPath, subPath)
			}
		}

		// Validate this element against the referenced type
		if err := validateType(elemData, templates, refType, elemPath, elemPath, visited); err != nil {
			return err
		}
	}

	return nil
}

// validatePath checks if a path exists in the data
func validatePath(data any, path, basePath string) error {
	parts := strings.Split(path, ".")

	// Handle array wildcard notation separately
	if strings.Contains(path, "[*]") {
		return validateArrayWildcard(data, parts, basePath)
	}

	// Normal path
	for i, part := range parts {
		// Skip empty parts
		if part == "" {
			continue
		}

		// Handle object
		obj, ok := data.(map[string]any)
		if !ok {
			return fmt.Errorf("%s: expected object, got %T", joinPath(basePath, strings.Join(parts[:i+1], ".")), data)
		}

		value, exists := obj[part]
		if !exists {
			return fmt.Errorf("%s: required field is missing", joinPath(basePath, strings.Join(parts[:i+1], ".")))
		}

		// Move to next part
		data = value
	}

	return nil
}

// validateArrayWildcard handles validation of array wildcard paths like "items[*].id"
func validateArrayWildcard(data any, parts []string, basePath string) error {
	if len(parts) == 0 {
		return nil
	}

	part := parts[0]

	// Not a wildcard part, process normally
	if !strings.Contains(part, "[*]") {
		obj, ok := data.(map[string]any)
		if !ok {
			return fmt.Errorf("%s: expected object, got %T", joinPath(basePath, part), data)
		}

		value, exists := obj[part]
		if !exists {
			return fmt.Errorf("%s: required field is missing", joinPath(basePath, part))
		}

		return validateArrayWildcard(value, parts[1:], joinPath(basePath, part))
	}

	// Extract field name from "field[*]"
	fieldName := strings.Split(part, "[")[0]

	obj, ok := data.(map[string]any)
	if !ok {
		return fmt.Errorf("%s: expected object, got %T", joinPath(basePath, fieldName), data)
	}

	arrayData, exists := obj[fieldName]
	if !exists {
		return fmt.Errorf("%s: required field is missing", joinPath(basePath, fieldName))
	}

	array, ok := arrayData.([]any)
	if !ok {
		return fmt.Errorf("%s: expected array, got %T", joinPath(basePath, fieldName), arrayData)
	}

	// If array is empty but we're validating field existence in array elements,
	// we should consider this valid (since there are no elements to validate)
	if len(array) == 0 && len(parts) > 1 {
		return nil
	}

	// Validate each array element
	for i, elem := range array {
		elemPath := fmt.Sprintf("%s[%d]", joinPath(basePath, fieldName), i)

		// Continue validation with the next path parts for this element
		if len(parts) > 1 {
			if err := validateArrayWildcard(elem, parts[1:], elemPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// resolvePath gets a value at the specified path
func resolvePath(data any, path string) (any, error) {
	if path == "" {
		return data, nil
	}

	parts := strings.Split(path, ".")

	for _, part := range parts {
		// Handle array wildcard - can't resolve a specific value
		if strings.Contains(part, "[*]") {
			return nil, fmt.Errorf("cannot resolve wildcard path")
		}

		// Handle object
		obj, ok := data.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected object, got %T", data)
		}

		value, exists := obj[part]
		if !exists {
			return nil, fmt.Errorf("field not found")
		}

		data = value
	}

	return data, nil
}

// joinPath combines path segments with dots, handling empty segments
func joinPath(base, path string) string {
	if base == "" {
		return path
	}
	if path == "" {
		return base
	}
	return base + "." + path
}
