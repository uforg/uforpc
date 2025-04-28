import type { SearchResult } from "minisearch";

export function markSearchHints(
  searchResult: SearchResult,
  textToMark: string,
) {
  const regexp = new RegExp(`(${searchResult.terms.join("|")})`, "gi");
  return textToMark.replace(regexp, "<mark>$1</mark>");
}
