<!-- 
  This component is responsible for initiating the rendering of a field based on its type.
  It handles four scenarios:

    1. Inline Single: A field defined inline and is not an array.
    2. Inline Array: A field defined inline and is an array.
    3. Named Single: A field that references a named type and is not an array.
    4. Named Array: A field that references a named type and is an array.

  If the named type is a custom type (not a primitive), it expands the underlying fields
  and treats it as an inline type for rendering purposes.
-->

<script lang="ts">
  import { primitiveTypes, store } from "$lib/store.svelte";
  import type { FieldDefinition } from "$lib/urpcTypes";

  import FieldInlineArray from "./FieldInlineArray.svelte";
  import FieldInlineSingle from "./FieldInlineSingle.svelte";
  import FieldNamedArray from "./FieldNamedArray.svelte";
  import FieldNamedSingle from "./FieldNamedSingle.svelte";

  interface Props {
    path: string;
    field: FieldDefinition;
    value: Record<string, any>;
  }

  let { field: originalField, value = $bindable(), path }: Props = $props();

  /**
   * Get fields of a custom type
   * @param typeName Name of the custom type
   */
  function getCustomTypeFields(typeName: string): FieldDefinition[] {
    for (const node of store.jsonSchema.nodes) {
      if (node.kind !== "type") continue;
      if (node.name !== typeName) continue;
      if (!node.fields) break;
      return node.fields;
    }

    return [];
  }

  /**
   * This is the field with expanded typeInline if it's a custom type
   */
  let field = $derived.by(() => {
    if (!originalField.typeName) return originalField;
    if (primitiveTypes.includes(originalField.typeName)) return originalField;

    return {
      ...originalField,
      typeName: undefined,
      typeInline: {
        fields: getCustomTypeFields(originalField.typeName),
      },
    } satisfies FieldDefinition;
  });

  let isInlineArray = $derived(field.typeInline && field.isArray);
  let isInlineSingle = $derived(field.typeInline && !field.isArray);
  let isNamedArray = $derived(field.typeName && field.isArray);
  let isNamedSingle = $derived(field.typeName && !field.isArray);
</script>

{#if isInlineArray}
  <FieldInlineArray {field} {path} bind:value />
{/if}

{#if isInlineSingle}
  <FieldInlineSingle {field} {path} bind:value />
{/if}

{#if isNamedArray}
  <FieldNamedArray {field} {path} bind:value />
{/if}

{#if isNamedSingle}
  <FieldNamedSingle {field} {path} bind:value />
{/if}
