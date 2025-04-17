package transpile

import (
	"fmt"

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

	// Process all nodes in the original order
	for _, child := range astSchema.Children {
		switch {
		case child.Version != nil:
			result.Version = child.Version.Number

		case child.Docstring != nil:
			docNode := &schema.NodeDoc{
				Kind:    "doc",
				Content: child.Docstring.Value,
			}
			result.Nodes = append(result.Nodes, docNode)

		case child.Rule != nil:
			ruleNode, err := convertRuleToJSON(child.Rule)
			if err != nil {
				return schema.Schema{}, fmt.Errorf("error converting rule '%s': %w", child.Rule.Name, err)
			}
			result.Nodes = append(result.Nodes, ruleNode)

		case child.Type != nil:
			typeNode, err := convertTypeToJSON(child.Type)
			if err != nil {
				return schema.Schema{}, fmt.Errorf("error converting type '%s': %w", child.Type.Name, err)
			}
			result.Nodes = append(result.Nodes, typeNode)

		case child.Proc != nil:
			procNode, err := convertProcToJSON(child.Proc)
			if err != nil {
				return schema.Schema{}, fmt.Errorf("error converting procedure '%s': %w", child.Proc.Name, err)
			}
			result.Nodes = append(result.Nodes, procNode)
		}
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
			errorMsg := child.Error.Error
			ruleNode.Error = &errorMsg
		}
	}

	return ruleNode, nil
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

// convertTypeToJSON converts an AST TypeDecl to a schema NodeType
func convertTypeToJSON(typeDecl *ast.TypeDecl) (*schema.NodeType, error) {
	typeNode := &schema.NodeType{
		Kind: "type",
		Name: typeDecl.Name,
	}

	// Add docstring if available
	if typeDecl.Docstring != nil {
		docValue := typeDecl.Docstring.Value
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

// convertFieldToJSON converts an AST Field to a schema FieldDefinition
func convertFieldToJSON(field *ast.Field) (schema.FieldDefinition, error) {
	fieldDef := schema.FieldDefinition{
		Name:     field.Name,
		Optional: field.Optional,
		Depth:    int(field.Type.Depth),
	}

	// Process field type
	if field.Type.Base.Named != nil {
		typeName := *field.Type.Base.Named
		fieldDef.TypeName = &typeName
	}

	if field.Type.Base.Object != nil {
		inlineType := &schema.InlineTypeDefinition{}

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

	// Early return if there's no rule body
	if fieldRule.Body == nil {
		return rule, nil
	}

	// Process error message if available
	if fieldRule.Body.Error != nil {
		errorMsg := *fieldRule.Body.Error
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
			param.ArrayValues = fieldRule.Body.ParamListString
		}
		if len(fieldRule.Body.ParamListInt) > 0 {
			param.Type = schema.ParamPrimitiveTypeInt
			param.IsArray = true
			param.ArrayValues = fieldRule.Body.ParamListInt
		}
		if len(fieldRule.Body.ParamListFloat) > 0 {
			param.Type = schema.ParamPrimitiveTypeFloat
			param.IsArray = true
			param.ArrayValues = fieldRule.Body.ParamListFloat
		}
		if len(fieldRule.Body.ParamListBoolean) > 0 {
			param.Type = schema.ParamPrimitiveTypeBoolean
			param.IsArray = true
			param.ArrayValues = fieldRule.Body.ParamListBoolean
		}

		rule.Param = param
	}

	return rule, nil
}

// extractParamValue extracts the parameter type and value from an AnyLiteral
func extractParamValue(literal ast.AnyLiteral) (schema.ParamPrimitiveType, string, error) {
	if literal.Str != nil {
		return schema.ParamPrimitiveTypeString, *literal.Str, nil
	}
	if literal.Int != nil {
		return schema.ParamPrimitiveTypeInt, *literal.Int, nil
	}
	if literal.Float != nil {
		return schema.ParamPrimitiveTypeFloat, *literal.Float, nil
	}
	if literal.True != nil {
		return schema.ParamPrimitiveTypeBoolean, "true", nil
	}
	if literal.False != nil {
		return schema.ParamPrimitiveTypeBoolean, "false", nil
	}

	return schema.ParamPrimitiveType{}, "", fmt.Errorf("empty literal value")
}

// convertProcToJSON converts an AST ProcDecl to a schema NodeProc
func convertProcToJSON(procDecl *ast.ProcDecl) (*schema.NodeProc, error) {
	procNode := &schema.NodeProc{
		Kind: "proc",
		Name: procDecl.Name,
	}

	// Add docstring if available
	if procDecl.Docstring != nil {
		docValue := procDecl.Docstring.Value
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
		}
		if child.Output != nil {
			for _, fieldOrComment := range child.Output.Children {
				if fieldOrComment.Field != nil {
					fieldDef, err := convertFieldToJSON(fieldOrComment.Field)
					if err != nil {
						return nil, fmt.Errorf("error converting output field '%s': %w", fieldOrComment.Field.Name, err)
					}
					procNode.Output = append(procNode.Output, fieldDef)
				}
			}
		}
		if child.Meta != nil {
			for _, metaChild := range child.Meta.Children {
				if metaChild.KV != nil {
					metaValue, err := convertMetaValueToJSON(metaChild.KV.Value)
					if err != nil {
						return nil, fmt.Errorf("error converting meta value for key '%s': %w", metaChild.KV.Key, err)
					}
					procNode.Meta = append(procNode.Meta, schema.MetaKeyValue{
						Key:   metaChild.KV.Key,
						Value: metaValue,
					})
				}
			}
		}
	}

	return procNode, nil
}

// convertMetaValueToJSON converts an AST AnyLiteral to a schema MetaValue
func convertMetaValueToJSON(literal ast.AnyLiteral) (schema.MetaValue, error) {
	var metaValue schema.MetaValue

	if literal.Str == nil && literal.Int == nil && literal.Float == nil && literal.True == nil && literal.False == nil {
		return schema.MetaValue{}, fmt.Errorf("empty meta value")
	}

	if literal.Str != nil {
		strValue := *literal.Str
		metaValue.StringVal = &strValue
	}
	if literal.Int != nil {
		var intValue int64
		if _, err := fmt.Sscanf(*literal.Int, "%d", &intValue); err != nil {
			return schema.MetaValue{}, fmt.Errorf("invalid integer value: %s", *literal.Int)
		}
		metaValue.IntVal = &intValue
	}
	if literal.Float != nil {
		var floatValue float64
		if _, err := fmt.Sscanf(*literal.Float, "%f", &floatValue); err != nil {
			return schema.MetaValue{}, fmt.Errorf("invalid float value: %s", *literal.Float)
		}
		metaValue.FloatVal = &floatValue
	}
	if literal.True != nil {
		boolValue := true
		metaValue.BoolVal = &boolValue
	}
	if literal.False != nil {
		boolValue := false
		metaValue.BoolVal = &boolValue
	}

	return metaValue, nil
}
