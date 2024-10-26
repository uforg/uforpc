// deno-lint-ignore-file require-await
import { assertEquals, assertRejects } from "@std/assert";
import {
  parseSchema,
  SchemaParsingError,
  SchemaValidationError,
} from "./index.ts";
import type { ArrayType } from "./types.ts";

Deno.test("parseSchema - basic valid schema", async () => {
  const validSchema = `{
    "procedures": [
      {
        "name": "GetUser",
        "type": "query",
        "input": {
          "id": "string"
        },
        "output": {
          "user": {
            "type": "object",
            "fields": {
              "id": "string",
              "name": "string"
            }
          }
        }
      }
    ]
  }`;

  const schema = parseSchema(validSchema);
  assertEquals(schema.procedures.length, 1);
  assertEquals(schema.procedures[0].name, "GetUser");
  assertEquals(schema.procedures[0].input?.id.type, "string");
});

Deno.test("parseSchema - complete schema with types", async () => {
  const validSchema = `{
    "types": [
      {
        "name": "User",
        "desc": "User type",
        "fields": {
          "id": "string",
          "name": "string",
          "age": {
            "type": "int",
            "desc": "User's age"
          },
          "roles": {
            "type": "string[][]",
            "desc": "User's roles"
          },
          "permissions": "string[]" 
        }
      }
    ],
    "procedures": [
      {
        "name": "CreateUser",
        "type": "mutation",
        "input": {
          "user": "User"
        },
        "output": {
          "id": "string"
        },
        "meta": {
          "requiresAuth": true,
          "rateLimit": 100,
          "ranking": 123.456
        }
      }
    ]
  }`;

  const schema = parseSchema(validSchema);
  assertEquals(schema.types?.length, 1);
  assertEquals(schema.types?.[0].fields.name.type, "string");
  assertEquals(schema.types?.[0].fields.age.type, "int");
  assertEquals(
    (schema.types?.[0].fields.roles.type as ArrayType).baseType,
    "string",
  );
  assertEquals(
    (schema.types?.[0].fields.roles.type as ArrayType).dimensions,
    2,
  );
  assertEquals(
    (schema.types?.[0].fields.permissions.type as ArrayType).baseType,
    "string",
  );
  assertEquals(
    (schema.types?.[0].fields.permissions.type as ArrayType).dimensions,
    1,
  );
  assertEquals(schema.types?.[0].name, "User");
  assertEquals(schema.procedures[0].meta?.requiresAuth, true);
  assertEquals(schema.procedures[0].meta?.rateLimit, 100);
  assertEquals(schema.procedures[0].meta?.ranking, 123.456);
});

Deno.test("parseSchema - array types", async () => {
  const schema = `{
    "types": [
      {
        "name": "Matrix",
        "fields": {
          "data": "int[][]",
          "labels": "string[]"
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
          "result": "int[][]"
        }
      }
    ]
  }`;

  const parsed = parseSchema(schema);
  const matrixType = parsed.types?.[0];
  assertEquals((matrixType?.fields.data.type as ArrayType).dimensions, 2);
  assertEquals((matrixType?.fields.labels.type as ArrayType).dimensions, 1);
});

Deno.test("parseSchema - nested objects", async () => {
  const schema = `{
    "procedures": [
      {
        "name": "CreateProfile",
        "type": "mutation",
        "input": {
          "profile": {
            "type": "object",
            "fields": {
              "name": "string",
              "address": {
                "type": "object",
                "fields": {
                  "street": "string",
                  "city": "string"
                }
              }
            }
          }
        }
      }
    ]
  }`;

  const parsed = parseSchema(schema);
  const procedure = parsed.procedures[0];
  const profile = procedure.input?.profile;
  assertEquals(profile?.type, "object");
  assertEquals(profile?.fields?.name.type, "string");
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
    async () => parseSchema(invalidJson),
    SchemaParsingError,
    "Invalid JSON",
  );
});

