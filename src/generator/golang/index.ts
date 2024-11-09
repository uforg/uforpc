import path from "node:path";
import type {
  FieldSchemaType,
  MainSchemaType,
  TypeSchemaType,
} from "@/schema/index.ts";
import { handlebarsCompileTemplate } from "@/generator/common/handlebars-compile-template.ts";
import { extractContentAfterMarker } from "@/generator/common/extract-content-after-marker.ts";
import { formatGoCode } from "@/generator/golang/format-go-code.ts";
import { isArrayType, parseArrayType } from "@/schema/helpers.ts";
import Handlebars from "handlebars";

function registerHelpers() {
  // Convierte un tipo primitivo a su equivalente en Go
  Handlebars.registerHelper(
    "goPrimitiveType",
    function (typeName: string): string {
      const typeMap: Record<string, string> = {
        "string": "string",
        "int": "int",
        "float": "float64",
        "boolean": "bool",
      };
      return typeMap[typeName] || "any";
    },
  );

  // Convierte un tipo UFO RPC a su equivalente en Go
  Handlebars.registerHelper(
    "goType",
    function (field: FieldSchemaType): string {
      if (typeof field.type !== "string") {
        return "any";
      }

      // Manejar arrays
      if (isArrayType(field)) {
        const parsed = parseArrayType(field);
        if (typeof parsed.type.type !== "string") {
          return "[]any";
        }

        const baseType = Handlebars.helpers["goPrimitiveType"](
          parsed.type.type,
        );
        return `[]${baseType}`;
      }

      // Manejar tipos personalizados que empiezan con mayúscula
      if (/^[A-Z]/.test(field.type)) {
        return field.type;
      }

      return Handlebars.helpers["goPrimitiveType"](field.type);
    },
  );

  // Determina si un campo es requerido
  Handlebars.registerHelper(
    "isRequiredField",
    function (field: FieldSchemaType): boolean {
      if (!field.rules) {
        return false;
      }

      for (const rule of field.rules) {
        if ("rule" in rule && rule.rule === "required") {
          return true;
        }
      }
      return false;
    },
  );

  // Genera comentarios en formato Go
  Handlebars.registerHelper(
    "goComment",
    function (text: string | undefined): string {
      if (!text) {
        return "";
      }
      // Asegurar que el comentario empiece con mayúscula y termine en punto
      const formattedText = text.charAt(0).toUpperCase() + text.slice(1);
      const finalText = formattedText.endsWith(".")
        ? formattedText
        : formattedText + ".";
      return finalText.split("\n").map((line) => `// ${line}`).join("\n");
    },
  );

  // Renderiza campos de struct en Go
  Handlebars.registerHelper(
    "renderGoFields",
    function (fields: Record<string, FieldSchemaType>): string {
      if (!fields) {
        return "";
      }

      let result = "";
      for (const key of Object.keys(fields)) {
        const field = fields[key];
        const fieldName = key[0].toUpperCase() + key.slice(1);
        const isRequired = Handlebars.helpers["isRequiredField"](field);

        // Agregar documentación del campo si existe
        if (field.desc) {
          result += `  ${Handlebars.helpers["goComment"](field.desc)}\n`;
        }

        let fieldType: string;
        if (field.fields) {
          // Es un objeto anidado
          fieldType = "struct {\n";
          fieldType += Handlebars.helpers["renderGoFields"](field.fields);
          fieldType += "  }";
        } else if (isArrayType(field)) {
          const parsed = parseArrayType(field);
          if (typeof parsed.type.type === "string") {
            if (/^[A-Z]/.test(parsed.type.type)) {
              // Es un array de tipo personalizado
              const baseType = !isRequired
                ? `Null${parsed.type.type}`
                : parsed.type.type;
              fieldType = `[]${baseType}`;
            } else {
              // Es un array de tipo primitivo
              const baseType = Handlebars.helpers["goPrimitiveType"](
                parsed.type.type,
              );
              fieldType = `[]${baseType}`;
              if (!isRequired) {
                fieldType = `Null[${fieldType}]`;
              }
            }
          } else {
            fieldType = "[]any";
          }
        } else {
          fieldType = Handlebars.helpers["goType"](field);
          if (!isRequired) {
            if (field.type === "string") {
              fieldType = "NullString";
            } else if (field.type === "int") {
              fieldType = "NullInt";
            } else if (field.type === "float") {
              fieldType = "NullFloat64";
            } else if (field.type === "boolean") {
              fieldType = "NullBool";
            } else if (
              typeof field.type === "string" && /^[A-Z]/.test(field.type)
            ) {
              fieldType = `Null${field.type}`;
            } else {
              fieldType = `Null[${fieldType}]`;
            }
          }
        }

        result += `  ${fieldName} ${fieldType} \`json:"${key}${
          !isRequired ? ",omitempty" : ""
        }"\`\n`;
      }
      return result;
    },
  );

  // Generar tipo y su versión Null
  Handlebars.registerHelper(
    "renderTypeWithNull",
    function (type: TypeSchemaType): string {
      let result = "";

      if (type.desc) {
        result += `${Handlebars.helpers["goComment"](type.desc)}\n`;
      }

      // Tipo base
      result +=
        `// ${type.name} represents a domain type in the UFO RPC system.\n`;
      result += `type ${type.name} struct {\n`;
      result += Handlebars.helpers["renderGoFields"](type.fields);
      result += `}\n\n`;

      // Tipo Null correspondiente
      result +=
        `// Null${type.name} represents a nullable version of ${type.name}.\n`;
      result += `type Null${type.name} = Null[${type.name}]\n\n`;

      return result;
    },
  );

  // Inferir tipo Go para metadatos
  Handlebars.registerHelper("goMetaType", function (value: unknown): string {
    if (typeof value === "string") return "string";
    if (typeof value === "number") {
      if (Number.isInteger(value)) return "int";
      return "float64";
    }
    if (typeof value === "boolean") return "bool";
    return "any";
  });
}

