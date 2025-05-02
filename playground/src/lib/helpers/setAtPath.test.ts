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

Deno.test("setAtPath - creates arrays when numeric keys are used in path", () => {
  const original = {};
  const result = setAtPath(original, "items.0", "first");
  // Should create an array, not an object with a "0" key
  assertEquals(result, { items: ["first"] });
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

Deno.test("setAtPath - removes null values by default", () => {
  const original: any = {
    name: "John",
    age: 30,
    address: { city: "New York", zip: "10001" },
  };

  // Setting a property to null should remove it
  const result = setAtPath(original, "name", null);

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

Deno.test("setAtPath - removes undefined values by default", () => {
  const original: any = { name: "John", age: 30 };

  // Setting a property to undefined should remove it
  const result = setAtPath(original, "age", undefined);

  assertEquals(result, { name: "John" });
  assertEquals(original, { name: "John", age: 30 });
});

Deno.test("setAtPath - removes null values in nested properties", () => {
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
  const result = setAtPath(original, "user.address.zip", null);

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

Deno.test("setAtPath - removes null values in arrays", () => {
  const original: any = { users: ["John", "Jane", "Bob"] };

  // Setting an array element to null should remove the element
  const result = setAtPath(original, "users.1", null);

  // The array should have the element removed and subsequent elements shifted
  assertEquals(result, { users: ["John", "Bob"] });
});

Deno.test("setAtPath - creates and populates nested arrays", () => {
  const original = {};

  // Create a new array with elements
  const result = setAtPath(original, "user.followers.0.name", "Alice");

  // Should create an array for "followers", not an object with a "0" key
  assertEquals(result, {
    user: {
      followers: [
        { name: "Alice" },
      ],
    },
  });
});

Deno.test("setAtPath - creates arrays with multiple items", () => {
  const original = {};

  // Set first entry
  let result: any = setAtPath(original, "items.0", "first");

  // Set third entry (index 2)
  result = setAtPath(result, "items.2", "third");

  // The second element should be undefined (not an empty item)
  assertEquals(result.items.length, 3);
  assertEquals(result.items[0], "first");
  assertEquals(result.items[1], undefined);
  assertEquals(result.items[2], "third");
});

Deno.test("setAtPath - creates complex nested arrays", () => {
  const original = {};

  // Create complex nested structure
  const result: any = setAtPath(
    original,
    "users.0.posts.0.comments.1.author",
    "Alice",
  );

  // Verify the structure instead of making a complete comparison
  assertEquals(result.users.length, 1);
  assertEquals(result.users[0].posts.length, 1);
  assertEquals(result.users[0].posts[0].comments.length, 2);
  assertEquals(result.users[0].posts[0].comments[0], undefined);
  assertEquals(result.users[0].posts[0].comments[1].author, "Alice");
});

Deno.test("setAtPath - correctly handles mixed array and object paths", () => {
  const original = {};

  // Create a path with both array indices and object properties
  const result: any = setAtPath(
    original,
    "departments.0.teams.1.members.0.skills.2",
    "JavaScript",
  );

  // Verify the structure instead of making a complete comparison
  assertEquals(result.departments.length, 1);
  assertEquals(result.departments[0].teams.length, 2);
  assertEquals(result.departments[0].teams[0], undefined);
  assertEquals(result.departments[0].teams[1].members.length, 1);
  assertEquals(result.departments[0].teams[1].members[0].skills.length, 3);
  assertEquals(result.departments[0].teams[1].members[0].skills[0], undefined);
  assertEquals(result.departments[0].teams[1].members[0].skills[1], undefined);
  assertEquals(
    result.departments[0].teams[1].members[0].skills[2],
    "JavaScript",
  );
});
