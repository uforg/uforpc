/**
 * Delete all level-1 headers from a markdown document
 *
 * @param markdown
 */
export function deleteMarkdownHeadings(markdown: string): string {
  // Use global flag (g) to remove all level-1 headers
  // Specifically target only single hash headers (level 1)
  return markdown.replace(/^#\s+.*\n/gm, "");
}
