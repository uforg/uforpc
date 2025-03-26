# UFO-RPC DSL (URPC) Specification v1.0

## 1. Overview

Domain-Specific Language (DSL) for defining RPC services with type validation
and business rules. Transpiles to UFO-RPC-compatible JSON Schema.

## 2. Basic Syntax

### 2.1 General Structure

```urpc
version: <number>

// <comment>

"""
<Type documentation>
"""
type <TypeName> {
  <field>[?]: <Type>
    [@<validationRule>(<params>, [error: <"message">])...]
}

"""
<Procedure documentation>
"""
proc <ProcedureName> {
  input {
    <field>[?]: <Type>
      [@<validationRule>(<params>, [error: <"message">])...]
  }
  
  output {
    <field>[?]: <Type>
      [@<validationRule>(<params>, [error: <"message">])...]
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

// Type reference
TypeName
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
    baz?: string[]
  }
}
```

## 4. Validation Rules

### 4.1 General Syntax

```urpc
@<ruleName>([value][, error: "message"])
```

### 4.2 Type-Specific Rules

#### String (`@rule`)

| Rule        | Parameters   | Example                      |
| ----------- | ------------ | ---------------------------- |
| `minLen`    | integer      | `@minLen(5)`                 |
| `maxLen`    | integer      | `@maxLen(100, error: "...")` |
| `enum`      | [string,...] | `@enum(["yes", "no"])`       |
| `email`     | -            | `@email`                     |
| `uuid`      | -            | `@uuid(error: "Invalid ID")` |
| `iso8601`   | -            | `@iso8601`                   |
| `lowercase` | -            | `@lowercase`                 |
| `uppercase` | -            | `@uppercase`                 |

#### Int (`@rule`)

| Rule   | Parameters    | Example            |
| ------ | ------------- | ------------------ |
| `min`  | integer       | `@min(18)`         |
| `max`  | integer       | `@max(100)`        |
| `enum` | [integer,...] | `@enum([1, 2, 3])` |

#### Float (`@rule`)

| Rule   | Parameters   | Example             |
| ------ | ------------ | ------------------- |
| `min`  | number       | `@min(0.5)`         |
| `max`  | number       | `@max(999.99)`      |
| `enum` | [number,...] | `@enum([1.1, 2.0])` |

#### Array (`@rule`)

| Rule       | Parameters | Example          |
| ---------- | ---------- | ---------------- |
| `minItems` | integer    | `@minItems(1)`   |
| `maxItems` | integer    | `@maxItems(100)` |

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
  <key>: string|number|boolean
  
  // Examples
  cache: true
  ttl: 300
  auth: "required"
}
```

## 7. Complete Example

```urpc
version: 1

"""
Represents a product in the catalog
"""
type Product {
  id: string
    @uuid
    @minLen(36)
  
  name: string
    @minLen(3)
    @maxLen(100)
  
  price: float
    @min(0.01)
  
  tags?: string[]
    @maxItems(5)
}

"""
Creates a new product in the system and returns the product id.
"""
proc CreateProduct {
  input {
    product: Product
    priority: int
      @enum([1, 2, 3])
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
