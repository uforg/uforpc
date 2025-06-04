# UFO-RPC DSL (URPC) Specification v1.0

## 1. Overview

The UFO-RPC DSL (URPC) is a domain-specific language designed to define RPC
services with strong type validation and business rules. It provides a
declarative syntax for defining data structures, validation rules, and
procedures that UFO RPC can interpret and generate code for.

The primary goal of URPC is to offer an intuitive, human-readable format that
ensures the best possible developer experience (DX) while maintaining strict
data integrity.

## 2. URPC Syntax

This is the syntax for the URPC DSL.

```urpc
version <number>

// <comment>

/*
  <multiline comment>
*/

"""
<Custom rule documentation>
"""
rule @<CustomRuleName> {
  for: <Type>
  param: <PrimitiveType> | <PrimitiveType>[]
  error: "<Default error message>"
}

"""
<Type documentation>
"""
type <CustomTypeName> {
  <field>[?]: <Type>
    [@<validationRule>(<param>, [error: <"message">])]
}

"""
<Procedure documentation>
"""
proc <ProcedureName> {
  input {
    <field>[?]: <PrimitiveType> | <CustomType>
      [@<validationRule>(<param>, [error: <"message">])]
  }

  output {
    <field>[?]: <PrimitiveType> | <CustomType>
  }

  meta {
    <key>: <value>
  }
}

"""
<Stream documentation>
"""
stream <StreamName> {
  input {
    <field>[?]: <PrimitiveType> | <CustomType>
      [@<validationRule>(<param>, [error: <"message">])]
  }

  output {
    <field>[?]: <PrimitiveType> | <CustomType>
  }

  meta {
    <key>: <value>
  }
}
```

## 3. Types

Types are the building blocks of your API. They define the structure of the data
that can be exchanged between the client and the server.

### 3.1 Primitive Types

Primitive types are the types that are built-in into the URPC DSL.

| DSL        | JSON Type | Description                           |
| ---------- | --------- | ------------------------------------- |
| `string`   | string    | UTF-8 text string                     |
| `int`      | integer   | 64-bit integer                        |
| `float`    | number    | Floating point number                 |
| `bool`     | boolean   | Either true or false                  |
| `datetime` | string    | Date and time value (ISO 8601 format) |

### 3.2 Composite Types

Composite types are types that are composed of other types. They can be used to
create more complex data structures.

```urpc
// Array
ElementType[]  // E.g.: string[]

// Inline object
{
  field1: Type
  field2: Type
}
```

### 3.3 Custom Types

You can define custom types additional of the primitive types provided by the
transpiler that you can use in the input and output of your procedures.

```urpc
"""
<Type documentation>
"""
type <CustomTypeName> {
  <field>[?]: <Type>
    [@<validationRule>(<param>, [error: <"message">])]
}
```

#### 3.3.1 Custom type documentation

You can add documentation to your custom types to help the developer understand
how to use them, they can include Markdown syntax that will be rendered in the
generated documentation.

#### 3.3.2 Type composition

To reuse fields from other types, use composition by including the type as a field:

```urpc
type BaseEntity {
  id: string
  createdAt: datetime
  updatedAt: datetime
}

type User {
  base: BaseEntity
  email: string
  name: string
}
```

#### 3.3.3 Optional fields

All fields of a type are required by default. To make a field optional, use the
`?` suffix.

```urpc
// Optional field
field?: Type
```

## 4. Defining Procedures

Procedures are the main building block of your API. They define the procedures
(AKA functions) that can be implemented on the server and called from the
client.

```urpc
"""
<Procedure documentation>
"""
proc <ProcedureName> {
  input {
    <field>[?]: <PrimitiveType> | <CustomType>
      [@<validationRule>(<param>, [error: <"message">])]
  }

  output {
    <field>[?]: <PrimitiveType> | <CustomType>
  }

  meta {
    <key>: <value>
  }
}
```

### 4.1 Procedure documentation

You can add documentation to your procedures to help the developer understand
how to use them, they can include Markdown syntax that will be rendered in the
generated documentation.

### 4.2 Procedure input/output

The input of a procedure defines the parameters that can be validated before
processing. The output defines the structure of the response data.

Validation rules can only be applied to input fields.

### 4.3 Procedure metadata

The metadata of a procedure is a map of key-value pairs that can be used to
provide additional information about the procedure.

This information will be available in the generated code and can be used for any
purpose you want.

There is no built-in metadata, it's completely up to you to define it.

You can only define values of the following types:

- string
- int
- float
- bool

```urpc
meta {
  // Allowed values
  <key>: string|int|float|bool

  // Examples
  cache: true
  ttl: 300
  auth: "required"
}
```

## 5. Defining Streams

Streams allow server-to-client real-time communication using Server-Sent Events
(SSE). They enable unidirectional data flow from the server to subscribed
clients.

```urpc
"""
<Stream documentation>
"""
stream <StreamName> {
  input {
    <field>[?]: <PrimitiveType> | <CustomType>
      [@<validationRule>(<param>, [error: <"message">])]
  }

  output {
    <field>[?]: <PrimitiveType> | <CustomType>
  }

  meta {
    <key>: <value>
  }
}
```

