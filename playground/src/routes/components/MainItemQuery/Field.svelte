<script lang="ts">
  import type { FieldDefinition } from "$lib/urpcTypes";
  import { store } from "$lib/store.svelte";
  import FieldNamed from "./FieldNamed.svelte";
  import FieldInline from "./FieldInline.svelte";
  import FieldArray from "./FieldArray.svelte";
  import Label from "./Label.svelte";
  import Fieldset from "./Fieldset.svelte";

  interface Props {
    fields: FieldDefinition | FieldDefinition[];
    path: string;
    value: any;
  }

  let { fields: fieldOrFields, path, value = $bindable() }: Props = $props();

  function getCustomTypeFields(typeName: string): FieldDefinition[] {
    for (const node of store.jsonSchema.nodes) {
      if (node.kind !== "type") continue;
      if (node.name !== typeName) continue;
      if (!node.fields) break;
      return node.fields;
    }

    return [];
  }

  const primitiveTypes = ["string", "int", "float", "boolean", "datetime"];

  let fieldsArray = $derived.by(() => {
    let arr = Array.isArray(fieldOrFields) ? fieldOrFields : [fieldOrFields];

    // Transform custom fields to inline fields
    arr = arr.map((field) => {
      if (!field.typeName) return field;
      if (primitiveTypes.includes(field.typeName)) return field;

      const newField: FieldDefinition = {
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
    <FieldNamed {field} path={`${path}.${field.name}`} bind:value />
  {/if}

  {#if !field.isArray && field.typeInline}
    <Fieldset>
      <legend class="fieldset-legend">
        <Label label={`${path}.${field.name}`} optional={field.optional} />
      </legend>
      <FieldInline
        fields={field.typeInline.fields}
        path={`${path}.${field.name}`}
        bind:value
      />
    </Fieldset>
  {/if}

  {#if field.isArray}
    <FieldArray {field} path={`${path}.${field.name}`} bind:value />
  {/if}
{/each}
