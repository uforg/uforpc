<script lang="ts">
  import type { FieldDefinitionWithLabel } from "./types";
  import Field from "./Field.svelte";

  interface Props {
    fields: FieldDefinitionWithLabel[];
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
    path={`${path}.${field.name}`}
    bind:value
  />
{/each}
