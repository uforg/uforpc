import { assertEquals } from "@std/assert";
import {
  getBaseType,
  isCustomType,
  isPrimitiveType,
  isValidArrayType,
  isValidType,
} from "./types.ts";

Deno.test("isPrimitiveType", async (t) => {
  await t.step("should identify primitive types", () => {
    assertEquals(isPrimitiveType("string"), true);
    assertEquals(isPrimitiveType("number"), true);
    assertEquals(isPrimitiveType("float"), true);
    assertEquals(isPrimitiveType("boolean"), true);
  });

  await t.step("should reject non-primitive types", () => {
    assertEquals(isPrimitiveType("object"), false);
    assertEquals(isPrimitiveType("User"), false);
    assertEquals(isPrimitiveType("string[]"), false);
  });
});

Deno.test("isValidArrayType", async (t) => {
  await t.step("should validate array types", () => {
    assertEquals(isValidArrayType("string[]"), true);
    assertEquals(isValidArrayType("string[][]"), true);
    assertEquals(isValidArrayType("User[]"), true);
    assertEquals(isValidArrayType("User[][]"), true);
    assertEquals(isValidArrayType("User[][][][][][]"), true);
  });

  await t.step("should reject invalid array types", () => {
    assertEquals(isValidArrayType("[]"), false);
    assertEquals(isValidArrayType("[string]"), false);
    assertEquals(isValidArrayType("string["), false);
    assertEquals(isValidArrayType("string]"), false);
  });
});

Deno.test("getBaseType", async (t) => {
  await t.step("should extract base type", () => {
    assertEquals(getBaseType("string[]"), "string");
    assertEquals(getBaseType("string[][][][][][]"), "string");
    assertEquals(getBaseType("User[]"), "User");
    assertEquals(getBaseType("User[][]"), "User");
    assertEquals(getBaseType("object[]"), "object");
  });

  await t.step("should return non-array types as-is", () => {
    assertEquals(getBaseType("string"), "string");
    assertEquals(getBaseType("User"), "User");
    assertEquals(getBaseType("object"), "object");
  });
});

Deno.test("isCustomType", async (t) => {
  await t.step("should identify custom types", () => {
    assertEquals(isCustomType("User"), true);
    assertEquals(isCustomType("User[]"), true);
    assertEquals(isCustomType("User[][]"), true);
    assertEquals(isCustomType("UserProfile"), true);
  });

  await t.step("should reject invalid custom types", () => {
    assertEquals(isCustomType("user"), false);
    assertEquals(isCustomType("123User"), false);
    assertEquals(isCustomType("_User"), false);
    assertEquals(isCustomType(""), false);
    assertEquals(isCustomType("string"), false);
  });
});

Deno.test("isValidType", async (t) => {
  await t.step("should validate types", () => {
    assertEquals(isValidType("string"), true);
    assertEquals(isValidType("string[]"), true);
    assertEquals(isValidType("string[][]"), true);
    assertEquals(isValidType("User"), true);
    assertEquals(isValidType("User[]"), true);
    assertEquals(isValidType("User[][]"), true);
    assertEquals(isValidType("object"), true);
    assertEquals(isValidType("object[]"), true);
  });

  await t.step("should reject invalid types", () => {
    assertEquals(isValidType(""), false);
    assertEquals(isValidType("[]"), false);
    assertEquals(isValidType("[string]"), false);
    assertEquals(isValidType("string["), false);
    assertEquals(isValidType("string]"), false);
    assertEquals(isValidType("123Type"), false);
  });
});
