/** Primitive types supported by the schema */
export type PrimitiveType = "string" | "int" | "float" | "boolean";

/** Represents the object type in the schema */
export type ObjectType = "object";

/** Structure for array types, containing base type and dimensions */
export interface ArrayType {
  baseType: FieldType;
  dimensions: number;
}

/** Represents a custom type name (must start with uppercase) */
export type CustomType = string;

/**
 * Unified type system that represents all possible field types in the schema
 * This includes primitive types, custom types, object type, and array types
 */
export type FieldType = PrimitiveType | ObjectType | ArrayType | CustomType;

/** Represents a primitive value that can be used in metadata */
export type MetaValue = string | number | boolean;

/**
 * Unified field definition that includes type information, description,
 * and optional nested fields for object types
 */
export interface DetailedField {
  type: FieldType;
  desc?: string;
  fields?: Record<string, DetailedField>;
}

/**
 * Represents a type definition in the schema, containing its name,
 * optional description, and field definitions
 */
export interface TypeDefinition {
  name: string;
  desc?: string;
  fields: Record<string, DetailedField>;
}

/** Defines the type of procedure (query for read operations, mutation for write operations) */
export type ProcedureType = "query" | "mutation";

/**
 * Represents a procedure definition in the schema, including its name, type,
 * optional description, input/output definitions, and metadata
 */
export interface Procedure {
  name: string;
  type: ProcedureType;
  desc?: string;
  input?: Record<string, DetailedField>;
  output?: Record<string, DetailedField>;
  meta?: Record<string, MetaValue>;
}

/**
 * Represents the complete schema structure, containing type definitions
 * and procedure definitions
 */
export interface Schema {
  types?: TypeDefinition[];
  procedures: Procedure[];
}

// Type checking functions

/**
 * Checks if the given type is a primitive type
 * @param type - The type to check
 * @returns True if the type is a primitive type, false otherwise
 */
export function isPrimitiveType(type: FieldType): type is PrimitiveType {
  return typeof type === "string" &&
    ["string", "int", "float", "boolean"].includes(type);
}

/**
 * Checks if the given type is a custom type (starts with uppercase letter)
 * @param type - The type to check
 * @returns True if the type is a custom type, false otherwise
 */
export function isCustomType(type: FieldType): boolean {
  return typeof type === "string" && /^[A-Z][a-zA-Z0-9]*$/.test(type);
}

/**
 * Type guard to check if the given type is an array type
 * @param type - The type to check
 * @returns True if the type is an array type, false otherwise
 */
export function isArrayType(type: FieldType): type is ArrayType {
  return typeof type === "object" && "baseType" in type && "dimensions" in type;
}

/**
 * Checks if the given type is an object type
 * @param type - The type to check
 * @returns True if the type is an object type, false otherwise
 */
export function isObjectType(type: FieldType): boolean {
  return type === "object";
}

/**
 * Checks if the given type is a valid field type
 * @param type - The type to check
 * @returns True if the type is valid, false otherwise
 */
export function isValidFieldType(type: FieldType): boolean {
  if (isPrimitiveType(type) || isObjectType(type)) return true;
  if (isCustomType(type)) return true;
  if (isArrayType(type)) return isValidFieldType(type.baseType);
  return false;
}

/**
 * Type guard to check if a value is a detailed field
 * @param value - The value to check
 * @returns True if the value is a detailed field, false otherwise
 */
export function isDetailedField(value: unknown): value is DetailedField {
  return typeof value === "object" && value !== null && "type" in value;
}

// Type parsing functions

/**
 * Parses a string representation of an array type into an ArrayType structure
 * @param typeStr - The string to parse
 * @returns An ArrayType object if the string represents a valid array type, null otherwise
 */
export function parseArrayType(typeStr: string): ArrayType | null {
  let dimensions = 0;
  let baseTypeStr = typeStr.trim();

  while (baseTypeStr.endsWith("[]")) {
    dimensions++;
    baseTypeStr = baseTypeStr.slice(0, -2).trim();
  }

  if (dimensions === 0) return null;

  const baseType = parseFieldType(baseTypeStr);
  if (!baseType) return null;

  return {
    baseType,
    dimensions,
  };
}

/**
 * Parses a string representation of any type into the corresponding FieldType
 * @param typeStr - The string to parse
 * @returns A FieldType if the string represents a valid type, null otherwise
 */
export function parseFieldType(typeStr: string): FieldType | null {
  const arrayType = parseArrayType(typeStr);
  if (arrayType) return arrayType;

  if (["string", "int", "float", "boolean", "object"].includes(typeStr)) {
    return typeStr as FieldType;
  }

  if (/^[A-Z][a-zA-Z0-9]*$/.test(typeStr)) {
    return typeStr;
  }

  return null;
}

/**
 * Creates an array type with the specified base type and dimensions
 * @param baseType - The base type of the array
 * @param dimensions - The number of array dimensions
 * @returns An ArrayType object
 */
export function createArrayType(
  baseType: FieldType,
  dimensions: number,
): ArrayType {
  return { baseType, dimensions };
}

/**
 * Converts a FieldType to its string representation
 * @param type - The type to convert
 * @returns The string representation of the type
 */
export function fieldTypeToString(type: FieldType): string {
  if (isArrayType(type)) {
    return `${fieldTypeToString(type.baseType)}${"[]".repeat(type.dimensions)}`;
  }
  return type;
}

/**
 * Parses a string into a DetailedField
 * @param value - The string to parse into a DetailedField
 * @returns A DetailedField object
 * @throws Error if the value contains an invalid type
 */
export function parseDetailedField(
  value: string | Partial<DetailedField>,
): DetailedField {
  if (typeof value === "string") {
    const parsedType = parseFieldType(value);
    if (!parsedType) throw new Error(`Invalid type: ${value}`);
    return { type: parsedType };
  }

  if (!value.type) throw new Error("DetailedField must have a type");

  const typeStr = typeof value.type === "string"
    ? value.type
    : fieldTypeToString(value.type);
  const parsedType = parseFieldType(typeStr);
  if (!parsedType) throw new Error(`Invalid type: ${typeStr}`);

  return {
    ...value,
    type: parsedType,
  };
}

/**
 * Gets the base type of a field type by removing all array layers
 * @param type - The field type to analyze
 * @returns The base type without array dimensions
 */
export function getBaseFieldType(
  type: FieldType,
): Exclude<FieldType, ArrayType> {
  return isArrayType(type) ? getBaseFieldType(type.baseType) : type;
}

/**
 * Gets the total dimensions of a field type by counting all nested array layers
 * @param type - The field type to analyze
 * @returns The total number of array dimensions
 */
export function getTotalArrayDimensions(type: FieldType): number {
  if (!isArrayType(type)) return 0;
  return type.dimensions + getTotalArrayDimensions(type.baseType);
}

/**
 * Flattens a nested array type into a single array type with combined dimensions
 * @param type - The field type to flatten
 * @returns A flattened array type or the original type if not an array
 */
export function flattenArrayType(type: FieldType): FieldType {
  if (!isArrayType(type)) return type;

  const baseType = getBaseFieldType(type);
  const dimensions = getTotalArrayDimensions(type);

  return dimensions > 0 ? createArrayType(baseType, dimensions) : baseType;
}
