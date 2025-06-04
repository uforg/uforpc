package transpile

import (
	"fmt"

	"github.com/uforg/uforpc/urpc/internal/schema"
	"github.com/uforg/uforpc/urpc/internal/urpc/ast"
)

// ToURPC transpiles an UFO-RPC JSON schema to it's AST representation.
//
// The resulting AST Schema will not include any imports, extends, external
// docstrings, comments nor comment blocks.
//
// To get the string representation of the AST Schema, you can use the
// formatter package.
func ToURPC(jsonSchema schema.Schema) (ast.Schema, error) {
	result := ast.Schema{}

	// Add version declaration first
	if jsonSchema.Version > 0 {
		result.Children = append(result.Children, &ast.SchemaChild{
			Version: &ast.Version{
				Number: jsonSchema.Version,
			},
		})
	}

	// Process all nodes in the order they appear in the JSON schema
	for _, node := range jsonSchema.Nodes {
		switch n := node.(type) {
		case *schema.NodeDoc:
			// Convert standalone documentation to docstring
			result.Children = append(result.Children, &ast.SchemaChild{
				Docstring: &ast.Docstring{
					Value: n.Content,
				},
			})
		case *schema.NodeRule:
			ruleDecl, err := convertRuleToURPC(n)
			if err != nil {
				return ast.Schema{}, fmt.Errorf("error converting rule '%s': %w", n.Name, err)
			}
			result.Children = append(result.Children, &ast.SchemaChild{
				Rule: ruleDecl,
			})
		case *schema.NodeType:
			typeDecl, err := convertTypeToURPC(n)
			if err != nil {
				return ast.Schema{}, fmt.Errorf("error converting type '%s': %w", n.Name, err)
			}
			result.Children = append(result.Children, &ast.SchemaChild{
				Type: typeDecl,
			})
		case *schema.NodeProc:
			procDecl, err := convertProcToURPC(n)
			if err != nil {
				return ast.Schema{}, fmt.Errorf("error converting procedure '%s': %w", n.Name, err)
			}
			result.Children = append(result.Children, &ast.SchemaChild{
				Proc: procDecl,
			})
		}
	}

	return result, nil
}

// convertRuleToURPC converts a schema NodeRule to an AST RuleDecl
func convertRuleToURPC(rule *schema.NodeRule) (*ast.RuleDecl, error) {
	ruleDecl := &ast.RuleDecl{
		Name: rule.Name,
	}

	// Add docstring if available
	if rule.Doc != nil && *rule.Doc != "" {
		ruleDecl.Docstring = &ast.Docstring{
			Value: *rule.Doc,
		}
	}

	// Add deprecated if available
	if rule.Deprecated != nil {
		deprecated := &ast.Deprecated{}
		if *rule.Deprecated != "" {
			deprecated.Message = rule.Deprecated
		}
		ruleDecl.Deprecated = deprecated
	}

	// Add 'for' clause
	if rule.For != nil {
		ruleDecl.Children = append(ruleDecl.Children, &ast.RuleDeclChild{
			For: &ast.RuleDeclChildFor{
				Type:    rule.For.Type,
				IsArray: rule.For.IsArray,
			},
		})
	}

	// Add 'param' clause if available
	if rule.Param != nil {
		paramType, err := convertParamTypeToURPC(rule.Param.Type)
		if err != nil {
			return nil, fmt.Errorf("invalid parameter type: %w", err)
		}

		ruleDecl.Children = append(ruleDecl.Children, &ast.RuleDeclChild{
			Param: &ast.RuleDeclChildParam{
				Param:   paramType,
				IsArray: rule.Param.IsArray,
			},
		})
	}

	// Add 'error' clause if available
	if rule.Error != nil && *rule.Error != "" {
		ruleDecl.Children = append(ruleDecl.Children, &ast.RuleDeclChild{
			Error: &ast.RuleDeclChildError{
				Error: *rule.Error,
			},
		})
	}

	return ruleDecl, nil
}

// convertParamTypeToURPC converts a schema ParamPrimitiveType to a string
func convertParamTypeToURPC(paramType schema.ParamPrimitiveType) (string, error) {
	switch paramType {
	case schema.ParamPrimitiveTypeString:
		return "string", nil
	case schema.ParamPrimitiveTypeInt:
		return "int", nil
	case schema.ParamPrimitiveTypeFloat:
		return "float", nil
	case schema.ParamPrimitiveTypeBool:
		return "bool", nil
	default:
		return "", fmt.Errorf("invalid parameter type: %s", paramType.Value)
	}
}

