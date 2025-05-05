<script lang="ts">
  import { BrushCleaning, Trash } from "@lucide/svelte";
  import type { FieldDefinition } from "$lib/urpcTypes";
  import { setAtPath } from "$lib/helpers/setAtPath";
  import { prettyLabel } from "./prettyLabel";
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

  let prettyPath = $derived(prettyLabel(path));
  let renderKey = $state(0);

  function clearObject() {
    renderKey++;
    value = setAtPath(value, path, {});
  }

  function deleteObject() {
    renderKey++;
    setTimeout(() => {
      value = setAtPath(value, path, null);
    }, 50);
  }
</script>

{#key renderKey}
  {#each fields as field}
    <Field
      fields={field}
      {path}
      bind:value
    />
  {/each}
{/key}

<div class="flex justify-end">
  <button
    class="btn btn-sm btn-ghost btn-square tooltip tooltip-left"
    data-tip={`Clear and reset ${prettyPath} to an empty object`}
    onclick={clearObject}
  >
    <BrushCleaning class="size-4" />
  </button>

  <button
    class="btn btn-sm btn-ghost btn-square tooltip tooltip-left"
    data-tip={`Delete ${prettyPath} from the JSON object`}
    onclick={deleteObject}
  >
    <Trash class="size-4" />
  </button>
</div>
