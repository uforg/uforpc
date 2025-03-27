package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/urpc/ast"
	"github.com/uforg/uforpc/internal/urpc/lexer"
)

func TestParser(t *testing.T) {
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
		require.Equal(t, expected, schema)
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

	t.Run("Parse type declaration basic", func(t *testing.T) {
		input := `
			type User {}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDeclaration{
				{
					Name: "User",
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
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
			Types: []ast.TypeDeclaration{
				{
					Doc:  "Product type documentation",
					Name: "Product",
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
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
			Types: []ast.TypeDeclaration{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "name",
							Optional: false,
							Type:     &ast.TypeString{},
						},
						{
							Name:     "age",
							Optional: true,
							Type:     &ast.TypeInt{},
						},
						{
							Name:     "height",
							Optional: false,
							Type:     &ast.TypeFloat{},
						},
						{
							Name:     "isActive",
							Optional: false,
							Type:     &ast.TypeBoolean{},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
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
			Types: []ast.TypeDeclaration{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "name",
							Optional: false,
							Type:     &ast.TypeString{},
						},
						{
							Name:     "address",
							Optional: true,
							Type:     &ast.TypeCustom{Name: "Address"},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
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
			Types: []ast.TypeDeclaration{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "objField",
							Optional: false,
							Type: &ast.TypeObject{
								Fields: []ast.Field{
									{
										Name:     "name",
										Optional: false,
										Type:     &ast.TypeString{},
									},
									{
										Name:     "age",
										Optional: false,
										Type:     &ast.TypeInt{},
									},
								},
							},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
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
			Types: []ast.TypeDeclaration{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "objField",
							Optional: false,
							Type: &ast.TypeObject{
								Fields: []ast.Field{
									{
										Name:     "name",
										Optional: false,
										Type:     &ast.TypeString{},
									},
									{
										Name:     "age",
										Optional: false,
										Type:     &ast.TypeInt{},
									},
									{
										Name:     "address",
										Optional: false,
										Type: &ast.TypeObject{
											Fields: []ast.Field{
												{
													Name:     "street",
													Optional: false,
													Type:     &ast.TypeString{},
												},
												{
													Name:     "city",
													Optional: false,
													Type:     &ast.TypeString{},
												},
												{
													Name:     "zip",
													Optional: false,
													Type:     &ast.TypeString{},
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
		require.Equal(t, expected, schema)
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
			Types: []ast.TypeDeclaration{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "objField",
							Optional: false,
							Type: &ast.TypeArray{
								ArrayType: &ast.TypeObject{
									Fields: []ast.Field{
										{
											Name:     "name",
											Optional: false,
											Type:     &ast.TypeString{},
										},
										{
											Name:     "age",
											Optional: false,
											Type:     &ast.TypeInt{},
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
		require.Equal(t, expected, schema)
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
			Types: []ast.TypeDeclaration{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "arrayField",
							Optional: false,
							Type:     &ast.TypeArray{ArrayType: &ast.TypeString{}},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
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
			Types: []ast.TypeDeclaration{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "stringField",
							Optional: false,
							Type:     &ast.TypeString{},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleWithValue{
									RuleName:     "equals",
									Value:        "Foo",
									ValueType:    ast.ValidationRuleValueTypeString,
									ErrorMessage: "",
								},
								&ast.ValidationRuleWithValue{
									RuleName:     "contains",
									Value:        "Bar",
									ValueType:    ast.ValidationRuleValueTypeString,
									ErrorMessage: "",
								},
								&ast.ValidationRuleWithValue{
									RuleName:     "minlen",
									Value:        "3",
									ValueType:    ast.ValidationRuleValueTypeInt,
									ErrorMessage: "",
								},
								&ast.ValidationRuleWithValue{
									RuleName:     "maxlen",
									Value:        "100",
									ValueType:    ast.ValidationRuleValueTypeInt,
									ErrorMessage: "",
								},
								&ast.ValidationRuleWithArray{
									RuleName:     "enum",
									Values:       []string{"Foo", "Bar"},
									ValueType:    ast.ValidationRuleValueTypeString,
									ErrorMessage: "",
								},
								&ast.ValidationRuleSimple{
									RuleName:     "email",
									ErrorMessage: "",
								},
								&ast.ValidationRuleSimple{
									RuleName:     "iso8601",
									ErrorMessage: "",
								},
								&ast.ValidationRuleSimple{
									RuleName:     "uuid",
									ErrorMessage: "",
								},
								&ast.ValidationRuleSimple{
									RuleName:     "json",
									ErrorMessage: "",
								},
								&ast.ValidationRuleSimple{
									RuleName:     "lowercase",
									ErrorMessage: "",
								},
								&ast.ValidationRuleSimple{
									RuleName:     "uppercase",
									ErrorMessage: "",
								},
							},
						},
						{
							Name:     "intField",
							Optional: false,
							Type:     &ast.TypeInt{},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleWithValue{
									RuleName:     "equals",
									Value:        "1",
									ValueType:    ast.ValidationRuleValueTypeInt,
									ErrorMessage: "",
								},
								&ast.ValidationRuleWithValue{
									RuleName:     "min",
									Value:        "0",
									ValueType:    ast.ValidationRuleValueTypeInt,
									ErrorMessage: "",
								},
								&ast.ValidationRuleWithValue{
									RuleName:     "max",
									Value:        "100",
									ValueType:    ast.ValidationRuleValueTypeInt,
									ErrorMessage: "",
								},
								&ast.ValidationRuleWithArray{
									RuleName:     "enum",
									Values:       []string{"1", "2", "3"},
									ValueType:    ast.ValidationRuleValueTypeInt,
									ErrorMessage: "",
								},
							},
						},
						{
							Name:     "floatField",
							Optional: false,
							Type:     &ast.TypeFloat{},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleWithValue{
									RuleName:     "min",
									Value:        "0.0",
									ValueType:    ast.ValidationRuleValueTypeFloat,
									ErrorMessage: "",
								},
								&ast.ValidationRuleWithValue{
									RuleName:     "max",
									Value:        "100.0",
									ValueType:    ast.ValidationRuleValueTypeFloat,
									ErrorMessage: "",
								},
							},
						},
						{
							Name:     "booleanField",
							Optional: false,
							Type:     &ast.TypeBoolean{},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleWithValue{
									RuleName:     "equals",
									Value:        "true",
									ValueType:    ast.ValidationRuleValueTypeBoolean,
									ErrorMessage: "",
								},
							},
						},
						{
							Name:     "arrayField",
							Optional: false,
							Type:     &ast.TypeArray{ArrayType: &ast.TypeString{}},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleWithValue{
									RuleName:     "minlen",
									Value:        "1",
									ValueType:    ast.ValidationRuleValueTypeInt,
									ErrorMessage: "",
								},
								&ast.ValidationRuleWithValue{
									RuleName:     "maxlen",
									Value:        "100",
									ValueType:    ast.ValidationRuleValueTypeInt,
									ErrorMessage: "",
								},
							},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
	})

	t.Run("Parse procedure declaration basic", func(t *testing.T) {
		input := `
			proc CreateUser {}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Procedures: []ast.ProcDeclaration{
				{
					Name:     "CreateUser",
					Input:    ast.ProcInput{},
					Output:   ast.ProcOutput{},
					Metadata: ast.ProcMeta{},
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
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
			Procedures: []ast.ProcDeclaration{
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
		require.Equal(t, expected, schema)
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
			Procedures: []ast.ProcDeclaration{
				{
					Name:     "CreateUser",
					Input:    ast.ProcInput{},
					Output:   ast.ProcOutput{},
					Metadata: ast.ProcMeta{},
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
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
			Procedures: []ast.ProcDeclaration{
				{
					Doc:  "Creates a product with the given stock and returns the product id.",
					Name: "CreateProduct",
					Input: ast.ProcInput{
						Fields: []ast.Field{
							{
								Name:     "product",
								Optional: false,
								Type:     &ast.TypeCustom{Name: "Product"},
							},
							{
								Name:     "stock",
								Optional: false,
								Type:     &ast.TypeInt{},
							},
						},
					},
					Output: ast.ProcOutput{
						Fields: []ast.Field{
							{
								Name:     "productId",
								Optional: false,
								Type:     &ast.TypeString{},
							},
						},
					},
					Metadata: ast.ProcMeta{
						Entries: []ast.ProcMetaKV{
							{
								Type:  ast.ProcMetaValueTypeString,
								Key:   "versionNumber",
								Value: "1.0.0",
							},
							{
								Type:  ast.ProcMetaValueTypeInt,
								Key:   "maxRetries",
								Value: "3",
							},
							{
								Type:  ast.ProcMetaValueTypeFloat,
								Key:   "waitMinutes",
								Value: "10.5",
							},
							{
								Type:  ast.ProcMetaValueTypeBoolean,
								Key:   "audit",
								Value: "true",
							},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
	})

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
			Types: []ast.TypeDeclaration{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "name",
							Optional: false,
							Type:     &ast.TypeString{},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleSimple{
									RuleName:     "required",
									ErrorMessage: "Name is required",
								},
							},
						},
						{
							Name:     "age",
							Optional: false,
							Type:     &ast.TypeInt{},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleWithValue{
									RuleName:     "min",
									Value:        "18",
									ValueType:    ast.ValidationRuleValueTypeInt,
									ErrorMessage: "Must be an adult",
								},
							},
						},
						{
							Name:     "email",
							Optional: false,
							Type:     &ast.TypeString{},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleSimple{
									RuleName:     "email",
									ErrorMessage: "Invalid email format",
								},
							},
						},
						{
							Name:     "options",
							Optional: false,
							Type:     &ast.TypeArray{ArrayType: &ast.TypeString{}},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleWithArray{
									RuleName:     "enum",
									Values:       []string{"a", "b", "c"},
									ValueType:    ast.ValidationRuleValueTypeString,
									ErrorMessage: "Invalid option selected",
								},
							},
						},
						{
							Name:     "tag",
							Optional: false,
							Type:     &ast.TypeString{},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleWithValue{
									RuleName:     "pattern",
									Value:        "^[a-z]+$",
									ValueType:    ast.ValidationRuleValueTypeString,
									ErrorMessage: "Only lowercase letters allowed",
								},
							},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
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
			Types: []ast.TypeDeclaration{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "name",
							Optional: false,
							Type:     &ast.TypeString{},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleSimple{
									RuleName:     "required",
									ErrorMessage: "This field cannot be empty",
								},
							},
						},
						{
							Name:     "email",
							Optional: false,
							Type:     &ast.TypeString{},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleSimple{
									RuleName:     "email",
									ErrorMessage: "Please enter a valid email address",
								},
							},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
	})
}

func TestParserFullExample(t *testing.T) {
	input := `
		// Version declaration
		version 1

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
			category: Category
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
		Types: []ast.TypeDeclaration{
			{
				Name: "Category",
				Doc:  "Category represents a product category in the system.\n\t\tThis type is used across the catalog module.",
				Fields: []ast.Field{
					{
						Name:     "id",
						Optional: false,
						Type:     &ast.TypeString{},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleSimple{
								RuleName:     "uuid",
								ErrorMessage: "Must be a valid UUID",
							},
							&ast.ValidationRuleWithValue{
								RuleName:     "minlen",
								Value:        "36",
								ValueType:    ast.ValidationRuleValueTypeInt,
								ErrorMessage: "",
							},
							&ast.ValidationRuleWithValue{
								RuleName:     "maxlen",
								Value:        "36",
								ValueType:    ast.ValidationRuleValueTypeInt,
								ErrorMessage: "UUID must be exactly 36 characters",
							},
						},
					},
					{
						Name:     "name",
						Optional: false,
						Type:     &ast.TypeString{},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleWithValue{
								RuleName:     "minlen",
								Value:        "3",
								ValueType:    ast.ValidationRuleValueTypeInt,
								ErrorMessage: "Name must be at least 3 characters long",
							},
						},
					},
					{
						Name:     "description",
						Optional: true,
						Type:     &ast.TypeString{},
					},
					{
						Name:     "isActive",
						Optional: false,
						Type:     &ast.TypeBoolean{},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleWithValue{
								RuleName:     "equals",
								Value:        "true",
								ValueType:    ast.ValidationRuleValueTypeBoolean,
								ErrorMessage: "",
							},
						},
					},
					{
						Name:     "parentId",
						Optional: true,
						Type:     &ast.TypeString{},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleSimple{
								RuleName:     "uuid",
								ErrorMessage: "",
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
						Type:     &ast.TypeString{},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleSimple{
								RuleName:     "uuid",
								ErrorMessage: "",
							},
						},
					},
					{
						Name:     "name",
						Optional: false,
						Type:     &ast.TypeString{},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleWithValue{
								RuleName:     "minlen",
								Value:        "2",
								ValueType:    ast.ValidationRuleValueTypeInt,
								ErrorMessage: "",
							},
							&ast.ValidationRuleWithValue{
								RuleName:     "maxlen",
								Value:        "100",
								ValueType:    ast.ValidationRuleValueTypeInt,
								ErrorMessage: "Name cannot exceed 100 characters",
							},
						},
					},
					{
						Name:     "price",
						Optional: false,
						Type:     &ast.TypeFloat{},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleWithValue{
								RuleName:     "min",
								Value:        "0.01",
								ValueType:    ast.ValidationRuleValueTypeFloat,
								ErrorMessage: "Price must be greater than zero",
							},
						},
					},
					{
						Name:     "stock",
						Optional: false,
						Type:     &ast.TypeInt{},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleWithValue{
								RuleName:     "min",
								Value:        "0",
								ValueType:    ast.ValidationRuleValueTypeInt,
								ErrorMessage: "",
							},
						},
					},
					{
						Name:     "category",
						Optional: false,
						Type:     &ast.TypeCustom{Name: "Category"},
					},
					{
						Name:     "tags",
						Optional: true,
						Type:     &ast.TypeArray{ArrayType: &ast.TypeString{}},
						ValidationRules: []ast.ValidationRule{
							&ast.ValidationRuleWithValue{
								RuleName:     "minlen",
								Value:        "1",
								ValueType:    ast.ValidationRuleValueTypeInt,
								ErrorMessage: "At least one tag is required",
							},
							&ast.ValidationRuleWithValue{
								RuleName:     "maxlen",
								Value:        "10",
								ValueType:    ast.ValidationRuleValueTypeInt,
								ErrorMessage: "",
							},
						},
					},
					{
						Name:     "details",
						Optional: false,
						Type: &ast.TypeObject{
							Fields: []ast.Field{
								{
									Name:     "dimensions",
									Optional: false,
									Type: &ast.TypeObject{
										Fields: []ast.Field{
											{
												Name:     "width",
												Optional: false,
												Type:     &ast.TypeFloat{},
												ValidationRules: []ast.ValidationRule{
													&ast.ValidationRuleWithValue{
														RuleName:     "min",
														Value:        "0.0",
														ValueType:    ast.ValidationRuleValueTypeFloat,
														ErrorMessage: "Width cannot be negative",
													},
												},
											},
											{
												Name:     "height",
												Optional: false,
												Type:     &ast.TypeFloat{},
												ValidationRules: []ast.ValidationRule{
													&ast.ValidationRuleWithValue{
														RuleName:     "min",
														Value:        "0.0",
														ValueType:    ast.ValidationRuleValueTypeFloat,
														ErrorMessage: "",
													},
												},
											},
											{
												Name:     "depth",
												Optional: true,
												Type:     &ast.TypeFloat{},
											},
										},
									},
								},
								{
									Name:     "weight",
									Optional: true,
									Type:     &ast.TypeFloat{},
								},
								{
									Name:     "colors",
									Optional: false,
									Type:     &ast.TypeArray{ArrayType: &ast.TypeString{}},
									ValidationRules: []ast.ValidationRule{
										&ast.ValidationRuleWithArray{
											RuleName:     "enum",
											Values:       []string{"red", "green", "blue", "black", "white"},
											ValueType:    ast.ValidationRuleValueTypeString,
											ErrorMessage: "Color must be one of the allowed values",
										},
									},
								},
								{
									Name:     "attributes",
									Optional: true,
									Type: &ast.TypeArray{
										ArrayType: &ast.TypeObject{
											Fields: []ast.Field{
												{
													Name:     "name",
													Optional: false,
													Type:     &ast.TypeString{},
												},
												{
													Name:     "value",
													Optional: false,
													Type:     &ast.TypeString{},
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
						Type: &ast.TypeArray{
							ArrayType: &ast.TypeObject{
								Fields: []ast.Field{
									{
										Name:     "sku",
										Optional: false,
										Type:     &ast.TypeString{},
									},
									{
										Name:     "price",
										Optional: false,
										Type:     &ast.TypeFloat{},
										ValidationRules: []ast.ValidationRule{
											&ast.ValidationRuleWithValue{
												RuleName:     "min",
												Value:        "0.01",
												ValueType:    ast.ValidationRuleValueTypeFloat,
												ErrorMessage: "Variation price must be greater than zero",
											},
										},
									},
									{
										Name:     "attributes",
										Optional: false,
										Type: &ast.TypeArray{
											ArrayType: &ast.TypeObject{
												Fields: []ast.Field{
													{
														Name:     "name",
														Optional: false,
														Type:     &ast.TypeString{},
													},
													{
														Name:     "value",
														Optional: false,
														Type:     &ast.TypeString{},
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
		Procedures: []ast.ProcDeclaration{
			{
				Name: "GetCategory",
				Doc:  "GetCategory retrieves a category by its ID.\n\t\tThis is a basic read operation.",
				Input: ast.ProcInput{
					Fields: []ast.Field{
						{
							Name:     "id",
							Optional: false,
							Type:     &ast.TypeString{},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleSimple{
									RuleName:     "uuid",
									ErrorMessage: "Category ID must be a valid UUID",
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
							Type:     &ast.TypeCustom{Name: "Category"},
						},
						{
							Name:     "exists",
							Optional: false,
							Type:     &ast.TypeBoolean{},
						},
					},
				},
				Metadata: ast.ProcMeta{
					Entries: []ast.ProcMetaKV{
						{
							Type:  ast.ProcMetaValueTypeBoolean,
							Key:   "cache",
							Value: "true",
						},
						{
							Type:  ast.ProcMetaValueTypeInt,
							Key:   "cacheTime",
							Value: "300",
						},
						{
							Type:  ast.ProcMetaValueTypeBoolean,
							Key:   "requiresAuth",
							Value: "false",
						},
						{
							Type:  ast.ProcMetaValueTypeString,
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
							Type:     &ast.TypeCustom{Name: "Product"},
						},
						{
							Name:     "options",
							Optional: true,
							Type: &ast.TypeObject{
								Fields: []ast.Field{
									{
										Name:     "draft",
										Optional: false,
										Type:     &ast.TypeBoolean{},
									},
									{
										Name:     "notify",
										Optional: false,
										Type:     &ast.TypeBoolean{},
									},
									{
										Name:     "scheduledFor",
										Optional: true,
										Type:     &ast.TypeString{},
										ValidationRules: []ast.ValidationRule{
											&ast.ValidationRuleSimple{
												RuleName:     "iso8601",
												ErrorMessage: "Must be a valid ISO8601 date",
											},
										},
									},
									{
										Name:     "tags",
										Optional: true,
										Type:     &ast.TypeArray{ArrayType: &ast.TypeString{}},
									},
								},
							},
						},
						{
							Name:     "validation",
							Optional: false,
							Type: &ast.TypeObject{
								Fields: []ast.Field{
									{
										Name:     "skipValidation",
										Optional: true,
										Type:     &ast.TypeBoolean{},
									},
									{
										Name:     "customRules",
										Optional: true,
										Type: &ast.TypeArray{
											ArrayType: &ast.TypeObject{
												Fields: []ast.Field{
													{
														Name:     "name",
														Optional: false,
														Type:     &ast.TypeString{},
													},
													{
														Name:     "severity",
														Optional: false,
														Type:     &ast.TypeInt{},
														ValidationRules: []ast.ValidationRule{
															&ast.ValidationRuleWithArray{
																RuleName:     "enum",
																Values:       []string{"1", "2", "3"},
																ValueType:    ast.ValidationRuleValueTypeInt,
																ErrorMessage: "Severity must be 1, 2, or 3",
															},
														},
													},
													{
														Name:     "message",
														Optional: false,
														Type:     &ast.TypeString{},
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
							Type:     &ast.TypeBoolean{},
						},
						{
							Name:     "productId",
							Optional: false,
							Type:     &ast.TypeString{},
							ValidationRules: []ast.ValidationRule{
								&ast.ValidationRuleSimple{
									RuleName:     "uuid",
									ErrorMessage: "Product ID must be a valid UUID",
								},
							},
						},
						{
							Name:     "errors",
							Optional: true,
							Type: &ast.TypeArray{
								ArrayType: &ast.TypeObject{
									Fields: []ast.Field{
										{
											Name:     "code",
											Optional: false,
											Type:     &ast.TypeInt{},
										},
										{
											Name:     "message",
											Optional: false,
											Type:     &ast.TypeString{},
										},
										{
											Name:     "field",
											Optional: true,
											Type:     &ast.TypeString{},
										},
									},
								},
							},
						},
						{
							Name:     "analytics",
							Optional: false,
							Type: &ast.TypeObject{
								Fields: []ast.Field{
									{
										Name:     "duration",
										Optional: false,
										Type:     &ast.TypeFloat{},
									},
									{
										Name:     "processingSteps",
										Optional: false,
										Type: &ast.TypeArray{
											ArrayType: &ast.TypeObject{
												Fields: []ast.Field{
													{
														Name:     "name",
														Optional: false,
														Type:     &ast.TypeString{},
													},
													{
														Name:     "duration",
														Optional: false,
														Type:     &ast.TypeFloat{},
													},
													{
														Name:     "success",
														Optional: false,
														Type:     &ast.TypeBoolean{},
													},
												},
											},
										},
									},
									{
										Name:     "serverInfo",
										Optional: false,
										Type: &ast.TypeObject{
											Fields: []ast.Field{
												{
													Name:     "id",
													Optional: false,
													Type:     &ast.TypeString{},
												},
												{
													Name:     "region",
													Optional: false,
													Type:     &ast.TypeString{},
												},
												{
													Name:     "load",
													Optional: false,
													Type:     &ast.TypeFloat{},
													ValidationRules: []ast.ValidationRule{
														&ast.ValidationRuleWithValue{
															RuleName:     "min",
															Value:        "0.0",
															ValueType:    ast.ValidationRuleValueTypeFloat,
															ErrorMessage: "",
														},
														&ast.ValidationRuleWithValue{
															RuleName:     "max",
															Value:        "1.0",
															ValueType:    ast.ValidationRuleValueTypeFloat,
															ErrorMessage: "Load factor cannot exceed 1.0",
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
							Type:  ast.ProcMetaValueTypeString,
							Key:   "auth",
							Value: "required",
						},
						{
							Type:  ast.ProcMetaValueTypeString,
							Key:   "roles",
							Value: "admin,product-manager",
						},
						{
							Type:  ast.ProcMetaValueTypeInt,
							Key:   "rateLimit",
							Value: "100",
						},
						{
							Type:  ast.ProcMetaValueTypeFloat,
							Key:   "timeout",
							Value: "30.5",
						},
						{
							Type:  ast.ProcMetaValueTypeBoolean,
							Key:   "audit",
							Value: "true",
						},
						{
							Type:  ast.ProcMetaValueTypeString,
							Key:   "apiVersion",
							Value: "1.2.0",
						},
					},
				},
			},
		},
	}

	require.NoError(t, err)
	require.Equal(t, expected, schema)
}
