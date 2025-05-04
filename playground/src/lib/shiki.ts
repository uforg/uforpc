// deno-lint-ignore-file require-await

import { bundledLanguages, createHighlighter } from "shiki";
import type { BundledTheme, Highlighter } from "shiki";

/**
 * Returns the provided language if it's supported, otherwise falls back to plain text.
 * @param {string} lang - The language identifier to check
 * @returns {string} The original language if supported, or "text" as fallback
 */
export const getOrFallbackLanguage = (lang: string) => {
  const langs = ["urpc", ...Object.keys(bundledLanguages)];
  if (langs.includes(lang)) return lang;
  return "text"; // https://shiki.matsu.io/languages#plain-text
};

const lightTheme: BundledTheme = "github-light";
const darkTheme: BundledTheme = "github-dark";
let highlighterInstance: Highlighter | null = null;
let highlighterPromise: Promise<Highlighter> | null = null;

/**
 * Returns a Shiki highlighter instance with URPC and bundled languages support.
 * This function implements a singleton pattern, returning an existing instance
 * or promise if available, otherwise creating a new highlighter.
 *
 * The highlighter is configured with both light and dark GitHub themes,
 * and includes URPC syntax highlighting loaded from a remote source in
 * addition to the bundled languages.
 *
 * @returns {Promise<Highlighter>} A promise that resolves to a Shiki highlighter instance
 */
export const getHighlighter = async (): Promise<Highlighter> => {
  if (highlighterInstance) {
    return highlighterInstance;
  }

  if (highlighterPromise) {
    return highlighterPromise;
  }

  highlighterPromise = (async () => {
    const urpcSyntaxUrl =
      "https://cdn.jsdelivr.net/gh/uforg/uforpc-vscode/syntaxes/urpc.tmLanguage.json";
    const urpcSyntax = await fetch(urpcSyntaxUrl);
    const urpcSyntaxJson = await urpcSyntax.json();
    urpcSyntaxJson.name = "urpc";

    highlighterInstance = await createHighlighter({
      langs: [urpcSyntaxJson, ...Object.keys(bundledLanguages)],
      themes: [lightTheme, darkTheme],
    });

    return highlighterInstance;
  })();

  return highlighterPromise;
};

export { darkTheme, lightTheme };
