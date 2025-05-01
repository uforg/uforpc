<script lang="ts">
  import { onMount } from "svelte";
  import type { FieldDefinition } from "$lib/urpcTypes";
  import QueryProcFieldNamed from "./QueryProcFieldNamed.svelte";
  import QueryProcFieldInline from "./QueryProcFieldInline.svelte";

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
      if (field.depth == 0 && field.typeInline) {
        (value as any)[field.name] = {};
      }
    }
    mounted = true;
  });
</script>

{#if mounted}
  {#each fields as field}
    {#if field.depth == 0 && field.typeName}
      <QueryProcFieldNamed
        {field}
        bind:value={(value as any)[field.name]}
      />
    {/if}

    {#if field.depth == 0 && field.typeInline}
      <fieldset class="fieldset border border-base-300 rounded-box w-full p-4 space-y-2">
        <legend class="fieldset-legend">{field.name}</legend>
        <QueryProcFieldInline
          fields={field.typeInline.fields}
          bind:value={(value as any)[field.name]}
        />
      </fieldset>
    {/if}
  {/each}
{/if}
