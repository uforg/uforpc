import { describe, expect, it } from "vitest";

import { deleteMarkdownHeadings } from "./deleteMarkdownHeadings.ts";

describe("deleteMarkdownHeadings", () => {
  it("removes level-1 headings", () => {
    const markdown = "# Title\nDocument content.";
    const result = deleteMarkdownHeadings(markdown);
    expect(result).toBe("Document content.");
  });

  it("keeps document unchanged if no headings exist", () => {
    const markdown = "Content without headings.";
    const result = deleteMarkdownHeadings(markdown);
    expect(result).toBe("Content without headings.");
  });

  it("removes all level-1 headings", () => {
    const markdown = "# First title\nContent\n# Second title\nMore content";
    const result = deleteMarkdownHeadings(markdown);
    expect(result).toBe("Content\nMore content");
  });

  it("preserves non-level-1 headings", () => {
    const markdown = "## Level 2 title\nDocument content.";
    const result = deleteMarkdownHeadings(markdown);
    expect(result).toBe("## Level 2 title\nDocument content.");
  });

  it("removes level-1 headings but keeps other levels", () => {
    const markdown =
      "# Main title\n## Subtitle\nContent\n# Another title\nMore content";
    const result = deleteMarkdownHeadings(markdown);
    expect(result).toBe("## Subtitle\nContent\nMore content");
  });

  it("handles empty documents", () => {
    const markdown = "";
    const result = deleteMarkdownHeadings(markdown);
    expect(result).toBe("");
  });

  it("removes level-1 headings with preceding newlines and spaces", () => {
    const markdown =
      "Some initial content\n\n\n   # Title with preceding newlines and spaces\nFollowing content";
    const result = deleteMarkdownHeadings(markdown);
    expect(result).toBe("Some initial content\n\n\nFollowing content");
  });
});
