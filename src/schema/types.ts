export type PrimitiveType = "string" | "number" | "float" | "boolean";

export type CustomType = string;

export type ValidType = PrimitiveType | CustomType | "object";

export interface DetailedField {
  type: string; // Supports nested arrays with multiple '[]' suffixes
  desc?: string;
  fields?: Record<string, Field>;
}

export type Field = string | DetailedField;

export interface Type {
  name: string;
  desc?: string;
  fields: Record<string, Field>;
}

export type ProcedureType = "query" | "mutation";

export interface Procedure {
  name: string;
  type: ProcedureType;
  desc?: string;
  input?: Record<string, Field>;
  output?: Record<string, Field>;
  meta?: Record<string, unknown>;
}

export interface Schema {
  types?: Type[];
  procedures: Procedure[];
}

export function isPrimitiveType(type: string): type is PrimitiveType {
  return ["string", "number", "float", "boolean"].includes(type);
}

export function isValidArrayType(type: string): boolean {
  // Remove all array suffixes
  const baseType = type.replace(/\[\]/g, "");
  return baseType.length > 0 && !baseType.includes("[") &&
    !baseType.includes("]");
}

export function isCustomType(type: string): boolean {
  // Get base type without array suffixes
  const baseType = type.replace(/\[\]/g, "");
  return /^[A-Z][a-zA-Z0-9]*$/.test(baseType);
}

export function getBaseType(type: string): string {
  return type.replace(/\[\]/g, "");
}

export function isValidType(type: string): boolean {
  const baseType = getBaseType(type);
  return (isPrimitiveType(baseType) || isCustomType(baseType) ||
    baseType === "object") && isValidArrayType(type);
}
