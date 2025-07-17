import { marked } from "marked";

import { normalizeIndent } from "./normalizeIndent";

export async function markdownToHtml(markdown: string): Promise<string> {
  return await marked.parse(normalizeIndent(markdown));
}
