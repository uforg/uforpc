import { zodToJsonSchema } from "zod-to-json-schema";
import path from "node:path";
import MainSchema from "../src/schema/schema.ts";

const jsonSchema = zodToJsonSchema(MainSchema, {
  name: "uforpcSchema",
});

const scriptsDir = import.meta.dirname ?? "";
const schemaPath = path.join(scriptsDir, "../src/schema/schema.json");

Deno.writeTextFileSync(
  schemaPath,
  JSON.stringify(jsonSchema, null, 2),
  { append: false },
);

console.log(`Wrote json schema to ${schemaPath}`);
