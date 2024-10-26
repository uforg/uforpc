import Ajv from "ajv";
import addFormats from "ajv-formats";
import {
  type DetailedField,
  parseDetailedField,
  type Procedure,
  type Schema,
  type TypeDefinition,
} from "./types.ts";

export class SchemaParsingError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "SchemaParsingError";
  }
}

export class SchemaValidationError extends Error {
  constructor(
    message: string,
    public errors: unknown[],
  ) {
    super(message);
    this.name = "SchemaValidationError";
  }
}

interface RawField {
  type?: string;
  desc?: string;
  fields?: Record<string, RawField | string>;
}

interface RawType {
  name: string;
  desc?: string;
  fields: Record<string, RawField | string>;
}

interface RawProcedure {
  name: string;
  type: "query" | "mutation";
  desc?: string;
  input?: Record<string, RawField | string>;
  output?: Record<string, RawField | string>;
  meta?: Record<string, string | number | boolean>;
}

interface RawSchema {
  types?: RawType[];
  procedures: RawProcedure[];
}

/**
 * Transforms a raw field (from JSON) into a DetailedField
 * @param field - Raw field from JSON
 * @returns Transformed DetailedField
 */
function transformField(field: RawField | string): DetailedField {
  if (typeof field === "string") {
    return parseDetailedField(field);
  }

  // Transform nested fields if they exist
  const transformedFields = field.fields
    ? Object.entries(field.fields).reduce(
      (acc, [key, value]) => ({
        ...acc,
        [key]: transformField(value),
      }),
      {} as Record<string, DetailedField>,
    )
    : undefined;

  return parseDetailedField({
    type: field.type!,
    desc: field.desc,
    fields: transformedFields,
  });
}

/**
 * Transforms a raw type definition into a TypeDefinition
 * @param type - Raw type from JSON
 * @returns Transformed TypeDefinition
 */
function transformType(type: RawType): TypeDefinition {
  const transformedFields = Object.entries(type.fields).reduce(
    (acc, [key, value]) => ({
      ...acc,
      [key]: transformField(value),
    }),
    {} as Record<string, DetailedField>,
  );

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
function transformProcedure(procedure: RawProcedure): Procedure {
  return {
    name: procedure.name,
    type: procedure.type,
    desc: procedure.desc,
    input: procedure.input
      ? Object.entries(procedure.input).reduce(
        (acc, [key, value]) => ({
          ...acc,
          [key]: transformField(value),
        }),
        {} as Record<string, DetailedField>,
      )
      : undefined,
    output: procedure.output
      ? Object.entries(procedure.output).reduce(
        (acc, [key, value]) => ({
          ...acc,
          [key]: transformField(value),
        }),
        {} as Record<string, DetailedField>,
      )
      : undefined,
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
export async function parseSchema(content: string): Promise<Schema> {
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

  // Setup Ajv
  const ajv = new Ajv({
    allErrors: true,
    strict: true,
    strictTypes: true,
    strictRequired: true,
    allowUnionTypes: true,
  });
  addFormats(ajv);

  // Load and validate against JSON Schema
  const schemaJson = JSON.parse(
    await Deno.readTextFile(new URL("./schema.json", import.meta.url)),
  );
  const validate = ajv.compile(schemaJson);

  if (!validate(rawSchema)) {
    throw new SchemaValidationError(
      "Schema validation failed",
      validate.errors ?? [],
    );
  }

  // Transform to our internal types
  const typedSchema = rawSchema as RawSchema;

  return {
    types: typedSchema.types?.map(transformType),
    procedures: typedSchema.procedures.map(transformProcedure),
  };
}
