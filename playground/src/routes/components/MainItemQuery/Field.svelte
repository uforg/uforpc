<script lang="ts">
  import { store } from "$lib/store.svelte";
  import type { FieldDefinitionWithLabel } from "./types";
  import FieldNamed from "./FieldNamed.svelte";
  import FieldInline from "./FieldInline.svelte";
  import FieldArray from "./FieldArray.svelte";

  interface Props {
    fields: FieldDefinitionWithLabel | FieldDefinitionWithLabel[];
    path: string;
    value: any;
  }

  let {
    fields: fieldOrFields,
    path,
    value = $bindable(),
  }: Props = $props();

  function getCustomTypeFields(
    typeName: string,
  ): FieldDefinitionWithLabel[] {
    for (const node of store.jsonSchema.nodes) {
      if (node.kind !== "type") continue;
      if (node.name !== typeName) continue;
      if (!node.fields) break;
      return node.fields;
    }

    return [];
  }

  const primitiveTypes = ["string", "int", "float", "bool", "datetime"];

  let fieldsArray = $derived.by(() => {
    let arr = Array.isArray(fieldOrFields)
      ? fieldOrFields
      : [fieldOrFields];

    // Transform custom fields to inline fields
    arr = arr.map((field) => {
      if (!field.typeName) return field;
      if (primitiveTypes.includes(field.typeName)) return field;

      const newField: FieldDefinitionWithLabel = {
        ...field,
        typeName: undefined,
        typeInline: {
          fields: getCustomTypeFields(field.typeName),
        },
      };

      return newField;
    });

    return arr;
  });
</script>

{#each fieldsArray as field}
  {#if !field.isArray && field.typeName}
    <FieldNamed
      {field}
      path={`${path}.${field.name}`}
      bind:value
    />
  {/if}

  {#if !field.isArray && field.typeInline}
    <fieldset class="fieldset border border-base-content/20 rounded-box w-full p-4 space-y-2">
      <legend class="fieldset-legend">{field.label ?? field.name}</legend>
      <FieldInline
        fields={field.typeInline.fields}
        path={`${path}.${field.name}`}
        bind:value
      />
    </fieldset>
  {/if}

  {#if field.isArray}
    <FieldArray
      {field}
      path={`${path}.${field.name}`}
      bind:value
    />
  {/if}
{/each}