### 5.1 Stream documentation

You can add documentation to your streams to help developers understand their
purpose and usage. Documentation can include Markdown syntax.

### 5.2 Stream input

The input section defines the parameters required to establish a stream
subscription. These parameters determine what data the client wants to receive
and can include validation rules.

### 5.3 Stream output

The output section defines the structure of events that will be emitted through
the stream. Each event sent to the client will conform to this structure.

### 5.4 Stream metadata

Stream metadata works the same way as procedure metadata, allowing you to attach
additional information to the stream definition.

### 5.5 Example

```urpc
"""
Stream of new messages in a specific chat room
"""
stream NewMessage {
  input {
    chatId: string
      @minlen(1)
  }

  output {
    id: string
    message: string
    userId: string
    timestamp: datetime
  }

  meta {
    auth: "required"
  }
}
```

## 6. Validation Rules

Validation rules can be applied to input fields to validate the value before
processing.

```urpc
@<ruleName>([param][, error: "message"])
```

### 6.1 Built-in Type-Specific Rules

Built-in rules are implemented in the generated code for you, so you can use
them as is without any additional implementation.

These rules are very carefully picked and designed to be implemented in any
language in the most deterministic way possible, so the built-in rules will work
the same regardless of the target language of the generated code.

If you need to implement a custom or more complex rule, you can do so by
declaring a custom rule and implementing it in your codebase.

#### 6.1.1 String built-in rules

Rules like `regex`, `email`, `uuid`, `iso8601`, `json` are not supported because
they are not deterministic and UFO RPC can't guarantee the same result in any
language. You can implement custom rules for these cases and handle them with
your own logic.

| Rule        | Parameters   | Example                               |
| ----------- | ------------ | ------------------------------------- |
| `equals`    | string       | `@equals("Foo")`                      |
| `contains`  | string       | `@contains("Bar")` (Case insensitive) |
| `minlen`    | integer      | `@minlen(3)`                          |
| `maxlen`    | integer      | `@maxlen(100)`                        |
| `enum`      | [string,...] | `@enum(["Foo", "Bar"])`               |
| `lowercase` | -            | `@lowercase`                          |
| `uppercase` | -            | `@uppercase`                          |

#### 6.1.2 Int built-in rules

| Rule     | Parameters    | Example            |
| -------- | ------------- | ------------------ |
| `equals` | integer       | `@equals(1)`       |
| `min`    | integer       | `@min(0)`          |
| `max`    | integer       | `@max(100)`        |
| `enum`   | [integer,...] | `@enum([1, 2, 3])` |

#### 6.1.3 Float built-in rules

Rules like `equals`, `enum` are not supported for float because of the nature of
how computers represent floating point numbers. UFO RPC can't guarantee the
precision of float numbers. You can implement custom rules for float with your
own logic to address this limitation.

| Rule  | Parameters | Example       |
| ----- | ---------- | ------------- |
| `min` | number     | `@min(0.0)`   |
| `max` | number     | `@max(100.0)` |

#### 6.1.4 Bool built-in rules

| Rule     | Parameters | Example         |
| -------- | ---------- | --------------- |
| `equals` | bool       | `@equals(true)` |

#### 6.1.5 Array built-in rules

| Rule     | Parameters | Example        |
| -------- | ---------- | -------------- |
| `minlen` | integer    | `@minlen(1)`   |
| `maxlen` | integer    | `@maxlen(100)` |

#### 6.1.6 Datetime built-in rules

The parameter for these rules should be a string that follows the ISO 8601
format.

| Rule  | Parameters | Example                        |
| ----- | ---------- | ------------------------------ |
| `min` | datetime   | `@min("2000-01-01T00:00:00Z")` |
| `max` | datetime   | `@max("2050-12-31T23:59:59Z")` |

### 6.2 Custom Rules

If the built-in rules don't cover your needs, you can define your own custom
rules.

Your custom defined rules should be defined in the DSL and implemented in your
codebase using the helpers included in the generated code, so you can validate
in any way you want with your own logic.

```urpc
"""
<Rule documentation>
"""
rule @<RuleName> {
  for: <Type>                // Type of the field this rule can be applied to
  param: <Type>              // Allowed: string|int|float|bool or array of only these types
  error: "<Default message>" // Optional
}
```

#### 6.2.1 Custom rule documentation

You can add documentation to your custom rules to help the developer understand
how to use them, they can include Markdown syntax that will be rendered in the
generated documentation.

#### 6.2.2 Example

```urpc
""" This is a custom rule that validates if the field matches a regular expression """
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
```

#### 6.2.3 Usage

```urpc
input {
  username: string
    @regex("^[a-z0-9_]+$", error: "Only lowercase allowed")

  age: int
    @range([18, 120])
}
```

### 6.3 Validation rules with types

If a custom type is used as a field type in another type, the validation rules
of the referenced type will be applied to the field.

