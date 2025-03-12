package schema

import (
	"encoding/json"
	"fmt"
)

// ParseSchema validates and converts a JSON schema string to its Go representation
func ParseSchema(schemaStr string) (*Schema, error) {
	if err := ValidateSchema(schemaStr); err != nil {
		return nil, fmt.Errorf("invalid schema: %w", err)
	}

	var schema Schema
	if err := json.Unmarshal([]byte(schemaStr), &schema); err != nil {
		return nil, fmt.Errorf("error decoding schema: %w", err)
	}

	if err := processSchema(&schema); err != nil {
		return nil, err
	}

	return &schema, nil
}

// processSchema processes and validates the Schema structure after decoding
func processSchema(schema *Schema) error {
	if schema.Version != 1 {
		return fmt.Errorf("unsupported schema version: %d", schema.Version)
	}

	for typeName, field := range schema.Types {
		if err := processField(typeName, &field, schema); err != nil {
			return err
		}
		schema.Types[typeName] = field
	}

	for procName, proc := range schema.Procedures {
		procType, ok := proc.GetProcedureType()
		if !ok {
			return fmt.Errorf("invalid procedure type for %s: %s", procName, proc.Type)
		}

		if proc.Input != nil {
			if err := processField(procName+".input", proc.Input, schema); err != nil {
				return err
			}
		}

		if proc.Output != nil {
			if err := processField(procName+".output", proc.Output, schema); err != nil {
				return err
			}
		}

		proc.Type = procType.Value
		schema.Procedures[procName] = proc
	}

	return nil
}

// processField processes and validates a field and its subfields recursively
func processField(path string, field *Field, schema *Schema) error {
	if field.IsBuiltInType() {
		fieldType, _ := field.GetFieldType()

		switch fieldType.Value {
		case FieldTypeObject.Value:
			if field.Fields == nil || len(field.Fields) == 0 {
				return fmt.Errorf("object field %s must have defined fields", path)
			}

			for fieldName, subField := range field.Fields {
				if err := processField(path+"."+fieldName, &subField, schema); err != nil {
					return err
				}
				field.Fields[fieldName] = subField
			}

		case FieldTypeArray.Value:
			if field.ArrayType == nil {
				return fmt.Errorf("array field %s must have a defined array type", path)
			}

			if err := processField(path+".arrayType", field.ArrayType, schema); err != nil {
				return err
			}
		}

		if field.Rules == nil {
			switch fieldType.Value {
			case FieldTypeString.Value:
				field.Rules = &StringRules{}
			case FieldTypeInt.Value:
				field.Rules = &IntRules{}
			case FieldTypeFloat.Value:
				field.Rules = &FloatRules{}
			case FieldTypeBoolean.Value:
				field.Rules = &BooleanRules{}
			case FieldTypeObject.Value:
				field.Rules = &ObjectRules{}
			case FieldTypeArray.Value:
				field.Rules = &ArrayRules{}
			}
		} else {
			convertRules(field)
		}
	} else {
		if _, exists := schema.Types[field.Type]; !exists {
			return fmt.Errorf("undefined custom type: %s in field %s", field.Type, path)
		}

		if field.Rules == nil {
			field.Rules = &CustomTypeRules{}
		} else {
			if _, ok := field.Rules.(*CustomTypeRules); !ok {
				if mapRules, ok := field.Rules.(map[string]any); ok {
					customRules := &CustomTypeRules{}
					if optional, ok := mapRules["optional"].(bool); ok {
						customRules.Optional = optional
					}
					field.Rules = customRules
				}
			}
		}
	}

	return nil
}

