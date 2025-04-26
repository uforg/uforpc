import { browser } from "$app/environment";

export function getCurrentHost() {
  if (!browser) return "";
  return `${globalThis.location.protocol}//${globalThis.location.host}`;
}
