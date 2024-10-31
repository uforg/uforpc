import type { FieldSchemaType } from "./schema.ts";

/**
 * Checks if the given type is a primitive type
 * @param type - The type to check
 * @returns True if the type is a primitive type, false otherwise
 */
export function isPrimitiveType(type: FieldSchemaType): boolean {
  return typeof type.type === "string" &&
    ["string", "int", "float", "boolean"].includes(type.type);
}

/**
 * Checks if the given type is a custom type (starts with uppercase letter)
 * @param type - The type to check
 * @returns True if the type is a custom type, false otherwise
 */
export function isCustomType(type: FieldSchemaType): boolean {
  return typeof type.type === "string" && /^[A-Z][a-zA-Z0-9]*$/.test(type.type);
}

/**
 * Checks if the given type is an object type
 * @param type - The type to check
 * @returns True if the type is an object type, false otherwise
 */
export function isObjectType(type: FieldSchemaType): boolean {
  return typeof type.type === "string" && type.type === "object";
}

/**
 * Returns the base type and dimensions of an array type
 * @param type - The array type to analyze
 * @returns The base type of the array
 */
export function parseArrayType(
  type: FieldSchemaType,
): { type: FieldSchemaType; dimensions: number } {
  if (!isArrayType(type)) return { type, dimensions: 0 };

  let dimensions = 0;
  let baseTypeStr = type.type.trim();

  while (baseTypeStr.endsWith("[]")) {
    dimensions++;
    baseTypeStr = baseTypeStr.slice(0, -2).trim();
  }

  return { type: { type: baseTypeStr }, dimensions };
}

/**
 * Type guard to check if the given type is an array type
 * @param type - The type to check
 * @returns True if the type is an array type, false otherwise
 */
export function isArrayType(type: FieldSchemaType): boolean {
  if (typeof type.type !== "string") return false;
  if (!type.type.endsWith("[]")) return false;

  const baseType = type.type.replaceAll("[]", "");
  if (isPrimitiveType({ type: baseType })) return true;
  if (isCustomType({ type: baseType })) return true;
  if (isObjectType({ type: baseType })) return true;

  return false;
}

/**
 * Checks if the given type is a valid field type
 * @param type - The type to check
 * @returns True if the type is valid, false otherwise
 */
export function isValidFieldType(type: FieldSchemaType): boolean {
  if (isPrimitiveType(type)) return true;
  if (isCustomType(type)) return true;
  if (isObjectType(type)) return true;
  if (isArrayType(type)) return true;
  return false;
}
