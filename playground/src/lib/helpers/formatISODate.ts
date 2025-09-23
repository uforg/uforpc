/**
 * Formats an ISO 8601 date string into a more human-readable format.
 *
 * @param isoDate - date in ISO 8601 format (e.g., "2000-10-05T14:48:00.000Z")
 * @returns formatted date like "2000-10-05 14:48:00"
 */
export function formatISODate(isoDate: string): string {
  return isoDate.replaceAll("T", " ").replaceAll("Z", "").split(".")[0];
}
