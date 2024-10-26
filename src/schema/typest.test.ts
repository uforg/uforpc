// deno-lint-ignore-file
import { assertEquals, assertThrows } from "@std/assert";
import {
  type ArrayType,
  type DetailedField,
  type FieldType,
  fieldTypeToString,
  flattenArrayType,
  getBaseFieldType,
  getTotalArrayDimensions,
  isArrayType,
  isCustomType,
  isDetailedField,
  isObjectType,
  isPrimitiveType,
  isValidFieldType,
  parseArrayType,
  parseDetailedField,
  parseFieldType,
} from "./types.ts";

Deno.test("isPrimitiveType", async (t) => {
  await t.step("should identify valid primitive types", () => {
    assertEquals(isPrimitiveType("string"), true);
    assertEquals(isPrimitiveType("number"), true);
    assertEquals(isPrimitiveType("float"), true);
    assertEquals(isPrimitiveType("boolean"), true);
  });

  await t.step("should reject non-primitive types", () => {
    assertEquals(isPrimitiveType("object"), false);
    assertEquals(isPrimitiveType("User"), false);
    assertEquals(isPrimitiveType({ baseType: "string", dimensions: 1 }), false);
  });

  await t.step("should handle edge cases", () => {
    assertEquals(isPrimitiveType("String"), false);
    assertEquals(isPrimitiveType("NUMBER"), false);
    assertEquals(isPrimitiveType(""), false);
    assertEquals(isPrimitiveType(" string "), false);
  });
});

Deno.test("isCustomType", async (t) => {
  await t.step("should identify valid custom types", () => {
    assertEquals(isCustomType("User"), true);
    assertEquals(isCustomType("UserProfile"), true);
    assertEquals(isCustomType("A"), true);
    assertEquals(isCustomType("ABC123"), true);
  });

  await t.step("should reject invalid custom types", () => {
    assertEquals(isCustomType("user"), false);
    assertEquals(isCustomType("123User"), false);
    assertEquals(isCustomType("User_Profile"), false);
    assertEquals(isCustomType("User-Profile"), false);
    assertEquals(isCustomType(""), false);
  });

  await t.step("should reject primitive and array types", () => {
    assertEquals(isCustomType("string"), false);
    assertEquals(isCustomType("object"), false);
    assertEquals(isCustomType({ baseType: "User", dimensions: 1 }), false);
  });
});

Deno.test("isArrayType", async (t) => {
  await t.step("should identify valid array types", () => {
    assertEquals(isArrayType({ baseType: "string", dimensions: 1 }), true);
    assertEquals(isArrayType({ baseType: "User", dimensions: 2 }), true);
    assertEquals(
      isArrayType({
        baseType: { baseType: "number", dimensions: 1 },
        dimensions: 1,
      }),
      true,
    );
  });

  await t.step("should reject non-array types", () => {
    assertEquals(isArrayType("string"), false);
    assertEquals(isArrayType("object"), false);
    assertEquals(isArrayType("User"), false);
  });

  await t.step("should handle edge cases", () => {
    assertEquals(isArrayType({ baseType: "string" } as any), false);
    assertEquals(isArrayType({ dimensions: 1 } as any), false);
    assertEquals(isArrayType({} as any), false);
  });
});

Deno.test("isObjectType", async (t) => {
  await t.step("should identify object type", () => {
    assertEquals(isObjectType("object"), true);
  });

  await t.step("should reject non-object types", () => {
    assertEquals(isObjectType("string"), false);
    assertEquals(isObjectType("User"), false);
    assertEquals(isObjectType({ baseType: "string", dimensions: 1 }), false);
  });
});

Deno.test("isValidFieldType", async (t) => {
  await t.step("should validate primitive types", () => {
    assertEquals(isValidFieldType("string"), true);
    assertEquals(isValidFieldType("number"), true);
    assertEquals(isValidFieldType("float"), true);
    assertEquals(isValidFieldType("boolean"), true);
  });

  await t.step("should validate object and custom types", () => {
    assertEquals(isValidFieldType("object"), true);
    assertEquals(isValidFieldType("User"), true);
    assertEquals(isValidFieldType("UserProfile"), true);
  });

  await t.step("should validate array types", () => {
    assertEquals(isValidFieldType({ baseType: "string", dimensions: 1 }), true);
    assertEquals(isValidFieldType({ baseType: "User", dimensions: 2 }), true);
    assertEquals(
      isValidFieldType({
        baseType: { baseType: "number", dimensions: 1 },
        dimensions: 1,
      }),
      true,
    );
  });

  await t.step("should reject invalid types", () => {
    assertEquals(isValidFieldType("invalid" as any), false);
    assertEquals(
      isValidFieldType({ baseType: "invalid", dimensions: 1 } as any),
      false,
    );
    assertEquals(isValidFieldType({} as any), false);
  });
});

