import { describe, expect, it } from "vitest";
import { getMarkdownTitle } from "./getMarkdownTitle.ts";

describe("getMarkdownTitle", () => {
  it('should return "Untitled: first line" if no title is found', () => {
    const markdown = "This is the first line\nThis is the second line";
    const result = getMarkdownTitle(markdown);
    expect(result).toBe("Untitled: This is the first line");
  });

  it("should return the first level-1 header without the hash", () => {
    const markdown = "Some text\n# This is a title\nMore text";
    const result = getMarkdownTitle(markdown);
    expect(result).toBe("This is a title");
  });

  it("should return the first level-1 header and trim it", () => {
    const markdown = "Some text\n#    Title with spaces    \nMore text";
    const result = getMarkdownTitle(markdown);
    expect(result).toBe("Title with spaces");
  });

  it("should not return nested headers if a level-1 header exists", () => {
    const markdown = "Some text\n# Main Title\n## Nested Title\nMore text";
    const result = getMarkdownTitle(markdown);
    expect(result).toBe("Main Title");
  });

  it("should return first nested header if no level-1 header exists", () => {
    const markdown = "Some text\n### Nested Title\nMore text";
    const result = getMarkdownTitle(markdown);
    expect(result).toBe("Nested Title");
  });
});
