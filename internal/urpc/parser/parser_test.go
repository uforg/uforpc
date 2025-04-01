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

		for i, procDecl := range ast.Procs {
			ast.Procs[i] = setEmptyPos(procDecl)
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

func TestParserProcDecl(t *testing.T) {
	t.Run("Minimum procedure declaration parsing", func(t *testing.T) {
		input := `
			proc MyProc {}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Procs: []*ast.ProcDecl{
				{
					Name: "MyProc",
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Procedure with docstring", func(t *testing.T) {
		input := `
			"""
			MyProc is a procedure that does something.
			"""
			proc MyProc {}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Procs: []*ast.ProcDecl{
				{
					Docstring: "MyProc is a procedure that does something.",
					Name:      "MyProc",
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Procedure with input", func(t *testing.T) {
		input := `
			proc MyProc {
				input {
					field1: string
				}
			}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Procs: []*ast.ProcDecl{
				{
					Name: "MyProc",
					Body: ast.ProcDeclBody{
						Input: []*ast.Field{
							{
								Name: "field1",
								Type: ast.FieldType{Base: &ast.FieldTypeBase{Named: ptr("string")}},
							},
						},
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Procedure with output", func(t *testing.T) {
		input := `
			proc MyProc {
				output {
					field1: int
				}
			}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Procs: []*ast.ProcDecl{
				{
					Name: "MyProc",
					Body: ast.ProcDeclBody{
						Output: []*ast.Field{
							{
								Name: "field1",
								Type: ast.FieldType{Base: &ast.FieldTypeBase{Named: ptr("int")}},
							},
						},
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
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
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Procs: []*ast.ProcDecl{
				{
					Name: "MyProc",
					Body: ast.ProcDeclBody{
						Meta: []*ast.ProcDeclBodyMetaKV{
							{Key: "key1", Value: "hello"},
							{Key: "key2", Value: "123"},
							{Key: "key3", Value: "1.23"},
							{Key: "key4", Value: "true"},
							{Key: "key5", Value: "false"},
						},
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Full procedure", func(t *testing.T) {
		input := `
			"""
			MyProc is a procedure that does something.
			"""
			proc MyProc {
				input {
					input1: string[][]
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
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Procs: []*ast.ProcDecl{
				{
					Docstring: "MyProc is a procedure that does something.",
					Name:      "MyProc",
					Body: ast.ProcDeclBody{
						Input: []*ast.Field{
							{
								Name: "input1",
								Type: ast.FieldType{
									Depth: 2,
									Base:  &ast.FieldTypeBase{Named: ptr("string")},
								},
							},
						},
						Output: []*ast.Field{
							{
								Name:     "output1",
								Optional: true,
								Type: ast.FieldType{
									Base: &ast.FieldTypeBase{Named: ptr("int")},
								},
							},
						},
						Meta: []*ast.ProcDeclBodyMetaKV{
							{
								Key:   "key1",
								Value: "hello",
							},
							{
								Key:   "key2",
								Value: "123",
							},
							{
								Key:   "key3",
								Value: "1.23",
							},
							{
								Key:   "key4",
								Value: "true",
							},
							{
								Key:   "key5",
								Value: "false",
							},
						},
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})
}

func TestParserFullSchema(t *testing.T) {
	input := `
		// Version declaration
		version 1

		/* Import other schema */
		import "my_sub_schema.urpc"

		//////////////////////////////
		// Custom rule declarations //
		//////////////////////////////

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

		""" Validate "category" with custom logic """
		rule @validateCategory {
			for: Category
			error: "Field \"category\" is not valid"
		}

		///////////////////////
		// Type declarations //
		///////////////////////

		type FirstDummyType {
			dummyField: datetime
				@min("1900-01-01T00:00:00Z")
				@max("2100-01-01T00:00:00Z")
		}

		type SecondDummyType {
			dummyField: int
		}

		type ThirdDummyType extends FirstDummyType, SecondDummyType {
			dummyField: string
		}

		"""
		Category represents a product category in the system.
		This type is used across the catalog module.
		"""
		type Category extends ThirdDummyType {
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

		////////////////////////////
		// Procedure declarations //
		////////////////////////////

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

	parsed, err := Parser.ParseString("schema.urpc", input)
	require.NoError(t, err)

	expected := &ast.URPCSchema{
		Version: &ast.Version{
			Number: 1,
		},
		Imports: []*ast.Import{
			{
				Path: "my_sub_schema.urpc",
			},
		},
		Rules: []*ast.RuleDecl{
			{
				Docstring: "This rule validates if a string matches a regular expression pattern.\n\t\tUseful for emails, URLs, and other formatted strings.",
				Name:      "regex",
				Body: ast.RuleDeclBody{
					For:   "string",
					Param: "string",
					Error: "Invalid format",
				},
			},
			{
				Docstring: "Validates if a value is within a specified range.",
				Name:      "range",
				Body: ast.RuleDeclBody{
					For:          "int",
					Param:        "int",
					ParamIsArray: true,
					Error:        "Value out of range",
				},
			},
			{
				Docstring: "Validate \"category\" with custom logic",
				Name:      "validateCategory",
				Body: ast.RuleDeclBody{
					For:   "Category",
					Error: "Field \"category\" is not valid",
				},
			},
		},
		Types: []*ast.TypeDecl{
			{
				Name: "FirstDummyType",
				Fields: []*ast.Field{
					{
						Name: "dummyField",
						Type: ast.FieldType{
							Base: &ast.FieldTypeBase{
								Named: ptr("datetime"),
							},
						},
						Rules: []*ast.FieldRule{
							{
								Name: "min",
								Body: ast.FieldRuleBody{
									ParamSingle: ptr("1900-01-01T00:00:00Z"),
								},
							},
							{
								Name: "max",
								Body: ast.FieldRuleBody{
									ParamSingle: ptr("2100-01-01T00:00:00Z"),
								},
							},
						},
					},
				},
			},
			{
				Name: "SecondDummyType",
				Fields: []*ast.Field{
					{
						Name: "dummyField",
						Type: ast.FieldType{
							Base: &ast.FieldTypeBase{
								Named: ptr("int"),
							},
						},
					},
				},
			},
			{
				Name: "ThirdDummyType",
				Extends: []string{
					"FirstDummyType",
					"SecondDummyType",
				},
				Fields: []*ast.Field{
					{
						Name: "dummyField",
						Type: ast.FieldType{
							Base: &ast.FieldTypeBase{
								Named: ptr("string"),
							},
						},
					},
				},
			},
			{
				Docstring: "Category represents a product category in the system.\n\t\tThis type is used across the catalog module.",
				Name:      "Category",
				Extends: []string{
					"ThirdDummyType",
				},
				Fields: []*ast.Field{
					{
						Name: "id",
						Type: ast.FieldType{
							Base: &ast.FieldTypeBase{
								Named: ptr("string"),
							},
						},
						Rules: []*ast.FieldRule{
							{
								Name: "uuid",
								Body: ast.FieldRuleBody{
									Error: "Must be a valid UUID",
								},
							},
							{
								Name: "minlen",
								Body: ast.FieldRuleBody{
									ParamSingle: ptr("36"),
								},
							},
							{
								Name: "maxlen",
								Body: ast.FieldRuleBody{
									ParamSingle: ptr("36"),
									Error:       "UUID must be exactly 36 characters",
								},
							},
						},
					},
					{
						Name: "name",
						Type: ast.FieldType{
							Base: &ast.FieldTypeBase{
								Named: ptr("string"),
							},
						},
						Rules: []*ast.FieldRule{
							{
								Name: "minlen",
								Body: ast.FieldRuleBody{
									ParamSingle: ptr("3"),
									Error:       "Name must be at least 3 characters long",
								},
							},
						},
					},
					{
						Name:     "description",
						Optional: true,
						Type: ast.FieldType{
							Base: &ast.FieldTypeBase{
								Named: ptr("string"),
							},
						},
					},
					{
						Name: "isActive",
						Type: ast.FieldType{
							Base: &ast.FieldTypeBase{
								Named: ptr("boolean"),
							},
						},
						Rules: []*ast.FieldRule{
							{
								Name: "equals",
								Body: ast.FieldRuleBody{
									ParamSingle: ptr("true"),
								},
							},
						},
					},
					{
						Name:     "parentId",
						Optional: true,
						Type: ast.FieldType{
							Base: &ast.FieldTypeBase{
								Named: ptr("string"),
							},
						},
						Rules: []*ast.FieldRule{
							{
								Name: "uuid",
							},
						},
					},
				},
			},
			{
				Docstring: "Product represents a sellable item in the store.\n\t\tProducts have complex validation rules and can be\n\t\tnested inside catalogs.",
				Name:      "Product",
				Fields: []*ast.Field{
					{
						Name: "id",
						Type: ast.FieldType{
							Base: &ast.FieldTypeBase{
								Named: ptr("string"),
							},
						},
						Rules: []*ast.FieldRule{
							{
								Name: "uuid",
							},
						},
					},
					{
						Name: "name",
						Type: ast.FieldType{
							Base: &ast.FieldTypeBase{
								Named: ptr("string"),
							},
						},
						Rules: []*ast.FieldRule{
							{
								Name: "minlen",
								Body: ast.FieldRuleBody{
									ParamSingle: ptr("2"),
								},
							},
							{
								Name: "maxlen",
								Body: ast.FieldRuleBody{
									ParamSingle: ptr("100"),
									Error:       "Name cannot exceed 100 characters",
								},
							},
						},
					},
					{
						Name: "price",
						Type: ast.FieldType{
							Base: &ast.FieldTypeBase{
								Named: ptr("float"),
							},
						},
						Rules: []*ast.FieldRule{
							{
								Name: "min",
								Body: ast.FieldRuleBody{
									ParamSingle: ptr("0.01"),
									Error:       "Price must be greater than zero",
								},
							},
						},
					},
					{
						Name: "stock",
						Type: ast.FieldType{
							Base: &ast.FieldTypeBase{
								Named: ptr("int"),
							},
						},
						Rules: []*ast.FieldRule{
							{
								Name: "min",
								Body: ast.FieldRuleBody{
									ParamSingle: ptr("0"),
								},
							},
							{
								Name: "range",
								Body: ast.FieldRuleBody{
									ParamList: []string{
										"0",
										"1000",
									},
									Error: "Stock must be between 0 and 1000",
								},
							},
						},
					},
					{
						Name: "category",
						Type: ast.FieldType{
							Base: &ast.FieldTypeBase{
								Named: ptr("Category"),
							},
						},
						Rules: []*ast.FieldRule{
							{
								Name: "validateCategory",
								Body: ast.FieldRuleBody{
									Error: "Invalid category custom message",
								},
							},
						},
					},
					{
						Name:     "tags",
						Optional: true,
						Type: ast.FieldType{
							Depth: 1,
							Base: &ast.FieldTypeBase{
								Named: ptr("string"),
							},
						},
						Rules: []*ast.FieldRule{
							{
								Name: "minlen",
								Body: ast.FieldRuleBody{
									ParamSingle: ptr("1"),
									Error:       "At least one tag is required",
								},
							},
							{
								Name: "maxlen",
								Body: ast.FieldRuleBody{
									ParamSingle: ptr("10"),
								},
							},
						},
					},
					{
						Name: "details",
						Type: ast.FieldType{
							Base: &ast.FieldTypeBase{
								Object: &ast.FieldTypeObject{
									Fields: []*ast.Field{
										{
											Name: "dimensions",
											Type: ast.FieldType{
												Base: &ast.FieldTypeBase{
													Object: &ast.FieldTypeObject{
														Fields: []*ast.Field{
															{
																Name: "width",
																Type: ast.FieldType{
																	Base: &ast.FieldTypeBase{
																		Named: ptr("float"),
																	},
																},
																Rules: []*ast.FieldRule{
																	{
																		Name: "min",
																		Body: ast.FieldRuleBody{
																			ParamSingle: ptr("0.0"),
																			Error:       "Width cannot be negative",
																		},
																	},
																},
															},
															{
																Name: "height",
																Type: ast.FieldType{
																	Base: &ast.FieldTypeBase{
																		Named: ptr("float"),
																	},
																},
																Rules: []*ast.FieldRule{
																	{
																		Name: "min",
																		Body: ast.FieldRuleBody{
																			ParamSingle: ptr("0.0"),
																		},
																	},
																},
															},
															{
																Name:     "depth",
																Optional: true,
																Type: ast.FieldType{
																	Base: &ast.FieldTypeBase{
																		Named: ptr("float"),
																	},
																},
															},
														},
													},
												},
											},
										},
										{
											Name:     "weight",
											Optional: true,
											Type: ast.FieldType{
												Base: &ast.FieldTypeBase{
													Named: ptr("float"),
												},
											},
										},
										{
											Name: "colors",
											Type: ast.FieldType{
												Depth: 1,
												Base: &ast.FieldTypeBase{
													Named: ptr("string"),
												},
											},
											Rules: []*ast.FieldRule{
												{
													Name: "enum",
													Body: ast.FieldRuleBody{
														ParamList: []string{
															"red",
															"green",
															"blue",
															"black",
															"white",
														},
														Error: "Color must be one of the allowed values",
													},
												},
											},
										},
										{
											Name:     "attributes",
											Optional: true,
											Type: ast.FieldType{
												Depth: 1,
												Base: &ast.FieldTypeBase{
													Object: &ast.FieldTypeObject{
														Fields: []*ast.Field{
															{
																Name: "name",
																Type: ast.FieldType{
																	Base: &ast.FieldTypeBase{
																		Named: ptr("string"),
																	},
																},
															},
															{
																Name: "value",
																Type: ast.FieldType{
																	Base: &ast.FieldTypeBase{
																		Named: ptr("string"),
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
						Name: "variations",
						Type: ast.FieldType{
							Depth: 1,
							Base: &ast.FieldTypeBase{
								Object: &ast.FieldTypeObject{
									Fields: []*ast.Field{
										{
											Name: "sku",
											Type: ast.FieldType{
												Base: &ast.FieldTypeBase{
													Named: ptr("string"),
												},
											},
										},
										{
											Name: "price",
											Type: ast.FieldType{
												Base: &ast.FieldTypeBase{
													Named: ptr("float"),
												},
											},
											Rules: []*ast.FieldRule{
												{
													Name: "min",
													Body: ast.FieldRuleBody{
														ParamSingle: ptr("0.01"),
														Error:       "Variation price must be greater than zero",
													},
												},
											},
										},
										{
											Name: "attributes",
											Type: ast.FieldType{
												Depth: 1,
												Base: &ast.FieldTypeBase{
													Object: &ast.FieldTypeObject{
														Fields: []*ast.Field{
															{
																Name: "name",
																Type: ast.FieldType{
																	Base: &ast.FieldTypeBase{
																		Named: ptr("string"),
																	},
																},
															},
															{
																Name: "value",
																Type: ast.FieldType{
																	Base: &ast.FieldTypeBase{
																		Named: ptr("string"),
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
		Procs: []*ast.ProcDecl{
			{
				Docstring: "GetCategory retrieves a category by its ID.\n\t\tThis is a basic read operation.",
				Name:      "GetCategory",
				Body: ast.ProcDeclBody{
					Input: []*ast.Field{
						{
							Name: "id",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{
									Named: ptr("string"),
								},
							},
							Rules: []*ast.FieldRule{
								{
									Name: "uuid",
									Body: ast.FieldRuleBody{
										Error: "Category ID must be a valid UUID",
									},
								},
							},
						},
					},
					Output: []*ast.Field{
						{
							Name: "category",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{
									Named: ptr("Category"),
								},
							},
						},
						{
							Name: "exists",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{
									Named: ptr("boolean"),
								},
							},
						},
					},
					Meta: []*ast.ProcDeclBodyMetaKV{
						{
							Key:   "cache",
							Value: "true",
						},
						{
							Key:   "cacheTime",
							Value: "300",
						},
						{
							Key:   "requiresAuth",
							Value: "false",
						},
						{
							Key:   "apiVersion",
							Value: "1.0.0",
						},
					},
				},
			},
			{
				Docstring: "CreateProduct adds a new product to the catalog.\n\t\tThis procedure handles complex validation and returns\n\t\tdetailed success information.",
				Name:      "CreateProduct",
				Body: ast.ProcDeclBody{
					Input: []*ast.Field{
						{
							Name: "product",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{
									Named: ptr("Product"),
								},
							},
						},
						{
							Name:     "options",
							Optional: true,
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{
									Object: &ast.FieldTypeObject{
										Fields: []*ast.Field{
											{
												Name: "draft",
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{
														Named: ptr("boolean"),
													},
												},
											},
											{
												Name: "notify",
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{
														Named: ptr("boolean"),
													},
												},
											},
											{
												Name:     "scheduledFor",
												Optional: true,
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{
														Named: ptr("string"),
													},
												},
												Rules: []*ast.FieldRule{
													{
														Name: "iso8601",
														Body: ast.FieldRuleBody{
															Error: "Must be a valid ISO8601 date",
														},
													},
												},
											},
											{
												Name:     "tags",
												Optional: true,
												Type: ast.FieldType{
													Depth: 1,
													Base: &ast.FieldTypeBase{
														Named: ptr("string"),
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Name: "validation",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{
									Object: &ast.FieldTypeObject{
										Fields: []*ast.Field{
											{
												Name:     "skipValidation",
												Optional: true,
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{
														Named: ptr("boolean"),
													},
												},
											},
											{
												Name:     "customRules",
												Optional: true,
												Type: ast.FieldType{
													Depth: 1,
													Base: &ast.FieldTypeBase{
														Object: &ast.FieldTypeObject{
															Fields: []*ast.Field{
																{
																	Name: "name",
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: ptr("string"),
																		},
																	},
																},
																{
																	Name: "severity",
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: ptr("int"),
																		},
																	},
																	Rules: []*ast.FieldRule{
																		{
																			Name: "enum",
																			Body: ast.FieldRuleBody{
																				ParamList: []string{
																					"1",
																					"2",
																					"3",
																				},
																				Error: "Severity must be 1, 2, or 3",
																			},
																		},
																	},
																},
																{
																	Name: "message",
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: ptr("string"),
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
					Output: []*ast.Field{
						{
							Name: "success",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{
									Named: ptr("boolean"),
								},
							},
						},
						{
							Name: "productId",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{
									Named: ptr("string"),
								},
							},
							Rules: []*ast.FieldRule{
								{
									Name: "uuid",
									Body: ast.FieldRuleBody{
										Error: "Product ID must be a valid UUID",
									},
								},
							},
						},
						{
							Name:     "errors",
							Optional: true,
							Type: ast.FieldType{
								Depth: 1,
								Base: &ast.FieldTypeBase{
									Object: &ast.FieldTypeObject{
										Fields: []*ast.Field{
											{
												Name: "code",
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{
														Named: ptr("int"),
													},
												},
											},
											{
												Name: "message",
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{
														Named: ptr("string"),
													},
												},
											},
											{
												Name:     "field",
												Optional: true,
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{
														Named: ptr("string"),
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Name: "analytics",
							Type: ast.FieldType{
								Base: &ast.FieldTypeBase{
									Object: &ast.FieldTypeObject{
										Fields: []*ast.Field{
											{
												Name: "duration",
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{
														Named: ptr("float"),
													},
												},
											},
											{
												Name: "processingSteps",
												Type: ast.FieldType{
													Depth: 1,
													Base: &ast.FieldTypeBase{
														Object: &ast.FieldTypeObject{
															Fields: []*ast.Field{
																{
																	Name: "name",
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: ptr("string"),
																		},
																	},
																},
																{
																	Name: "duration",
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: ptr("float"),
																		},
																	},
																},
																{
																	Name: "success",
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: ptr("boolean"),
																		},
																	},
																},
															},
														},
													},
												},
											},
											{
												Name: "serverInfo",
												Type: ast.FieldType{
													Base: &ast.FieldTypeBase{
														Object: &ast.FieldTypeObject{
															Fields: []*ast.Field{
																{
																	Name: "id",
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: ptr("string"),
																		},
																	},
																},
																{
																	Name: "region",
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: ptr("string"),
																		},
																	},
																},
																{
																	Name: "load",
																	Type: ast.FieldType{
																		Base: &ast.FieldTypeBase{
																			Named: ptr("float"),
																		},
																	},
																	Rules: []*ast.FieldRule{
																		{
																			Name: "min",
																			Body: ast.FieldRuleBody{
																				ParamSingle: ptr("0.0"),
																			},
																		},
																		{
																			Name: "max",
																			Body: ast.FieldRuleBody{
																				ParamSingle: ptr("1.0"),
																				Error:       "Load factor cannot exceed 1.0",
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
					Meta: []*ast.ProcDeclBodyMetaKV{
						{
							Key:   "auth",
							Value: "required",
						},
						{
							Key:   "roles",
							Value: "admin,product-manager",
						},
						{
							Key:   "rateLimit",
							Value: "100",
						},
						{
							Key:   "timeout",
							Value: "30.5",
						},
						{
							Key:   "audit",
							Value: "true",
						},
						{
							Key:   "apiVersion",
							Value: "1.2.0",
						},
					},
				},
			},
		},
	}

	equalNoPos(t, expected, parsed)
}
