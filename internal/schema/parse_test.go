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
					"description": "A user of the system",
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
					"input": {
						"id": {
							"type": "string",
							"rules": {
								"uuid": {}
							}
						}
					},
					"output": {
						"user": {
							"type": "User"
						}
					}
				},
				"CreateUser": {
					"input": {
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
					},
					"output": {
						"user": {
							"type": "User"
						}
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

		userType := parsedSchema.Types["User"]
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

		addressesField := parsedSchema.Procedures["CreateUser"].Input["addresses"]
		assert.Equal(t, "array", addressesField.Type)
		assert.NotNil(t, addressesField.ArrayType)
		assert.Equal(t, "Address", addressesField.ArrayType.Type)
	})

	t.Run("schema with undefined custom type", func(t *testing.T) {
		customSchema := `{
			"version": 1,
			"types": {
				"User": {
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
					"input": {
						"id": {
							"type": "string"
						}
					},
					"output": {
						"user": {
							"type": "User"
						}
					}
				}
			}
		}`

		_, err := schema.ParseSchema(customSchema)
		isErr := err != nil
		assert.True(t, isErr)
		assert.Contains(t, err.Error(), "undefined custom type: Profile")
	})

	t.Run("schema with circular reference", func(t *testing.T) {
		customSchema := `{
			"version": 1,
			"types": {
				"User": {
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
					"input": {
						"id": {
							"type": "string"
						}
					},
					"output": {
						"user": {
							"type": "User"
						}
					}
				}
			}
		}`
		parsedSchema, err := schema.ParseSchema(customSchema)
		assert.NoError(t, err)
		assert.NotEqual(t, schema.Schema{}, parsedSchema)

		// Circular references are allowed and should parse correctly
		userType := parsedSchema.Types["User"]
		friendField := userType.Fields["friend"]
		assert.Equal(t, "User", friendField.Type)
	})

	t.Run("schema with nested array of custom type", func(t *testing.T) {
		customSchema := `{
			"version": 1,
			"types": {
				"Tag": {
					"fields": {
						"name": {
							"type": "string"
						}
					}
				},
				"Post": {
					"fields": {
						"id": {
							"type": "string"
						},
						"tags": {
							"type": "array",
							"arrayType": {
								"type": "Tag"
							}
						}
					}
				}
			},
			"procedures": {
				"GetPost": {
					"input": {
						"id": {
							"type": "string"
						}
					},
					"output": {
						"post": {
							"type": "Post"
						}
					}
				}
			}
		}`

		parsedSchema, err := schema.ParseSchema(customSchema)
		assert.NoError(t, err)
		assert.NotEqual(t, schema.Schema{}, parsedSchema)

		postType := parsedSchema.Types["Post"]
		tagsField := postType.Fields["tags"]
		assert.Equal(t, "array", tagsField.Type)
		assert.NotNil(t, tagsField.ArrayType)
		assert.Equal(t, "Tag", tagsField.ArrayType.Type)
	})

	t.Run("schema with undeclared type in array", func(t *testing.T) {
		schemaWithUndeclaredArrayType := `{
			"version": 1,
			"types": {
				"Collection": {
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
					"input": {},
					"output": {
						"collection": {
							"type": "Collection"
						}
					}
				}
			}
		}`

		_, err := schema.ParseSchema(schemaWithUndeclaredArrayType)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "undefined custom type: UndeclaredItem")
	})
}
