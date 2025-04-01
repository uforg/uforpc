package analyzer

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/urpc/parser"
)

func TestAnalyzer_ValidVersion(t *testing.T) {
	input := `
		version 1
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestAnalyzer_InvalidVersion(t *testing.T) {
	input := `
		version 2 // Only version 1 is valid
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "version must be 1")
}

func TestAnalyzer_ValidImport(t *testing.T) {
	input := `
		version 1

		import "other.urpc"
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestAnalyzer_EmptyImportPath(t *testing.T) {
	input := `
		version 1

		import ""
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)

	errors, err := analyzer.Analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "import path is required")
}

func TestAnalyzer_InvalidImportExtension(t *testing.T) {
	input := `
		version 1

		import "other.txt" // Not a .urpc file
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "import path must end with .urpc")
}

func TestAnalyzer_DuplicateImportPath(t *testing.T) {
	input := `
		version 1

		import "other.urpc"
		import "other.urpc" // Duplicate
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "import path must be unique")
}

func TestAnalyzer_ValidCustomRule(t *testing.T) {
	input := `
		version 1

		rule @customStringRule {
		  for: string
		  error: "Invalid string format"
		}

		rule @customIntRule {
		  for: int
		  error: "Invalid int format"
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestAnalyzer_DuplicateCustomRule(t *testing.T) {
	input := `
		version 1

		rule @customStringRule {
		  for: string
		  error: "Invalid string format"
		}

		rule @customStringRule { // Duplicate rule name
		  for: string
		  error: "Another invalid string format"
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "custom rule name is already defined")
}

func TestAnalyzer_InvalidCustomRuleName(t *testing.T) {
	input := `
		version 1

		rule @InvalidRuleName { // PascalCase, should be camelCase
		  for: string
		  error: "Invalid string format"
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "custom rule name must be in camelCase")
}

func TestAnalyzer_ValidCustomType(t *testing.T) {
	input := `
		version 1

		type User {
		  id: string
		  age: int
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestAnalyzer_DuplicateCustomType(t *testing.T) {
	input := `
		version 1

		type User {
		  id: string
		  age: int
		}

		type User { // Duplicate type name
		  name: string
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "custom type name is already defined")
}

func TestAnalyzer_InvalidCustomTypeName(t *testing.T) {
	input := `
		version 1

		type invalidTypeName { // camelCase, should be PascalCase
		  id: string
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "custom type name must be in PascalCase")
}

func TestAnalyzer_ValidProcedure(t *testing.T) {
	input := `
		version 1

		type User {
		  id: string
		  age: int
		}

		proc GetUser {
		  input {
		    userId: string
		  }
		  
		  output {
		    user: User
		  }
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestAnalyzer_DuplicateProcedure(t *testing.T) {
	input := `
		version 1

		proc GetUser {
		  input {
		    userId: string
		  }
		}

		proc GetUser { // Duplicate procedure name
		  input {
		    id: string
		  }
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "procedure is already defined")
}

func TestAnalyzer_InvalidProcedureName(t *testing.T) {
	input := `
		version 1

		proc getUser { // camelCase, should be PascalCase
		  input {
		    userId: string
		  }
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "custom procedure name must be in PascalCase")
}

func TestAnalyzer_NonExistentRuleReference(t *testing.T) {
	input := `
		version 1

		type User {
		  id: string
		    @nonExistentRule // This rule doesn't exist
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "referenced rule \"nonExistentRule\" in type \"User\" is not defined")
}

func TestAnalyzer_RuleTypeMismatch(t *testing.T) {
	input := `
		version 1

		rule @customStringRule {
		  for: string
		}

		type User {
		  id: int
		    @customStringRule // Rule is for string, but field is int
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "rule \"customStringRule\" in type \"User\" cannot be applied to type \"int\"")
}

func TestAnalyzer_NonExistentTypeReference(t *testing.T) {
	input := `
		version 1

		proc GetPost {
		  output {
		    post: Post // This type doesn't exist
		  }
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "referenced type \"Post\" in output of procedure \"GetPost\" is not defined")
}

func TestAnalyzer_InvalidTypeExtends(t *testing.T) {
	input := `
		version 1

		type UserExtended extends NonExistentType {
		  email: string
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "type \"UserExtended\" extends undefined type \"NonExistentType\"")
}

func TestAnalyzer_ValidArrayWithRules(t *testing.T) {
	input := `
		version 1

		type User {
		  tags: string[]
		    @minlen(1) // Valid rule for arrays
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestAnalyzer_InvalidArrayRule(t *testing.T) {
	input := `
		version 1

		type User {
		  tags: string[]
		    @contains("tag") // This rule is not valid for arrays (only for strings)
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "rule \"contains\" in type \"User\" cannot be applied to array type \"string[]\"")
}

func TestAnalyzer_BuiltInRuleValidation(t *testing.T) {
	input := `
		version 1

		type Validation {
		  email: string
		    @minlen(3)
		    @maxlen(100)
		    @contains("@")
		  
		  age: int
		    @min(18)
		    @max(120)
		  
		  price: float
		    @min(0.0)
		    @max(999.99)
		  
		  isActive: boolean
		    @equals(true)
		  
		  tags: string[]
		    @minlen(1)
		    @maxlen(10)
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}
