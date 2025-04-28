import { assertEquals } from "@std/assert";
import { markdownToText } from "./markdownToText.ts";

Deno.test("markdownToText - converts basic markdown to plain text", async () => {
  const markdown = "# Title\nThis is a paragraph.";
  const result = await markdownToText(markdown);
  assertEquals(result, "Title\nThis is a paragraph.");
});

Deno.test("markdownToText - handles multiple paragraphs", async () => {
  const markdown = "First paragraph.\n\nSecond paragraph.";
  const result = await markdownToText(markdown);
  assertEquals(result, "First paragraph.\nSecond paragraph.");
});

Deno.test("markdownToText - removes markdown formatting", async () => {
  const markdown = "**Bold** and *italic* text";
  const result = await markdownToText(markdown);
  assertEquals(result, "Bold and italic text");
});

Deno.test("markdownToText - handles links", async () => {
  const markdown = "[Link text](https://example.com)";
  const result = await markdownToText(markdown);
  assertEquals(result, "Link text");
});

Deno.test("markdownToText - handles lists", async () => {
  const markdown = "- Item 1\n- Item 2\n- Item 3";
  const result = await markdownToText(markdown);
  assertEquals(result, "Item 1\nItem 2\nItem 3");
});

Deno.test("markdownToText - handles empty input", async () => {
  const markdown = "";
  const result = await markdownToText(markdown);
  assertEquals(result, "");
});

Deno.test("markdownToText - handles code blocks", async () => {
  const markdown = "```\nconst x = 1;\n```";
  const result = await markdownToText(markdown);
  assertEquals(result, "const x = 1;");
});
