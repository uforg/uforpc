import { describe, expect, it } from "vitest";

import { extractNodeFromSchema } from "./extractNodeFromSchema.ts";

describe("extractNodeFromSchema", () => {
  it("should return null if node is not found", () => {
    const schema = `
    type User {
      name: string
    }
  `;

    const result = extractNodeFromSchema(schema, "type", "NonExistent");
    expect(result).toBe(null);
  });

  it("should extract a simple type node", () => {
    const schema = `
    type User {
      name: string
    }
  `;

    const expected = `    type User {
      name: string
    }`;

    const result = extractNodeFromSchema(schema, "type", "User");
    expect(result).toBe(expected);
  });

  it("should extract a proc node", () => {
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
    expect(result).toBe(expected);
  });

  it("should handle nested braces correctly", () => {
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
    expect(result).toBe(expected);
  });

  it("should not be confused by similar node names", () => {
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
    expect(result).toBe(expected);
  });

  it("should handle nodes with braces on the same line", () => {
    const schema = `
    type Empty {}
    
    type User {
      name: string
    }
  `;

    const expected = `    type Empty {}`;

    const result = extractNodeFromSchema(schema, "type", "Empty");
    expect(result).toBe(expected);
  });

  it("should handle complex schema with multiple node types", () => {
    const schema = `
    version 1
    
  
    
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
    expect(result).toBe(expected);
  });

  it("should extract deprecated node with message on separate line", () => {
    const schema = `
    deprecated("This is deprecated")
    type One {
      two: Two
    }
    
    type User {
      name: string
    }
  `;

    const expected = `    deprecated("This is deprecated")
    type One {
      two: Two
    }`;

    const result = extractNodeFromSchema(schema, "type", "One");
    expect(result).toBe(expected);
  });

  it("should extract deprecated node with keyword on same line", () => {
    const schema = `
    deprecated type OldUser {
      name: string
      email: string
    }
    
    type User {
      name: string
    }
  `;

    const expected = `    deprecated type OldUser {
      name: string
      email: string
    }`;

    const result = extractNodeFromSchema(schema, "type", "OldUser");
    expect(result).toBe(expected);
  });

  it("should extract deprecated proc with message", () => {
    const schema = `
    deprecated("Use NewGetUser instead")
    proc GetUser {
      input {
        id: string
      }
      output {
        user: User
      }
    }
  `;

    const expected = `    deprecated("Use NewGetUser instead")
    proc GetUser {
      input {
        id: string
      }
      output {
        user: User
      }
    }`;

    const result = extractNodeFromSchema(schema, "proc", "GetUser");
    expect(result).toBe(expected);
  });

  it("should extract deprecated stream with same line format", () => {
    const schema = `
    deprecated stream OldEvents {
      event: Event
    }
  `;

    const expected = `    deprecated stream OldEvents {
      event: Event
    }`;

    const result = extractNodeFromSchema(schema, "stream", "OldEvents");
    expect(result).toBe(expected);
  });

  it("should handle deprecated with empty braces", () => {
    const schema = `
    deprecated("No longer used")
    type Empty {}
  `;

    const expected = `    deprecated("No longer used")
    type Empty {}`;

    const result = extractNodeFromSchema(schema, "type", "Empty");
    expect(result).toBe(expected);
  });

  it("should not confuse deprecated with regular nodes", () => {
    const schema = `
    deprecated("Old version")
    type UserV1 {
      name: string
    }
    
    type User {
      name: string
      email: string
    }
  `;

    const expected = `    type User {
      name: string
      email: string
    }`;

    const result = extractNodeFromSchema(schema, "type", "User");
    expect(result).toBe(expected);
  });
});
