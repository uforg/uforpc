package schema_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/schema"
)

func TestMarshalUnmarshalRule(t *testing.T) {
	jsonData := `{
		"name": "equals",
		"value": "expected-value",
		"message": "Must equal expected-value"
	}`

	var rule schema.RuleEquals
	err := json.Unmarshal([]byte(jsonData), &rule)
	require.NoError(t, err)

	assert.Equal(t, schema.RuleNameEquals, rule.Name)
	assert.Equal(t, "expected-value", rule.Value)
	assert.Equal(t, "Must equal expected-value", rule.Message)

	marshaledData, err := json.Marshal(rule)
	require.NoError(t, err)

	assert.JSONEq(t, jsonData, string(marshaledData))
}

func TestParseImplementedRuleFromJSON(t *testing.T) {
	// Test parsing a rule that is implemented in ToSpecificRule
	jsonData := `{
		"name": "equals",
		"value": "expected-value",
		"message": "Must equal expected-value"
	}`

	var rule schema.RuleCatchAll
	err := json.Unmarshal([]byte(jsonData), &rule)
	require.NoError(t, err)

	// Verify parsing was correct
	assert.Equal(t, schema.RuleNameEquals, rule.Name)
	assert.Equal(t, "expected-value", rule.Value)
	assert.Equal(t, "Must equal expected-value", rule.Message)

	// Convert to specific rule
	specificRule := rule.ToSpecificRule()
	equalsRule, ok := specificRule.(schema.RuleEquals)
	require.True(t, ok)
	assert.Equal(t, schema.RuleNameEquals, equalsRule.Name)
	assert.Equal(t, "expected-value", equalsRule.Value)
	assert.Equal(t, "Must equal expected-value", equalsRule.Message)
}

func TestParseUnimplementedRuleFromJSON(t *testing.T) {
	// Test parsing a rule that exists but is not implemented in ToSpecificRule
	jsonData := `{
		"name": "notImplemented"
	}`

	var rule schema.RuleCatchAll
	err := json.Unmarshal([]byte(jsonData), &rule)
	require.NoError(t, err)

	// Verify parsing was correct
	assert.Equal(t, schema.RuleName{"notImplemented"}, rule.Name)

	// Try to convert to specific rule
	specificRule := rule.ToSpecificRule()
	assert.Nil(t, specificRule, "Unimplemented rule should return nil from ToSpecificRule")
}

func TestParseMultipleRulesFromJSON(t *testing.T) {
	// Test parsing multiple rules in a list
	jsonData := `[
		{
			"name": "optional",
			"message": "Field is optional"
		},
		{
			"name": "equals",
			"value": "test-value",
			"message": "Must equal test-value"
		},
		{
			"name": "email",
			"message": "Must be a valid email"
		}
	]`

	var rules []schema.RuleCatchAll
	err := json.Unmarshal([]byte(jsonData), &rules)
	require.NoError(t, err)

	// Verify correct number of rules
	require.Len(t, rules, 3)

	// Verify rule details
	assert.Equal(t, schema.RuleNameOptional, rules[0].Name)
	assert.Equal(t, "", rules[0].Value)
	assert.Equal(t, "Field is optional", rules[0].Message)

	assert.Equal(t, schema.RuleNameEquals, rules[1].Name)
	assert.Equal(t, "test-value", rules[1].Value)
	assert.Equal(t, "Must equal test-value", rules[1].Message)

	assert.Equal(t, schema.RuleNameEmail, rules[2].Name)
	assert.Equal(t, "", rules[2].Value)
	assert.Equal(t, "Must be a valid email", rules[2].Message)
}
