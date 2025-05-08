import { describe, expect, it } from "vitest";

import { markdownToText } from "./markdownToText.ts";

describe("markdownToText", () => {
  it("converts basic markdown to plain text", async () => {
    const markdown = "# Title\nThis is a paragraph.";
    const result = await markdownToText(markdown);
    expect(result).toBe("Title\nThis is a paragraph.");
  });

  it("handles multiple paragraphs", async () => {
    const markdown = "First paragraph.\n\nSecond paragraph.";
    const result = await markdownToText(markdown);
    expect(result).toBe("First paragraph.\nSecond paragraph.");
  });

  it("removes markdown formatting", async () => {
    const markdown = "**Bold** and *italic* text";
    const result = await markdownToText(markdown);
    expect(result).toBe("Bold and italic text");
  });

  it("handles links", async () => {
    const markdown = "[Link text](https://example.com)";
    const result = await markdownToText(markdown);
    expect(result).toBe("Link text");
  });

  it("handles lists", async () => {
    const markdown = "- Item 1\n- Item 2\n- Item 3";
    const result = await markdownToText(markdown);
    expect(result).toBe("Item 1\nItem 2\nItem 3");
  });

  it("handles empty input", async () => {
    const markdown = "";
    const result = await markdownToText(markdown);
    expect(result).toBe("");
  });

  it("handles code blocks", async () => {
    const markdown = "```\nconst x = 1;\n```";
    const result = await markdownToText(markdown);
    expect(result).toBe("const x = 1;");
  });
});
