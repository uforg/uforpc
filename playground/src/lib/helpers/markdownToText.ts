import { markdownToHtml } from "./markdownToHtml.ts";

/**
 * Convert markdown to text
 * @param markdown - The markdown to convert
 * @returns The text
 */
export async function markdownToText(markdown: string): Promise<string> {
  const html = await markdownToHtml(markdown);
  return html.replace(/<(?:.|\n)*?>/gm, "").trim();
}
