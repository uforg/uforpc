import { assertEquals } from "@std/assert";
import { extractNodeFromSchema } from "./extractNodeFromSchema.ts";

Deno.test("should return null if node is not found", () => {
  const schema = `
    type User {
      name: string
    }
  `;

  const result = extractNodeFromSchema(schema, "type", "NonExistent");
  assertEquals(result, null);
});

Deno.test("should extract a simple type node", () => {
  const schema = `
    type User {
      name: string
    }
  `;

  const expected = `    type User {
      name: string
    }`;

  const result = extractNodeFromSchema(schema, "type", "User");
  assertEquals(result, expected);
});

Deno.test("should extract a rule node with @ prefix", () => {
  const schema = `
    rule @minLength {
      for: string
      param: int
      error: "String is too short"
    }
  `;

  const expected = `    rule @minLength {
      for: string
      param: int
      error: "String is too short"
    }`;

  const result = extractNodeFromSchema(schema, "rule", "minLength");
  assertEquals(result, expected);
});

Deno.test("should extract a proc node", () => {
  const schema = `
    proc GetUser {
      input {
        id: string
      }
      output {
        user: User
      }
    }
  `;

  const expected = `    proc GetUser {
      input {
        id: string
      }
      output {
        user: User
      }
    }`;

  const result = extractNodeFromSchema(schema, "proc", "GetUser");
  assertEquals(result, expected);
});

Deno.test("should handle nested braces correctly", () => {
  const schema = `
    type ComplexType {
      nested: {
        field1: string
        field2: {
          subfield: int
        }
      }
    }
  `;

  const expected = `    type ComplexType {
      nested: {
        field1: string
        field2: {
          subfield: int
        }
      }
    }`;

  const result = extractNodeFromSchema(schema, "type", "ComplexType");
  assertEquals(result, expected);
});

Deno.test("should not be confused by similar node names", () => {
  const schema = `
    type UserBase {
      id: string
    }
    
    type User {
      id: string
      name: string
    }
  `;

  const expected = `    type User {
      id: string
      name: string
    }`;

  const result = extractNodeFromSchema(schema, "type", "User");
  assertEquals(result, expected);
});

Deno.test("should handle nodes with braces on the same line", () => {
  const schema = `
    type Empty {}
    
    type User {
      name: string
    }
  `;

  const expected = `    type Empty {}`;

  const result = extractNodeFromSchema(schema, "type", "Empty");
  assertEquals(result, expected);
});

Deno.test("should handle complex schema with multiple node types", () => {
  const schema = `
    version 1
    
    rule @email {
      for: string
      error: "Invalid email format"
    }
    
    type User {
      email: string
        @email
      verified: boolean
    }
    
    proc VerifyUser {
      input {
        userId: string
      }
      output {
        success: boolean
      }
    }
  `;

  const expected = `    proc VerifyUser {
      input {
        userId: string
      }
      output {
        success: boolean
      }
    }`;

  const result = extractNodeFromSchema(schema, "proc", "VerifyUser");
  assertEquals(result, expected);
});
