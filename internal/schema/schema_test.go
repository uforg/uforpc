package schema_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/schema"
)

func TestMarshalEmptySchema(t *testing.T) {
	s := schema.Schema{}

	data, err := json.Marshal(s)
	require.NoError(t, err)

	expected := `{}`
	assert.JSONEq(t, expected, string(data))
}

func TestUnmarshalEmptySchema(t *testing.T) {
	jsonData := `{"types":{},"procedures":{}}`

	var s schema.Schema
	err := json.Unmarshal([]byte(jsonData), &s)
	require.NoError(t, err)

	assert.NotNil(t, s.Types)
	assert.Empty(t, s.Types)
	assert.NotNil(t, s.Procedures)
	assert.Empty(t, s.Procedures)
}

func TestMarshalSimpleSchema(t *testing.T) {
	s := schema.Schema{
		Types: map[string]schema.Field{
			"User": {
				Type: "object",
				Fields: map[string]schema.Field{
					"id": {Type: "string"},
				},
			},
		},
		Procedures: map[string]schema.Procedure{
			"GetUser": {
				Type: schema.ProcedureTypeQuery,
				Output: schema.Field{
					Type: "User",
				},
			},
		},
	}

	data, err := json.Marshal(s)
	require.NoError(t, err)

	var unmarshalled schema.Schema
	err = json.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	assert.Equal(t, 1, len(unmarshalled.Types))
	assert.Equal(t, 1, len(unmarshalled.Procedures))
	assert.Equal(t, "object", unmarshalled.Types["User"].Type)
	assert.Contains(t, unmarshalled.Procedures, "GetUser")
	assert.Equal(t, schema.ProcedureTypeQuery, unmarshalled.Procedures["GetUser"].Type)
}

func TestUnmarshalComplexSchema(t *testing.T) {
	jsonData := `{
		"types": {
			"User": {
				"type": "object",
				"description": "User entity",
				"fields": {
					"id": {
						"type": "string",
						"rules": [
							{
								"name": "uuid",
								"message": "Must be a valid UUID"
							}
						]
					},
					"email": {
						"type": "string",
						"rules": [
							{
								"name": "email",
								"message": "Must be a valid email"
							},
							{
								"name": "minLength",
								"value": "5",
								"message": "Email must be at least 5 characters"
							}
						]
					},
					"roles": {
						"type": "string[]",
						"description": "User roles"
					}
				}
			},
			"Post": {
				"type": "object",
				"fields": {
					"title": {
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
							"type": "string"
						}
					}
				},
				"output": {
					"type": "User"
				},
				"meta": {
					"auth": true
				}
			},
			"UpdateUser": {
				"type": "mutation",
				"input": {
					"type": "object",
					"fields": {
						"user": {
							"type": "User"
						}
					}
				},
				"output": {
					"type": "User"
				}
			}
		}
	}`

	var s schema.Schema
	err := json.Unmarshal([]byte(jsonData), &s)
	require.NoError(t, err)

	// Verify types
	assert.Equal(t, 2, len(s.Types))
	assert.Contains(t, s.Types, "User")
	assert.Contains(t, s.Types, "Post")

	// Verify User type
	user := s.Types["User"]
	assert.Equal(t, "object", user.Type)
	assert.Equal(t, "User entity", user.Description)
	assert.Equal(t, 3, len(user.Fields))

	// Verify email field rules
	emailField := user.Fields["email"]
	assert.Equal(t, 2, len(emailField.Rules))
	assert.Equal(t, schema.RuleNameEmail, emailField.Rules[0].Name)
	assert.Equal(t, schema.RuleNameMinLength, emailField.Rules[1].Name)
	assert.Equal(t, "5", emailField.Rules[1].Value)

	// Verify array field
	rolesField := user.Fields["roles"]
	assert.Equal(t, "string[]", rolesField.Type)
	assert.Equal(t, "User roles", rolesField.Description)

	// Verify procedures
	assert.Equal(t, 2, len(s.Procedures))
	assert.Contains(t, s.Procedures, "GetUser")
	assert.Contains(t, s.Procedures, "UpdateUser")

	// Verify procedure meta
	getUserProc := s.Procedures["GetUser"]
	assert.Equal(t, schema.ProcedureTypeQuery, getUserProc.Type)
	assert.Contains(t, getUserProc.Meta, "auth")
	assert.Equal(t, true, getUserProc.Meta["auth"])

	// Verify update procedure
	updateUserProc := s.Procedures["UpdateUser"]
	assert.Equal(t, schema.ProcedureTypeMutation, updateUserProc.Type)
}

