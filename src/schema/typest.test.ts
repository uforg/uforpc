import { assertEquals } from "@std/assert";
import {
  type Field,
  getBaseType,
  isArrayType,
  isCustomType,
  isDetailedField,
  isPrimitiveType,
  isValidType,
} from "./types.ts";

Deno.test("isDetailedField", async (t) => {
  await t.step("should identify simple fields", () => {
    const simpleField: Field = "string";
    assertEquals(isDetailedField(simpleField), false);
  });

  await t.step("should identify detailed fields", () => {
    const detailedField: Field = {
      type: "string",
      desc: "A test field",
    };
    assertEquals(isDetailedField(detailedField), true);
  });

  await t.step("should identify detailed fields with nested fields", () => {
    const nestedField: Field = {
      type: "object",
      fields: {
        subField: "string",
      },
    };
    assertEquals(isDetailedField(nestedField), true);
  });
});

Deno.test("isPrimitiveType", async (t) => {
  await t.step("should identify all primitive types", () => {
    assertEquals(isPrimitiveType("string"), true);
    assertEquals(isPrimitiveType("number"), true);
    assertEquals(isPrimitiveType("float"), true);
    assertEquals(isPrimitiveType("boolean"), true);
  });

  await t.step("should reject non-primitive types", () => {
    assertEquals(isPrimitiveType("object"), false);
    assertEquals(isPrimitiveType("User"), false);
    assertEquals(isPrimitiveType("string[]"), false);
    assertEquals(isPrimitiveType(""), false);
  });
});

Deno.test("isArrayType", async (t) => {
  await t.step("should identify array types", () => {
    assertEquals(isArrayType("string[]"), true);
    assertEquals(isArrayType("User[]"), true);
    assertEquals(isArrayType("number[]"), true);
  });

  await t.step("should reject non-array types", () => {
    assertEquals(isArrayType("string"), false);
    assertEquals(isArrayType("User"), false);
    assertEquals(isArrayType("[]"), false);
    assertEquals(isArrayType("[string]"), false);
    assertEquals(isArrayType("array"), false);
  });

  await t.step("should handle edge cases", () => {
    assertEquals(isArrayType(""), false);
    assertEquals(isArrayType("string[][]"), true);
    assertEquals(isArrayType("[]string"), false);
  });
});

Deno.test("getBaseType", async (t) => {
  await t.step("should extract base type from array types", () => {
    assertEquals(getBaseType("string[]"), "string");
    assertEquals(getBaseType("User[]"), "User");
    assertEquals(getBaseType("number[]"), "number");
  });

  await t.step("should return non-array types as-is", () => {
    assertEquals(getBaseType("string"), "string");
    assertEquals(getBaseType("User"), "User");
    assertEquals(getBaseType("object"), "object");
  });

  await t.step("should handle edge cases", () => {
    assertEquals(getBaseType(""), "");
  });
});

Deno.test("isCustomType", async (t) => {
  await t.step("should identify valid custom types", () => {
    assertEquals(isCustomType("User"), true);
    assertEquals(isCustomType("UserProfile"), true);
    assertEquals(isCustomType("ABC"), true);
    assertEquals(isCustomType("A"), true);
  });

  await t.step("should reject invalid custom types", () => {
    assertEquals(isCustomType("user"), false);
    assertEquals(isCustomType("123User"), false);
    assertEquals(isCustomType("_User"), false);
    assertEquals(isCustomType(""), false);
  });

  await t.step("should handle array types", () => {
    assertEquals(isCustomType("User[]"), true);
    assertEquals(isCustomType("user[]"), false);
  });

  await t.step("should reject primitive types", () => {
    assertEquals(isCustomType("string"), false);
    assertEquals(isCustomType("number"), false);
    assertEquals(isCustomType("boolean"), false);
    assertEquals(isCustomType("float"), false);
  });
});

Deno.test("isValidType", async (t) => {
  await t.step("should validate primitive types", () => {
    assertEquals(isValidType("string"), true);
    assertEquals(isValidType("number"), true);
    assertEquals(isValidType("float"), true);
    assertEquals(isValidType("boolean"), true);
  });

  await t.step("should validate custom types", () => {
    assertEquals(isValidType("User"), true);
    assertEquals(isValidType("UserProfile"), true);
  });

  await t.step("should validate array types", () => {
    assertEquals(isValidType("string[]"), true);
    assertEquals(isValidType("User[]"), true);
    assertEquals(isValidType("number[]"), true);
  });

  await t.step("should validate object type", () => {
    assertEquals(isValidType("object"), true);
  });

  await t.step("should reject invalid types", () => {
    assertEquals(isValidType(""), false);
    assertEquals(isValidType("invalid"), false);
    assertEquals(isValidType("123Type"), false);
    assertEquals(isValidType("_Invalid"), false);
  });

  await t.step("should handle edge cases", () => {
    assertEquals(isValidType("[]"), false);
    assertEquals(isValidType("string[][]"), false);
    assertEquals(isValidType("Object"), true);
    assertEquals(isValidType("CONSTANT"), true);
  });
});

Deno.test("type combinations", async (t) => {
  await t.step("should validate array of custom types", () => {
    const type = "User[]";
    assertEquals(isArrayType(type), true);
    assertEquals(isCustomType(getBaseType(type)), true);
    assertEquals(isValidType(type), true);
  });

  await t.step("should validate array of primitive types", () => {
    const type = "string[]";
    assertEquals(isArrayType(type), true);
    assertEquals(isPrimitiveType(getBaseType(type)), true);
    assertEquals(isValidType(type), true);
  });

  await t.step("should handle complex scenarios", () => {
    // Custom type that looks like a primitive
    assertEquals(isValidType("String"), true);
    assertEquals(isPrimitiveType("String"), false);
    assertEquals(isCustomType("String"), true);

    // Array of objects
    assertEquals(isValidType("object[]"), true);
    assertEquals(isArrayType("object[]"), true);
    assertEquals(getBaseType("object[]"), "object");
  });
});
