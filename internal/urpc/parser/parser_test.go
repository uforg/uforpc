package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/urpc/ast"
	"github.com/uforg/uforpc/internal/urpc/lexer"
)

// equalWithoutPos compara dos esquemas ignorando las posiciones
func equalWithoutPos(t *testing.T, expected, actual ast.Schema) {
	// Comparar Version
	require.Equal(t, expected.Version.IsSet, actual.Version.IsSet)
	require.Equal(t, expected.Version.Value, actual.Version.Value)

	// Comparar CustomRules
	require.Equal(t, len(expected.CustomRules), len(actual.CustomRules))
	for i := range expected.CustomRules {
		require.Equal(t, expected.CustomRules[i].Name, actual.CustomRules[i].Name)
		require.Equal(t, expected.CustomRules[i].Doc, actual.CustomRules[i].Doc)
		require.Equal(t, expected.CustomRules[i].Error, actual.CustomRules[i].Error)
		equalTypeWithoutPos(t, expected.CustomRules[i].For, actual.CustomRules[i].For)

		// Param
		require.Equal(t, expected.CustomRules[i].Param.IsArray, actual.CustomRules[i].Param.IsArray)
		require.Equal(t, expected.CustomRules[i].Param.Type, actual.CustomRules[i].Param.Type)
	}

	// Comparar Types
	require.Equal(t, len(expected.Types), len(actual.Types))
	for i := range expected.Types {
		require.Equal(t, expected.Types[i].Name, actual.Types[i].Name)
		require.Equal(t, expected.Types[i].Doc, actual.Types[i].Doc)

		// Fields
		require.Equal(t, len(expected.Types[i].Fields), len(actual.Types[i].Fields))
		for j := range expected.Types[i].Fields {
			equalFieldWithoutPos(t, expected.Types[i].Fields[j], actual.Types[i].Fields[j])
		}
	}

	// Comparar Procedures
	require.Equal(t, len(expected.Procedures), len(actual.Procedures))
	for i := range expected.Procedures {
		require.Equal(t, expected.Procedures[i].Name, actual.Procedures[i].Name)
		require.Equal(t, expected.Procedures[i].Doc, actual.Procedures[i].Doc)

		// Input
		require.Equal(t, len(expected.Procedures[i].Input.Fields), len(actual.Procedures[i].Input.Fields))
		for j := range expected.Procedures[i].Input.Fields {
			equalFieldWithoutPos(t, expected.Procedures[i].Input.Fields[j], actual.Procedures[i].Input.Fields[j])
		}

		// Output
		require.Equal(t, len(expected.Procedures[i].Output.Fields), len(actual.Procedures[i].Output.Fields))
		for j := range expected.Procedures[i].Output.Fields {
			equalFieldWithoutPos(t, expected.Procedures[i].Output.Fields[j], actual.Procedures[i].Output.Fields[j])
		}

		// Metadata
		require.Equal(t, len(expected.Procedures[i].Metadata.Entries), len(actual.Procedures[i].Metadata.Entries))
		for j := range expected.Procedures[i].Metadata.Entries {
			require.Equal(t, expected.Procedures[i].Metadata.Entries[j].Key, actual.Procedures[i].Metadata.Entries[j].Key)
			require.Equal(t, expected.Procedures[i].Metadata.Entries[j].Type, actual.Procedures[i].Metadata.Entries[j].Type)
			require.Equal(t, expected.Procedures[i].Metadata.Entries[j].Value, actual.Procedures[i].Metadata.Entries[j].Value)
		}
	}
}

// equalTypeWithoutPos compara dos tipos ignorando las posiciones
func equalTypeWithoutPos(t *testing.T, expected, actual ast.Type) {
	switch expectedType := expected.(type) {
	case ast.TypePrimitive:
		actualType, ok := actual.(ast.TypePrimitive)
		require.True(t, ok, "Types don't match: expected TypePrimitive, got %T", actual)
		require.Equal(t, expectedType.Name, actualType.Name)
	case ast.TypeCustom:
		actualType, ok := actual.(ast.TypeCustom)
		require.True(t, ok, "Types don't match: expected TypeCustom, got %T", actual)
		require.Equal(t, expectedType.Name, actualType.Name)
	case ast.TypeArray:
		actualType, ok := actual.(ast.TypeArray)
		require.True(t, ok, "Types don't match: expected TypeArray, got %T", actual)
		equalTypeWithoutPos(t, expectedType.ElementsType, actualType.ElementsType)
	case ast.TypeObject:
		actualType, ok := actual.(ast.TypeObject)
		require.True(t, ok, "Types don't match: expected TypeObject, got %T", actual)
		require.Equal(t, len(expectedType.Fields), len(actualType.Fields))
		for i := range expectedType.Fields {
			equalFieldWithoutPos(t, expectedType.Fields[i], actualType.Fields[i])
		}
	}
}

// equalFieldWithoutPos compara dos campos ignorando las posiciones
func equalFieldWithoutPos(t *testing.T, expected, actual ast.Field) {
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.Optional, actual.Optional)
	equalTypeWithoutPos(t, expected.Type, actual.Type)

	// ValidationRules
	require.Equal(t, len(expected.ValidationRules), len(actual.ValidationRules))
	for i := range expected.ValidationRules {
		equalValidationRuleWithoutPos(t, expected.ValidationRules[i], actual.ValidationRules[i])
	}
}

