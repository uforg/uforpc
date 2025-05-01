<script lang="ts">
  import type { ProcedureDefinitionNode } from "$lib/urpcTypes";
  import { mergeClasses } from "$lib/helpers/mergeClasses";
  import type { ClassValue } from "$lib/helpers/mergeClasses";
  import QueryProcFieldInline from "./QueryProcFieldInline.svelte";

  interface Props {
    proc: ProcedureDefinitionNode;
    class?: ClassValue;
  }

  const { proc, class: className }: Props = $props();

  let value = $state({});
</script>

{#if proc.input}
  <div class={mergeClasses("space-y-2", className)}>
    <QueryProcFieldInline fields={proc.input} bind:value />

    <pre>{JSON.stringify(value, null, 2)}</pre>
  </div>
{/if}
