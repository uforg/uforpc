package parser

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/urpc/ast"
)

// setEmptyPos uses reflection to check if the node has Pos or EndPos fields
// and sets them to the empty position (emptyPos).
//
// It handles pointers and non-pointers correctly.
func setEmptyPos[T any](node T) T {
	valueOfNode, valueOfEmptyPos := reflect.ValueOf(node), reflect.ValueOf(ast.Position{})
	fix := func(s reflect.Value) {
		if f := s.FieldByName("Pos"); f.IsValid() && f.CanSet() && f.Type() == valueOfEmptyPos.Type() {
			f.Set(valueOfEmptyPos)
		}
		if f := s.FieldByName("EndPos"); f.IsValid() && f.CanSet() && f.Type() == valueOfEmptyPos.Type() {
			f.Set(valueOfEmptyPos)
		}
	}
	switch valueOfNode.Kind() {
	case reflect.Ptr:
		if !valueOfNode.IsNil() && valueOfNode.Elem().Kind() == reflect.Struct {
			fix(valueOfNode.Elem())
		}
		return node
	case reflect.Struct:
		s := reflect.New(valueOfNode.Type()).Elem()
		s.Set(valueOfNode)
		fix(s)
		return s.Interface().(T)
	default:
		return node
	}
}

// equal compares two URPC structs and fails if they are not equal.
// The validation includes the positions of the AST nodes.
func equal(t *testing.T, expected, actual *ast.URPCSchema) {
	t.Helper()
	require.Equal(t, expected, actual)
}

// equalNoPos compares two URPC structs and fails if they are not equal.
// It ignores the positions of any nested AST nodes.
func equalNoPos(t *testing.T, expected, actual *ast.URPCSchema) {
	t.Helper()

	cleanPositions := func(ast *ast.URPCSchema) *ast.URPCSchema {
		ast = setEmptyPos(ast)

		if ast.Version != nil {
			ast.Version = setEmptyPos(ast.Version)
		}

		for i, importStmt := range ast.Imports {
			ast.Imports[i] = setEmptyPos(importStmt)
		}

		for i, rule := range ast.Rules {
			ast.Rules[i] = setEmptyPos(rule)
		}

		for i, typeDecl := range ast.Types {
			ast.Types[i] = setEmptyPos(typeDecl)
		}

		return ast
	}

	expected = cleanPositions(expected)
	actual = cleanPositions(actual)
	equal(t, expected, actual)
}

// ptr creates a pointer to the given value.
func ptr[T any](v T) *T {
	return &v
}

func TestParserPositions(t *testing.T) {
	t.Run("Version position", func(t *testing.T) {
		input := `version 1`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)
		require.NotNil(t, parsed)

		expected := &ast.URPCSchema{
			Version: &ast.Version{
				Pos: ast.Position{
					Filename: "schema.urpc",
					Line:     1,
					Offset:   0,
					Column:   1,
				},
				EndPos: ast.Position{
					Filename: "schema.urpc",
					Line:     1,
					Offset:   9,
					Column:   10,
				},
				Number: 1,
			},
		}

		equal(t, expected, parsed)
	})
}

