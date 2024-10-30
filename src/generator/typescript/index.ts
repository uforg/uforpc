import Handlebars from "handlebars";
import prettier from "prettier";
import { isArrayType, Schema } from "@/schema/types.ts";
import type { DetailedField, FieldType } from "@/schema/types.ts";

function registerHelpers() {
  Handlebars.registerHelper("tsType", function (type: FieldType): string {
    if (typeof type === "string") {
      if (type === "int" || type === "float") return "number";
      if (/^[A-Z]/.test(type)) return `T${type}`;
      return type;
    }
    if (isArrayType(type)) {
      const baseType = Handlebars.helpers["tsType"](type.baseType);
      return `${baseType}${"[]".repeat(type.dimensions)}`;
    }
    return "unknown";
  });

  Handlebars.registerHelper("httpMethod", function (type: string): string {
    return type === "query" ? "GET" : "POST";
  });

  Handlebars.registerHelper("json", function (value: unknown): string {
    return JSON.stringify(value);
  });

  Handlebars.registerHelper("inferMetaType", function (value: unknown): string {
    if (typeof value === "string") return "string";
    if (typeof value === "number") return "number";
    if (typeof value === "boolean") return "boolean";
    return "unknown";
  });

  Handlebars.registerHelper(
    "renderFields",
    function (fields: Record<string, DetailedField>): string {
      if (!fields) return "";
      let result = "";

      for (const [key, field] of Object.entries(fields)) {
        const optional = field.optional ? "?" : "";

        let fieldType = "";
        if (field.fields) {
          fieldType = "{\n" + Handlebars.helpers["renderFields"](field.fields) +
            "}";
          if (isArrayType(field.type)) {
            fieldType += "[]".repeat(field.type.dimensions);
          }
        } else {
          fieldType = Handlebars.helpers["tsType"](field.type);
        }

        if (field.desc) {
          result += `/** ${field.desc} */\n`;
        }
        result += `${key}${optional}: ${fieldType};\n`;
      }

      return result;
    },
  );

  Handlebars.registerHelper(
    "renderValidationSchemaFields",
    function (fields: Record<string, DetailedField>): string {
      if (!fields) return "";
      let result = "";

      function getBaseSchema(key: string, type: string): string {
        if (type === "int" || type === "float") {
          return `validationSchema.number('${key} must be a number')`;
        }
        if (type === "string") {
          return `validationSchema.string('${key} must be a string')`;
        }
        if (type === "boolean") {
          return `validationSchema.boolean('${key} must be a boolean')`;
        }
        if (/^[A-Z]/.test(type)) {
          return `validationSchema.lazy(() => T${type}ValidationSchema, '${key} must be a ${type}')`;
        }
        return "";
      }

      for (const [key, field] of Object.entries(fields)) {
        let schemaType = "";

        if (field.fields) {
          const nestedFields = Handlebars.helpers
            ["renderValidationSchemaFields"](field.fields);
          schemaType = `validationSchema.object({\n${nestedFields}})`;
        } else if (isArrayType(field.type)) {
          schemaType = getBaseSchema(key, field.type.baseType as string);
          for (let i = 0; i < field.type.dimensions; i++) {
            schemaType = `validationSchema.array(${schemaType})`;
          }
        } else if (typeof field.type === "string") {
          schemaType = getBaseSchema(key, field.type);
        }

        if (!field.optional) {
          schemaType += `.required('${key} is required')`;
        }
        result += `    ${key}: ${schemaType},\n`;
      }

      return result;
    },
  );
}

const coreTypesTemplate = `
// This file has been generated using UFO RPC. DO NOT EDIT.
// If you edit this file, it will be overwritten the next time it is generated

// -----------------------------------------------------------------------------
// Core Types
// -----------------------------------------------------------------------------

/** Represents an HTTP method */
export type UFOHTTPMethod = "GET" | "POST";

/** Represents the output of an error in the UFO RPC system */
export interface UFOErrorOutput {
  message: string;
  details?: Record<string, unknown>;
}

/** Represents the output of a UFO RPC request */
export type UFOResponse<T> =
  | {
    readonly ok: true;
    readonly output: T;
    readonly error?: never;
  }
  | {
    readonly ok: false;
    readonly output?: never;
    readonly error: UFOErrorOutput;
  };

/** Represents an error in the UFO RPC system */
export class UFOError extends Error {
  constructor(message: string, public details?: Record<string, unknown>) {
    super(message);
    this.name = "UFOError";
  }
}

/** Gets the error output for a given error */
function getErrorOutput(err: unknown): UFOErrorOutput {
  if (err instanceof UFOError) {
    return {
      message: err.message,
      details: err.details,
    };
  }

  if (err instanceof Error) {
    return {
      message: err.message,
    };
  }

  return {
    message: "Unknown error",
  };
}
`;

