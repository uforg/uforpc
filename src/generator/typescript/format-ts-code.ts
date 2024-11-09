import prettier from "prettier";

/**
 * Formats an string of typescript code
 */
export async function formatTsCode(code: string): Promise<string> {
  return await prettier.format(code, { parser: "typescript" });
}
