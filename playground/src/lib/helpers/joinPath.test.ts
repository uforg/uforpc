import { describe, expect, it } from "vitest";

import { joinPath } from "./joinPath.ts";

describe("joinPath", () => {
  it("should handle empty array", () => {
    const parts: string[] = [];
    const result = joinPath(parts);
    expect(result).toBe("");
  });

  it("should handle single path part", () => {
    const parts = ["home"];
    const result = joinPath(parts);
    expect(result).toBe("home");
  });

  it("should join multiple path parts with separator", () => {
    const parts = ["home", "user", "documents"];
    const result = joinPath(parts);
    expect(result).toBe("home/user/documents");
  });

  it("should remove duplicate separators between parts", () => {
    const parts = ["home/", "/user", "documents"];
    const result = joinPath(parts);
    expect(result).toBe("home/user/documents");
  });

  it("should handle parts with leading separators", () => {
    const parts = ["/home", "/user", "/documents"];
    const result = joinPath(parts);
    expect(result).toBe("/home/user/documents");
  });

  it("should handle parts with trailing separators", () => {
    const parts = ["home/", "user/", "documents/"];
    const result = joinPath(parts);
    expect(result).toBe("home/user/documents/");
  });

  it("should handle empty strings in the array", () => {
    const parts = ["home", "", "user", "", "documents"];
    const result = joinPath(parts);
    expect(result).toBe("home/user/documents");
  });

  it("should handle multiple consecutive separators", () => {
    const parts = ["home//", "//user///", "documents"];
    const result = joinPath(parts);
    expect(result).toBe("home/user/documents");
  });

  it("should handle all empty strings", () => {
    const parts = ["", "", ""];
    const result = joinPath(parts);
    expect(result).toBe("/");
  });

  it("should handle single separator parts", () => {
    const parts = ["home", "/", "user", "/", "documents"];
    const result = joinPath(parts);
    expect(result).toBe("home/user/documents");
  });

  it("should handle absolute path starting with separator", () => {
    const parts = ["/", "home", "user", "documents"];
    const result = joinPath(parts);
    expect(result).toBe("/home/user/documents");
  });

  it("should handle complex mixed separators", () => {
    const parts = ["/home///", "//user//", "/documents///"];
    const result = joinPath(parts);
    expect(result).toBe("/home/user/documents/");
  });

  it("should handle path with file extension", () => {
    const parts = ["home", "user", "documents", "file.txt"];
    const result = joinPath(parts);
    expect(result).toBe("home/user/documents/file.txt");
  });

  it("should handle relative path parts", () => {
    const parts = [".", "home", "..", "user", "documents"];
    const result = joinPath(parts);
    expect(result).toBe("./home/../user/documents");
  });

  it("should handle parts with spaces", () => {
    const parts = ["my folder", "sub folder", "file name.txt"];
    const result = joinPath(parts);
    expect(result).toBe("my folder/sub folder/file name.txt");
  });

  it("should handle only separators", () => {
    const parts = ["/", "//", "///"];
    const result = joinPath(parts);
    expect(result).toBe("/");
  });

  it("should handle mixed content with multiple separators", () => {
    const parts = ["api///v1//", "//users///", "/123//"];
    const result = joinPath(parts);
    expect(result).toBe("api/v1/users/123/");
  });

  it("should handle single character parts", () => {
    const parts = ["a", "b", "c", "d"];
    const result = joinPath(parts);
    expect(result).toBe("a/b/c/d");
  });

  it("should preserve query parameters and fragments", () => {
    const parts = ["api", "users", "123?name=john&age=30#profile"];
    const result = joinPath(parts);
    expect(result).toBe("api/users/123?name=john&age=30#profile");
  });

  // URL protocol tests
  it("should handle HTTP URLs", () => {
    const parts = ["http://example.com", "api", "v1", "users"];
    const result = joinPath(parts);
    expect(result).toBe("http://example.com/api/v1/users");
  });

  it("should handle HTTPS URLs", () => {
    const parts = ["https://api.example.com", "v2", "users", "123"];
    const result = joinPath(parts);
    expect(result).toBe("https://api.example.com/v2/users/123");
  });

  it("should handle WebSocket URLs", () => {
    const parts = ["ws://localhost:8080", "chat", "room", "123"];
    const result = joinPath(parts);
    expect(result).toBe("ws://localhost:8080/chat/room/123");
  });

  it("should handle secure WebSocket URLs", () => {
    const parts = ["wss://secure.example.com", "api", "websocket"];
    const result = joinPath(parts);
    expect(result).toBe("wss://secure.example.com/api/websocket");
  });

  it("should handle FTP URLs", () => {
    const parts = [
      "ftp://files.example.com",
      "public",
      "documents",
      "file.pdf",
    ];
    const result = joinPath(parts);
    expect(result).toBe("ftp://files.example.com/public/documents/file.pdf");
  });

  it("should handle URLs with ports", () => {
    const parts = ["http://localhost:3000", "api", "health"];
    const result = joinPath(parts);
    expect(result).toBe("http://localhost:3000/api/health");
  });

  it("should handle URLs with paths and trailing slashes", () => {
    const parts = ["https://example.com/base/", "/api/", "/users/"];
    const result = joinPath(parts);
    expect(result).toBe("https://example.com/base/api/users/");
  });

  it("should handle complex URLs with multiple separators", () => {
    const parts = ["https://api.example.com///", "//v1//", "users///", "123"];
    const result = joinPath(parts);
    expect(result).toBe("https://api.example.com/v1/users/123");
  });

  it("should handle URLs with query parameters", () => {
    const parts = ["https://example.com", "search?q=test", "results"];
    const result = joinPath(parts);
    expect(result).toBe("https://example.com/search?q=test/results");
  });

  it("should handle file protocol URLs", () => {
    const parts = ["file://", "home", "user", "documents", "file.txt"];
    const result = joinPath(parts);
    expect(result).toBe("file://home/user/documents/file.txt");
  });

  it("should handle localhost URLs with different ports", () => {
    const parts = ["http://127.0.0.1:8080", "admin", "dashboard"];
    const result = joinPath(parts);
    expect(result).toBe("http://127.0.0.1:8080/admin/dashboard");
  });

  it("should handle URLs with subdirectories", () => {
    const parts = ["https://cdn.example.com/assets/", "/images/", "logo.png"];
    const result = joinPath(parts);
    expect(result).toBe("https://cdn.example.com/assets/images/logo.png");
  });
});