// equalValidationRuleWithoutPos compara dos reglas de validación ignorando las posiciones
func equalValidationRuleWithoutPos(t *testing.T, expected, actual ast.ValidationRule) {
	switch expectedRule := expected.(type) {
	case *ast.ValidationRuleSimple:
		actualRule, ok := actual.(*ast.ValidationRuleSimple)
		require.True(t, ok, "Rules don't match: expected ValidationRuleSimple, got %T", actual)
		require.Equal(t, expectedRule.Name, actualRule.Name)
		require.Equal(t, expectedRule.Error, actualRule.Error)
	case *ast.ValidationRuleWithValue:
		actualRule, ok := actual.(*ast.ValidationRuleWithValue)
		require.True(t, ok, "Rules don't match: expected ValidationRuleWithValue, got %T", actual)
		require.Equal(t, expectedRule.Name, actualRule.Name)
		require.Equal(t, expectedRule.Error, actualRule.Error)
		require.Equal(t, expectedRule.Value, actualRule.Value)
		require.Equal(t, expectedRule.ValueType, actualRule.ValueType)
	case *ast.ValidationRuleWithArray:
		actualRule, ok := actual.(*ast.ValidationRuleWithArray)
		require.True(t, ok, "Rules don't match: expected ValidationRuleWithArray, got %T", actual)
		require.Equal(t, expectedRule.Name, actualRule.Name)
		require.Equal(t, expectedRule.Error, actualRule.Error)
		require.Equal(t, expectedRule.Values, actualRule.Values)
		require.Equal(t, expectedRule.ValueType, actualRule.ValueType)
	}
}

