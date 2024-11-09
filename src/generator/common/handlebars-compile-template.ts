import Handlebars from "handlebars";

/**
 * Compiles a Handlebars template with proper configuration for TypeScript generation
 */
export function handlebarsCompileTemplate(
  template: string,
): HandlebarsTemplateDelegate {
  return Handlebars.compile(template, {
    noEscape: true,
    strict: true,
  });
}
