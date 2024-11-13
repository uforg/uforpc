import path from "node:path";
import Handlebars from "handlebars";
import { isArrayType, parseArrayType } from "@/schema/index.ts";
import type { FieldSchemaType, MainSchemaType } from "@/schema/index.ts";
import { handlebarsCompileTemplate } from "../common/handlebars-compile-template.ts";
import { formatTsCode } from "@/generator/typescript/format-ts-code.ts";

function registerHelpers() {
  Handlebars.registerHelper("tsType", function (type: FieldSchemaType): string {
    if (typeof type === "string") {
      if (type === "int" || type === "float") return "number";
      if (/^[A-Z]/.test(type)) return `T${type}`;
      return type;
    }
    if (isArrayType(type)) {
      const parsed = parseArrayType(type);
      const baseType = Handlebars.helpers["tsType"](parsed.type.type);
      return `${baseType}${"[]".repeat(parsed.dimensions)}`;
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
    function (fields: Record<string, FieldSchemaType>): string {
      if (!fields) return "";
      let result = "";

      for (const [key, field] of Object.entries(fields)) {
        const isRequired = field.rules?.some((rule) => {
          return rule.rule === "required";
        });

        const optional = isRequired ? "" : "?";

        let fieldType = "";
        if (field.fields) {
          fieldType = "{\n" + Handlebars.helpers["renderFields"](field.fields) +
            "}";
          if (isArrayType(field)) {
            const parsed = parseArrayType(field);
            fieldType += "[]".repeat(parsed.dimensions);
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
    function (fields: Record<string, FieldSchemaType>): string {
      if (!fields) return "";
      let result = "";

      function getBaseSchema(key: string, type: string): string {
        if (type === "int" || type === "float") {
          return `schValidator.number('${key} must be a number')`;
        }
        if (type === "string") {
          return `schValidator.string('${key} must be a string')`;
        }
        if (type === "boolean") {
          return `schValidator.boolean('${key} must be a boolean')`;
        }
        if (/^[A-Z]/.test(type)) {
          return `schValidator.lazy(() => T${type}ValidationSchema, '${key} must be a ${type}')`;
        }
        return "";
      }

      for (const [key, field] of Object.entries(fields)) {
        let schemaType = "";

        if (field.fields) {
          const nestedFields = Handlebars.helpers
            ["renderValidationSchemaFields"](field.fields);
          schemaType = `schValidator.object({\n${nestedFields}})`;
        } else if (isArrayType(field)) {
          const parsed = parseArrayType(field);
          schemaType = getBaseSchema(key, parsed.type.type);
          for (let i = 0; i < parsed.dimensions; i++) {
            schemaType = `schValidator.array(${schemaType})`;
          }
        } else if (typeof field.type === "string") {
          schemaType = getBaseSchema(key, field.type);
        }

        for (const rule of field.rules || []) {
          switch (rule.rule) {
            case "required": {
              const msg = rule.message || `${key} is required`;
              schemaType += `.required('${msg}')`;
              break;
            }
            case "regex": {
              const msg = rule.message || `${key} must match ${rule.pattern}`;
              schemaType += `.regex(${rule.pattern}, '${msg}')`;
              break;
            }
            case "contains": {
              const msg = rule.message || `${key} must contain ${rule.value}`;
              schemaType += `.contains('${rule.value}', '${msg}')`;
              break;
            }
            case "equals": {
              const msg = rule.message || `${key} must equal to ${rule.value}`;
              schemaType += `.equals('${rule.value}', '${msg}')`;
              break;
            }
            case "enum": {
              const msg = rule.message ||
                `${key} must be one of ${rule.values.join(", ")}`;
              schemaType += `.enum(${JSON.stringify(rule.values)}, '${msg}')`;
              break;
            }
            case "email": {
              const msg = rule.message || `${key} must be a valid email`;
              schemaType += `.email('${msg}')`;
              break;
            }
            case "iso8601": {
              const msg = rule.message ||
                `${key} must be a valid ISO 8601 date`;
              schemaType += `.iso8601('${msg}')`;
              break;
            }
            case "uuid": {
              const msg = rule.message || `${key} must be a valid UUID`;
              schemaType += `.uuid('${msg}')`;
              break;
            }
            case "json": {
              const msg = rule.message || `${key} must be a valid JSON string`;
              schemaType += `.json('${msg}')`;
              break;
            }
            case "length": {
              const msg = rule.message ||
                `${key} must have length ${rule.value}`;
              schemaType += `.length(${rule.value}, '${msg}')`;
              break;
            }
            case "minLength": {
              const msg = rule.message ||
                `${key} must have a minimum length of ${rule.value}`;
              schemaType += `.minLength(${rule.value}, '${msg}')`;
              break;
            }
            case "maxLength": {
              const msg = rule.message ||
                `${key} must have a maximum length of ${rule.value}`;
              schemaType += `.maxLength(${rule.value}, '${msg}')`;
              break;
            }
            case "lowercase": {
              const msg = rule.message || `${key} must be lowercase`;
              schemaType += `.lowercase('${msg}')`;
              break;
            }
            case "uppercase": {
              const msg = rule.message || `${key} must be uppercase`;
              schemaType += `.uppercase('${msg}')`;
              break;
            }
            case "min": {
              const msg = rule.message ||
                `${key} must be greater than ${rule.value}`;
              schemaType += `.min(${rule.value}, '${msg}')`;
              break;
            }
            case "max": {
              const msg = rule.message ||
                `${key} must be less than ${rule.value}`;
              schemaType += `.max(${rule.value}, '${msg}')`;
              break;
            }
          }
        }

        result += `    ${key}: ${schemaType},\n`;
      }

      return result;
    },
  );
}

function createCoreTypesTemplate() {
  return `
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
      }
      | {
        readonly ok: false;
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
}

function createDomainTypesTemplate() {
  return `
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

    {{/each}}
  `;
}

function createProcedureTypesTemplate() {
  return `
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
    export type P{{name}}Input = void;
    {{/if}}

    /** Represents the output for the **{{name}}** procedure. */
    {{#if output}}
    export interface P{{name}}Output {
      {{renderFields output}}
    }
    {{else}}
    export type P{{name}}Output = void;
    {{/if}}

    /** Represents the metadata for the **{{name}}** procedure. */
    {{#if meta}}
    export interface P{{name}}Meta {
      {{#each meta}}
        {{@key}}: {{inferMetaType this}};
      {{/each}}
    }
    {{else}}
    export type P{{name}}Meta = void;
    {{/if}}

    {{/each}}

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
    export type UFOProcedureNames = keyof UFOProcedures;
  `;
}

function createValidationSchemaTemplate(opts: GenerateTypescriptOpts): string {
  if (!opts.includeClient && !opts.includeServer) return "";
  if (opts.omitClientRequestValidation && opts.omitServerRequestValidation) {
    return "";
  }

  let validatorContent = Deno.readTextFileSync(path.join(
    import.meta.dirname || "./",
    "./validator/validator.ts",
  ));

  validatorContent = validatorContent.replaceAll("export class", "class");
  validatorContent = validatorContent.replaceAll("export const", "const");

  const template = `
    {{#each types}}

      /** Schema to validate the **T{{name}}** custom type. */
      const T{{name}}ValidationSchema = schValidator.object({
        {{renderValidationSchemaFields fields}}
      })

    {{/each}}

    {{#each procedures}}

      {{#if input}}
        /** Schema to validate the input for the **{{name}}** procedure. */
        const P{{name}}InputValidationSchema = schValidator.object({
          {{renderValidationSchemaFields input}}
        })
      {{/if}}

    {{/each}}

    /** All validation schemas for procedures */
    const AllValidationSchemas: Record<
      string,
      {hasSchValidator: true, schValidator: SchemaValidator<unknown>} |
      {hasSchValidator: false}
    > = {
      {{#each procedures}}
        {{name}}: {
          {{#if input}}
            hasSchValidator: true,
            schValidator: P{{name}}InputValidationSchema,
          {{else}}
            hasSchValidator: false,
          {{/if}}
        },
      {{/each}}
    };
  `;

  return `${validatorContent}\n\n${template}`;
}

function createServerTemplate(opts: GenerateTypescriptOpts): string {
  if (!opts.includeServer) return "";

  let validationLogic = "const isValid = true;";
  if (!opts.omitServerRequestValidation) {
    validationLogic = `
      let isValid = true;
      const valSchema = AllValidationSchemas[procedureName];
      if (valSchema?.hasSchValidator && valSchema?.schValidator) {
        const valRes = valSchema.schValidator.validate(request.input);
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
      readonly method: UFOHTTPMethod;
      readonly context: UFORequestContext;
      readonly procedure: string;
      readonly input: unknown;
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
      if (valSchema?.hasSchValidator && valSchema?.schValidator) {
        const valRes = valSchema.schValidator.validate(input);
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

        if (method === "GET" && input) {
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

export interface GenerateTypescriptOpts {
  includeServer?: boolean;
  includeClient?: boolean;
  omitServerRequestValidation?: boolean;
  omitClientRequestValidation?: boolean;
  omitClientDefaultFetch?: boolean;
}

/**
 * Generates TypeScript code from a UFO RPC schema, including types, server, and
 * client implementations
 */
export async function generateTypeScript(
  schema: MainSchemaType,
  opts: GenerateTypescriptOpts,
): Promise<string> {
  registerHelpers();

  const templates = [
    createCoreTypesTemplate(),
    createDomainTypesTemplate(),
    createProcedureTypesTemplate(),
    createValidationSchemaTemplate(opts),
    createServerTemplate(opts),
    createClientTemplate(opts),
  ];

  const compiled = templates.map(handlebarsCompileTemplate);
  const generated = compiled.map((template) => template(schema)).join("\n\n");

  const formatted = await formatTsCode(generated);

  return formatted;
}