const validationSchemaTemplate = `
  // -----------------------------------------------------------------------------
  // Schema validator
  // -----------------------------------------------------------------------------

  /** Available schema types for validation */
  type ValidationSchemaType =
    | "string"
    | "number"
    | "boolean"
    | "array"
    | "object";

  /** Result of schema validation containing validity status and optional error message */
  type ValidationSchemaResult = {
    isValid: boolean;
    error?: string;
  };

  /**
   * Schema class for type-safe validation
   * @template T - The type of value being validated
   */
  class ValidationSchema<T> {
    private type: ValidationSchemaType;
    private isRequired = false;
    private pattern?: RegExp;
    private arraySchema?: ValidationSchema<unknown>;
    private objectSchema?: Record<string, ValidationSchema<unknown>>;
    private errorMessage?: string;

    /**
     * Creates a new validation schema
     * @param type - The type of value to validate
     * @param errorMessage - Optional custom error message
     */
    constructor(type: ValidationSchemaType, errorMessage?: string) {
      this.type = type;
      this.errorMessage = errorMessage;
    }

    /**
     * Makes the schema required (non-nullable/undefined)
     * @param errorMessage - Optional custom error message for required validation
     */
    required(errorMessage?: string): ValidationSchema<T> {
      this.isRequired = true;
      this.errorMessage = errorMessage || this.errorMessage;
      return this;
    }

    /**
     * Adds regex validation for string schemas
     * @param pattern - RegExp to test against string values
     * @param errorMessage - Optional custom error message for pattern validation
     */
    regex(pattern: RegExp, errorMessage?: string): ValidationSchema<T> {
      if (this.type !== "string") {
        throw new Error("Regex validation only applies to string schemas");
      }
      this.pattern = pattern;
      this.errorMessage = errorMessage || this.errorMessage;
      return this;
    }

    /**
     * Creates an array schema
     * @param schema - Schema for array elements
     * @param errorMessage - Optional custom error message for array validation
     */
    array<U>(
      schema: ValidationSchema<U>,
      errorMessage?: string,
    ): ValidationSchema<U[]> {
      const newSchema = new ValidationSchema<U[]>("array", errorMessage);
      newSchema.arraySchema = schema;
      return newSchema;
    }

    /**
     * Creates an object schema
     * @param schema - Record of property schemas
     * @param errorMessage - Optional custom error message for object validation
     */
    object<U extends Record<string, unknown>>(
      schema: { [K in keyof U]: ValidationSchema<U[K]> },
      errorMessage?: string,
    ): ValidationSchema<U> {
      const newSchema = new ValidationSchema<U>("object", errorMessage);
      newSchema.objectSchema = schema as Record<
        string,
        ValidationSchema<unknown>
      >;
      return newSchema;
    }

    /**
     * Creates a lazy schema for recursive validation
     * @param schema - Function returning the validation schema
     * @param errorMessage - Optional custom error message for lazy validation
     */
    static lazy<T>(
      schema: () => ValidationSchema<T>,
      errorMessage?: string,
    ): ValidationSchema<T> {
      const lazySchema = new ValidationSchema<T>("object", errorMessage);
      lazySchema.validate = (value: unknown): ValidationSchemaResult =>
        schema().validate(value);
      return lazySchema;
    }

    /**
     * Validates a value against the schema
     * @param value - Value to validate
     * @returns Validation result with boolean and optional error message
     */
    validate(value: unknown): ValidationSchemaResult {
      if (value === undefined || value === null) {
        return {
          isValid: !this.isRequired,
          error: this.isRequired
            ? this.errorMessage || "Field is required"
            : undefined,
        };
      }

      if (!this.validateType(value)) {
        return {
          isValid: false,
          error: this.errorMessage || \`Invalid type, expected \${this.type}\`,
        };
      }

      if (this.type === "string" && typeof value === "string") {
        if (this.pattern && !this.pattern.test(value)) {
          return {
            isValid: false,
            error: this.errorMessage || "String does not match pattern",
          };
        }
      }

      if (this.type === "array" && Array.isArray(value)) {
        if (this.arraySchema) {
          for (const item of value) {
            const result = this.arraySchema.validate(item);
            if (!result.isValid) return result;
          }
        }
      }

      if (
        this.type === "object" && typeof value === "object" &&
        !Array.isArray(value)
      ) {
        if (this.objectSchema) {
          for (const [key, schema] of Object.entries(this.objectSchema)) {
            const result = schema.validate(
              (value as Record<string, unknown>)[key],
            );
            if (!result.isValid) return result;
          }
        }
      }

      return { isValid: true };
    }

    /**
     * Validates the type of a value
     * @param value - Value to validate type of
     * @returns Whether the value matches the schema type
     */
    private validateType(value: unknown): boolean {
      switch (this.type) {
        case "string":
          return typeof value === "string";
        case "number":
          return typeof value === "number";
        case "boolean":
          return typeof value === "boolean";
        case "array":
          return Array.isArray(value);
        case "object":
          return typeof value === "object" && !Array.isArray(value);
        default:
          return false;
      }
    }
  }

  /** Factory object for creating validation schemas */
  export const validationSchema = {
    /** Creates a string validation schema
     * @param errorMessage - Optional custom error message
     */
    string: (errorMessage?: string) =>
      new ValidationSchema<string>("string", errorMessage),

    /** Creates a number validation schema
     * @param errorMessage - Optional custom error message
     */
    number: (errorMessage?: string) =>
      new ValidationSchema<number>("number", errorMessage),

    /** Creates a boolean validation schema
     * @param errorMessage - Optional custom error message
     */
    boolean: (errorMessage?: string) =>
      new ValidationSchema<boolean>("boolean", errorMessage),

    /** Creates an array validation schema
     * @param schema - Schema for array elements
     * @param errorMessage - Optional custom error message
     */
    array: <T>(schema: ValidationSchema<T>, errorMessage?: string) =>
      new ValidationSchema<T[]>("array", errorMessage).array(schema),

    /** Creates an object validation schema
     * @param schema - Record of property schemas
     * @param errorMessage - Optional custom error message
     */
    object: <T extends Record<string, unknown>>(
      schema: { [K in keyof T]: ValidationSchema<T[K]> },
      errorMessage?: string,
    ) => new ValidationSchema<T>("object", errorMessage).object(schema),

    /** Creates a lazy validation schema for recursive validation
     * @param schema - Function returning the validation schema
     * @param errorMessage - Optional custom error message
     */
    lazy: <T>(
      schema: () => ValidationSchema<T>,
      errorMessage?: string,
    ): ValidationSchema<T> => ValidationSchema.lazy(schema, errorMessage),
  };
`;

