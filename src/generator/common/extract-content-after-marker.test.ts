import { assertEquals } from "@std/assert";
import { extractContentAfterMarker } from "./extract-content-after-marker.ts";

Deno.test("extractContentAfterMarker function", async (t) => {
  await t.step("should extract content after the marker", () => {
    const fileContent = `Preamble text
/** START FROM HERE **/
Content after marker.`;
    const expected = `
Content after marker.`;
    const result = extractContentAfterMarker(fileContent);
    assertEquals(result, expected);
  });

  await t.step("should return empty string if marker is at the end", () => {
    const fileContent = `Preamble text
/** START FROM HERE **/`;
    const expected = "";
    const result = extractContentAfterMarker(fileContent);
    assertEquals(result, expected);
  });

  await t.step("should return full string if marker is not found", () => {
    const fileContent = "No marker in this content.";
    const result = extractContentAfterMarker(fileContent);
    assertEquals(result, fileContent);
  });

  await t.step("should handle multiple occurrences of the marker", () => {
    const fileContent = `Preamble
/** START FROM HERE **/
First content
/** START FROM HERE **/
Second content.`;
    const expected = `
First content
/** START FROM HERE **/
Second content.`;
    const result = extractContentAfterMarker(fileContent);
    assertEquals(result, expected);
  });

  await t.step("should handle empty file content", () => {
    const fileContent = ``;
    const expected = "";
    const result = extractContentAfterMarker(fileContent);
    assertEquals(result, expected);
  });
});
