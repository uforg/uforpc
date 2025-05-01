// deno-lint-ignore-file no-explicit-any

import { assertEquals } from "@std/assert";
import { setAtPath } from "./setAtPath.ts";

Deno.test("setAtPath - modifies a simple property", () => {
  const original = { name: "John", age: 30 };
  const result = setAtPath(original, "name", "Jane");
  assertEquals(result, { name: "Jane", age: 30 });
  // Original should remain unchanged
  assertEquals(original, { name: "John", age: 30 });
});

Deno.test("setAtPath - creates nested property if it doesn't exist", () => {
  const original: any = { user: { name: "John" } };
  const result = setAtPath(original, "user.age", 25);
  assertEquals(result, { user: { name: "John", age: 25 } });
});

Deno.test("setAtPath - modifies nested property", () => {
  const original = { user: { name: "John", details: { age: 30 } } };
  const result = setAtPath(original, "user.details.age", 31);
  assertEquals(result, { user: { name: "John", details: { age: 31 } } });
});

Deno.test("setAtPath - modifies array element", () => {
  const original = { users: ["John", "Jane", "Bob"] };
  const result = setAtPath(original, "users.1", "Sarah");
  assertEquals(result, { users: ["John", "Sarah", "Bob"] });
});

Deno.test("setAtPath - modifies property in array of objects", () => {
  const original = { users: [{ name: "John" }, { name: "Jane" }] };
  const result = setAtPath(original, "users.0.name", "Robert");
  assertEquals(result, { users: [{ name: "Robert" }, { name: "Jane" }] });
});

