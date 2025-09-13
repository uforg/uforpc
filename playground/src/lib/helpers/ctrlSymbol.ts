/**
 * Returns the appropriate control key symbol based on the user's platform.
 * On macOS, it returns "⌘", otherwise "CTRL".
 *
 * @returns The control key symbol as a string.
 */
export function ctrlSymbol() {
  const isMac = /mac/.test(navigator.userAgent.toLowerCase());
  return isMac ? "⌘" : "CTRL";
}