// convertTypeToURPC converts a schema NodeType to an AST TypeDecl
func convertTypeToURPC(typeNode *schema.NodeType) (*ast.TypeDecl, error) {
	typeDecl := &ast.TypeDecl{
		Name: typeNode.Name,
	}

	// Add docstring if available
	if typeNode.Doc != nil && *typeNode.Doc != "" {
		typeDecl.Docstring = &ast.Docstring{
			Value: *typeNode.Doc,
		}
	}

	// Add deprecated if available
	if typeNode.Deprecated != nil {
		deprecated := &ast.Deprecated{}
		if *typeNode.Deprecated != "" {
			deprecated.Message = typeNode.Deprecated
		}
		typeDecl.Deprecated = deprecated
	}

	// Process fields
	for _, field := range typeNode.Fields {
		fieldNode, err := convertFieldToURPC(field)
		if err != nil {
			return nil, fmt.Errorf("error converting field '%s': %w", field.Name, err)
		}

		typeDecl.Children = append(typeDecl.Children, &ast.FieldOrComment{
			Field: fieldNode,
		})
	}

	return typeDecl, nil
}

// convertFieldToURPC converts a schema FieldDefinition to an AST Field
func convertFieldToURPC(fieldDef schema.FieldDefinition) (*ast.Field, error) {
	field := &ast.Field{
		Name:     fieldDef.Name,
		Optional: fieldDef.Optional,
	}

	// Process field type
	fieldType := ast.FieldType{
		IsArray: fieldDef.IsArray,
		Base:    &ast.FieldTypeBase{},
	}

	if fieldDef.IsNamed() {
		fieldType.Base.Named = fieldDef.TypeName
	}

	if fieldDef.IsInline() {
		object := &ast.FieldTypeObject{}

		// Process inline object fields
		for _, inlineField := range fieldDef.TypeInline.Fields {
			inlineFieldNode, err := convertFieldToURPC(inlineField)
			if err != nil {
				return nil, fmt.Errorf("error converting inline field '%s': %w", inlineField.Name, err)
			}

			object.Children = append(object.Children, &ast.FieldOrComment{
				Field: inlineFieldNode,
			})
		}

		fieldType.Base.Object = object
	}

	field.Type = fieldType

	// Process field rules
	if len(fieldDef.Rules) > 0 {
		for _, rule := range fieldDef.Rules {
			ruleNode, err := convertAppliedRuleToURPC(rule)
			if err != nil {
				return nil, fmt.Errorf("error converting rule '%s': %w", rule.Rule, err)
			}

			field.Children = append(field.Children, &ast.FieldChild{
				Rule: ruleNode,
			})
		}
	}

	return field, nil
}

// convertAppliedRuleToURPC converts a schema AppliedRule to an AST FieldRule
func convertAppliedRuleToURPC(rule schema.AppliedRule) (*ast.FieldRule, error) {
	fieldRule := &ast.FieldRule{
		Name: rule.Rule,
	}

	// Process rule parameters and error message if available
	if rule.Param != nil || (rule.Error != nil && *rule.Error != "") {
		body := &ast.FieldRuleBody{}

		// Process error message
		if rule.Error != nil && *rule.Error != "" {
			body.Error = rule.Error
		}

		// Process parameters
		if rule.Param != nil {
			if !rule.Param.IsArray {
				// Handle single parameter
				literal, err := convertParamValueToURPC(rule.Param.Type, rule.Param.Value)
				if err != nil {
					return nil, err
				}
				body.ParamSingle = &literal
			}

			if rule.Param.IsArray {
				// Handle array parameters
				switch rule.Param.Type.Value {
				case "string":
					body.ParamListString = rule.Param.ArrayValues
				case "int":
					body.ParamListInt = rule.Param.ArrayValues
				case "float":
					body.ParamListFloat = rule.Param.ArrayValues
				case "bool":
					body.ParamListBool = rule.Param.ArrayValues
				default:
					return nil, fmt.Errorf("unsupported array parameter type: %s", rule.Param.Type.Value)
				}
			}
		}

		fieldRule.Body = body
	}

	return fieldRule, nil
}

