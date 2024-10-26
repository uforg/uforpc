import Ajv from "ajv";
import addFormats from "ajv-formats";
import { isPrimitiveType, isValidType } from "./types.ts";
import type { Field, Schema } from "./types.ts";

export class SchemaValidationError extends Error {
  constructor(
    message: string,
    public errors: unknown[],
  ) {
    super(message);
    this.name = "SchemaValidationError";
  }
}

export class SchemaParsingError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "SchemaParsingError";
  }
}

function validateTypeReferences(
  fields: Record<string, Field>,
  definedTypes: Set<string>,
  path: string,
): string[] {
  return Object.entries(fields).flatMap(([fieldName, field]) => {
    const currentPath = `${path}.${fieldName}`;
    const errors: string[] = [];

    const fieldType = typeof field === "string" ? field : field.type;
    const baseType = fieldType.replace(/\[\]/g, "");

    if (
      !isPrimitiveType(baseType) && baseType !== "object" &&
      !definedTypes.has(baseType)
    ) {
      errors.push(
        `Type "${baseType}" referenced at "${currentPath}" is not defined`,
      );
    }

    if (!isValidType(fieldType)) {
      errors.push(`Invalid type "${fieldType}" at "${currentPath}"`);
    }

    if (typeof field === "object" && field.fields) {
      errors.push(
        ...validateTypeReferences(field.fields, definedTypes, currentPath),
      );
    }

    return errors;
  });
}

function validateProcedureNames(procedures: Schema["procedures"]): string[] {
  return procedures.flatMap((proc) => {
    if (!/^[A-Z][a-zA-Z0-9]*$/.test(proc.name)) {
      return [
        `Procedure name "${proc.name}" must start with uppercase letter and contain only alphanumeric characters`,
      ];
    }
    return [];
  });
}

function validateTypeNames(types: Schema["types"] = []): string[] {
  return types.flatMap((type) => {
    if (!/^[A-Z][a-zA-Z0-9]*$/.test(type.name)) {
      return [
        `Type name "${type.name}" must start with uppercase letter and contain only alphanumeric characters`,
      ];
    }
    return [];
  });
}

function validateSchema(schema: Schema): void {
  const definedTypes = new Set(schema.types?.map((type) => type.name) ?? []);
  const errors: string[] = [];

  // Validate type names
  errors.push(...validateTypeNames(schema.types));

  // Validate procedure names
  errors.push(...validateProcedureNames(schema.procedures));

  // Validate type references in custom types
  schema.types?.forEach((type) => {
    errors.push(
      ...validateTypeReferences(type.fields, definedTypes, `type:${type.name}`),
    );
  });

  // Validate type references in procedures
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
    throw new SchemaValidationError("Schema validation failed", errors);
  }
}

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

export { isPrimitiveType, isValidType };
export type { Field, Schema };