func TestRoundTripSerialization(t *testing.T) {
	// Create original schema
	original := schema.Schema{
		Types: map[string]schema.Field{
			"SimpleType": {
				Type: "object",
				Fields: map[string]schema.Field{
					"name": {
						Type: "string",
						Rules: []schema.RuleCatchAll{
							{
								Name:    schema.RuleNameMinLength,
								Value:   "3",
								Message: "Name must be at least 3 characters",
							},
						},
					},
					"active": {
						Type: "boolean",
					},
				},
			},
		},
		Procedures: map[string]schema.Procedure{
			"SimpleQuery": {
				Type: schema.ProcedureTypeQuery,
				Input: schema.Field{
					Type: "object",
					Fields: map[string]schema.Field{
						"id": {Type: "string"},
					},
				},
				Output: schema.Field{
					Type: "SimpleType",
				},
				Meta: map[string]any{
					"cached": true,
					"ttl":    60,
				},
			},
		},
	}

	// Serialize
	data, err := json.Marshal(original)
	require.NoError(t, err, "Marshal should not fail")

	// Deserialize
	var result schema.Schema
	err = json.Unmarshal(data, &result)
	require.NoError(t, err, "Unmarshal should not fail")

	// Verify types
	assert.Equal(t, len(original.Types), len(result.Types))
	assert.Contains(t, result.Types, "SimpleType")

	originalType := original.Types["SimpleType"]
	resultType := result.Types["SimpleType"]
	assert.Equal(t, originalType.Type, resultType.Type)

	// Verify fields
	assert.Contains(t, resultType.Fields, "name")
	assert.Contains(t, resultType.Fields, "active")

	// Verify rules
	nameField := resultType.Fields["name"]
	assert.Equal(t, 1, len(nameField.Rules))
	assert.Equal(t, schema.RuleNameMinLength, nameField.Rules[0].Name)
	assert.Equal(t, "3", nameField.Rules[0].Value)

	// Verify procedures
	assert.Equal(t, len(original.Procedures), len(result.Procedures))
	assert.Contains(t, result.Procedures, "SimpleQuery")
	proc := result.Procedures["SimpleQuery"]
	assert.Equal(t, schema.ProcedureTypeQuery, proc.Type)

	// Verify metadata
	assert.Equal(t, original.Procedures["SimpleQuery"].Meta["cached"], result.Procedures["SimpleQuery"].Meta["cached"])
	assert.EqualValues(t, original.Procedures["SimpleQuery"].Meta["ttl"], result.Procedures["SimpleQuery"].Meta["ttl"])
}

func TestUnmarshalInvalidJSON(t *testing.T) {
	invalidJson := `{"types": {"User": {"type": "object", "fields": {"id": {"type": "string",}}}}}` // Extra comma

	var s schema.Schema
	err := json.Unmarshal([]byte(invalidJson), &s)
	assert.Error(t, err)
}

func TestProcedureUpperCamelCaseKeys(t *testing.T) {
	jsonData := `{
		"procedures": {
			"GetUser": {"type": "query"},
			"UpdateProfile": {"type": "mutation"},
			"DeleteAccount": {"type": "mutation"}
		}
	}`

	var s schema.Schema
	err := json.Unmarshal([]byte(jsonData), &s)
	require.NoError(t, err)

	assert.Equal(t, 3, len(s.Procedures))
	assert.Contains(t, s.Procedures, "GetUser")
	assert.Contains(t, s.Procedures, "UpdateProfile")
	assert.Contains(t, s.Procedures, "DeleteAccount")
}

