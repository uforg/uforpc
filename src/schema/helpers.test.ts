// deno-lint-ignore-file
import { assertEquals, assertThrows } from "@std/assert";
import {
  isArrayType,
  isCustomType,
  isObjectType,
  isPrimitiveType,
  isValidFieldType,
  parseArrayType,
} from "./helpers.ts";

Deno.test("isPrimitiveType", async (t) => {
  await t.step("should identify valid primitive types", () => {
    assertEquals(isPrimitiveType({ type: "string" }), true);
    assertEquals(isPrimitiveType({ type: "int" }), true);
    assertEquals(isPrimitiveType({ type: "float" }), true);
    assertEquals(isPrimitiveType({ type: "boolean" }), true);
  });

  await t.step("should reject non-primitive types", () => {
    assertEquals(isPrimitiveType({ type: "object" }), false);
    assertEquals(isPrimitiveType({ type: "User" }), false);
    assertEquals(isPrimitiveType({ type: "string[]" }), false);
  });

  await t.step("should handle edge cases", () => {
    assertEquals(isPrimitiveType({ type: "String" }), false);
    assertEquals(isPrimitiveType({ type: "INT" }), false);
    assertEquals(isPrimitiveType({ type: "" }), false);
    assertEquals(isPrimitiveType({ type: " string " }), false);
  });
});

Deno.test("isCustomType", async (t) => {
  await t.step("should identify valid custom types", () => {
    assertEquals(isCustomType({ type: "User" }), true);
    assertEquals(isCustomType({ type: "UserProfile" }), true);
    assertEquals(isCustomType({ type: "A" }), true);
    assertEquals(isCustomType({ type: "ABC123" }), true);
  });

  await t.step("should reject invalid custom types", () => {
    assertEquals(isCustomType({ type: "user" }), false);
    assertEquals(isCustomType({ type: "123User" }), false);
    assertEquals(isCustomType({ type: "User_Profile" }), false);
    assertEquals(isCustomType({ type: "User-Profile" }), false);
    assertEquals(isCustomType({ type: "" }), false);
  });

  await t.step("should reject primitive and array types", () => {
    assertEquals(isCustomType({ type: "string" }), false);
    assertEquals(isCustomType({ type: "object" }), false);
    assertEquals(isCustomType({ type: "User[]" }), false);
  });
});

Deno.test("isArrayType", async (t) => {
  await t.step("should identify valid array types", () => {
    assertEquals(isArrayType({ type: "string[]" }), true);
    assertEquals(isArrayType({ type: "User[][]" }), true);
    assertEquals(isArrayType({ type: "int[][]" }), true);
  });

  await t.step("should reject non-array types", () => {
    assertEquals(isArrayType({ type: "string" }), false);
    assertEquals(isArrayType({ type: "object" }), false);
    assertEquals(isArrayType({ type: "User" }), false);
  });
});

Deno.test("isObjectType", async (t) => {
  await t.step("should identify object type", () => {
    assertEquals(isObjectType({ type: "object" }), true);
  });

  await t.step("should reject non-object types", () => {
    assertEquals(isObjectType({ type: "string" }), false);
    assertEquals(isObjectType({ type: "User" }), false);
    assertEquals(isObjectType({ type: "string[]" }), false);
  });
});

Deno.test("isValidFieldType", async (t) => {
  await t.step("should validate primitive types", () => {
    assertEquals(isValidFieldType({ type: "string" }), true);
    assertEquals(isValidFieldType({ type: "int" }), true);
    assertEquals(isValidFieldType({ type: "float" }), true);
    assertEquals(isValidFieldType({ type: "boolean" }), true);
  });

  await t.step("should validate object and custom types", () => {
    assertEquals(isValidFieldType({ type: "object" }), true);
    assertEquals(isValidFieldType({ type: "User" }), true);
    assertEquals(isValidFieldType({ type: "UserProfile" }), true);
  });

  await t.step("should validate array types", () => {
    assertEquals(isValidFieldType({ type: "string[]" }), true);
    assertEquals(isValidFieldType({ type: "User[][]" }), true);
    assertEquals(isValidFieldType({ type: "int[][]" }), true);
  });

  await t.step("should reject invalid types", () => {
    assertEquals(isValidFieldType({ type: "invalid" }), false);
    assertEquals(isValidFieldType({ type: "invalid[]" }), false);
    assertEquals(isValidFieldType({} as any), false);
  });
});

Deno.test("parseArrayType", async (t) => {
  await t.step("should parse valid array type strings", () => {
    assertEquals(parseArrayType({ type: "string[]" }), {
      type: { type: "string" },
      dimensions: 1,
    });
    assertEquals(parseArrayType({ type: "int[][]" }), {
      type: { type: "int" },
      dimensions: 2,
    });
    assertEquals(parseArrayType({ type: "User[][][]" }), {
      type: { type: "User" },
      dimensions: 3,
    });
    assertEquals(parseArrayType({ type: "object[][]" }), {
      type: { type: "object" },
      dimensions: 2,
    });
  });

  await t.step("should parse nested array types", () => {
    const result = parseArrayType({ type: "string[][]" });
    assertEquals(result?.type.type, "string");
    assertEquals(result?.dimensions, 2);
  });

  await t.step("should return null for non-array types", () => {
    assertEquals(parseArrayType({ type: "string" }), {
      type: { type: "string" },
      dimensions: 0,
    });
    assertEquals(parseArrayType({ type: "User" }), {
      type: { type: "User" },
      dimensions: 0,
    });
    assertEquals(parseArrayType({ type: "object" }), {
      type: { type: "object" },
      dimensions: 0,
    });
  });
});
