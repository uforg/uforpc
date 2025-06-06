package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/urpc/internal/urpc/ast"
	"github.com/uforg/uforpc/urpc/internal/util/testutil"
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
				field4: bool
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
										Base: &ast.FieldTypeBase{Named: testutil.Pointer("bool")},
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
						Children: []*ast.ProcOrStreamDeclChild{
							{
								Input: &ast.ProcOrStreamDeclChildInput{
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
						Children: []*ast.ProcOrStreamDeclChild{
							{
								Output: &ast.ProcOrStreamDeclChildOutput{
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
						Children: []*ast.ProcOrStreamDeclChild{
							{
								Meta: &ast.ProcOrStreamDeclChildMeta{
									Children: []*ast.ProcOrStreamDeclChildMetaChild{
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{Key: "key1", Value: ast.AnyLiteral{Str: testutil.Pointer("hello")}},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{Key: "key2", Value: ast.AnyLiteral{Int: testutil.Pointer("123")}},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{Key: "key3", Value: ast.AnyLiteral{Float: testutil.Pointer("1.23")}},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{Key: "key4", Value: ast.AnyLiteral{True: testutil.Pointer("true")}},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{Key: "key5", Value: ast.AnyLiteral{False: testutil.Pointer("false")}},
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
						Children: []*ast.ProcOrStreamDeclChild{
							{
								Input: &ast.ProcOrStreamDeclChildInput{
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
								Output: &ast.ProcOrStreamDeclChildOutput{
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
								Meta: &ast.ProcOrStreamDeclChildMeta{
									Children: []*ast.ProcOrStreamDeclChildMetaChild{
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{
												Key:   "key1",
												Value: ast.AnyLiteral{Str: testutil.Pointer("hello")},
											},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{
												Key:   "key2",
												Value: ast.AnyLiteral{Int: testutil.Pointer("123")},
											},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{
												Key:   "key3",
												Value: ast.AnyLiteral{Float: testutil.Pointer("1.23")},
											},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{
												Key:   "key4",
												Value: ast.AnyLiteral{True: testutil.Pointer("true")},
											},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{
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

func TestParserStreamDecl(t *testing.T) {
	t.Run("Minimum stream declaration parsing", func(t *testing.T) {
		input := `
			stream MyStream {}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Stream: &ast.StreamDecl{
						Name: "MyStream",
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Stream with docstring", func(t *testing.T) {
		input := `
			""" MyStream is a stream that does something. """
			stream MyStream {}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Stream: &ast.StreamDecl{
						Docstring: &ast.Docstring{
							Value: " MyStream is a stream that does something. ",
						},
						Name: "MyStream",
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Stream with deprecated", func(t *testing.T) {
		input := `
			deprecated stream MyStream {}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Stream: &ast.StreamDecl{
						Deprecated: &ast.Deprecated{},
						Name:       "MyStream",
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Stream with deprecated with message", func(t *testing.T) {
		input := `
			deprecated("Use newStream instead")
			stream MyStream {}
		`
		parsed, err := ParserInstance.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.Schema{
			Children: []*ast.SchemaChild{
				{
					Stream: &ast.StreamDecl{
						Deprecated: &ast.Deprecated{
							Message: testutil.Pointer("Use newStream instead"),
						},
						Name: "MyStream",
					},
				},
			},
		}

		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Stream with input", func(t *testing.T) {
		input := `
			stream MyStream {
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
					Stream: &ast.StreamDecl{
						Name: "MyStream",
						Children: []*ast.ProcOrStreamDeclChild{
							{
								Input: &ast.ProcOrStreamDeclChildInput{
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

	t.Run("Stream with output", func(t *testing.T) {
		input := `
			stream MyStream {
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
					Stream: &ast.StreamDecl{
						Name: "MyStream",
						Children: []*ast.ProcOrStreamDeclChild{
							{
								Output: &ast.ProcOrStreamDeclChildOutput{
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

	t.Run("Stream with meta", func(t *testing.T) {
		input := `
			stream MyStream {
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
					Stream: &ast.StreamDecl{
						Name: "MyStream",
						Children: []*ast.ProcOrStreamDeclChild{
							{
								Meta: &ast.ProcOrStreamDeclChildMeta{
									Children: []*ast.ProcOrStreamDeclChildMetaChild{
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{Key: "key1", Value: ast.AnyLiteral{Str: testutil.Pointer("hello")}},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{Key: "key2", Value: ast.AnyLiteral{Int: testutil.Pointer("123")}},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{Key: "key3", Value: ast.AnyLiteral{Float: testutil.Pointer("1.23")}},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{Key: "key4", Value: ast.AnyLiteral{True: testutil.Pointer("true")}},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{Key: "key5", Value: ast.AnyLiteral{False: testutil.Pointer("false")}},
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

	t.Run("Full stream", func(t *testing.T) {
		input := `
			""" MyStream is a stream that does something. """
			deprecated("Use NewStream instead")
			stream MyStream {
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
					Stream: &ast.StreamDecl{
						Docstring: &ast.Docstring{
							Value: " MyStream is a stream that does something. ",
						},
						Deprecated: &ast.Deprecated{
							Message: testutil.Pointer("Use NewStream instead"),
						},
						Name: "MyStream",
						Children: []*ast.ProcOrStreamDeclChild{
							{
								Input: &ast.ProcOrStreamDeclChildInput{
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
								Output: &ast.ProcOrStreamDeclChildOutput{
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
								Meta: &ast.ProcOrStreamDeclChildMeta{
									Children: []*ast.ProcOrStreamDeclChildMetaChild{
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{
												Key:   "key1",
												Value: ast.AnyLiteral{Str: testutil.Pointer("hello")},
											},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{
												Key:   "key2",
												Value: ast.AnyLiteral{Int: testutil.Pointer("123")},
											},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{
												Key:   "key3",
												Value: ast.AnyLiteral{Float: testutil.Pointer("1.23")},
											},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{
												Key:   "key4",
												Value: ast.AnyLiteral{True: testutil.Pointer("true")},
											},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{
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
			/* Type comment */
			type MyType { field: int }
			// Proc comment
			proc MyProc {}
			// Stream comment
			stream MyStream {}
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
					Comment: &ast.Comment{Simple: testutil.Pointer(" Stream comment")},
				},
				{
					Stream: &ast.StreamDecl{Name: "MyStream"},
				},
				{
					Comment: &ast.Comment{Block: testutil.Pointer(" Trailing comment ")},
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
						Children: []*ast.ProcOrStreamDeclChild{
							{
								Comment: &ast.Comment{Simple: testutil.Pointer(" Before input")},
							},
							{
								Input: &ast.ProcOrStreamDeclChildInput{
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
								Output: &ast.ProcOrStreamDeclChildOutput{
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
								Meta: &ast.ProcOrStreamDeclChildMeta{
									Children: []*ast.ProcOrStreamDeclChildMetaChild{
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{
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
						Children: []*ast.ProcOrStreamDeclChild{
							{
								Input: &ast.ProcOrStreamDeclChildInput{
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
						Children: []*ast.ProcOrStreamDeclChild{
							{
								Output: &ast.ProcOrStreamDeclChildOutput{
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
						Children: []*ast.ProcOrStreamDeclChild{
							{
								Meta: &ast.ProcOrStreamDeclChildMeta{

									Children: []*ast.ProcOrStreamDeclChildMetaChild{
										{
											Comment: &ast.Comment{Simple: testutil.Pointer(" Before key1")},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{
												Key:   "key1",
												Value: ast.AnyLiteral{Str: testutil.Pointer("value1")},
											},
										},
										{
											Comment: &ast.Comment{Block: testutil.Pointer(" Between key1 and key2 ")},
										},
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{
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

	t.Run("Comments within StreamDecl Input block", func(t *testing.T) {
		input := `
			stream MyStream {
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
					Stream: &ast.StreamDecl{
						Name: "MyStream",
						Children: []*ast.ProcOrStreamDeclChild{
							{
								Input: &ast.ProcOrStreamDeclChildInput{
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

	t.Run("End-of-line comments", func(t *testing.T) {
		input := `
			version 1 // EOL on version
			type MyType { // EOL on type start
				field: string // EOL on field
			} // EOL on type end
			proc MyProc { // EOL on proc start
				input { f: int } // EOL on input
				output { o: int } // EOL on output
				meta { k: "v" } // EOL on meta
			} // EOL on proc end
			stream MyStream { // EOL on stream start
				input { f: int } // EOL on input
				output { o: int } // EOL on output
				meta { k: "v" } // EOL on meta
			} // EOL on stream end
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
								},
							},
							{Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on field")}},
						},
					},
				},
				{
					Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on type end")},
				}, // Comment after the block
				{
					Proc: &ast.ProcDecl{
						Name: "MyProc",
						Children: []*ast.ProcOrStreamDeclChild{
							{Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on proc start")}}, // Comment inside the block
							{
								Input: &ast.ProcOrStreamDeclChildInput{
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
								Output: &ast.ProcOrStreamDeclChildOutput{
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
								Meta: &ast.ProcOrStreamDeclChildMeta{
									Children: []*ast.ProcOrStreamDeclChildMetaChild{
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{
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
				{
					Stream: &ast.StreamDecl{
						Name: "MyStream",
						Children: []*ast.ProcOrStreamDeclChild{
							{Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on stream start")}}, // Comment inside the block
							{
								Input: &ast.ProcOrStreamDeclChildInput{
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
								Output: &ast.ProcOrStreamDeclChildOutput{
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
								Meta: &ast.ProcOrStreamDeclChildMeta{
									Children: []*ast.ProcOrStreamDeclChildMetaChild{
										{
											KV: &ast.ProcOrStreamDeclChildMetaKV{
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
					Comment: &ast.Comment{Simple: testutil.Pointer(" EOL on stream end")},
				}, // Comment after the block
			},
		}
		testutil.ASTEqualNoPos(t, expected, parsed)
	})

	t.Run("Comments inside empty blocks", func(t *testing.T) {
		input := `
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
						Children: []*ast.ProcOrStreamDeclChild{
							{
								Comment: &ast.Comment{Block: testutil.Pointer(" Proc Comment ")},
							},
							{
								Input: &ast.ProcOrStreamDeclChildInput{
									Children: []*ast.FieldOrComment{
										{
											Comment: &ast.Comment{Block: testutil.Pointer(" Input Comment ")},
										},
									},
								},
							},
							{
								Output: &ast.ProcOrStreamDeclChildOutput{
									Children: []*ast.FieldOrComment{
										{
											Comment: &ast.Comment{Simple: testutil.Pointer(" Output Comment")},
										},
									},
								},
							},
							{
								Meta: &ast.ProcOrStreamDeclChildMeta{
									Children: []*ast.ProcOrStreamDeclChildMetaChild{
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

		// Type declarations

		type FirstDummyType {
			dummyField: datetime
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
			name: string
			description?: string
			isActive: bool
			parentId?: string
		}

		"""
		Product represents a sellable item in the store.
		Products have complex validation rules and can be
		nested inside catalogs.
		"""
		type Product {
			id: string
			name: string
			price: float
			stock: int
			category: Category
			tags?: string[]

			details: {
				dimensions: {
					width: float
					height: float
					depth?: float
				}
				weight?: float
				colors: string[]
				attributes?: {
					name: string
					value: string
				}[]
			}

			variations: {
				sku: string
				price: float
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
			}

			output {
				category: Category
				exists: bool
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
					draft: bool
					notify: bool
					scheduledFor?: string
					tags?: string[]
				}

				validation: {
					skipValidation?: bool
					customRules?: {
						name: string
						severity: int
						message: string
					}[]
				}
			}

			output {
				success: bool
				productId: string
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
						success: bool
					}[]
					serverInfo: {
						id: string
						region: string
						load: float
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

		// Stream declarations
		deprecated stream MyStream {
			input {
				fieldIn1: string
			}
			output {
				fieldOut1: string
			}
			meta {
				cache: true
				cacheTime: 300
			}
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
										Named: testutil.Pointer("bool"),
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
					Children: []*ast.ProcOrStreamDeclChild{
						{
							Input: &ast.ProcOrStreamDeclChildInput{
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
								},
							},
						},
						{
							Output: &ast.ProcOrStreamDeclChildOutput{
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
													Named: testutil.Pointer("bool"),
												},
											},
										},
									},
								},
							},
						},
						{
							Meta: &ast.ProcOrStreamDeclChildMeta{
								Children: []*ast.ProcOrStreamDeclChildMetaChild{
									{
										KV: &ast.ProcOrStreamDeclChildMetaKV{
											Key:   "cache",
											Value: ast.AnyLiteral{True: testutil.Pointer("true")},
										},
									},
									{
										KV: &ast.ProcOrStreamDeclChildMetaKV{
											Key:   "cacheTime",
											Value: ast.AnyLiteral{Int: testutil.Pointer("300")},
										},
									},
									{
										KV: &ast.ProcOrStreamDeclChildMetaKV{
											Key:   "requiresAuth",
											Value: ast.AnyLiteral{False: testutil.Pointer("false")},
										},
									},
									{
										KV: &ast.ProcOrStreamDeclChildMetaKV{
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
					Children: []*ast.ProcOrStreamDeclChild{
						{
							Input: &ast.ProcOrStreamDeclChildInput{
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
																			Named: testutil.Pointer("bool"),
																		},
																	},
																},
															},
															{
																Field: &ast.Field{
																	Name: "notify",
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: testutil.Pointer("bool"),
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
																			Named: testutil.Pointer("bool"),
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
							Output: &ast.ProcOrStreamDeclChildOutput{
								Children: []*ast.FieldOrComment{
									{
										Field: &ast.Field{
											Name: "success",
											Type: ast.FieldType{
												Base: &ast.FieldTypeBase{
													Named: testutil.Pointer("bool"),
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
																									Named: testutil.Pointer("bool"),
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
							Meta: &ast.ProcOrStreamDeclChildMeta{
								Children: []*ast.ProcOrStreamDeclChildMetaChild{
									{
										KV: &ast.ProcOrStreamDeclChildMetaKV{
											Key:   "auth",
											Value: ast.AnyLiteral{Str: testutil.Pointer("required")},
										},
									},
									{
										KV: &ast.ProcOrStreamDeclChildMetaKV{
											Key:   "roles",
											Value: ast.AnyLiteral{Str: testutil.Pointer("admin,product-manager")},
										},
									},
									{
										KV: &ast.ProcOrStreamDeclChildMetaKV{
											Key:   "rateLimit",
											Value: ast.AnyLiteral{Int: testutil.Pointer("100")},
										},
									},
									{
										KV: &ast.ProcOrStreamDeclChildMetaKV{
											Key:   "timeout",
											Value: ast.AnyLiteral{Float: testutil.Pointer("30.5")},
										},
									},
									{
										KV: &ast.ProcOrStreamDeclChildMetaKV{
											Key:   "audit",
											Value: ast.AnyLiteral{True: testutil.Pointer("true")},
										},
									},
									{
										KV: &ast.ProcOrStreamDeclChildMetaKV{
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
				Comment: &ast.Comment{Simple: testutil.Pointer(" Stream declarations")},
			},
			{
				Stream: &ast.StreamDecl{
					Deprecated: &ast.Deprecated{},
					Name:       "MyStream",
					Children: []*ast.ProcOrStreamDeclChild{
						{
							Input: &ast.ProcOrStreamDeclChildInput{
								Children: []*ast.FieldOrComment{
									{
										Field: &ast.Field{
											Name: "fieldIn1",
											Type: ast.FieldType{Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")}},
										},
									},
								},
							},
						},
						{
							Output: &ast.ProcOrStreamDeclChildOutput{
								Children: []*ast.FieldOrComment{
									{
										Field: &ast.Field{
											Name: "fieldOut1",
											Type: ast.FieldType{Base: &ast.FieldTypeBase{Named: testutil.Pointer("string")}},
										},
									},
								},
							},
						},
						{
							Meta: &ast.ProcOrStreamDeclChildMeta{
								Children: []*ast.ProcOrStreamDeclChildMetaChild{
									{
										KV: &ast.ProcOrStreamDeclChildMetaKV{
											Key:   "cache",
											Value: ast.AnyLiteral{True: testutil.Pointer("true")},
										},
									},
									{
										KV: &ast.ProcOrStreamDeclChildMetaKV{
											Key:   "cacheTime",
											Value: ast.AnyLiteral{Int: testutil.Pointer("300")},
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
}