Deno.test("isDetailedField", async (t) => {
  await t.step("should identify valid detailed fields", () => {
    assertEquals(isDetailedField({ type: "string" }), true);
    assertEquals(isDetailedField({ type: "User", desc: "A user" }), true);
    assertEquals(
      isDetailedField({ type: { baseType: "string", dimensions: 1 } }),
      true,
    );
  });

  await t.step("should reject invalid values", () => {
    assertEquals(isDetailedField(null), false);
    assertEquals(isDetailedField(undefined), false);
    assertEquals(isDetailedField("string"), false);
    assertEquals(isDetailedField({}), false);
    assertEquals(isDetailedField({ desc: "missing type" }), false);
  });
});

Deno.test("parseArrayType", async (t) => {
  await t.step("should parse valid array type strings", () => {
    assertEquals(parseArrayType("string[]"), {
      baseType: "string",
      dimensions: 1,
    });
    assertEquals(parseArrayType("number[][]"), {
      baseType: "number",
      dimensions: 2,
    });
    assertEquals(parseArrayType("User[][][]"), {
      baseType: "User",
      dimensions: 3,
    });
    assertEquals(parseArrayType("object[][]"), {
      baseType: "object",
      dimensions: 2,
    });
  });

  await t.step("should parse nested array types", () => {
    const result = parseArrayType("string[][]");
    assertEquals(result?.baseType, "string");
    assertEquals(result?.dimensions, 2);
  });

  await t.step("should return null for non-array types", () => {
    assertEquals(parseArrayType("string"), null);
    assertEquals(parseArrayType("User"), null);
    assertEquals(parseArrayType("object"), null);
  });

  await t.step("should handle edge cases", () => {
    assertEquals(parseArrayType(""), null);
    assertEquals(parseArrayType("[]"), null);
    assertEquals(parseArrayType("[string]"), null);
    assertEquals(parseArrayType("string["), null);
    assertEquals(parseArrayType("string]"), null);
    assertEquals(parseArrayType("string[][]a"), null);
  });
});

Deno.test("parseFieldType", async (t) => {
  await t.step("should parse primitive types", () => {
    assertEquals(parseFieldType("string"), "string");
    assertEquals(parseFieldType("number"), "number");
    assertEquals(parseFieldType("float"), "float");
    assertEquals(parseFieldType("boolean"), "boolean");
  });

  await t.step("should parse object type", () => {
    assertEquals(parseFieldType("object"), "object");
  });

  await t.step("should parse custom types", () => {
    assertEquals(parseFieldType("User"), "User");
    assertEquals(parseFieldType("UserProfile"), "UserProfile");
  });

  await t.step("should parse array types", () => {
    assertEquals(parseFieldType("string[]"), {
      baseType: "string",
      dimensions: 1,
    });
    assertEquals(parseFieldType("User[][]"), {
      baseType: "User",
      dimensions: 2,
    });
  });

  await t.step("should handle edge cases", () => {
    assertEquals(parseFieldType(""), null);
    assertEquals(parseFieldType("[]"), null);
    assertEquals(parseFieldType("[string]"), null);
    assertEquals(parseFieldType("123Type"), null);
    assertEquals(parseFieldType("Type_123"), null);
  });
});

Deno.test("parseDetailedField", async (t) => {
  await t.step("should parse string to detailed field", () => {
    assertEquals(parseDetailedField("string"), { type: "string" });
    assertEquals(parseDetailedField("User"), { type: "User" });
    assertEquals(parseDetailedField("string[]"), {
      type: { baseType: "string", dimensions: 1 },
    });
  });

  await t.step("should parse partial detailed field", () => {
    assertEquals(
      parseDetailedField({ type: "string", desc: "A string field" }),
      { type: "string", desc: "A string field" },
    );
  });

  await t.step("should parse detailed field with fields", () => {
    const field = {
      type: "object",
      desc: "An object field",
      fields: {
        name: { type: "string" },
        age: { type: "number" },
      },
    };
    assertEquals(parseDetailedField(field), field);
  });

  await t.step("should throw error for invalid types", () => {
    assertThrows(
      () => parseDetailedField("invalid"),
      Error,
      "Invalid type: invalid",
    );

    assertThrows(
      () => parseDetailedField({ type: "invalid" }),
      Error,
      "Invalid type: invalid",
    );

    assertThrows(
      () => parseDetailedField({} as any),
      Error,
      "DetailedField must have a type",
    );
  });
});

