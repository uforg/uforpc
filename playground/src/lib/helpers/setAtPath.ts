// deno-lint-ignore-file no-explicit-any

/**
 * Configuration options for the setAtPath function
 */
export interface SetAtPathOptions {
  /**
   * If true, properties with null or undefined values will be removed from the object
   * rather than being set to null or undefined.
   * For arrays, elements will be removed entirely, causing subsequent elements to shift left.
   * @default false
   */
  removeNullOrUndefined?: boolean;
}

/**
 * Updates a JSON object at the specified dot notation path, returning a new object without mutating the original.
 *
 * @typeParam T - The type of the original JSON object.
 * @param originalJson - The original JSON object to update.
 * @param path - The dot-separated path to the property to update (e.g., `"a.b.c"` or `"arr.0.prop"`).
 * @param value - The value to set at the specified path.
 * @param options - Configuration options for how the update is performed.
 * @returns A new JSON object with the updated value.
 */
export function setAtPath<T extends object>(
  originalJson: T,
  path: string,
  value: any,
  options?: SetAtPathOptions,
): T {
  const { removeNullOrUndefined = false } = options || {};

  // Split the path into individual keys
  const keys = path.split(".");

  // Recursive function to update nested properties
  const updateNestedProperty = (obj: any, keyIndex: number): any => {
    // Base case: we've processed all keys, return the new value
    if (keyIndex === keys.length) {
      return value;
    }

    // Get the current key we're processing
    const currentKey = keys[keyIndex];

    // Get the current value at this key (if it exists and is an object)
    const currentValue = obj != null && typeof obj === "object"
      ? obj[currentKey]
      : undefined;

    // Recursively process the next level
    const updatedValue = updateNestedProperty(currentValue, keyIndex + 1);

    // Handle arrays differently than objects
    if (Array.isArray(obj)) {
      // Create a copy of the array
      const arrayCopy = obj.slice();

      // If we're at the parent of the target property and the value is null/undefined
      // and removeNullOrUndefined is true, then remove the array element completely
      if (
        keyIndex === keys.length - 1 &&
        removeNullOrUndefined &&
        (updatedValue === null || updatedValue === undefined)
      ) {
        arrayCopy.splice(Number(currentKey), 1);
        return arrayCopy;
      }

      // Otherwise, update the specific index
      arrayCopy[Number(currentKey)] = updatedValue;
      return arrayCopy;
    }

    // If we're at the parent of the target property and the value is null/undefined
    // and removeNullOrUndefined is true, don't include the property in the result
    if (
      keyIndex === keys.length - 1 &&
      removeNullOrUndefined &&
      (updatedValue === null || updatedValue === undefined)
    ) {
      const result = { ...obj };
      delete result[currentKey];
      return result;
    }

    // For objects, create a new object with the updated property
    return {
      ...obj,
      [currentKey]: updatedValue,
    };
  };

  // Start the recursive update from the root object
  return updateNestedProperty(originalJson, 0) as T;
}
