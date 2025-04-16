package transpile

import (
	"fmt"
	"strconv"

	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/urpc/ast"
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

	// Add version declaration
	if jsonSchema.Version > 0 {
		result.Children = append(result.Children, &ast.SchemaChild{
			Version: &ast.Version{
				Number: jsonSchema.Version,
			},
		})
	}

	// Process all nodes
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
		default:
			return ast.Schema{}, fmt.Errorf("unknown node type: %T", node)
		}
	}

	return result, nil
}

// convertRuleToURPC converts a schema NodeRule to an AST RuleDecl
func convertRuleToURPC(rule *schema.NodeRule) (*ast.RuleDecl, error) {
	ruleDecl := &ast.RuleDecl{
		Name:     rule.Name,
		Children: []*ast.RuleDeclChild{},
	}

	// Add docstring if available
	if rule.Doc != nil && *rule.Doc != "" {
		ruleDecl.Docstring = &ast.Docstring{
			Value: *rule.Doc,
		}
	}

	// Add 'for' clause
	if rule.For != "" {
		ruleDecl.Children = append(ruleDecl.Children, &ast.RuleDeclChild{
			For: &ast.RuleDeclChildFor{
				For: rule.For,
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
				Error: strconv.Quote(*rule.Error),
			},
		})
	}

	return ruleDecl, nil
}

