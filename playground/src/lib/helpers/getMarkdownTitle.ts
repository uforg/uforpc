/**
 * Get the title of a markdown document, which is the first #-header
 *
 * If no #-header is found, the title is "Untitled ${first-line}"
 *
 * If multiple #-headers are found, the first one is used
 *
 * @param markdown
 */
export function getMarkdownTitle(markdown: string): string {
  const lines = markdown
    .split("\n")
    .map((line) => line.trim())
    .filter((line) => line.length > 0);

  // First look for level-1 headers (starts with single #)
  const titleLevel1 = lines.find((line) =>
    line.startsWith("# ") || line === "#"
  );

  if (titleLevel1) {
    return titleLevel1.slice(1).trim();
  }

  // If no level-1 header found, look for any header (starts with any number of #)
  const anyHeader = lines.find((line) => /^#{2,}\s/.test(line));

  if (anyHeader) {
    // Remove all # symbols from the beginning and trim
    return anyHeader.replace(/^#+/, "").trim();
  }

  // If no header at all, return "Untitled: first line"
  return `Untitled: ${lines[0]}`;
}