func TestParserVersionDeclaration(t *testing.T) {
	t.Run("Parse version", func(t *testing.T) {
		lexer := lexer.NewLexer("test.urpc", "version 2")
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Version: ast.Version{
				IsSet: true,
				Value: 2,
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse version invalid", func(t *testing.T) {
		lexer := lexer.NewLexer("test.urpc", "version foobar")
		parser := New(lexer)
		_, _, err := parser.Parse()
		require.Error(t, err)
		require.Contains(t, err.Error(), "version expected to be an integer")
	})

	t.Run("Parse version already set", func(t *testing.T) {
		lexer := lexer.NewLexer("test.urpc", "version 2 version 3")
		parser := New(lexer)
		_, _, err := parser.Parse()
		require.Error(t, err)
		require.Contains(t, err.Error(), "version already set")
	})
}

func TestParserTypeDeclaration(t *testing.T) {
	t.Run("Parse type declaration basic", func(t *testing.T) {
		input := `
			type User {}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDecl{
				{
					Name: "User",
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse type declaration with docstring", func(t *testing.T) {
		input := `
			""" Product type documentation """
			type Product {}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDecl{
				{
					Doc:  "Product type documentation",
					Name: "Product",
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse type declaration with primitive type fields", func(t *testing.T) {
		input := `
			type User {
				name: string
				age?: int
				height: float
				isActive: boolean
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDecl{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "name",
							Optional: false,
							Type: ast.TypePrimitive{
								Name: ast.PrimitiveTypeString,
							},
						},
						{
							Name:     "age",
							Optional: true,
							Type: ast.TypePrimitive{
								Name: ast.PrimitiveTypeInt,
							},
						},
						{
							Name:     "height",
							Optional: false,
							Type: ast.TypePrimitive{
								Name: ast.PrimitiveTypeFloat,
							},
						},
						{
							Name:     "isActive",
							Optional: false,
							Type: ast.TypePrimitive{
								Name: ast.PrimitiveTypeBoolean,
							},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse type declaration with field using other custom type", func(t *testing.T) {
		input := `
			type User {
				name: string
				address?: Address
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDecl{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "name",
							Optional: false,
							Type: ast.TypePrimitive{
								Name: ast.PrimitiveTypeString,
							},
						},
						{
							Name:     "address",
							Optional: true,
							Type: ast.TypeCustom{
								Name: "Address",
							},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse type declaration with object type field", func(t *testing.T) {
		input := `
			type User {
				objField: {
					name: string
					age: int
				}
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDecl{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "objField",
							Optional: false,
							Type: ast.TypeObject{
								Fields: []ast.Field{
									{
										Name:     "name",
										Optional: false,
										Type: ast.TypePrimitive{
											Name: ast.PrimitiveTypeString,
										},
									},
									{
										Name:     "age",
										Optional: false,
										Type: ast.TypePrimitive{
											Name: ast.PrimitiveTypeInt,
										},
									},
								},
							},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse type declaration with nested object type field", func(t *testing.T) {
		input := `
			type User {
				objField: {
					name: string
					age: int
					address: {
						street: string
						city: string
						zip: string
					}
				}
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDecl{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "objField",
							Optional: false,
							Type: ast.TypeObject{
								Fields: []ast.Field{
									{
										Name:     "name",
										Optional: false,
										Type: ast.TypePrimitive{
											Name: ast.PrimitiveTypeString,
										},
									},
									{
										Name:     "age",
										Optional: false,
										Type: ast.TypePrimitive{
											Name: ast.PrimitiveTypeInt,
										},
									},
									{
										Name:     "address",
										Optional: false,
										Type: ast.TypeObject{
											Fields: []ast.Field{
												{
													Name:     "street",
													Optional: false,
													Type: ast.TypePrimitive{
														Name: ast.PrimitiveTypeString,
													},
												},
												{
													Name:     "city",
													Optional: false,
													Type: ast.TypePrimitive{
														Name: ast.PrimitiveTypeString,
													},
												},
												{
													Name:     "zip",
													Optional: false,
													Type: ast.TypePrimitive{
														Name: ast.PrimitiveTypeString,
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

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse type declaration with array of objects", func(t *testing.T) {
		input := `
			type User {
				objField: {
					name: string
					age: int
				}[]
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDecl{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "objField",
							Optional: false,
							Type: ast.TypeArray{
								ElementsType: ast.TypeObject{
									Fields: []ast.Field{
										{
											Name:     "name",
											Optional: false,
											Type: ast.TypePrimitive{
												Name: ast.PrimitiveTypeString,
											},
										},
										{
											Name:     "age",
											Optional: false,
											Type: ast.TypePrimitive{
												Name: ast.PrimitiveTypeInt,
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

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse type declaration with array type field", func(t *testing.T) {
		input := `
			type User {
				arrayField: string[]
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDecl{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "arrayField",
							Optional: false,
							Type:     ast.TypeArray{ElementsType: ast.TypePrimitive{Name: ast.PrimitiveTypeString}},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse type declaration with multidimensional array type field", func(t *testing.T) {
		input := `
			type User {
				arrayField: string[][][]
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDecl{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "arrayField",
							Optional: false,
							Type: ast.TypeArray{
								ElementsType: ast.TypeArray{
									ElementsType: ast.TypeArray{
										ElementsType: ast.TypePrimitive{Name: ast.PrimitiveTypeString},
									},
								},
							},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse type declaration with fields containing validation rules", func(t *testing.T) {
		input := `
			type User {
				stringField: string
					@equals("Foo")
					@contains("Bar")
					@minlen(3)
					@maxlen(100)
					@enum(["Foo", "Bar"])
					@email
					@iso8601
					@uuid
					@json
					@lowercase
					@uppercase

				intField: int
					@equals(1)
					@min(0)
					@max(100)
					@enum([1, 2, 3])

				floatField: float
					@min(0.0)
					@max(100.0)

				booleanField: boolean
					@equals(true)

				arrayField: string[]
					@minlen(1)
					@maxlen(100)
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDecl{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "stringField",
							Optional: false,
							Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleWithValue{
									Name:      "equals",
									Value:     "Foo",
									ValueType: ast.ValidationRuleValueTypeString,
								},
								&ast.ValidationRuleWithValue{
									Name:      "contains",
									Value:     "Bar",
									ValueType: ast.ValidationRuleValueTypeString,
								},
								&ast.ValidationRuleWithValue{
									Name:      "minlen",
									Value:     "3",
									ValueType: ast.ValidationRuleValueTypeInt,
								},
								&ast.ValidationRuleWithValue{
									Name:      "maxlen",
									Value:     "100",
									ValueType: ast.ValidationRuleValueTypeInt,
								},
								&ast.ValidationRuleWithArray{
									Name:      "enum",
									Values:    []string{"Foo", "Bar"},
									ValueType: ast.ValidationRuleValueTypeString,
								},
								&ast.ValidationRuleSimple{
									Name: "email",
								},
								&ast.ValidationRuleSimple{
									Name: "iso8601",
								},
								&ast.ValidationRuleSimple{
									Name: "uuid",
								},
								&ast.ValidationRuleSimple{
									Name: "json",
								},
								&ast.ValidationRuleSimple{
									Name: "lowercase",
								},
								&ast.ValidationRuleSimple{
									Name: "uppercase",
								},
							},
						},
						{
							Name:     "intField",
							Optional: false,
							Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeInt},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleWithValue{
									Name:      "equals",
									Value:     "1",
									ValueType: ast.ValidationRuleValueTypeInt,
								},
								&ast.ValidationRuleWithValue{
									Name:      "min",
									Value:     "0",
									ValueType: ast.ValidationRuleValueTypeInt,
								},
								&ast.ValidationRuleWithValue{
									Name:      "max",
									Value:     "100",
									ValueType: ast.ValidationRuleValueTypeInt,
								},
								&ast.ValidationRuleWithArray{
									Name:      "enum",
									Values:    []string{"1", "2", "3"},
									ValueType: ast.ValidationRuleValueTypeInt,
								},
							},
						},
						{
							Name:     "floatField",
							Optional: false,
							Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeFloat},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleWithValue{
									Name:      "min",
									Value:     "0.0",
									ValueType: ast.ValidationRuleValueTypeFloat,
								},
								&ast.ValidationRuleWithValue{
									Name:      "max",
									Value:     "100.0",
									ValueType: ast.ValidationRuleValueTypeFloat,
								},
							},
						},
						{
							Name:     "booleanField",
							Optional: false,
							Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeBoolean},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleWithValue{
									Name:      "equals",
									Value:     "true",
									ValueType: ast.ValidationRuleValueTypeBoolean,
								},
							},
						},
						{
							Name:     "arrayField",
							Optional: false,
							Type:     ast.TypeArray{ElementsType: ast.TypePrimitive{Name: ast.PrimitiveTypeString}},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleWithValue{
									Name:      "minlen",
									Value:     "1",
									ValueType: ast.ValidationRuleValueTypeInt,
								},
								&ast.ValidationRuleWithValue{
									Name:      "maxlen",
									Value:     "100",
									ValueType: ast.ValidationRuleValueTypeInt,
								},
							},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})
}

func TestParserProcedureDeclaration(t *testing.T) {
	t.Run("Parse procedure declaration basic", func(t *testing.T) {
		input := `
			proc GetUser {
				input {}
				output {}
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Procedures: []ast.ProcDecl{
				{
					Name: "GetUser",
					Doc:  "",
					Input: ast.ProcInput{
						Fields: []ast.Field{},
					},
					Output: ast.ProcOutput{
						Fields: []ast.Field{},
					},
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse procedure declaration with docstring", func(t *testing.T) {
		input := `
			""" Create user procedure documentation """
			proc CreateUser {}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Procedures: []ast.ProcDecl{
				{
					Doc:      "Create user procedure documentation",
					Name:     "CreateUser",
					Input:    ast.ProcInput{},
					Output:   ast.ProcOutput{},
					Metadata: ast.ProcMeta{},
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse procedure declaration with empty input, output and meta", func(t *testing.T) {
		input := `
			proc CreateUser {
				input {}
				output {}
				meta {}
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Procedures: []ast.ProcDecl{
				{
					Name:     "CreateUser",
					Input:    ast.ProcInput{},
					Output:   ast.ProcOutput{},
					Metadata: ast.ProcMeta{},
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse procedure declaration with docstring, input, output and meta", func(t *testing.T) {
		input := `
			""" Creates a product with the given stock and returns the product id. """
			proc CreateProduct {
				input {
					product: Product
					stock: int
				}
				
				output {
					productId: string
				}
				
				meta {
					versionNumber: "1.0.0"
					maxRetries: 3
					waitMinutes: 10.5
					audit: true
				}
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Procedures: []ast.ProcDecl{
				{
					Doc:  "Creates a product with the given stock and returns the product id.",
					Name: "CreateProduct",
					Input: ast.ProcInput{
						Fields: []ast.Field{
							{
								Name:     "product",
								Optional: false,
								Type:     ast.TypeCustom{Name: "Product"},
							},
							{
								Name:     "stock",
								Optional: false,
								Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeInt},
							},
						},
					},
					Output: ast.ProcOutput{
						Fields: []ast.Field{
							{
								Name:     "productId",
								Optional: false,
								Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
							},
						},
					},
					Metadata: ast.ProcMeta{
						Entries: []ast.ProcMetaKV{
							{
								Type:  ast.PrimitiveTypeString,
								Key:   "versionNumber",
								Value: "1.0.0",
							},
							{
								Type:  ast.PrimitiveTypeInt,
								Key:   "maxRetries",
								Value: "3",
							},
							{
								Type:  ast.PrimitiveTypeFloat,
								Key:   "waitMinutes",
								Value: "10.5",
							},
							{
								Type:  ast.PrimitiveTypeBoolean,
								Key:   "audit",
								Value: "true",
							},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})
}

func TestParserValidationRules(t *testing.T) {
	t.Run("Parse validation rule with error message", func(t *testing.T) {
		input := `
			type User {
				name: string @required(error: "Name is required")
				age: int @min(18, error: "Must be an adult")
				email: string @email(error: "Invalid email format")
				options: string[] @enum(["a", "b", "c"], error: "Invalid option selected")
				tag: string @pattern("^[a-z]+$", error: "Only lowercase letters allowed")
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDecl{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "name",
							Optional: false,
							Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleSimple{
									Name:  "required",
									Error: "Name is required",
								},
							},
						},
						{
							Name:     "age",
							Optional: false,
							Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeInt},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleWithValue{
									Name:      "min",
									Value:     "18",
									ValueType: ast.ValidationRuleValueTypeInt,
									Error:     "Must be an adult",
								},
							},
						},
						{
							Name:     "email",
							Optional: false,
							Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleSimple{
									Name:  "email",
									Error: "Invalid email format",
								},
							},
						},
						{
							Name:     "options",
							Optional: false,
							Type:     ast.TypeArray{ElementsType: ast.TypePrimitive{Name: ast.PrimitiveTypeString}},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleWithArray{
									Name:      "enum",
									Values:    []string{"a", "b", "c"},
									ValueType: ast.ValidationRuleValueTypeString,
									Error:     "Invalid option selected",
								},
							},
						},
						{
							Name:     "tag",
							Optional: false,
							Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleWithValue{
									Name:      "pattern",
									Value:     "^[a-z]+$",
									ValueType: ast.ValidationRuleValueTypeString,
									Error:     "Only lowercase letters allowed",
								},
							},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse validation rule with error only", func(t *testing.T) {
		input := `
			type User {
				name: string @required(error: "This field cannot be empty")
				email: string @email(error: "Please enter a valid email address")
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDecl{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "name",
							Optional: false,
							Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleSimple{
									Name:  "required",
									Error: "This field cannot be empty",
								},
							},
						},
						{
							Name:     "email",
							Optional: false,
							Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleSimple{
									Name:  "email",
									Error: "Please enter a valid email address",
								},
							},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})
}

func TestParserCustomRuleDeclaration(t *testing.T) {
	t.Run("Parse custom rule declaration basic", func(t *testing.T) {
		input := `
			rule @minlen {
				for: string
				param: int
				error: "Value is too short"
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			CustomRules: []ast.CustomRuleDecl{
				{
					Name: "minlen",
					Doc:  "",
					For:  ast.TypePrimitive{Name: ast.PrimitiveTypeString},
					Param: ast.CustomRuleDeclParamType{
						IsArray: false,
						Type:    ast.PrimitiveTypeInt,
					},
					Error: "Value is too short",
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse custom rule declaration with array param", func(t *testing.T) {
		input := `
			rule @enum {
				for: string
				param: string[]
				error: "Value must be one of the allowed options"
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			CustomRules: []ast.CustomRuleDecl{
				{
					Name: "enum",
					Doc:  "",
					For:  ast.TypePrimitive{Name: ast.PrimitiveTypeString},
					Param: ast.CustomRuleDeclParamType{
						IsArray: true,
						Type:    ast.PrimitiveTypeString,
					},
					Error: "Value must be one of the allowed options",
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse custom rule declaration with docstring", func(t *testing.T) {
		input := `
			"""
			Validates if a string matches a regular expression pattern.
			This rule is useful for format validation like emails, etc.
			"""
			rule @regex {
				for: string
				param: string
				error: "Invalid format"
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			CustomRules: []ast.CustomRuleDecl{
				{
					Name: "regex",
					Doc:  "Validates if a string matches a regular expression pattern.\n\t\t\tThis rule is useful for format validation like emails, etc.",
					For:  ast.TypePrimitive{Name: ast.PrimitiveTypeString},
					Param: ast.CustomRuleDeclParamType{
						IsArray: false,
						Type:    ast.PrimitiveTypeString,
					},
					Error: "Invalid format",
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse custom rule declaration with multiple rules", func(t *testing.T) {
		input := `
			rule @lowercase {
				for: string
				error: "Must be lowercase"
			}

			rule @uppercase {
				for: string
				error: "Must be uppercase"
			}

			rule @range {
				for: int
				param: int[]
				error: "Value out of range"
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			CustomRules: []ast.CustomRuleDecl{
				{
					Name:  "lowercase",
					Doc:   "",
					For:   ast.TypePrimitive{Name: ast.PrimitiveTypeString},
					Param: ast.CustomRuleDeclParamType{},
					Error: "Must be lowercase",
				},
				{
					Name:  "uppercase",
					Doc:   "",
					For:   ast.TypePrimitive{Name: ast.PrimitiveTypeString},
					Param: ast.CustomRuleDeclParamType{},
					Error: "Must be uppercase",
				},
				{
					Name: "range",
					Doc:  "",
					For:  ast.TypePrimitive{Name: ast.PrimitiveTypeInt},
					Param: ast.CustomRuleDeclParamType{
						IsArray: true,
						Type:    ast.PrimitiveTypeInt,
					},
					Error: "Value out of range",
				},
			},
		}

		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})

	t.Run("Parse custom rule declaration invalid for type", func(t *testing.T) {
		input := `
			rule @invalid {
				for: invalidType
				param: int
				error: "Test error"
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		_, _, err := parser.Parse()

		require.Error(t, err)
		require.Contains(t, err.Error(), "must be in PascalCase")
	})

	t.Run("Parse custom rule declaration invalid param type", func(t *testing.T) {
		input := `
			rule @invalid {
				for: string
				param: invalidType
				error: "Test error"
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		_, _, err := parser.Parse()

		require.Error(t, err)
		require.Contains(t, err.Error(), `invalid "invalidType" param type`)
	})

	t.Run("Parse custom rule declaration for complex types", func(t *testing.T) {
		input := `
			type CustomType {
				field1: string
				field2: int
				field3: float
			}

			rule @rule1 {
				for: string
			}
			
			rule @rule2 {
				for: int
			}
			
			rule @rule3 {
				for: float
			}

			rule @rule4 {
				for: boolean
			}

			rule @rule5 {
				for: string[]
			}

			rule @rule6 {
				for: CustomType
			}
			
			rule @rule7 {
				for: CustomType[]
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDecl{
				{
					Name: "CustomType",
					Fields: []ast.Field{
						{Name: "field1", Optional: false, Type: ast.TypePrimitive{Name: ast.PrimitiveTypeString}},
						{Name: "field2", Optional: false, Type: ast.TypePrimitive{Name: ast.PrimitiveTypeInt}},
						{Name: "field3", Optional: false, Type: ast.TypePrimitive{Name: ast.PrimitiveTypeFloat}},
					},
				},
			},
			CustomRules: []ast.CustomRuleDecl{
				{
					Name: "rule1",
					For:  ast.TypePrimitive{Name: ast.PrimitiveTypeString},
				},
				{
					Name: "rule2",
					For:  ast.TypePrimitive{Name: ast.PrimitiveTypeInt},
				},
				{
					Name: "rule3",
					For:  ast.TypePrimitive{Name: ast.PrimitiveTypeFloat},
				},
				{
					Name: "rule4",
					For:  ast.TypePrimitive{Name: ast.PrimitiveTypeBoolean},
				},
				{
					Name: "rule5",
					For:  ast.TypeArray{ElementsType: ast.TypePrimitive{Name: ast.PrimitiveTypeString}},
				},
				{
					Name: "rule6",
					For:  ast.TypeCustom{Name: "CustomType"},
				},
				{
					Name: "rule7",
					For:  ast.TypeArray{ElementsType: ast.TypeCustom{Name: "CustomType"}},
				},
			},
		}
		require.NoError(t, err)
		equalWithoutPos(t, expected, schema)
	})
}

func TestParserFullExample(t *testing.T) {
	input := `
		// Version declaration
		version 1

		// Custom rule declarations
		"""
		This rule validates if a string matches a regular expression pattern.
		Useful for emails, URLs, and other formatted strings.
		"""
		rule @regex {
			for: string
			param: string
			error: "Invalid format"
		}

		"""
		Validates if a value is within a specified range.
		"""
		rule @range {
			for: int
			param: int[]
			error: "Value out of range"
		}

		// Simple type with documentation
		"""
		Category represents a product category in the system.
		This type is used across the catalog module.
		"""
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

		""" Validate category with custom logic """
		rule @validateCategory {
			for: Category
			error: "Invalid category"
		}

		// Type with nested objects and arrays
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

		// Simple procedure 
		"""
		GetCategory retrieves a category by its ID.
		This is a basic read operation.
		"""
		proc GetCategory {
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

		// Complex procedure with nested types
		"""
		CreateProduct adds a new product to the catalog.
		This procedure handles complex validation and returns
		detailed success information.
		"""
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
	`

	lexer := lexer.NewLexer("comprehensive.urpc", input)
	parser := New(lexer)
	schema, _, err := parser.Parse()

	expected := ast.Schema{
		Version: ast.Version{
			IsSet: true,
			Value: 1,
		},
		CustomRules: []ast.CustomRuleDecl{
			{
				Name: "regex",
				Doc:  "This rule validates if a string matches a regular expression pattern.\n\t\tUseful for emails, URLs, and other formatted strings.",
				For:  ast.TypePrimitive{Name: ast.PrimitiveTypeString},
				Param: ast.CustomRuleDeclParamType{
					IsArray: false,
					Type:    ast.PrimitiveTypeString,
				},
				Error: "Invalid format",
			},
			{
				Name: "range",
				Doc:  "Validates if a value is within a specified range.",
				For:  ast.TypePrimitive{Name: ast.PrimitiveTypeInt},
				Param: ast.CustomRuleDeclParamType{
					IsArray: true,
					Type:    ast.PrimitiveTypeInt,
				},
				Error: "Value out of range",
			},
			{
				Name:  "validateCategory",
				Doc:   "Validate category with custom logic",
				For:   ast.TypeCustom{Name: "Category"},
				Error: "Invalid category",
			},
		},
		Types: []ast.TypeDecl{
			{
				Name: "Category",
				Doc:  "Category represents a product category in the system.\n\t\tThis type is used across the catalog module.",
				Fields: []ast.Field{
					{
						Name:     "id",
						Optional: false,
						Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleSimple{
								Name:  "uuid",
								Error: "Must be a valid UUID",
							},
							&ast.ValidationRuleWithValue{
								Name:      "minlen",
								Value:     "36",
								ValueType: ast.ValidationRuleValueTypeInt,
								Error:     "",
							},
							&ast.ValidationRuleWithValue{
								Name:      "maxlen",
								Value:     "36",
								ValueType: ast.ValidationRuleValueTypeInt,
								Error:     "UUID must be exactly 36 characters",
							},
						},
					},
					{
						Name:     "name",
						Optional: false,
						Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleWithValue{
								Name:      "minlen",
								Value:     "3",
								ValueType: ast.ValidationRuleValueTypeInt,
								Error:     "Name must be at least 3 characters long",
							},
						},
					},
					{
						Name:     "description",
						Optional: true,
						Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
					},
					{
						Name:     "isActive",
						Optional: false,
						Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeBoolean},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleWithValue{
								Name:      "equals",
								Value:     "true",
								ValueType: ast.ValidationRuleValueTypeBoolean,
								Error:     "",
							},
						},
					},
					{
						Name:     "parentId",
						Optional: true,
						Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleSimple{
								Name:  "uuid",
								Error: "",
							},
						},
					},
				},
			},
			{
				Name: "Product",
				Doc:  "Product represents a sellable item in the store.\n\t\tProducts have complex validation rules and can be\n\t\tnested inside catalogs.",
				Fields: []ast.Field{
					{
						Name:     "id",
						Optional: false,
						Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleSimple{
								Name:  "uuid",
								Error: "",
							},
						},
					},
					{
						Name:     "name",
						Optional: false,
						Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleWithValue{
								Name:      "minlen",
								Value:     "2",
								ValueType: ast.ValidationRuleValueTypeInt,
								Error:     "",
							},
							&ast.ValidationRuleWithValue{
								Name:      "maxlen",
								Value:     "100",
								ValueType: ast.ValidationRuleValueTypeInt,
								Error:     "Name cannot exceed 100 characters",
							},
						},
					},
					{
						Name:     "price",
						Optional: false,
						Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeFloat},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleWithValue{
								Name:      "min",
								Value:     "0.01",
								ValueType: ast.ValidationRuleValueTypeFloat,
								Error:     "Price must be greater than zero",
							},
						},
					},
					{
						Name:     "stock",
						Optional: false,
						Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeInt},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleWithValue{
								Name:      "min",
								Value:     "0",
								ValueType: ast.ValidationRuleValueTypeInt,
								Error:     "",
							},
							&ast.ValidationRuleWithArray{
								Name:      "range",
								Values:    []string{"0", "1000"},
								ValueType: ast.ValidationRuleValueTypeInt,
								Error:     "Stock must be between 0 and 1000",
							},
						},
					},
					{
						Name:     "category",
						Optional: false,
						Type:     ast.TypeCustom{Name: "Category"},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleSimple{
								Name:  "validateCategory",
								Error: "Invalid category custom message",
							},
						},
					},
					{
						Name:     "tags",
						Optional: true,
						Type:     ast.TypeArray{ElementsType: ast.TypePrimitive{Name: ast.PrimitiveTypeString}},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleWithValue{
								Name:      "minlen",
								Value:     "1",
								ValueType: ast.ValidationRuleValueTypeInt,
								Error:     "At least one tag is required",
							},
							&ast.ValidationRuleWithValue{
								Name:      "maxlen",
								Value:     "10",
								ValueType: ast.ValidationRuleValueTypeInt,
								Error:     "",
							},
						},
					},
					{
						Name:     "details",
						Optional: false,
						Type: ast.TypeObject{
							Fields: []ast.Field{
								{
									Name:     "dimensions",
									Optional: false,
									Type: ast.TypeObject{
										Fields: []ast.Field{
											{
												Name:     "width",
												Optional: false,
												Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeFloat},
												ValidationRules: []ast.ValidationRule{
													&ast.ValidationRuleWithValue{
														Name:      "min",
														Value:     "0.0",
														ValueType: ast.ValidationRuleValueTypeFloat,
														Error:     "Width cannot be negative",
													},
												},
											},
											{
												Name:     "height",
												Optional: false,
												Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeFloat},
												ValidationRules: []ast.ValidationRule{
													&ast.ValidationRuleWithValue{
														Name:      "min",
														Value:     "0.0",
														ValueType: ast.ValidationRuleValueTypeFloat,
														Error:     "",
													},
												},
											},
											{
												Name:     "depth",
												Optional: true,
												Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeFloat},
											},
										},
									},
								},
								{
									Name:     "weight",
									Optional: true,
									Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeFloat},
								},
								{
									Name:     "colors",
									Optional: false,
									Type:     ast.TypeArray{ElementsType: ast.TypePrimitive{Name: ast.PrimitiveTypeString}},
									ValidationRules: []ast.ValidationRule{
										&ast.ValidationRuleWithArray{
											Name:      "enum",
											Values:    []string{"red", "green", "blue", "black", "white"},
											ValueType: ast.ValidationRuleValueTypeString,
											Error:     "Color must be one of the allowed values",
										},
									},
								},
								{
									Name:     "attributes",
									Optional: true,
									Type: ast.TypeArray{
										ElementsType: ast.TypeObject{
											Fields: []ast.Field{
												{
													Name:     "name",
													Optional: false,
													Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
												},
												{
													Name:     "value",
													Optional: false,
													Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						Name:     "variations",
						Optional: false,
						Type: ast.TypeArray{
							ElementsType: ast.TypeObject{
								Fields: []ast.Field{
									{
										Name:     "sku",
										Optional: false,
										Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
									},
									{
										Name:     "price",
										Optional: false,
										Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeFloat},
										ValidationRules: []ast.ValidationRule{
											&ast.ValidationRuleWithValue{
												Name:      "min",
												Value:     "0.01",
												ValueType: ast.ValidationRuleValueTypeFloat,
												Error:     "Variation price must be greater than zero",
											},
										},
									},
									{
										Name:     "attributes",
										Optional: false,
										Type: ast.TypeArray{
											ElementsType: ast.TypeObject{
												Fields: []ast.Field{
													{
														Name:     "name",
														Optional: false,
														Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
													},
													{
														Name:     "value",
														Optional: false,
														Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
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
		Procedures: []ast.ProcDecl{
			{
				Name: "GetCategory",
				Doc:  "GetCategory retrieves a category by its ID.\n\t\tThis is a basic read operation.",
				Input: ast.ProcInput{
					Fields: []ast.Field{
						{
							Name:     "id",
							Optional: false,
							Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleSimple{
									Name:  "uuid",
									Error: "Category ID must be a valid UUID",
								},
							},
						},
					},
				},
				Output: ast.ProcOutput{
					Fields: []ast.Field{
						{
							Name:     "category",
							Optional: false,
							Type:     ast.TypeCustom{Name: "Category"},
						},
						{
							Name:     "exists",
							Optional: false,
							Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeBoolean},
						},
					},
				},
				Metadata: ast.ProcMeta{
					Entries: []ast.ProcMetaKV{
						{
							Type:  ast.PrimitiveTypeBoolean,
							Key:   "cache",
							Value: "true",
						},
						{
							Type:  ast.PrimitiveTypeInt,
							Key:   "cacheTime",
							Value: "300",
						},
						{
							Type:  ast.PrimitiveTypeBoolean,
							Key:   "requiresAuth",
							Value: "false",
						},
						{
							Type:  ast.PrimitiveTypeString,
							Key:   "apiVersion",
							Value: "1.0.0",
						},
					},
				},
			},
			{
				Name: "CreateProduct",
				Doc:  "CreateProduct adds a new product to the catalog.\n\t\tThis procedure handles complex validation and returns\n\t\tdetailed success information.",
				Input: ast.ProcInput{
					Fields: []ast.Field{
						{
							Name:     "product",
							Optional: false,
							Type:     ast.TypeCustom{Name: "Product"},
						},
						{
							Name:     "options",
							Optional: true,
							Type: ast.TypeObject{
								Fields: []ast.Field{
									{
										Name:     "draft",
										Optional: false,
										Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeBoolean},
									},
									{
										Name:     "notify",
										Optional: false,
										Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeBoolean},
									},
									{
										Name:     "scheduledFor",
										Optional: true,
										Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
										ValidationRules: []ast.ValidationRule{
											&ast.ValidationRuleSimple{
												Name:  "iso8601",
												Error: "Must be a valid ISO8601 date",
											},
										},
									},
									{
										Name:     "tags",
										Optional: true,
										Type:     ast.TypeArray{ElementsType: ast.TypePrimitive{Name: ast.PrimitiveTypeString}},
									},
								},
							},
						},
						{
							Name:     "validation",
							Optional: false,
							Type: ast.TypeObject{
								Fields: []ast.Field{
									{
										Name:     "skipValidation",
										Optional: true,
										Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeBoolean},
									},
									{
										Name:     "customRules",
										Optional: true,
										Type: ast.TypeArray{
											ElementsType: ast.TypeObject{
												Fields: []ast.Field{
													{
														Name:     "name",
														Optional: false,
														Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
													},
													{
														Name:     "severity",
														Optional: false,
														Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeInt},
														ValidationRules: []ast.ValidationRule{
															&ast.ValidationRuleWithArray{
																Name:      "enum",
																Values:    []string{"1", "2", "3"},
																ValueType: ast.ValidationRuleValueTypeInt,
																Error:     "Severity must be 1, 2, or 3",
															},
														},
													},
													{
														Name:     "message",
														Optional: false,
														Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
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
				Output: ast.ProcOutput{
					Fields: []ast.Field{
						{
							Name:     "success",
							Optional: false,
							Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeBoolean},
						},
						{
							Name:     "productId",
							Optional: false,
							Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleSimple{
									Name:  "uuid",
									Error: "Product ID must be a valid UUID",
								},
							},
						},
						{
							Name:     "errors",
							Optional: true,
							Type: ast.TypeArray{
								ElementsType: ast.TypeObject{
									Fields: []ast.Field{
										{
											Name:     "code",
											Optional: false,
											Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeInt},
										},
										{
											Name:     "message",
											Optional: false,
											Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
										},
										{
											Name:     "field",
											Optional: true,
											Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
										},
									},
								},
							},
						},
						{
							Name:     "analytics",
							Optional: false,
							Type: ast.TypeObject{
								Fields: []ast.Field{
									{
										Name:     "duration",
										Optional: false,
										Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeFloat},
									},
									{
										Name:     "processingSteps",
										Optional: false,
										Type: ast.TypeArray{
											ElementsType: ast.TypeObject{
												Fields: []ast.Field{
													{
														Name:     "name",
														Optional: false,
														Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
													},
													{
														Name:     "duration",
														Optional: false,
														Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeFloat},
													},
													{
														Name:     "success",
														Optional: false,
														Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeBoolean},
													},
												},
											},
										},
									},
									{
										Name:     "serverInfo",
										Optional: false,
										Type: ast.TypeObject{
											Fields: []ast.Field{
												{
													Name:     "id",
													Optional: false,
													Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
												},
												{
													Name:     "region",
													Optional: false,
													Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeString},
												},
												{
													Name:     "load",
													Optional: false,
													Type:     ast.TypePrimitive{Name: ast.PrimitiveTypeFloat},
													ValidationRules: []ast.ValidationRule{
														&ast.ValidationRuleWithValue{
															Name:      "min",
															Value:     "0.0",
															ValueType: ast.ValidationRuleValueTypeFloat,
															Error:     "",
														},
														&ast.ValidationRuleWithValue{
															Name:      "max",
															Value:     "1.0",
															ValueType: ast.ValidationRuleValueTypeFloat,
															Error:     "Load factor cannot exceed 1.0",
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
				Metadata: ast.ProcMeta{
					Entries: []ast.ProcMetaKV{
						{
							Type:  ast.PrimitiveTypeString,
							Key:   "auth",
							Value: "required",
						},
						{
							Type:  ast.PrimitiveTypeString,
							Key:   "roles",
							Value: "admin,product-manager",
						},
						{
							Type:  ast.PrimitiveTypeInt,
							Key:   "rateLimit",
							Value: "100",
						},
						{
							Type:  ast.PrimitiveTypeFloat,
							Key:   "timeout",
							Value: "30.5",
						},
						{
							Type:  ast.PrimitiveTypeBoolean,
							Key:   "audit",
							Value: "true",
						},
						{
							Type:  ast.PrimitiveTypeString,
							Key:   "apiVersion",
							Value: "1.2.0",
						},
					},
				},
			},
		},
	}

	require.NoError(t, err)
	equalWithoutPos(t, expected, schema)
}
