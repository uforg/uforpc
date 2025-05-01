// deno-lint-ignore-file no-explicit-any

import { assertEquals } from "@std/assert";
import { modifyObjectWithPath } from "./modifyObjectWithPath.ts";

Deno.test("modifyObjectWithPath - modifies a simple property", () => {
  const original = { name: "John", age: 30 };
  const result = modifyObjectWithPath(original, "name", "Jane");
  assertEquals(result, { name: "Jane", age: 30 });
  // Original should remain unchanged
  assertEquals(original, { name: "John", age: 30 });
});

Deno.test("modifyObjectWithPath - creates nested property if it doesn't exist", () => {
  const original: any = { user: { name: "John" } };
  const result = modifyObjectWithPath(original, "user.age", 25);
  assertEquals(result, { user: { name: "John", age: 25 } });
});

Deno.test("modifyObjectWithPath - modifies nested property", () => {
  const original = { user: { name: "John", details: { age: 30 } } };
  const result = modifyObjectWithPath(original, "user.details.age", 31);
  assertEquals(result, { user: { name: "John", details: { age: 31 } } });
});

Deno.test("modifyObjectWithPath - modifies array element", () => {
  const original = { users: ["John", "Jane", "Bob"] };
  const result = modifyObjectWithPath(original, "users.1", "Sarah");
  assertEquals(result, { users: ["John", "Sarah", "Bob"] });
});

Deno.test("modifyObjectWithPath - modifies property in array of objects", () => {
  const original = { users: [{ name: "John" }, { name: "Jane" }] };
  const result = modifyObjectWithPath(original, "users.0.name", "Robert");
  assertEquals(result, { users: [{ name: "Robert" }, { name: "Jane" }] });
});

Deno.test("modifyObjectWithPath - handles deeply nested array properties", () => {
  const original = {
    departments: [
      {
        name: "Engineering",
        teams: [
          { name: "Frontend", members: ["Alice", "Bob"] },
        ],
      },
    ],
  };
  const result = modifyObjectWithPath(
    original,
    "departments.0.teams.0.members.1",
    "Charlie",
  );
  assertEquals(result, {
    departments: [
      {
        name: "Engineering",
        teams: [
          { name: "Frontend", members: ["Alice", "Charlie"] },
        ],
      },
    ],
  });
});

Deno.test("modifyObjectWithPath - adds property to empty object", () => {
  const original = {};
  const result = modifyObjectWithPath(original, "newProp", "value");
  assertEquals(result, { newProp: "value" });
});

Deno.test("modifyObjectWithPath - creates full path with nested objects", () => {
  const original = {};
  const result = modifyObjectWithPath(original, "a.b.c", "value");
  assertEquals(result, { a: { b: { c: "value" } } });
});

Deno.test("modifyObjectWithPath - creates arrays when numeric keys are used", () => {
  const original = {};
  const result = modifyObjectWithPath(original, "items.0", "first");
  assertEquals(result, { items: { "0": "first" } });
});

Deno.test("modifyObjectWithPath - handles null values", () => {
  const original: any = { user: null };
  const result = modifyObjectWithPath(original, "user.name", "John");
  assertEquals(result, { user: { name: "John" } });
});

Deno.test("modifyObjectWithPath - overwrites primitive with object", () => {
  const original: any = { count: 5 };
  const result = modifyObjectWithPath(original, "count.value", 10);
  assertEquals(result, { count: { value: 10 } });
});

Deno.test("modifyObjectWithPath - handles empty string path", () => {
  const original: any = { name: "John" };
  const result = modifyObjectWithPath(original, "", { replaced: true });
  assertEquals(result, { "": { replaced: true }, name: "John" });
});

Deno.test("modifyObjectWithPath - sets falsy values correctly", () => {
  const original = { active: true, count: 1 };
  const result1 = modifyObjectWithPath(original, "active", false);
  const result2 = modifyObjectWithPath(original, "count", 0);
  assertEquals(result1, { active: false, count: 1 });
  assertEquals(result2, { active: true, count: 0 });
});

Deno.test("modifyObjectWithPath - preserves object type", () => {
  class User {
    name: string;
    age: number;
    constructor(name: string, age: number) {
      this.name = name;
      this.age = age;
    }
  }

  const original = new User("John", 30);
  const result = modifyObjectWithPath(original, "age", 31);
  assertEquals(result.age, 31);
  assertEquals(result.name, "John");
});

Deno.test("modifyObjectWithPath - works with array at root", () => {
  const original = [1, 2, 3];
  const result = modifyObjectWithPath(original, "1", 42);
  assertEquals(result, [1, 42, 3]);
});

Deno.test("modifyObjectWithPath - preserves everything else in complex objects", () => {
  const original = {
    user: {
      name: "John",
      address: {
        city: "New York",
        zip: "10001",
      },
      hobbies: ["reading", "sports"],
    },
    settings: {
      theme: "dark",
      notifications: true,
    },
  };

  const result = modifyObjectWithPath(original, "user.address.city", "Boston");

  assertEquals(result.user.name, "John");
  assertEquals(result.user.address.city, "Boston");
  assertEquals(result.user.address.zip, "10001");
  assertEquals(result.user.hobbies, ["reading", "sports"]);
  assertEquals(result.settings.theme, "dark");
  assertEquals(result.settings.notifications, true);
});

Deno.test("modifyObjectWithPath - Sets null and undefined values correctly", () => {
  const original: any = { name: "John", age: 30 };
  let result = modifyObjectWithPath(original, "name", null);
  result = modifyObjectWithPath(result, "age", undefined);

  assertEquals(result, { name: null, age: undefined });
  assertEquals(original, { name: "John", age: 30 });
});
