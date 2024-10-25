import Ajv, { ErrorObject } from "ajv";
import addFormats from "ajv-formats";
import { isPrimitiveType } from "./types.ts";
import type { Field, Schema } from "./types.ts";

/**
 * Custom error for schema validation failures
 */
export class SchemaValidationError extends Error {
  constructor(
    message: string,
    public errors: ErrorObject[] | string[],
  ) {
    super(message);
    this.name = "SchemaValidationError";
  }
}

/**
 * Custom error for JSON parsing failures
 */
export class SchemaParsingError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "SchemaParsingError";
  }
}

/**
 * Type guard to check if a field is of type Field
 */
function isField(field: unknown): field is Field {
  return typeof field === "string" || (
    typeof field === "object" &&
    field !== null &&
    "type" in field &&
    typeof field.type === "string"
  );
}

/**
 * Validates custom type references in the schema
 * @param fields - Record of fields to validate
 * @param definedTypes - Set of defined type names
 * @param path - Current path in the schema for error reporting
 * @returns Array of validation errors
 */
function validateTypeReferences(
  fields: Record<string, Field>,
  definedTypes: Set<string>,
  path: string,
): string[] {
  return Object.entries(fields).flatMap(([fieldName, field]) => {
    const currentPath = `${path}.${fieldName}`;
    const errors: string[] = [];

    if (typeof field === "string") {
      const baseType = field.endsWith("[]") ? field.slice(0, -2) : field;
      if (
        !isPrimitiveType(baseType) && baseType !== "object" &&
        !definedTypes.has(baseType)
      ) {
        errors.push(
          `Type "${baseType}" referenced at "${currentPath}" is not defined`,
        );
      }
    } else if (isField(field)) {
      const baseType = field.type.endsWith("[]")
        ? field.type.slice(0, -2)
        : field.type;
      if (
        !isPrimitiveType(baseType) && baseType !== "object" &&
        !definedTypes.has(baseType)
      ) {
        errors.push(
          `Type "${baseType}" referenced at "${currentPath}" is not defined`,
        );
      }
      if (field.fields) {
        errors.push(
          ...validateTypeReferences(field.fields, definedTypes, currentPath),
        );
      }
    }

    return errors;
  });
}

/**
 * Validates all custom type references in a schema
 * @param schema - The schema to validate
 * @throws {SchemaValidationError} If any type references are invalid
 */
function validateSchema(schema: Schema): void {
  const definedTypes = new Set(schema.types?.map((type) => type.name) ?? []);
  const errors: string[] = [];

  // Validate types
  schema.types?.forEach((type) => {
    errors.push(
      ...validateTypeReferences(type.fields, definedTypes, `type:${type.name}`),
    );
  });

  // Validate procedures
  schema.procedures.forEach((proc) => {
    if (proc.input) {
      errors.push(
        ...validateTypeReferences(
          proc.input,
          definedTypes,
          `procedure:${proc.name}:input`,
        ),
      );
    }
    if (proc.output) {
      errors.push(
        ...validateTypeReferences(
          proc.output,
          definedTypes,
          `procedure:${proc.name}:output`,
        ),
      );
    }
  });

  if (errors.length > 0) {
    throw new SchemaValidationError("Invalid type references found", errors);
  }
}

/**
 * Parses and validates a UFO RPC schema
 * @param content - The schema content as a string
 * @returns Parsed and validated schema
 * @throws {SchemaParsingError} If JSON parsing fails
 * @throws {SchemaValidationError} If schema validation fails
 */
export async function parseSchema(content: string): Promise<Schema> {
  let parsed: unknown;

  try {
    parsed = JSON.parse(content);
  } catch (error) {
    throw new SchemaParsingError(
      `Invalid JSON: ${
        error instanceof Error ? error.message : "unknown error"
      }`,
    );
  }

  const ajv = new Ajv({
    allErrors: true,
    strict: true,
    strictTypes: true,
    strictRequired: true,
  });
  addFormats(ajv);

  const schemaJson = JSON.parse(
    await Deno.readTextFile(new URL("./schema.json", import.meta.url)),
  );

  const validate = ajv.compile(schemaJson);

  if (!validate(parsed)) {
    throw new SchemaValidationError(
      "Schema validation failed",
      validate.errors ?? [],
    );
  }

  const typedSchema = parsed as Schema;
  validateSchema(typedSchema);

  return typedSchema;
}

export { isPrimitiveType } from "./types.ts";