const domainTypesTemplate = `
// -----------------------------------------------------------------------------
// Domain Types
// -----------------------------------------------------------------------------

{{#each types}}

{{#if desc}}
/** {{desc}} */
{{/if}}
export interface T{{name}} {
  {{renderFields fields}}
}

/** Schema to validate the **T{{name}}** custom type. */
const T{{name}}ValidationSchema = validationSchema.object({
  {{renderValidationSchemaFields fields}}
})

{{/each}}`;

const procedureTypesTemplate = `
// -----------------------------------------------------------------------------
// Procedure Types
// -----------------------------------------------------------------------------

{{#each procedures}}

/** Represents the input for the **{{name}}** procedure. */
{{#if input}}
export interface P{{name}}Input {
  {{renderFields input}}
}
{{else}}
export type P{{name}}Input = never;
{{/if}}

{{#if input}}
/** Schema to validate the input for the **{{name}}** procedure. */
const P{{name}}InputValidationSchema = validationSchema.object({
  {{renderValidationSchemaFields input}}
})
{{/if}}

/** Represents the output for the **{{name}}** procedure. */
{{#if output}}
export interface P{{name}}Output {
  {{renderFields output}}
}
{{else}}
export type P{{name}}Output = never;
{{/if}}

/** Represents the metadata for the **{{name}}** procedure. */
{{#if meta}}
export interface P{{name}}Meta {
  {{#each meta}}
    {{@key}}: {{inferMetaType this}};
  {{/each}}
}
{{else}}
export type P{{name}}Meta = never;
{{/if}}

{{/each}}

/** All validation schemas for procedures */
const AllValidationSchemas = {
  {{#each procedures}}
    {{name}}: {
      {{#if input}}
        hasValidationSchema: true,
        validationSchema: P{{name}}InputValidationSchema,
      {{else}}
        hasValidationSchema: false,
      {{/if}}
    },
  {{/each}}
};

/** Types for all procedures */
export interface UFOProcedures {
  {{#each procedures}}
    {{name}}: {
      type: "{{type}}";
      input: P{{name}}Input;
      output: P{{name}}Output;
      meta: P{{name}}Meta;
    };
  {{/each}}
}

/** Names for all procedures */
export type UFOProcedureNames = keyof UFOProcedures;`;

