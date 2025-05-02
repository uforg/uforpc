<script lang="ts">
  import type { FieldDefinition } from "$lib/urpcTypes";
  import Field from "./Field.svelte";

  interface Props {
    fields: FieldDefinition[];
    parentPath: string;
    value: any;
  }

  let {
    fields,
    parentPath,
    value = $bindable(),
  }: Props = $props();
</script>

{#each fields as field}
  <Field
    fields={field}
    {parentPath}
    bind:value
  />
{/each}
