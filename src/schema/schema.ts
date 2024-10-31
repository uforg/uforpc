/**
 * ⚠️ IMPORTANT ⚠️
 * If you modify this file, make sure to run the following command in order
 * to auto-generate the corresponding json schema from this zod schema:
 * deno task gen:schema
 */

import { z } from "zod";
import type { IssueData, ZodType } from "zod";

/**
 * Type definitions for validation rules
 */
type ValidationRule =
  | "required"
  | "regex"
  | "equals"
  | "contains"
  | "length"
  | "minLength"
  | "maxLength"
  | "enum"
  | "email"
  | "iso8601"
  | "json"
  | "lowercase"
  | "uppercase"
  | "min"
  | "max";

type RuleConfig = {
  string: ValidationRule[];
  int: ValidationRule[];
  float: ValidationRule[];
  boolean: ValidationRule[];
  object: ValidationRule[];
};

/**
 * Valid rule configurations per field type
 */
const VALID_RULES: RuleConfig = {
  string: [
    "required",
    "regex",
    "equals",
    "contains",
    "length",
    "minLength",
    "maxLength",
    "enum",
    "email",
    "iso8601",
    "json",
    "lowercase",
    "uppercase",
  ],
  int: [
    "required",
    "equals",
    "min",
    "max",
    "enum",
  ],
  float: [
    "required",
    "equals",
    "min",
    "max",
  ],
  boolean: [
    "required",
    "equals",
  ],
  object: [
    "required",
    "json",
  ],
};

/**
 * Schema for validation rule types
 */
const ValidationRuleSchema = z.discriminatedUnion("rule", [
  // Required validation
  z.object({
    rule: z.literal("required"),
    message: z.string().optional(),
  }),

  // String pattern validations
  z.object({
    rule: z.literal("regex"),
    pattern: z.string(),
    message: z.string().optional(),
  }),
  z.object({
    rule: z.literal("equals"),
    value: z.union([z.string(), z.number(), z.boolean()]),
    message: z.string().optional(),
  }),
  z.object({
    rule: z.literal("contains"),
    value: z.string(),
    message: z.string().optional(),
  }),

  // Format validations
  z.object({
    rule: z.literal("email"),
    message: z.string().optional(),
  }),
  z.object({
    rule: z.literal("iso8601"),
    message: z.string().optional(),
  }),
  z.object({
    rule: z.literal("json"),
    message: z.string().optional(),
  }),
  z.object({
    rule: z.literal("lowercase"),
    message: z.string().optional(),
  }),
  z.object({
    rule: z.literal("uppercase"),
    message: z.string().optional(),
  }),

  // Numeric validations
  z.object({
    rule: z.literal("min"),
    value: z.number(),
    message: z.string().optional(),
  }),
  z.object({
    rule: z.literal("max"),
    value: z.number(),
    message: z.string().optional(),
  }),

  // Length validations
  z.object({
    rule: z.literal("length"),
    value: z.number(),
    message: z.string().optional(),
  }),
  z.object({
    rule: z.literal("minLength"),
    value: z.number(),
    message: z.string().optional(),
  }),
  z.object({
    rule: z.literal("maxLength"),
    value: z.number(),
    message: z.string().optional(),
  }),

  // Enumeration validation
  z.object({
    rule: z.literal("enum"),
    values: z.array(z.union([z.string(), z.number()])),
    message: z.string().optional(),
  }),
]);

const TypeRegexSchema = z.string().regex(
  /^(string|int|float|boolean|object|.*\[\]|[A-Z][a-zA-Z0-9]*)$/,
  {
    message:
      'The "type" field must be one of the allowed types: string, int, float, boolean, object, an array (e.g., string[] or MyType[][]), or a custom type name starting with an uppercase letter.',
  },
);

const UpperNameSchema = z.string().regex(
  /^[A-Z][a-zA-Z0-9]*$/,
  {
    message:
      'The "name" field must start with an uppercase letter and contain only alphanumeric characters.',
  },
);

/**
 * Validation core types and utilities
 */
interface IField {
  type: string;
  desc?: string;
  rules?: z.infer<typeof ValidationRuleSchema>[];
  fields?: Record<string, IField>;
}

class FieldValidator {
  constructor(private field: IField) {}

