package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/urpc/ast"
	"github.com/uforg/uforpc/internal/util/testutil"
)

////////////////
// TEST CASES //
////////////////

func TestParserPositions(t *testing.T) {
	t.Run("Version position", func(t *testing.T) {
		input := `version 1`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)
		require.NotNil(t, parsed)

		expected := &ast.Schema{
			Positions: ast.Positions{
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
			},
			Children: []*ast.SchemaChild{
				{
					Positions: ast.Positions{
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
					},
					Version: &ast.Version{
						Positions: ast.Positions{
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
						},
						Number: 1,
					},
				},
			},
		}

		testutil.ASTEqual(t, expected, parsed)
	})
}

func TestParserVersion(t *testing.T) {
	t.Run("Correct version parsing", func(t *testing.T) {
		input := `version 1`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)

		require.NoError(t, err)
		require.NotNil(t, parsed)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Version: &ast.Version{
						Number: 1,
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("More than one version should fail", func(t *testing.T) {
		input := `version 1 version: 2`
		_, err := ParserInstance.ParseString("schema.urpc", input)
		require.Error(t, err)
	})

	t.Run("Version as float should fail", func(t *testing.T) {
		input := `version 1.0`
		_, err := ParserInstance.ParseString("schema.urpc", input)
		require.Error(t, err)
	})

	t.Run("Version as identifier should fail", func(t *testing.T) {
		input := `version: version`
		_, err := ParserInstance.ParseString("schema.urpc", input)
		require.Error(t, err)
	})

	t.Run("Version as string should fail", func(t *testing.T) {
		input := `version: "1"`
		_, err := ParserInstance.ParseString("schema.urpc", input)
		require.Error(t, err)
	})
}

func TestParserImport(t *testing.T) {
	t.Run("Import parsing", func(t *testing.T) {
		input := `import "../../my_sub_schema.urpc"`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Import: &ast.Import{
						Path: "../../my_sub_schema.urpc",
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})
}

func TestParserRuleDecl(t *testing.T) {
	t.Run("Minimum rule declaration parsing", func(t *testing.T) {
		input := `
				rule @myRule {
					for: string
				}
			`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Rule: &ast.RuleDecl{
						Name: "myRule",
						Children: []*ast.RuleDeclChild{
							{
								For: &ast.RuleDeclChildFor{
									For: "string",
								},
							},
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Rule with no body not allowed", func(t *testing.T) {
		input := `rule @myRule`
		_, err := ParserInstance.ParseString("schema.urpc", input)
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
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Rule: &ast.RuleDecl{
						Docstring: &ast.Docstring{
							Value: "\n\t\t\t\tMy rule description\n\t\t\t\t",
						},
						Name: "myRule",
						Children: []*ast.RuleDeclChild{
							{
								For: &ast.RuleDeclChildFor{
									For: "string",
								},
							},
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Deprecated rule", func(t *testing.T) {
		input := `
			deprecated rule @myRule {
				for: string
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Rule: &ast.RuleDecl{
						Deprecated: &ast.Deprecated{},
						Name:       "myRule",
						Children: []*ast.RuleDeclChild{
							{
								For: &ast.RuleDeclChildFor{
									For: "string",
								},
							},
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Deprecated with message rule", func(t *testing.T) {
		input := `
			deprecated("Use newRule instead")
			rule @myRule {
				for: string
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Rule: &ast.RuleDecl{
						Deprecated: &ast.Deprecated{
							Message: testutil.Pointer("Use newRule instead"),
						},
						Name: "myRule",
						Children: []*ast.RuleDeclChild{
							{
								For: &ast.RuleDeclChildFor{
									For: "string",
								},
							},
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Rule with array param", func(t *testing.T) {
		input := `
				rule @myRule {
					for: string
					param: string[]
				}
			`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Rule: &ast.RuleDecl{
						Name: "myRule",
						Children: []*ast.RuleDeclChild{
							{
								For: &ast.RuleDeclChildFor{
									For: "string",
								},
							},
							{
								Param: &ast.RuleDeclChildParam{
									Param:   "string",
									IsArray: true,
								},
							},
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Rule with array for", func(t *testing.T) {
		input := `
				rule @myRule {
					for: MyType[]
				}
			`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Rule: &ast.RuleDecl{
						Name: "myRule",
						Children: []*ast.RuleDeclChild{
							{
								For: &ast.RuleDeclChildFor{
									For:     "MyType",
									IsArray: true,
								},
							},
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Rule with all options", func(t *testing.T) {
		input := `
				""" My rule description """
				deprecated("This rule is deprecated")
				rule @myRule {
					for: MyType[]
					param: float[]
					error: "My error message"
				}
			`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Rule: &ast.RuleDecl{
						Docstring: &ast.Docstring{
							Value: " My rule description ",
						},
						Deprecated: &ast.Deprecated{
							Message: testutil.Pointer("This rule is deprecated"),
						},
						Name: "myRule",
						Children: []*ast.RuleDeclChild{
							{
								For: &ast.RuleDeclChildFor{
									For:     "MyType",
									IsArray: true,
								},
							},
							{
								Param: &ast.RuleDeclChildParam{
									Param:   "float",
									IsArray: true,
								},
							},
							{
								Error: &ast.RuleDeclChildError{
									Error: "My error message",
								},
							},
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})
}

func TestParserTypeDecl(t *testing.T) {
	t.Run("Minimum type declaration parsing", func(t *testing.T) {
		input := `
			type MyType {
				field: string
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Type: &ast.TypeDecl{
						Name: "MyType",
						Children: []*ast.FieldOrComment{
							{
								Field: &ast.Field{
									Name: "field",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")},
									},
								},
							},
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Type declaration With Docstring", func(t *testing.T) {
		input := `
			""" My type description """
			type MyType {
				field: string
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Type: &ast.TypeDecl{
						Docstring: &ast.Docstring{
							Value: " My type description ",
						},
						Name: "MyType",
						Children: []*ast.FieldOrComment{
							{
								Field: &ast.Field{
									Name: "field",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")},
									},
								},
							},
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Deprecated type", func(t *testing.T) {
		input := `
			deprecated type MyType {
				field: string
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Type: &ast.TypeDecl{
						Deprecated: &ast.Deprecated{},
						Name:       "MyType",
						Children: []*ast.FieldOrComment{
							{
								Field: &ast.Field{
									Name: "field",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")},
									},
								},
							},
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Deprecated with message type", func(t *testing.T) {
		input := `
			deprecated("Use newType instead")
			type MyType {
				field: string
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Type: &ast.TypeDecl{
						Deprecated: &ast.Deprecated{
							Message: testutil.Pointer("Use newType instead"),
						},
						Name: "MyType",
						Children: []*ast.FieldOrComment{
							{
								Field: &ast.Field{
									Name: "field",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")},
									},
								},
							},
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Type declaration with docstring and deprecated", func(t *testing.T) {
		input := `
			""" My type description """
			deprecated("This type is deprecated")
			type MyType {
				field: string
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Type: &ast.TypeDecl{
						Docstring: &ast.Docstring{
							Value: " My type description ",
						},
						Deprecated: &ast.Deprecated{
							Message: testutil.Pointer("This type is deprecated"),
						},
						Name: "MyType",
						Children: []*ast.FieldOrComment{
							{
								Field: &ast.Field{
									Name: "field",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")},
									},
								},
							},
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Type declaration with custom type field", func(t *testing.T) {
		input := `
			type MyType {
				field: MyCustomType
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Type: &ast.TypeDecl{
						Name: "MyType",
						Children: []*ast.FieldOrComment{
							{
								Field: &ast.Field{
									Name: "field",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("MyCustomType")},
									},
								},
							},
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
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
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Type: &ast.TypeDecl{
						Name: "MyType",
						Children: []*ast.FieldOrComment{
							{
								Field: &ast.Field{
									Name: "field1",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")},
									},
								},
							},
							{
								Field: &ast.Field{
									Name: "field2",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("int")},
									},
								},
							},
							{
								Field: &ast.Field{
									Name: "field3",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("float")},
									},
								},
							},
							{
								Field: &ast.Field{
									Name: "field4",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("boolean")},
									},
								},
							},
							{
								Field: &ast.Field{
									Name: "field5",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("datetime")},
									},
								},
							},
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Fields with custom types", func(t *testing.T) {
		input := `
			type MyType {
				field1: MyCustomType
				field2: MyOtherCustomType
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Type: &ast.TypeDecl{
						Name: "MyType",
						Children: []*ast.FieldOrComment{
							{
								Field: &ast.Field{
									Name: "field1",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("MyCustomType")},
									},
								},
							},
							{
								Field: &ast.Field{
									Name: "field2",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("MyOtherCustomType")},
									},
								},
							},
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Optional fields", func(t *testing.T) {
		input := `
			type MyType {
				field1?: string
				field2?: MyCustomType
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Type: &ast.TypeDecl{
						Name: "MyType",
						Children: []*ast.FieldOrComment{
							{
								Field: &ast.Field{
									Name: "field1",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")},
									},
									Optional: true,
								},
							},
							{
								Field: &ast.Field{
									Name: "field2",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("MyCustomType")},
									},
									Optional: true,
								},
							},
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Complex array and nested object fields", func(t *testing.T) {
		input := `
			type MyType {
				field1: string[]
				field2: {
					subfield: string
				}
				field3: int[]
				field4?: {
					subfield: {
						subsubfield: datetime[]
					}[]
				}[]
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Type: &ast.TypeDecl{
						Name: "MyType",
						Children: []*ast.FieldOrComment{
							{
								Field: &ast.Field{
									Name: "field1",
									Type: ast.FieldType{
										IsArray: true,
										Base: &ast.FieldTypeBase{
											Named: testutil.Pointer("string"),
										},
									},
								},
							},
							{
								Field: &ast.Field{
									Name: "field2",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{
											Object: &ast.FieldTypeObject{
												Children: []*ast.FieldOrComment{
													{
														Field: &ast.Field{
															Name: "subfield",
															Type: ast.FieldType{
																Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")},
															},
														},
													},
												},
											},
										},
									},
								},
							},
							{
								Field: &ast.Field{
									Name: "field3",
									Type: ast.FieldType{
										IsArray: true,
										Base: &ast.FieldTypeBase{
											Named: testutil.Pointer("int"),
										},
									},
								},
							},
							{
								Field: &ast.Field{
									Name:     "field4",
									Optional: true,
									Type: ast.FieldType{
										IsArray: true,
										Base: &ast.FieldTypeBase{
											Object: &ast.FieldTypeObject{
												Children: []*ast.FieldOrComment{
													{
														Field: &ast.Field{
															Name: "subfield",
															Type: ast.FieldType{
																IsArray: true,
																Base: &ast.FieldTypeBase{
																	Object: &ast.FieldTypeObject{
																		Children: []*ast.FieldOrComment{
																			{
																				Field: &ast.Field{
																					Name: "subsubfield",
																					Type: ast.FieldType{
																						IsArray: true,
																						Base: &ast.FieldTypeBase{
																							Named: testutil.Pointer("datetime"),
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
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
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
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Type: &ast.TypeDecl{
						Name: "MyType",
						Children: []*ast.FieldOrComment{
							{
								Field: &ast.Field{
									Name: "field1",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")},
									},
									Children: []*ast.FieldChild{
										{
											Rule: &ast.FieldRule{
												Name: "uppercase",
											},
										},
										{
											Rule: &ast.FieldRule{
												Name: "uppercase",
												Body: &ast.FieldRuleBody{},
											},
										},
										{
											Rule: &ast.FieldRule{
												Name: "uppercase",
												Body: &ast.FieldRuleBody{
													Error: testutil.Pointer("Field must be uppercase"),
												},
											},
										},
										{
											Rule: &ast.FieldRule{
												Name: "contains",
												Body: &ast.FieldRuleBody{
													ParamSingle: &ast.AnyLiteral{Str: testutil.Pointer("hello")},
													Error:       testutil.Pointer("Field must contain 'hello'"),
												},
											},
										},
										{
											Rule: &ast.FieldRule{
												Name: "enum",
												Body: &ast.FieldRuleBody{
													ParamListString: []string{"hello", "world"},
													Error:           testutil.Pointer("Field must be 'hello' or 'world'"),
												},
											},
										},
										{
											Rule: &ast.FieldRule{
												Name: "enum",
												Body: &ast.FieldRuleBody{
													ParamListInt: []string{"1", "2", "3"},
												},
											},
										},
										{
											Rule: &ast.FieldRule{
												Name: "enum",
												Body: &ast.FieldRuleBody{
													ParamListFloat: []string{"1.1", "2.2", "3.3"},
												},
											},
										},
										{
											Rule: &ast.FieldRule{
												Name: "enum",
												Body: &ast.FieldRuleBody{
													ParamListBoolean: []string{"true", "false"},
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

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Rules with array parameters of multiple types not allowed", func(t *testing.T) {
		input := `
			type MyType {
				field1: string
					@enum(["hello", 1])
					@enum([1, 1.1])
					@enum([1.1, true])
			}
		`

		_, err := ParserInstance.ParseString("schema.urpc", input)
		require.Error(t, err)
	})
}

func TestParserProcDecl(t *testing.T) {
	t.Run("Minimum procedure declaration parsing", func(t *testing.T) {
		input := `
			proc MyProc {}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Proc: &ast.ProcDecl{
						Name: "MyProc",
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Procedure with docstring", func(t *testing.T) {
		input := `
			""" MyProc is a procedure that does something. """
			proc MyProc {}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Proc: &ast.ProcDecl{
						Docstring: &ast.Docstring{
							Value: " MyProc is a procedure that does something. ",
						},
						Name: "MyProc",
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Procedure with deprecated", func(t *testing.T) {
		input := `
			deprecated proc MyProc {}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Proc: &ast.ProcDecl{
						Deprecated: &ast.Deprecated{},
						Name:       "MyProc",
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Procedure with deprecated with message", func(t *testing.T) {
		input := `
			deprecated("Use newProc instead")
			proc MyProc {}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Proc: &ast.ProcDecl{
						Deprecated: &ast.Deprecated{
							Message: testutil.Pointer("Use newProc instead"),
						},
						Name: "MyProc",
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Procedure with input", func(t *testing.T) {
		input := `
			proc MyProc {
				input {
					field1: string
				}
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Proc: &ast.ProcDecl{
						Name: "MyProc",
						Children: []*ast.ProcDeclChild{
							{
								Input: &ast.ProcDeclChildInput{
									Children: []*ast.FieldOrComment{
										{
											Field: &ast.Field{
												Name: "field1",
												Type: ast.FieldType{Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")}},
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

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Procedure with output", func(t *testing.T) {
		input := `
			proc MyProc {
				output {
					field1: int
				}
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Proc: &ast.ProcDecl{
						Name: "MyProc",
						Children: []*ast.ProcDeclChild{
							{
								Output: &ast.ProcDeclChildOutput{
									Children: []*ast.FieldOrComment{
										{
											Field: &ast.Field{
												Name: "field1",
												Type: ast.FieldType{Base: &ast.FieldTypeBase{Named: testutil.Pointer("int")}},
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

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Procedure with meta", func(t *testing.T) {
		input := `
			proc MyProc {
				meta {
					key1: "hello"
					key2: 123
					key3: 1.23
					key4: true
					key5: false
				}
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Proc: &ast.ProcDecl{
						Name: "MyProc",
						Children: []*ast.ProcDeclChild{
							{
								Meta: &ast.ProcDeclChildMeta{
									Children: []*ast.ProcDeclChildMetaChild{
										{
											KV: &ast.ProcDeclChildMetaKV{Key: "key1", Value: ast.AnyLiteral{Str: testutil.Pointer("hello")}},
										},
										{
											KV: &ast.ProcDeclChildMetaKV{Key: "key2", Value: ast.AnyLiteral{Int: testutil.Pointer("123")}},
										},
										{
											KV: &ast.ProcDeclChildMetaKV{Key: "key3", Value: ast.AnyLiteral{Float: testutil.Pointer("1.23")}},
										},
										{
											KV: &ast.ProcDeclChildMetaKV{Key: "key4", Value: ast.AnyLiteral{True: testutil.Pointer("true")}},
										},
										{
											KV: &ast.ProcDeclChildMetaKV{Key: "key5", Value: ast.AnyLiteral{False: testutil.Pointer("false")}},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Full procedure", func(t *testing.T) {
		input := `
			""" MyProc is a procedure that does something. """
			deprecated("Use newProc instead")
			proc MyProc {
				input {
					input1: string[]
				}
				output {
					output1?: int
				}
				meta {
					key1: "hello"
					key2: 123
					key3: 1.23
					key4: true
					key5: false
				}
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Proc: &ast.ProcDecl{
						Docstring: &ast.Docstring{
							Value: " MyProc is a procedure that does something. ",
						},
						Deprecated: &ast.Deprecated{
							Message: testutil.Pointer("Use newProc instead"),
						},
						Name: "MyProc",
						Children: []*ast.ProcDeclChild{
							{
								Input: &ast.ProcDeclChildInput{
									Children: []*ast.FieldOrComment{
										{
											Field: &ast.Field{
												Name: "input1",
												Type: ast.FieldType{
													IsArray: true,
													Base:    &ast.FieldTypeBase{Named: testutil.Pointer("string")},
												},
											},
										},
									},
								},
							},
							{
								Output: &ast.ProcDeclChildOutput{
									Children: []*ast.FieldOrComment{
										{
											Field: &ast.Field{
												Name:     "output1",
												Optional: true,
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{Named: testutil.Pointer("int")},
												},
											},
										},
									},
								},
							},
							{
								Meta: &ast.ProcDeclChildMeta{
									Children: []*ast.ProcDeclChildMetaChild{
										{
											KV: &ast.ProcDeclChildMetaKV{
												Key:   "key1",
												Value: ast.AnyLiteral{Str: testutil.Pointer("hello")},
											},
										},
										{
											KV: &ast.ProcDeclChildMetaKV{
												Key:   "key2",
												Value: ast.AnyLiteral{Int: testutil.Pointer("123")},
											},
										},
										{
											KV: &ast.ProcDeclChildMetaKV{
												Key:   "key3",
												Value: ast.AnyLiteral{Float: testutil.Pointer("1.23")},
											},
										},
										{
											KV: &ast.ProcDeclChildMetaKV{
												Key:   "key4",
												Value: ast.AnyLiteral{True: testutil.Pointer("true")},
											},
										},
										{
											KV: &ast.ProcDeclChildMetaKV{
												Key:   "key5",
												Value: ast.AnyLiteral{False: testutil.Pointer("false")},
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

		testutil.ASTEqualNoPos(t, expected, parsed)
	})
}

func TestParserComments(t *testing.T) {
	t.Run("Top level comments between declarations", func(t *testing.T) {
		input := `
			// Version comment
			version 1
			/* Import comment */
			import "path/to/file.urpc"
			// Rule comment
			rule @myRule { for: string }
			/* Type comment */
			type MyType { field: int }
			// Proc comment
			proc MyProc {}
			/* Trailing comment */
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Comment: &ast.Comment{Simple: testutil.Pointer(" Version comment")},
				},
				{
					Version: &ast.Version{Number: 1},
				},
				{
					Comment: &ast.Comment{Block: testutil.Pointer(" Import comment ")},
				},
				{
					Import: &ast.Import{Path: "path/to/file.urpc"},
				},
				{
					Comment: &ast.Comment{Simple: testutil.Pointer(" Rule comment")},
				},
				{
					Rule: &ast.RuleDecl{
						Name: "myRule",
						Children: []*ast.RuleDeclChild{
							{
								For: &ast.RuleDeclChildFor{For: "string"},
							},
						},
					},
				},
				{
					Comment: &ast.Comment{Block: testutil.Pointer(" Type comment ")},
				},
				{
					Type: &ast.TypeDecl{
						Name: "MyType",
						Children: []*ast.FieldOrComment{
							{
								Field: &ast.Field{
									Name: "field",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("int")},
									},
								},
							},
						},
					},
				},
				{
					Comment: &ast.Comment{Simple: testutil.Pointer(" Proc comment")},
				},
				{
					Proc: &ast.ProcDecl{Name: "MyProc"},
				},
				{
					Comment: &ast.Comment{Block: testutil.Pointer(" Trailing comment ")},
				},
			},
		}
		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Comments within RuleDecl", func(t *testing.T) {
		input := `
			rule @myRule {
				// Before for
				for: string
				/* Between for and param */
				param: int
				// Before error
				error: "msg"
				// Trailing comment in rule
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Rule: &ast.RuleDecl{
						Name: "myRule",
						Children: []*ast.RuleDeclChild{
							{
								Comment: &ast.Comment{Simple: testutil.Pointer(" Before for")},
							},
							{
								For: &ast.RuleDeclChildFor{For: "string"},
							},
							{
								Comment: &ast.Comment{Block: testutil.Pointer(" Between for and param ")},
							},
							{
								Param: &ast.RuleDeclChildParam{Param: "int"},
							},
							{
								Comment: &ast.Comment{Simple: testutil.Pointer(" Before error")},
							},
							{
								Error: &ast.RuleDeclChildError{Error: "msg"},
							},
							{
								Comment: &ast.Comment{Simple: testutil.Pointer(" Trailing comment in rule")},
							},
						},
					},
				},
			},
		}
		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Comments within TypeDecl", func(t *testing.T) {
		input := `
			type MyType {
				// Before field1
				field1: string
				/* Between field1 and field2 */
				field2?: int
				// Trailing comment in type
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Type: &ast.TypeDecl{
						Name: "MyType",
						Children: []*ast.FieldOrComment{
							{
								Comment: &ast.Comment{Simple: testutil.Pointer(" Before field1")},
							},
							{
								Field: &ast.Field{
									Name: "field1",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")},
									},
								},
							},
							{
								Comment: &ast.Comment{Block: testutil.Pointer(" Between field1 and field2 ")},
							},
							{
								Field: &ast.Field{
									Name:     "field2",
									Optional: true,
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("int")},
									},
								},
							},
							{
								Comment: &ast.Comment{Simple: testutil.Pointer(" Trailing comment in type")},
							},
						},
					},
				},
			},
		}
		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Comments within ProcDecl (between blocks)", func(t *testing.T) {
		input := `
			proc MyProc {
				// Before input
				input { fieldIn: string }
				/* Between input and output */
				output { fieldOut: int }
				// Between output and meta
				meta { key: "value" }
				// Trailing comment in proc
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Proc: &ast.ProcDecl{
						Name: "MyProc",
						Children: []*ast.ProcDeclChild{
							{
								Comment: &ast.Comment{Simple: testutil.Pointer(" Before input")},
							},
							{
								Input: &ast.ProcDeclChildInput{
									Children: []*ast.FieldOrComment{
										{
											Field: &ast.Field{
												Name: "fieldIn",
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")},
												},
											},
										},
									},
								},
							},
							{
								Comment: &ast.Comment{Block: testutil.Pointer(" Between input and output ")},
							},
							{
								Output: &ast.ProcDeclChildOutput{
									Children: []*ast.FieldOrComment{
										{
											Field: &ast.Field{
												Name: "fieldOut",
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{Named: testutil.Pointer("int")},
												},
											},
										},
									},
								},
							},
							{
								Comment: &ast.Comment{Simple: testutil.Pointer(" Between output and meta")},
							},
							{
								Meta: &ast.ProcDeclChildMeta{
									Children: []*ast.ProcDeclChildMetaChild{
										{
											KV: &ast.ProcDeclChildMetaKV{
												Key:   "key",
												Value: ast.AnyLiteral{Str: testutil.Pointer("value")},
											},
										},
									},
								},
							},
							{
								Comment: &ast.Comment{Simple: testutil.Pointer(" Trailing comment in proc")},
							},
						},
					},
				},
			},
		}
		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Comments within ProcDecl Input block", func(t *testing.T) {
		input := `
			proc MyProc {
				input {
					// Before fieldIn1
					fieldIn1: string
					/* Between fieldIn1 and fieldIn2 */
					fieldIn2: int
					// Trailing comment in input
				}
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Proc: &ast.ProcDecl{
						Name: "MyProc",
						Children: []*ast.ProcDeclChild{
							{
								Input: &ast.ProcDeclChildInput{
									Children: []*ast.FieldOrComment{
										{
											Comment: &ast.Comment{Simple: testutil.Pointer(" Before fieldIn1")},
										},
										{
											Field: &ast.Field{
												Name: "fieldIn1",
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")},
												},
											},
										},
										{
											Comment: &ast.Comment{Block: testutil.Pointer(" Between fieldIn1 and fieldIn2 ")},
										},
										{
											Field: &ast.Field{
												Name: "fieldIn2",
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{Named: testutil.Pointer("int")},
												},
											},
										},
										{
											Comment: &ast.Comment{Simple: testutil.Pointer(" Trailing comment in input")},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Comments within ProcDecl Output block", func(t *testing.T) {
		input := `
			proc MyProc {
				output {
					// Before fieldOut1
					fieldOut1: string
					/* Between fieldOut1 and fieldOut2 */
					fieldOut2: int
					// Trailing comment in output
				}
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Proc: &ast.ProcDecl{
						Name: "MyProc",
						Children: []*ast.ProcDeclChild{
							{
								Output: &ast.ProcDeclChildOutput{
									Children: []*ast.FieldOrComment{
										{
											Comment: &ast.Comment{Simple: testutil.Pointer(" Before fieldOut1")},
										},
										{
											Field: &ast.Field{
												Name: "fieldOut1",
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")},
												},
											},
										},
										{
											Comment: &ast.Comment{Block: testutil.Pointer(" Between fieldOut1 and fieldOut2 ")},
										},
										{
											Field: &ast.Field{
												Name: "fieldOut2",
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{Named: testutil.Pointer("int")},
												},
											},
										},
										{
											Comment: &ast.Comment{Simple: testutil.Pointer(" Trailing comment in output")},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Comments within ProcDecl Meta block", func(t *testing.T) {
		input := `
			proc MyProc {
				meta {
					// Before key1
					key1: "value1"
					/* Between key1 and key2 */
					key2: 123
					// Trailing comment in meta
				}
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Proc: &ast.ProcDecl{
						Name: "MyProc",
						Children: []*ast.ProcDeclChild{
							{
								Meta: &ast.ProcDeclChildMeta{

									Children: []*ast.ProcDeclChildMetaChild{
										{
											Comment: &ast.Comment{Simple: testutil.Pointer(" Before key1")},
										},
										{
											KV: &ast.ProcDeclChildMetaKV{
												Key:   "key1",
												Value: ast.AnyLiteral{Str: testutil.Pointer("value1")},
											},
										},
										{
											Comment: &ast.Comment{Block: testutil.Pointer(" Between key1 and key2 ")},
										},
										{
											KV: &ast.ProcDeclChildMetaKV{
												Key:   "key2",
												Value: ast.AnyLiteral{Int: testutil.Pointer("123")},
											},
										},
										{
											Comment: &ast.Comment{Simple: testutil.Pointer(" Trailing comment in meta")},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Comments within FieldTypeObject (nested type)", func(t *testing.T) {
		input := `
			type MyType {
				nested: {
					// Before sub1
					sub1: string
					/* Between sub1 and sub2 */
					sub2: int
					// Trailing comment in nested
				}
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Type: &ast.TypeDecl{
						Name: "MyType",
						Children: []*ast.FieldOrComment{
							{
								Field: &ast.Field{
									Name: "nested",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{
											Object: &ast.FieldTypeObject{
												Children: []*ast.FieldOrComment{
													{
														Comment: &ast.Comment{Simple: testutil.Pointer(" Before sub1")},
													},
													{
														Field: &ast.Field{
															Name: "sub1",
															Type: ast.FieldType{
																Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")},
															},
														},
													},
													{
														Comment: &ast.Comment{Block: testutil.Pointer(" Between sub1 and sub2 ")},
													},
													{
														Field: &ast.Field{
															Name: "sub2",
															Type: ast.FieldType{
																Base: &ast.FieldTypeBase{Named: testutil.Pointer("int")},
															},
														},
													},
													{
														Comment: &ast.Comment{Simple: testutil.Pointer(" Trailing comment in nested")},
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
		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Comments between Field rules (are part of the field children)", func(t *testing.T) {
		input := `
			type MyType {
				field: string
					// Before rule1
					@rule1
					/* Between rule1 and rule2 */
					@rule2("param")

				// Trailing comment after rules
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Type: &ast.TypeDecl{
						Name: "MyType",
						Children: []*ast.FieldOrComment{
							{
								Field: &ast.Field{
									Name: "field",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")},
									},
									Children: []*ast.FieldChild{
										{
											Comment: &ast.Comment{Simple: testutil.Pointer(" Before rule1")},
										},
										{
											Rule: &ast.FieldRule{Name: "rule1"},
										},
										{
											Comment: &ast.Comment{Block: testutil.Pointer(" Between rule1 and rule2 ")},
										},
										{
											Rule: &ast.FieldRule{
												Name: "rule2",
												Body: &ast.FieldRuleBody{
													ParamSingle: &ast.AnyLiteral{Str: testutil.Pointer("param")},
												},
											},
										},
									},
								},
							},
							{
								Comment: &ast.Comment{Simple: testutil.Pointer(" Trailing comment after rules")},
							},
						},
					},
				},
			},
		}
		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("End-of-line comments", func(t *testing.T) {
		input := `
			version 1 // EOL on version
			import "path" // EOL on import
			rule @myRule { // EOL on rule start
				for: string // EOL on for
				param: int // EOL on param
			} // EOL on rule end
			type MyType { // EOL on type start
				field: string // EOL on field
					@rule1 // EOL on rule
			} // EOL on type end
			proc MyProc { // EOL on proc start
				input { f: int } // EOL on input
				output { o: int } // EOL on output
				meta { k: "v" } // EOL on meta
			} // EOL on proc end
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Version: &ast.Version{Number: 1},
				},
				{
					Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on version")},
				},
				{
					Import: &ast.Import{Path: "path"},
				},
				{
					Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on import")},
				},
				{
					Rule: &ast.RuleDecl{
						Name: "myRule",
						Children: []*ast.RuleDeclChild{
							{Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on rule start")}}, // Comment inside the block
							{For: &ast.RuleDeclChildFor{For: "string"}},
							{Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on for")}},
							{Param: &ast.RuleDeclChildParam{Param: "int"}},
							{Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on param")}},
						},
					},
				},
				{
					Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on rule end")},
				}, // Comment after the block
				{
					Type: &ast.TypeDecl{
						Name: "MyType",
						Children: []*ast.FieldOrComment{
							{Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on type start")}}, // Comment inside the block
							{
								Field: &ast.Field{
									Name: "field",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")},
									},
									Children: []*ast.FieldChild{
										{Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on field")}}, // Comment after type, before rule
										{Rule: &ast.FieldRule{Name: "rule1"}},
									},
								},
							},
							{Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on rule")}}, // Comment after rule
						},
					},
				},
				{
					Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on type end")},
				}, // Comment after the block
				{
					Proc: &ast.ProcDecl{
						Name: "MyProc",
						Children: []*ast.ProcDeclChild{
							{Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on proc start")}}, // Comment inside the block
							{
								Input: &ast.ProcDeclChildInput{
									Children: []*ast.FieldOrComment{
										{
											Field: &ast.Field{
												Name: "f",
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{Named: testutil.Pointer("int")},
												},
											},
										},
									},
								},
							},
							{Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on input")}},
							{
								Output: &ast.ProcDeclChildOutput{
									Children: []*ast.FieldOrComment{
										{
											Field: &ast.Field{
												Name: "o",
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{Named: testutil.Pointer("int")},
												},
											},
										},
									},
								},
							},
							{Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on output")}},
							{
								Meta: &ast.ProcDeclChildMeta{
									Children: []*ast.ProcDeclChildMetaChild{
										{
											KV: &ast.ProcDeclChildMetaKV{
												Key:   "k",
												Value: ast.AnyLiteral{Str: testutil.Pointer("v")},
											},
										},
									},
								},
							},
							{Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on meta")}},
						},
					},
				},
				{
					Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on proc end")},
				}, // Comment after the block
			},
		}
		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Comments inside empty blocks", func(t *testing.T) {
		input := `
			rule @emptyRule { /* Rule Comment */ }
			type EmptyType { // Type Comment
			}
			proc EmptyProc {
				/* Proc Comment */
				input { /* Input Comment */ }
				output { // Output Comment
				}
				meta { /* Meta Comment */ }
			}
			type NestedEmpty {
				field: { /* Nested Comment */ }
			}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Rule: &ast.RuleDecl{
						Name: "emptyRule",
						Children: []*ast.RuleDeclChild{
							{
								Comment: &ast.Comment{Block: testutil.Pointer(" Rule Comment ")},
							},
						},
					},
				},
				{
					Type: &ast.TypeDecl{
						Name: "EmptyType",
						Children: []*ast.FieldOrComment{
							{
								Comment: &ast.Comment{Simple: testutil.Pointer(" Type Comment")},
							},
						},
					},
				},
				{
					Proc: &ast.ProcDecl{
						Name: "EmptyProc",
						Children: []*ast.ProcDeclChild{
							{
								Comment: &ast.Comment{Block: testutil.Pointer(" Proc Comment ")},
							},
							{
								Input: &ast.ProcDeclChildInput{
									Children: []*ast.FieldOrComment{
										{
											Comment: &ast.Comment{Block: testutil.Pointer(" Input Comment ")},
										},
									},
								},
							},
							{
								Output: &ast.ProcDeclChildOutput{
									Children: []*ast.FieldOrComment{
										{
											Comment: &ast.Comment{Simple: testutil.Pointer(" Output Comment")},
										},
									},
								},
							},
							{
								Meta: &ast.ProcDeclChildMeta{
									Children: []*ast.ProcDeclChildMetaChild{
										{
											Comment: &ast.Comment{Block: testutil.Pointer(" Meta Comment ")},
										},
									},
								},
							},
						},
					},
				},
				{
					Type: &ast.TypeDecl{
						Name: "NestedEmpty",
						Children: []*ast.FieldOrComment{
							{
								Field: &ast.Field{
									Name: "field",
									Type: ast.FieldType{
										Base: &ast.FieldTypeBase{
											Object: &ast.FieldTypeObject{
												Children: []*ast.FieldOrComment{
													{
														Comment: &ast.Comment{Block: testutil.Pointer(" Nested Comment ")},
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
		testutil.ASTEqualNoPos(t, expected, parsed)
	})
}

func TestParserDocstrings(t *testing.T) {
	t.Run("Standalone docstrings", func(t *testing.T) {
		input := `
			""" This is a standalone docstring. """
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Docstring: &ast.Docstring{Value: " This is a standalone docstring. "},
				},
			},
		}
		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Multiple standalone docstrings", func(t *testing.T) {
		input := `
			""" This is a standalone docstring. """
			""" This is a standalone docstring. """
			""" This is a standalone docstring. """
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{Docstring: &ast.Docstring{Value: " This is a standalone docstring. "}},
				{Docstring: &ast.Docstring{Value: " This is a standalone docstring. "}},
				{Docstring: &ast.Docstring{Value: " This is a standalone docstring. "}},
			},
		}
		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Standalone docstrings and associated docstrings", func(t *testing.T) {
		input := `
			""" This is a standalone docstring. """
			""" This is a standalone docstring. """
			""" This is a standalone docstring. """
			""" This is an associated docstring. """
			type MyType {}
			""" This is a standalone docstring. """
			""" This is a standalone docstring. """
			""" This is a standalone docstring. """""" This is an associated docstring. """type MyType {}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{Docstring: &ast.Docstring{Value: " This is a standalone docstring. "}},
				{Docstring: &ast.Docstring{Value: " This is a standalone docstring. "}},
				{Docstring: &ast.Docstring{Value: " This is a standalone docstring. "}},
				{
					Type: &ast.TypeDecl{
						Docstring: &ast.Docstring{Value: " This is an associated docstring. "},
						Name:      "MyType",
					},
				},
				{Docstring: &ast.Docstring{Value: " This is a standalone docstring. "}},
				{Docstring: &ast.Docstring{Value: " This is a standalone docstring. "}},
				{Docstring: &ast.Docstring{Value: " This is a standalone docstring. "}},
				{
					Type: &ast.TypeDecl{
						Docstring: &ast.Docstring{Value: " This is an associated docstring. "},
						Name:      "MyType",
					},
				},
			},
		}
		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Standalone docstrings should not associate if there is a blank line", func(t *testing.T) {
		input := `
			""" This is a standalone docstring. """
			""" This is a standalone docstring. """
			""" This is a standalone docstring. """

			type MyType {}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{Docstring: &ast.Docstring{Value: " This is a standalone docstring. "}},
				{Docstring: &ast.Docstring{Value: " This is a standalone docstring. "}},
				{Docstring: &ast.Docstring{Value: " This is a standalone docstring. "}},
				{Type: &ast.TypeDecl{Name: "MyType"}},
			},
		}
		testutil.ASTEqualNoPos(t, expected, parsed)
	})
}

func TestParserFullSchema(t *testing.T) {
	input := `
		version 1

		/* Import other schema */
		import "my_sub_schema.urpc"

		// Custom rule declarations

		"""
		This rule validates if a string matches a regular expression pattern.
		Useful for emails, URLs, and other formatted strings.
		"""
		deprecated rule @regex {
			for: string
			param: string
			error: "Invalid format"
		}

		""" Validate "category" with custom logic """
		deprecated("Deprecated")
		rule @validateCategory {
			for: Category
			error: "Field \"category\" is not valid"
		}

		// Type declarations

		type FirstDummyType {
			dummyField: datetime
				@min("1900-01-01T00:00:00Z")
				@max("2100-01-01T00:00:00Z")
		}

		type SecondDummyType {
			dummyField: int
		}

		deprecated type ThirdDummyType {
			dummyField: string
		}

		"""
		Category represents a product category in the system.
		This type is used across the catalog module.
		"""
		deprecated("Deprecated")
		type Category {
			id: string
				@uuid(error: "Must be a valid UUID")
				@minlen(36)
				@maxlen(36, error: "UUID must be exactly 36 characters")
			name: string
				@minlen(3, error: "Name must be at least 3 characters long")
			description?: string
			isActive: boolean
				@equals(true)
			parentId?: string
				@uuid
		}

		"""
		Product represents a sellable item in the store.
		Products have complex validation rules and can be
		nested inside catalogs.
		"""
		type Product {
			id: string
				@uuid
			name: string
				@minlen(2)
				@maxlen(100, error: "Name cannot exceed 100 characters")
			price: float
				@min(0.01, error: "Price must be greater than zero")
			stock: int
				@min(0)
				@range([0, 1000], error: "Stock must be between 0 and 1000")
			category: Category
				@validateCategory(error: "Invalid category custom message")
			tags?: string[]
				@minlen(1, error: "At least one tag is required")
				@maxlen(10)

			details: {
				dimensions: {
					width: float
						@min(0.0, error: "Width cannot be negative")
					height: float
						@min(0.0)
					depth?: float
				}
				weight?: float
				colors: string[]
					@enum(["red", "green", "blue", "black", "white"], error: "Color must be one of the allowed values")
				attributes?: {
					name: string
					value: string
				}[]
			}

			variations: {
				sku: string
				price: float
					@min(0.01, error: "Variation price must be greater than zero")
				attributes: {
					name: string
					value: string
				}[]
			}[]
		}

		// Procedure declarations

		"""
		GetCategory retrieves a category by its ID.
		This is a basic read operation.
		"""
		deprecated proc GetCategory {
			input {
				id: string
					@uuid(error: "Category ID must be a valid UUID")
			}

			output {
				category: Category
				exists: boolean
			}

			meta {
				cache: true
				cacheTime: 300
				requiresAuth: false
				apiVersion: "1.0.0"
			}
		}

		"""
		CreateProduct adds a new product to the catalog.
		This procedure handles complex validation and returns
		detailed success information.
		"""
		deprecated("Deprecated")
		proc CreateProduct {
			input {
				product: Product
				options?: {
					draft: boolean
					notify: boolean
					scheduledFor?: string
						@iso8601(error: "Must be a valid ISO8601 date")
					tags?: string[]
				}

				validation: {
					skipValidation?: boolean
					customRules?: {
						name: string
						severity: int
							@enum([1, 2, 3], error: "Severity must be 1, 2, or 3")
						message: string
					}[]
				}
			}

			output {
				success: boolean
				productId: string
					@uuid(error: "Product ID must be a valid UUID")
				errors?: {
					code: int
					message: string
					field?: string
				}[]

				analytics: {
					duration: float
					processingSteps: {
						name: string
						duration: float
						success: boolean
					}[]
					serverInfo: {
						id: string
						region: string
						load: float
							@min(0.0)
							@max(1.0, error: "Load factor cannot exceed 1.0")
					}
				}
			}

			meta {
				auth: "required"
				roles: "admin,product-manager"
				rateLimit: 100
				timeout: 30.5
				audit: true
				apiVersion: "1.2.0"
			}
		}

		"""
		Validates if a value is within a specified range.
		"""
		rule @range {
			for: int
			param: int[]
			error: "Value out of range"
		}
	`

	parsed, err := ParserInstance.ParseString("schema.urpc", input)
	require.NoError(t, err)

	expected := &ast.Schema{
		Children: []*ast.SchemaChild{
			{
				Version: &ast.Version{
					Number: 1,
				},
			},
			{
				Comment: &ast.Comment{
					Block: testutil.Pointer(" Import other schema "),
				},
			},
			{
				Import: &ast.Import{
					Path: "my_sub_schema.urpc",
				},
			},
			{
				Comment: &ast.Comment{
					Simple: testutil.Pointer(" Custom rule declarations"),
				},
			},
			{
				Rule: &ast.RuleDecl{
					Docstring: &ast.Docstring{
						Value: "\n\t\tThis rule validates if a string matches a regular expression pattern.\n\t\tUseful for emails, URLs, and other formatted strings.\n\t\t",
					},
					Deprecated: &ast.Deprecated{},
					Name:       "regex",
					Children: []*ast.RuleDeclChild{
						{
							For: &ast.RuleDeclChildFor{
								For: "string",
							},
						},
						{
							Param: &ast.RuleDeclChildParam{
								Param: "string",
							},
						},
						{
							Error: &ast.RuleDeclChildError{
								Error: "Invalid format",
							},
						},
					},
				},
			},
			{
				Rule: &ast.RuleDecl{
					Docstring: &ast.Docstring{
						Value: " Validate \"category\" with custom logic ",
					},
					Deprecated: &ast.Deprecated{
						Message: testutil.Pointer("Deprecated"),
					},
					Name: "validateCategory",
					Children: []*ast.RuleDeclChild{
						{
							For: &ast.RuleDeclChildFor{
								For: "Category",
							},
						},
						{
							Error: &ast.RuleDeclChildError{
								Error: "Field \"category\" is not valid",
							},
						},
					},
				},
			},
			{
				Comment: &ast.Comment{
					Simple: testutil.Pointer(" Type declarations"),
				},
			},
			{
				Type: &ast.TypeDecl{
					Name: "FirstDummyType",
					Children: []*ast.FieldOrComment{
						{
							Field: &ast.Field{
								Name: "dummyField",
								Type: ast.FieldType{
									Base: &ast.FieldTypeBase{
										Named: testutil.Pointer("datetime"),
									},
								},
								Children: []*ast.FieldChild{
									{
										Rule: &ast.FieldRule{
											Name: "min",
											Body: &ast.FieldRuleBody{
												ParamSingle: &ast.AnyLiteral{Str: testutil.Pointer("1900-01-01T00:00:00Z")},
											},
										},
									},
									{
										Rule: &ast.FieldRule{
											Name: "max",
											Body: &ast.FieldRuleBody{
												ParamSingle: &ast.AnyLiteral{Str: testutil.Pointer("2100-01-01T00:00:00Z")},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				Type: &ast.TypeDecl{
					Name: "SecondDummyType",
					Children: []*ast.FieldOrComment{
						{
							Field: &ast.Field{
								Name: "dummyField",
								Type: ast.FieldType{
									Base: &ast.FieldTypeBase{
										Named: testutil.Pointer("int"),
									},
								},
							},
						},
					},
				},
			},
			{
				Type: &ast.TypeDecl{
					Deprecated: &ast.Deprecated{},
					Name:       "ThirdDummyType",
					Children: []*ast.FieldOrComment{
						{
							Field: &ast.Field{
								Name: "dummyField",
								Type: ast.FieldType{
									Base: &ast.FieldTypeBase{
										Named: testutil.Pointer("string"),
									},
								},
							},
						},
					},
				},
			},
			{
				Type: &ast.TypeDecl{
					Docstring: &ast.Docstring{
						Value: "\n\t\tCategory represents a product category in the system.\n\t\tThis type is used across the catalog module.\n\t\t",
					},
					Deprecated: &ast.Deprecated{
						Message: testutil.Pointer("Deprecated"),
					},
					Name: "Category",
					Children: []*ast.FieldOrComment{
						{
							Field: &ast.Field{
								Name: "id",
								Type: ast.FieldType{
									Base: &ast.FieldTypeBase{
										Named: testutil.Pointer("string"),
									},
								},
								Children: []*ast.FieldChild{
									{
										Rule: &ast.FieldRule{
											Name: "uuid",
											Body: &ast.FieldRuleBody{
												Error: testutil.Pointer("Must be a valid UUID"),
											},
										},
									},
									{
										Rule: &ast.FieldRule{
											Name: "minlen",
											Body: &ast.FieldRuleBody{
												ParamSingle: &ast.AnyLiteral{Int: testutil.Pointer("36")},
											},
										},
									},
									{
										Rule: &ast.FieldRule{
											Name: "maxlen",
											Body: &ast.FieldRuleBody{
												ParamSingle: &ast.AnyLiteral{Int: testutil.Pointer("36")},
												Error:       testutil.Pointer("UUID must be exactly 36 characters"),
											},
										},
									},
								},
							},
						},
						{
							Field: &ast.Field{
								Name: "name",
								Type: ast.FieldType{
									Base: &ast.FieldTypeBase{
										Named: testutil.Pointer("string"),
									},
								},
								Children: []*ast.FieldChild{
									{
										Rule: &ast.FieldRule{
											Name: "minlen",
											Body: &ast.FieldRuleBody{
												ParamSingle: &ast.AnyLiteral{Int: testutil.Pointer("3")},
												Error:       testutil.Pointer("Name must be at least 3 characters long"),
											},
										},
									},
								},
							},
						},
						{
							Field: &ast.Field{
								Name:     "description",
								Optional: true,
								Type: ast.FieldType{
									Base: &ast.FieldTypeBase{
										Named: testutil.Pointer("string"),
									},
								},
							},
						},
						{
							Field: &ast.Field{
								Name: "isActive",
								Type: ast.FieldType{
									Base: &ast.FieldTypeBase{
										Named: testutil.Pointer("boolean"),
									},
								},
								Children: []*ast.FieldChild{
									{
										Rule: &ast.FieldRule{
											Name: "equals",
											Body: &ast.FieldRuleBody{
												ParamSingle: &ast.AnyLiteral{True: testutil.Pointer("true")},
											},
										},
									},
								},
							},
						},
						{
							Field: &ast.Field{
								Name:     "parentId",
								Optional: true,
								Type: ast.FieldType{
									Base: &ast.FieldTypeBase{
										Named: testutil.Pointer("string"),
									},
								},
								Children: []*ast.FieldChild{
									{
										Rule: &ast.FieldRule{
											Name: "uuid",
										},
									},
								},
							},
						},
					},
				},
			},
			{
				Type: &ast.TypeDecl{
					Docstring: &ast.Docstring{
						Value: "\n\t\tProduct represents a sellable item in the store.\n\t\tProducts have complex validation rules and can be\n\t\tnested inside catalogs.\n\t\t",
					},
					Name: "Product",
					Children: []*ast.FieldOrComment{
						{
							Field: &ast.Field{
								Name: "id",
								Type: ast.FieldType{
									Base: &ast.FieldTypeBase{
										Named: testutil.Pointer("string"),
									},
								},
								Children: []*ast.FieldChild{
									{
										Rule: &ast.FieldRule{
											Name: "uuid",
										},
									},
								},
							},
						},
						{
							Field: &ast.Field{
								Name: "name",
								Type: ast.FieldType{
									Base: &ast.FieldTypeBase{
										Named: testutil.Pointer("string"),
									},
								},
								Children: []*ast.FieldChild{
									{
										Rule: &ast.FieldRule{
											Name: "minlen",
											Body: &ast.FieldRuleBody{
												ParamSingle: &ast.AnyLiteral{Int: testutil.Pointer("2")},
											},
										},
									},
									{
										Rule: &ast.FieldRule{
											Name: "maxlen",
											Body: &ast.FieldRuleBody{
												ParamSingle: &ast.AnyLiteral{Int: testutil.Pointer("100")},
												Error:       testutil.Pointer("Name cannot exceed 100 characters"),
											},
										},
									},
								},
							},
						},
						{
							Field: &ast.Field{
								Name: "price",
								Type: ast.FieldType{
									Base: &ast.FieldTypeBase{
										Named: testutil.Pointer("float"),
									},
								},
								Children: []*ast.FieldChild{
									{
										Rule: &ast.FieldRule{
											Name: "min",
											Body: &ast.FieldRuleBody{
												ParamSingle: &ast.AnyLiteral{Float: testutil.Pointer("0.01")},
												Error:       testutil.Pointer("Price must be greater than zero"),
											},
										},
									},
								},
							},
						},
						{
							Field: &ast.Field{
								Name: "stock",
								Type: ast.FieldType{
									Base: &ast.FieldTypeBase{
										Named: testutil.Pointer("int"),
									},
								},
								Children: []*ast.FieldChild{
									{
										Rule: &ast.FieldRule{
											Name: "min",
											Body: &ast.FieldRuleBody{
												ParamSingle: &ast.AnyLiteral{Int: testutil.Pointer("0")},
											},
										},
									},
									{
										Rule: &ast.FieldRule{
											Name: "range",
											Body: &ast.FieldRuleBody{
												ParamListInt: []string{"0", "1000"},
												Error:        testutil.Pointer("Stock must be between 0 and 1000"),
											},
										},
									},
								},
							},
						},
						{
							Field: &ast.Field{
								Name: "category",
								Type: ast.FieldType{
									Base: &ast.FieldTypeBase{
										Named: testutil.Pointer("Category"),
									},
								},
								Children: []*ast.FieldChild{
									{
										Rule: &ast.FieldRule{
											Name: "validateCategory",
											Body: &ast.FieldRuleBody{
												Error: testutil.Pointer("Invalid category custom message"),
											},
										},
									},
								},
							},
						},
						{
							Field: &ast.Field{
								Name:     "tags",
								Optional: true,
								Type: ast.FieldType{
									IsArray: true,
									Base: &ast.FieldTypeBase{
										Named: testutil.Pointer("string"),
									},
								},
								Children: []*ast.FieldChild{
									{
										Rule: &ast.FieldRule{
											Name: "minlen",
											Body: &ast.FieldRuleBody{
												ParamSingle: &ast.AnyLiteral{Int: testutil.Pointer("1")},
												Error:       testutil.Pointer("At least one tag is required"),
											},
										},
									},
									{
										Rule: &ast.FieldRule{
											Name: "maxlen",
											Body: &ast.FieldRuleBody{
												ParamSingle: &ast.AnyLiteral{Int: testutil.Pointer("10")},
											},
										},
									},
								},
							},
						},
						{
							Field: &ast.Field{
								Name: "details",
								Type: ast.FieldType{
									Base: &ast.FieldTypeBase{
										Object: &ast.FieldTypeObject{
											Children: []*ast.FieldOrComment{
												{
													Field: &ast.Field{
														Name: "dimensions",
														Type: ast.FieldType{
															Base: &ast.FieldTypeBase{
																Object: &ast.FieldTypeObject{
																	Children: []*ast.FieldOrComment{
																		{
																			Field: &ast.Field{
																				Name: "width",
																				Type: ast.FieldType{
																					Base: &ast.FieldTypeBase{
																						Named: testutil.Pointer("float"),
																					},
																				},
																				Children: []*ast.FieldChild{
																					{
																						Rule: &ast.FieldRule{
																							Name: "min",
																							Body: &ast.FieldRuleBody{
																								ParamSingle: &ast.AnyLiteral{Float: testutil.Pointer("0.0")},
																								Error:       testutil.Pointer("Width cannot be negative"),
																							},
																						},
																					},
																				},
																			},
																		},
																		{
																			Field: &ast.Field{
																				Name: "height",
																				Type: ast.FieldType{
																					Base: &ast.FieldTypeBase{
																						Named: testutil.Pointer("float"),
																					},
																				},
																				Children: []*ast.FieldChild{
																					{
																						Rule: &ast.FieldRule{
																							Name: "min",
																							Body: &ast.FieldRuleBody{
																								ParamSingle: &ast.AnyLiteral{Float: testutil.Pointer("0.0")},
																							},
																						},
																					},
																				},
																			},
																		},
																		{
																			Field: &ast.Field{
																				Name:     "depth",
																				Optional: true,
																				Type: ast.FieldType{
																					Base: &ast.FieldTypeBase{
																						Named: testutil.Pointer("float"),
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
												{
													Field: &ast.Field{
														Name:     "weight",
														Optional: true,
														Type: ast.FieldType{
															Base: &ast.FieldTypeBase{
																Named: testutil.Pointer("float"),
															},
														},
													},
												},
												{
													Field: &ast.Field{
														Name: "colors",
														Type: ast.FieldType{
															IsArray: true,
															Base: &ast.FieldTypeBase{
																Named: testutil.Pointer("string"),
															},
														},
														Children: []*ast.FieldChild{
															{
																Rule: &ast.FieldRule{
																	Name: "enum",
																	Body: &ast.FieldRuleBody{
																		ParamListString: []string{
																			"red",
																			"green",
																			"blue",
																			"black",
																			"white",
																		},
																		Error: testutil.Pointer("Color must be one of the allowed values"),
																	},
																},
															},
														},
													},
												},
												{
													Field: &ast.Field{
														Name:     "attributes",
														Optional: true,
														Type: ast.FieldType{
															IsArray: true,
															Base: &ast.FieldTypeBase{
																Object: &ast.FieldTypeObject{
																	Children: []*ast.FieldOrComment{
																		{
																			Field: &ast.Field{
																				Name: "name",
																				Type: ast.FieldType{
																					Base: &ast.FieldTypeBase{
																						Named: testutil.Pointer("string"),
																					},
																				},
																			},
																		},
																		{
																			Field: &ast.Field{
																				Name: "value",
																				Type: ast.FieldType{
																					Base: &ast.FieldTypeBase{
																						Named: testutil.Pointer("string"),
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
						{
							Field: &ast.Field{
								Name: "variations",
								Type: ast.FieldType{
									IsArray: true,
									Base: &ast.FieldTypeBase{
										Object: &ast.FieldTypeObject{
											Children: []*ast.FieldOrComment{
												{
													Field: &ast.Field{
														Name: "sku",
														Type: ast.FieldType{
															Base: &ast.FieldTypeBase{
																Named: testutil.Pointer("string"),
															},
														},
													},
												},
												{
													Field: &ast.Field{
														Name: "price",
														Type: ast.FieldType{
															Base: &ast.FieldTypeBase{
																Named: testutil.Pointer("float"),
															},
														},
														Children: []*ast.FieldChild{
															{
																Rule: &ast.FieldRule{
																	Name: "min",
																	Body: &ast.FieldRuleBody{
																		ParamSingle: &ast.AnyLiteral{Float: testutil.Pointer("0.01")},
																		Error:       testutil.Pointer("Variation price must be greater than zero"),
																	},
																},
															},
														},
													},
												},
												{
													Field: &ast.Field{
														Name: "attributes",
														Type: ast.FieldType{
															IsArray: true,
															Base: &ast.FieldTypeBase{
																Object: &ast.FieldTypeObject{
																	Children: []*ast.FieldOrComment{
																		{
																			Field: &ast.Field{
																				Name: "name",
																				Type: ast.FieldType{
																					Base: &ast.FieldTypeBase{
																						Named: testutil.Pointer("string"),
																					},
																				},
																			},
																		},
																		{
																			Field: &ast.Field{
																				Name: "value",
																				Type: ast.FieldType{
																					Base: &ast.FieldTypeBase{
																						Named: testutil.Pointer("string"),
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
					},
				},
			},
			{
				Comment: &ast.Comment{
					Simple: testutil.Pointer(" Procedure declarations"),
				},
			},
			{
				Proc: &ast.ProcDecl{
					Docstring: &ast.Docstring{
						Value: "\n\t\tGetCategory retrieves a category by its ID.\n\t\tThis is a basic read operation.\n\t\t",
					},
					Deprecated: &ast.Deprecated{},
					Name:       "GetCategory",
					Children: []*ast.ProcDeclChild{
						{
							Input: &ast.ProcDeclChildInput{
								Children: []*ast.FieldOrComment{
									{
										Field: &ast.Field{
											Name: "id",
											Type: ast.FieldType{
												Base: &ast.FieldTypeBase{
													Named: testutil.Pointer("string"),
												},
											},
											Children: []*ast.FieldChild{
												{
													Rule: &ast.FieldRule{
														Name: "uuid",
														Body: &ast.FieldRuleBody{
															Error: testutil.Pointer("Category ID must be a valid UUID"),
														},
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Output: &ast.ProcDeclChildOutput{
								Children: []*ast.FieldOrComment{
									{
										Field: &ast.Field{
											Name: "category",
											Type: ast.FieldType{
												Base: &ast.FieldTypeBase{
													Named: testutil.Pointer("Category"),
												},
											},
										},
									},
									{
										Field: &ast.Field{
											Name: "exists",
											Type: ast.FieldType{
												Base: &ast.FieldTypeBase{
													Named: testutil.Pointer("boolean"),
												},
											},
										},
									},
								},
							},
						},
						{
							Meta: &ast.ProcDeclChildMeta{
								Children: []*ast.ProcDeclChildMetaChild{
									{
										KV: &ast.ProcDeclChildMetaKV{
											Key:   "cache",
											Value: ast.AnyLiteral{True: testutil.Pointer("true")},
										},
									},
									{
										KV: &ast.ProcDeclChildMetaKV{
											Key:   "cacheTime",
											Value: ast.AnyLiteral{Int: testutil.Pointer("300")},
										},
									},
									{
										KV: &ast.ProcDeclChildMetaKV{
											Key:   "requiresAuth",
											Value: ast.AnyLiteral{False: testutil.Pointer("false")},
										},
									},
									{
										KV: &ast.ProcDeclChildMetaKV{
											Key:   "apiVersion",
											Value: ast.AnyLiteral{Str: testutil.Pointer("1.0.0")},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				Proc: &ast.ProcDecl{
					Docstring: &ast.Docstring{
						Value: "\n\t\tCreateProduct adds a new product to the catalog.\n\t\tThis procedure handles complex validation and returns\n\t\tdetailed success information.\n\t\t",
					},
					Deprecated: &ast.Deprecated{
						Message: testutil.Pointer("Deprecated"),
					},
					Name: "CreateProduct",
					Children: []*ast.ProcDeclChild{
						{
							Input: &ast.ProcDeclChildInput{
								Children: []*ast.FieldOrComment{
									{
										Field: &ast.Field{
											Name: "product",
											Type: ast.FieldType{
												Base: &ast.FieldTypeBase{
													Named: testutil.Pointer("Product"),
												},
											},
										},
									},
									{
										Field: &ast.Field{
											Name:     "options",
											Optional: true,
											Type: ast.FieldType{
												Base: &ast.FieldTypeBase{
													Object: &ast.FieldTypeObject{
														Children: []*ast.FieldOrComment{
															{
																Field: &ast.Field{
																	Name: "draft",
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: testutil.Pointer("boolean"),
																		},
																	},
																},
															},
															{
																Field: &ast.Field{
																	Name: "notify",
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: testutil.Pointer("boolean"),
																		},
																	},
																},
															},
															{
																Field: &ast.Field{
																	Name:     "scheduledFor",
																	Optional: true,
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: testutil.Pointer("string"),
																		},
																	},
																	Children: []*ast.FieldChild{
																		{
																			Rule: &ast.FieldRule{
																				Name: "iso8601",
																				Body: &ast.FieldRuleBody{
																					Error: testutil.Pointer("Must be a valid ISO8601 date"),
																				},
																			},
																		},
																	},
																},
															},
															{
																Field: &ast.Field{
																	Name:     "tags",
																	Optional: true,
																	Type: ast.FieldType{
																		IsArray: true,
																		Base: &ast.FieldTypeBase{
																			Named: testutil.Pointer("string"),
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
									{
										Field: &ast.Field{
											Name: "validation",
											Type: ast.FieldType{
												Base: &ast.FieldTypeBase{
													Object: &ast.FieldTypeObject{
														Children: []*ast.FieldOrComment{
															{
																Field: &ast.Field{
																	Name:     "skipValidation",
																	Optional: true,
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: testutil.Pointer("boolean"),
																		},
																	},
																},
															},
															{
																Field: &ast.Field{
																	Name:     "customRules",
																	Optional: true,
																	Type: ast.FieldType{
																		IsArray: true,
																		Base: &ast.FieldTypeBase{
																			Object: &ast.FieldTypeObject{
																				Children: []*ast.FieldOrComment{
																					{
																						Field: &ast.Field{
																							Name: "name",
																							Type: ast.FieldType{
																								Base: &ast.FieldTypeBase{
																									Named: testutil.Pointer("string"),
																								},
																							},
																						},
																					},
																					{
																						Field: &ast.Field{
																							Name: "severity",
																							Type: ast.FieldType{
																								Base: &ast.FieldTypeBase{
																									Named: testutil.Pointer("int"),
																								},
																							},
																							Children: []*ast.FieldChild{
																								{
																									Rule: &ast.FieldRule{
																										Name: "enum",
																										Body: &ast.FieldRuleBody{
																											ParamListInt: []string{
																												"1",
																												"2",
																												"3",
																											},
																											Error: testutil.Pointer("Severity must be 1, 2, or 3"),
																										},
																									},
																								},
																							},
																						},
																					},
																					{
																						Field: &ast.Field{
																							Name: "message",
																							Type: ast.FieldType{
																								Base: &ast.FieldTypeBase{
																									Named: testutil.Pointer("string"),
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
								},
							},
						},
						{
							Output: &ast.ProcDeclChildOutput{
								Children: []*ast.FieldOrComment{
									{
										Field: &ast.Field{
											Name: "success",
											Type: ast.FieldType{
												Base: &ast.FieldTypeBase{
													Named: testutil.Pointer("boolean"),
												},
											},
										},
									},
									{
										Field: &ast.Field{
											Name: "productId",
											Type: ast.FieldType{
												Base: &ast.FieldTypeBase{
													Named: testutil.Pointer("string"),
												},
											},
											Children: []*ast.FieldChild{
												{
													Rule: &ast.FieldRule{
														Name: "uuid",
														Body: &ast.FieldRuleBody{
															Error: testutil.Pointer("Product ID must be a valid UUID"),
														},
													},
												},
											},
										},
									},
									{
										Field: &ast.Field{
											Name:     "errors",
											Optional: true,
											Type: ast.FieldType{
												IsArray: true,
												Base: &ast.FieldTypeBase{
													Object: &ast.FieldTypeObject{
														Children: []*ast.FieldOrComment{
															{
																Field: &ast.Field{
																	Name: "code",
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: testutil.Pointer("int"),
																		},
																	},
																},
															},
															{
																Field: &ast.Field{
																	Name: "message",
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: testutil.Pointer("string"),
																		},
																	},
																},
															},
															{
																Field: &ast.Field{
																	Name:     "field",
																	Optional: true,
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: testutil.Pointer("string"),
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
									{
										Field: &ast.Field{
											Name: "analytics",
											Type: ast.FieldType{
												Base: &ast.FieldTypeBase{
													Object: &ast.FieldTypeObject{
														Children: []*ast.FieldOrComment{
															{
																Field: &ast.Field{
																	Name: "duration",
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: testutil.Pointer("float"),
																		},
																	},
																},
															},
															{
																Field: &ast.Field{
																	Name: "processingSteps",
																	Type: ast.FieldType{
																		IsArray: true,
																		Base: &ast.FieldTypeBase{
																			Object: &ast.FieldTypeObject{
																				Children: []*ast.FieldOrComment{
																					{
																						Field: &ast.Field{
																							Name: "name",
																							Type: ast.FieldType{
																								Base: &ast.FieldTypeBase{
																									Named: testutil.Pointer("string"),
																								},
																							},
																						},
																					},
																					{
																						Field: &ast.Field{
																							Name: "duration",
																							Type: ast.FieldType{
																								Base: &ast.FieldTypeBase{
																									Named: testutil.Pointer("float"),
																								},
																							},
																						},
																					},
																					{
																						Field: &ast.Field{
																							Name: "success",
																							Type: ast.FieldType{
																								Base: &ast.FieldTypeBase{
																									Named: testutil.Pointer("boolean"),
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
															{
																Field: &ast.Field{
																	Name: "serverInfo",
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Object: &ast.FieldTypeObject{
																				Children: []*ast.FieldOrComment{
																					{
																						Field: &ast.Field{
																							Name: "id",
																							Type: ast.FieldType{
																								Base: &ast.FieldTypeBase{
																									Named: testutil.Pointer("string"),
																								},
																							},
																						},
																					},
																					{
																						Field: &ast.Field{
																							Name: "region",
																							Type: ast.FieldType{
																								Base: &ast.FieldTypeBase{
																									Named: testutil.Pointer("string"),
																								},
																							},
																						},
																					},
																					{
																						Field: &ast.Field{
																							Name: "load",
																							Type: ast.FieldType{
																								Base: &ast.FieldTypeBase{
																									Named: testutil.Pointer("float"),
																								},
																							},
																							Children: []*ast.FieldChild{
																								{
																									Rule: &ast.FieldRule{
																										Name: "min",
																										Body: &ast.FieldRuleBody{
																											ParamSingle: &ast.AnyLiteral{Float: testutil.Pointer("0.0")},
																										},
																									},
																								},
																								{
																									Rule: &ast.FieldRule{
																										Name: "max",
																										Body: &ast.FieldRuleBody{
																											ParamSingle: &ast.AnyLiteral{Float: testutil.Pointer("1.0")},
																											Error:       testutil.Pointer("Load factor cannot exceed 1.0"),
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
										},
									},
								},
							},
						},
						{
							Meta: &ast.ProcDeclChildMeta{
								Children: []*ast.ProcDeclChildMetaChild{
									{
										KV: &ast.ProcDeclChildMetaKV{
											Key:   "auth",
											Value: ast.AnyLiteral{Str: testutil.Pointer("required")},
										},
									},
									{
										KV: &ast.ProcDeclChildMetaKV{
											Key:   "roles",
											Value: ast.AnyLiteral{Str: testutil.Pointer("admin,product-manager")},
										},
									},
									{
										KV: &ast.ProcDeclChildMetaKV{
											Key:   "rateLimit",
											Value: ast.AnyLiteral{Int: testutil.Pointer("100")},
										},
									},
									{
										KV: &ast.ProcDeclChildMetaKV{
											Key:   "timeout",
											Value: ast.AnyLiteral{Float: testutil.Pointer("30.5")},
										},
									},
									{
										KV: &ast.ProcDeclChildMetaKV{
											Key:   "audit",
											Value: ast.AnyLiteral{True: testutil.Pointer("true")},
										},
									},
									{
										KV: &ast.ProcDeclChildMetaKV{
											Key:   "apiVersion",
											Value: ast.AnyLiteral{Str: testutil.Pointer("1.2.0")},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				Rule: &ast.RuleDecl{
					Docstring: &ast.Docstring{
						Value: "\n\t\tValidates if a value is within a specified range.\n\t\t",
					},
					Name: "range",
					Children: []*ast.RuleDeclChild{
						{
							For: &ast.RuleDeclChildFor{
								For: "int",
							},
						},
						{
							Param: &ast.RuleDeclChildParam{
								Param:   "int",
								IsArray: true,
							},
						},
						{
							Error: &ast.RuleDeclChildError{
								Error: "Value out of range",
							},
						},
					},
				},
			},
		},
	}

	testutil.ASTEqualNoPos(t, expected, parsed)
}