function createServerTemplate(opts: GenerateTypescriptOpts): string {
  if (!opts.includeServer) return "";

  let validationLogic = "const isValid = true;";
  if (!opts.omitServerRequestValidation) {
    validationLogic = `
      let isValid = true;
      const valSchema = AllValidationSchemas[procedureName];
      if (valSchema.hasValidationSchema) {
        const valRes = valSchema.validationSchema.validate(request.input);
        isValid = valRes.isValid;
        if (!isValid) {
          response = {
            ok: false,
            error: {
              message: valRes.error ?? "Invalid input",
            },
          };
        }
      }
    `;
  }

  return `
    // -----------------------------------------------------------------------------
    // Server Types
    // -----------------------------------------------------------------------------

    export interface UFOServerRPCRequest<UFORequestContext> {
      readonly procedure: UFOProcedureNames | string;
      readonly method: UFOHTTPMethod;
      readonly input: UFOProcedures[UFOProcedureNames]["input"];
      readonly context: UFORequestContext;
    }

    export interface UFOServerProcedureContext<
      UFOProcedureNames,
      TInput,
      TMeta,
      UFORequestContext
    > {
      readonly procedure: UFOProcedureNames;
      readonly input: TInput;
      readonly meta: TMeta;
      readonly context: UFORequestContext;
    }

    {{#each procedures}}

    {{#if desc}}
    /** {{desc}} */
    {{/if}}
    export type P{{name}}Handler<UFORequestContext> = (
      ctx: UFOServerProcedureContext<
        "{{name}}",
        P{{name}}Input,
        P{{name}}Meta,
        UFORequestContext
      >
    ) => Promise<P{{name}}Output>;

    {{/each}}

    export interface UFOServerMiddleware<UFORequestContext> {
      before?(context: UFORequestContext): Promise<UFORequestContext>;
      after?(
        context: UFORequestContext,
        response: UFOResponse<UFOProcedures[UFOProcedureNames]["output"]>,
      ): Promise<typeof response>;
    }

    // -----------------------------------------------------------------------------
    // Server Implementation
    // -----------------------------------------------------------------------------

    export class UFOServer<UFORequestContext> {
      private readonly handlers = new Map<
        UFOProcedureNames,
        (
          ctx: UFOServerProcedureContext<
            UFOProcedureNames,
            unknown,
            unknown,
            UFORequestContext
          >,
        ) => Promise<unknown>
      >();
      private readonly middleware: UFOServerMiddleware<UFORequestContext>[] = [];

      private readonly methodMap: Record<UFOProcedureNames, UFOHTTPMethod> = {
        {{#each procedures}}
          {{name}}: "{{httpMethod type}}"{{#unless @last}},{{/unless}}
        {{/each}}
      };

      private readonly metaMap: Partial<Record<UFOProcedureNames, unknown>> = {
        {{#each procedures}}
          {{#if meta}}
            {{name}}: {{{json meta}}}{{#unless @last}},{{/unless}}
          {{/if}}
        {{/each}}
      };

      defineHandler<P extends UFOProcedureNames>(
        procedure: P,
        handler: (
          ctx: UFOServerProcedureContext<
            P,
            UFOProcedures[P]["input"],
            UFOProcedures[P]["meta"],
            UFORequestContext
          >
        ) => Promise<UFOProcedures[P]["output"]>
      ): this {
        this.handlers.set(
          procedure,
          handler as (
            ctx: UFOServerProcedureContext<
              UFOProcedureNames,
              unknown,
              unknown,
              UFORequestContext
            >
          ) => Promise<unknown>
        );
        return this;
      }

      defineMiddleware(middleware: UFOServerMiddleware<UFORequestContext>): this {
        this.middleware.push(middleware);
        return this;
      }

      async handleRequest(
        request: UFOServerRPCRequest<UFORequestContext>
      ): Promise<UFOResponse<UFOProcedures[UFOProcedureNames]["output"]>> {
        type ProcedureOutput = UFOProcedures[UFOProcedureNames]["output"];
        const procedureName = request.procedure as UFOProcedureNames;

        if (!procedureName) {
          return {
            ok: false,
            error: {
              message: "Procedure not defined"
            }
          };
        }

        const handler = this.handlers.get(procedureName);
        if (!handler) {
          return {
            ok: false,
            error: {
              message: \`Handler not defined for procedure \${request.procedure}\`
            }
          };
        }

        const expectedMethod = this.methodMap[procedureName];
        if (request.method !== expectedMethod) {
          return {
            ok: false,
            error: {
              message: \`Method \${request.method} not allowed for \${request.procedure} procedure\`
            }
          };
        }

        try {
          let currentUFORequestContext = request.context;

          for await (const m of this.middleware) {
            if (m.before) {
              currentUFORequestContext = await m.before(currentUFORequestContext);
            }
          }

          let response: UFOResponse<UFOProcedures[UFOProcedureNames]["output"]> = {
            ok: false,
            error: {
              message: "Unknown error"
            }
          }
          
          ${validationLogic}

          if (isValid) {
            try {
              const output = (await handler({
                procedure: procedureName,
                input: request.input,
                meta: this.metaMap[procedureName],
                context: currentUFORequestContext,
              })) as ProcedureOutput;
              response = {
                ok: true,
                output,
              };
            } catch (err) {
              response = {
                ok: false,
                error: getErrorOutput(err),
              };
            }
          }

          for await (const m of this.middleware) {
            if (m.after) {
              response = (await m.after(
                currentUFORequestContext,
                response,
              )) as typeof response;
            }
          }

          return response
        } catch (err) {
          return {
            ok: false,
            error: getErrorOutput(err)
          };
        }
      }
    }
  `;
}

