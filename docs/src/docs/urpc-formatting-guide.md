---
title: Formatting Guide
description: A guide to correctly format UFO RPC DSL (URPC)
---

This document specifies the standard formatting rules for the UFO-RPC DSL
(URPC). Consistent formatting enhances code readability, maintainability, and
collaboration. The primary goal is to produce clean, predictable, and
aesthetically pleasing URPC code.

> **⚠️ Reference Only:** All style conventions are automatically enforced by the official URPC formatter. Run it manually with `urpc fmt ./schema.urpc`, or let the built-in LSP formatter (bundled with the VS Code extension and configurable for other editors) format files on save.

## 1. General Principles

- **Encoding:** UTF-8.
- **Line Endings:** Use newline characters (`\n`).
- **Trailing Whitespace:** None.
- **Final Newline:** End non-empty files with one newline.

## 2. Indentation

- Use **2 spaces** per indentation level.
- Do not use tabs.

_Example:_

```urpc
type Example {
  field: string
}
```

## 3. Top-Level Elements

Top-level elements include `version`, `type`, `proc`, `stream`, and
standalone comments.

- **Default:** Separate each top-level element with one blank line.
- **Exceptions:**
  - **First Element:** No blank line before the very first element.
  - **Consecutive Comments:** Do not insert extra blank lines between
    consecutive standalone comments.
  - **Following a Comment:** When an element follows a standalone comment, do
    not add an extra blank line unless the source intentionally contains one.
- **Preservation:** Intentionally placed blank lines in the source (e.g. between
  comments) are respected.

_Example:_

```urpc
version 1

// A standalone comment
// Another standalone comment
type TypeA {
  field: string
}

type TypeB {
  field: int
}
```

## 4. Fields and Blocks

### 4.1 Fields in a Type

This section applies to fields in a type block, as well as fields in a
procedure's input, output, or inline object, and stream's input, output.

- Each field is placed on its own line.
- **Field Separation:** For simple fields without complex formatting, fields may
  be placed consecutively without blank lines. When a field has a docstring,
  separate it from the preceding field with one blank line for readability.

_Recommended for simple fields:_

Do not use blank lines between fields.

```urpc
address: {
  street: string
  city: string
  zip: string
}
```

_Recommended for fields with docstrings:_

Use one blank line to separate fields that have docstrings.

```urpc
type Product {
  """ The unique identifier. """
  id: string

  """ The product's name. """
  name: string
}
```

### 4.2 Blocks (Type, Input/Output)

- Opening braces (`{`) are on the same line as the declaration header (preceded
  by one space).
- Contents inside non-empty blocks always start on a new, indented line.
- The closing brace (`}`) is placed on its own line, aligned with the opening
  line.
- In procedure and stream bodies, separate the `input`, and `output` blocks with
  one blank line.

## 5. Spacing

- **Colons (`:`):** No space before; one space after (e.g. `field: string`).
- **Commas (`,`):** No space before; one space after.
- **Braces (`{` and `}`):** One space before `{` in declarations; inside blocks,
  use newlines and proper indentation.
- **Brackets (`[]`):** No spaces for array types (e.g. `string[]`); no extra
  interior spacing.
- **Parentheses (`()`):** No extra spaces inside the parentheses.
- **Optional Marker (`?`):** Immediately follows the field name (e.g.
  `email?: string`).

## 6. Comments

Comment content is preserved exactly (including internal whitespace).

- **Standalone Comments:** Use `//` or `/* … */` on their own lines; indent to
  the current block.
- **End-of-Line (EOL) Comments:** Can use either `//` or block style (`/* … */`)
  following code on the same line, with at least one space separating them.

_Example:_

```urpc
version 1 // EOL comment

type Example {
  field: string // Inline comment for field
}
```

## 7. Docstrings

- Place docstrings immediately above the `type`, `proc`, `stream`, or `field` they
  document.
- They are enclosed in triple quotes (`"""`), preserving internal newlines and
  formatting.

_Example:_

```urpc
"""
Docstring for MyType.
Can be multi-line.
"""
type MyType {
  // ...
}
```

### 7.1 Field Docstrings

- For fields, prefer concise, single-line docstrings.
- Place the docstring on the line immediately before the field it documents,
  indented to the same level.
- When multiple fields have docstrings, separate each one with a blank line to
  improve readability.

_Example:_

```urpc
type User {
  """ The user's unique identifier. """
  id: string

  """ The user's email address. Must be unique. """
  email: string

  name?: string // A field without a docstring does not require a blank line before it.
}
```

## 8. Deprecation

The `deprecated` keyword is used to mark types, procedures, or streams as
deprecated.

- Place the `deprecated` keyword on its own line immediately before the element
  definition
- If a docstring exists, place the `deprecated` keyword between the docstring
  and the element definition
- For deprecation with a message, use parentheses with the message in quotes

### 8.1 Basic Deprecation

_Example:_

```urpc
deprecated type MyType {
  // type definition
}

"""
Documentation for MyProc
"""
deprecated proc MyProc {
  // procedure definition
}

deprecated stream MyStream {
  // stream definition
}
```

### 8.2 Deprecation with Message

_Example:_

```urpc
"""
Documentation for MyType
"""
deprecated("Replaced by ImprovedType")
type MyType {
  // type definition
}

deprecated("Use NewStream instead")
stream MyStream {
  // stream definition
}
```

## 9. Naming Conventions

### 9.1 Type, Procedure, and Stream Names

- Use **strict PascalCase** (also known as UpperCamelCase). Each word starts with an uppercase letter with no underscores or consecutive capital letters.
- Acronyms longer than two letters should be treated as regular words (e.g. `HttpRequest`, not `HTTPRequest`).

_Example:_

```urpc
// Correct
type FooBar {
  myField: string
}

// Incorrect
type FooBAR {
  myField: string
}
```

### 9.2 Field Names

- Use **strict camelCase**. The first word is lowercase and each subsequent word starts with an uppercase letter. Do not use underscores or all-caps abbreviations.

_Example:_

```urpc
// Correct
myInput: FooBar
zipCode: string

// Incorrect
MyINPUT: FooBar
zip_code: string
```
