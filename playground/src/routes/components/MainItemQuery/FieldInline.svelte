<script lang="ts">
  import { onMount } from "svelte";
  import type { FieldDefinition } from "$lib/urpcTypes";
  import FieldNamed from "./FieldNamed.svelte";
  import FieldInline from "./FieldInline.svelte";

  interface Props {
    fields: FieldDefinition[];
    value: unknown;
  }

  let {
    fields,
    value = $bindable(),
  }: Props = $props();

  let mounted = $state(false);
  onMount(() => {
    for (const field of fields) {
      if (!field.isArray && field.typeInline) {
        (value as any)[field.name] = {};
      }
    }
    mounted = true;
  });
</script>

{#if mounted}
  {#each fields as field}
    {#if !field.isArray && field.typeName}
      <FieldNamed
        {field}
        bind:value={(value as any)[field.name]}
      />
    {/if}

    {#if !field.isArray && field.typeInline}
      <fieldset class="fieldset border border-base-300 rounded-box w-full p-4 space-y-2">
        <legend class="fieldset-legend">{field.name}</legend>
        <FieldInline
          fields={field.typeInline.fields}
          bind:value={(value as any)[field.name]}
        />
      </fieldset>
    {/if}
  {/each}
{/if}
