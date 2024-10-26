# UFO RPC Schema Module

The Schema module handles the parsing, validation, and transformation of UFO RPC
schemas into a normalized internal representation optimized for code generation.

## Process

1. **Schema Parsing**: Takes a JSON string input and parses it into a raw
   JavaScript object, ensuring the input is valid JSON.

2. **Schema Validation**: Validates the parsed object against the UFO RPC JSON
   Schema specification, ensuring:
   - Valid structure
   - Correct types and patterns
   - Required fields
   - Type naming conventions
   - Valid array notations

3. **Schema Transformation**: Transforms the validated schema into a normalized
   internal representation that:
   - Unifies field representations using DetailedField
   - Normalizes array types with explicit dimensions
   - Maintains type safety
   - Simplifies nested structures

## Key Features

- **Type Safety**: Full TypeScript support with comprehensive type definitions
- **Error Handling**: Clear error types for parsing and validation issues
- **Utility Functions**: Helper functions for type checking and manipulation
- **Array Support**: Robust handling of array types and nested arrays
- **Object Nesting**: Full support for nested object types and complex
  structures

## Usage

The module exports a single `parseSchema` function that handles the entire
process:

```typescript
const schema = await parseSchema(jsonString);
```

For detailed type definitions, validation rules, and utility functions, please
refer to the source code.