function createPackageAndCoreTypesTemplate(opts: GenerateGolangOpts) {
  return `
    // This file has been generated using UFO RPC. DO NOT EDIT.
    // If you edit this file, it will be overwritten the next time it is generated

    // Package ${opts.packageName} contains the generated code for the UFO RPC schema
    package ${opts.packageName}

    import (
      "errors"
      "encoding/json"
    )

    // -----------------------------------------------------------------------------
    // Core Types
    // -----------------------------------------------------------------------------

    // UFOHTTPMethod represents an HTTP method.
    type UFOHTTPMethod string

    const (
      // GET represents the HTTP GET method.
      GET UFOHTTPMethod = "GET"
      // POST represents the HTTP POST method.
      POST UFOHTTPMethod = "POST"
    )

    // UFOResponse represents the response of a UFO RPC call.
    type UFOResponse[T any] struct {
      Ok     bool           \`json:"ok"\`
      Output T              \`json:"output,omitempty"\`
      Error  UFOErrorOutput \`json:"error,omitempty"\`
    }

    // UFOErrorOutput represents an error output in the UFO RPC system.
    type UFOErrorOutput struct {
      Message string                 \`json:"message"\`
      Details map[string]any \`json:"details,omitempty"\`
    }

    // UFOError represents an error in the UFO RPC system.
    type UFOError struct {
      Message string
      Details map[string]any
    }

    // Error implements the error interface.
    func (e *UFOError) Error() string {
      return e.Message
    }

    // getErrorOutput returns the UFOErrorOutput for a given error.
    func getErrorOutput(err error) UFOErrorOutput {
      var ufoErr *UFOError
      if errors.As(err, &ufoErr) {
        return UFOErrorOutput{
          Message: ufoErr.Message,
          Details: ufoErr.Details,
        }
      }
      return UFOErrorOutput{
        Message: err.Error(),
      }
    }
  `;
}

function createNullTypeTemplate() {
  const fileContent = Deno.readTextFileSync(path.join(
    import.meta.dirname || "./",
    "./null/null.go",
  ));

  return extractContentAfterMarker(fileContent);
}

function createDomainTypesTemplate() {
  return `
    // -----------------------------------------------------------------------------
    // Domain Types
    // -----------------------------------------------------------------------------

    {{#each types}}
    {{renderTypeWithNull this}}
    {{/each}}
  `;
}

function createProcedureTypesTemplate() {
  return `
    // -----------------------------------------------------------------------------
    // Procedure Types
    // -----------------------------------------------------------------------------

    {{#each procedures}}

    // P{{name}}Input represents the input parameters for the {{name}} procedure.
    {{#if input}}
    type P{{name}}Input struct {
      {{renderGoFields input}}
    }
    {{else}}
    type P{{name}}Input struct{}
    {{/if}}

    // P{{name}}Output represents the output results for the {{name}} procedure.
    {{#if output}}
    type P{{name}}Output struct {
      {{renderGoFields output}}
    }
    {{else}}
    type P{{name}}Output struct{}
    {{/if}}

    // P{{name}}Meta represents the metadata for the {{name}} procedure.
    {{#if meta}}
    type P{{name}}Meta struct {
      {{#each meta}}
      {{@key}} {{goMetaType this}} \`json:"{{@key}}"\`
      {{/each}}
    }
    {{else}}
    type P{{name}}Meta struct{}
    {{/if}}

    {{/each}}

    // ProcedureTypes defines the interface for all procedure types.
    type ProcedureTypes interface {
      {{#each procedures}}
        // {{name}} implements the {{name}} procedure.
        {{name}}(input P{{name}}Input) (UFOResponse[P{{name}}Output], error)
      {{/each}}
    }
  `;
}

export interface GenerateGolangOpts {
  packageName: string;
  includeServer?: boolean;
  includeClient?: boolean;
  omitServerRequestValidation?: boolean;
  omitClientRequestValidation?: boolean;
}

/**
 * Generates Golang code from a UFO RPC schema.
 *
 * @param schema - The UFO RPC schema to generate code from
 * @param opts - Options for code generation
 * @returns Generated Golang code as a string
 */
export async function generateGolang(
  schema: MainSchemaType,
  opts: GenerateGolangOpts,
): Promise<string> {
  registerHelpers();

  const templates = [
    createPackageAndCoreTypesTemplate(opts),
    createNullTypeTemplate(),
    createDomainTypesTemplate(),
    createProcedureTypesTemplate(),
  ];

  const compiled = templates.map(handlebarsCompileTemplate);
  const generated = compiled.map(
    (template) => template(schema).trim(),
  ).join("\n\n");

  return await formatGoCode(generated);
}
