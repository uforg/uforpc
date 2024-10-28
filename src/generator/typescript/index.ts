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
}

const coreTypesTemplate = `
// This file has been generated using UFO RPC. DO NOT EDIT.
// If you edit this file, it will be overwritten the next time it is generated

// -----------------------------------------------------------------------------
// Core Types
// -----------------------------------------------------------------------------

export type UFOHTTPMethod = "GET" | "POST";

export class UFOError extends Error {
  constructor(message: string, public details?: Record<string, unknown>) {
    super(message);
    this.name = "UFOError";
  }
}`;

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

{{/each}}`;

const procedureTypesTemplate = `
// -----------------------------------------------------------------------------
// Procedure Types
// -----------------------------------------------------------------------------

{{#each procedures}}

{{#if input}}
/** Represents the input for the **{{name}}** procedure. */
export interface P{{name}}Input {
  {{renderFields input}}
}
{{/if}}

{{#if output}}
/** Represents the output for the **{{name}}** procedure. */
export interface P{{name}}Output {
  {{renderFields output}}
}
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

/** Unified types for all procedures */
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

export type UFOProcedureNames = keyof UFOProcedures;`;

const serverTemplate = `
// -----------------------------------------------------------------------------
// Server Types
// -----------------------------------------------------------------------------

export interface UFOProcedureContext<TInput, TMeta, UFORequestContext> {
  readonly input: TInput;
  readonly meta: TMeta;
  readonly context: UFORequestContext;
}

{{#each procedures}}

{{#if desc}}
/** {{desc}} */
{{/if}}
export type P{{name}}Handler<UFORequestContext> = (
  ctx: UFOProcedureContext<P{{name}}Input, P{{name}}Meta, UFORequestContext>
) => Promise<P{{name}}Output>;

{{/each}}

export interface UFOServerMiddleware<UFORequestContext> {
  before?(context: UFORequestContext): Promise<UFORequestContext>;
  after?(context: UFORequestContext, result: unknown): Promise<unknown>;
  error?(context: UFORequestContext, error: Error): Promise<void>;
}

// -----------------------------------------------------------------------------
// Server Implementation
// -----------------------------------------------------------------------------

export class UFOServer<UFORequestContext> {
  private readonly handlers = new Map<
    UFOProcedureNames,
    (
      ctx: UFOProcedureContext<unknown, unknown, UFORequestContext>
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
      ctx: UFOProcedureContext<
        UFOProcedures[P]["input"],
        UFOProcedures[P]["meta"],
        UFORequestContext
      >
    ) => Promise<UFOProcedures[P]["output"]>
  ): this {
    this.handlers.set(
      procedure,
      handler as (
        ctx: UFOProcedureContext<unknown, unknown, UFORequestContext>
      ) => Promise<unknown>
    );
    return this;
  }

  defineMiddleware(middleware: UFOServerMiddleware<UFORequestContext>): this {
    this.middleware.push(middleware);
    return this;
  }

  async handleRequest<P extends UFOProcedureNames>(
    procedure: P,
    method: UFOHTTPMethod,
    input: UFOProcedures[P]["input"],
    context: UFORequestContext
  ): Promise<UFOProcedures[P]["output"]> {
    const expectedMethod = this.methodMap[procedure];
    if (method !== expectedMethod) {
      throw new UFOError(
        \`\${procedure} requires \${expectedMethod} method, got \${method}\`
      );
    }

    const handler = this.handlers.get(procedure);
    if (!handler) {
      throw new UFOError(\`Handler not defined for procedure: \${procedure}\`);
    }

    try {
      let currentUFORequestContext = context;

      for (const m of this.middleware) {
        if (m.before) {
          currentUFORequestContext = await m.before(currentUFORequestContext);
        }
      }

      let result = await handler({
        input,
        meta: this.metaMap[procedure] as UFOProcedures[P]["meta"],
        context: currentUFORequestContext,
      });

      for (const m of this.middleware) {
        if (m.after) {
          result = await m.after(currentUFORequestContext, result);
        }
      }

      return result as UFOProcedures[P]["output"];
    } catch (err) {
      const error = err instanceof UFOError
        ? err
        : new UFOError(err instanceof Error ? err.message : "Unknown error");

      for (const m of this.middleware) {
        if (m.error) await m.error(context, error);
      }

      throw error;
    }
  }
}`;

function createClientTemplate(opts: GenerateTypescriptOpts): string {
  const emitFetch = !opts.omitClientDefaultFetch;

  let fetchClient = "";
  if (emitFetch) {
    fetchClient = `
      /** Default UFO RPC Fetch HTTP Client Implementation */
      export class UFOFetchClient implements UFOHTTPClient {
        constructor(
          private readonly fetch: typeof globalThis.fetch = globalThis.fetch
        ) {}

        async request<T>(request: UFOCientRequest): Promise<UFOClientResponse<T>> {
          const options: RequestInit = {
            method: request.method,
            headers: request.headers,
          };

          if (request.body) {
            options.body = JSON.stringify(request.body);
          }

          const response = await this.fetch(request.url, options);
          const data = await response.json();

          if (!response.ok) {
            return {
              ok: false,
              data: data,
              error: {
                message: data.error?.message ?? "Unknown error",
                details: data.error?.details,
              },
            };
          }

          return {
            ok: true,
            data,
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

    export interface UFOCientRequest {
      url: string;
      method: UFOHTTPMethod;
      body?: unknown;
      headers?: Record<string, string>;
    }

    export interface UFOClientResponse<T = unknown> {
      ok: boolean;
      data: T;
      error?: {
        message: string;
        details?: Record<string, unknown>;
      };
    }

    export interface UFOHTTPClient {
      request<T>(request: UFOCientRequest): Promise<UFOClientResponse<T>>;
    }

    export interface UFOClientMiddleware {
      before?(request: UFOCientRequest): Promise<UFOCientRequest>;
      after?(response: UFOClientResponse): Promise<UFOClientResponse>;
      error?(error: UFOError): Promise<never>;
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
        input: UFOProcedures[P]["input"]
      ): Promise<UFOProcedures[P]["output"]> {
        let request: UFOCientRequest = {
          url: this.buildUrl(procedure, method, input),
          method,
          headers: method === "POST"
            ? { "Content-Type": "application/json" }
            : undefined,
          ...(method === "POST" && { body: input }),
        };

        try {
          for (const m of this.middleware) {
            if (m.before) request = await m.before(request);
          }

          let response = await this.httpClient.request<UFOProcedures[P]["output"]>(
            request
          );

          for (const m of this.middleware) {
            if (m.after) {
              response = await m.after(response) as UFOClientResponse<
                UFOProcedures[P]["output"]
              >;
            }
          }

          if (!response.ok) {
            throw new UFOError(
              response.error?.message ?? "Unknown error",
              response.error?.details
            );
          }

          return response.data;
        } catch (err) {
          const error = err instanceof UFOError
            ? err
            : new UFOError(err instanceof Error ? err.message : "Unknown error");

          for (const m of this.middleware) {
            if (m.error) await m.error(error);
          }

          throw error;
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
    domainTypesTemplate,
    procedureTypesTemplate,
  ];

  if (opts.includeServer) {
    templates.push(serverTemplate);
  }

  if (opts.includeClient) {
    templates.push(createClientTemplate(opts));
  }

  const generated = templates.map((template) =>
    compileTemplate(template)(schema)
  ).join("\n");
  const formatted = await formatCode(generated);

  return formatted;
}
