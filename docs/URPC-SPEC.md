# UFO-RPC DSL (URPC) Specification v1.0

## 1. Overview

The UFO-RPC DSL (URPC) is a domain-specific language designed to define RPC
services with strong type validation and business rules. It provides a
declarative syntax for defining data structures, validation rules, and
procedures that transpile into UFO-RPC-compatible JSON Schema.

The primary goal of URPC is to offer an intuitive, human-readable format that
ensures the best possible developer experience (DX) while maintaining strict
data integrity.

## 2. URPC Syntax

This is the syntax for the URPC DSL.

```urpc
version: <number>

// <comment>

/*
  <multiline comment>
*/

import "<schema_path>"

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
type <CustomTypeName> [extends <OtherCustomTypeName>, <OtherCustomTypeName>, ...] {
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
      [@<validationRule>(<param>, [error: <"message">])]
  }
  
  meta {
    <key>: <value>
  }
}
```

## 3. Importing Schemas

To achieve modularity, you can import external schema files with a simple
syntax.

All paths should be relative to the file that is importing the other schema.

```urpc
import "<schema_path>"
```

### 3.1 Conflict resolution

If there are conflicting custom validaion rules, custom types or procedures, the
transpiler will throw an error. Make sure to use unique names across schemas to
avoid issues.

## 4. Types

Types are the building blocks of your API. They define the structure of the data
that can be exchanged between the client and the server.

### 4.1 Primitive Types

Primitive types are the types that are built-in into the URPC DSL.

| DSL        | JSON Type | Description                           |
| ---------- | --------- | ------------------------------------- |
| `string`   | string    | UTF-8 text string                     |
| `int`      | integer   | 64-bit integer                        |
| `float`    | number    | Floating point number                 |
| `boolean`  | boolean   | Either true or false                  |
| `datetime` | string    | Date and time value (ISO 8601 format) |

### 4.2 Composite Types

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

### 4.3 Custom Types

You can define custom types additional of the primitive types provided by the
transpiler that you can use in the input and output of your procedures.

```urpc
"""
<Type documentation>
"""
type <CustomTypeName> [extends <OtherCustomTypeName>, <OtherCustomTypeName>, ...] {
  <field>[?]: <Type>
    [@<validationRule>(<param>, [error: <"message">])]
}
```

#### 4.3.1 Custom type documentation

You can add documentation to your custom types to help the developer understand
how to use them, they can include Markdown syntax that will be rendered in the
generated documentation.

#### 4.3.2 Type inheritance

Types can extend other types, inheriting their fields and validation rules.

Fields or validation rules defined in the parent type can't be overridden by the
child type.

```urpc
type ExtendedType extends BaseType {
  // Additional fields and rules
}
```

#### 4.3.3 Optional fields

All fields of a type are required by default. To make a field optional, use the
`?` suffix.

```urpc
// Optional field
field?: Type
```

## 5. Defining Procedures

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
      [@<validationRule>(<param>, [error: <"message">])]
  }
  
  meta {
    <key>: <value>
  }
}
```

### 5.1 Procedure documentation

You can add documentation to your procedures to help the developer understand
how to use them, they can include Markdown syntax that will be rendered in the
generated documentation.

### 5.2 Procedure input/output

The input and output of a procedure is defined in the same way as the types.

### 5.3 Procedure metadata

The metadata of a procedure is a map of key-value pairs that can be used to
provide additional information about the procedure.

This information will be available in the generated code and can be used for any
purpose you want.

There is no built-in metadata, it's completely up to you to define it.

You can only define values of the following types:

- string
- int
- float
- boolean

```urpc
meta {
  // Allowed values
  <key>: string|int|float|boolean
  
  // Examples
  cache: true
  ttl: 300
  auth: "required"
}
```

## 6. Validation Rules

Validation rules can be applied to fields of a type to validate the value of the
field.

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

#### 6.1.4 Boolean built-in rules

| Rule     | Parameters | Example         |
| -------- | ---------- | --------------- |
| `equals` | boolean    | `@equals(true)` |

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
  param: <Type>              // Allowed: string|int|float|boolean or array of only these types
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

### 6.3 Validation rules and type inheritance

- Any type that extends another type will inherit the validation rules of the
  parent type.
- If a custom type is used as a field type in another type, the validation rules
  of the referenced type will be applied to the field.

## 7. Complete Example

```urpc
version: 1

import "path/to/other_schema.urpc"

// Custom Rules
rule @regex {
  for: string
  param: string
  error: "Invalid format"
}

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
Represents a product in the catalog
"""
type Product {
  id: string
    @uuid

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
Represents a product with its reviews
"""
type ReviewedProduct extends Product {
  reviews: Review[]
}

"""
Creates a new product in the system and returns the product id.
"""
proc CreateProduct {
  input {
    product: Product
  }

  output {
    success: boolean
    productId: string
      @uuid
  }

  meta {
    requiresAuth: true
    maxRetries: 3
  }
}

"""
Get a product by id, it returns a product with its reviews.
"""
proc GetProduct {
  input {
    productId: string
      @uuid
  }

  output {
    product: ReviewedProduct
  }
}
```

## 8. Known Limitations

1. Keywords can't be used as identifiers
2. Custom validators require external implementation
