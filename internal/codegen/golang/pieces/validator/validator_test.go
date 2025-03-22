package validator

import (
	"encoding/json"
	"testing"
)

func TestStringValidation(t *testing.T) {
	t.Run("should validate required string", func(t *testing.T) {
		schema := schValidator.String("").Required("")
		result := schema.Validate("Hola Mundo")
		if !result.IsValid {
			t.Error("Expected string to be valid")
		}

		result2 := schema.Validate("")
		if !result2.IsValid {
			t.Error("Expected empty string to be valid")
		}

		result3 := schema.Validate(" ")
		if !result3.IsValid {
			t.Error("Expected space string to be valid")
		}
	})

	t.Run("should invalidate undefined when string is required", func(t *testing.T) {
		schema := schValidator.String("").Required("")

		result1 := schema.Validate(nil)
		if result1.IsValid {
			t.Error("Expected nil to be invalid")
		}
	})

	t.Run("should invalidate required with custom message", func(t *testing.T) {
		customMsg := "My custom error message"
		schema := schValidator.String("").Required(customMsg)

		result := schema.Validate(nil)
		if result.IsValid {
			t.Error("Expected nil to be invalid")
		}
		if result.Error != customMsg {
			t.Errorf("Expected error message '%s', got '%s'", customMsg, result.Error)
		}
	})

	t.Run("should validate string matching regex pattern", func(t *testing.T) {
		schema := schValidator.String("").Regex("^[A-Z]+$", "")
		result := schema.Validate("HOLA")
		if !result.IsValid {
			t.Error("Expected uppercase string to be valid")
		}
	})

	t.Run("should invalidate string not matching regex pattern", func(t *testing.T) {
		schema := schValidator.String("").Regex("^[A-Z]+$", "")
		result := schema.Validate("Hola")
		if result.IsValid {
			t.Error("Expected mixed case string to be invalid")
		}
	})

	t.Run("should validate string with length", func(t *testing.T) {
		schema := schValidator.String("").Length(5, "")
		result := schema.Validate("Hola!")
		if !result.IsValid {
			t.Error("Expected string with correct length to be valid")
		}
	})

	t.Run("should invalidate string different than length", func(t *testing.T) {
		schema := schValidator.String("").Length(5, "")
		result := schema.Validate("Hola Mundo!!")
		if result.IsValid {
			t.Error("Expected string with incorrect length to be invalid")
		}
	})

	t.Run("should validate string with minLength", func(t *testing.T) {
		schema := schValidator.String("").MinLength(5, "")
		result := schema.Validate("Hola Mundo")
		if !result.IsValid {
			t.Error("Expected string with sufficient length to be valid")
		}
	})

	t.Run("should invalidate string shorter than minLength", func(t *testing.T) {
		schema := schValidator.String("").MinLength(10, "")
		result := schema.Validate("Hola")
		if result.IsValid {
			t.Error("Expected short string to be invalid")
		}
	})

	t.Run("should validate email format - basic case", func(t *testing.T) {
		schema := schValidator.String("").Email("")
		tests := []string{
			"test@example.com",
			"user.name@domain.com",
			"user+label@domain.co.uk",
			"first.last@subdomain.domain.org",
		}

		for _, test := range tests {
			result := schema.Validate(test)
			if !result.IsValid {
				t.Errorf("Expected email '%s' to be valid", test)
			}
		}
	})

	t.Run("should invalidate incorrect email formats", func(t *testing.T) {
		schema := schValidator.String("").Email("")
		tests := []string{
			"test@example",
			"test@.com",
			"@domain.com",
			"test@domain..com",
			"test@dom ain.com",
			"te st@domain.com",
		}

		for _, test := range tests {
			result := schema.Validate(test)
			if result.IsValid {
				t.Errorf("Expected email '%s' to be invalid", test)
			}
		}
	})
}

func TestNumberValidation(t *testing.T) {
	t.Run("should validate required number", func(t *testing.T) {
		schema := schValidator.Number("").Required("")
		result := schema.Validate(100.0)
		if !result.IsValid {
			t.Error("Expected number to be valid")
		}
	})

	t.Run("should invalidate undefined when number is required", func(t *testing.T) {
		schema := schValidator.Number("").Required("")
		result := schema.Validate(nil)
		if result.IsValid {
			t.Error("Expected nil to be invalid")
		}
	})

	t.Run("should validate number with min value", func(t *testing.T) {
		minVal := 50.0
		schema := schValidator.Number("")
		schema.MinValue = &minVal
		result := schema.Validate(75.0)
		if !result.IsValid {
			t.Error("Expected number above minimum to be valid")
		}
	})

	t.Run("should invalidate number less than min value", func(t *testing.T) {
		minVal := 50.0
		schema := schValidator.Number("")
		schema.MinValue = &minVal
		result := schema.Validate(25.0)
		if result.IsValid {
			t.Error("Expected number below minimum to be invalid")
		}
	})
}

