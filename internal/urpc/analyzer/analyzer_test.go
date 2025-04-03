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

		// Duplicate rule name
		rule @customStringRule {
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

		// PascalCase, should be camelCase
		rule @InvalidRuleName {
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

		// Duplicate type name
		type User {
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

		// camelCase, should be PascalCase
		type invalidTypeName {
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

		// Duplicate procedure name
		proc GetUser {
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

		// camelCase, should be PascalCase
		proc getUser {
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

		// This rule doesn't exist
		type User {
		  id: string
		    @nonExistentRule
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

		// Rule is for string, but field is int
		type User {
		  id: int
		    @customStringRule
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

		// This type doesn't exist
		proc GetPost {
		  output {
		    post: Post
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

		// Valid rule for arrays
		type User {
		  tags: string[]
		    @minlen(1)
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

		// This rule is not valid for arrays (only for strings)
		type User {
		  tags: string[]
		    @contains("tag")
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

func TestAnalyzer_ValidTypeExtends(t *testing.T) {
	input := `
		version 1

		type BaseUser {
		  id: string
		  username: string
		}

		type ExtendedUser extends BaseUser {
		  email: string
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

func TestAnalyzer_MultipleTypeExtends(t *testing.T) {
	input := `
		version 1

		type Base1 {
		  field1: string
		}

		type Base2 {
		  field2: int
		}

		type Combined extends Base1, Base2 {
		  field3: boolean
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestAnalyzer_ValidProcedureMeta(t *testing.T) {
	input := `
		version 1

		proc DoSomething {
		  input {
		    data: string
		  }

		  meta {
		    requiresAuth: true
		    maxRetries: 3
		    timeout: 60
		    description: "This is a test procedure"
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

func TestAnalyzer_ValidInlineObject(t *testing.T) {
	input := `
		version 1

		type User {
		  id: string
		  address: {
		    street: string
		    city: string
		    zipCode: string
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

func TestAnalyzer_ValidNestedInlineObject(t *testing.T) {
	input := `
		version 1

		type User {
		  id: string
		  contact: {
		    email: string
		    address: {
		      street: string
		      city: string
		      country: string
		    }
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

func TestAnalyzer_ValidInlineObjectWithRules(t *testing.T) {
	input := `
		version 1

		type User {
		  id: string
		  contact: {
		    email: string
		      @contains("@")
		      @minlen(5)
		    phone: string
		      @minlen(10)
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

func TestAnalyzer_InvalidRuleInInlineObject(t *testing.T) {
	input := `
		version 1

		// This rule doesn't exist
		type User {
		  id: string
		  contact: {
		    email: string
		      @invalidRule
		  }
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "referenced rule \"invalidRule\" in inline object in field \"contact\" is not defined")
}

func TestAnalyzer_ValidArrayOfArrays(t *testing.T) {
	input := `
		version 1

		type Matrix {
		  data: int[][]
		    @minlen(1)
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestAnalyzer_ValidArrayOfObjects(t *testing.T) {
	input := `
		version 1

		type User {
		  id: string
		  addresses: {
		    street: string
		    city: string
		  }[]
		    @minlen(1)
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestAnalyzer_CustomRuleWithParameters(t *testing.T) {
	input := `
		version 1

		rule @regex {
		  for: string
		  param: string
		  error: "Invalid format"
		}

		type User {
		  email: string
		    @regex("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", error: "Invalid email format")
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestAnalyzer_CustomRuleWithArrayParameter(t *testing.T) {
	input := `
		version 1

		rule @range {
		  for: int
		  param: int[]
		  error: "Value out of range"
		}

		type Product {
		  price: int
		    @range([10, 1000], error: "Price must be between 10 and 1000")
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestAnalyzer_EnumValidation(t *testing.T) {
	input := `
		version 1

		type Product {
		  status: string
		    @enum(["pending", "approved", "rejected"])
		    
		  priority: int
		    @enum([1, 2, 3])
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestAnalyzer_DatetimeValidation(t *testing.T) {
	input := `
		version 1

		type Event {
		  startDate: datetime
		    @min("2023-01-01T00:00:00Z")
		    
		  endDate: datetime
		    @max("2030-12-31T23:59:59Z")
		}
	`
	schema, err := parser.Parser.ParseString("test.urpc", input)
	require.NoError(t, err)

	analyzer := NewAnalyzer(schema)
	errors, err := analyzer.Analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestAnalyzer_OptionalFields(t *testing.T) {
	input := `
		version 1

		type User {
		  id: string
		  email: string
		  phone?: string
		  address?: {
		    street: string
		    city: string
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

func TestAnalyzer_CompleteSchema(t *testing.T) {
	input := `
		version 1

		rule @regex {
		  for: string
		  param: string
		  error: "Invalid format"
		}

		rule @range {
		  for: int
		  param: int[]
		  error: "Value out of range"
		}

		type Address {
		  street: string
		    @minlen(3)
		  city: string
		    @minlen(2)
		  zipCode: string
		    @regex("^\\d{5}$", error: "Zip code must be 5 digits")
		}

		type BaseUser {
		  id: string
		    @minlen(10)
		  username: string
		    @minlen(3)
		    @maxlen(30)
		    @regex("^[a-zA-Z0-9_]+$", error: "Username can only contain letters, numbers, and underscores")
		}

		type User extends BaseUser {
		  email: string
		    @contains("@")
		    @minlen(5)
		  password: string
		    @minlen(8)
		  age: int
		    @range([18, 120], error: "Age must be between 18 and 120")
		  isActive: boolean
		    @equals(true)
		  address: Address
		  tags: string[]
		    @minlen(1)
		    @maxlen(10)
		  metadata: {
		    lastLogin: datetime
		      @min("2020-01-01T00:00:00Z")
		    preferences: {
		      theme: string
		        @enum(["light", "dark", "system"])
		      notifications: boolean
		    }
		  }
		}

		proc CreateUser {
		  input {
		    user: User
		  }

		  output {
		    success: boolean
		    userId: string
		    errors: string[]
		  }

		  meta {
		    requiresAuth: false
		    rateLimit: 10
		    description: "Creates a new user in the system"
		  }
		}

		proc GetUser {
		  input {
		    userId: string
		      @minlen(10)
		  }

		  output {
		    user: User
		  }

		  meta {
		    requiresAuth: true
		    cacheTTL: 300
		    description: "Retrieves a user by ID"
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
