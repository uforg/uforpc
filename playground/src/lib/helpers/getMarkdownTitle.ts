/**
 * Get the title of a markdown document, which is the first level-1 heading (#-heading)
 *
 * If no level-1 heading is found, the title is "Untitled ${first-line}"
 *
 * If multiple level-1 headings are found, the first one is used
 *
 * @param markdown
 */
export function getMarkdownTitle(markdown: string): string {
  const lines = markdown
    .split("\n")
    .map((line) => line.trim())
    .filter((line) => line.length > 0);

  // Look for level-1 headings - starts with single # followed by a character that is not #
  const titleLevel1 = lines.find((line) => /^#([^#]|$)/.test(line));

  if (titleLevel1) {
    return titleLevel1.replace(/^#\s*/, "").trim();
  }

  // If no level-1 heading found, return "Untitled: first line"
  return `Untitled: ${lines[0]}`;
}