// convertRules converts a field's rules to the correct type
func convertRules(field *Field) {
	fieldType, ok := field.GetFieldType()
	if !ok {
		return
	}

	mapRules, ok := field.Rules.(map[string]any)
	if !ok {
		return
	}

	switch fieldType.Value {
	case FieldTypeString.Value:
		stringRules := &StringRules{}

		if optional, ok := mapRules["optional"].(bool); ok {
			stringRules.Optional = optional
		}

		if equals, ok := mapRules["equals"].(map[string]any); ok {
			stringRules.Equals = processStringValueRule(equals)
		}

		if contains, ok := mapRules["contains"].(map[string]any); ok {
			stringRules.Contains = processStringValueRule(contains)
		}

		if regex, ok := mapRules["regex"].(map[string]any); ok {
			stringRules.Regex = processStringValueRule(regex)
		}

		if minLen, ok := mapRules["minLen"].(map[string]any); ok {
			stringRules.MinLen = processIntValueRule(minLen)
		}

		if maxLen, ok := mapRules["maxLen"].(map[string]any); ok {
			stringRules.MaxLen = processIntValueRule(maxLen)
		}

		if enum, ok := mapRules["enum"].(map[string]any); ok {
			stringRules.Enum = processStringArrayRule(enum)
		}

		if email, ok := mapRules["email"].(map[string]any); ok {
			stringRules.Email = processSimpleRule(email)
		}

		if iso8601, ok := mapRules["iso8601"].(map[string]any); ok {
			stringRules.ISO8601 = processSimpleRule(iso8601)
		}

		if uuid, ok := mapRules["uuid"].(map[string]any); ok {
			stringRules.UUID = processSimpleRule(uuid)
		}

		if json, ok := mapRules["json"].(map[string]any); ok {
			stringRules.JSON = processSimpleRule(json)
		}

		if lowercase, ok := mapRules["lowercase"].(map[string]any); ok {
			stringRules.Lowercase = processSimpleRule(lowercase)
		}

		if uppercase, ok := mapRules["uppercase"].(map[string]any); ok {
			stringRules.Uppercase = processSimpleRule(uppercase)
		}

		field.Rules = stringRules

	case FieldTypeInt.Value:
		intRules := &IntRules{}

		if optional, ok := mapRules["optional"].(bool); ok {
			intRules.Optional = optional
		}

		if equals, ok := mapRules["equals"].(map[string]any); ok {
			intRules.Equals = processIntValueRule(equals)
		}

		if min, ok := mapRules["min"].(map[string]any); ok {
			intRules.Min = processIntValueRule(min)
		}

		if max, ok := mapRules["max"].(map[string]any); ok {
			intRules.Max = processIntValueRule(max)
		}

		if enum, ok := mapRules["enum"].(map[string]any); ok {
			intRules.Enum = processIntArrayRule(enum)
		}

		field.Rules = intRules

	case FieldTypeFloat.Value:
		floatRules := &FloatRules{}

		if optional, ok := mapRules["optional"].(bool); ok {
			floatRules.Optional = optional
		}

		if equals, ok := mapRules["equals"].(map[string]any); ok {
			floatRules.Equals = processNumberValueRule(equals)
		}

		if min, ok := mapRules["min"].(map[string]any); ok {
			floatRules.Min = processNumberValueRule(min)
		}

		if max, ok := mapRules["max"].(map[string]any); ok {
			floatRules.Max = processNumberValueRule(max)
		}

		if enum, ok := mapRules["enum"].(map[string]any); ok {
			floatRules.Enum = processNumberArrayRule(enum)
		}

		field.Rules = floatRules

	case FieldTypeBoolean.Value:
		boolRules := &BooleanRules{}

		if optional, ok := mapRules["optional"].(bool); ok {
			boolRules.Optional = optional
		}

		if equals, ok := mapRules["equals"].(map[string]any); ok {
			boolRules.Equals = processBooleanValueRule(equals)
		}

		field.Rules = boolRules

	case FieldTypeObject.Value:
		objectRules := &ObjectRules{}

		if optional, ok := mapRules["optional"].(bool); ok {
			objectRules.Optional = optional
		}

		field.Rules = objectRules

	case FieldTypeArray.Value:
		arrayRules := &ArrayRules{}

		if optional, ok := mapRules["optional"].(bool); ok {
			arrayRules.Optional = optional
		}

		if minLen, ok := mapRules["minLen"].(map[string]any); ok {
			arrayRules.MinLen = processIntValueRule(minLen)
		}

		if maxLen, ok := mapRules["maxLen"].(map[string]any); ok {
			arrayRules.MaxLen = processIntValueRule(maxLen)
		}

		field.Rules = arrayRules
	}
}

