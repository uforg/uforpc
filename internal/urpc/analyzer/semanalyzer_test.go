package analyzer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/urpc/ast"
	"github.com/uforg/uforpc/internal/urpc/parser"
)

// createCombinedSchema creates a CombinedSchema from the given schema.
// this is a test helper function.
func parseToCombinedSchema(input string) (CombinedSchema, error) {
	schema, err := parser.ParserInstance.ParseString("test.urpc", input)
	if err != nil {
		return CombinedSchema{}, err
	}

	// Collect all declarations from the combined schema
	ruleDecls := make(map[string]*ast.RuleDecl)
	typeDecls := make(map[string]*ast.TypeDecl)
	procDecls := make(map[string]*ast.ProcDecl)

	// Collect rule declarations
	for _, rule := range schema.GetRules() {
		ruleDecls[rule.Name] = rule
	}

	// Collect type declarations
	for _, typeDecl := range schema.GetTypes() {
		typeDecls[typeDecl.Name] = typeDecl
	}

	// Collect procedure declarations
	for _, proc := range schema.GetProcs() {
		procDecls[proc.Name] = proc
	}

	return CombinedSchema{
		Schema:    schema,
		RuleDecls: ruleDecls,
		TypeDecls: typeDecls,
		ProcDecls: procDecls,
	}, nil
}

func TestSemanalyzer_ValidVersion(t *testing.T) {
	input := `
		version 1
	`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_ValidCustomRule(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_DuplicateCustomRule(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "already declared")
}

func TestSemanalyzer_InvalidCustomRuleName(t *testing.T) {
	input := `
		version 1

		// PascalCase, should be camelCase
		rule @InvalidRuleName {
		  for: string
		  error: "Invalid string format"
		}
	`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "must be in camelCase")
}

func TestSemanalyzer_ValidCustomType(t *testing.T) {
	input := `
		version 1

		type User {
		  id: string
		  age: int
		}
	`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_DuplicateCustomType(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "is already declared")
}

func TestSemanalyzer_InvalidCustomTypeName(t *testing.T) {
	input := `
		version 1

		// camelCase, should be PascalCase
		type invalidTypeName {
		  id: string
		}
	`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "must be in PascalCase")
}

func TestSemanalyzer_ValidProcedure(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_DuplicateProcedure(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "is already declared")
}

func TestSemanalyzer_InvalidProcedureName(t *testing.T) {
	input := `
		version 1

		// camelCase, should be PascalCase
		proc getUser {
		  input {
		    userId: string
		  }
		}
	`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "must be in PascalCase")
}

func TestSemanalyzer_NonExistentRuleReference(t *testing.T) {
	input := `
		version 1

		// This rule doesn't exist
		type User {
		  id: string
		    @nonExistentRule
		}
	`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "is not declared")
}

func TestSemanalyzer_RuleTypeMismatch(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "cannot be applied to type")
}

func TestSemanalyzer_NonExistentTypeReference(t *testing.T) {
	input := `
		version 1

		// This type doesn't exist
		proc GetPost {
		  output {
		    post: Post
		  }
		}
	`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "is not declared")
}

func TestSemanalyzer_InvalidTypeExtends(t *testing.T) {
	input := `
		version 1

		type UserExtended extends NonExistentType {
		  email: string
		}
	`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "type \"UserExtended\" extends non-declared type \"NonExistentType\"")
}

func TestSemanalyzer_ValidArrayWithRules(t *testing.T) {
	input := `
		version 1

		// Valid rule for arrays
		type User {
		  tags: string[]
		    @minlen(1)
		}
	`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_InvalidArrayRule(t *testing.T) {
	input := `
		version 1

		// This rule is not valid for arrays (only for strings)
		type User {
		  tags: string[]
		    @contains("tag")
		}
	`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "cannot be applied to array type")
}

