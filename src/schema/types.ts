//TODO: Support nested arrays (e.g., "string[][]")

/**
 * Represents primitive types supported by UFO RPC
 */
export type PrimitiveType = "string" | "number" | "float" | "boolean";

/**
 * Represents a custom type name that must start with an uppercase letter
 */
export type CustomType = string;

/**
 * Represents an array type notation (e.g., "string[]", "User[]")
 */
export type ArrayType = `${PrimitiveType | CustomType}[]`;

/**
 * Union type of all valid types in the schema
 */
export type ValidType = PrimitiveType | CustomType | ArrayType | "object";

/**
 * Represents a simple field definition using just the type name
 */
export type SimpleField = ValidType;

/**
 * Represents a detailed field definition with additional metadata
 * @interface DetailedField
 * @property {ValidType} type - The type of the field
 * @property {string} [desc] - Optional description of the field
 * @property {Record<string, Field>} [fields] - Optional nested fields for object types
 */
export interface DetailedField {
  type: ValidType;
  desc?: string;
  fields?: Record<string, Field>;
}

/**
 * Union type representing either a simple field or a detailed field definition
 */
export type Field = SimpleField | DetailedField;

/**
 * Represents a custom type definition in the schema
 * @interface Type
 * @property {string} name - The name of the type (must start with uppercase)
 * @property {string} [desc] - Optional description of the type
 * @property {Record<string, Field>} fields - The fields that compose this type
 */
export interface Type {
  name: string;
  desc?: string;
  fields: Record<string, Field>;
}

/**
 * Represents the type of procedure (query or mutation)
 */
export type ProcedureType = "query" | "mutation";

/**
 * Represents a procedure definition in the schema
 * @interface Procedure
 * @property {string} name - The name of the procedure (must start with lowercase)
 * @property {ProcedureType} type - The type of procedure
 * @property {string} [desc] - Optional description of the procedure
 * @property {Record<string, Field>} [input] - Optional input parameters
 * @property {Record<string, Field>} [output] - Optional output parameters
 * @property {Record<string, unknown>} [meta] - Optional metadata for the procedure
 */
export interface Procedure {
  name: string;
  type: ProcedureType;
  desc?: string;
  input?: Record<string, Field>;
  output?: Record<string, Field>;
  meta?: Record<string, unknown>;
}

/**
 * Represents the complete schema structure
 * @interface Schema
 * @property {Type[]} [types] - Optional array of custom type definitions
 * @property {Procedure[]} procedures - Array of procedure definitions
 */
export interface Schema {
  types?: Type[];
  procedures: Procedure[];
}

/**
 * Type guard to check if a field is a detailed field definition
 * @param {Field} field - The field to check
 * @returns {field is DetailedField} True if the field is a detailed field
 */
export function isDetailedField(field: Field): field is DetailedField {
  return typeof field === "object" && "type" in field;
}

/**
 * Checks if a type name is a primitive type
 * @param {string} type - The type name to check
 * @returns {type is PrimitiveType} True if the type is a primitive type
 */
export function isPrimitiveType(type: string): type is PrimitiveType {
  return ["string", "number", "float", "boolean"].includes(type);
}

/**
 * Checks if a type name represents an array type
 * @param {string} type - The type name to check
 * @returns {type is ArrayType} True if the type is an array type
 */
export function isArrayType(type: string): type is ArrayType {
  return type != "[]" && type.endsWith("[]");
}

/**
 * Gets the base type from a type name (removes array notation if present)
 * @param {string} type - The type name to process
 * @returns {string} The base type name
 */
export function getBaseType(type: string): string {
  return isArrayType(type) ? type.slice(0, -2) : type;
}

/**
 * Checks if a type name represents a custom type
 * @param {string} type - The type name to check
 * @returns {boolean} True if the type is a custom type
 */
export function isCustomType(type: string): boolean {
  const baseType = getBaseType(type);
  return /^[A-Z][a-zA-Z0-9]*$/.test(baseType);
}

/**
 * Checks if a type name is valid according to the schema rules
 * @param {string} type - The type name to check
 * @returns {type is ValidType} True if the type is valid
 */
export function isValidType(type: string): type is ValidType {
  const baseType = getBaseType(type);
  return isPrimitiveType(baseType) || isCustomType(baseType) ||
    baseType === "object";
}
