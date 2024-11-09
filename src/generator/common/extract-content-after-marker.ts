/**
 * Extracts and returns the content after the **\/\*\* START FROM HERE \*\*\/**
 * marker.
 *
 * It's used to get the content of a file that is generated from a template
 * and has a marker to indicate where the generated content starts.
 *
 * @param fileContent - The full content of the file.
 * @returns The content after the marker.
 */
export function extractContentAfterMarker(fileContent: string): string {
  const marker = "/** START FROM HERE **/";
  const index = fileContent.indexOf(marker);

  if (index === -1) return fileContent;
  return fileContent.substring(index + marker.length);
}