func TestNonStandardTypesHandling(t *testing.T) {
	jsonData := `{
		"types": {
			"CustomType": {
				"type": "OtherType",
				"fields": {
					"value": {"type": "any"}
				}
			}
		}
	}`

	var s schema.Schema
	err := json.Unmarshal([]byte(jsonData), &s)
	require.NoError(t, err)

	assert.Contains(t, s.Types, "CustomType")
	assert.Equal(t, "OtherType", s.Types["CustomType"].Type)
	assert.Equal(t, "any", s.Types["CustomType"].Fields["value"].Type)
}

func TestNestedPrimitiveArrays(t *testing.T) {
	jsonData := `{
		"types": {
			"ComplexArray": {
				"type": "object",
				"fields": {
					"matrix": {"type": "int[][]"},
					"cube": {"type": "string[][][]"}
				}
			}
		}
	}`

	var s schema.Schema
	err := json.Unmarshal([]byte(jsonData), &s)
	require.NoError(t, err)

	matrixField := s.Types["ComplexArray"].Fields["matrix"]
	assert.Equal(t, "int[][]", matrixField.Type)
	assert.True(t, matrixField.IsArray())
	assert.Equal(t, 2, matrixField.GetArrayDepth())
	assert.Equal(t, "int", matrixField.GetBaseType())

	cubeField := s.Types["ComplexArray"].Fields["cube"]
	assert.Equal(t, "string[][][]", cubeField.Type)
	assert.True(t, cubeField.IsArray())
	assert.Equal(t, 3, cubeField.GetArrayDepth())
	assert.Equal(t, "string", cubeField.GetBaseType())
}

func TestSelfReferencingTypes(t *testing.T) {
	jsonData := `{
		"types": {
			"Node": {
				"type": "object",
				"fields": {
					"value": {"type": "string"},
					"parent": {"type": "Node"},
					"children": {"type": "Node[]"}
				}
			}
		}
	}`

	var s schema.Schema
	err := json.Unmarshal([]byte(jsonData), &s)
	require.NoError(t, err)

	nodeType := s.Types["Node"]
	assert.Equal(t, "Node", nodeType.Fields["parent"].Type)
	assert.Equal(t, "Node[]", nodeType.Fields["children"].Type)
	assert.Equal(t, "Node", nodeType.Fields["children"].GetBaseType())
}

func TestDescriptionHandling(t *testing.T) {
	// Test with multi-line descriptions and special characters
	jsonData := `{
		"types": {
			"Documented": {
				"type": "object",
				"description": "This is a type with\nmulti-line description\nand \"special\" characters: & < >",
				"fields": {
					"field1": {
						"type": "string",
						"description": "Field with emoji 😀"
					}
				}
			}
		}
	}`

	var s schema.Schema
	err := json.Unmarshal([]byte(jsonData), &s)
	require.NoError(t, err)

	assert.Contains(t, s.Types, "Documented")
	documented := s.Types["Documented"]
	assert.Contains(t, documented.Description, "multi-line description")
	assert.Contains(t, documented.Description, "\"special\" characters")
	assert.Contains(t, documented.Fields["field1"].Description, "emoji 😀")
}

func TestProcedureWithNoInputOutput(t *testing.T) {
	jsonData := `{
		"procedures": {
			"EmptyProc": {
				"type": "mutation",
				"description": "A procedure with no input or output"
			}
		}
	}`

	var s schema.Schema
	err := json.Unmarshal([]byte(jsonData), &s)
	require.NoError(t, err)

	assert.Contains(t, s.Procedures, "EmptyProc")
	emptyProc := s.Procedures["EmptyProc"]
	assert.Equal(t, schema.ProcedureTypeMutation, emptyProc.Type)
	assert.Equal(t, "A procedure with no input or output", emptyProc.Description)
	assert.Empty(t, emptyProc.Input.Type)
	assert.Empty(t, emptyProc.Output.Type)
}

