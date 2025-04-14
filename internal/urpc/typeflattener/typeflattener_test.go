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
