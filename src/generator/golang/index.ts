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

  Handlebars.registerHelper(
    "renderValidationSchemaFields",
    function (fields: Record<string, FieldSchemaType>): string {
      if (!fields) return "";
      let result = "";

      function getBaseSchema(key: string, type: string): string {
        if (type === "int" || type === "float") {
          return `schValidator.Number("${key} must be a number")`;
        }
        if (type === "string") {
          return `schValidator.String("${key} must be a string")`;
        }
        if (type === "boolean") {
          return `schValidator.Boolean("${key} must be a boolean")`;
        }
        if (/^[A-Z]/.test(type)) {
          return `schValidator.Lazy(func() *schemaValidator { return vs${type}ValidationSchema }, "${key} must be a ${type}")`;
        }
        return "";
      }

      for (const [key, field] of Object.entries(fields)) {
        let schemaType = "";

        if (field.fields) {
          const nestedFields = Handlebars.helpers
            ["renderValidationSchemaFields"](field.fields);
          schemaType =
            `schValidator.Object(map[string]*schemaValidator{\n${nestedFields}}, "")`;
        } else if (isArrayType(field)) {
          const parsed = parseArrayType(field);
          schemaType = getBaseSchema(key, parsed.type.type);
          for (let i = 0; i < parsed.dimensions; i++) {
            schemaType = `schValidator.Array(${schemaType}, "")`;
          }
        } else if (typeof field.type === "string") {
          schemaType = getBaseSchema(key, field.type);
        }

        for (const rule of field.rules || []) {
          switch (rule.rule) {
            case "required": {
              const msg = rule.message || `${key} is required`;
              schemaType += `.Required("${msg}")`;
              break;
            }
            case "regex": {
              const msg = rule.message || `${key} must match ${rule.pattern}`;
              schemaType += `.Regex("${rule.pattern}", "${msg}")`;
              break;
            }
            case "contains": {
              const msg = rule.message || `${key} must contain ${rule.value}`;
              schemaType += `.Contains("${rule.value}", "${msg}")`;
              break;
            }
            case "equals": {
              const msg = rule.message || `${key} must equal to ${rule.value}`;
              schemaType += `.Equals(${JSON.stringify(rule.value)}, "${msg}")`;
              break;
            }
            case "enum": {
              const msg = rule.message ||
                `${key} must be one of ${rule.values.join(", ")}`;

              let goValues = "[]any{";
              for (const value of rule.values) {
                goValues += JSON.stringify(value) + ", ";
              }
              goValues += "}";

              schemaType += `.Enum(${goValues}, "${msg}")`;
              break;
            }
            case "email": {
              const msg = rule.message || `${key} must be a valid email`;
              schemaType += `.Email("${msg}")`;
              break;
            }
            case "iso8601": {
              const msg = rule.message ||
                `${key} must be a valid ISO 8601 date`;
              schemaType += `.Iso8601("${msg}")`;
              break;
            }
            case "uuid": {
              const msg = rule.message || `${key} must be a valid UUID`;
              schemaType += `.UUID("${msg}")`;
              break;
            }
            case "json": {
              const msg = rule.message || `${key} must be a valid JSON string`;
              schemaType += `.JSON("${msg}")`;
              break;
            }
            case "length": {
              const msg = rule.message ||
                `${key} must have length ${rule.value}`;
              schemaType += `.Length(${rule.value}, "${msg}")`;
              break;
            }
            case "minLength": {
              const msg = rule.message ||
                `${key} must have a minimum length of ${rule.value}`;
              schemaType += `.MinLength(${rule.value}, "${msg}")`;
              break;
            }
            case "maxLength": {
              const msg = rule.message ||
                `${key} must have a maximum length of ${rule.value}`;
              schemaType += `.MaxLength(${rule.value}, "${msg}")`;
              break;
            }
            case "lowercase": {
              const msg = rule.message || `${key} must be lowercase`;
              schemaType += `.Lowercase("${msg}")`;
              break;
            }
            case "uppercase": {
              const msg = rule.message || `${key} must be uppercase`;
              schemaType += `.Uppercase("${msg}")`;
              break;
            }
            case "min": {
              const msg = rule.message ||
                `${key} must be greater than ${rule.value}`;
              schemaType += `.Min(${rule.value}, "${msg}")`;
              break;
            }
            case "max": {
              const msg = rule.message ||
                `${key} must be less than ${rule.value}`;
              schemaType += `.Max(${rule.value}, "${msg}")`;
              break;
            }
          }
        }

        result += `    "${key}": ${schemaType},\n`;
      }

      return result;
    },
  );
}

