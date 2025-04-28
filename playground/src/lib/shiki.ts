import { createHighlighter } from "shiki";
import type { BundledTheme, Highlighter } from "shiki";

const lightTheme: BundledTheme = "github-light";
const darkTheme: BundledTheme = "github-dark";
let highlighterInstance: Highlighter | null = null;

export const getHighlighter = async (): Promise<Highlighter> => {
  if (highlighterInstance) {
    return highlighterInstance;
  }

  const urpcSyntaxUrl =
    "https://cdn.jsdelivr.net/gh/uforg/uforpc-vscode/syntaxes/urpc.tmLanguage.json";
  const urpcSyntax = await fetch(urpcSyntaxUrl);
  const urpcSyntaxJson = await urpcSyntax.json();
  urpcSyntaxJson.name = "urpc";

  highlighterInstance = await createHighlighter({
    langs: [urpcSyntaxJson],
    themes: [lightTheme, darkTheme],
  });

  return highlighterInstance;
};

export { darkTheme, lightTheme };
