/**
 * Extracts a node (type, proc, or stream) from a URPC schema
 *
 * @param schema - The URPC schema string to search in
 * @param kind - The kind of node to extract ('type', 'proc', or 'stream')
 * @param nodeName - The name of the node to find
 * @returns The complete node definition as a string, or null if not found
 */
export function extractNodeFromSchema(
  schema: string,
  kind: "type" | "proc" | "stream",
  nodeName: string,
): string | null {
  const lines = schema.split("\n");
  const patterns = createPatterns(kind, nodeName);
  const state = createInitialState();

  for (const line of lines) {
    if (state.foundNode) {
      const result = processFoundNodeLine(line, state);
      if (result) return result;
    } else {
      const nodeFound = processSearchLine(line, patterns, state);
      if (nodeFound) {
        const result = handleSingleLineNode(line, state);
        if (result) return result;
      }
    }
  }

  return null;
}

/**
 * Creates regex patterns for detecting different node types
 */
function createPatterns(kind: string, nodeName: string) {
  return {
    regular: new RegExp(`^\\s*${kind}\\s+${nodeName}\\s*{`),
    deprecated: new RegExp(`^\\s*deprecated\\s+${kind}\\s+${nodeName}\\s*{`),
    deprecatedMessage: /^\s*deprecated\s*\([^)]*\)\s*$/,
  };
}

/**
 * Creates the initial state for the extraction process
 */
function createInitialState() {
  return {
    openBraces: 0,
    foundNode: false,
    pendingDeprecated: false,
    deprecatedLine: "",
    result: [] as string[],
  };
}

/**
 * Processes a line when we've already found the target node
 * @returns The complete node if extraction is finished, null otherwise
 */
function processFoundNodeLine(
  line: string,
  state: ReturnType<typeof createInitialState>,
): string | null {
  state.result.push(line);

  const braceCount = countBraces(line);
  state.openBraces += braceCount;

  if (state.openBraces === 0) {
    return buildResult(state.result, line);
  }

  return null;
}

/**
 * Processes a line when searching for the target node
 * @returns true if the node was found, false otherwise
 */
function processSearchLine(
  line: string,
  patterns: ReturnType<typeof createPatterns>,
  state: ReturnType<typeof createInitialState>,
): boolean {
  // Check for deprecated node on same line
  if (patterns.deprecated.test(line)) {
    state.foundNode = true;
    state.result.push(line);
    state.openBraces = countBraces(line);
    return true;
  }

  // Check for deprecated message line
  if (patterns.deprecatedMessage.test(line)) {
    state.pendingDeprecated = true;
    state.deprecatedLine = line;
    return false;
  }

  // Check for regular node
  if (patterns.regular.test(line)) {
    state.foundNode = true;

    // Include pending deprecated line if exists
    if (state.pendingDeprecated) {
      state.result.push(state.deprecatedLine);
      resetDeprecatedState(state);
    }

    state.result.push(line);
    state.openBraces = countBraces(line);
    return true;
  }

  // Reset pending deprecated if we encounter a non-matching non-empty line
  if (state.pendingDeprecated && line.trim() !== "") {
    resetDeprecatedState(state);
  }

  return false;
}

/**
 * Handles the case where a node definition fits on a single line
 * @returns The complete node if it's a single line, null otherwise
 */
function handleSingleLineNode(
  line: string,
  state: ReturnType<typeof createInitialState>,
): string | null {
  if (state.openBraces === 0 && line.includes("{") && line.includes("}")) {
    return state.result.join("\n");
  }
  return null;
}

/**
 * Counts the net number of opening braces in a line
 */
function countBraces(line: string): number {
  const openBraces = (line.match(/{/g) || []).length;
  const closeBraces = (line.match(/}/g) || []).length;
  return openBraces - closeBraces;
}

/**
 * Builds the final result string from collected lines
 */
function buildResult(result: string[], lastLine: string): string {
  // Single line case
  if (result.length === 1) {
    return result[0];
  }

  // Multi-line case with proper closing
  if (lastLine.trim() === "}" || lastLine.trim().endsWith("}")) {
    return result.join("\n");
  }

  return result.join("\n");
}

/**
 * Resets the deprecated state tracking
 */
function resetDeprecatedState(
  state: ReturnType<typeof createInitialState>,
): void {
  state.pendingDeprecated = false;
  state.deprecatedLine = "";
}
