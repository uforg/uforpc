import { assertEquals } from "jsr:@std/assert";
import { getMarkdownTitle } from "./getMarkdownTitle.ts";

Deno.test("should return 'Untitled: first line' if no title is found", () => {
  const markdown = "This is the first line\nThis is the second line";
  const result = getMarkdownTitle(markdown);
  assertEquals(result, "Untitled: This is the first line");
});

Deno.test("should return the first level-1 header without the hash", () => {
  const markdown = "Some text\n# This is a title\nMore text";
  const result = getMarkdownTitle(markdown);
  assertEquals(result, "This is a title");
});

Deno.test("should return the first level-1 header and trim it", () => {
  const markdown = "Some text\n#    Title with spaces    \nMore text";
  const result = getMarkdownTitle(markdown);
  assertEquals(result, "Title with spaces");
});

Deno.test("should not return nested headers if a level-1 header exists", () => {
  const markdown = "Some text\n# Main Title\n## Nested Title\nMore text";
  const result = getMarkdownTitle(markdown);
  assertEquals(result, "Main Title");
});

Deno.test("should return first nested header if no level-1 header exists", () => {
  const markdown = "Some text\n### Nested Title\nMore text";
  const result = getMarkdownTitle(markdown);
  assertEquals(result, "Nested Title");
});
