package cleaner

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/urpc/parser"
	"github.com/uforg/uforpc/internal/util/testutil"
)

func assertClean(t *testing.T, input string, expected string) {
	t.Helper()

	// Parse the input
	schema, err := parser.ParserInstance.ParseString("test.urpc", input)
	require.NoError(t, err)

	// Clean the schema
	cleaned := Clean(schema)

	// Parse the expected output
	expectedSchema, err := parser.ParserInstance.ParseString("test.urpc", expected)
	require.NoError(t, err)

	// Compare the cleaned schema with the expected schema
	testutil.ASTEqualNoPos(t, expectedSchema, cleaned)
}

func TestClean(t *testing.T) {
	t.Run("removes unused rules and types", func(t *testing.T) {
		input := `
			version 1
			
			rule @usedRule { // Used rule
				for: string
				error: "Used rule error"
			}

			rule @unusedRule { // Unused rule
				for: string
				error: "Unused rule error"
			}

			type UsedType { // Used type
				field: string
					@usedRule
			}

			type UnusedType { // Unused type
				field: string
			}

			type TypeWithReference { // Type that references another type
				field: UsedType
			}

			proc MyProc { // Procedure that uses a type
				input {
					field: TypeWithReference
				}
			}
		`

		expected := `
			version 1
			
			rule @usedRule { // Used rule
				for: string
				error: "Used rule error"
			}

			type UsedType { // Used type
				field: string
					@usedRule
			}

			type TypeWithReference { // Type that references another type
				field: UsedType
			}

			proc MyProc { // Procedure that uses a type
				input {
					field: TypeWithReference
				}
			}
		`

		assertClean(t, input, expected)
	})

	t.Run("remove even when referenced in unused", func(t *testing.T) {
		input := `
			version 1

			rule @validateCustomType { // Rule for custom type
				for: CustomType
				error: "Invalid custom type"
			}

			type CustomType { // Custom type referenced by rule
				field: string
			}

			type AnotherType { // Type that uses the rule
				field: CustomType
					@validateCustomType
			}
		`

		expected := `
			version 1
		`

		assertClean(t, input, expected)
	})

	t.Run("handles type extensions", func(t *testing.T) {
		input := `
			version 1

			type BaseType { // Base type
				field: string
			}

			type ExtendingType extends BaseType { // Type that extends base type
				anotherField: string
			}

			type UnusedType { // Unused type
				field: string
			}

			proc MyProc { // Procedure that uses the extending type
				input {
					field: ExtendingType
				}
			}
		`

		expected := `
			version 1

			type BaseType { // Base type
				field: string
			}

			type ExtendingType extends BaseType { // Type that extends base type
				anotherField: string
			}

			proc MyProc { // Procedure that uses the extending type
				input {
					field: ExtendingType
				}
			}
		`

		assertClean(t, input, expected)
	})

	t.Run("handles inline objects", func(t *testing.T) {
		input := `
			version 1

			rule @usedRule { // Used rule
				for: string
				error: "Used rule error"
			}

			rule @unusedRule { // Unused rule
				for: string
				error: "Unused rule error"
			}

			type UsedType { // Used type
				field: string
			}

			type UnusedType { // Unused type
				field: string
			}

			type TypeWithInline { // Type with inline object that references another type
				inlineField: {
					nestedField: UsedType
						@usedRule
				}
			}
			
			proc MyProc { // Procedure that uses a type
				input {
					field: TypeWithInline
				}
			}
		`

		expected := `
			version 1

			rule @usedRule { // Used rule
				for: string
				error: "Used rule error"
			}

			type UsedType { // Used type
				field: string
			}

			type TypeWithInline { // Type with inline object that references another type
				inlineField: {
					nestedField: UsedType
						@usedRule
				}
			}
			
			proc MyProc { // Procedure that uses a type
				input {
					field: TypeWithInline
				}
			}
		`

		assertClean(t, input, expected)
	})

	t.Run("handles procedure input and output", func(t *testing.T) {
		input := `
			version 1

			rule @usedRule { // Used rule
				for: string
				error: "Used rule error"
			}

			rule @unusedRule { // Unused rule
				for: string
				error: "Unused rule error"
			}

			type InputType { // Used type in input
				field: string
					@usedRule
			}

			type OutputType { // Used type in output
				field: string
			}

			type UnusedType { // Unused type
				field: string
			}
			
			proc MyProc { // Procedure that uses types in input and output
				input {
					field: InputType
				}
				output {
					field: OutputType
				}
			}
		`

		expected := `
			version 1

			rule @usedRule { // Used rule
				for: string
				error: "Used rule error"
			}

			type InputType { // Used type in input
				field: string
					@usedRule
			}

			type OutputType { // Used type in output
				field: string
			}
			
			proc MyProc { // Procedure that uses types in input and output
				input {
					field: InputType
				}
				output {
					field: OutputType
				}
			}
		`

		assertClean(t, input, expected)
	})
}