function createClientTemplate(opts: GenerateTypescriptOpts): string {
  if (!opts.includeClient) return "";

  const emitFetch = !opts.omitClientDefaultFetch;
  let fetchClient = "";
  if (emitFetch) {
    fetchClient = `
      /** Default UFO RPC Fetch HTTP Client Implementation */
      export class UFOFetchClient implements UFOHTTPClient {
        constructor(
          private readonly fetch: typeof globalThis.fetch = globalThis.fetch
        ) {}

        async request<T>(request: UFOClientHTTPRequest): Promise<UFOResponse<T>> {
          const options: RequestInit = {
            method: request.method,
            headers: request.headers,
          };

          if (request.body) {
            options.body = JSON.stringify(request.body);
          }

          const response = await this.fetch(request.url, options);
          const data = await response.json();

          if (typeof data.ok === "boolean" && (data.output || data.error)) {
            return data;
          }

          return {
            ok: false,
            error: {
              message: "Invalid response from server",
              details: {
                status: response.status,
                statusText: response.statusText,
                data,
              },
            },
          }
        }
      }
    `;
  }

  let validationLogic = "const isValid = true;";
  if (!opts.omitServerRequestValidation) {
    validationLogic = `
      let isValid = true;
      const valSchema = AllValidationSchemas[procedure];
      if (valSchema.hasValidationSchema) {
        const valRes = valSchema.validationSchema.validate(input);
        isValid = valRes.isValid;
        if (!isValid) {
          response = {
            ok: false,
            error: {
              message: valRes.error ?? "Invalid input",
            },
          };
        }
      }
    `;
  }

  let httpClientInit = "this.httpClient = config.httpClient;";
  if (emitFetch) {
    httpClientInit =
      "this.httpClient = config.httpClient ?? new UFOFetchClient();";
  }

  const template = `
    // -----------------------------------------------------------------------------
    // Client Implementation
    // -----------------------------------------------------------------------------

    export interface UFOClientHTTPRequest {
      url: string;
      method: UFOHTTPMethod;
      body?: unknown;
      headers?: Record<string, string>;
    }

    export interface UFOHTTPClient {
      request<T>(request: UFOClientHTTPRequest): Promise<UFOResponse<T>>;
    }

    export interface UFOClientMiddleware {
      before?(request: UFOClientHTTPRequest): Promise<UFOClientHTTPRequest>;
      after?(response: UFOResponse<UFOProcedures[UFOProcedureNames]["output"]>): Promise<typeof response>;
    }

    ${fetchClient}

    export interface UFOClientConfig {
      baseUrl: string;
      ${emitFetch ? "httpClient?" : "httpClient"}: UFOHTTPClient;
    }

    export class UFOClient {
      private readonly middleware: UFOClientMiddleware[] = [];
      private readonly httpClient: UFOHTTPClient;

      constructor(private readonly config: UFOClientConfig) {
        ${httpClientInit}
      }

      defineMiddleware(middleware: UFOClientMiddleware): this {
        this.middleware.push(middleware);
        return this;
      }

      private async request<P extends UFOProcedureNames>(
        procedure: P,
        method: UFOHTTPMethod,
        input: UFOProcedures[P]["input"],
      ): Promise<UFOResponse<UFOProcedures[P]["output"]>> {
        let request: UFOClientHTTPRequest = {
          url: this.buildUrl(procedure, method, input),
          method,
        };
        if (method === "POST") {
          request.headers = { "Content-Type": "application/json" };
          request.body = input;
        }

        try {
          for await (const m of this.middleware) {
            if (m.before) request = await m.before(request);
          }

          let response: UFOResponse<UFOProcedures[P]["output"]> = {
            ok: false,
            error: {
              message: "Unknown error"
            }
          };

          ${validationLogic}

          if (isValid) {
            try {
              response = await this.httpClient.request<UFOProcedures[P]["output"]>(
                request,
              );
            } catch (err) {
              response = {
                ok: false,
                error: getErrorOutput(err),
              };
            }
          }

          for await (const m of this.middleware) {
            if (m.after) {
              response = (await m.after(response)) as UFOResponse<
                UFOProcedures[P]["output"]
              >;
            }
          }

          return response;
        } catch (err) {
          return {
            ok: false,
            error: getErrorOutput(err),
          };
        }
      }

      private buildUrl(
        procedure: string,
        method: UFOHTTPMethod,
        input: unknown
      ): string {
        const url = new URL(\`\${this.config.baseUrl}/\${procedure}\`);

        if (method === "GET") {
          url.searchParams.append("input", JSON.stringify(input));
        }

        return url.toString();
      }

      {{#each procedures}}
      {{name}} = (input: P{{name}}Input) =>
        this.request("{{name}}", "{{httpMethod type}}", input);
      {{/each}}
    }`;

  return template;
}

