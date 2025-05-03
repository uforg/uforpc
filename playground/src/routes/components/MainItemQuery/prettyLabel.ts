/**
 * Remove the root. prefix from a label and convert
 * array indexes from .n to [n]
 */
export function prettyLabel(label: string) {
  return label.replace(/^root\./, "").replace(/\.(\d+)/g, "[$1]");
}
