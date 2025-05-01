// deno-lint-ignore-file no-explicit-any

/**
 * Updates a JSON object at the specified dot notation path, returning a new object without mutating the original.
 *
 * @typeParam T - The type of the original JSON object.
 * @param originalJson - The original JSON object to update.
 * @param path - The dot-separated path to the property to update (e.g., `"a.b.c"` or `"arr.0.prop"`).
 * @param value - The value to set at the specified path.
 * @returns A new JSON object with the updated value.
 */
export function setAtPath<T extends object>(
  originalJson: T,
  path: string,
  value: any,
): T {
  const keys = path.split(".");
  const updater = (obj: any, idx: number): any => {
    if (idx === keys.length) return value;
    const key = keys[idx];
    const current = obj != null && typeof obj === "object"
      ? obj[key]
      : undefined;
    const updated = updater(current, idx + 1);
    if (Array.isArray(obj)) {
      const i = Number(key);
      const copy = obj.slice();
      copy[i] = updated;
      return copy;
    }
    return { ...obj, [key]: updated };
  };
  return updater(originalJson, 0) as T;
}
