package pieces

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicValidation(t *testing.T) {
	templates := map[string][]string{
		"Person": {
			"name",
			"age",
		},
	}

	t.Run("Valid simple JSON", func(t *testing.T) {
		json := `{"name": "John", "age": 30, "extra": "field"}`
		err := ValidateJSONPaths([]byte(json), templates, "Person")
		assert.NoError(t, err)
	})

	t.Run("Missing required field", func(t *testing.T) {
		json := `{"name": "John"}`
		err := ValidateJSONPaths([]byte(json), templates, "Person")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "age: required field is missing")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		json := `{"name": "John", "age": 30`
		err := ValidateJSONPaths([]byte(json), templates, "Person")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid JSON")
	})

	t.Run("Non-existent template", func(t *testing.T) {
		json := `{"name": "John", "age": 30}`
		err := ValidateJSONPaths([]byte(json), templates, "NonExistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no template defined for type")
	})
}

func TestArrayValidation(t *testing.T) {
	templates := map[string][]string{
		"Collection": {
			"items",
			"items[*].id",
			"items[*].name",
		},
	}

	t.Run("Valid array with items", func(t *testing.T) {
		json := `{
			"items": [
				{"id": 1, "name": "Item 1", "extra": "data"},
				{"id": 2, "name": "Item 2"}
			]
		}`
		err := ValidateJSONPaths([]byte(json), templates, "Collection")
		assert.NoError(t, err)
	})

	t.Run("Empty array is valid", func(t *testing.T) {
		json := `{"items": []}`
		err := ValidateJSONPaths([]byte(json), templates, "Collection")
		assert.NoError(t, err)
	})

	t.Run("Missing field in array item", func(t *testing.T) {
		json := `{
			"items": [
				{"id": 1, "name": "Item 1"},
				{"id": 2}
			]
		}`
		err := ValidateJSONPaths([]byte(json), templates, "Collection")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name: required field is missing")
	})

	t.Run("Array expected but got something else", func(t *testing.T) {
		json := `{"items": "not an array"}`
		err := ValidateJSONPaths([]byte(json), templates, "Collection")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected array")
	})
}

func TestNestedObjectValidation(t *testing.T) {
	templates := map[string][]string{
		"User": {
			"profile",
			"profile.name",
			"profile.details.age",
			"profile.details.address.city",
		},
	}

	t.Run("Valid nested objects", func(t *testing.T) {
		json := `{
			"profile": {
				"name": "John Doe",
				"details": {
					"age": 30,
					"address": {
						"city": "New York",
						"country": "USA"
					}
				}
			}
		}`
		err := ValidateJSONPaths([]byte(json), templates, "User")
		assert.NoError(t, err)
	})

	t.Run("Missing deep nested field", func(t *testing.T) {
		json := `{
			"profile": {
				"name": "John Doe",
				"details": {
					"age": 30,
					"address": {
						"country": "USA"
					}
				}
			}
		}`
		err := ValidateJSONPaths([]byte(json), templates, "User")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "city: required field is missing")
	})

	t.Run("Intermediate object missing", func(t *testing.T) {
		json := `{
			"profile": {
				"name": "John Doe"
			}
		}`
		err := ValidateJSONPaths([]byte(json), templates, "User")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "details: required field is missing")
	})
}

func TestTypeReferenceValidation(t *testing.T) {
	templates := map[string][]string{
		"Order": {
			"id",
			"customer->Customer",
			"items",
			"items[*]->Product",
		},
		"Customer": {
			"id",
			"name",
			"email",
		},
		"Product": {
			"id",
			"name",
			"price",
		},
	}

	t.Run("Valid object with type references", func(t *testing.T) {
		json := `{
			"id": "ORD-001",
			"customer": {
				"id": "CUST-001",
				"name": "John Doe",
				"email": "john@example.com",
				"phone": "555-1234"
			},
			"items": [
				{
					"id": "PROD-001",
					"name": "Laptop",
					"price": 999.99,
					"description": "High performance laptop"
				},
				{
					"id": "PROD-002",
					"name": "Mouse",
					"price": 24.99
				}
			]
		}`
		err := ValidateJSONPaths([]byte(json), templates, "Order")
		assert.NoError(t, err)
	})

	t.Run("Missing field in referenced type", func(t *testing.T) {
		json := `{
			"id": "ORD-001",
			"customer": {
				"id": "CUST-001",
				"name": "John Doe"
			},
			"items": [
				{
					"id": "PROD-001",
					"name": "Laptop",
					"price": 999.99
				}
			]
		}`
		err := ValidateJSONPaths([]byte(json), templates, "Order")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email: required field is missing")
	})

	t.Run("Missing field in referenced array item", func(t *testing.T) {
		json := `{
			"id": "ORD-001",
			"customer": {
				"id": "CUST-001",
				"name": "John Doe",
				"email": "john@example.com"
			},
			"items": [
				{
					"id": "PROD-001",
					"name": "Laptop"
				}
			]
		}`
		err := ValidateJSONPaths([]byte(json), templates, "Order")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "price: required field is missing")
	})

	t.Run("Empty array with type references is valid", func(t *testing.T) {
		json := `{
			"id": "ORD-001",
			"customer": {
				"id": "CUST-001",
				"name": "John Doe",
				"email": "john@example.com"
			},
			"items": []
		}`
		err := ValidateJSONPaths([]byte(json), templates, "Order")
		assert.NoError(t, err)
	})
}