func TestUnmarshalWithSpacesInTypeNames(t *testing.T) {
	jsonData := `{
		"types": {
			" SpacedType ": {
				"type": " object ",
				"fields": {
					" fieldWithSpace ": {"type": " string "}
				}
			}
		}
	}`

	var s schema.Schema
	err := json.Unmarshal([]byte(jsonData), &s)
	require.NoError(t, err)

	// The spaces in the type name are preserved in the map key
	assert.Contains(t, s.Types, " SpacedType ")

	// The type value should have spaces preserved
	assert.Equal(t, " object ", s.Types[" SpacedType "].Type)

	// The field key should preserve spaces
	assert.Contains(t, s.Types[" SpacedType "].Fields, " fieldWithSpace ")

	// The field type should preserve spaces
	field := s.Types[" SpacedType "].Fields[" fieldWithSpace "]
	assert.Equal(t, " string ", field.Type)

	// Verify our utility methods handle spaces correctly
	assert.False(t, field.IsArray())
	assert.Equal(t, 0, field.GetArrayDepth())
	assert.Equal(t, "string", field.GetBaseType())
}

func TestEmptyAndNilMapHandling(t *testing.T) {
	// Test nil vs empty map behaviors
	s1 := schema.Schema{}
	s2 := schema.Schema{
		Types:      map[string]schema.Field{},
		Procedures: map[string]schema.Procedure{},
	}

	// Marshal both
	data1, err := json.Marshal(s1)
	require.NoError(t, err)

	data2, err := json.Marshal(s2)
	require.NoError(t, err)

	// Nil maps should be omitted, empty maps should be included
	assert.JSONEq(t, `{}`, string(data1))
	assert.JSONEq(t, `{"types":{},"procedures":{}}`, string(data2))
}

func TestComplexRulesArrayUnmarshal(t *testing.T) {
	jsonData := `{
		"types": {
			"ValidatedType": {
				"type": "object",
				"fields": {
					"password": {
						"type": "string",
						"rules": [
							{
								"name": "minLength",
								"value": "8",
								"message": "Password must be at least 8 characters"
							},
							{
								"name": "regex",
								"value": "^(?=.*[A-Za-z])(?=.*\\d)[A-Za-z\\d]{8,}$",
								"message": "Password must contain letters and numbers"
							},
							{
								"name": "custom-rule",
								"value": "complex-value",
								"message": "Custom validation rule"
							}
						]
					}
				}
			}
		}
	}`

	var s schema.Schema
	err := json.Unmarshal([]byte(jsonData), &s)
	require.NoError(t, err)

	passwordField := s.Types["ValidatedType"].Fields["password"]
	assert.Equal(t, 3, len(passwordField.Rules))

	// Check standard rules
	assert.Equal(t, schema.RuleNameMinLength, passwordField.Rules[0].Name)
	assert.Equal(t, "8", passwordField.Rules[0].Value)

	// Check regex rule with special characters
	assert.Equal(t, schema.RuleNameRegex, passwordField.Rules[1].Name)
	assert.Equal(t, "^(?=.*[A-Za-z])(?=.*\\d)[A-Za-z\\d]{8,}$", passwordField.Rules[1].Value)

	// Check custom rule
	assert.Equal(t, schema.RuleName{"custom-rule"}, passwordField.Rules[2].Name)
	assert.Equal(t, "complex-value", passwordField.Rules[2].Value)
}

func TestUnmarshalMetaWithNumericValues(t *testing.T) {
	jsonData := `{
		"procedures": {
			"NumericMeta": {
				"type": "query",
				"meta": {
					"intValue": 42,
					"floatValue": 3.14,
					"zeroValue": 0,
					"negativeValue": -1,
					"largeValue": 9007199254740991
				}
			}
		}
	}`

	var s schema.Schema
	err := json.Unmarshal([]byte(jsonData), &s)
	require.NoError(t, err)

	meta := s.Procedures["NumericMeta"].Meta

	// JSON numbers are unmarshaled as float64 in Go
	assert.Equal(t, float64(42), meta["intValue"])
	assert.Equal(t, 3.14, meta["floatValue"])
	assert.Equal(t, float64(0), meta["zeroValue"])
	assert.Equal(t, float64(-1), meta["negativeValue"])
	assert.Equal(t, float64(9007199254740991), meta["largeValue"])
}

