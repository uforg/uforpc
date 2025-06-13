# UFO-RPC DSL (URPC) Specification v1.0

## 1. Overview

The UFO-RPC DSL (URPC) is a domain-specific language designed to define RPC
services with strong typing. It provides a declarative syntax for defining data
structures and procedures that UFO RPC can interpret and generate code for.

The primary goal of URPC is to offer an intuitive, human-readable format that
ensures the best possible developer experience (DX) while maintaining type
safety.

## 2. URPC Syntax

This is the syntax for the URPC DSL.

```urpc
version <number>

// <comment>

/*
  <multiline comment>
*/

"""
<Type documentation>
"""
type <CustomTypeName> {
  <field>[?]: <Type>
}

"""
<Procedure documentation>
"""
proc <ProcedureName> {
  input {
    <field>[?]: <PrimitiveType> | <CustomType>
  }

  output {
    <field>[?]: <PrimitiveType> | <CustomType>
  }
}

"""
<Stream documentation>
"""
stream <StreamName> {
  input {
    <field>[?]: <PrimitiveType> | <CustomType>
  }

  output {
    <field>[?]: <PrimitiveType> | <CustomType>
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
  }

  output {
    <field>[?]: <PrimitiveType> | <CustomType>
  }
}
```

### 4.1 Procedure documentation

You can add documentation to your procedures to help the developer understand
how to use them, they can include Markdown syntax that will be rendered in the
generated documentation.

### 4.2 Procedure input/output

The input of a procedure defines the parameters that are sent to the server for
processing. The output defines the structure of the response data.

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
  }

  output {
    <field>[?]: <PrimitiveType> | <CustomType>
  }
}
```

### 5.1 Stream documentation

You can add documentation to your streams to help developers understand their
purpose and usage. Documentation can include Markdown syntax.

### 5.2 Stream input

The input section defines the parameters required to establish a stream
subscription. These parameters determine what data the client wants to receive.

### 5.3 Stream output

The output section defines the structure of events that will be emitted through
the stream. Each event sent to the client will conform to this structure.

### 5.5 Example

```urpc
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
}
```

## 6. Documentation

### 6.1 Docstrings

Docstrings can be used in two ways: associated with specific elements (types,
procedures, or streams) or as standalone documentation.

1. Associated docstrings: These are placed immediately before a type, procedure,
   or stream definition and provide specific documentation for that element.

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

### 6.2 External Documentation Files

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

## 7. Deprecation

URPC provides a mechanism to mark types, procedures, and streams as deprecated,
indicating they should no longer be used in new code and may be removed in
future versions.

### 7.1 Basic Deprecation

To mark an element as deprecated without a specific message, use the
`deprecated` keyword before the element definition:

```urpc
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

### 7.2 Deprecation with Message

To provide additional information about the deprecation, include a message in
parentheses:

```urpc
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

### 7.3 Placement

The `deprecated` keyword must be placed between any docstring and the element
definition (type, proc, or stream):

```urpc
"""
Documentation for MyType
"""
deprecated("Use NewType instead")
type MyType {
  // type definition
}
```

### 7.4 Effects

Deprecated elements will:

- Be displayed with special styling in the playground to discourage their use
- Generate warning comments in the output code to discourage their use
- Not change their behavior in the generated code, it's just a warning

## 8. Complete Example

```urpc
version 1

""" ./docs/welcome.md """
""" ./docs/authentication.md """

"""
Base entity with common fields
"""
type BaseEntity {
  id: string
  createdAt: datetime
  updatedAt: datetime
}

"""
Represents a product in the catalog
"""
type Product {
  base: BaseEntity
  name: string
  price: float
  availabilityDate: datetime
  tags?: string[]
}

"""
Represents a review of a product
"""
type Review {
  rating: int
  comment: string
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
}

"""
Get a product by id with its reviews.
"""
proc GetProduct {
  input {
    productId: string
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
}
```

## 9. Known Limitations

1. Keywords can't be used as identifiers
2. Complex validation logic requires implementation via input processors
3. Circular type dependencies are not allowed