Deno.test("parseSchema - schema validation errors", async (t) => {
  await t.step("invalid procedure name", async () => {
    const invalidSchema = `{
      "procedures": [
        {
          "name": "invalidName",
          "type": "query"
        }
      ]
    }`;

    await assertRejects(
      async () => parseSchema(invalidSchema),
      SchemaValidationError,
      "Schema validation failed",
    );
  });

  await t.step("invalid type pattern", async () => {
    const invalidSchema = `{
      "procedures": [
        {
          "name": "Test",
          "type": "query",
          "input": {
            "test": "[int]"
          }
        }
      ]
    }`;

    await assertRejects(
      async () => parseSchema(invalidSchema),
      SchemaValidationError,
      "Schema validation failed",
    );
  });

  await t.step("invalid meta value type", async () => {
    const invalidSchema = `{
      "procedures": [
        {
          "name": "Test",
          "type": "query",
          "meta": {
            "test": { "invalid": "object" }
          }
        }
      ]
    }`;

    await assertRejects(
      async () => parseSchema(invalidSchema),
      SchemaValidationError,
      "Schema validation failed",
    );
  });
});

Deno.test("parseSchema - object fields validations", async (t) => {
  await t.step("Rejects schema when object fields are missing", async () => {
    const schema = `{
      "procedures": [
        {
          "name": "GetUser",
          "type": "query",
          "input": {
            "user": "object"
          }
        }
      ]
    }`;

    await assertRejects(
      async () => parseSchema(schema),
      SchemaValidationError,
      "Schema validation failed",
    );
  });

  await t.step("Rejects schema when object fields are empty", async () => {
    const schema = `{
      "procedures": [
        {
          "name": "GetUser",
          "type": "query",
          "input": {
            "user": {
              "type": "object",
              "fields": {}
            }
          }
        }
      ]
    }`;

    await assertRejects(
      async () => parseSchema(schema),
      SchemaValidationError,
      "Schema validation failed",
    );
  });

  await t.step(
    "Rejects schema when type is not object but fields are provided",
    async () => {
      const schema = `{
        "procedures": [
          {
            "name": "GetUser",
            "type": "query",
            "input": {
              "user": {
                "type": "string",
                "fields": {
                  "name": "string"
                }
              }
            }
          }
        ]
      }`;

      await assertRejects(
        async () => parseSchema(schema),
        SchemaValidationError,
        "Schema validation failed",
      );
    },
  );

  await t.step(
    "Parses schema when object and fields are correctly defined",
    async () => {
      const schema = `{
        "procedures": [
          {
            "name": "GetUser",
            "type": "query",
            "input": {
              "user": {
                "type": "object",
                "fields": {
                  "name": "string"
                }
              }
            }
          }
        ]
      }`;

      const parsed = parseSchema(schema);
      assertEquals(
        parsed.procedures[0].input?.user?.fields?.name.type,
        "string",
      );
    },
  );
});

Deno.test("parseSchema - edge cases", async (t) => {
  await t.step("empty arrays", async () => {
    const schema = `{
      "types": [],
      "procedures": []
    }`;

    await assertRejects(
      async () => parseSchema(schema),
      SchemaValidationError,
      "Schema validation failed",
    );
  });

  await t.step("deeply nested arrays", async () => {
    const schema = `{
      "procedures": [
        {
          "name": "Test",
          "type": "query",
          "input": {
            "data": "int[][][][][]"
          }
        }
      ]
    }`;

    const parsed = parseSchema(schema);
    assertEquals(
      (parsed.procedures[0].input?.data.type as ArrayType).dimensions,
      5,
    );
  });

  await t.step("complex nested structure", async () => {
    const schema = `{
      "types": [
        {
          "name": "Complex",
          "fields": {
            "data": {
              "type": "object",
              "fields": {
                "matrix": "int[][]",
                "metadata": {
                  "type": "object",
                  "fields": {
                    "tags": "string[]"
                  }
                }
              }
            }
          }
        }
      ],
      "procedures": [
        {
          "name": "Process",
          "type": "query",
          "input": {
            "complex": "Complex"
          }
        }
      ]
    }`;

    const parsed = parseSchema(schema);
    const complexType = parsed.types?.[0];
    assertEquals(complexType?.name, "Complex");
    assertEquals(complexType?.fields.data.type, "object");
  });
});
