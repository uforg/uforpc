<script lang="ts">
  import type { FieldDefinition } from "$lib/urpcTypes";

  import Field from "./Field.svelte";

  interface Props {
    input: FieldDefinition[];
    value: Record<string, any>;
  }

  let { input, value = $bindable() }: Props = $props();
</script>

{#each input as field}
  <Field {field} path={field.name} bind:value />
{/each}
