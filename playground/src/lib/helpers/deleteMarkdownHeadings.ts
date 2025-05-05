/**
 * Delete all level-1 headers from a markdown document
 *
 * @param markdown
 */
export function deleteMarkdownHeadings(markdown: string): string {
  const lines = markdown.split("\n");

  // Filter out level-1 headers
  const filteredLines = lines.filter((line) => {
    const trimmedLine = line.trim();
    if (!trimmedLine.startsWith("#")) return true;

    if (trimmedLine.length === 1) return false;
    if (trimmedLine[1] !== "#") return false;

    return true;
  });

  // Join the remaining lines back together
  return filteredLines.join("\n");
}
