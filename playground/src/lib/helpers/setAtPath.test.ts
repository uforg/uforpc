import { describe, expect, it } from "vitest";

import { setAtPath } from "./setAtPath.ts";

describe("setAtPath", () => {
  it("modifies a simple property", () => {
    const original = { name: "John", age: 30 };
    const result = setAtPath(original, "name", "Jane");
    expect(result).toEqual({ name: "Jane", age: 30 });
    // Original should remain unchanged
    expect(original).toEqual({ name: "John", age: 30 });
  });

  it("creates nested property if it doesn't exist", () => {
    const original: any = { user: { name: "John" } };
    const result = setAtPath(original, "user.age", 25);
    expect(result).toEqual({ user: { name: "John", age: 25 } });
  });

  it("modifies nested property", () => {
    const original = { user: { name: "John", details: { age: 30 } } };
    const result = setAtPath(original, "user.details.age", 31);
    expect(result).toEqual({ user: { name: "John", details: { age: 31 } } });
  });

  it("modifies array element", () => {
    const original = { users: ["John", "Jane", "Bob"] };
    const result = setAtPath(original, "users.1", "Sarah");
    expect(result).toEqual({ users: ["John", "Sarah", "Bob"] });
  });

  it("modifies property in array of objects", () => {
    const original = { users: [{ name: "John" }, { name: "Jane" }] };
    const result = setAtPath(original, "users.0.name", "Robert");
    expect(result).toEqual({ users: [{ name: "Robert" }, { name: "Jane" }] });
  });

  it("handles deeply nested array properties", () => {
    const original = {
      departments: [
        {
          name: "Engineering",
          teams: [{ name: "Frontend", members: ["Alice", "Bob"] }],
        },
      ],
    };
    const result = setAtPath(
      original,
      "departments.0.teams.0.members.1",
      "Charlie",
    );
    expect(result).toEqual({
      departments: [
        {
          name: "Engineering",
          teams: [{ name: "Frontend", members: ["Alice", "Charlie"] }],
        },
      ],
    });
  });

  it("adds property to empty object", () => {
    const original = {};
    const result = setAtPath(original, "newProp", "value");
    expect(result).toEqual({ newProp: "value" });
  });

  it("creates full path with nested objects", () => {
    const original = {};
    const result = setAtPath(original, "a.b.c", "value");
    expect(result).toEqual({ a: { b: { c: "value" } } });
  });

  it("creates arrays when numeric keys are used in path", () => {
    const original = {};
    const result = setAtPath(original, "items.0", "first");
    // Should create an array, not an object with a "0" key
    expect(result).toEqual({ items: ["first"] });
  });

  it("handles null values", () => {
    const original: any = { user: null };
    const result = setAtPath(original, "user.name", "John");
    expect(result).toEqual({ user: { name: "John" } });
  });

  it("overwrites primitive with object", () => {
    const original: any = { count: 5 };
    const result = setAtPath(original, "count.value", 10);
    expect(result).toEqual({ count: { value: 10 } });
  });

  it("handles empty string path", () => {
    const original: any = { name: "John" };
    const result = setAtPath(original, "", { replaced: true });
    expect(result).toEqual({ "": { replaced: true }, name: "John" });
  });

  it("sets falsy values correctly", () => {
    const original = { active: true, count: 1 };
    const result1 = setAtPath(original, "active", false);
    const result2 = setAtPath(original, "count", 0);
    expect(result1).toEqual({ active: false, count: 1 });
    expect(result2).toEqual({ active: true, count: 0 });
  });

  it("preserves object type", () => {
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
    expect(result.age).toEqual(31);
    expect(result.name).toEqual("John");
  });

  it("works with array at root", () => {
    const original = [1, 2, 3];
    const result = setAtPath(original, "1", 42);
    expect(result).toEqual([1, 42, 3]);
  });

  it("preserves everything else in complex objects", () => {
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

    expect(result.user.name).toEqual("John");
    expect(result.user.address.city).toEqual("Boston");
    expect(result.user.address.zip).toEqual("10001");
    expect(result.user.hobbies).toEqual(["reading", "sports"]);
    expect(result.settings.theme).toEqual("dark");
    expect(result.settings.notifications).toEqual(true);
  });

  it("removes null values by default", () => {
    const original: any = {
      name: "John",
      age: 30,
      address: { city: "New York", zip: "10001" },
    };

    // Setting a property to null should remove it
    const result = setAtPath(original, "name", null);

    expect(result).toEqual({
      age: 30,
      address: { city: "New York", zip: "10001" },
    });
    // Original should be unchanged
    expect(original).toEqual({
      name: "John",
      age: 30,
      address: { city: "New York", zip: "10001" },
    });
  });

  it("removes undefined values by default", () => {
    const original: any = { name: "John", age: 30 };

    // Setting a property to undefined should remove it
    const result = setAtPath(original, "age", undefined);

    expect(result).toEqual({ name: "John" });
    expect(original).toEqual({ name: "John", age: 30 });
  });

  it("removes null values in nested properties", () => {
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

    expect(result).toEqual({
      user: {
        name: "John",
        address: {
          city: "New York",
          country: "USA",
        },
      },
    });
  });

  it("removes null values in arrays", () => {
    const original: any = { users: ["John", "Jane", "Bob"] };

    // Setting an array element to null should remove the element
    const result = setAtPath(original, "users.1", null);

    // The array should have the element removed and subsequent elements shifted
    expect(result).toEqual({ users: ["John", "Bob"] });
  });

  it("creates and populates nested arrays", () => {
    const original = {};

    // Create a new array with elements
    const result = setAtPath(original, "user.followers.0.name", "Alice");

    // Should create an array for "followers", not an object with a "0" key
    expect(result).toEqual({
      user: {
        followers: [{ name: "Alice" }],
      },
    });
  });

  it("creates arrays with multiple items", () => {
    const original = {};

    // Set first entry
    let result: any = setAtPath(original, "items.0", "first");

    // Set third entry (index 2)
    result = setAtPath(result, "items.2", "third");

    // The second element should be undefined (not an empty item)
    expect(result.items.length).toEqual(3);
    expect(result.items[0]).toEqual("first");
    expect(result.items[1]).toEqual(undefined);
    expect(result.items[2]).toEqual("third");
  });

  it("creates complex nested arrays", () => {
    const original = {};

    // Create complex nested structure
    const result: any = setAtPath(
      original,
      "users.0.posts.0.comments.1.author",
      "Alice",
    );

    // Verify the structure instead of making a complete comparison
    expect(result.users.length).toEqual(1);
    expect(result.users[0].posts.length).toEqual(1);
    expect(result.users[0].posts[0].comments.length).toEqual(2);
    expect(result.users[0].posts[0].comments[0]).toEqual(undefined);
    expect(result.users[0].posts[0].comments[1].author).toEqual("Alice");
  });

  it("correctly handles mixed array and object paths", () => {
    const original = {};

    // Create a path with both array indices and object properties
    const result: any = setAtPath(
      original,
      "departments.0.teams.1.members.0.skills.2",
      "JavaScript",
    );

    // Verify the structure instead of making a complete comparison
    expect(result.departments.length).toEqual(1);
    expect(result.departments[0].teams.length).toEqual(2);
    expect(result.departments[0].teams[0]).toEqual(undefined);
    expect(result.departments[0].teams[1].members.length).toEqual(1);
    expect(result.departments[0].teams[1].members[0].skills.length).toEqual(3);
    expect(result.departments[0].teams[1].members[0].skills[0]).toEqual(
      undefined,
    );
    expect(result.departments[0].teams[1].members[0].skills[1]).toEqual(
      undefined,
    );
    expect(result.departments[0].teams[1].members[0].skills[2]).toEqual(
      "JavaScript",
    );
  });
});
