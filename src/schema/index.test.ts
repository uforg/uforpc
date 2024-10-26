import { assertEquals, assertRejects } from "@std/assert";
import {
  parseSchema,
  SchemaParsingError,
  SchemaValidationError,
} from "./index.ts";

Deno.test("parseSchema - valid basic schema", async () => {
  const validSchema = `{
    "procedures": [
      {
        "name": "Hello",
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
  assertEquals(schema.procedures[0].name, "Hello");
});

Deno.test("parseSchema - nested arrays", async () => {
  const schema = `{
    "types": [
      {
        "name": "Matrix",
        "fields": {
          "data": "number[][]",
          "labels": "string[][]"
        }
      }
    ],
    "procedures": [
      {
        "name": "ProcessMatrix",
        "type": "query",
        "input": {
          "matrix": "Matrix"
        },
        "output": {
          "result": "number[][]"
        }
      }
    ]
  }`;

  const parsed = await parseSchema(schema);
  assertEquals(parsed.types?.length, 1);
  assertEquals(parsed.procedures.length, 1);
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

Deno.test("parseSchema - procedure name validation", async () => {
  const invalidSchema = `{
    "procedures": [
      {
        "name": "invalidName",
        "type": "query"
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
        "name": "GetUser",
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
  );
});

Deno.test("parseSchema - invalid array type", async () => {
  const invalidArraySchema = `{
    "types": [
      {
        "name": "Invalid",
        "fields": {
          "bad": "[string]"
        }
      }
    ],
    "procedures": []
  }`;

  await assertRejects(
    () => parseSchema(invalidArraySchema),
    SchemaValidationError,
    "Schema validation failed",
  );
});