func TestRuleNameMarshalUnmarshal(t *testing.T) {
	// Standard rule names
	standardRules := []schema.RuleName{
		schema.RuleNameOptional,
		schema.RuleNameEquals,
		schema.RuleNameMinLength,
		schema.RuleNameEmail,
		schema.RuleNameUuid,
	}

	for _, rule := range standardRules {
		t.Run(rule.Value, func(t *testing.T) {
			data, err := json.Marshal(rule)
			require.NoError(t, err)

			expectedJson := fmt.Sprintf(`"%s"`, rule.Value)
			assert.Equal(t, expectedJson, string(data))

			var unmarshalled schema.RuleName
			err = json.Unmarshal(data, &unmarshalled)
			require.NoError(t, err)

			assert.Equal(t, rule, unmarshalled)
		})
	}

	// Custom rule name
	customRule := schema.RuleName{"custom-validator"}
	data, err := json.Marshal(customRule)
	require.NoError(t, err)

	assert.Equal(t, `"custom-validator"`, string(data))

	var unmarshalled schema.RuleName
	err = json.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	assert.Equal(t, customRule, unmarshalled)
}

func TestProcedureTypeMarshalUnmarshal(t *testing.T) {
	// Test both procedure types
	procedureTypes := []schema.ProcedureType{
		schema.ProcedureTypeQuery,
		schema.ProcedureTypeMutation,
	}

	for _, procType := range procedureTypes {
		t.Run(procType.Value, func(t *testing.T) {
			data, err := json.Marshal(procType)
			require.NoError(t, err)

			expectedJson := fmt.Sprintf(`"%s"`, procType.Value)
			assert.Equal(t, expectedJson, string(data))

			var unmarshalled schema.ProcedureType
			err = json.Unmarshal(data, &unmarshalled)
			require.NoError(t, err)

			assert.Equal(t, procType, unmarshalled)
		})
	}

	// Custom procedure type
	customType := schema.ProcedureType{"subscription"}
	data, err := json.Marshal(customType)
	require.NoError(t, err)

	assert.Equal(t, `"subscription"`, string(data))

	var unmarshalled schema.ProcedureType
	err = json.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	assert.Equal(t, customType, unmarshalled)
}

func TestPrimitiveTypeMarshalUnmarshal(t *testing.T) {
	// Test all primitive types
	primitiveTypes := []schema.PrimitiveType{
		schema.PrimitiveTypeString,
		schema.PrimitiveTypeInt,
		schema.PrimitiveTypeFloat,
		schema.PrimitiveTypeBoolean,
	}

	for _, primType := range primitiveTypes {
		t.Run(primType.Value, func(t *testing.T) {
			data, err := json.Marshal(primType)
			require.NoError(t, err)

			expectedJson := fmt.Sprintf(`"%s"`, primType.Value)
			assert.Equal(t, expectedJson, string(data))

			var unmarshalled schema.PrimitiveType
			err = json.Unmarshal(data, &unmarshalled)
			require.NoError(t, err)

			assert.Equal(t, primType, unmarshalled)
		})
	}
}

func TestRuleCatchAllToSpecificRule(t *testing.T) {
	// Test conversion of RuleCatchAll to specific rule types

	// Optional rule (no value needed)
	optionalRule := schema.RuleCatchAll{
		Name:    schema.RuleNameOptional,
		Message: "This field is optional",
	}

	specificOptional := optionalRule.ToSpecificRule()
	require.IsType(t, schema.RuleOptional{}, specificOptional)
	typedRule := specificOptional.(schema.RuleOptional)
	assert.Equal(t, schema.RuleNameOptional, typedRule.Name)
	assert.Equal(t, "This field is optional", typedRule.Message)

	// Equals rule (needs a value)
	equalsRule := schema.RuleCatchAll{
		Name:    schema.RuleNameEquals,
		Value:   "admin",
		Message: "Must equal admin",
	}

	specificEquals := equalsRule.ToSpecificRule()
	require.IsType(t, schema.RuleEquals{}, specificEquals)
	typedEqualsRule := specificEquals.(schema.RuleEquals)
	assert.Equal(t, schema.RuleNameEquals, typedEqualsRule.Name)
	assert.Equal(t, "admin", typedEqualsRule.Value)
	assert.Equal(t, "Must equal admin", typedEqualsRule.Message)

	// Unimplemented rule (should return nil)
	unimplementedRule := schema.RuleCatchAll{
		Name:    schema.RuleName{"notImplemented"},
		Message: "Must be a valid email",
	}

	specific := unimplementedRule.ToSpecificRule()
	assert.Nil(t, specific)
}
