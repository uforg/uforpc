package transpile

import (
	"fmt"
	"strings"

	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/urpc/ast"
)

// ToJSON transpiles an UFO-RPC AST schema to it's JSON representation.
//
// The imports, extends and external docstrings of the AST Schema are expected
// to be already resolved.
//
// If there are any unresolved imports or extends, the transpiler
// will ignore them.
//
// If there are any unresolved external docstrings, the transpiler will
// treat them literally as strings.
//
// All comments and comment blocks will be ignored.
//
// To get the string representation of the JSON Schema, you can use the
// json.Marshal function.
func ToJSON(astSchema ast.Schema) (schema.Schema, error) {
	result := schema.Schema{
		Version: 1, // Default version is 1
		Nodes:   []schema.Node{},
	}

	// Get the version from the schema if available
	versions := astSchema.GetVersions()
	if len(versions) > 0 {
		result.Version = versions[0].Number
	}

	// Process standalone docstrings
	for _, docstring := range astSchema.GetDocstrings() {
		docNode := &schema.NodeDoc{
			Kind:    "doc",
			Content: strings.TrimSpace(docstring.Value),
		}
		result.Nodes = append(result.Nodes, docNode)
	}

	// Process rules
	for _, rule := range astSchema.GetRules() {
		ruleNode, err := convertRuleToJSON(rule)
		if err != nil {
			return schema.Schema{}, fmt.Errorf("error converting rule '%s': %w", rule.Name, err)
		}
		result.Nodes = append(result.Nodes, ruleNode)
	}

	// Process types
	for _, typeDecl := range astSchema.GetTypes() {
		typeNode, err := convertTypeToJSON(typeDecl)
		if err != nil {
			return schema.Schema{}, fmt.Errorf("error converting type '%s': %w", typeDecl.Name, err)
		}
		result.Nodes = append(result.Nodes, typeNode)
	}

	// Process procedures
	for _, procDecl := range astSchema.GetProcs() {
		procNode, err := convertProcToJSON(procDecl)
		if err != nil {
			return schema.Schema{}, fmt.Errorf("error converting procedure '%s': %w", procDecl.Name, err)
		}
		result.Nodes = append(result.Nodes, procNode)
	}

	return result, nil
}

// convertRuleToJSON converts an AST RuleDecl to a schema NodeRule
func convertRuleToJSON(rule *ast.RuleDecl) (*schema.NodeRule, error) {
	ruleNode := &schema.NodeRule{
		Kind: "rule",
		Name: rule.Name,
	}

	// Add docstring if available
	if rule.Docstring != nil {
		docValue := rule.Docstring.Value
		ruleNode.Doc = &docValue
	}

	// Process rule children
	for _, child := range rule.Children {
		if child.For != nil {
			ruleNode.For = child.For.For
		}
		if child.Param != nil {
			paramType, err := convertParamTypeToJSON(child.Param.Param)
			if err != nil {
				return nil, fmt.Errorf("invalid parameter type: %w", err)
			}

			ruleNode.Param = &schema.ParamDefinition{
				Type:    paramType,
				IsArray: child.Param.IsArray,
			}
		}
		if child.Error != nil {
			errorMsg := strings.Trim(child.Error.Error, `"`)
			ruleNode.Error = &errorMsg
		}
	}

	return ruleNode, nil
}

// convertTypeToJSON converts an AST TypeDecl to a schema NodeType
func convertTypeToJSON(typeDecl *ast.TypeDecl) (*schema.NodeType, error) {
	typeNode := &schema.NodeType{
		Kind:   "type",
		Name:   typeDecl.Name,
		Fields: []schema.FieldDefinition{},
	}

	// Add docstring if available
	if typeDecl.Docstring != nil {
		docValue := strings.TrimSpace(typeDecl.Docstring.Value)
		typeNode.Doc = &docValue
	}

	// Process fields
	for _, child := range typeDecl.Children {
		if child.Field != nil {
			fieldDef, err := convertFieldToJSON(child.Field)
			if err != nil {
				return nil, fmt.Errorf("error converting field '%s': %w", child.Field.Name, err)
			}
			typeNode.Fields = append(typeNode.Fields, fieldDef)
		}
	}

	return typeNode, nil
}

