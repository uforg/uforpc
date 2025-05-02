<script lang="ts">
  import type { FieldDefinition } from "$lib/urpcTypes";
  import FieldNamed from "./FieldNamed.svelte";
  import FieldInline from "./FieldInline.svelte";

  interface Props {
    fields: FieldDefinition | FieldDefinition[];
    parentPath: string;
    value: any;
  }

  let {
    fields: fieldOrFields,
    parentPath,
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
      {parentPath}
      bind:value
    />
  {/if}

  {#if !field.isArray && field.typeInline}
    <fieldset class="fieldset border border-base-300 rounded-box w-full p-4 space-y-2">
      <legend class="fieldset-legend">{field.name}</legend>
      <FieldInline
        fields={field.typeInline.fields}
        parentPath={`${parentPath}.${field.name}`}
        bind:value
      />
    </fieldset>
  {/if}
{/each}
