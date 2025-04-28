/**
 * Extracts a node (rule, type, or procedure) from a URPC schema
 *
 * @param schema - The URPC schema string to search in
 * @param kind - The kind of node to extract ('rule', 'type', or 'proc')
 * @param nodeName - The name of the node to find
 * @returns The complete node definition as a string, or null if not found
 */
export function extractNodeFromSchema(
  schema: string,
  kind: "rule" | "type" | "proc",
  nodeName: string,
): string | null {
  const lines = schema.split("\n");
  const nodePattern = kind === "rule"
    ? new RegExp(`^\\s*${kind}\\s+@${nodeName}\\s*{`)
    : new RegExp(`^\\s*${kind}\\s+${nodeName}\\s*{`);

  let openBraces = 0;
  let foundNode = false;
  const result: string[] = [];

  for (const line of lines) {
    // If we found the node previously, collect lines until matching closing brace
    if (foundNode) {
      result.push(line);

      // Count braces
      openBraces += (line.match(/{/g) || []).length -
        (line.match(/}/g) || []).length;

      // Found closing brace - return the complete node
      if (openBraces === 0 && line.trim() === "}") {
        return result.join("\n");
      }
    } // Check if this line matches the node pattern
    else if (nodePattern.test(line)) {
      foundNode = true;
      result.push(line);
      openBraces = (line.match(/{/g) || []).length -
        (line.match(/}/g) || []).length;
    }
  }

  return null;
}
