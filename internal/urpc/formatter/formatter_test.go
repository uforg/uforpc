package formatter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormat_Version(t *testing.T) {
	input := `version 1`
	expected := "version 1\n"

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_Imports(t *testing.T) {
	input := `
version 1 import "path/to/schema.urpc" import   "another/schema.urpc"
`

	expected := `version 1

import "path/to/schema.urpc"
import "another/schema.urpc"
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_EmptySchema(t *testing.T) {
	input := ``
	expected := "\n"

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_VersionOnly(t *testing.T) {
	input := `
  version  42 
`
	expected := "version 42\n"

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_CustomRule(t *testing.T) {
	input := `
version 1

rule @customRule {
for: string
param: string
error: "This is an error message"
}
`

	expected := `version 1

rule @customRule {
  for: string
  param: string
  error: "This is an error message"
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_CustomRuleWithArrayParam(t *testing.T) {
	input := `
version 1

rule @rangeRule {
for: int
param: int[]
error: "Value out of range"
}
`

	expected := `version 1

rule @rangeRule {
  for: int
  param: int[]
  error: "Value out of range"
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_CustomRuleNoParam(t *testing.T) {
	input := `
version 1

rule @noParamRule {
for: string
error: "Error without parameter"
}
`

	expected := `version 1

rule @noParamRule {
  for: string
  error: "Error without parameter"
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_CustomRuleNoError(t *testing.T) {
	input := `
version 1

rule @noErrorRule {
for: boolean
param: string
}
`

	expected := `version 1

rule @noErrorRule {
  for: boolean
  param: string
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_CustomType(t *testing.T) {
	input := `
version 1

type User {
id: string
name: string
age: int
}
`

	expected := `version 1

type User {
  id: string
  name: string
  age: int
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_CustomTypeWithExtends(t *testing.T) {
	input := `
version 1

type BaseUser {
id: string
name: string
}

type ExtendedUser extends BaseUser {
email: string
phone?: string
}
`

	expected := `version 1

type BaseUser {
  id: string
  name: string
}

type ExtendedUser extends BaseUser {
  email: string
  phone?: string
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_CustomTypeWithMultipleExtends(t *testing.T) {
	input := `
version 1

type BaseUser {
id: string
}

type WithName {
name: string
}

type WithEmail {
email: string
}

type CompleteUser extends BaseUser, 
WithName, 			WithEmail {
age: int
}
`

	expected := `version 1

type BaseUser {
  id: string
}

type WithName {
  name: string
}

type WithEmail {
  email: string
}

type CompleteUser extends BaseUser, WithName, WithEmail {
  age: int
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_EmptyType(t *testing.T) {
	input := `
version 1

type EmptyType {
dummy: string
}
`

	expected := `version 1

type EmptyType {
  dummy: string
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_TypeWithRules(t *testing.T) {
	input := `
version 1

rule @minValue {
for: int
param: int
error: "Value too small"
}

type User {
id: string
@minlen(3)
@maxlen(50)
age: int
@minValue(18, error: "Must be at least 18 years old")
}
`

	expected := `version 1

rule @minValue {
  for: int
  param: int
  error: "Value too small"
}

type User {
  id: string
    @minlen(3)
    @maxlen(50)
  age: int
    @minValue(18, error: "Must be at least 18 years old")
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_TypeWithParameterlessRule(t *testing.T) {
	input := `
version 1

rule @noParams {
for: string
}

type Contact {
email: string
@noParams
}
`

	expected := `version 1

rule @noParams {
  for: string
}

type Contact {
  email: string
    @noParams
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_TypeWithArrays(t *testing.T) {
	input := `
version 1

type Collection {
stringArray: string[]
nestedArray: int[][]
}
`

	expected := `version 1

type Collection {
  stringArray: string[]
  nestedArray: int[][]
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_TypeWithDeepNestedArrays(t *testing.T) {
	input := `
version 1

type DeepNested {
matrix3d: int[][][]
complex: string[][][][][]
}
`

	expected := `version 1

type DeepNested {
  matrix3d: int[][][]
  complex: string[][][][][]
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_TypeWithInlineObject(t *testing.T) {
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

	expected := `version 1

type User {
  id: string
  address: {
    street: string
    city: string
    zipCode: string
  }
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_TypeWithEmptyInlineObject(t *testing.T) {
	input := `
version 1

type Config {
options: {
  dummy: string
}
}
`

	expected := `version 1

type Config {
  options: {
    dummy: string
  }
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_TypeWithNestedInlineObjects(t *testing.T) {
	input := `
version 1

type Nested {
user: {  name: string
address: {
	street: string
	location: {lat: float
		lng: float
	}
  }}                            }
`

	expected := `version 1

type Nested {
  user: {
    name: string
    address: {
      street: string
      location: {
        lat: float
        lng: float
      }
    }
  }
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_TypeWithInlineObjectAndRules(t *testing.T) {
	input := `
version 1

type ComplexType {
settings: {
  theme: string
  @enum(["light", "dark"])
  notifications: boolean
  @equals(true)
}
}
`

	expected := `version 1

type ComplexType {
  settings: {
    theme: string
      @enum(["light", "dark"])
    notifications: boolean
      @equals(true)
  }
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_ProcedureWithInputOutput(t *testing.T) {
	input := `
version 1

proc GetUser {
input {
  userId: string
}

output {
  name: string
  age: int
}
}
`

	expected := `version 1

proc GetUser {
  input {
    userId: string
  }

  output {
    name: string
    age: int
  }
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_ProcedureEmpty(t *testing.T) {
	input := `
version 1

proc EmptyProc {
input {
  dummy: string
}
}
`

	expected := `version 1

proc EmptyProc {
  input {
    dummy: string
  }
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_ProcedureWithOnlyInput(t *testing.T) {
	input := `
version 1

proc InputOnly {
input {
  command: string
}
}
`

	expected := `version 1

proc InputOnly {
  input {
    command: string
  }
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_ProcedureWithOnlyOutput(t *testing.T) {
	input := `
version 1

proc OutputOnly {
output {
  result: boolean
}
}
`

	expected := `version 1

proc OutputOnly {
  output {
    result: boolean
  }
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_ProcedureWithMeta(t *testing.T) {
	input := `
version 1

proc GetUser {
input {
  userId: string
}

output {
  name: string
  age: int
}

meta {
  requiresAuth: true
  maxRetries: 3
  timeout: 60
  description: "Gets a user by ID"
}
}
`

	expected := `version 1

proc GetUser {
  input {
    userId: string
  }

  output {
    name: string
    age: int
  }

  meta {
    requiresAuth: true
    maxRetries: 3
    timeout: 60
    description: "Gets a user by ID"
  }
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_ProcedureWithOnlyMeta(t *testing.T) {
	input := `
version 1

proc MetaOnly {
meta {
  internal: true
  debug: false
}
}
`

	expected := `version 1

proc MetaOnly {
  meta {
    internal: true
    debug: false
  }
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_RuleWithDocstring(t *testing.T) {
	input := `
version 1

"""
This is a custom validation rule for emails
""" rule @email {
  for: string
  error: "Invalid email format"
}
`

	expected := `version 1

"""
This is a custom validation rule for emails
"""
rule @email {
  for: string
  error: "Invalid email format"
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_MultilineDocstring(t *testing.T) {
	input := `
version 1

"""
This is a multiline
docstring with several
lines of documentation
that should be preserved
as is.
"""


type Documented {
  id: string
}
`

	expected := `version 1

"""
This is a multiline
docstring with several
lines of documentation
that should be preserved
as is.
"""
type Documented {
  id: string
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_DocstringConversion(t *testing.T) {
	input := `
version 1

"""
This is a docstring that will be formatted properly
"""
rule @withDocstring {
  for: string
}
`

	expected := `version 1

"""
This is a docstring that will be formatted properly
"""
rule @withDocstring {
  for: string
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_EnumArrayRule(t *testing.T) {
	input := `
version 1

type Status {
  state: string
    @enum(["pending", "active", "completed"])
}
`

	expected := `version 1

type Status {
  state: string
    @enum(["pending", "active", "completed"])
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_BooleanArrayRule(t *testing.T) {
	input := `
version 1

type Flags {
  options: boolean[]
    @enum([true, false])
}
`

	expected := `version 1

type Flags {
  options: boolean[]
    @enum([true, false])
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_IntArrayRule(t *testing.T) {
	input := `
version 1

type Ratings {
  stars: int
    @enum([1, 2, 3, 4, 5])
}
`

	expected := `version 1

type Ratings {
  stars: int
    @enum([1, 2, 3, 4, 5])
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_CompleteSchema(t *testing.T) {
	input := `
version 1

import "common.urpc"

"""
Custom validation rule for URL format
"""
rule @url {
  for: string
  param: string
  error: "Invalid URL format"
}

"""
Base user type with common properties
"""
type BaseUser {
  id: string
    @minlen(10)
  createdAt: datetime
}

"""
Extended user with additional properties
"""
type User extends BaseUser {
  username: string
    @minlen(3)
    @maxlen(50)
  email: string
    @contains("@")
  address?: {
    street: string
    city: string
    zipCode: string
      @minlen(5)
  }
  status: string
    @enum(["active", "inactive", "suspended"])
  tags: string[]
    @minlen(1)
  metadata: {
    lastLogin?: datetime
    preferences: {
      theme: string
        @enum(["light", "dark"])
      notifications: boolean
    }
  }
}

"""
Creates a new user in the system
"""
proc CreateUser {
  input {
    username: string
      @minlen(3)
    email: string
      @contains("@")
    password: string
      @minlen(8)
  }

  output {
    userId: string
    success: boolean
    errors: string[]
  }

  meta {
    requiresAuth: false
    rateLimit: 10
    description: "Creates a new user account"
  }
}

"""
Gets a user by ID
"""
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
  }
}
`

	expected := `version 1

import "common.urpc"

"""
Custom validation rule for URL format
"""
rule @url {
  for: string
  param: string
  error: "Invalid URL format"
}

"""
Base user type with common properties
"""
type BaseUser {
  id: string
    @minlen(10)
  createdAt: datetime
}

"""
Extended user with additional properties
"""
type User extends BaseUser {
  username: string
    @minlen(3)
    @maxlen(50)
  email: string
    @contains("@")
  address?: {
    street: string
    city: string
    zipCode: string
      @minlen(5)
  }
  status: string
    @enum(["active", "inactive", "suspended"])
  tags: string[]
    @minlen(1)
  metadata: {
    lastLogin?: datetime
    preferences: {
      theme: string
        @enum(["light", "dark"])
      notifications: boolean
    }
  }
}

"""
Creates a new user in the system
"""
proc CreateUser {
  input {
    username: string
      @minlen(3)
    email: string
      @contains("@")
    password: string
      @minlen(8)
  }

  output {
    userId: string
    success: boolean
    errors: string[]
  }

  meta {
    requiresAuth: false
    rateLimit: 10
    description: "Creates a new user account"
  }
}

"""
Gets a user by ID
"""
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
  }
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_MalformedButParseable(t *testing.T) {
	input := `
	version 1

 type  InconsistentSpacing{
id:string
 name :  string
age   :int
}

proc   MessyProc{
input {userId :string}
	
 output   {name:string
age:int}
}
`

	expected := `version 1

type InconsistentSpacing {
  id: string
  name: string
  age: int
}

proc MessyProc {
  input {
    userId: string
  }

  output {
    name: string
    age: int
  }
}
`

	formatted, err := Format(input)

	require.NoError(t, err)
	require.Equal(t, expected, formatted)
}

func TestFormat_InvalidInput(t *testing.T) {
	input := `invalid urpc schema`

	_, err := Format(input)

	require.Error(t, err)
}
