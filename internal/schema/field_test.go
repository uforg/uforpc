package schema_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uforg/uforpc/internal/schema"
)

func TestIsArray(t *testing.T) {
	tests := []struct {
		name      string
		fieldType string
		expected  bool
	}{
		// Positive cases
		{"Simple array type", "string[]", true},
		{"Nested array type", "string[][]", true},
		{"Object array type", "User[]", true},
		{"Array type with spaces", "int[] ", true},

		// Negative cases
		{"Primitive string type", "string", false},
		{"Primitive int type", "int", false},
		{"Object type", "User", false},
		{"Type with brackets in middle", "arr[5]ays", false},
		{"Type with opening bracket", "[string", false},
		{"Type with single bracket", "string[", false},
		{"Empty type", "", false},
		{"Malformed array type", "[]string", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := schema.Field{Type: tt.fieldType}
			result := field.IsArray()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetArrayDepth(t *testing.T) {
	tests := []struct {
		name      string
		fieldType string
		expected  int
	}{
		// Arrays of different depths
		{"Not an array", "string", 0},
		{"Simple array", "string[]", 1},
		{"Two-dimensional array", "string[][]", 2},
		{"Three-dimensional array", "string[][][]", 3},
		{"Custom type array", "User[]", 1},
		{"Two-dimensional custom type array", "User[][]", 2},

		// Edge cases
		{"Empty type", "", 0},
		{"Incomplete bracket", "string[", 0},
		{"Wrong format", "[]string", 0},
		{"Type with trailing space", "int[] ", 1},
		{"Non-contiguous brackets", "string[]object[]", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := schema.Field{Type: tt.fieldType}
			result := field.GetArrayDepth()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetBaseType(t *testing.T) {
	tests := []struct {
		name      string
		fieldType string
		expected  string
	}{
		// Basic types
		{"String type", "string", "string"},
		{"Integer type", "int", "int"},
		{"Float type", "float", "float"},
		{"Boolean type", "boolean", "boolean"},
		{"Custom type", "User", "User"},

		// Arrays of different depths
		{"Simple array", "string[]", "string"},
		{"Two-dimensional array", "string[][]", "string"},
		{"Three-dimensional array", "string[][][]", "string"},
		{"Custom type array", "User[]", "User"},
		{"Two-dimensional custom type array", "User[][]", "User"},

		// Edge cases
		{"Empty type", "", ""},
		{"Type with spaces", "  int[]  ", "int"},
		{"Incomplete bracket", "string[", "string["},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := schema.Field{Type: tt.fieldType}
			result := field.GetBaseType()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFieldMethodChaining(t *testing.T) {
	tests := []struct {
		name          string
		fieldType     string
		expectedBase  string
		expectedDepth int
	}{
		{"Simple array", "string[]", "string", 1},
		{"Two-dimensional array", "int[][]", "int", 2},
		{"Normal type", "boolean", "boolean", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := schema.Field{Type: tt.fieldType}
			baseType := field.GetBaseType()
			depth := field.GetArrayDepth()

			assert.Equal(t, tt.expectedBase, baseType)
			assert.Equal(t, tt.expectedDepth, depth)
			// Verify original field remains unchanged
			assert.Equal(t, tt.fieldType, field.Type)
		})
	}
}

func TestFieldArrayOperationsWithRules(t *testing.T) {
	field := schema.Field{
		Type: "string[]",
		Rules: []schema.RuleCatchAll{
			{
				Name:    schema.RuleNameMinLength,
				Value:   "3",
				Message: "Must be at least 3 characters",
			},
		},
	}

	assert.True(t, field.IsArray())
	assert.Equal(t, 1, field.GetArrayDepth())
	assert.Equal(t, "string", field.GetBaseType())

	assert.Len(t, field.Rules, 1)
	assert.Equal(t, schema.RuleNameMinLength, field.Rules[0].Name)
}

func TestFieldArrayEdgeCases(t *testing.T) {
	// Extremely nested array
	deepArray := schema.Field{Type: "int[][][][][][][][][][]"} // 10 levels of nesting
	assert.True(t, deepArray.IsArray())
	assert.Equal(t, 10, deepArray.GetArrayDepth())
	assert.Equal(t, "int", deepArray.GetBaseType())

	// Type with special characters
	specialChars := schema.Field{Type: "map<string,int>[]"}
	assert.True(t, specialChars.IsArray())
	assert.Equal(t, 1, specialChars.GetArrayDepth())
	assert.Equal(t, "map<string,int>", specialChars.GetBaseType())

	// Type with brackets in name
	bracketInName := schema.Field{Type: "Array[T][]"}
	assert.True(t, bracketInName.IsArray())
	assert.Equal(t, 1, bracketInName.GetArrayDepth())
	assert.Equal(t, "Array[T]", bracketInName.GetBaseType())
}

func TestFieldWithNestedFields(t *testing.T) {
	field := schema.Field{
		Type: "object[]",
		Fields: map[string]schema.Field{
			"name":   {Type: "string"},
			"scores": {Type: "int[]"},
		},
	}

	assert.True(t, field.IsArray())
	assert.Equal(t, 1, field.GetArrayDepth())
	assert.Equal(t, "object", field.GetBaseType())

	assert.Equal(t, "string", field.Fields["name"].Type)
	assert.False(t, field.Fields["name"].IsArray())

	assert.Equal(t, "int[]", field.Fields["scores"].Type)
	assert.True(t, field.Fields["scores"].IsArray())
	assert.Equal(t, 1, field.Fields["scores"].GetArrayDepth())
	assert.Equal(t, "int", field.Fields["scores"].GetBaseType())
}

func TestFieldTypeModifications(t *testing.T) {
	original := schema.Field{Type: "string[][]"}
	modified := original

	assert.True(t, original.IsArray())
	assert.Equal(t, 2, original.GetArrayDepth())
	assert.Equal(t, "string", original.GetBaseType())

	baseType := modified.GetBaseType()
	assert.Equal(t, "string", baseType)

	// Verify original field remains unchanged
	assert.Equal(t, "string[][]", original.Type)
	assert.Equal(t, "string[][]", modified.Type)

	assert.True(t, modified.IsArray())
	assert.Equal(t, 2, modified.GetArrayDepth())
	assert.Equal(t, "string", modified.GetBaseType())
}
