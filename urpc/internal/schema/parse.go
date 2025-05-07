package schema

import (
	"encoding/json"
	"fmt"
)

// ParseSchema parses and validates a JSON schema string into a Schema struct
func ParseSchema(schemaStr string) (Schema, error) {
	if err := validateSchema(schemaStr); err != nil {
		return Schema{}, fmt.Errorf("error validating against JSON schema: %w", err)
	}

	var schema Schema
	if err := json.Unmarshal([]byte(schemaStr), &schema); err != nil {
		return Schema{}, fmt.Errorf("error decoding schema: %w", err)
	}

	if err := validateStructure(schema); err != nil {
		return Schema{}, fmt.Errorf("error validating schema structure: %w", err)
	}

	return schema, nil
}