function createPackageAndCoreTypesTemplate(opts: GenerateGolangOpts) {
  return `
    // This file has been generated using UFO RPC. DO NOT EDIT.
    // If you edit this file, it will be overwritten the next time it is generated

    // Package ${opts.packageName} contains the generated code for the UFO RPC schema
    package ${opts.packageName}

    import (
      "encoding/json"
      "fmt"
      "regexp"
      "strings"
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
      Ok     bool             \`json:"ok"\`
      Output T                \`json:"output,omitempty"\`
      Error  UFOError         \`json:"error,omitempty"\`
    }

    // UFOError represents a standardized error in the UFO RPC system.
    //
    // It provides structured information about errors that occur within the system,
    // enabling consistent error handling across servers and clients.
    //
    // Fields:
    //   - Message: A human-readable description of the error.
    //   - Category: Optional. Categorizes the error by its nature or source (e.g., "ValidationError", "DatabaseError").
    //   - Code: Optional. A machine-readable identifier for the specific error condition (e.g., "INVALID_EMAIL").
    //   - Details: Optional. Additional information about the error.
    //
    // The struct implements the error interface.
    type UFOError struct {
      // Message provides a human-readable description of the error.
      //
      // This message can be displayed to end-users or used for logging and debugging purposes.
      //
      // Use Cases:
      //   1. If localization is not implemented, Message can be directly shown to the user to inform them of the issue.
      //   2. Developers can use Message in logs to diagnose problems during development or in production.
      Message string \`json:"message"\`

      // Category categorizes the error by its nature or source.
      //
      // Examples:
      //   - "ValidationError" for input validation errors.
      //   - "DatabaseError" for errors originating from database operations.
      //   - "AuthenticationError" for authentication-related issues.
      //
      // Use Cases:
      //   1. In middleware, you can use Category to determine how to handle the error.
      //      For instance, you might log "InternalError" types and return a generic message to the client.
      //   2. Clients can inspect the Category to decide whether to prompt the user for action,
      //      such as re-authentication if the Category is "AuthenticationError".
      Category string \`json:"category,omitempty"\`

      // Code is a machine-readable identifier for the specific error condition.
      //
      // Examples:
      //   - "INVALID_EMAIL" when an email address fails validation.
      //   - "USER_NOT_FOUND" when a requested user does not exist.
      //   - "RATE_LIMIT_EXCEEDED" when a client has made too many requests.
      //
      // Use Cases:
      //   1. Clients can map Codes to localized error messages for internationalization (i18n),
      //      displaying appropriate messages based on the user's language settings.
      //   2. Clients or middleware can implement specific logic based on the Code,
      //      such as retry mechanisms for "TEMPORARY_FAILURE" or showing captcha for "RATE_LIMIT_EXCEEDED".
      Code string \`json:"code,omitempty"\`

      // Details contains optional additional information about the error.
      //
      // This field can include any relevant data that provides more context about the error.
      // The contents should be serializable to JSON.
      //
      // Use Cases:
      //   1. Providing field-level validation errors, e.g., Details could be:
      //      {"fields": {"email": "Email is invalid", "password": "Password is too short"}}
      //   2. Including diagnostic information such as timestamps, request IDs, or stack traces
      //      (ensure sensitive information is not exposed to clients).
      Details map[string]any \`json:"details,omitempty"\`
    }

    // Error implements the error interface, returning the error message.
    func (e UFOError) Error() string {
      return e.Message
    }

    // String implements the fmt.Stringer interface, returning the error message.
    func (e UFOError) String() string {
      return e.Message
    }
    
    // ToJSON returns the UFOError as a JSON-formatted string including all its fields.
    // This is useful for logging and debugging purposes.
    //
    // Example usage:
    //   err := UFOError{
    //     Category: "ValidationError",
    //     Code:     "INVALID_EMAIL",
    //     Message:  "The email address provided is invalid.",
    //     Details:  map[string]any{
    //       "field": "email",
    //     },
    //   }
    //   log.Println(err.ToJSON())
    func (e UFOError) ToJSON() string {
      b, err := json.Marshal(e)
      if err != nil {
        return fmt.Sprintf(
          \`{"message":%q,"error":"Failed to marshal UFOError: %s"}\`,
          e.Message, err.Error(),
        )
      }
      return string(b)
    }

    // asUFOError converts any error into a UFOError.
    // If the provided error is already a UFOError, it returns it as is.
    // Otherwise, it wraps the error message into a new UFOError.
    //
    // This function ensures that all errors conform to the UFOError structure,
    // facilitating consistent error handling across the system.
    func asUFOError(err error) UFOError {
      switch e := err.(type) {
      case UFOError:
        return e
      case *UFOError:
        return *e
      default:
        return UFOError{
          Message: err.Error(),
        }
      }
    }
  `;
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

    {{/each}}

    // ProcedureTypes defines the interface for all procedure types.
    type ProcedureTypes interface {
      {{#each procedures}}
        // {{name}} implements the {{name}} procedure.
        {{name}}(input P{{name}}Input) (UFOResponse[P{{name}}Output], error)
      {{/each}}
    }
    
    type UFOProcedureName string

    var UFOProcedureNames = struct {
      {{#each procedures}}
        {{name}} UFOProcedureName
      {{/each}}
    }{
      {{#each procedures}}
        {{name}}: "{{name}}",
      {{/each}}
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

function createValidationSchemaTemplate(opts: GenerateGolangOpts): string {
  if (opts.omitClientRequestValidation && opts.omitServerRequestValidation) {
    return "";
  }

  let validatorContent = Deno.readTextFileSync(path.join(
    import.meta.dirname || "./",
    "./validator/validator.go",
  ));

  validatorContent = extractContentAfterMarker(validatorContent);

  const template = `
    {{#each types}}
    
    // vs{{name}}ValidationSchema defines the validation rules for the {{name}} type
    var vs{{name}}ValidationSchema = schValidator.Object(map[string]*schemaValidator{
      {{renderValidationSchemaFields fields}}
    }, "")

    {{/each}}

    {{#each procedures}}

    {{#if input}}
    // vs{{name}}InputValidationSchema defines the validation rules for the {{name}} procedure input
    var vs{{name}}InputValidationSchema = schValidator.Object(map[string]*schemaValidator{
      {{renderValidationSchemaFields input}}
    }, "")
    {{/if}}

    {{/each}}

    // validationSchemas contains all available validation schemas for procedures
    var validationSchemas = map[string]struct {
      HasValidator bool
      Validator    *schemaValidator
    }{
      {{#each procedures}}
      "{{name}}": {
        {{#if input}}
        HasValidator: true,
        Validator:    vs{{name}}InputValidationSchema,
        {{else}}
        HasValidator: false,
        Validator:    nil,
        {{/if}}
      },
      {{/each}}
    }
  `;

  return `${validatorContent}\n\n${template}`;
}

function createServerTemplate(opts: GenerateGolangOpts): string {
  if (!opts.includeServer) return "";

  let validationLogic = "isValid := true";
  if (!opts.omitServerRequestValidation) {
    validationLogic = `
      isValid := true
      if valSchema, exists := validationSchemas[string(procedureName)]; exists && valSchema.HasValidator {
        valRes := valSchema.Validator.Validate(request.Input)
        isValid = valRes.IsValid
        if !isValid {
          response = UFOResponse[any]{
            Ok: false,
            Error: UFOError{
              Message: valRes.Error,
            },
          }
        }
      }
    `;
  }

  return `
    // -----------------------------------------------------------------------------
    // Server Types
    // -----------------------------------------------------------------------------

    // UFOServerRequest represents an incoming RPC request
    type UFOServerRequest[T any] struct {
      Method     UFOHTTPMethod
      Context    T
      Procedure  string
      Input      any
    }

    // UFOMiddlewareBefore represents a function that runs before request processing
    type UFOMiddlewareBefore[T any] func(context T) (T, error)

    // UFOMiddlewareAfter represents a function that runs after request processing
    type UFOMiddlewareAfter[T any] func(context T, response UFOResponse[any]) UFOResponse[any]

    // UFOServer handles RPC requests
    type UFOServer[T any] struct {
      handlers         map[UFOProcedureName]func(context T, input any) (any, error)
      beforeMiddleware []UFOMiddlewareBefore[T]
      afterMiddleware  []UFOMiddlewareAfter[T]
      methodMap       map[UFOProcedureName]UFOHTTPMethod
    }

    // NewUFOServer creates a new UFO RPC server
    func NewUFOServer[T any]() *UFOServer[T] {
      return &UFOServer[T]{
        handlers:         make(map[UFOProcedureName]func(T, any) (any, error)),
        beforeMiddleware: make([]UFOMiddlewareBefore[T], 0),
        afterMiddleware:  make([]UFOMiddlewareAfter[T], 0),
        methodMap: map[UFOProcedureName]UFOHTTPMethod{
          {{#each procedures}}
          UFOProcedureNames.{{name}}: "{{httpMethod type}}",
          {{/each}}
        },
      }
    }

    {{#each procedures}}
    // Set{{name}}Handler registers the handler for the {{name}} procedure
    func (s *UFOServer[T]) Set{{name}}Handler(
      handler func(context T, input P{{name}}Input) (P{{name}}Output, error),
    ) *UFOServer[T] {
      s.handlers[UFOProcedureNames.{{name}}] = func(context T, input any) (any, error) {
        typedInput, ok := input.(P{{name}}Input)
        if !ok {
          return nil, &UFOError{Message: "Invalid input type for {{name}}"}
        }
        return handler(context, typedInput)
      }
      return s
    }
    {{/each}}

    // AddMiddlewareBefore adds a middleware function that runs before the handler
    func (s *UFOServer[T]) AddMiddlewareBefore(fn UFOMiddlewareBefore[T]) *UFOServer[T] {
      s.beforeMiddleware = append(s.beforeMiddleware, fn)
      return s
    }

    // AddMiddlewareAfter adds a middleware function that runs after the handler
    func (s *UFOServer[T]) AddMiddlewareAfter(fn UFOMiddlewareAfter[T]) *UFOServer[T] {
      s.afterMiddleware = append(s.afterMiddleware, fn)
      return s
    }

    // HandleRequest processes an incoming RPC request
    func (s *UFOServer[T]) HandleRequest(request UFOServerRequest[T]) (UFOResponse[any], error) {
      procedureName := UFOProcedureName(request.Procedure)
      currentContext := request.Context
      response := UFOResponse[any]{Ok: true}
      shouldSkipHandler := false

      // Initial validation for procedure and method
      if _, exists := s.handlers[procedureName]; !exists {
        response = UFOResponse[any]{
          Ok: false,
          Error: UFOError{
            Message: fmt.Sprintf("Handler not defined for procedure %s", request.Procedure),
          },
        }
        shouldSkipHandler = true
      } else if expectedMethod := s.methodMap[procedureName]; expectedMethod != request.Method {
        response = UFOResponse[any]{
          Ok: false,
          Error: UFOError{
            Message: fmt.Sprintf("Method %s not allowed for %s procedure", request.Method, request.Procedure),
          },
        }
        shouldSkipHandler = true
      }

      // Execute Before middleware if we haven't encountered an error yet
      if !shouldSkipHandler {
        // Execute Before middleware
        for _, fn := range s.beforeMiddleware {
          var err error
          if currentContext, err = fn(currentContext); err != nil {
            response = UFOResponse[any]{
              Ok: false,
              Error: asUFOError(err),
            }
            shouldSkipHandler = true
            break
          }
        }
      }

      // Run handler if no errors have occurred
      if !shouldSkipHandler {
        // Validate input if required
        ${validationLogic}

        if isValid {
          // Execute handler
          if output, err := s.handlers[procedureName](currentContext, request.Input); err != nil {
            response = UFOResponse[any]{
              Ok:    false,
              Error: asUFOError(err),
            }
          } else {
            response = UFOResponse[any]{
              Ok:     true,
              Output: output,
            }
          }
        }
      }

      // Always execute After middleware, regardless of any previous errors
      for _, fn := range s.afterMiddleware {
        response = fn(currentContext, response)
      }

      return response, nil
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
    createDomainTypesTemplate(),
    createProcedureTypesTemplate(),
    createNullTypeTemplate(),
    createValidationSchemaTemplate(opts),
    createServerTemplate(opts),
  ];

  const compiled = templates.map(handlebarsCompileTemplate);
  const generated = compiled.map(
    (template) => template(schema).trim(),
  ).join("\n\n");

  return await formatGoCode(generated);
}
