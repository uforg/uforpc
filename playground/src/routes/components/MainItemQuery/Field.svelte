<script lang="ts">
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

  let fieldsArray = $derived(
    Array.isArray(fieldOrFields) ? fieldOrFields : [fieldOrFields],
  );
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
    <fieldset class="fieldset border border-base-300 rounded-box w-full p-4 space-y-2">
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
