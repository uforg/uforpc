export * from "./types.ts";
import {
  type DetailedField,
  parseDetailedField,
  type Procedure,
  type Schema,
  type TypeDefinition,
} from "./types.ts";
import MainSchema from "./schema.ts";
import type {
  FieldSchemaType,
  ProcedureSchemaType,
  TypeSchemaType,
} from "./schema.ts";

export class SchemaParsingError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "SchemaParsingError";
  }
}

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
 * Transforms a raw field (from JSON) into a DetailedField
 * @param field - Raw field from JSON
 * @returns Transformed DetailedField
 */
function transformField(field: FieldSchemaType): DetailedField {
  if (typeof field === "string") {
    return parseDetailedField(field);
  }

  let transformedFields: Record<string, DetailedField> | undefined;

  if (field.fields) {
    transformedFields = {};
    for (const [key, value] of Object.entries(field.fields)) {
      transformedFields[key] = transformField(value);
    }
  }

  return parseDetailedField({
    type: field.type,
    optional: field.optional,
    desc: field.desc,
    fields: transformedFields,
  });
}

/**
 * Transforms a raw type definition into a TypeDefinition
 * @param type - Raw type from JSON
 * @returns Transformed TypeDefinition
 */
function transformType(type: TypeSchemaType): TypeDefinition {
  const transformedFields: Record<string, DetailedField> = {};

  for (const [key, value] of Object.entries(type.fields)) {
    transformedFields[key] = transformField(value);
  }

  return {
    name: type.name,
    desc: type.desc,
    fields: transformedFields,
  };
}

/**
 * Transforms a raw procedure into a Procedure
 * @param procedure - Raw procedure from JSON
 * @returns Transformed Procedure
 */
function transformProcedure(procedure: ProcedureSchemaType): Procedure {
  let transformedInput: Record<string, DetailedField> | undefined;
  let transformedOutput: Record<string, DetailedField> | undefined;

  if (procedure.input) {
    transformedInput = {};
    for (const [key, value] of Object.entries(procedure.input)) {
      transformedInput[key] = transformField(value);
    }
  }

  if (procedure.output) {
    transformedOutput = {};
    for (const [key, value] of Object.entries(procedure.output)) {
      transformedOutput[key] = transformField(value);
    }
  }

  return {
    name: procedure.name,
    type: procedure.type,
    desc: procedure.desc,
    input: transformedInput,
    output: transformedOutput,
    meta: procedure.meta,
  };
}

/**
 * Parses and validates a schema string, then transforms it into our internal Schema type
 * @param content - The schema string to parse
 * @returns Parsed and transformed Schema
 * @throws SchemaParsingError if the JSON is invalid
 * @throws SchemaValidationError if the schema is invalid
 */
export function parseSchema(content: string): Schema {
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

  const parsed = MainSchema.safeParse(rawSchema);
  if (!parsed.success) {
    const errs: string[] = [];

    for (const issue of parsed.error.issues ?? []) {
      errs.push(issue.path.join(".") + ": " + issue.message);
    }

    throw new SchemaValidationError("Schema validation failed", errs);
  }

  return {
    types: parsed.data.types?.map(transformType),
    procedures: parsed.data.procedures.map(transformProcedure),
  };
}
