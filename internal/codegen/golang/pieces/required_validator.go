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
	if data == nil {
		return fmt.Errorf("%s: required field is missing", currentPath)
	}

	paths, ok := templates[typeName]
	if !ok {
		return fmt.Errorf("no template defined for type: %s", typeName)
	}

	obj, ok := data.(map[string]any)
	if !ok {
		return fmt.Errorf("%s: expected object, got %T", currentPath, data)
	}

	// Create a copy of visited map for this type
	typeVisited := make(map[string]bool)
	for k, v := range visited {
		typeVisited[k] = v
	}

	// Check if we've already visited this type at this path
	visitKey := fmt.Sprintf("%s:%s", typeName, currentPath)
	if typeVisited[visitKey] {
		// If we've already visited this type at this path, we still need to validate required fields
		// but we don't need to validate type references again to prevent infinite recursion
		for _, path := range paths {
			if !strings.Contains(path, "->") {
				err := validatePath(obj, path, currentPath)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	typeVisited[visitKey] = true

	// First validate required fields
	for _, path := range paths {
		if !strings.Contains(path, "->") {
			err := validatePath(obj, path, currentPath)
			if err != nil {
				return err
			}
		}
	}

	// Then validate type references
	for _, path := range paths {
		if !strings.Contains(path, "->") {
			continue
		}

		parts := strings.Split(path, "->")
		if len(parts) != 2 {
			return fmt.Errorf("invalid type reference: %s", path)
		}

		fieldPath := parts[0]
		refType := parts[1]

		if strings.Contains(fieldPath, "[*]") {
			arrayField := strings.Split(fieldPath, "[*]")[0]
			value, exists := obj[arrayField]
			if !exists {
				continue
			}

			err := validateWildcardTypeRef(value, arrayField, refType, templates, currentPath, typeVisited)
			if err != nil {
				return err
			}
		} else {
			value, exists := obj[fieldPath]
			if !exists {
				continue
			}

			// Si el valor es null, lo ignoramos (es opcional)
			if value == nil {
				continue
			}

			// Si el valor es un array, lo validamos como un array de referencias de tipo
			if array, ok := value.([]any); ok {
				for i, item := range array {
					if item == nil {
						continue
					}
					err := validateType(item, templates, refType, currentPath, fmt.Sprintf("%s[%d]", fieldPath, i), typeVisited)
					if err != nil {
						return err
					}
				}
				continue
			}

			err := validateType(value, templates, refType, currentPath, joinPath(currentPath, fieldPath), typeVisited)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// validateWildcardTypeRef handles type references with wildcards like "posts[*]->Post"
func validateWildcardTypeRef(data any, fieldPath, refType string, templates map[string][]string,
	basePath string, visited map[string]bool) error {
	if data == nil {
		return nil // Skip validation for nil fields
	}

	array, ok := data.([]any)
	if !ok {
		return fmt.Errorf("%s: expected array, got %T", fieldPath, data)
	}

	if len(array) == 0 {
		return nil
	}

	for i, item := range array {
		// Si el elemento es null, lo ignoramos
		if item == nil {
			continue
		}

		err := validateType(item, templates, refType, basePath, fmt.Sprintf("%s[%d]", fieldPath, i), visited)
		if err != nil {
			return err
		}
	}

	return nil
}

// validatePath checks if a path exists in the data
func validatePath(data any, path, basePath string) error {
	if data == nil {
		return fmt.Errorf("%s: required field is missing", joinPath(basePath, path))
	}

	obj, ok := data.(map[string]any)
	if !ok {
		if _, isArray := data.([]any); isArray {
			return fmt.Errorf("%s: expected object, got array", joinPath(basePath, path))
		}
		return fmt.Errorf("%s: expected object, got %T", joinPath(basePath, path), data)
	}

	parts := strings.Split(path, ".")
	for i, part := range parts {
		if strings.Contains(part, "[*]") {
			arrayField := strings.Split(part, "[*]")[0]
			value, exists := obj[arrayField]
			if !exists {
				return fmt.Errorf("%s: required field is missing", joinPath(basePath, strings.Join(parts[:i+1], ".")))
			}

			array, ok := value.([]any)
			if !ok {
				return fmt.Errorf("%s: expected array, got %T", joinPath(basePath, strings.Join(parts[:i+1], ".")), value)
			}

			if i < len(parts)-1 {
				for idx, elem := range array {
					elemObj, ok := elem.(map[string]any)
					if !ok {
						return fmt.Errorf("%s[%d]: expected object, got %T", joinPath(basePath, strings.Join(parts[:i+1], ".")), idx, elem)
					}

					if err := validatePath(elemObj, strings.Join(parts[i+1:], "."), joinPath(basePath, fmt.Sprintf("%s[%d]", arrayField, idx))); err != nil {
						return err
					}
				}
			}
			return nil
		}

		if i == len(parts)-1 {
			if _, exists := obj[part]; !exists {
				return fmt.Errorf("%s: required field is missing", joinPath(basePath, path))
			}
			return nil
		}

		value, exists := obj[part]
		if !exists {
			return fmt.Errorf("%s: required field is missing", joinPath(basePath, strings.Join(parts[:i+1], ".")))
		}

		nextObj, ok := value.(map[string]any)
		if !ok {
			if _, isArray := value.([]any); isArray {
				return fmt.Errorf("%s: expected object, got array", joinPath(basePath, strings.Join(parts[:i+1], ".")))
			}
			return fmt.Errorf("%s: expected object, got %T", joinPath(basePath, strings.Join(parts[:i+1], ".")), value)
		}
		obj = nextObj
	}

	return nil
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