// convertTypeToURPC converts a schema NodeType to an AST TypeDecl
func convertTypeToURPC(typeNode *schema.NodeType) (*ast.TypeDecl, error) {
	typeDecl := &ast.TypeDecl{
		Name:     typeNode.Name,
		Children: []*ast.FieldOrComment{},
	}

	// Add docstring if available
	if typeNode.Doc != nil && *typeNode.Doc != "" {
		typeDecl.Docstring = &ast.Docstring{
			Value: *typeNode.Doc,
		}
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

// convertProcToURPC converts a schema NodeProc to an AST ProcDecl
func convertProcToURPC(procNode *schema.NodeProc) (*ast.ProcDecl, error) {
	procDecl := &ast.ProcDecl{
		Name:     procNode.Name,
		Children: []*ast.ProcDeclChild{},
	}

	// Add docstring if available
	if procNode.Doc != nil && *procNode.Doc != "" {
		procDecl.Docstring = &ast.Docstring{
			Value: *procNode.Doc,
		}
	}

	// Process input fields if any
	if len(procNode.Input) > 0 {
		inputChild := &ast.ProcDeclChildInput{
			Children: []*ast.FieldOrComment{},
		}

		for _, field := range procNode.Input {
			fieldNode, err := convertFieldToURPC(field)
			if err != nil {
				return nil, fmt.Errorf("error converting input field '%s': %w", field.Name, err)
			}

			inputChild.Children = append(inputChild.Children, &ast.FieldOrComment{
				Field: fieldNode,
			})
		}

		procDecl.Children = append(procDecl.Children, &ast.ProcDeclChild{
			Input: inputChild,
		})
	}

	// Process output fields if any
	if len(procNode.Output) > 0 {
		outputChild := &ast.ProcDeclChildOutput{
			Children: []*ast.FieldOrComment{},
		}

		for _, field := range procNode.Output {
			fieldNode, err := convertFieldToURPC(field)
			if err != nil {
				return nil, fmt.Errorf("error converting output field '%s': %w", field.Name, err)
			}

			outputChild.Children = append(outputChild.Children, &ast.FieldOrComment{
				Field: fieldNode,
			})
		}

		procDecl.Children = append(procDecl.Children, &ast.ProcDeclChild{
			Output: outputChild,
		})
	}

	// Process meta fields if any
	if len(procNode.Meta) > 0 {
		metaChild := &ast.ProcDeclChildMeta{
			Children: []*ast.ProcDeclChildMetaChild{},
		}

		for key, value := range procNode.Meta {
			literal, err := convertMetaValueToURPC(value)
			if err != nil {
				return nil, fmt.Errorf("error converting meta value for key '%s': %w", key, err)
			}

			metaChild.Children = append(metaChild.Children, &ast.ProcDeclChildMetaChild{
				KV: &ast.ProcDeclChildMetaKV{
					Key:   key,
					Value: literal,
				},
			})
		}

		procDecl.Children = append(procDecl.Children, &ast.ProcDeclChild{
			Meta: metaChild,
		})
	}

	return procDecl, nil
}

// convertFieldToURPC converts a schema FieldDefinition to an AST Field
func convertFieldToURPC(fieldDef schema.FieldDefinition) (*ast.Field, error) {
	// Always initialize Children as an empty slice to avoid nil vs empty slice comparison issues
	field := &ast.Field{
		Name:     fieldDef.Name,
		Optional: fieldDef.Optional,
		Children: []*ast.FieldChild{},
	}

	// Process field type
	fieldType := ast.FieldType{
		Depth: ast.FieldTypeDepth(fieldDef.Depth),
		Base:  &ast.FieldTypeBase{},
	}

	if fieldDef.IsNamed() {
		fieldType.Base.Named = fieldDef.TypeName
	} else if fieldDef.IsInline() {
		object := &ast.FieldTypeObject{
			Children: []*ast.FieldOrComment{},
		}

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
	for _, rule := range fieldDef.Rules {
		ruleNode, err := convertAppliedRuleToURPC(rule)
		if err != nil {
			return nil, fmt.Errorf("error converting rule '%s': %w", rule.Rule, err)
		}

		field.Children = append(field.Children, &ast.FieldChild{
			Rule: ruleNode,
		})
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
			errorMsg := strconv.Quote(*rule.Error)
			body.Error = &errorMsg
		}

		// Process parameters
		if rule.Param != nil {
			if rule.Param.IsArray {
				// Handle array parameters
				switch rule.Param.Type.Value {
				case "string":
					body.ParamListString = make([]string, len(rule.Param.ArrayValues))
					for i, val := range rule.Param.ArrayValues {
						body.ParamListString[i] = strconv.Quote(val)
					}
				case "int":
					body.ParamListInt = rule.Param.ArrayValues
				case "float":
					body.ParamListFloat = rule.Param.ArrayValues
				case "boolean":
					body.ParamListBoolean = rule.Param.ArrayValues
				default:
					return nil, fmt.Errorf("unsupported array parameter type: %s", rule.Param.Type.Value)
				}
			} else {
				// Handle single parameter
				literal, err := convertParamValueToURPC(rule.Param.Type, rule.Param.Value)
				if err != nil {
					return nil, err
				}
				body.ParamSingle = &literal
			}
		}

		fieldRule.Body = body
	}

	return fieldRule, nil
}

// convertMetaValueToURPC converts a schema MetaValue to an AST AnyLiteral
func convertMetaValueToURPC(value schema.MetaValue) (ast.AnyLiteral, error) {
	var literal ast.AnyLiteral

	if value.StringVal != nil {
		quoted := strconv.Quote(*value.StringVal)
		literal.Str = &quoted
	} else if value.IntVal != nil {
		intStr := fmt.Sprintf("%d", *value.IntVal)
		literal.Int = &intStr
	} else if value.FloatVal != nil {
		floatStr := fmt.Sprintf("%g", *value.FloatVal)
		literal.Float = &floatStr
	} else if value.BoolVal != nil {
		if *value.BoolVal {
			true := "true"
			literal.True = &true
		} else {
			false := "false"
			literal.False = &false
		}
	} else {
		return ast.AnyLiteral{}, fmt.Errorf("empty meta value")
	}

	return literal, nil
}

// convertParamValueToURPC converts a parameter value to an AST AnyLiteral
func convertParamValueToURPC(paramType schema.ParamPrimitiveType, value string) (ast.AnyLiteral, error) {
	var literal ast.AnyLiteral

	switch paramType.Value {
	case "string":
		quoted := strconv.Quote(value)
		literal.Str = &quoted
	case "int":
		literal.Int = &value
	case "float":
		literal.Float = &value
	case "boolean":
		if value == "true" {
			true := "true"
			literal.True = &true
		} else if value == "false" {
			false := "false"
			literal.False = &false
		} else {
			return ast.AnyLiteral{}, fmt.Errorf("invalid boolean value: %s", value)
		}
	default:
		return ast.AnyLiteral{}, fmt.Errorf("unsupported parameter type: %s", paramType.Value)
	}

	return literal, nil
}

// convertParamTypeToURPC converts a schema ParamPrimitiveType to a string
func convertParamTypeToURPC(paramType schema.ParamPrimitiveType) (string, error) {
	switch paramType.Value {
	case "string":
		return "string", nil
	case "int":
		return "int", nil
	case "float":
		return "float", nil
	case "boolean":
		return "boolean", nil
	default:
		return "", fmt.Errorf("invalid parameter type: %s", paramType.Value)
	}
}