Deno.test("getBaseFieldType", async (t) => {
  await t.step("should return base type for non-array types", () => {
    assertEquals(getBaseFieldType("string"), "string");
    assertEquals(getBaseFieldType("User"), "User");
    assertEquals(getBaseFieldType("object"), "object");
  });

  await t.step("should return base type from array types", () => {
    assertEquals(
      getBaseFieldType({ baseType: "string", dimensions: 1 }),
      "string",
    );
    assertEquals(getBaseFieldType({ baseType: "User", dimensions: 2 }), "User");
  });

  await t.step("should handle nested array types", () => {
    const nestedType: ArrayType = {
      baseType: { baseType: "string", dimensions: 1 },
      dimensions: 1,
    };
    assertEquals(getBaseFieldType(nestedType), "string");
  });
});

Deno.test("getTotalArrayDimensions", async (t) => {
  await t.step("should return 0 for non-array types", () => {
    assertEquals(getTotalArrayDimensions("string"), 0);
    assertEquals(getTotalArrayDimensions("User"), 0);
    assertEquals(getTotalArrayDimensions("object"), 0);
  });

  await t.step("should return correct dimensions for array types", () => {
    assertEquals(
      getTotalArrayDimensions({ baseType: "string", dimensions: 1 }),
      1,
    );
    assertEquals(
      getTotalArrayDimensions({ baseType: "User", dimensions: 2 }),
      2,
    );
  });

  await t.step("should sum dimensions for nested array types", () => {
    const nestedType: ArrayType = {
      baseType: { baseType: "string", dimensions: 2 },
      dimensions: 3,
    };
    assertEquals(getTotalArrayDimensions(nestedType), 5);
  });
});

Deno.test("flattenArrayType", async (t) => {
  await t.step("should return non-array types as is", () => {
    assertEquals(flattenArrayType("string"), "string");
    assertEquals(flattenArrayType("User"), "User");
    assertEquals(flattenArrayType("object"), "object");
  });

  await t.step("should flatten simple array types", () => {
    assertEquals(flattenArrayType({ baseType: "string", dimensions: 1 }), {
      baseType: "string",
      dimensions: 1,
    });
  });

  await t.step("should flatten nested array types", () => {
    const nestedType: ArrayType = {
      baseType: { baseType: "string", dimensions: 2 },
      dimensions: 3,
    };
    assertEquals(flattenArrayType(nestedType), {
      baseType: "string",
      dimensions: 5,
    });
  });
});

Deno.test("Complex type scenarios", async (t) => {
  await t.step("should handle deeply nested array types", () => {
    const complexType = "User[][][][][]";
    const parsed = parseFieldType(complexType);
    assertEquals(parsed, { baseType: "User", dimensions: 5 });
    assertEquals(fieldTypeToString(parsed as FieldType), complexType);
  });

  await t.step("should handle detailed fields with nested array types", () => {
    const detailedField: DetailedField = {
      type: "object",
      desc: "Complex object",
      fields: {
        matrix: { type: { baseType: "number", dimensions: 2 } },
        users: { type: { baseType: "User", dimensions: 1 } },
        metadata: {
          type: "object",
          fields: {
            tags: { type: { baseType: "string", dimensions: 1 } },
          },
        },
      },
    };

    const parsed = parseDetailedField(detailedField);
    assertEquals(parsed.type, "object");
    assertEquals(parsed.fields?.metadata?.type, "object");
    assertEquals(
      (parsed.fields?.metadata?.fields?.tags?.type as ArrayType)?.baseType,
      "string",
    );
  });

  await t.step(
    "should maintain type consistency through parse-stringify-parse cycle",
    () => {
      const originalType: FieldType = { baseType: "User", dimensions: 2 };
      const stringified = fieldTypeToString(originalType);
      const reparsed = parseFieldType(stringified);
      assertEquals(reparsed, originalType);
    },
  );
});
