/**
 * Normalizes the indentation of a multi-line string.
 *
 * It determines the base indentation from the first non-empty line
 * and removes that same amount of leading whitespace from every line.
 *
 * @param text The input string with potentially inconsistent indentation.
 * @returns The text with normalized indentation.
 */
export function normalizeIndent(text: string): string {
  // 1. Split the text into an array of lines.
  const lines = text.split("\n");

  // 2. Find the first line that is not empty.
  const firstNonEmptyLine = lines.find((line) => line.trim() !== "");

  // If there's no content, return the original text.
  if (!firstNonEmptyLine) {
    return text;
  }

  // 3. Capture the leading whitespace (tabs or spaces) from that first line.
  const indentMatch = firstNonEmptyLine.match(/^(\s*)/);
  const indentation = indentMatch ? indentMatch[0] : "";

  // If there's no indentation on the first line, there's nothing to do.
  if (indentation.length === 0) {
    return text;
  }

  // 4. Map over each line to remove the base indentation.
  const dedentedLines = lines.map((line) => {
    if (line.startsWith(indentation)) {
      // If the line starts with the base indentation, remove it.
      return line.substring(indentation.length);
    }
    // Otherwise, return the line as is.
    return line;
  });

  // 5. Join the lines back into a single string.
  return dedentedLines.join("\n");
}
