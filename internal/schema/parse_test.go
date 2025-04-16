package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/util/testutil"
)

func TestParseSchema(t *testing.T) {
	t.Run("Valid empty schema", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": []
		}`

		schema, err := ParseSchema(input)
		require.NoError(t, err)
		require.Equal(t, 1, schema.Version)
		require.Empty(t, schema.Nodes)
	})

	t.Run("Schema with doc node", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "doc",
					"content": "This is documentation"
				}
			]
		}`

		schema, err := ParseSchema(input)
		require.NoError(t, err)
		require.Equal(t, 1, schema.Version)
		require.Len(t, schema.Nodes, 1)

		docNode, ok := schema.Nodes[0].(*NodeDoc)
		require.True(t, ok, "Node should be a NodeDoc")
		require.Equal(t, "doc", docNode.Kind)
		require.Equal(t, "This is documentation", docNode.Content)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
		}`

		_, err := ParseSchema(input)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to unmarshal input schema")
	})

}

func TestParamDefinition(t *testing.T) {
	t.Run("Unmarshal ParamDefinition", func(t *testing.T) {
		input := `{
			"type": "string",
			"isArray": false
		}`

		var paramDef ParamDefinition
		err := json.Unmarshal([]byte(input), &paramDef)
		require.NoError(t, err)
		require.Equal(t, ParamPrimitiveTypeString, paramDef.Type)
		require.False(t, paramDef.IsArray)
	})

	t.Run("Unmarshal ParamDefinition with array", func(t *testing.T) {
		input := `{
			"type": "int",
			"isArray": true
		}`

		var paramDef ParamDefinition
		err := json.Unmarshal([]byte(input), &paramDef)
		require.NoError(t, err)
		require.Equal(t, ParamPrimitiveTypeInt, paramDef.Type)
		require.True(t, paramDef.IsArray)
	})

	t.Run("Marshal ParamDefinition", func(t *testing.T) {
		paramDef := ParamDefinition{
			Type:    ParamPrimitiveTypeString,
			IsArray: false,
		}

		data, err := json.Marshal(paramDef)
		require.NoError(t, err)

		var result map[string]any
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)
		require.Equal(t, "string", result["type"])
		require.Equal(t, false, result["isArray"])
	})
}

func TestAppliedParam(t *testing.T) {
	t.Run("Unmarshal AppliedParam with single value", func(t *testing.T) {
		input := `{
			"type": "string",
			"isArray": false,
			"value": "test"
		}`

		var appliedParam AppliedParam
		err := json.Unmarshal([]byte(input), &appliedParam)
		require.NoError(t, err)
		require.Equal(t, ParamPrimitiveTypeString, appliedParam.Type)
		require.False(t, appliedParam.IsArray)
		require.Equal(t, "test", appliedParam.Value)
		require.Empty(t, appliedParam.ArrayValues)
	})

	t.Run("Unmarshal AppliedParam with array values", func(t *testing.T) {
		input := `{
			"type": "int",
			"isArray": true,
			"arrayValues": ["1", "2", "3"]
		}`

		var appliedParam AppliedParam
		err := json.Unmarshal([]byte(input), &appliedParam)
		require.NoError(t, err)
		require.Equal(t, ParamPrimitiveTypeInt, appliedParam.Type)
		require.True(t, appliedParam.IsArray)
		require.Empty(t, appliedParam.Value)
		require.Equal(t, []string{"1", "2", "3"}, appliedParam.ArrayValues)
	})

	t.Run("Marshal AppliedParam with single value", func(t *testing.T) {
		appliedParam := AppliedParam{
			Type:    ParamPrimitiveTypeString,
			IsArray: false,
			Value:   "test",
		}

		data, err := json.Marshal(appliedParam)
		require.NoError(t, err)

		var result map[string]any
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)
		require.Equal(t, "string", result["type"])
		require.Equal(t, false, result["isArray"])
		require.Equal(t, "test", result["value"])
		require.Nil(t, result["arrayValues"])
	})

	t.Run("Marshal AppliedParam with array values", func(t *testing.T) {
		appliedParam := AppliedParam{
			Type:        ParamPrimitiveTypeInt,
			IsArray:     true,
			ArrayValues: []string{"1", "2", "3"},
		}

		data, err := json.Marshal(appliedParam)
		require.NoError(t, err)

		var result map[string]any
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)
		require.Equal(t, "int", result["type"])
		require.Equal(t, true, result["isArray"])
		require.Nil(t, result["value"])

		arrayValues, ok := result["arrayValues"].([]any)
		require.True(t, ok)
		require.Len(t, arrayValues, 3)
		require.Equal(t, "1", arrayValues[0])
		require.Equal(t, "2", arrayValues[1])
		require.Equal(t, "3", arrayValues[2])
	})
}

func TestPrimitiveType(t *testing.T) {
	t.Run("Marshal PrimitiveType", func(t *testing.T) {
		primitiveType := PrimitiveTypeString

		data, err := json.Marshal(primitiveType)
		require.NoError(t, err)
		require.Equal(t, `"string"`, string(data))
	})

	t.Run("Unmarshal PrimitiveType", func(t *testing.T) {
		input := `"string"`

		var primitiveType PrimitiveType
		err := json.Unmarshal([]byte(input), &primitiveType)
		require.NoError(t, err)
		require.Equal(t, PrimitiveTypeString, primitiveType)
	})
}

func TestNodeKind(t *testing.T) {
	t.Run("NodeDoc.NodeKind", func(t *testing.T) {
		node := NodeDoc{
			Kind:    "doc",
			Content: "Documentation",
		}
		require.Equal(t, "doc", node.NodeKind())
	})

	t.Run("NodeRule.NodeKind", func(t *testing.T) {
		node := NodeRule{
			Kind: "rule",
			Name: "required",
			For:  "string",
		}
		require.Equal(t, "rule", node.NodeKind())
	})

	t.Run("NodeType.NodeKind", func(t *testing.T) {
		node := NodeType{
			Kind: "type",
			Name: "User",
		}
		require.Equal(t, "type", node.NodeKind())
	})

	t.Run("NodeProc.NodeKind", func(t *testing.T) {
		node := NodeProc{
			Kind: "proc",
			Name: "GetUser",
		}
		require.Equal(t, "proc", node.NodeKind())
	})
}

func TestBasicSchemaUnmarshal(t *testing.T) {
	t.Run("Valid empty schema", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": []
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)
		require.Equal(t, 1, schema.Version)
		require.Empty(t, schema.Nodes)
	})

	t.Run("Schema with doc node", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "doc",
					"content": "This is documentation"
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)
		require.Equal(t, 1, schema.Version)
		require.Len(t, schema.Nodes, 1)

		docNode, ok := schema.Nodes[0].(*NodeDoc)
		require.True(t, ok, "Node should be a NodeDoc")
		require.Equal(t, "doc", docNode.Kind)
		require.Equal(t, "This is documentation", docNode.Content)
	})

	t.Run("Schema with rule node", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "rule",
					"name": "required",
					"for": "string",
					"doc": "Required field rule",
					"error": "Field is required"
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)
		require.Equal(t, 1, schema.Version)
		require.Len(t, schema.Nodes, 1)

		ruleNode, ok := schema.Nodes[0].(*NodeRule)
		require.True(t, ok, "Node should be a NodeRule")
		require.Equal(t, "rule", ruleNode.Kind)
		require.Equal(t, "required", ruleNode.Name)
		require.Equal(t, "string", ruleNode.For)
		require.NotNil(t, ruleNode.Doc)
		require.Equal(t, "Required field rule", *ruleNode.Doc)
		require.NotNil(t, ruleNode.Error)
		require.Equal(t, "Field is required", *ruleNode.Error)
	})

	t.Run("Schema with type node", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "type",
					"name": "User",
					"doc": "User type",
					"fields": [
						{
							"name": "id",
							"typeName": "string",
							"depth": 0,
							"optional": false,
							"rules": []
						},
						{
							"name": "name",
							"typeName": "string",
							"depth": 0,
							"optional": false,
							"rules": []
						}
					]
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)
		require.Equal(t, 1, schema.Version)
		require.Len(t, schema.Nodes, 1)

		// Check that the node is a type node
		typeNode, ok := schema.Nodes[0].(*NodeType)
		require.True(t, ok, "Node should be a NodeType")
		require.Equal(t, "type", typeNode.Kind)
		require.Equal(t, "User", typeNode.Name)
		require.NotNil(t, typeNode.Doc)
		require.Equal(t, "User type", *typeNode.Doc)
		require.Len(t, typeNode.Fields, 2)

		// Check the fields
		require.Equal(t, "id", typeNode.Fields[0].Name)
		require.NotNil(t, typeNode.Fields[0].TypeName)
		require.Equal(t, "string", *typeNode.Fields[0].TypeName)
		require.Equal(t, 0, typeNode.Fields[0].Depth)
		require.False(t, typeNode.Fields[0].Optional)
		require.Empty(t, typeNode.Fields[0].Rules)

		require.Equal(t, "name", typeNode.Fields[1].Name)
		require.NotNil(t, typeNode.Fields[1].TypeName)
		require.Equal(t, "string", *typeNode.Fields[1].TypeName)
		require.Equal(t, 0, typeNode.Fields[1].Depth)
		require.False(t, typeNode.Fields[1].Optional)
		require.Empty(t, typeNode.Fields[1].Rules)
	})

	t.Run("Schema with proc node", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "proc",
					"name": "GetUser",
					"doc": "Get user by ID",
					"input": [
						{
							"name": "id",
							"typeName": "string",
							"depth": 0,
							"optional": false,
							"rules": []
						}
					],
					"output": [
						{
							"name": "user",
							"typeName": "User",
							"depth": 0,
							"optional": false,
							"rules": []
						}
					],
					"meta": {
						"http.method": "GET",
						"http.path": "/users/{id}",
						"auth": true
					}
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)
		require.Equal(t, 1, schema.Version)
		require.Len(t, schema.Nodes, 1)

		// Check that the node is a proc node
		procNode, ok := schema.Nodes[0].(*NodeProc)
		require.True(t, ok, "Node should be a NodeProc")
		require.Equal(t, "proc", procNode.Kind)
		require.Equal(t, "GetUser", procNode.Name)
		require.NotNil(t, procNode.Doc)
		require.Equal(t, "Get user by ID", *procNode.Doc)

		// Check input fields
		require.Len(t, procNode.Input, 1)
		require.Equal(t, "id", procNode.Input[0].Name)
		require.NotNil(t, procNode.Input[0].TypeName)
		require.Equal(t, "string", *procNode.Input[0].TypeName)

		// Check output fields
		require.Len(t, procNode.Output, 1)
		require.Equal(t, "user", procNode.Output[0].Name)
		require.NotNil(t, procNode.Output[0].TypeName)
		require.Equal(t, "User", *procNode.Output[0].TypeName)

		// Check metadata
		require.Len(t, procNode.Meta, 3)

		// Check http.method
		httpMethod, ok := procNode.Meta["http.method"]
		require.True(t, ok)
		require.NotNil(t, httpMethod.StringVal)
		require.Equal(t, "GET", *httpMethod.StringVal)

		// Check http.path
		httpPath, ok := procNode.Meta["http.path"]
		require.True(t, ok)
		require.NotNil(t, httpPath.StringVal)
		require.Equal(t, "/users/{id}", *httpPath.StringVal)

		// Check auth
		auth, ok := procNode.Meta["auth"]
		require.True(t, ok)
		require.NotNil(t, auth.BoolVal)
		require.True(t, *auth.BoolVal)
	})
}

func TestGetNodeMethods(t *testing.T) {
	input := `{
		"version": 1,
		"nodes": [
			{
				"kind": "doc",
				"content": "Documentation"
			},
			{
				"kind": "rule",
				"name": "required",
				"for": "string"
			},
			{
				"kind": "type",
				"name": "User",
				"fields": [
					{
						"name": "id",
						"typeName": "string",
						"depth": 0,
						"optional": false,
						"rules": []
					}
				]
			},
			{
				"kind": "proc",
				"name": "GetUser",
				"input": [
					{
						"name": "id",
						"typeName": "string",
						"depth": 0,
						"optional": false,
						"rules": []
					}
				],
				"output": [
					{
						"name": "user",
						"typeName": "User",
						"depth": 0,
						"optional": false,
						"rules": []
					}
				],
				"meta": {}
			}
		]
	}`

	var schema Schema
	err := json.Unmarshal([]byte(input), &schema)
	require.NoError(t, err)

	// Test GetDocNodes
	docNodes := schema.GetDocNodes()
	require.Len(t, docNodes, 1)
	require.Equal(t, "Documentation", docNodes[0].Content)

	// Test GetRuleNodes
	ruleNodes := schema.GetRuleNodes()
	require.Len(t, ruleNodes, 1)
	require.Equal(t, "required", ruleNodes[0].Name)

	// Test GetTypeNodes
	typeNodes := schema.GetTypeNodes()
	require.Len(t, typeNodes, 1)
	require.Equal(t, "User", typeNodes[0].Name)

	// Test GetProcNodes
	procNodes := schema.GetProcNodes()
	require.Len(t, procNodes, 1)
	require.Equal(t, "GetUser", procNodes[0].Name)
}

func TestFieldDefinitionHelperMethods(t *testing.T) {
	// Test IsNamed
	namedField := FieldDefinition{
		Name:     "user",
		TypeName: testutil.Pointer("User"),
	}
	require.True(t, namedField.IsNamed())
	require.False(t, namedField.IsInline())

	// Test IsInline
	inlineField := FieldDefinition{
		Name: "address",
		TypeInline: &InlineTypeDefinition{
			Fields: []FieldDefinition{},
		},
	}
	require.False(t, inlineField.IsNamed())
	require.True(t, inlineField.IsInline())

	// Test neither named nor inline
	emptyField := FieldDefinition{
		Name: "empty",
	}
	require.False(t, emptyField.IsNamed())
	require.False(t, emptyField.IsInline())
}

func TestBasicMetaValue(t *testing.T) {
	// Test string value
	t.Run("String value", func(t *testing.T) {
		strVal := "test"
		mv := MetaValue{StringVal: &strVal}

		data, err := mv.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, `"test"`, string(data))

		var newMV MetaValue
		err = newMV.UnmarshalJSON(data)
		require.NoError(t, err)
		require.NotNil(t, newMV.StringVal)
		require.Equal(t, strVal, *newMV.StringVal)
	})

	// Test integer value
	t.Run("Integer value", func(t *testing.T) {
		intVal := int64(42)
		mv := MetaValue{IntVal: &intVal}

		data, err := mv.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, `42`, string(data))

		var newMV MetaValue
		err = newMV.UnmarshalJSON(data)
		require.NoError(t, err)
		require.NotNil(t, newMV.IntVal)
		require.Equal(t, intVal, *newMV.IntVal)
	})

	t.Run("Float value", func(t *testing.T) {
		floatVal := 3.14
		mv := MetaValue{FloatVal: &floatVal}

		data, err := mv.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, `3.14`, string(data))

		var newMV MetaValue
		err = newMV.UnmarshalJSON(data)
		require.NoError(t, err)
		require.NotNil(t, newMV.FloatVal)
		require.Equal(t, floatVal, *newMV.FloatVal)
	})

	t.Run("Boolean value", func(t *testing.T) {
		boolVal := true
		mv := MetaValue{BoolVal: &boolVal}

		data, err := mv.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, `true`, string(data))

		var newMV MetaValue
		err = newMV.UnmarshalJSON(data)
		require.NoError(t, err)
		require.NotNil(t, newMV.BoolVal)
		require.Equal(t, boolVal, *newMV.BoolVal)
	})
}
