/**
 * Updates a JSON object at the specified dot notation path, returning a new object without mutating the original.
 * Properties with null or undefined values will be removed from the object.
 * For arrays, elements will be removed entirely, causing subsequent elements to shift left.
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

    // Check if the next key in the path is numeric (indicating we need an array)
    const nextKeyIsNumeric =
      keyIndex < keys.length - 1 && /^\d+$/.test(keys[keyIndex + 1]);

    // Get the current value at this key (if it exists and is an object)
    const currentValue =
      obj != null && typeof obj === "object" ? obj[currentKey] : undefined;

    // Create the proper container for the next level if it doesn't exist
    // Use an array if the next key is numeric, otherwise use an object
    let nextValue = currentValue;
    if (nextValue === undefined) {
      nextValue = nextKeyIsNumeric ? [] : {};
    }

    // Recursively process the next level
    const updatedValue = updateNestedProperty(nextValue, keyIndex + 1);

    // Handle arrays differently than objects
    if (Array.isArray(obj)) {
      // Create a copy of the array
      const arrayCopy = obj.slice();

      // If the updated value is null/undefined, remove the array element completely
      if (
        keyIndex === keys.length - 1 &&
        (updatedValue === null || updatedValue === undefined)
      ) {
        arrayCopy.splice(Number(currentKey), 1);
        return arrayCopy;
      }

      // Otherwise, update the specific index
      arrayCopy[Number(currentKey)] = updatedValue;
      return arrayCopy;
    }

    // If the updated value is null/undefined, don't include the property in the result
    if (
      keyIndex === keys.length - 1 &&
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
