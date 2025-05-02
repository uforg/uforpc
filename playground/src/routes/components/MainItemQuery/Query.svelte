<script lang="ts">
  import type { ProcedureDefinitionNode } from "$lib/urpcTypes";
  import Field from "./Field.svelte";
  import Snippets from "./Snippets.svelte";

  interface Props {
    proc: ProcedureDefinitionNode;
  }

  const { proc }: Props = $props();

  let value = $state({ root: {} });
</script>

{#if proc.input}
  <div class="flex space-x-4">
    <div class="space-y-2 flex-grow">
      <Field fields={proc.input} path="root" bind:value />
    </div>
    <Snippets {value} />
  </div>
{/if}
