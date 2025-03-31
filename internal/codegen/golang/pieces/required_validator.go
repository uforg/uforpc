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

// validateRequiredJSONPaths validates JSON data against a set of required paths.
// It supports recursive type references with the syntax field->Type.
func validateRequiredJSONPaths(data []byte, templates map[string][]string, rootType string) error {
	var jsonData any
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}

	// Helper function to check if a value is an array
	isArray := func(v any) bool {
		_, ok := v.([]any)
		return ok
	}

	// Helper function to join path segments
	joinPath := func(base, path string) string {
		if base == "" || path == "" {
			return base + path
		}
		return base + "." + path
	}

	// Helper function to validate a path in an object
	var validatePath func(obj map[string]any, path, basePath string) error

	// Main validation function
	var validateType func(data any, typeName, currentPath string, visited map[string]bool) error

	validatePath = func(obj map[string]any, path, basePath string) error {
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
				if isArray(value) {
					return fmt.Errorf("%s: expected object, got array", joinPath(basePath, strings.Join(parts[:i+1], ".")))
				}
				return fmt.Errorf("%s: expected object, got %T", joinPath(basePath, strings.Join(parts[:i+1], ".")), value)
			}
			obj = nextObj
		}
		return nil
	}

	validateType = func(data any, typeName, currentPath string, visited map[string]bool) error {
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

		// Check if we've already visited this type at this path
		visitKey := fmt.Sprintf("%s:%s", typeName, currentPath)
		if visited[visitKey] {
			// Only validate required fields if already visited
			for _, path := range paths {
				if !strings.Contains(path, "->") {
					if err := validatePath(obj, path, currentPath); err != nil {
						return err
					}
				}
			}
			return nil
		}
		visited[visitKey] = true

		// Validate all fields
		for _, path := range paths {
			if strings.Contains(path, "->") {
				parts := strings.Split(path, "->")
				if len(parts) != 2 {
					return fmt.Errorf("invalid type reference: %s", path)
				}

				fieldPath, refType := parts[0], parts[1]
				value, exists := obj[strings.Split(fieldPath, "[*]")[0]]
				if !exists || value == nil {
					continue
				}

				if strings.Contains(fieldPath, "[*]") || isArray(value) {
					array, ok := value.([]any)
					if !ok {
						return fmt.Errorf("%s: expected array, got %T", fieldPath, value)
					}
					for i, item := range array {
						if item != nil {
							if err := validateType(item, refType, fmt.Sprintf("%s[%d]", strings.Split(fieldPath, "[*]")[0], i), visited); err != nil {
								return err
							}
						}
					}
					continue
				}

				if err := validateType(value, refType, joinPath(currentPath, fieldPath), visited); err != nil {
					return err
				}
			} else if err := validatePath(obj, path, currentPath); err != nil {
				return err
			}
		}
		return nil
	}

	return validateType(jsonData, rootType, "", make(map[string]bool))
}