func TestParserVersion(t *testing.T) {
	t.Run("Correct version parsing", func(t *testing.T) {
		input := `version 1`
		parsed, err := Parser.ParseString("schema.urpc", input)

		require.NoError(t, err)
		require.NotNil(t, parsed)

		expected := &ast.URPCSchema{
			Version: &ast.Version{
				Number: 1,
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("More than one version should fail", func(t *testing.T) {
		input := `version 1 version: 2`
		_, err := Parser.ParseString("schema.urpc", input)
		require.Error(t, err)
	})

	t.Run("Version as float should fail", func(t *testing.T) {
		input := `version 1.0`
		_, err := Parser.ParseString("schema.urpc", input)
		require.Error(t, err)
	})

	t.Run("Version as identifier should fail", func(t *testing.T) {
		input := `version: version`
		_, err := Parser.ParseString("schema.urpc", input)
		require.Error(t, err)
	})

	t.Run("Version as string should fail", func(t *testing.T) {
		input := `version: "1"`
		_, err := Parser.ParseString("schema.urpc", input)
		require.Error(t, err)
	})
}

func TestParserImport(t *testing.T) {
	t.Run("Import parsing", func(t *testing.T) {
		input := `import "../../my_sub_schema.urpc"`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Imports: []*ast.Import{
				{Path: "../../my_sub_schema.urpc"},
			},
		}

		equalNoPos(t, expected, parsed)
	})
}

func TestParserRuleDecl(t *testing.T) {
	t.Run("Minimum rule declaration parsing", func(t *testing.T) {
		input := `
				rule @myRule {
					for: string
				}
			`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Rules: []*ast.RuleDecl{
				{
					Name: "myRule",
					Body: ast.RuleDeclBody{
						For: "string",
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Empty rule not allowed", func(t *testing.T) {
		input := `rule @myRule {}`
		_, err := Parser.ParseString("schema.urpc", input)
		require.Error(t, err)
	})

	t.Run("Rule with no body not allowed", func(t *testing.T) {
		input := `rule @myRule`
		_, err := Parser.ParseString("schema.urpc", input)
		require.Error(t, err)
	})

	t.Run("For in rule body is required", func(t *testing.T) {
		input := `rule @myRule { param: string }`
		_, err := Parser.ParseString("schema.urpc", input)
		require.Error(t, err)
	})

	t.Run("Rule with docstring", func(t *testing.T) {
		input := `
				"""
				My rule description
				"""
				rule @myRule {
					for: string
				}
			`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Rules: []*ast.RuleDecl{
				{
					Docstring: "My rule description",
					Name:      "myRule",
					Body: ast.RuleDeclBody{
						For: "string",
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Rule with array param", func(t *testing.T) {
		input := `
				rule @myRule {
					for: string
					param: string[]
				}
			`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Rules: []*ast.RuleDecl{
				{
					Name: "myRule",
					Body: ast.RuleDeclBody{
						For:          "string",
						Param:        "string",
						ParamIsArray: true,
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Rule with all options", func(t *testing.T) {
		input := `
				"""
				My rule description
				"""
				rule @myRule {
					for: MyType
					param: float[]
					error: "My error message"
				}
			`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Rules: []*ast.RuleDecl{
				{
					Docstring: "My rule description",
					Name:      "myRule",
					Body: ast.RuleDeclBody{
						For:          "MyType",
						Param:        "float",
						ParamIsArray: true,
						Error:        "My error message",
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})
}

func TestParserTypeDecl(t *testing.T) {
	t.Run("Minimum type declaration parsing", func(t *testing.T) {
		input := `
			type MyType {
				field: string
			}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Types: []*ast.TypeDecl{
				{
					Name: "MyType",
					Fields: []*ast.Field{
						{
							Name: "field",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{Named: ptr("string")},
							},
						},
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Type declaration With Docstring", func(t *testing.T) {
		input := `
			"""
			My type description
			"""
			type MyType {
				field: string
			}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Types: []*ast.TypeDecl{
				{
					Docstring: "My type description",
					Name:      "MyType",
					Fields: []*ast.Field{
						{
							Name: "field",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{Named: ptr("string")},
							},
						},
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Type declaration with extends", func(t *testing.T) {
		input := `
			type MyType extends OtherType {
				field: string
			}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Types: []*ast.TypeDecl{
				{
					Name:    "MyType",
					Extends: []string{"OtherType"},
					Fields: []*ast.Field{
						{
							Name: "field",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{Named: ptr("string")},
							},
						},
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Type declaration with multiple extends", func(t *testing.T) {
		input := `
			type MyType extends OtherType, AnotherType, YetAnotherType {
				field: string
			}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Types: []*ast.TypeDecl{
				{
					Name:    "MyType",
					Extends: []string{"OtherType", "AnotherType", "YetAnotherType"},
					Fields: []*ast.Field{
						{
							Name: "field",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{Named: ptr("string")},
							},
						},
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Type declaration with extends and docstring", func(t *testing.T) {
		input := `
			"""
			My type description
			"""
			type MyType extends OtherType {
				field: string
			}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Types: []*ast.TypeDecl{
				{
					Docstring: "My type description",
					Name:      "MyType",
					Extends:   []string{"OtherType"},
					Fields: []*ast.Field{
						{
							Name: "field",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{Named: ptr("string")},
							},
						},
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Empty type not allowed", func(t *testing.T) {
		input := `type MyType {}`
		_, err := Parser.ParseString("schema.urpc", input)
		require.Error(t, err)
	})

	t.Run("Type declaration with custom type field", func(t *testing.T) {
		input := `
			type MyType {
				field: MyCustomType
			}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Types: []*ast.TypeDecl{
				{
					Name: "MyType",
					Fields: []*ast.Field{
						{
							Name: "field",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{Named: ptr("MyCustomType")},
							},
						},
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})
}

func TestParserField(t *testing.T) {
	t.Run("Fields with primitive types", func(t *testing.T) {
		input := `
			type MyType {
				field1: string
				field2: int
				field3: float
				field4: boolean
				field5: datetime
			}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Types: []*ast.TypeDecl{
				{
					Name: "MyType",
					Fields: []*ast.Field{
						{
							Name: "field1",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{Named: ptr("string")},
							},
						},
						{
							Name: "field2",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{Named: ptr("int")},
							},
						},
						{
							Name: "field3",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{Named: ptr("float")},
							},
						},
						{
							Name: "field4",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{Named: ptr("boolean")},
							},
						},
						{
							Name: "field5",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{Named: ptr("datetime")},
							},
						},
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Fields with custom types", func(t *testing.T) {
		input := `
			type MyType {
				field1: MyCustomType
				field2: MyOtherCustomType
			}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Types: []*ast.TypeDecl{
				{
					Name: "MyType",
					Fields: []*ast.Field{
						{
							Name: "field1",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{Named: ptr("MyCustomType")},
							},
						},
						{
							Name: "field2",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{Named: ptr("MyOtherCustomType")},
							},
						},
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Optional fields", func(t *testing.T) {
		input := `
			type MyType {
				field1?: string
				field2?: MyCustomType
			}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Types: []*ast.TypeDecl{
				{
					Name: "MyType",
					Fields: []*ast.Field{
						{
							Name: "field1",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{Named: ptr("string")},
							},
							Optional: true,
						},
						{
							Name: "field2",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{Named: ptr("MyCustomType")},
							},
							Optional: true,
						},
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Complex array and nested object fields", func(t *testing.T) {
		input := `
			type MyType {
				field1: string[]
				field2: {
					subfield: string
				}
				field3: int[][]
				field4?: {
					subfield: {
						subsubfield: datetime[][][]
					}[][]
				}[][][][][][][]
			}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Types: []*ast.TypeDecl{
				{
					Name: "MyType",
					Fields: []*ast.Field{
						{
							Name: "field1",
							Type: ast.FieldType{
								Depth: 1,
								Base: &ast.FieldTypeBase{
									Named: ptr("string"),
								},
							},
						},
						{
							Name: "field2",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{
									Object: &ast.FieldTypeObject{
										Fields: []*ast.Field{
											{
												Name: "subfield",
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{Named: ptr("string")},
												},
											},
										},
									},
								},
							},
						},
						{
							Name: "field3",
							Type: ast.FieldType{
								Depth: 2,
								Base: &ast.FieldTypeBase{
									Named: ptr("int"),
								},
							},
						},
						{
							Name:     "field4",
							Optional: true,
							Type: ast.FieldType{
								Depth: 7,
								Base: &ast.FieldTypeBase{
									Object: &ast.FieldTypeObject{
										Fields: []*ast.Field{
											{
												Name: "subfield",
												Type: ast.FieldType{
													Depth: 2,
													Base: &ast.FieldTypeBase{
														Object: &ast.FieldTypeObject{
															Fields: []*ast.Field{
																{
																	Name: "subsubfield",
																	Type: ast.FieldType{
																		Depth: 3,
																		Base: &ast.FieldTypeBase{
																			Named: ptr("datetime"),
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Field with rules", func(t *testing.T) {
		input := `
			type MyType {
				field1: string
					@uppercase
					@uppercase()
					@uppercase(error: "Field must be uppercase")
					@contains("hello", error: "Field must contain 'hello'")
					@enum(["hello", "world"], error: "Field must be 'hello' or 'world'")
					@enum([1, 2, 3])
					@enum([1.1, 2.2, 3.3])
					@enum([true, false])
			}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Types: []*ast.TypeDecl{
				{
					Name: "MyType",
					Fields: []*ast.Field{
						{
							Name: "field1",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{Named: ptr("string")},
							},
							Rules: []*ast.FieldRule{
								{
									Name: "uppercase",
								},
								{
									Name: "uppercase",
								},
								{
									Name: "uppercase",
									Body: ast.FieldRuleBody{
										Error: "Field must be uppercase",
									},
								},
								{
									Name: "contains",
									Body: ast.FieldRuleBody{
										ParamSingle: ptr("hello"),
										Error:       "Field must contain 'hello'",
									},
								},
								{
									Name: "enum",
									Body: ast.FieldRuleBody{
										ParamList: []string{"hello", "world"},
										Error:     "Field must be 'hello' or 'world'",
									},
								},
								{
									Name: "enum",
									Body: ast.FieldRuleBody{
										ParamList: []string{"1", "2", "3"},
									},
								},
								{
									Name: "enum",
									Body: ast.FieldRuleBody{
										ParamList: []string{"1.1", "2.2", "3.3"},
									},
								},
								{
									Name: "enum",
									Body: ast.FieldRuleBody{
										ParamList: []string{"true", "false"},
									},
								},
							},
						},
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})
}