// convertProcToJSON converts an AST ProcDecl to a schema NodeProc
func convertProcToJSON(procDecl *ast.ProcDecl) (*schema.NodeProc, error) {
	procNode := &schema.NodeProc{
		Kind:   "proc",
		Name:   procDecl.Name,
		Input:  []schema.FieldDefinition{},
		Output: []schema.FieldDefinition{},
		Meta:   make(map[string]schema.MetaValue),
	}

	// Add docstring if available
	if procDecl.Docstring != nil {
		docValue := strings.TrimSpace(procDecl.Docstring.Value)
		procNode.Doc = &docValue
	}

	// Process procedure children
	for _, child := range procDecl.Children {
		if child.Input != nil {
			for _, fieldOrComment := range child.Input.Children {
				if fieldOrComment.Field != nil {
					fieldDef, err := convertFieldToJSON(fieldOrComment.Field)
					if err != nil {
						return nil, fmt.Errorf("error converting input field '%s': %w", fieldOrComment.Field.Name, err)
					}
					procNode.Input = append(procNode.Input, fieldDef)
				}
			}
		} else if child.Output != nil {
			for _, fieldOrComment := range child.Output.Children {
				if fieldOrComment.Field != nil {
					fieldDef, err := convertFieldToJSON(fieldOrComment.Field)
					if err != nil {
						return nil, fmt.Errorf("error converting output field '%s': %w", fieldOrComment.Field.Name, err)
					}
					procNode.Output = append(procNode.Output, fieldDef)
				}
			}
		} else if child.Meta != nil {
			for _, metaChild := range child.Meta.Children {
				if metaChild.KV != nil {
					metaValue, err := convertMetaValueToJSON(metaChild.KV.Value)
					if err != nil {
						return nil, fmt.Errorf("error converting meta value for key '%s': %w", metaChild.KV.Key, err)
					}
					procNode.Meta[metaChild.KV.Key] = metaValue
				}
			}
		}
	}

	return procNode, nil
}

// convertFieldToJSON converts an AST Field to a schema FieldDefinition
func convertFieldToJSON(field *ast.Field) (schema.FieldDefinition, error) {
	fieldDef := schema.FieldDefinition{
		Name:     field.Name,
		Optional: field.Optional,
		Depth:    int(field.Type.Depth),
		Rules:    []schema.AppliedRule{},
	}

	// Process field type
	if field.Type.Base.Named != nil {
		typeName := *field.Type.Base.Named
		fieldDef.TypeName = &typeName
	} else if field.Type.Base.Object != nil {
		inlineType := &schema.InlineTypeDefinition{
			Fields: []schema.FieldDefinition{},
		}

		// Process inline object fields
		for _, child := range field.Type.Base.Object.Children {
			if child.Field != nil {
				inlineField, err := convertFieldToJSON(child.Field)
				if err != nil {
					return schema.FieldDefinition{}, fmt.Errorf("error converting inline field '%s': %w", child.Field.Name, err)
				}
				inlineType.Fields = append(inlineType.Fields, inlineField)
			}
		}

		fieldDef.TypeInline = inlineType
	}

	// Process field rules
	for _, child := range field.Children {
		if child.Rule != nil {
			rule, err := convertFieldRuleToJSON(child.Rule)
			if err != nil {
				return schema.FieldDefinition{}, fmt.Errorf("error converting rule '%s': %w", child.Rule.Name, err)
			}
			fieldDef.Rules = append(fieldDef.Rules, rule)
		}
	}

	return fieldDef, nil
}

