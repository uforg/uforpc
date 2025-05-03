<script lang="ts">
  import type { FieldDefinition } from "$lib/urpcTypes";
  import Field from "./Field.svelte";

  interface Props {
    fields: FieldDefinition[];
    path: string;
    value: any;
  }

  let {
    fields,
    path,
    value = $bindable(),
  }: Props = $props();
</script>

{#each fields as field}
  <Field
    fields={field}
    {path}
    bind:value
  />
{/each}
