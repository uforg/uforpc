package schema

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateStructure(t *testing.T) {
	t.Run("Valid empty schema", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": []
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.NoError(t, err)
	})

	t.Run("Valid schema with all node types", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "doc",
					"content": "Documentation"
				},
				{
					"kind": "rule",
					"name": "required",
					"for": {
						"type": "string",
						"isArray": false
					}
				},
				{
					"kind": "type",
					"name": "User",
					"fields": [
						{
							"name": "id",
							"typeName": "string",
							"isArray": false,
							"optional": false,
							"rules": []
						}
					]
				},
				{
					"kind": "proc",
					"name": "GetUser",
					"input": [
						{
							"name": "id",
							"typeName": "string",
							"isArray": false,
							"optional": false,
							"rules": []
						}
					],
					"output": [
						{
							"name": "user",
							"typeName": "User",
							"isArray": false,
							"optional": false,
							"rules": []
						}
					],
					"meta": []
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.NoError(t, err)
	})

	t.Run("Duplicate rule names", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "rule",
					"name": "required",
					"for": {
						"type": "string",
						"isArray": false
					}
				},
				{
					"kind": "rule",
					"name": "required",
					"for": {
						"type": "int",
						"isArray": false
					}
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.Error(t, err)
		require.Contains(t, err.Error(), "duplicate rule name: required")
	})

	t.Run("Duplicate type names", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "type",
					"name": "User",
					"fields": []
				},
				{
					"kind": "type",
					"name": "User",
					"fields": []
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.Error(t, err)
		require.Contains(t, err.Error(), "duplicate type name: User")
	})

	t.Run("Rule references undefined type", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "rule",
					"name": "validateUser",
					"for": {
						"type": "User",
						"isArray": false
					}
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.Error(t, err)
		require.Contains(t, err.Error(), "rule 'validateUser' references undefined type: User")
	})

	t.Run("Rule references invalid primitive type", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "rule",
					"name": "validateUser",
					"for": {
						"type": "invalid",
						"isArray": false
					}
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.Error(t, err)
		require.Contains(t, err.Error(), "rule 'validateUser' references invalid primitive type: invalid")
	})

	t.Run("Field references undefined type", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "type",
					"name": "User",
					"fields": [
						{
							"name": "address",
							"typeName": "Address",
							"isArray": false,
							"optional": false,
							"rules": []
						}
					]
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.Error(t, err)
		require.Contains(t, err.Error(), "field 'address' in type 'User' references undefined type: Address")
	})

	t.Run("Field references invalid primitive type", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "type",
					"name": "User",
					"fields": [
						{
							"name": "id",
							"typeName": "invalid",
							"isArray": false,
							"optional": false,
							"rules": []
						}
					]
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.Error(t, err)
		require.Contains(t, err.Error(), "field 'id' in type 'User' references invalid primitive type: invalid")
	})

	t.Run("Inline field references undefined type", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "type",
					"name": "User",
					"fields": [
						{
							"name": "address",
							"typeInline": {
								"fields": [
									{
										"name": "city",
										"typeName": "City",
										"isArray": false,
										"optional": false,
										"rules": []
									}
								]
							},
							"isArray": false,
							"optional": false,
							"rules": []
						}
					]
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.Error(t, err)
		require.Contains(t, err.Error(), "field 'city' in inline type in field 'address' of type 'User' references undefined type: City")
	})

	t.Run("Procedure input field references undefined type", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "proc",
					"name": "GetUser",
					"input": [
						{
							"name": "filter",
							"typeName": "UserFilter",
							"isArray": false,
							"optional": false,
							"rules": []
						}
					],
					"output": [],
					"meta": []
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.Error(t, err)
		require.Contains(t, err.Error(), "input field 'filter' in procedure 'GetUser' references undefined type: UserFilter")
	})

	t.Run("Procedure output field references undefined type", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "proc",
					"name": "GetUser",
					"input": [],
					"output": [
						{
							"name": "user",
							"typeName": "User",
							"isArray": false,
							"optional": false,
							"rules": []
						}
					],
					"meta": []
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.Error(t, err)
		require.Contains(t, err.Error(), "output field 'user' in procedure 'GetUser' references undefined type: User")
	})

	t.Run("Circular type reference without optional fields", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "type",
					"name": "User",
					"fields": [
						{
							"name": "address",
							"typeName": "Address",
							"isArray": false,
							"optional": false,
							"rules": []
						}
					]
				},
				{
					"kind": "type",
					"name": "Address",
					"fields": [
						{
							"name": "user",
							"typeName": "User",
							"isArray": false,
							"optional": false,
							"rules": []
						}
					]
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "circular type reference detected without optional fields"))
	})

	t.Run("Circular type reference with optional field", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "type",
					"name": "User",
					"fields": [
						{
							"name": "address",
							"typeName": "Address",
							"isArray": false,
							"optional": false,
							"rules": []
						}
					]
				},
				{
					"kind": "type",
					"name": "Address",
					"fields": [
						{
							"name": "user",
							"typeName": "User",
							"isArray": false,
							"optional": true,
							"rules": []
						}
					]
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.NoError(t, err)
	})

	t.Run("Complex circular reference with optional field", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "type",
					"name": "User",
					"fields": [
						{
							"name": "department",
							"typeName": "Department",
							"isArray": false,
							"optional": false,
							"rules": []
						}
					]
				},
				{
					"kind": "type",
					"name": "Department",
					"fields": [
						{
							"name": "company",
							"typeName": "Company",
							"isArray": false,
							"optional": false,
							"rules": []
						}
					]
				},
				{
					"kind": "type",
					"name": "Company",
					"fields": [
						{
							"name": "employees",
							"typeName": "User",
							"depth": 1,
							"optional": true,
							"rules": []
						}
					]
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.NoError(t, err)
	})

	t.Run("Circular reference in inline types without optional field", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "type",
					"name": "User",
					"fields": [
						{
							"name": "details",
							"typeInline": {
								"fields": [
									{
										"name": "preferences",
										"typeInline": {
											"fields": [
												{
													"name": "user",
													"typeName": "User",
													"isArray": false,
													"optional": false,
													"rules": []
												}
											]
										},
										"isArray": false,
										"optional": false,
										"rules": []
									}
								]
							},
							"isArray": false,
							"optional": false,
							"rules": []
						}
					]
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "circular type reference detected without optional fields"))
	})

	t.Run("Circular reference in inline types with optional field", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "type",
					"name": "User",
					"fields": [
						{
							"name": "details",
							"typeInline": {
								"fields": [
									{
										"name": "preferences",
										"typeInline": {
											"fields": [
												{
													"name": "user",
													"typeName": "User",
													"isArray": false,
													"optional": true,
													"rules": []
												}
											]
										},
										"isArray": false,
										"optional": false,
										"rules": []
									}
								]
							},
							"isArray": false,
							"optional": false,
							"rules": []
						}
					]
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.NoError(t, err)
	})

	t.Run("Optional parent field in inline type makes circular reference valid", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "type",
					"name": "User",
					"fields": [
						{
							"name": "details",
							"typeInline": {
								"fields": [
									{
										"name": "preferences",
										"typeInline": {
											"fields": [
												{
													"name": "user",
													"typeName": "User",
													"isArray": false,
													"optional": false,
													"rules": []
												}
											]
										},
										"isArray": false,
										"optional": true,
										"rules": []
									}
								]
							},
							"isArray": false,
							"optional": false,
							"rules": []
						}
					]
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.NoError(t, err)
	})

	t.Run("Field uses undefined rule", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "type",
					"name": "User",
					"fields": [
						{
							"name": "id",
							"typeName": "string",
							"isArray": false,
							"optional": false,
							"rules": [
								{
									"rule": "uuid",
									"param": null,
									"error": null
								}
							]
						}
					]
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.Error(t, err)
		require.Contains(t, err.Error(), "field 'id' in type 'User' uses undefined rule: uuid")
	})

	t.Run("Inline field uses undefined rule", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "type",
					"name": "User",
					"fields": [
						{
							"name": "address",
							"typeInline": {
								"fields": [
									{
										"name": "city",
										"typeName": "string",
										"isArray": false,
										"optional": false,
										"rules": [
											{
												"rule": "minlen",
												"param": null,
												"error": null
											}
										]
									}
								]
							},
							"isArray": false,
							"optional": false,
							"rules": []
						}
					]
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.Error(t, err)
		require.Contains(t, err.Error(), "field 'city' in inline type in field 'address' of type 'User' uses undefined rule: minlen")
	})

	t.Run("Procedure input field uses undefined rule", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "proc",
					"name": "GetUser",
					"input": [
						{
							"name": "id",
							"typeName": "string",
							"isArray": false,
							"optional": false,
							"rules": [
								{
									"rule": "required",
									"param": null,
									"error": null
								}
							]
						}
					],
					"output": [],
					"meta": []
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.Error(t, err)
		require.Contains(t, err.Error(), "input field 'id' in procedure 'GetUser' uses undefined rule: required")
	})

	t.Run("Procedure output field uses undefined rule", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "proc",
					"name": "GetUser",
					"input": [],
					"output": [
						{
							"name": "success",
							"typeName": "bool",
							"isArray": false,
							"optional": false,
							"rules": [
								{
									"rule": "equals",
									"param": null,
									"error": null
								}
							]
						}
					],
					"meta": []
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.Error(t, err)
		require.Contains(t, err.Error(), "output field 'success' in procedure 'GetUser' uses undefined rule: equals")
	})

	t.Run("Valid schema with defined rules", func(t *testing.T) {
		input := `{
			"version": 1,
			"nodes": [
				{
					"kind": "rule",
					"name": "uuid",
						"for": {
							"type": "string",
							"isArray": false
						}
				},
				{
					"kind": "type",
					"name": "User",
					"fields": [
						{
							"name": "id",
							"typeName": "string",
							"isArray": false,
							"optional": false,
							"rules": [
								{
									"rule": "uuid",
									"param": null,
									"error": null
								}
							]
						}
					]
				}
			]
		}`

		var schema Schema
		err := json.Unmarshal([]byte(input), &schema)
		require.NoError(t, err)

		err = validateStructure(schema)
		require.NoError(t, err)
	})
}
