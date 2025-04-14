package typeflattener

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/urpc/parser"
	"github.com/uforg/uforpc/internal/util/testutil"
)

// assertFlattenedSchema receives two strings representing the expected
// flattened schema and the actual unflattened schema.
//
// It parses both strings, flattens the actual schema using the Flatten function,
// and then compares the flattened schema with the expected flattened schema.
func assertFlattenedSchema(t *testing.T, expected string, actual string) {
	t.Helper()

	// Parse the expected schema
	expectedSchema, err := parser.ParserInstance.ParseString("test.urpc", expected)
	require.NoError(t, err)

	// Parse the actual schema
	actualSchema, err := parser.ParserInstance.ParseString("test.urpc", actual)
	require.NoError(t, err)

	// Flatten the actual schema
	actualSchema = Flatten(actualSchema)

	// Compare the schemas directly, ignoring positions
	testutil.ASTEqualNoPos(t, expectedSchema, actualSchema)
}

func TestTypeFlattener_SimpleExtend(t *testing.T) {
	actual := `
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

	expected := `
		version 1

		type BaseUser {
		  id: string
		  username: string
		}

		type ExtendedUser {
		  id: string
		  username: string
		  email: string
		  age: int
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_MultipleExtends(t *testing.T) {
	// Test case: Type that extends multiple types
	actual := `
		version 1

		type BaseUser {
		  id: string
		}

		type UserWithName {
		  username: string
		}

		type ExtendedUser extends BaseUser, UserWithName {
		  email: string
		}
	`

	expected := `
		version 1

		type BaseUser {
		  id: string
		}

		type UserWithName {
		  username: string
		}

		type ExtendedUser {
		  id: string
		  username: string
		  email: string
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_NestedExtends(t *testing.T) {
	// Test case: Type that extends another type that in turn extends another type
	actual := `
		version 1

		type BaseUser {
		  id: string
		}

		type UserWithName extends BaseUser {
		  username: string
		}

		type ExtendedUser extends UserWithName {
		  email: string
		}
	`

	expected := `
		version 1

		type BaseUser {
		  id: string
		}

		type UserWithName {
		  id: string
		  username: string
		}

		type ExtendedUser {
		  id: string
		  username: string
		  email: string
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_PreserveComments(t *testing.T) {
	// Test case: Verify that comments are preserved when flattening types
	actual := `
		version 1

		type BaseUser {
		  // ID del usuario
		  id: string
		  /* Nombre de usuario */
		  username: string
		}

		type ExtendedUser extends BaseUser {
		  // Email del usuario
		  email: string
		}
	`

	expected := `
		version 1

		type BaseUser {
		  // ID del usuario
		  id: string
		  /* Nombre de usuario */
		  username: string
		}

		type ExtendedUser {
		  // ID del usuario
		  id: string
		  /* Nombre de usuario */
		  username: string
		  // Email del usuario
		  email: string
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_PreserveFieldRules(t *testing.T) {
	// Test case: Verify that validation rules are preserved when flattening types
	actual := `
		version 1

		type BaseUser {
		  id: string
		    @uuid
		  username: string
		    @minlen(3)
		    @maxlen(50)
		}

		type ExtendedUser extends BaseUser {
		  email: string
		    @email
		}
	`

	expected := `
		version 1

		type BaseUser {
		  id: string
		    @uuid
		  username: string
		    @minlen(3)
		    @maxlen(50)
		}

		type ExtendedUser {
		  id: string
		    @uuid
		  username: string
		    @minlen(3)
		    @maxlen(50)
		  email: string
		    @email
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_NoExtends(t *testing.T) {
	// Test case: Verify that types without extends don't change
	actual := `
		version 1

		type User {
		  id: string
		  username: string
		}

		type Product {
		  id: string
		  name: string
		  price: float
		}
	`

	expected := `
		version 1

		type User {
		  id: string
		  username: string
		}

		type Product {
		  id: string
		  name: string
		  price: float
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_CircularExtends(t *testing.T) {
	// A circular extends should not be flattened
	actual := `
		version 1

		type A extends B {
		  fieldA: string
		}

		type B extends A {
		  fieldB: string
		}
	`

	expected := `
		version 1

		type A {
		  fieldA: string
		}

		type B {
		  fieldB: string
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_ExtendNonExistentType(t *testing.T) {
	// Test case: Verify that the flattener correctly handles extensions to non-existent types
	// This case should not occur in practice because the semantic analyzer
	// already detects and reports extensions to non-existent types, but we test
	// that the type flattener handles this case correctly.
	actual := `
		version 1

		type User extends NonExistentType {
		  username: string
		}
	`

	expected := `
		version 1

		type User {
		  username: string
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_FieldOrder(t *testing.T) {
	// Test case: Verify that the field order is maintained correctly
	actual := `
		version 1

		type A {
		  fieldA1: string
		  fieldA2: string
		}

		type B {
		  fieldB1: string
		  fieldB2: string
		}

		type C extends A, B {
		  fieldC1: string
		  fieldC2: string
		}
	`

	expected := `
		version 1

		type A {
		  fieldA1: string
		  fieldA2: string
		}

		type B {
		  fieldB1: string
		  fieldB2: string
		}

		type C {
		  fieldA1: string
		  fieldA2: string
		  fieldB1: string
		  fieldB2: string
		  fieldC1: string
		  fieldC2: string
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_ComplexNestedExtends(t *testing.T) {
	// Test case: Verify that the flattener correctly handles complex nested extensions
	actual := `
		version 1

		type Base {
		  id: string
		}

		type Level1A extends Base {
		  fieldA: string
		}

		type Level1B extends Base {
		  fieldB: string
		}

		type Level2 extends Level1A, Level1B {
		  fieldC: string
		}

		type Level3 extends Level2 {
		  fieldD: string
		}
	`

	expected := `
		version 1

		type Base {
		  id: string
		}

		type Level1A {
		  id: string
		  fieldA: string
		}

		type Level1B {
		  id: string
		  fieldB: string
		}

		type Level2 {
		  id: string
		  fieldA: string
			id: string
		  fieldB: string
		  fieldC: string
		}

		type Level3 {
		  id: string
		  fieldA: string
			id: string
		  fieldB: string
		  fieldC: string
		  fieldD: string
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_InlineObjectFields(t *testing.T) {
	// Test case: Verify that the flattener correctly handles fields with inline objects
	actual := `
		version 1

		type BaseAddress {
		  street: string
		  city: string
		}

		type User extends BaseAddress {
		  name: string
		  contact: {
		    email: string
		    phone: string
		  }
		}
	`

	expected := `
		version 1

		type BaseAddress {
		  street: string
		  city: string
		}

		type User {
		  street: string
		  city: string
		  name: string
		  contact: {
		    email: string
		    phone: string
		  }
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_ArrayFields(t *testing.T) {
	// Test case: Verify that the flattener correctly handles array type fields
	actual := `
		version 1

		type BaseCollection {
		  items: string[]
		}

		type ExtendedCollection extends BaseCollection {
		  moreItems: int[]
		  complexItems: {
		    name: string
		    value: float
		  }[]
		}
	`

	expected := `
		version 1

		type BaseCollection {
		  items: string[]
		}

		type ExtendedCollection {
		  items: string[]
		  moreItems: int[]
		  complexItems: {
		    name: string
		    value: float
		  }[]
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_OptionalFields(t *testing.T) {
	// Test case: Verify that the flattener correctly handles optional fields
	actual := `
		version 1

		type BaseUser {
		  id: string
		  email?: string
		}

		type ExtendedUser extends BaseUser {
		  name: string
		  phone?: string
		}
	`

	expected := `
		version 1

		type BaseUser {
		  id: string
		  email?: string
		}

		type ExtendedUser {
		  id: string
		  email?: string
		  name: string
		  phone?: string
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_ComplexNestedObjects(t *testing.T) {
	// Test case: Verify that the flattener correctly handles deeply nested object structures
	actual := `
		version 1

		type BaseConfig {
		  settings: {
		    display: {
		      theme: string
		      fontSize: int
		      colors: {
		        primary: string
		        secondary: string
		        accent: string
		      }
		    }
		    notifications: {
		      enabled: boolean
		      frequency: string
		    }
		  }
		}

		type ExtendedConfig extends BaseConfig {
		  settings: {
		    advanced: {
		      debug: boolean
		      experimental: {
		        features: boolean
		      }
		    }
		  }
		  customSettings: {
		    userDefined: boolean
		  }
		}
	`

	expected := `
		version 1

		type BaseConfig {
		  settings: {
		    display: {
		      theme: string
		      fontSize: int
		      colors: {
		        primary: string
		        secondary: string
		        accent: string
		      }
		    }
		    notifications: {
		      enabled: boolean
		      frequency: string
		    }
		  }
		}

		type ExtendedConfig {
		  settings: {
		    display: {
		      theme: string
		      fontSize: int
		      colors: {
		        primary: string
		        secondary: string
		        accent: string
		      }
		    }
		    notifications: {
		      enabled: boolean
		      frequency: string
		    }
		  }
		  settings: {
		    advanced: {
		      debug: boolean
		      experimental: {
		        features: boolean
		      }
		    }
		  }
		  customSettings: {
		    userDefined: boolean
		  }
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_FieldNameConflicts(t *testing.T) {
	// Test case: Verify that the flattener correctly handles field name conflicts in multiple inheritance
	actual := `
		version 1

		type Base1 {
		  id: string
		  name: string
		  common: int
		}

		type Base2 {
		  id: int
		  email: string
		  common: string
		}

		type Combined extends Base1, Base2 {
		  id: boolean
		  custom: float
		}
	`

	expected := `
		version 1

		type Base1 {
		  id: string
		  name: string
		  common: int
		}

		type Base2 {
		  id: int
		  email: string
		  common: string
		}

		type Combined {
		  id: string
		  name: string
		  common: int
		  email: string
		  id: boolean
		  custom: float
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_WithDocstrings(t *testing.T) {
	// Test case: Verify that the flattener preserves docstrings
	actual := `
		version 1

		"""Base entity with common fields"""
		type BaseEntity {
		  // Unique identifier
		  id: string
		  // Creation timestamp
		  createdAt: datetime
		}

		"""User entity with authentication details"""
		type User extends BaseEntity {
		  // User's email address
		  email: string
		  // Hashed password
		  password: string
		}
	`

	expected := `
		version 1

		"""Base entity with common fields"""
		type BaseEntity {
		  // Unique identifier
		  id: string
		  // Creation timestamp
		  createdAt: datetime
		}

		"""User entity with authentication details"""
		type User {
		  // Unique identifier
		  id: string
		  // Creation timestamp
		  createdAt: datetime
		  // User's email address
		  email: string
		  // Hashed password
		  password: string
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_MultiDimensionalArrays(t *testing.T) {
	// Test case: Verify that the flattener correctly handles multi-dimensional arrays
	actual := `
		version 1

		type BaseMatrix {
		  values: int[][]
		  labels: string[]
		}

		type ExtendedMatrix extends BaseMatrix {
		  tensors: float[][][]
		  metadata: {
		    name: string
		    tags: string[]
		  }[]
		}
	`

	expected := `
		version 1

		type BaseMatrix {
		  values: int[][]
		  labels: string[]
		}

		type ExtendedMatrix {
		  values: int[][]
		  labels: string[]
		  tensors: float[][][]
		  metadata: {
		    name: string
		    tags: string[]
		  }[]
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_ComplexValidationRules(t *testing.T) {
	// Test case: Verify that the flattener correctly handles complex validation rules
	actual := `
		version 1

		type BaseValidation {
		  id: string
		    @uuid
		    @required
		  score: float
		    @min(0.0)
		    @max(100.0)
		    @required(error: "Score is required")
		}

		type ExtendedValidation extends BaseValidation {
		  email: string
		    @email
		    @required
		  tags: string[]
		    @minlen(1, error: "At least one tag is required")
		    @maxlen(10)
		  config: {
		    option1: boolean
		      @equals(true)
		    option2: string
		      @enum(["value1", "value2", "value3"])
		  }
		}
	`

	expected := `
		version 1

		type BaseValidation {
		  id: string
		    @uuid
		    @required
		  score: float
		    @min(0.0)
		    @max(100.0)
		    @required(error: "Score is required")
		}

		type ExtendedValidation {
		  id: string
		    @uuid
		    @required
		  score: float
		    @min(0.0)
		    @max(100.0)
		    @required(error: "Score is required")
		  email: string
		    @email
		    @required
		  tags: string[]
		    @minlen(1, error: "At least one tag is required")
		    @maxlen(10)
		  config: {
		    option1: boolean
		      @equals(true)
		    option2: string
		      @enum(["value1", "value2", "value3"])
		  }
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_CustomTypeFields(t *testing.T) {
	// Test case: Verify that the flattener correctly handles fields with custom types
	actual := `
		version 1

		type Address {
		  street: string
		  city: string
		  country: string
		}

		type Contact {
		  email: string
		  phone: string
		}

		type BaseUser {
		  id: string
		  address: Address
		}

		type ExtendedUser extends BaseUser {
		  name: string
		  contact: Contact
		  alternateAddresses: Address[]
		}
	`

	expected := `
		version 1

		type Address {
		  street: string
		  city: string
		  country: string
		}

		type Contact {
		  email: string
		  phone: string
		}

		type BaseUser {
		  id: string
		  address: Address
		}

		type ExtendedUser {
		  id: string
		  address: Address
		  name: string
		  contact: Contact
		  alternateAddresses: Address[]
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_ComplexCircularExtends(t *testing.T) {
	// Test case: Verify that the flattener correctly handles complex circular dependencies
	// When we have circular dependencies, the fields from extended types are not included
	// but the type's own fields are preserved
	actual := `
		version 1

		type A extends B {
		  fieldA: string
		}

		type B extends C {
		  fieldB: string
		}

		type C extends A {
		  fieldC: string
		}

		// D is outside the circle
		type D extends A {
		  fieldD: string
		}
	`

	expected := `
		version 1

		type A {
		  fieldA: string
		}

		type B {
		  fieldB: string
		}

		type C {
		  fieldC: string
		}

		// D is outside the circle
		type D {
			fieldA: string
		  fieldD: string
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_MixedInlineAndCustomTypes(t *testing.T) {
	// Test case: Verify that the flattener correctly handles a mix of inline objects and custom types
	actual := `
		version 1

		type Metadata {
		  createdAt: datetime
		  updatedAt: datetime
		}

		type BaseEntity {
		  id: string
		  metadata: Metadata
		  config: {
		    enabled: boolean
		    visible: boolean
		  }
		}

		type Product extends BaseEntity {
		  name: string
		  price: float
		  details: {
		    description: string
		    specifications: {
		      weight: float
		      dimensions: {
		        width: float
		        height: float
		        depth: float
		      }
		    }
		  }
		  categories: string[]
		}
	`

	expected := `
		version 1

		type Metadata {
		  createdAt: datetime
		  updatedAt: datetime
		}

		type BaseEntity {
		  id: string
		  metadata: Metadata
		  config: {
		    enabled: boolean
		    visible: boolean
		  }
		}

		type Product {
		  id: string
		  metadata: Metadata
		  config: {
		    enabled: boolean
		    visible: boolean
		  }
		  name: string
		  price: float
		  details: {
		    description: string
		    specifications: {
		      weight: float
		      dimensions: {
		        width: float
		        height: float
		        depth: float
		      }
		    }
		  }
		  categories: string[]
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_NestedOptionalFields(t *testing.T) {
	// Test case: Verify that the flattener correctly handles nested optional fields
	actual := `
		version 1

		type BaseUser {
		  id: string
		  profile?: {
		    name?: string
		    bio?: string
		  }
		}

		type ExtendedUser extends BaseUser {
		  settings?: {
		    theme?: string
		    notifications?: {
		      email?: boolean
		      push?: boolean
		    }
		  }
		  friends?: string[]
		}
	`

	expected := `
		version 1

		type BaseUser {
		  id: string
		  profile?: {
		    name?: string
		    bio?: string
		  }
		}

		type ExtendedUser {
		  id: string
		  profile?: {
		    name?: string
		    bio?: string
		  }
		  settings?: {
		    theme?: string
		    notifications?: {
		      email?: boolean
		      push?: boolean
		    }
		  }
		  friends?: string[]
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_EmptySchema(t *testing.T) {
	// Test case: Verify that the flattener correctly handles an empty schema
	actual := `
		version 1
	`

	expected := `
		version 1
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_MultipleVersions(t *testing.T) {
	// Test case: Verify that the flattener correctly handles schemas with multiple version declarations
	actual := `
		version 1

		type BaseUser {
		  id: string
		}

		version 1

		type ExtendedUser extends BaseUser {
		  name: string
		}
	`

	expected := `
		version 1

		type BaseUser {
		  id: string
		}

		version 1

		type ExtendedUser {
		  id: string
		  name: string
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_ExtendWithComments(t *testing.T) {
	// Test case: Verify that the flattener correctly handles types with comments
	actual := `
		version 1

		type BaseUser {
		  id: string
		}

		// This type extends BaseUser
		type ExtendedUser extends BaseUser {
		  // User's name
		  name: string
		}
	`

	expected := `
		version 1

		type BaseUser {
		  id: string
		}

		// This type extends BaseUser
		type ExtendedUser {
		  id: string
		  // User's name
		  name: string
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_ExtendWithRulesAndComments(t *testing.T) {
	// Test case: Verify that the flattener correctly handles types with rules and comments
	actual := `
		version 1

		type BaseValidation {
		  // ID field with validation
		  id: string
		    @uuid
		    // Required field
		    @required
		}

		type ExtendedValidation extends BaseValidation {
		  // Email field with validation
		  email: string
		    /* Must be a valid email */
		    @email
		    // Required field
		    @required
		}
	`

	expected := `
		version 1

		type BaseValidation {
		  // ID field with validation
		  id: string
		    @uuid
		    // Required field
		    @required
		}

		type ExtendedValidation {
		  // ID field with validation
		  id: string
		    @uuid
		    // Required field
		    @required
		  // Email field with validation
		  email: string
		    /* Must be a valid email */
		    @email
		    // Required field
		    @required
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_ExtendWithMultipleInheritanceLevels(t *testing.T) {
	// Test case: Verify that the flattener correctly handles multiple levels of inheritance with complex fields
	actual := `
		version 1

		type Entity {
		  id: string
		  createdAt: datetime
		}

		type Person extends Entity {
		  name: string
		  age: int
		}

		type Employee extends Person {
		  employeeId: string
		  department: string
		}

		type Manager extends Employee {
		  managedDepartments: string[]
		  reports: Employee[]
		}
	`

	expected := `
		version 1

		type Entity {
		  id: string
		  createdAt: datetime
		}

		type Person {
		  id: string
		  createdAt: datetime
		  name: string
		  age: int
		}

		type Employee {
		  id: string
		  createdAt: datetime
		  name: string
		  age: int
		  employeeId: string
		  department: string
		}

		type Manager {
		  id: string
		  createdAt: datetime
		  name: string
		  age: int
		  employeeId: string
		  department: string
		  managedDepartments: string[]
		  reports: Employee[]
		}
	`

	assertFlattenedSchema(t, expected, actual)
}

func TestTypeFlattener_ExtendWithDuplicateFieldNames(t *testing.T) {
	// Test case: Verify that the flattener correctly handles duplicate field names
	// The first occurrence of a field should be used
	actual := `
		version 1

		type Base {
		  id: string
		  name: string
		  description: string
		}

		type Extended extends Base {
		  // This field should override the one from Base
		  name: int
		  // This is a new field
		  code: string
		}
	`

	expected := `
		version 1

		type Base {
		  id: string
		  name: string
		  description: string
		}

		type Extended {
		  id: string
		  name: string
		  description: string
		  // This field should override the one from Base
		  name: int
		  // This is a new field
		  code: string
		}
	`

	assertFlattenedSchema(t, expected, actual)
}