Deno.test("setAtPath - handles deeply nested array properties", () => {
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
  const result = setAtPath(
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

Deno.test("setAtPath - adds property to empty object", () => {
  const original = {};
  const result = setAtPath(original, "newProp", "value");
  assertEquals(result, { newProp: "value" });
});

Deno.test("setAtPath - creates full path with nested objects", () => {
  const original = {};
  const result = setAtPath(original, "a.b.c", "value");
  assertEquals(result, { a: { b: { c: "value" } } });
});

Deno.test("setAtPath - creates arrays when numeric keys are used", () => {
  const original = {};
  const result = setAtPath(original, "items.0", "first");
  assertEquals(result, { items: { "0": "first" } });
});

Deno.test("setAtPath - handles null values", () => {
  const original: any = { user: null };
  const result = setAtPath(original, "user.name", "John");
  assertEquals(result, { user: { name: "John" } });
});

Deno.test("setAtPath - overwrites primitive with object", () => {
  const original: any = { count: 5 };
  const result = setAtPath(original, "count.value", 10);
  assertEquals(result, { count: { value: 10 } });
});

Deno.test("setAtPath - handles empty string path", () => {
  const original: any = { name: "John" };
  const result = setAtPath(original, "", { replaced: true });
  assertEquals(result, { "": { replaced: true }, name: "John" });
});

Deno.test("setAtPath - sets falsy values correctly", () => {
  const original = { active: true, count: 1 };
  const result1 = setAtPath(original, "active", false);
  const result2 = setAtPath(original, "count", 0);
  assertEquals(result1, { active: false, count: 1 });
  assertEquals(result2, { active: true, count: 0 });
});

Deno.test("setAtPath - preserves object type", () => {
  class User {
    name: string;
    age: number;
    constructor(name: string, age: number) {
      this.name = name;
      this.age = age;
    }
  }

  const original = new User("John", 30);
  const result = setAtPath(original, "age", 31);
  assertEquals(result.age, 31);
  assertEquals(result.name, "John");
});

Deno.test("setAtPath - works with array at root", () => {
  const original = [1, 2, 3];
  const result = setAtPath(original, "1", 42);
  assertEquals(result, [1, 42, 3]);
});

Deno.test("setAtPath - preserves everything else in complex objects", () => {
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

  const result = setAtPath(original, "user.address.city", "Boston");

  assertEquals(result.user.name, "John");
  assertEquals(result.user.address.city, "Boston");
  assertEquals(result.user.address.zip, "10001");
  assertEquals(result.user.hobbies, ["reading", "sports"]);
  assertEquals(result.settings.theme, "dark");
  assertEquals(result.settings.notifications, true);
});

Deno.test("setAtPath - Sets null and undefined values correctly", () => {
  const original: any = { name: "John", age: 30 };
  let result = setAtPath(original, "name", null);
  result = setAtPath(result, "age", undefined);

  assertEquals(result, { name: null, age: undefined });
  assertEquals(original, { name: "John", age: 30 });
});

Deno.test("setAtPath - removeNullOrUndefined option removes null properties", () => {
  const original: any = {
    name: "John",
    age: 30,
    address: { city: "New York", zip: "10001" },
  };

  // Setting a property to null should remove it when removeNullOrUndefined is true
  const result = setAtPath(original, "name", null, {
    removeNullOrUndefined: true,
  });

  assertEquals(result, {
    age: 30,
    address: { city: "New York", zip: "10001" },
  });
  // Original should be unchanged
  assertEquals(original, {
    name: "John",
    age: 30,
    address: { city: "New York", zip: "10001" },
  });
});

Deno.test("setAtPath - removeNullOrUndefined option removes undefined properties", () => {
  const original: any = { name: "John", age: 30 };

  // Setting a property to undefined should remove it when removeNullOrUndefined is true
  const result = setAtPath(original, "age", undefined, {
    removeNullOrUndefined: true,
  });

  assertEquals(result, { name: "John" });
  assertEquals(original, { name: "John", age: 30 });
});

Deno.test("setAtPath - removeNullOrUndefined works with nested properties", () => {
  const original: any = {
    user: {
      name: "John",
      address: {
        city: "New York",
        zip: "10001",
        country: "USA",
      },
    },
  };

  // Remove a nested property
  const result = setAtPath(
    original,
    "user.address.zip",
    null,
    { removeNullOrUndefined: true },
  );

  assertEquals(result, {
    user: {
      name: "John",
      address: {
        city: "New York",
        country: "USA",
      },
    },
  });
});

Deno.test("setAtPath - removeNullOrUndefined works with arrays", () => {
  const original: any = { users: ["John", "Jane", "Bob"] };

  // Setting an array element to null should remove the element when removeNullOrUndefined is true
  const result = setAtPath(original, "users.1", null, {
    removeNullOrUndefined: true,
  });

  // The array should have the element removed and subsequent elements shifted
  assertEquals(result, { users: ["John", "Bob"] });
});

Deno.test("setAtPath - removeNullOrUndefined doesn't affect non-null values", () => {
  const original = { name: "John", age: 30 };

  // Regular values should still be set normally
  const result = setAtPath(original, "name", "Jane", {
    removeNullOrUndefined: true,
  });

  assertEquals(result, { name: "Jane", age: 30 });
});

Deno.test("setAtPath - removeNullOrUndefined works with complex nested arrays", () => {
  const original: any = {
    departments: [
      {
        name: "Engineering",
        teams: [
          {
            name: "Frontend",
            members: [
              { id: 1, name: "Alice" },
              { id: 2, name: "Bob" },
              { id: 3, name: "Charlie" },
            ],
          },
          {
            name: "Backend",
            members: [
              { id: 4, name: "Dave" },
              { id: 5, name: "Eve" },
            ],
          },
        ],
      },
      {
        name: "Marketing",
        teams: [
          {
            name: "Social Media",
            members: [
              { id: 6, name: "Frank" },
            ],
          },
        ],
      },
    ],
  };

  // Remove a deeply nested array element
  const result1 = setAtPath(
    original,
    "departments.0.teams.0.members.1",
    null,
    { removeNullOrUndefined: true },
  );

  // Bob should be removed, Charlie should shift to index 1
  assertEquals(result1.departments[0].teams[0].members.length, 2);
  assertEquals(result1.departments[0].teams[0].members[0].name, "Alice");
  assertEquals(result1.departments[0].teams[0].members[1].name, "Charlie");

  // Remove an entire team
  const result2 = setAtPath(
    original,
    "departments.0.teams.0",
    null,
    { removeNullOrUndefined: true },
  );

  // Backend team should be removed
  assertEquals(result2.departments[0].teams.length, 1);
  assertEquals(result2.departments[0].teams[0].name, "Backend");

  // Remove an entire department
  const result3 = setAtPath(
    original,
    "departments.1",
    null,
    { removeNullOrUndefined: true },
  );

  // Marketing department should be removed
  assertEquals(result3.departments.length, 1);
  assertEquals(result3.departments[0].name, "Engineering");
});

Deno.test("setAtPath - removeNullOrUndefined removes entire nested subtree", () => {
  const original: any = {
    project: {
      name: "Example",
      details: {
        stages: [
          {
            phase: "Planning",
            tasks: [
              { id: 1, description: "Research" },
              { id: 2, description: "Design" },
            ],
          },
          {
            phase: "Implementation",
            tasks: [
              { id: 3, description: "Coding" },
            ],
          },
        ],
        budget: {
          allocated: 10000,
          spent: 5000,
        },
      },
    },
  };

  // Remove the entire details object
  const result = setAtPath(
    original,
    "project.details",
    null,
    { removeNullOrUndefined: true },
  );

  // The entire subtree should be gone
  assertEquals(result, {
    project: {
      name: "Example",
    },
  });

  // Original should remain unchanged
  assertEquals(original.project.details.stages.length, 2);
  assertEquals(original.project.details.budget.allocated, 10000);
});