## 7. Documentation

### 7.1 Docstrings

Docstrings can be used in two ways: associated with specific elements (rules,
types, or procedures) or as standalone documentation.

1. Associated docstrings: These are placed immediately before a rule, type, or
   procedure definition and provide specific documentation for that element.

   ```urpc
   """
   This is documentation for MyType.
   """
   type MyType {
     // ...
   }
   ```

2. Standalone docstrings: These provide general documentation for the schema and
   are not associated with any specific element. To create a standalone
   docstring, ensure there is at least one blank line between the docstring and
   any following element.

   ```urpc
   """
   This is general documentation for the entire schema.
   It can include multiple paragraphs and Markdown formatting.
   """

   // At least one blank line here

   type MyType {
     // ...
   }
   ```

Docstrings support Markdown syntax, allowing you to format your documentation
with headings, lists, code blocks, and more.

### 7.2 External Documentation Files

For extensive documentation, you can reference external Markdown files:

```urpc
version 1

// Standalone documentation
""" ./docs/welcome.md """
""" ./docs/authentication.md """

// Associated documentation
""" ./docs/myproc.md """
proc MyProc {
  // ...
}
```

When a docstring contains only a valid path to a Markdown file, the content of
that file will be used as documentation. This approach helps maintain clean and
focused schema files while allowing for detailed documentation in separate
files.

Remember to keep external documentation files up to date with your schema
changes.

## 8. Deprecation

URPC provides a mechanism to mark rules, types, and procedures as deprecated,
indicating they should no longer be used in new code and may be removed in
future versions.

### 8.1 Basic Deprecation

To mark an element as deprecated without a specific message, use the
`deprecated` keyword before the element definition:

```urpc
deprecated rule @myRule {
  // rule definition
}

deprecated type MyType {
  // type definition
}

deprecated proc MyProc {
  // procedure definition
}

deprecated stream MyStream {
  // stream definition
}
```

### 8.2 Deprecation with Message

To provide additional information about the deprecation, include a message in
parentheses:

```urpc
deprecated("Use newRule instead")
rule @myRule {
  // rule definition
}

deprecated("Replaced by ImprovedType")
type MyType {
  // type definition
}

deprecated("This procedure will be removed in v2.0")
proc MyProc {
  // procedure definition
}

deprecated("Use NewMessageStream instead")
stream MyStream {
  // stream definition
}
```

### 8.3 Placement

The `deprecated` keyword must be placed between any docstring and the element
definition (rule, type, proc, or stream):

```urpc
"""
Documentation for MyType
"""
deprecated("Use NewType instead")
type MyType {
  // type definition
}
```

### 8.4 Effects

Deprecated elements will:

- Be displayed with special styling in the playground to discourage their use
- Generate warning comments in the output code to discourage their use
- Not change their behavior in the generated code, it's just a warning

## 9. Complete Example

```urpc
version 1

""" ./docs/welcome.md """
""" ./docs/authentication.md """

// Custom Rules
rule @regex {
  for: string
  param: string
  error: "Invalid format"
}

"""
Validates if a string is a valid UUID
"""
deprecated("This rule will be removed in v2.0")
rule @uuid {
  for: string
  error: "Invalid UUID format"
}

rule @priceRange {
  for: float
  param: float[]
  error: "Price must be between the required range"
}

"""
Base entity with common fields
"""
type BaseEntity {
  id: string
    @uuid
  createdAt: datetime
  updatedAt: datetime
}

"""
Represents a product in the catalog
"""
type Product {
  base: BaseEntity

  name: string
    @minlen(3)
    @regex("^[A-Za-z ]+$")

  price: float
    @priceRange([0.01, 9999.99])

  availabilityDate: datetime
    @min("2020-01-01T00:00:00Z")
    @max("2030-12-31T23:59:59Z")

  tags?: string[]
    @maxlen(5)
}

"""
Represents a review of a product
"""
type Review {
  rating: int
    @enum([1, 2, 3, 4, 5])

  comment: string
    @minlen(10)
    @maxlen(500)
}

"""
Creates a new product in the system and returns the product id.
"""
proc CreateProduct {
  input {
    product: Product
  }

  output {
    success: bool
    productId: string
  }

  meta {
    requiresAuth: true
    maxRetries: 3
  }
}

"""
Get a product by id with its reviews.
"""
proc GetProduct {
  input {
    productId: string
      @uuid
  }

  output {
    product: Product
    reviews: Review[]
  }
}

"""
Sends a message to a chat room
"""
proc SendMessage {
  input {
    chatId: string
    message: string
      @maxlen(1000)
  }

  output {
    messageId: string
    timestamp: datetime
  }
}

"""
Stream of new messages in a specific chat room
"""
stream NewMessage {
  input {
    chatId: string
  }

  output {
    id: string
    message: string
    userId: string
    timestamp: datetime
  }

  meta {
    auth: "required"
  }
}
```

## 10. Known Limitations

1. Keywords can't be used as identifiers
2. Custom validators require external implementation
3. Circular type dependencies are not allowed
