import { assertEquals } from "@std/assert";
import { deleteMarkdownHeadings } from "./deleteMarkdownHeadings.ts";

Deno.test("deleteMarkdownHeadings - removes level-1 headings", () => {
  const markdown = "# Title\nDocument content.";
  const result = deleteMarkdownHeadings(markdown);
  assertEquals(result, "Document content.");
});

Deno.test("deleteMarkdownHeadings - keeps document unchanged if no headings exist", () => {
  const markdown = "Content without headings.";
  const result = deleteMarkdownHeadings(markdown);
  assertEquals(result, "Content without headings.");
});

Deno.test("deleteMarkdownHeadings - removes all level-1 headings", () => {
  const markdown = "# First title\nContent\n# Second title\nMore content";
  const result = deleteMarkdownHeadings(markdown);
  assertEquals(result, "Content\nMore content");
});

Deno.test("deleteMarkdownHeadings - preserves non-level-1 headings", () => {
  const markdown = "## Level 2 title\nDocument content.";
  const result = deleteMarkdownHeadings(markdown);
  assertEquals(result, "## Level 2 title\nDocument content.");
});

Deno.test("deleteMarkdownHeadings - removes level-1 headings but keeps other levels", () => {
  const markdown =
    "# Main title\n## Subtitle\nContent\n# Another title\nMore content";
  const result = deleteMarkdownHeadings(markdown);
  assertEquals(result, "## Subtitle\nContent\nMore content");
});

Deno.test("deleteMarkdownHeadings - handles empty documents", () => {
  const markdown = "";
  const result = deleteMarkdownHeadings(markdown);
  assertEquals(result, "");
});

Deno.test("deleteMarkdownHeadings - removes level-1 headings with preceding newlines and spaces", () => {
  const markdown =
    "Some initial content\n\n\n   # Title with preceding newlines and spaces\nFollowing content";
  const result = deleteMarkdownHeadings(markdown);
  assertEquals(result, "Some initial content\n\n\nFollowing content");
});
