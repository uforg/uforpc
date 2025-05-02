<script lang="ts">
  import type { ProcedureDefinitionNode } from "$lib/urpcTypes";
  import { mergeClasses } from "$lib/helpers/mergeClasses";
  import type { ClassValue } from "$lib/helpers/mergeClasses";
  import { setAtPath } from "$lib/helpers/setAtPath";
  import Field from "./Field.svelte";

  interface Props {
    proc: ProcedureDefinitionNode;
    class?: ClassValue;
  }

  const { proc, class: className }: Props = $props();

  let value = $state({ root: {} });
  // const setValue = (path: string, newValue: unknown) => {
  //   // value.root = setAtPath(value, path, newValue).root; // Si comento esto no pasa el bucle
  // };
</script>

{#if proc.input}
  <div class={mergeClasses("space-y-2", className)}>
    <Field fields={proc.input} parentPath="root" bind:value />

    <pre>{JSON.stringify(value, null, 2)}</pre>
  </div>
{/if}