/**
 * Compiles a Handlebars template with proper configuration for TypeScript generation
 */
function compileTemplate(template: string): HandlebarsTemplateDelegate {
  return Handlebars.compile(template, {
    noEscape: true,
    strict: true,
  });
}

/**
 * Formats an string of typescript code
 */
async function formatCode(code: string): Promise<string> {
  return await prettier.format(code, { parser: "typescript" });
}

export interface GenerateTypescriptOpts {
  includeServer?: boolean;
  includeClient?: boolean;
  omitServerRequestValidation?: boolean;
  omitClientRequestValidation?: boolean;
  omitClientDefaultFetch?: boolean;
}

/**
 * Generates TypeScript code from a UFO RPC schema, including types, server, and client implementations
 */
export async function generateTypeScript(
  schema: Schema,
  opts: GenerateTypescriptOpts,
): Promise<string> {
  registerHelpers();

  const templates = [
    coreTypesTemplate,
    validationSchemaTemplate,
    domainTypesTemplate,
    procedureTypesTemplate,
    createServerTemplate(opts),
    createClientTemplate(opts),
  ];

  const compiled = templates.map(compileTemplate);
  const generated = compiled.map((template) => template(schema)).join("\n\n");

  const formatted = await formatCode(generated);

  return formatted;
}
