# UFO-RPC DSL (URPC) Specification v1.0

## 1. Overview

The UFO-RPC DSL (URPC) is a domain-specific language designed to define RPC
services with strong type validation and business rules. It provides a
declarative syntax for defining data structures, validation rules, and
procedures that transpile into UFO-RPC-compatible JSON Schema.

The primary goal of URPC is to offer an intuitive, human-readable format that
ensures the best possible developer experience (DX) while maintaining strict
data integrity.

## 2. Basic Syntax

### 2.1 General Structure

```urpc
version: <number>

// <comment>

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
type <TypeName> {
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

## 3. Data Types

### 3.1 Primitive Types

| DSL       | JSON Type | Description           |
| --------- | --------- | --------------------- |
| `string`  | string    | UTF-8 text string     |
| `int`     | integer   | 64-bit integer        |
| `float`   | number    | Floating point number |
| `boolean` | boolean   | True/False value      |

### 3.2 Composite Types

```urpc
// Array
ElementType[]  // E.g.: string[] 

// Inline object
{
  field1: Type
  field2: Type
}
```

### 3.3 Optionals

All fields are required by default. To make a field optional, use the `?`
suffix.

```urpc
// Optional field
field?: Type

// Example
input {
  foo?: int
  bar?: {
    baz?: MyCustomType[]
  }
}
```

## 4. Validation Rules

### 4.1 General Syntax

```urpc
@<ruleName>([param][, error: "message"])
```

### 4.2 Built-in Type-Specific Rules

Built-in rules are implemented in the generated code for you, so you can use
them as is without any additional implementation.

These rules are very carefully picked and designed to be implemented in any
language in the most deterministic way possible, so the built-in rules will work
the same regardless of the target language of the generated code.

If you need to implement a custom or more complex rule, you can do so by
declaring a custom rule and implementing it in your codebase.

#### String built-in rules

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

#### Int built-in rules

| Rule     | Parameters    | Example            |
| -------- | ------------- | ------------------ |
| `equals` | integer       | `@equals(1)`       |
| `min`    | integer       | `@min(0)`          |
| `max`    | integer       | `@max(100)`        |
| `enum`   | [integer,...] | `@enum([1, 2, 3])` |

#### Float built-in rules

Rules like `equals`, `enum` are not supported for float because of the nature of
how computers represent floating point numbers UFO RPC can't guarantee the
precision of the float number. You can implement custom rules for float with
your own logic to address this limitation.

| Rule  | Parameters | Example       |
| ----- | ---------- | ------------- |
| `min` | number     | `@min(0.0)`   |
| `max` | number     | `@max(100.0)` |

#### Boolean built-in rules

| Rule     | Parameters | Example         |
| -------- | ---------- | --------------- |
| `equals` | boolean    | `@equals(true)` |

#### Array built-in rules

| Rule     | Parameters | Example        |
| -------- | ---------- | -------------- |
| `minlen` | integer    | `@minlen(1)`   |
| `maxlen` | integer    | `@maxlen(100)` |

### 4.3 Custom Rules

Define reusable validation rules. Rules must be declared **before** usage.

Your custom defined rules, should be defined in the DSL and implemented in your
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

#### Example

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

#### Usage

```urpc
input {
  username: string
    @regex("^[a-z0-9_]+$", error: "Only lowercase allowed")

  age: int
    @range([18, 120])
}
```

## 5. Procedures

### 5.1 Structure

```urpc
"""
<Procedure documentation>
"""
proc <Name> {
  input {
    <field>: <Type>
  }

  output {
    <field>: <Type>
  }

  meta {
    <key>: <value>
  }
}
```

## 6. Metadata

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

## 7. Complete Example

```urpc
version: 1

// Custom Rules
rule @regex {
  for: string
  param: string
  error: "Invalid format"
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
  
  tags?: string[]
    @maxlen(5)
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
```

## 8. Known Limitations

1. No type inheritance support
2. Custom validators require external implementation
3. No multi-line comment support