  private createIssue(message: string, path?: (string | number)[]): IssueData {
    return {
      code: z.ZodIssueCode.custom,
      message,
      path,
    };
  }

  private getBaseType(): string {
    return this.field.type.replaceAll("[]", "");
  }

  private getValidRules(baseType: string): ValidationRule[] {
    return (baseType in VALID_RULES
      ? VALID_RULES[baseType as keyof RuleConfig]
      : VALID_RULES.string);
  }

  validateObjectFields(): IssueData[] {
    const issues: IssueData[] = [];
    const isObjectType = this.field.type.startsWith("object");
    const hasFields = this.field.fields &&
      Object.keys(this.field.fields).length > 0;

    if (isObjectType && !hasFields) {
      issues.push(this.createIssue(
        "'fields' is required when 'type' is 'object'",
        ["fields"],
      ));
    }

    if (!isObjectType && hasFields) {
      issues.push(this.createIssue(
        "'fields' is not allowed when 'type' is not 'object'",
        ["fields"],
      ));
    }

    return issues;
  }

  validateRuleTypeCompatibility(
    rule: z.infer<typeof ValidationRuleSchema>,
    ruleIndex: number,
  ): IssueData[] {
    const baseType = this.getBaseType();
    const validRules = this.getValidRules(baseType);

    if (!validRules.includes(rule.rule as ValidationRule)) {
      return [this.createIssue(
        `The rule "${rule.rule}" is not valid for type "${this.field.type}". Valid rules are: ${
          validRules.join(", ")
        }`,
        ["rules", ruleIndex],
      )];
    }

    return [];
  }

  validateEqualsRule(
    rule: z.infer<typeof ValidationRuleSchema>,
    ruleIndex: number,
  ): IssueData[] {
    if (rule.rule !== "equals") return [];

    const baseType = this.getBaseType();
    const valueType = typeof rule.value;
    const typeValidations: Record<
      string,
      { expected: string; check: () => boolean }
    > = {
      string: {
        expected: "string",
        check: () => valueType === "string",
      },
      int: {
        expected: "number",
        check: () => valueType === "number",
      },
      float: {
        expected: "number",
        check: () => valueType === "number",
      },
      boolean: {
        expected: "boolean",
        check: () => valueType === "boolean",
      },
    };

    const validation = typeValidations[baseType];
    if (validation && !validation.check()) {
      return [this.createIssue(
        `equals rule for ${baseType} type must have a ${validation.expected} value`,
        ["rules", ruleIndex, "value"],
      )];
    }

    return [];
  }

  validate(): IssueData[] {
    const issues: IssueData[] = [
      ...this.validateObjectFields(),
    ];

    if (this.field.rules) {
      this.field.rules.forEach((rule, index) => {
        issues.push(
          ...this.validateRuleTypeCompatibility(rule, index),
          ...this.validateEqualsRule(rule, index),
        );
      });
    }

    return issues;
  }
}

/**
 * Schema definitions
 */
const FieldSchema: ZodType<IField> = z
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
    rules: z.array(ValidationRuleSchema).optional()
      .describe("Validation rules for the field"),
    fields: z
      .record(z.lazy(() => FieldSchema))
      .optional()
      .describe("Optional nested fields"),
  })
  .strict()
  .superRefine((data, ctx) => {
    const validator = new FieldValidator(data);
    const issues = validator.validate();

    for (const issue of issues) {
      ctx.addIssue(issue);
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
     * Must adhere to the FieldSchema.
     * If not provided, the procedure does not accept any input.
     */
    input: z.record(FieldSchema).optional().describe(
      "Optional input for the procedure",
    ),

    /**
     * Optional output for the procedure.
     * Must adhere to the FieldSchema.
     * If not provided, the procedure does not return any output.
     */
    output: z.record(FieldSchema).optional().describe(
      "Optional output for the procedure",
    ),

    /**
     * Optional metadata for the procedure.
     * Must be a record with string, number, or boolean values.
     */
    meta: z
      .record(z.union([z.string(), z.number(), z.boolean()]))
      .optional()
      .describe("Optional metadata for the procedure"),
  })
  .strict();

/**
 * Zod schema for the main UFO RPC schema.
 * Contains a list of custom types and procedures.
 * At least one procedure must be defined.
 * Custom types are optional.
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
export type ValidationRuleType = z.infer<typeof ValidationRuleSchema>;
