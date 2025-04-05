# UFO-RPC DSL (URPC) Formatting Guide

This document specifies the standard formatting rules for the UFO-RPC DSL
(URPC). Consistent formatting enhances code readability, maintainability, and
collaboration. The primary goal is to produce clean, predictable, and
aesthetically pleasing URPC code.

This specification is enforced by the official URPC formatter.

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
    @rule 
}
```

## 3. Top-Level Elements

Top-level elements include `version`, `import`, `rule`, `type`, `proc`, and
standalone comments.

- **Default:** Separate each top-level element with one blank line.
- **Exceptions:**
  - **First Element:** No blank line before the very first element.
  - **Consecutive Comments:** Do not insert extra blank lines between
    consecutive standalone comments.
  - **Consecutive Imports:** Successive `import` statements share no blank line
    between them.
  - **Imports vs. Others:** An `import` must be separated from any non-import
    element by a blank line.
  - **Following a Comment:** When an element follows a standalone comment, do
    not add an extra blank line unless the source intentionally contains one.
- **Preservation:** Intentionally placed blank lines in the source (e.g. between
  comments or imports) are respected.

_Example:_

```urpc
version 1

import "a.urpc"
import "b.urpc"

// A standalone comment
// Another standalone comment
type TypeA {
  field: string
}

import "c.urpc"  // (Here, an import following a type must be separated from the type block by a blank line)
```

## 4. Fields and Blocks

### 4.1 Fields in a Type

This section applies to fields in a type block, fields in an input, output or
meta of a procedure, and fields in an inline object.

- Each field within a type block is placed on its own line.
- **Separate fields with one blank line.** This ensures a clear visual
  separation between consecutive fields.
- Field rules and inline comments (end-of-line comments) appear immediately
  after the field definition (without an extra blank line).

_Incorrect:_

```urpc
address: { // Inline object field
  street: string
  city: string
  zip: string
}
```

_Correct:_

```urpc
address: { // Inline object field
  street: string

  city: string

  zip: string
}
```

### 4.2 Blocks (Type, Input/Output, Meta)

- Opening braces (`{`) are on the same line as the declaration header (preceded
  by one space).
- Contents inside non-empty blocks always start on a new, indented line.
- The closing brace (`}`) is placed on its own line, aligned with the opening
  line.
- In procedure bodies, separate the `input`, `output`, and `meta` blocks with
  one blank line.

## 5. Spacing

- **Colons (`:`):** No space before; one space after (e.g. `field: string`).
- **Commas (`,`):** No space before; one space after (e.g. `extends A, B`).
- **Braces (`{` and `}`):** One space before `{` in declarations; inside blocks,
  use newlines and proper indentation.
- **Brackets (`[]`):** No spaces for array types (e.g. `string[]`); no extra
  interior spacing.
- **Parentheses (`()`):** No extra spaces inside the parentheses.
- **At Symbol (`@`):** Immediately followed by the rule name (e.g. `@minlen`).
- **Optional Marker (`?`):** Immediately follows the field name (e.g.
  `email?: string`).

## 6. Comments

Comment content is preserved exactly (including internal whitespace).

- **Standalone Comments:**\
  Use `//` or `/* … */` on their own lines; indent to the current block.

- **End-of-Line (EOL) Comments:**\
  Can use either `//` or block style (`/* … */`) following code on the same
  line, with at least one space separating them.

_Example:_

```urpc
version 1 // EOL comment

import "path" /* EOL comment */

type Example {
  field: string // Inline comment for field
    @rule /* Inline rule comment */
}
```

## 7. Docstrings

- Place docstrings immediately above the `rule`, `type`, or `proc` they
  document.
- They are enclosed in triple quotes (`"""`), preserving internal newlines and
  formatting.

_Example:_

```urpc
"""
Docstring for MyType.
"""
type MyType {
  // ...
}
```

## 8. Rules on Fields

Field validation rules are indented one level deeper than the field they apply
to and are placed immediately after the field definition. Each rule begins on a
new line within the field block.

_Example:_

```urpc
type User {
  name: string
    @minlen(3)
    @maxlen(50, error: "Name too long")
}
```