func TestRecursiveArrayValidation(t *testing.T) {
	templates := map[string][]string{
		"User": {
			"id",
			"name",
			"posts",
			"posts[*]->Post",
		},
		"Post": {
			"id",
			"title",
			"author->User",
		},
	}

	t.Run("Valid recursive structure", func(t *testing.T) {
		json := `{
			"id": 1,
			"name": "John Doe",
			"posts": [
				{
					"id": 101,
					"title": "First Post",
					"author": {
						"id": 1,
						"name": "John Doe",
						"posts": []
					}
				}
			]
		}`
		err := ValidateJSONPaths([]byte(json), templates, "User")
		assert.NoError(t, err)
	})

	t.Run("Recursive structure with deep nesting", func(t *testing.T) {
		json := `{
			"id": 1,
			"name": "John",
			"posts": [
				{
					"id": 101,
					"title": "Post 1",
					"author": {
						"id": 2,
						"name": "Jane",
						"posts": [
							{
								"id": 201,
								"title": "Jane's post",
								"author": {
									"id": 1,
									"name": "John",
									"posts": []
								}
							}
						]
					}
				}
			]
		}`
		err := ValidateJSONPaths([]byte(json), templates, "User")
		assert.NoError(t, err)
	})

	t.Run("Cyclic reference with deep nesting and missing field", func(t *testing.T) {
		json := `{
			"id": 1,
			"name": "John",
			"posts": [
				{
					"id": 101,
					"title": "Post 1",
					"author": {
						"id": 2,
						"name": "Jane",
						"posts": [
							{
								"id": 201,
								"author": {
									"id": 1,
									"name": "John",
									"posts": []
								}
							}
						]
					}
				}
			]
		}`
		err := ValidateJSONPaths([]byte(json), templates, "User")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title: required field is missing")
	})

	t.Run("Large nested array (1000 levels)", func(t *testing.T) {
		var genPost func(depth int) string
		genPost = func(depth int) string {
			if depth == 0 {
				return `{"id": 1, "title": "Final Post", "author": {"id": 1, "name": "John Doe", "posts": []}}`
			}
			return fmt.Sprintf(`{"id": %d, "title": "Post %d", "author": {"id": 1, "name": "John Doe", "posts": [%s]}}`, depth, depth, genPost(depth-1))
		}

		json := fmt.Sprintf(`{"id": 1, "name": "John Doe", "posts": [%s]}`, genPost(1000))
		err := ValidateJSONPaths([]byte(json), templates, "User")
		assert.NoError(t, err)
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("Empty object with no template should pass", func(t *testing.T) {
		templates := map[string][]string{
			"Empty": {},
		}
		json := `{}`
		err := ValidateJSONPaths([]byte(json), templates, "Empty")
		assert.NoError(t, err)
	})

	t.Run("Deeply nested path should work", func(t *testing.T) {
		templates := map[string][]string{
			"Deep": {
				"level1.level2.level3.level4.level5.value",
			},
		}
		json := `{
			"level1": {
				"level2": {
					"level3": {
						"level4": {
							"level5": {
								"value": "deep"
							}
						}
					}
				}
			}
		}`
		err := ValidateJSONPaths([]byte(json), templates, "Deep")
		assert.NoError(t, err)
	})

	t.Run("Complex array path with multi-level wildcards", func(t *testing.T) {
		templates := map[string][]string{
			"Complex": {
				"users[*].roles[*].permissions[*].name",
			},
		}
		json := `{
			"users": [
				{
					"roles": [
						{
							"permissions": [
								{"name": "read"},
								{"name": "write"}
							]
						}
					]
				},
				{
					"roles": [
						{
							"permissions": [
								{"name": "execute"}
							]
						}
					]
				}
			]
		}`
		err := ValidateJSONPaths([]byte(json), templates, "Complex")
		assert.NoError(t, err)
	})

	t.Run("Multiple errors should report first occurrence", func(t *testing.T) {
		templates := map[string][]string{
			"MultiError": {
				"field1",
				"field2",
			},
		}
		json := `{}`
		err := ValidateJSONPaths([]byte(json), templates, "MultiError")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required field is missing")
	})
}

func TestComplexTypeReferences(t *testing.T) {
	templates := map[string][]string{
		"API": {
			"endpoints[*]->Endpoint",
		},
		"Endpoint": {
			"path",
			"method",
			"responses[*]->Response",
			"parameters[*]->Parameter",
		},
		"Response": {
			"status",
			"body->Schema",
		},
		"Parameter": {
			"name",
			"type",
		},
		"Schema": {
			"type",
			"properties",
		},
	}

	t.Run("Complex nested type references", func(t *testing.T) {
		json := `{
			"endpoints": [
				{
					"path": "/users",
					"method": "GET",
					"responses": [
						{
							"status": 200,
							"body": {
								"type": "array",
								"properties": {
									"items": {
										"type": "object"
									}
								}
							}
						}
					],
					"parameters": [
						{
							"name": "limit",
							"type": "integer"
						}
					]
				}
			]
		}`
		err := ValidateJSONPaths([]byte(json), templates, "API")
		assert.NoError(t, err)
	})

	t.Run("Missing nested type reference field", func(t *testing.T) {
		json := `{
			"endpoints": [
				{
					"path": "/users",
					"method": "GET",
					"responses": [
						{
							"status": 200,
							"body": {
								"properties": {
									"items": {
										"type": "object"
									}
								}
							}
						}
					],
					"parameters": [
						{
							"name": "limit",
							"type": "integer"
						}
					]
				}
			]
		}`
		err := ValidateJSONPaths([]byte(json), templates, "API")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "type: required field is missing")
	})
}

func TestDataTypeVariations(t *testing.T) {
	templates := map[string][]string{
		"AllTypes": {
			"string_value",
			"number_value",
			"boolean_value",
			"null_value",
			"array_value",
			"object_value",
		},
	}

	t.Run("All JSON data types should be valid", func(t *testing.T) {
		json := `{
			"string_value": "text",
			"number_value": 42.5,
			"boolean_value": true,
			"null_value": null,
			"array_value": [1, 2, 3],
			"object_value": {"key": "value"}
		}`
		err := ValidateJSONPaths([]byte(json), templates, "AllTypes")
		assert.NoError(t, err)
	})

	t.Run("Special character values", func(t *testing.T) {
		json := `{
			"string_value": "Special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?",
			"number_value": -123.456e+78,
			"boolean_value": false,
			"null_value": null,
			"array_value": [],
			"object_value": {}
		}`
		err := ValidateJSONPaths([]byte(json), templates, "AllTypes")
		assert.NoError(t, err)
	})

	t.Run("Unicode strings", func(t *testing.T) {
		json := `{
			"string_value": "Unicode: 你好, Привет, こんにちは, مرحبا",
			"number_value": 0,
			"boolean_value": false,
			"null_value": null,
			"array_value": [],
			"object_value": {}
		}`
		err := ValidateJSONPaths([]byte(json), templates, "AllTypes")
		assert.NoError(t, err)
	})
}

func TestInvalidPathSyntax(t *testing.T) {
	t.Run("Invalid type reference format", func(t *testing.T) {
		templates := map[string][]string{
			"Invalid": {
				"field->->Type",
			},
		}
		json := `{"field": {}}`
		err := ValidateJSONPaths([]byte(json), templates, "Invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid type reference")
	})

	t.Run("Non-existent referenced type", func(t *testing.T) {
		templates := map[string][]string{
			"HasRef": {
				"field->NonExistent",
			},
		}
		json := `{"field": {}}`
		err := ValidateJSONPaths([]byte(json), templates, "HasRef")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no template defined for type")
	})
}

func TestLargeDataSets(t *testing.T) {
	// Skip this if we're not doing performance testing
	if testing.Short() {
		t.Skip("Skipping large dataset test in short mode")
	}

	// Create a template for large datasets
	templates := map[string][]string{
		"Large": {
			"items[*].id",
			"items[*].value",
		},
	}

	// Generate a large JSON dataset
	generateLargeJSON := func(size int) []byte {
		json := `{"items": [`
		for i := range size {
			if i > 0 {
				json += ","
			}
			json += `{"id": 1, "value": "item"}`
		}
		json += `]}`
		return []byte(json)
	}

	t.Run("Medium dataset (100 items)", func(t *testing.T) {
		json := generateLargeJSON(100)
		err := ValidateJSONPaths(json, templates, "Large")
		assert.NoError(t, err)
	})

	t.Run("Large dataset (1000 items)", func(t *testing.T) {
		json := generateLargeJSON(1000)
		err := ValidateJSONPaths(json, templates, "Large")
		assert.NoError(t, err)
	})
}

func TestBenchmarkValidation(t *testing.T) {
	// Skip benchmarks unless we're explicitly running them
	if testing.Short() {
		t.Skip("Skipping benchmarks in short mode")
	}

	templates := map[string][]string{
		"User": {
			"id",
			"name",
			"posts[*]->Post",
		},
		"Post": {
			"id",
			"title",
			"author->User",
		},
	}

	json := `{
		"id": 1,
		"name": "John Doe",
		"posts": [
			{
				"id": 101,
				"title": "First Post",
				"author": {
					"id": 1,
					"name": "John Doe",
					"posts": []
				}
			}
		]
	}`

	// Simple benchmark to ensure validation performance is reasonable
	b := require.New(t)
	for i := 0; i < 1000; i++ {
		err := ValidateJSONPaths([]byte(json), templates, "User")
		b.NoError(err)
	}
}
