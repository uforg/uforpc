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
				assert.NotNil(t, parsedSchema)
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
				assert.Nil(t, parsedSchema)
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
		require.NotNil(t, parsedSchema)

		assert.Equal(t, 1, parsedSchema.Version)

		assert.Len(t, parsedSchema.Types, 2)
		assert.Contains(t, parsedSchema.Types, "User")
		assert.Contains(t, parsedSchema.Types, "Address")

		assert.Len(t, parsedSchema.Procedures, 2)
		assert.Contains(t, parsedSchema.Procedures, "GetUser")
		assert.Contains(t, parsedSchema.Procedures, "CreateUser")

		getUserProc := parsedSchema.Procedures["GetUser"]
		assert.Equal(t, "query", getUserProc.Type)

		createUserProc := parsedSchema.Procedures["CreateUser"]
		assert.Equal(t, "mutation", createUserProc.Type)

		userType := parsedSchema.Types["User"]
		assert.Equal(t, "object", userType.Type)
		assert.Len(t, userType.Fields, 4)
		assert.Contains(t, userType.Fields, "id")
		assert.Contains(t, userType.Fields, "name")
		assert.Contains(t, userType.Fields, "age")
		assert.Contains(t, userType.Fields, "isActive")

		nameField := userType.Fields["name"]
		nameRules, ok := nameField.Rules.(*schema.StringRules)
		assert.True(t, ok)
		assert.Equal(t, 3, nameRules.MinLen.Value)
		assert.Equal(t, "Name must be at least 3 characters", nameRules.MinLen.ErrorMessage)

		addressesField := createUserProc.Input.Fields["addresses"]
		assert.Equal(t, "array", addressesField.Type)
		assert.NotNil(t, addressesField.ArrayType)
		assert.Equal(t, "Address", addressesField.ArrayType.Type)
	})
}
