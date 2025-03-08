package schema_test

import (
	"encoding/json"
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
