/**
 * Join path parts with a separator
 *
 * @param parts - The parts to join
 * @returns The joined path
 */
export function joinPath(parts: string[]) {
  const separator = "/";

  // Handle empty array
  if (parts.length === 0) {
    return "";
  }

  // Join all parts with separator
  const joined = parts.join(separator);

  // First, preserve protocols by temporarily replacing them
  const protocolMarker = "___PROTOCOL___";
  const protocolRegex = /(\w+):\/\//g;

  // Store protocols with a temporary marker
  const withMarkers = joined.replace(protocolRegex, `$1:${protocolMarker}`);

  // Now replace multiple consecutive separators with single separator
  const normalized = withMarkers.replace(/\/+/g, separator);

  // Restore the protocols
  let result = normalized.replace(new RegExp(`${protocolMarker}`, "g"), "//");

  // Handle special case where protocol is followed by path parts
  // e.g., "file://" + "home" should become "file://home" not "file:///home"
  result = result.replace(/(:\/\/)\/+/g, "://");

  return result;
}
