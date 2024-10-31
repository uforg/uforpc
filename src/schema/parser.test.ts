// deno-lint-ignore-file require-await
import { assertEquals, assertRejects } from "@std/assert";
import {
  parseSchema,
  SchemaParsingError,
  SchemaValidationError,
} from "./parser.ts";
import { parseArrayType } from "@/schema/helpers.ts";

Deno.test("parseSchema - basic valid schema", async () => {
  const validSchema = `{
    "procedures": [
      {
        "name": "GetUser",
        "type": "query",
        "input": {
          "id": {
            "type": "string"
          }
        },
        "output": {
          "user": {
            "type": "object",
            "fields": {
              "id": {
                "type": "string"
              },
              "name": {
                "type": "string"
              }
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
          "id": {
            "type": "string"
          },
          "name": {
            "type": "string"
          },
          "age": {
            "type": "int",
            "desc": "User's age"
          },
          "roles": {
            "type": "string[][]",
            "desc": "User's roles"
          },
          "permissions": {
            "type": "string[]"
          } 
        }
      }
    ],
    "procedures": [
      {
        "name": "CreateUser",
        "type": "mutation",
        "input": {
          "user": {
            "type": "User"
          }
        },
        "output": {
          "id": {
            "type": "string"
          }
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
    parseArrayType(schema.types?.[0].fields.roles!).type.type,
    "string",
  );
  assertEquals(
    parseArrayType(schema.types?.[0].fields.roles!).dimensions,
    2,
  );
  assertEquals(
    parseArrayType(schema.types?.[0].fields.permissions!).type.type,
    "string",
  );
  assertEquals(
    parseArrayType(schema.types?.[0].fields.permissions!).dimensions,
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
          "data": {
            "type": "int[][]"
          },
          "labels": {
            "type": "string[]"
          }
        }
      }
    ],
    "procedures": [
      {
        "name": "ProcessMatrix",
        "type": "query",
        "input": {
          "matrix": {
            "type": "Matrix"
          }
        },
        "output": {
          "result": {
            "type": "int[][]"
          }
        }
      }
    ]
  }`;

  const parsed = parseSchema(schema);
  const matrixType = parsed.types?.[0];

  assertEquals(
    parseArrayType(matrixType?.fields.data!).type.type,
    "int",
  );
  assertEquals(
    parseArrayType(matrixType?.fields.labels!).type.type,
    "string",
  );

  assertEquals(
    parseArrayType(matrixType?.fields.data!).dimensions,
    2,
  );
  assertEquals(
    parseArrayType(matrixType?.fields.labels!).dimensions,
    1,
  );
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
              "name": {
                "type": "string"
              },
              "address": {
                "type": "object",
                "fields": {
                  "street": {
                    "type": "string"
                  },
                  "city": {
                    "type": "string"
                  }
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
          "id": {
            "type": "string"
          },
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
            "test": {
              "type": "[int]"
            }
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
            "user": {
              "type": "object"
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
                  "name": {
                    "type": "string"
                  }
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
            "data": {
              "type": "int[][][][][]"
            }
          }
        }
      ]
    }`;

    const parsed = parseSchema(schema);
    assertEquals(
      parseArrayType(parsed.procedures[0].input?.data!).type.type,
      "int",
    );
    assertEquals(
      parseArrayType(parsed.procedures[0].input?.data!).dimensions,
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
                "matrix": {
                  "type": "int[][]"
                },
                "metadata": {
                  "type": "object",
                  "fields": {
                    "tags": {
                      "type": "string[]"
                    }
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
            "complex": {
              "type": "Complex"
            }
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

Deno.test("parseSchema - duplicate type names", async () => {
  const schema = `{
    "types": [
      {
        "name": "User",
        "fields": {
          "id": {
            "type": "string"
          }
        }
      },
      {
        "name": "User",
        "fields": {
          "name": {
            "type": "string"
          }
        }
      }
    ],
    "procedures": [
      {
        "name": "GetUser",
        "type": "query",
        "input": {
          "id": {
            "type": "string"
          }
        },
        "output": {
          "user": {
            "type": "User"
          }
        }
      }
    ]
  }`;

  await assertRejects(
    async () => parseSchema(schema),
    SchemaValidationError,
    "Duplicate type name",
  );
});

Deno.test("parseSchema - undefined custom types", async () => {
  const schema = `{
    "procedures": [
      {
        "name": "GetUser",
        "type": "query",
        "input": {
          "id": {
            "type": "string"
          }
        },
        "output": {
          "user": {
            "type": "User"
          }
        }
      }
    ]
  }`;

  await assertRejects(
    async () => parseSchema(schema),
    SchemaValidationError,
    "Custom type User is not defined",
  );
});
