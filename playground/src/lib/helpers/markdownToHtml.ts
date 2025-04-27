import { marked } from "marked";

export async function markdownToHtml(markdown: string): Promise<string> {
  return await marked.parse(markdown);
}
