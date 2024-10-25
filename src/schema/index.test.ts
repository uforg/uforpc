import { assertEquals, assertRejects } from "@std/assert";
import {
  parseSchema,
  SchemaParsingError,
  SchemaValidationError,
} from "./index.ts";
import type { Schema } from "./types.ts";

Deno.test("parseSchema - valid basic schema", async () => {
  const validSchema = `{
    "procedures": [
      {
        "name": "hello",
        "type": "query",
        "input": {
          "name": "string"
        },
        "output": {
          "message": "string"
        }
      }
    ]
  }`;

  const schema = await parseSchema(validSchema);
  assertEquals(schema.procedures.length, 1);
  assertEquals(schema.procedures[0].name, "hello");
});

Deno.test("parseSchema - complete type system", async () => {
  const completeSchema = `{
    "types": [
      {
        "name": "Address",
        "fields": {
          "street": "string",
          "number": "number",
          "zipCode": "string"
        }
      },
      {
        "name": "User",
        "fields": {
          "id": "string",
          "name": "string",
          "age": "number",
          "active": "boolean",
          "score": "float",
          "address": "Address",
          "tags": "string[]",
          "metadata": {
            "type": "object",
            "fields": {
              "lastLogin": "string",
              "loginCount": "number"
            }
          }
        }
      }
    ],
    "procedures": [
      {
        "name": "createUser",
        "type": "mutation",
        "input": {
          "user": "User"
        },
        "output": {
          "id": "string",
          "success": "boolean"
        }
      }
    ]
  }`;

  const schema = await parseSchema(completeSchema) as Schema;
  assertEquals(schema.types?.length, 2);
  assertEquals(schema.procedures.length, 1);
});

Deno.test("parseSchema - invalid JSON", async () => {
  const invalidJson = `{
    "types": [
      {
        "name": "User",
        "fields": {
          "id": "string",
        }
      }
    ],
  }`;

  await assertRejects(
    () => parseSchema(invalidJson),
    SchemaParsingError,
    "Invalid JSON",
  );
});

Deno.test("parseSchema - missing required procedure type", async () => {
  const invalidSchema = `{
    "procedures": [
      {
        "name": "test"
      }
    ]
  }`;

  await assertRejects(
    () => parseSchema(invalidSchema),
    SchemaValidationError,
    "Schema validation failed",
  );
});

Deno.test("parseSchema - undefined custom type", async () => {
  const schemaWithUndefinedType = `{
    "procedures": [
      {
        "name": "getUser",
        "type": "query",
        "output": {
          "user": "NonExistentType"
        }
      }
    ]
  }`;

  await assertRejects(
    () => parseSchema(schemaWithUndefinedType),
    SchemaValidationError,
    "Invalid type references found",
  );
});

Deno.test("parseSchema - circular type references", async () => {
  const circularSchema = `{
    "types": [
      {
        "name": "A",
        "fields": {
          "b": "B"
        }
      },
      {
        "name": "B",
        "fields": {
          "a": "A"
        }
      }
    ],
    "procedures": []
  }`;

  const schema = await parseSchema(circularSchema);
  assertEquals(schema.types?.length, 2);
});

Deno.test("parseSchema - array edge cases", async () => {
  const arraySchema = `{
    "types": [
      {
        "name": "Complex",
        "fields": {
          "strings": "string[]",
          "numbers": "number[]",
          "nested": {
            "type": "object[]",
            "fields": {
              "value": "string"
            }
          },
          "complex": "Complex[]"
        }
      }
    ],
    "procedures": []
  }`;

  const schema = await parseSchema(arraySchema);
  assertEquals(schema.types?.length, 1);
});

Deno.test("parseSchema - deep nested objects", async () => {
  const deepSchema = `{
    "types": [
      {
        "name": "Deep",
        "fields": {
          "level1": {
            "type": "object",
            "fields": {
              "level2": {
                "type": "object",
                "fields": {
                  "level3": {
                    "type": "object",
                    "fields": {
                      "value": "string"
                    }
                  }
                }
              }
            }
          }
        }
      }
    ],
    "procedures": []
  }`;

  const schema = await parseSchema(deepSchema);
  assertEquals(schema.types?.length, 1);
});

Deno.test("parseSchema - empty input/output procedures", async () => {
  const emptySchema = `{
    "procedures": [
      {
        "name": "noIO",
        "type": "query"
      },
      {
        "name": "onlyInput",
        "type": "mutation",
        "input": {
          "value": "string"
        }
      },
      {
        "name": "onlyOutput",
        "type": "query",
        "output": {
          "value": "string"
        }
      }
    ]
  }`;

  const schema = await parseSchema(emptySchema);
  assertEquals(schema.procedures.length, 3);
});

Deno.test("parseSchema - procedure naming conventions", async () => {
  const namingSchema = `{
    "procedures": [
      {
        "name": "validName",
        "type": "query"
      }
    ]
  }`;

  const schema = await parseSchema(namingSchema);
  assertEquals(schema.procedures[0].name, "validName");

  const invalidNameSchema = `{
    "procedures": [
      {
        "name": "InvalidName",
        "type": "query"
      }
    ]
  }`;

  await assertRejects(
    () => parseSchema(invalidNameSchema),
    SchemaValidationError,
  );
});

Deno.test("parseSchema - all primitive types", async () => {
  const primitivesSchema = `{
    "types": [
      {
        "name": "Primitives",
        "fields": {
          "string": "string",
          "number": "number",
          "float": "float",
          "boolean": "boolean"
        }
      }
    ],
    "procedures": []
  }`;

  const schema = await parseSchema(primitivesSchema);
  assertEquals(Object.keys(schema.types![0].fields).length, 4);
});