func TestSemanalyzer_BuiltInRuleValidation(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_ValidTypeExtends(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_MultipleTypeExtends(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_ValidProcedureMeta(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_ValidInlineObject(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_ValidNestedInlineObject(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_ValidInlineObjectWithRules(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_InvalidRuleInInlineObject(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "is not declared")
}

func TestSemanalyzer_ValidArrayOfArrays(t *testing.T) {
	input := `
		version 1

		type Matrix {
		  data: int[][]
		    @minlen(1)
		}
	`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_ValidArrayOfObjects(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_CustomRuleWithParameters(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_CustomRuleWithArrayParameter(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_EnumValidation(t *testing.T) {
	input := `
		version 1

		type Product {
		  status: string
		    @enum(["pending", "approved", "rejected"])

		  priority: int
		    @enum([1, 2, 3])
		}
	`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_DatetimeValidation(t *testing.T) {
	input := `
		version 1

		type Event {
		  startDate: datetime
		    @min("2023-01-01T00:00:00Z")

		  endDate: datetime
		    @max("2030-12-31T23:59:59Z")
		}
	`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_OptionalFields(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_RecursiveTypeExtension(t *testing.T) {
	input := `
			version 1

			// Direct recursive extension
			type User extends User {
			  id: string
			}
		`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.NotEmpty(t, errors)

	// Check that at least one error contains the expected message
	found := false
	for _, diag := range errors {
		if strings.Contains(diag.Message, "recursive type extension detected") {
			found = true
			break
		}
	}
	require.True(t, found, "Expected to find an error about recursive type extension")
}

func TestSemanalyzer_IndirectRecursiveTypeExtension(t *testing.T) {
	input := `
			version 1

			// Indirect recursive extension: A extends B, B extends C, C extends A
			type A extends B {
			  id: string
			}

			type B extends C {
			  name: string
			}

			type C extends A {
			  age: int
			}
		`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.NotEmpty(t, errors)
	require.Contains(t, errors[0].Message, "recursive type extension detected")
}

func TestSemanalyzer_DuplicateFieldsInTypeExtension(t *testing.T) {
	input := `
			version 1

			type BaseUser {
			  id: string
			  email: string
			}

			// Duplicate field 'email' from BaseUser
			type User extends BaseUser {
			  name: string
			  email: string  // This is a duplicate
			}
		`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "is already defined in extended type")
}

func TestSemanalyzer_CircularTypeDependency(t *testing.T) {
	input := `
			version 1

			// Circular dependency: User -> Post -> User
			type User {
			  id: string
			  posts: Post[]
			}

			type Post {
			  id: string
			  author: User  // This creates a circular dependency
			}
		`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.NotEmpty(t, errors)

	// Check that at least one error contains the expected message
	found := false
	for _, diag := range errors {
		if strings.Contains(diag.Message, "circular dependency detected between types") {
			found = true
			break
		}
	}
	require.True(t, found, "Expected to find an error about circular dependency")
}

func TestSemanalyzer_CircularTypeDependencyWithOptionalField(t *testing.T) {
	input := `
			version 1

			// Circular dependency with optional field
			type User {
			  id: string
			  posts: Post[]
			}

			type Post {
			  id: string
			  author?: User  // Optional field breaks the circular dependency
			}
		`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}

func TestSemanalyzer_RuleWithoutForClause(t *testing.T) {
	input := `
			version 1

			// Rule without 'for' clause
			rule @invalidRule {
			  error: "Invalid value"
			}
		`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "must have exactly one 'for' clause")
}

func TestSemanalyzer_RuleWithMultipleForClauses(t *testing.T) {
	input := `
			version 1

			// Rule with multiple 'for' clauses
			rule @invalidRule {
			  for: string
			  for: int  // Duplicate 'for' clause
			  error: "Invalid value"
			}
		`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "cannot have more than one 'for' clause")
}

func TestSemanalyzer_RuleWithMultipleParamClauses(t *testing.T) {
	input := `
			version 1

			// Rule with multiple 'param' clauses
			rule @invalidRule {
			  for: string
			  param: string
			  param: int  // Duplicate 'param' clause
			  error: "Invalid value"
			}
		`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "cannot have more than one 'param' clause")
}

func TestSemanalyzer_RuleWithMultipleErrorClauses(t *testing.T) {
	input := `
			version 1

			// Rule with multiple 'error' clauses
			rule @invalidRule {
			  for: string
			  error: "Invalid value"
			  error: "Another error message"  // Duplicate 'error' clause
			}
		`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "cannot have more than one 'error' clause")
}

func TestSemanalyzer_ProcWithMultipleInputSections(t *testing.T) {
	input := `
			version 1

			// Procedure with multiple 'input' sections
			proc InvalidProc {
			  input {
			    id: string
			  }

			  input {  // Duplicate 'input' section
			    name: string
			  }
			}
		`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "cannot have more than one 'input' section")
}

func TestSemanalyzer_ProcWithMultipleOutputSections(t *testing.T) {
	input := `
			version 1

			// Procedure with multiple 'output' sections
			proc InvalidProc {
			  output {
			    success: boolean
			  }

			  output {  // Duplicate 'output' section
			    message: string
			  }
			}
		`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "cannot have more than one 'output' section")
}

func TestSemanalyzer_ProcWithMultipleMetaSections(t *testing.T) {
	input := `
			version 1

			// Procedure with multiple 'meta' sections
			proc InvalidProc {
			  meta {
			    requiresAuth: true
			  }

			  meta {  // Duplicate 'meta' section
			    timeout: 30
			  }
			}
		`
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.Error(t, err)
	require.Len(t, errors, 1)
	require.Contains(t, errors[0].Message, "cannot have more than one 'meta' section")
}

func TestSemanalyzer_CompleteSchema(t *testing.T) {
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
	combinedSchema, err := parseToCombinedSchema(input)
	require.NoError(t, err)

	analyzer := newSemanalyzer(combinedSchema)
	errors, err := analyzer.analyze()

	require.NoError(t, err)
	require.Empty(t, errors)
}
