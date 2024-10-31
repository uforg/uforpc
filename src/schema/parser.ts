import type { FieldSchemaType } from "@/schema/schema.ts";
import MainSchema from "./schema.ts";
import type { MainSchemaType } from "./schema.ts";
import { isCustomType } from "@/schema/helpers.ts";

/** Custom error classes for schema parsing */
export class SchemaParsingError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "SchemaParsingError";
  }
}

/** Custom error classes for schema validation */
export class SchemaValidationError extends Error {
  constructor(
    message: string,
    public errors: string[],
  ) {
    super(`${message}: ${errors.join(", ")}`);
    this.name = "SchemaValidationError";
  }
}

/**
 * Parses and validates a schema string, then transforms it into our internal Schema type
 * @param content - The schema string to parse
 * @returns Parsed and transformed Schema
 * @throws SchemaParsingError if the JSON is invalid
 * @throws SchemaValidationError if the schema is invalid
 */
export function parseSchema(content: string): MainSchemaType {
  // Parse JSON
  let rawSchema: unknown;
  try {
    rawSchema = JSON.parse(content);
  } catch (error) {
    throw new SchemaParsingError(
      `Invalid JSON: ${
        error instanceof Error ? error.message : "unknown error"
      }`,
    );
  }

  const parsedSchema = MainSchema.safeParse(rawSchema);
  if (!parsedSchema.success) {
    const errs: string[] = [];

    for (const issue of parsedSchema.error.issues ?? []) {
      errs.push(issue.path.join(".") + ": " + issue.message);
    }

    throw new SchemaValidationError("Schema validation failed", errs);
  }

  assertCustomTypesUniqueness(parsedSchema.data);
  assertCustomTypeDefinitions(parsedSchema.data);

  return parsedSchema.data;
}

/**
 * Asserts that custom type names are unique
 * @param schema - Schema to validate
 * @throws SchemaValidationError if a custom type name is not unique
 */
function assertCustomTypesUniqueness(schema: MainSchemaType) {
  const typeNames: Record<string, boolean> = {};
  for (const type of schema.types ?? []) {
    if (typeNames[type.name]) {
      throw new SchemaValidationError(`Duplicate type name: ${type.name}`, []);
    }
    typeNames[type.name] = true;
  }
}

/**
 * Asserts that all used custom types are defined in the schema
 * @param schema - Schema to validate
 * @throws SchemaValidationError if a custom type is not defined
 */
function assertCustomTypeDefinitions(schema: MainSchemaType) {
  const findCustomTypes = (field: FieldSchemaType): string[] => {
    const customTypes: string[] = [];

    if (field.type.startsWith("object")) {
      for (const key in field.fields) {
        const customTypes = findCustomTypes(field.fields[key]);
        customTypes.push(...customTypes);
      }
    }
    if (isCustomType(field)) customTypes.push(field.type);

    return customTypes;
  };

  const desiredCustomTypes: string[] = [];
  for (const procedure of schema.procedures) {
    for (const input in procedure.input) {
      desiredCustomTypes.push(...findCustomTypes(procedure.input[input]));
    }
    for (const output in procedure.output) {
      desiredCustomTypes.push(...findCustomTypes(procedure.output[output]));
    }
  }

  const existingCustomTypes = (schema.types ?? []).map((type) => type.name);
  for (const desired of desiredCustomTypes) {
    if (!existingCustomTypes.includes(desired)) {
      throw new SchemaValidationError(
        `Custom type ${desired} is not defined in the schema`,
        [],
      );
    }
  }
}