// convertFieldRuleToJSON converts an AST FieldRule to a schema AppliedRule
func convertFieldRuleToJSON(fieldRule *ast.FieldRule) (schema.AppliedRule, error) {
	rule := schema.AppliedRule{
		Rule: fieldRule.Name,
	}

	// Process rule body if available
	if fieldRule.Body != nil {
		// Process error message if available
		if fieldRule.Body.Error != nil {
			errorMsg := strings.Trim(*fieldRule.Body.Error, `"`)
			rule.Error = &errorMsg
		}

		// Process parameters
		if fieldRule.Body.ParamSingle != nil ||
			len(fieldRule.Body.ParamListString) > 0 ||
			len(fieldRule.Body.ParamListInt) > 0 ||
			len(fieldRule.Body.ParamListFloat) > 0 ||
			len(fieldRule.Body.ParamListBoolean) > 0 {

			param := &schema.AppliedParam{}

			// Handle single parameter
			if fieldRule.Body.ParamSingle != nil {
				paramType, paramValue, err := extractParamValue(*fieldRule.Body.ParamSingle)
				if err != nil {
					return schema.AppliedRule{}, err
				}
				param.Type = paramType
				param.IsArray = false
				param.Value = paramValue
			}

			// Handle array parameters
			if len(fieldRule.Body.ParamListString) > 0 {
				param.Type = schema.ParamPrimitiveTypeString
				param.IsArray = true
				param.ArrayValues = make([]string, len(fieldRule.Body.ParamListString))
				for i, val := range fieldRule.Body.ParamListString {
					param.ArrayValues[i] = strings.Trim(val, `"`)
				}
			} else if len(fieldRule.Body.ParamListInt) > 0 {
				param.Type = schema.ParamPrimitiveTypeInt
				param.IsArray = true
				param.ArrayValues = fieldRule.Body.ParamListInt
			} else if len(fieldRule.Body.ParamListFloat) > 0 {
				param.Type = schema.ParamPrimitiveTypeFloat
				param.IsArray = true
				param.ArrayValues = fieldRule.Body.ParamListFloat
			} else if len(fieldRule.Body.ParamListBoolean) > 0 {
				param.Type = schema.ParamPrimitiveTypeBoolean
				param.IsArray = true
				param.ArrayValues = fieldRule.Body.ParamListBoolean
			}

			rule.Param = param
		}
	}

	return rule, nil
}

// convertMetaValueToJSON converts an AST AnyLiteral to a schema MetaValue
func convertMetaValueToJSON(literal ast.AnyLiteral) (schema.MetaValue, error) {
	var metaValue schema.MetaValue

	if literal.Str != nil {
		strValue := strings.Trim(*literal.Str, `"`)
		metaValue.StringVal = &strValue
	} else if literal.Int != nil {
		var intValue int64
		if _, err := fmt.Sscanf(*literal.Int, "%d", &intValue); err != nil {
			return schema.MetaValue{}, fmt.Errorf("invalid integer value: %s", *literal.Int)
		}
		metaValue.IntVal = &intValue
	} else if literal.Float != nil {
		var floatValue float64
		if _, err := fmt.Sscanf(*literal.Float, "%f", &floatValue); err != nil {
			return schema.MetaValue{}, fmt.Errorf("invalid float value: %s", *literal.Float)
		}
		metaValue.FloatVal = &floatValue
	} else if literal.True != nil {
		boolValue := true
		metaValue.BoolVal = &boolValue
	} else if literal.False != nil {
		boolValue := false
		metaValue.BoolVal = &boolValue
	} else {
		return schema.MetaValue{}, fmt.Errorf("empty literal value")
	}

	return metaValue, nil
}

// extractParamValue extracts the parameter type and value from an AnyLiteral
func extractParamValue(literal ast.AnyLiteral) (schema.ParamPrimitiveType, string, error) {
	if literal.Str != nil {
		return schema.ParamPrimitiveTypeString, strings.Trim(*literal.Str, `"`), nil
	} else if literal.Int != nil {
		return schema.ParamPrimitiveTypeInt, *literal.Int, nil
	} else if literal.Float != nil {
		return schema.ParamPrimitiveTypeFloat, *literal.Float, nil
	} else if literal.True != nil {
		return schema.ParamPrimitiveTypeBoolean, "true", nil
	} else if literal.False != nil {
		return schema.ParamPrimitiveTypeBoolean, "false", nil
	}

	return schema.ParamPrimitiveType{}, "", fmt.Errorf("empty literal value")
}

// convertParamTypeToJSON converts a string parameter type to a schema ParamPrimitiveType
func convertParamTypeToJSON(paramType string) (schema.ParamPrimitiveType, error) {
	switch paramType {
	case "string":
		return schema.ParamPrimitiveTypeString, nil
	case "int":
		return schema.ParamPrimitiveTypeInt, nil
	case "float":
		return schema.ParamPrimitiveTypeFloat, nil
	case "boolean":
		return schema.ParamPrimitiveTypeBoolean, nil
	default:
		return schema.ParamPrimitiveType{}, fmt.Errorf("invalid parameter type: %s", paramType)
	}
}
