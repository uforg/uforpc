<script lang="ts">
  import { BrushCleaning, Trash } from "@lucide/svelte";

  import { setAtPath } from "$lib/helpers/setAtPath";
  import type { FieldDefinition } from "$lib/urpcTypes";

  import Tooltip from "$lib/components/Tooltip.svelte";

  import Field from "./Field.svelte";
  import { prettyLabel } from "./prettyLabel";

  interface Props {
    fields: FieldDefinition[];
    path: string;
    // biome-ignore lint/suspicious/noExplicitAny: it's too dynamic to determine the type
    value: any;
  }

  let { fields, path, value = $bindable() }: Props = $props();

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
    <Field fields={field} {path} bind:value />
  {/each}
{/key}

<div class="flex justify-end">
  <Tooltip
    content={`Clear and reset ${prettyPath} to an empty object`}
    placement="left"
  >
    <button class="btn btn-sm btn-ghost btn-square" onclick={clearObject}>
      <BrushCleaning class="size-4" />
    </button>
  </Tooltip>

  <Tooltip
    content={`Delete ${prettyPath} from the JSON object`}
    placement="left"
  >
    <button class="btn btn-sm btn-ghost btn-square" onclick={deleteObject}>
      <Trash class="size-4" />
    </button>
  </Tooltip>
</div>