// Helper functions for processing rules

func processSimpleRule(rule map[string]any) RuleSimple {
	result := RuleSimple{}

	if errorMsg, ok := rule["errorMessage"].(string); ok {
		result.ErrorMessage = errorMsg
	}

	return result
}

func processStringValueRule(rule map[string]any) RuleWithStringValue {
	result := RuleWithStringValue{}

	if value, ok := rule["value"].(string); ok {
		result.Value = value
	}

	if errorMsg, ok := rule["errorMessage"].(string); ok {
		result.ErrorMessage = errorMsg
	}

	return result
}

func processIntValueRule(rule map[string]any) RuleWithIntValue {
	result := RuleWithIntValue{}

	if value, ok := rule["value"].(float64); ok {
		result.Value = int(value)
	} else if value, ok := rule["value"].(int); ok {
		result.Value = value
	}

	if errorMsg, ok := rule["errorMessage"].(string); ok {
		result.ErrorMessage = errorMsg
	}

	return result
}

func processNumberValueRule(rule map[string]any) RuleWithNumberValue {
	result := RuleWithNumberValue{}

	if value, ok := rule["value"].(float64); ok {
		result.Value = value
	}

	if errorMsg, ok := rule["errorMessage"].(string); ok {
		result.ErrorMessage = errorMsg
	}

	return result
}

func processBooleanValueRule(rule map[string]any) RuleWithBooleanValue {
	result := RuleWithBooleanValue{}

	if value, ok := rule["value"].(bool); ok {
		result.Value = value
	}

	if errorMsg, ok := rule["errorMessage"].(string); ok {
		result.ErrorMessage = errorMsg
	}

	return result
}

func processStringArrayRule(rule map[string]any) RuleWithStringArray {
	result := RuleWithStringArray{
		Values: []string{},
	}

	if values, ok := rule["values"].([]any); ok {
		for _, v := range values {
			if strValue, ok := v.(string); ok {
				result.Values = append(result.Values, strValue)
			}
		}
	}

	if errorMsg, ok := rule["errorMessage"].(string); ok {
		result.ErrorMessage = errorMsg
	}

	return result
}

func processIntArrayRule(rule map[string]any) RuleWithIntArray {
	result := RuleWithIntArray{
		Values: []int{},
	}

	if values, ok := rule["values"].([]any); ok {
		for _, v := range values {
			if floatValue, ok := v.(float64); ok {
				result.Values = append(result.Values, int(floatValue))
			} else if intValue, ok := v.(int); ok {
				result.Values = append(result.Values, intValue)
			}
		}
	}

	if errorMsg, ok := rule["errorMessage"].(string); ok {
		result.ErrorMessage = errorMsg
	}

	return result
}

func processNumberArrayRule(rule map[string]any) RuleWithNumberArray {
	result := RuleWithNumberArray{
		Values: []float64{},
	}

	if values, ok := rule["values"].([]any); ok {
		for _, v := range values {
			if floatValue, ok := v.(float64); ok {
				result.Values = append(result.Values, floatValue)
			}
		}
	}

	if errorMsg, ok := rule["errorMessage"].(string); ok {
		result.ErrorMessage = errorMsg
	}

	return result
}

// UnmarshalJSON implements the json.Unmarshaler interface for Field
func (f *Field) UnmarshalJSON(data []byte) error {
	type FieldTemp Field
	var temp FieldTemp

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*f = Field(temp)

	if f.Type == FieldTypeObject.Value && f.Fields == nil {
		f.Fields = make(map[string]Field)
	}

	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for Schema
func (s *Schema) UnmarshalJSON(data []byte) error {
	type SchemaTemp Schema
	var temp SchemaTemp

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*s = Schema(temp)

	if s.Types == nil {
		s.Types = make(map[string]Field)
	}
	if s.Procedures == nil {
		s.Procedures = make(map[string]Procedure)
	}

	return nil
}
