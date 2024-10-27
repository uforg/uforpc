/**
 * ⚠️ IMPORTANT ⚠️
 * If you modify this file, make sure to run the following command in order
 * to auto-generate the corresponding json schema from this zod schema:
 * deno task gen:schema
 */

import { z, type ZodType } from "zod";

/**
 * Schema with regular expression pattern to validate the "type" field.
 * It allows:
 * - Basic types: string, int, float, boolean, object
 * - Custom types starting with an uppercase letter followed by alphanumerics
 * - Array types: e.g., string[], MyType[][], etc.
 */
const TypeRegexSchema = z.string().regex(
  /^(string|int|float|boolean|object|.*\[\]|[A-Z][a-zA-Z0-9]*)$/,
  {
    message:
      'The "type" field must be one of the allowed types: string, int, float, boolean, object, an array (e.g., string[] or MyType[][]), or a custom type name starting with an uppercase letter.',
  },
);

/**
 * Schema with regular expression pattern to validate the "name" field.
 * It requires the name to start with an uppercase letter followed
 * by alphanumerics.
 */
const UpperNameSchema = z.string().regex(
  /^[A-Z][a-zA-Z0-9]*$/,
  {
    message:
      'The "name" field must start with an uppercase letter and contain only alphanumeric characters.',
  },
);

/**
 * Zod schema for Field.
 * Represents a field that can either be a simple type string or
 * a DetailedField.
 */
const FieldSchema = z.lazy(() =>
  z
    .union([
      TypeRegexSchema,
      DetailedFieldSchema,
    ])
    .describe(
      "Definition of a field, which can be a type string or a detailed field definition",
    )
    .superRefine((data, ctx) => {
      if (typeof data === "string" && data === "object") {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message:
            'The field type "object" must be a detailed field definition',
        });
      }
    })
);

/**
 * TypeScript interface for DetailedField.
 * Required to define a recursive structure in zod
 * https://zod.dev/?id=recursive-types
 */
interface IDetailedField {
  type: string;
  desc?: string;
  optional?: boolean;
  fields?: Record<string, string | IDetailedField>;
}

/**
 * Zod schema for DetailedField.
 * Represents a detailed definition of a field, including its
 * type, description, and nested fields.
 */
const DetailedFieldSchema: ZodType<IDetailedField> = z
  .object({
    /**
     * The type of the field.
     * Must match one of the allowed types defined by the pattern.
     */
    type: TypeRegexSchema,

    /**
     * Optional description of the field.
     */
    desc: z.string().optional().describe("Optional description of the field"),

    /**
     * Optional flag to indicate if the field is optional.
     */
    optional: z.boolean().optional().describe(
      "Optional flag to indicate if the field is optional",
    ),

    /**
     * Optional nested fields within the field definition.
     * Each nested field must adhere to the FieldSchema.
     */
    fields: z
      .record(FieldSchema)
      .optional()
      .describe("Optional nested fields"),
  })
  .strict()
  .superRefine((data, ctx) => {
    const isObjectType = data.type === "object";
    const hasFields = data.fields && Object.keys(data.fields).length > 0;

    if (isObjectType && !hasFields) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "'fields' is required when 'type' is 'object'",
        path: ["fields"],
      });
    }

    if (!isObjectType && hasFields) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "'fields' is not allowed when 'type' is not 'object'",
        path: ["fields"],
      });
    }
  });

/**
 * Zod schema for Type.
 * Represents a custom type with a name, optional description, and a set of fields.
 */
const TypeSchema = z
  .object({
    /**
     * The name of the custom type.
     * Must start with an uppercase letter and contain only alphanumeric characters.
     */
    name: UpperNameSchema,

    /**
     * Optional description of the custom type.
     */
    desc: z.string().optional().describe("Optional description of the type"),

    /**
     * Fields of the custom type.
     * Each field must adhere to the FieldSchema.
     */
    fields: z
      .record(FieldSchema, {
        errorMap: () => ({
          message:
            'Each field in "fields" must be a valid definition according to the schema.',
        }),
      })
      .describe("Fields of the custom type"),
  })
  .strict();

/**
 * Zod schema for Procedure.
 * Represents a procedure with a name, type, optional description, input, output, and metadata.
 */
const ProcedureSchema = z
  .object({
    /**
     * The name of the procedure.
     * Must start with an uppercase letter and contain only alphanumeric characters.
     */
    name: UpperNameSchema,

    /**
     * The type of the procedure.
     * Must be either "query" or "mutation".
     */
    type: z.enum(["query", "mutation"], {
      errorMap: () => ({
        message: 'The "type" field must be either "query" or "mutation".',
      }),
    }),

    /**
     * Optional description of the procedure.
     */
    desc: z.string().optional().describe(
      "Optional description of the procedure",
    ),

    /**
     * Optional input for the procedure.
     * Each input field must adhere to the FieldSchema.
     */
    input: z
      .record(FieldSchema)
      .optional()
      .describe("Optional input for the procedure"),

    /**
     * Optional output for the procedure.
     * Each output field must adhere to the FieldSchema.
     */
    output: z
      .record(FieldSchema)
      .optional()
      .describe("Optional output for the procedure"),

    /**
     * Optional metadata for the procedure.
     * Each metadata property can be a string, number (int/float), or boolean.
     */
    meta: z
      .record(z.union([z.string(), z.number(), z.boolean()]))
      .optional()
      .describe("Optional metadata for the procedure"),
  })
  .strict();

/**
 * Zod schema for the main object.
 * Requires at least one procedure and optionally includes custom types.
 */
const MainSchema = z
  .object({
    /**
     * Optional list of custom types.
     */
    types: z
      .array(TypeSchema)
      .optional()
      .describe("Optional list of custom types"),

    /**
     * List of procedures.
     * Must contain at least one procedure.
     */
    procedures: z
      .array(ProcedureSchema)
      .min(1, {
        message: 'There must be at least one procedure in "procedures".',
      })
      .describe("List of procedures"),
  })
  .describe("UFO RPC Schema")
  .strict();

export default MainSchema;
export type TypeSchemaType = z.infer<typeof TypeSchema>;
export type ProcedureSchemaType = z.infer<typeof ProcedureSchema>;
export type FieldSchemaType = z.infer<typeof FieldSchema>;
