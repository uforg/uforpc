import type { FieldDefinition } from "../../../lib/urpcTypes.ts";

export interface FieldDefinitionWithLabel extends FieldDefinition {
  label?: string;
}