func TestBooleanValidation(t *testing.T) {
	t.Run("should validate boolean true", func(t *testing.T) {
		schema := schValidator.Boolean("")
		result := schema.Validate(true)
		if !result.IsValid {
			t.Error("Expected true to be valid")
		}
	})

	t.Run("should validate boolean false", func(t *testing.T) {
		schema := schValidator.Boolean("")
		result := schema.Validate(false)
		if !result.IsValid {
			t.Error("Expected false to be valid")
		}
	})

	t.Run("should invalidate non-boolean value", func(t *testing.T) {
		schema := schValidator.Boolean("")
		result := schema.Validate("true")
		if result.IsValid {
			t.Error("Expected string to be invalid")
		}
	})
}

func TestArrayValidation(t *testing.T) {
	t.Run("should validate array of integers", func(t *testing.T) {
		intSchema := schValidator.Int("")
		schema := schValidator.Array(intSchema, "")
		result := schema.Validate([]interface{}{1.0, 2.0, 3.0}) // Note: JSON numbers are float64
		if !result.IsValid {
			t.Error("Expected array of integers to be valid")
		}
	})

	t.Run("should invalidate array with invalid element", func(t *testing.T) {
		intSchema := schValidator.Int("")
		schema := schValidator.Array(intSchema, "")
		result := schema.Validate([]interface{}{1.0, 2.0, "3"})
		if result.IsValid {
			t.Error("Expected array with string element to be invalid")
		}
	})

	t.Run("should invalidate non-array value", func(t *testing.T) {
		intSchema := schValidator.Int("")
		schema := schValidator.Array(intSchema, "")
		result := schema.Validate("not an array")
		if result.IsValid {
			t.Error("Expected string to be invalid")
		}
	})
}

func TestObjectValidation(t *testing.T) {
	t.Run("should validate object with defined schema", func(t *testing.T) {
		schema := schValidator.Object(map[string]*schemaValidator{
			"nombre": schValidator.String("").Required(""),
			"edad":   schValidator.Int(""),
		}, "")

		jsonStr := `{"nombre": "Juan", "edad": 30}`
		var data map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &data)

		result := schema.Validate(data)
		if !result.IsValid {
			t.Error("Expected valid object to pass validation")
		}
	})

	t.Run("should invalidate object with missing required field", func(t *testing.T) {
		schema := schValidator.Object(map[string]*schemaValidator{
			"nombre": schValidator.String("").Required(""),
			"edad":   schValidator.Int(""),
		}, "")

		jsonStr := `{"edad": 30}`
		var data map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &data)

		result := schema.Validate(data)
		if result.IsValid {
			t.Error("Expected object with missing required field to fail validation")
		}
	})

	t.Run("should validate nested object", func(t *testing.T) {
		direccionSchema := schValidator.Object(map[string]*schemaValidator{
			"calle":  schValidator.String("").Required(""),
			"ciudad": schValidator.String("").Required(""),
		}, "")

		schema := schValidator.Object(map[string]*schemaValidator{
			"nombre":    schValidator.String("").Required(""),
			"direccion": direccionSchema,
		}, "")

		jsonStr := `{
			"nombre": "Juan",
			"direccion": {
				"calle": "Calle Principal 123",
				"ciudad": "Ciudad"
			}
		}`
		var data map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &data)

		result := schema.Validate(data)
		if !result.IsValid {
			t.Error("Expected valid nested object to pass validation")
		}
	})
}

func TestEnumValidation(t *testing.T) {
	t.Run("should validate value in enum", func(t *testing.T) {
		schema := schValidator.String("").Enum([]interface{}{"rojo", "verde", "azul"}, "")
		result := schema.Validate("rojo")
		if !result.IsValid {
			t.Error("Expected enum value to be valid")
		}
	})

	t.Run("should invalidate value not in enum", func(t *testing.T) {
		schema := schValidator.String("").Enum([]interface{}{"rojo", "verde", "azul"}, "")
		result := schema.Validate("amarillo")
		if result.IsValid {
			t.Error("Expected non-enum value to be invalid")
		}
	})
}

func TestLazyValidation(t *testing.T) {
	t.Run("should validate recursive structures using lazy schema", func(t *testing.T) {
		var nodoSchema *schemaValidator
		nodoSchema = schValidator.Lazy(func() *schemaValidator {
			return schValidator.Object(map[string]*schemaValidator{
				"valor": schValidator.Number("").Required(""),
				"siguiente": schValidator.Lazy(func() *schemaValidator {
					return nodoSchema
				}, ""),
			}, "")
		}, "")

		jsonStr := `{
			"valor": 1,
			"siguiente": {
				"valor": 2,
				"siguiente": {
					"valor": 3
				}
			}
		}`
		var data map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &data)

		result := nodoSchema.Validate(data)
		if !result.IsValid {
			t.Error("Expected valid recursive structure to pass validation")
		}
	})

	t.Run("should invalidate recursive structure with invalid data", func(t *testing.T) {
		var nodoSchema *schemaValidator
		nodoSchema = schValidator.Lazy(func() *schemaValidator {
			return schValidator.Object(map[string]*schemaValidator{
				"valor": schValidator.Number("").Required(""),
				"siguiente": schValidator.Lazy(func() *schemaValidator {
					return nodoSchema
				}, ""),
			}, "")
		}, "")

		jsonStr := `{
			"valor": 1,
			"siguiente": {
				"valor": "invalido"
			}
		}`
		var data map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &data)

		result := nodoSchema.Validate(data)
		if result.IsValid {
			t.Error("Expected invalid recursive structure to fail validation")
		}
	})
}
