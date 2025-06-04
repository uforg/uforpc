# UFO-RPC DSL (URPC) Formatting Guide

This document specifies the standard formatting rules for the UFO-RPC DSL
(URPC). Consistent formatting enhances code readability, maintainability, and
collaboration. The primary goal is to produce clean, predictable, and
aesthetically pleasing URPC code.

This specification is enforced by the official URPC formatter.

## 1. General Principles

This guide is for reference only. All the formatting rules are enforced by the
official URPC formatter included in the UFO-RPC CLI so you don't have to worry
about it, continue reading for reference only.

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

Top-level elements include `version`, `rule`, `type`, `proc`, `stream`, and
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
procedure's input, output, meta, or inline object, and stream's input, output,
meta.

- Each field is placed on its own line.
- **Field Separation:** It is recommended to leave one blank line between fields
  that include validation rules or complex formatting, to visually separate
  them. For simpler fields (with no rules or extra comments), fields may be
  placed consecutively without a blank line.

_Recommended for simple fields:_

Do not use blank lines between fields.

```urpc
address: {
  street: string
  city: string
  zip: string
}
```

_Recommended for fields with rules or comments:_

Use blank lines between fields to visually separate them.

```urpc
address: {
  street: string // Inline comment
    @rule1
    @rule2(error: "Invalid street")

  city: string
    @rule3
    @rule4(error: "Invalid city") // Other comment

  zip: string
    @rule5
    @rule6(error: "Invalid zip")
}
```

### 4.2 Blocks (Type, Input/Output, Meta)

- Opening braces (`{`) are on the same line as the declaration header (preceded
  by one space).
- Contents inside non-empty blocks always start on a new, indented line.
- The closing brace (`}`) is placed on its own line, aligned with the opening
  line.
- In procedure and stream bodies, separate the `input`, `output`, and `meta` blocks with
  one blank line.

## 5. Spacing

- **Colons (`:`):** No space before; one space after (e.g. `field: string`).
- **Commas (`,`):** No space before; one space after.
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

- **Standalone Comments:** Use `//` or `/* … */` on their own lines; indent to
  the current block.
- **End-of-Line (EOL) Comments:** Can use either `//` or block style (`/* … */`)
  following code on the same line, with at least one space separating them.

_Example:_

```urpc
version 1 // EOL comment

type Example {
  field: string // Inline comment for field
    @rule /* Inline rule comment */
}
```

## 7. Docstrings

- Place docstrings immediately above the `rule`, `type`, `proc`, or `stream` they
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

"""
Docstring for MyStream.
"""
stream MyStream {
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

## 9. Deprecation

The `deprecated` keyword is used to mark rules, types, procedures, or streams as
deprecated.

- Place the `deprecated` keyword on its own line immediately before the element
  definition
- If a docstring exists, place the `deprecated` keyword between the docstring
  and the element definition
- For deprecation with a message, use parentheses with the message in quotes

### 9.1 Basic Deprecation

_Example:_

```urpc
deprecated rule @myRule {
  // rule definition
}

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

### 9.2 Deprecation with Message

_Example:_

```urpc
deprecated("Use newRule instead")
rule @myRule {
  // rule definition
}

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
