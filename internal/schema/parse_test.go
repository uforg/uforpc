package schema_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/schema"
)

func TestParseSchema(t *testing.T) {
	t.Run("valid schema", func(t *testing.T) {
		testCases, err := getSchemasFromFS(validSchemasFS, "examples/valid")
		require.NoError(t, err)

		for _, testCase := range testCases {
			t.Run(testCase.fileName, func(t *testing.T) {
				parsedSchema, err := schema.ParseSchema(testCase.schema)

				assert.NoError(t, err)
				assert.Equal(t, 1, parsedSchema.Version)
				assert.NotNil(t, parsedSchema.Types)
				assert.NotNil(t, parsedSchema.Procedures)
			})
		}
	})

	t.Run("invalid schema", func(t *testing.T) {
		testCases, err := getSchemasFromFS(invalidSchemasFS, "examples/invalid")
		require.NoError(t, err)

		for _, testCase := range testCases {
			t.Run(testCase.fileName, func(t *testing.T) {
				parsedSchema, err := schema.ParseSchema(testCase.schema)

				assert.Error(t, err)
				assert.Equal(t, schema.Schema{}, parsedSchema)
			})
		}
	})

	t.Run("custom schema", func(t *testing.T) {
		customSchema := `{
			"version": 1,
			"types": {
				"User": {
					"type": "object",
					"fields": {
						"id": {
							"type": "string",
							"rules": {
								"uuid": {}
							}
						},
						"name": {
							"type": "string",
							"rules": {
								"minLen": {
									"value": 3,
									"errorMessage": "Name must be at least 3 characters"
								}
							}
						},
						"age": {
							"type": "int",
							"rules": {
								"min": {
									"value": 18,
									"errorMessage": "Must be at least 18 years old"
								}
							}
						},
						"isActive": {
							"type": "boolean"
						}
					}
				},
				"Address": {
					"type": "object",
					"fields": {
						"street": {
							"type": "string"
						},
						"city": {
							"type": "string"
						}
					}
				}
			},
			"procedures": {
				"GetUser": {
					"type": "query",
					"input": {
						"type": "object",
						"fields": {
							"id": {
								"type": "string",
								"rules": {
									"uuid": {}
								}
							}
						}
					},
					"output": {
						"type": "User"
					}
				},
				"CreateUser": {
					"type": "mutation",
					"input": {
						"type": "object",
						"fields": {
							"name": {
								"type": "string"
							},
							"age": {
								"type": "int"
							},
							"addresses": {
								"type": "array",
								"arrayType": {
									"type": "Address"
								}
							}
						}
					},
					"output": {
						"type": "User"
					}
				}
			}
		}`

		parsedSchema, err := schema.ParseSchema(customSchema)

		require.NoError(t, err)
		assert.Equal(t, 1, parsedSchema.Version)

		assert.Len(t, parsedSchema.Types, 2)
		assert.Contains(t, parsedSchema.Types, "User")
		assert.Contains(t, parsedSchema.Types, "Address")

		assert.Len(t, parsedSchema.Procedures, 2)
		assert.Contains(t, parsedSchema.Procedures, "GetUser")
		assert.Contains(t, parsedSchema.Procedures, "CreateUser")

		getUserProc := parsedSchema.Procedures["GetUser"]
		assert.Equal(t, schema.ProcedureTypeQuery, getUserProc.Type)

		createUserProc := parsedSchema.Procedures["CreateUser"]
		assert.Equal(t, schema.ProcedureTypeMutation, createUserProc.Type)

		userType := parsedSchema.Types["User"]
		assert.Equal(t, "object", userType.Type)
		assert.Len(t, userType.Fields, 4)
		assert.Contains(t, userType.Fields, "id")
		assert.Contains(t, userType.Fields, "name")
		assert.Contains(t, userType.Fields, "age")
		assert.Contains(t, userType.Fields, "isActive")

		nameField := userType.Fields["name"]
		stringRules, ok := nameField.ProcessedRules.(schema.StringRules)
		assert.True(t, ok)
		assert.Equal(t, 3, stringRules.MinLen.Value)
		assert.Equal(t, "Name must be at least 3 characters", stringRules.MinLen.ErrorMessage)

		addressesField := createUserProc.Input.Fields["addresses"]
		assert.Equal(t, "array", addressesField.Type)
		assert.NotNil(t, addressesField.ArrayType)
		assert.Equal(t, "Address", addressesField.ArrayType.Type)
	})

	t.Run("schema with undefined custom type", func(t *testing.T) {
		schemaWithUndefinedType := `{
			"version": 1,
			"types": {
				"User": {
					"type": "object",
					"fields": {
						"id": {
							"type": "string"
						},
						"profile": {
							"type": "Profile"
						}
					}
				}
			},
			"procedures": {
				"GetUser": {
					"type": "query",
					"output": {
						"type": "User"
					}
				}
			}
		}`

		_, err := schema.ParseSchema(schemaWithUndefinedType)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "undefined custom type: Profile")
	})

	t.Run("schema with circular reference", func(t *testing.T) {
		schemaWithCircularReference := `{
			"version": 1,
			"types": {
				"User": {
					"type": "object",
					"fields": {
						"id": {
							"type": "string"
						},
						"friend": {
							"type": "User"
						}
					}
				}
			},
			"procedures": {
				"GetUser": {
					"type": "query",
					"output": {
						"type": "User"
					}
				}
			}
		}`

		parsedSchema, err := schema.ParseSchema(schemaWithCircularReference)
		assert.NoError(t, err)
		assert.NotEqual(t, schema.Schema{}, parsedSchema)

		// Circular references are allowed and should parse correctly
		userType := parsedSchema.Types["User"]
		friendField := userType.Fields["friend"]
		assert.Equal(t, "User", friendField.Type)
	})

	t.Run("schema with nested array of custom type", func(t *testing.T) {
		schemaWithNestedArray := `{
			"version": 1,
			"types": {
				"Item": {
					"type": "object",
					"fields": {
						"name": {
							"type": "string"
						}
					}
				},
				"Collection": {
					"type": "object",
					"fields": {
						"items": {
							"type": "array",
							"arrayType": {
								"type": "array",
								"arrayType": {
									"type": "Item"
								}
							}
						}
					}
				}
			},
			"procedures": {
				"GetCollection": {
					"type": "query",
					"output": {
						"type": "Collection"
					}
				}
			}
		}`

		parsedSchema, err := schema.ParseSchema(schemaWithNestedArray)
		assert.NoError(t, err)
		assert.NotEqual(t, schema.Schema{}, parsedSchema)

		// Should validate nested array types correctly
		collectionType := parsedSchema.Types["Collection"]
		itemsField := collectionType.Fields["items"]
		assert.Equal(t, "array", itemsField.Type)
		assert.NotNil(t, itemsField.ArrayType)
		assert.Equal(t, "array", itemsField.ArrayType.Type)
		assert.NotNil(t, itemsField.ArrayType.ArrayType)
		assert.Equal(t, "Item", itemsField.ArrayType.ArrayType.Type)
	})

	t.Run("schema with undeclared type in array", func(t *testing.T) {
		schemaWithUndeclaredArrayType := `{
			"version": 1,
			"types": {
				"Collection": {
					"type": "object",
					"fields": {
						"items": {
							"type": "array",
							"arrayType": {
								"type": "UndeclaredItem"
							}
						}
					}
				}
			},
			"procedures": {
				"GetCollection": {
					"type": "query",
					"output": {
						"type": "Collection"
					}
				}
			}
		}`

		_, err := schema.ParseSchema(schemaWithUndeclaredArrayType)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "undefined custom type: UndeclaredItem")
	})
}
