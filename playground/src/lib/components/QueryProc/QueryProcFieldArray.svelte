<script lang="ts">
  import { onMount } from "svelte";
  import type { FieldDefinition } from "$lib/urpcTypes";
  import QueryProcFieldNamed from "./QueryProcFieldNamed.svelte";
  import QueryProcFieldInline from "./QueryProcFieldInline.svelte";

  interface Props {
    field: FieldDefinition;
    value: any;
  }

  let {
    field,
    value = $bindable([]),
  }: Props = $props();

  let mounted = $state(false);
  onMount(() => {
    mounted = true;
  });
</script>

{#if mounted}
  <fieldset>
    <legend>{field.name}</legend>
  </fieldset>
{/if}
