import path from "node:path";
import { MainSchemaType } from "@/schema/index.ts";
import { handlebarsCompileTemplate } from "../common/handlebars-compile-template.ts";
import { extractContentAfterMarker } from "@/generator/common/extract-content-after-marker.ts";
import { formatGoCode } from "@/generator/golang/format-go-code.ts";

function registerHelpers() {
}

function createPackageAndCoreTypesTemplate(opts: GenerateGolangOpts) {
  return `
    // This file has been generated using UFO RPC. DO NOT EDIT.
    // If you edit this file, it will be overwritten the next time it is generated

    // Package ${opts.packageName} contains the generated code for the UFO RPC
    // schema
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
      GET UFOHTTPMethod = "GET"
      POST UFOHTTPMethod = "POST"
    )

    // UFOResponse represents the response of a UFO RPC call.
    type UFOResponse[T any] struct {
      Ok     bool           \`json:"ok"\`
      Output *T             \`json:"output,omitempty"\`
      Error  UFOErrorOutput \`json:"error,omitempty"\`
    }

    // UFOErrorOutput represents an error output in the UFO RPC system.
    type UFOErrorOutput struct {
      Message string                 \`json:"message"\`
      Details map[string]interface{} \`json:"details,omitempty"\`
    }

    // UFOError represents an error in the UFO RPC system.
    type UFOError struct {
      Message string
      Details map[string]interface{}
    }

    // Error implements the error interface.
    func (e *UFOError) Error() string {
      return e.Message
    }

    // getErrorOutput gets the UFOErrorOutput from an error.
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

export interface GenerateGolangOpts {
  packageName: string;
  includeServer?: boolean;
  includeClient?: boolean;
  omitServerRequestValidation?: boolean;
  omitClientRequestValidation?: boolean;
}

/**
 * Generates Golang code from a UFO RPC schema, including types, server, and
 * client implementations
 */
export async function generateGolang(
  schema: MainSchemaType,
  opts: GenerateGolangOpts,
): Promise<string> {
  registerHelpers();

  const templates = [
    createPackageAndCoreTypesTemplate(opts),
    createNullTypeTemplate(),
  ];

  const compiled = templates.map(handlebarsCompileTemplate);
  const generated = compiled.map(
    (template) => template(schema).trim(),
  ).join("\n\n");

  return await formatGoCode(generated);
}