// convertParamValueToURPC converts a parameter value to an AST AnyLiteral
func convertParamValueToURPC(paramType schema.ParamPrimitiveType, value string) (ast.AnyLiteral, error) {
	var literal ast.AnyLiteral

	switch paramType.Value {
	case "string":
		literal.Str = &value
	case "int":
		literal.Int = &value
	case "float":
		literal.Float = &value
	case "bool":
		if value == "true" {
			t := "true"
			literal.True = &t
		}
		if value == "false" {
			f := "false"
			literal.False = &f
		}
	default:
		return ast.AnyLiteral{}, fmt.Errorf("unsupported parameter type: %s", paramType.Value)
	}

	return literal, nil
}

// convertProcToURPC converts a schema NodeProc to an AST ProcDecl
func convertProcToURPC(procNode *schema.NodeProc) (*ast.ProcDecl, error) {
	procDecl := &ast.ProcDecl{
		Name: procNode.Name,
	}

	// Add docstring if available
	if procNode.Doc != nil && *procNode.Doc != "" {
		procDecl.Docstring = &ast.Docstring{
			Value: *procNode.Doc,
		}
	}

	// Add deprecated if available
	if procNode.Deprecated != nil {
		deprecated := &ast.Deprecated{}
		if *procNode.Deprecated != "" {
			deprecated.Message = procNode.Deprecated
		}
		procDecl.Deprecated = deprecated
	}

	// Process input fields if any
	if len(procNode.Input) > 0 {
		inputChild := &ast.ProcOrStreamDeclChildInput{}

		for _, field := range procNode.Input {
			fieldNode, err := convertFieldToURPC(field)
			if err != nil {
				return nil, fmt.Errorf("error converting input field '%s': %w", field.Name, err)
			}

			inputChild.Children = append(inputChild.Children, &ast.FieldOrComment{
				Field: fieldNode,
			})
		}

		procDecl.Children = append(procDecl.Children, &ast.ProcOrStreamDeclChild{
			Input: inputChild,
		})
	}

	// Process output fields if any
	if len(procNode.Output) > 0 {
		outputChild := &ast.ProcOrStreamDeclChildOutput{}

		for _, field := range procNode.Output {
			fieldNode, err := convertFieldToURPC(field)
			if err != nil {
				return nil, fmt.Errorf("error converting output field '%s': %w", field.Name, err)
			}

			outputChild.Children = append(outputChild.Children, &ast.FieldOrComment{
				Field: fieldNode,
			})
		}

		procDecl.Children = append(procDecl.Children, &ast.ProcOrStreamDeclChild{
			Output: outputChild,
		})
	}

	// Process meta fields if any
	if len(procNode.Meta) > 0 {
		metaChild := &ast.ProcOrStreamDeclChildMeta{}

		// Process any remaining keys that weren't in the predefined order
		for _, metaKV := range procNode.Meta {
			key := metaKV.Key
			value := metaKV.Value

			literal, err := convertMetaValueToURPC(value)
			if err != nil {
				return nil, fmt.Errorf("error converting meta value for key '%s': %w", key, err)
			}

			metaChild.Children = append(metaChild.Children, &ast.ProcOrStreamDeclChildMetaChild{
				KV: &ast.ProcOrStreamDeclChildMetaKV{
					Key:   key,
					Value: literal,
				},
			})
		}

		procDecl.Children = append(procDecl.Children, &ast.ProcOrStreamDeclChild{
			Meta: metaChild,
		})
	}

	return procDecl, nil
}

// convertMetaValueToURPC converts a schema MetaValue to an AST AnyLiteral
func convertMetaValueToURPC(value schema.MetaValue) (ast.AnyLiteral, error) {
	var literal ast.AnyLiteral

	if value.StringVal != nil {
		quoted := *value.StringVal
		literal.Str = &quoted
	} else if value.IntVal != nil {
		intStr := fmt.Sprintf("%d", *value.IntVal)
		literal.Int = &intStr
	} else if value.FloatVal != nil {
		floatStr := fmt.Sprintf("%g", *value.FloatVal)
		literal.Float = &floatStr
	} else if value.BoolVal != nil {
		if *value.BoolVal {
			t := "true"
			literal.True = &t
		} else {
			f := "false"
			literal.False = &f
		}
	} else {
		return ast.AnyLiteral{}, fmt.Errorf("empty meta value")
	}

	return literal, nil
}
