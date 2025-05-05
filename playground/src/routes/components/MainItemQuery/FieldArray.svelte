<script lang="ts">
  import {
    BrushCleaning,
    Minus,
    PackageOpen,
    Plus,
    Trash,
  } from "@lucide/svelte";
  import { setAtPath } from "$lib/helpers/setAtPath";
  import type { FieldDefinition } from "$lib/urpcTypes";
  import Tooltip from "$lib/components/Tooltip.svelte";
  import { prettyLabel } from "./prettyLabel";
  import Label from "./Label.svelte";
  import FieldNamed from "./FieldNamed.svelte";
  import FieldInline from "./FieldInline.svelte";
  import Fieldset from "./Fieldset.svelte";

  interface Props {
    field: FieldDefinition;
    path: string;
    value: any;
  }

  let {
    field,
    path,
    value = $bindable(),
  }: Props = $props();

  let indexes: number[] = $state([]);
  let lastIndex = $derived(indexes[indexes.length - 1]);
  let indexesLen = $derived(indexes.length);

  let prettyPath = $derived(prettyLabel(path));

  function clearArray() {
    value = setAtPath(value, path, []);
    indexes = [];
  }

  function deleteArray() {
    value = setAtPath(value, path, null);
    indexes = [];
  }

  function removeItem() {
    if (indexes.length <= 0) return;
    value = setAtPath(value, `${path}.${lastIndex}`, null);
    indexes.pop();
  }

  function addItem() {
    indexes.push(indexesLen);
  }
</script>

<Fieldset>
  <legend class="fieldset-legend">
    <Label optional={field.optional} label={path} />
  </legend>

  {#if indexesLen == 0}
    <PackageOpen class="size-6 mx-auto" />
    <p class="text-sm text-center italic">
      No items, add one using the button below
    </p>
  {/if}

  {#each indexes as index}
    {#if field.typeName}
      <FieldNamed
        {field}
        path={`${path}.${index}`}
        bind:value
      />
    {/if}

    {#if field.typeInline}
      <Fieldset>
        <legend class="fieldset-legend">
          <Label optional={field.optional} label={`${path}.${index}`} />
        </legend>
        <FieldInline
          fields={field.typeInline.fields}
          path={`${path}.${index}`}
          bind:value
        />
      </Fieldset>
    {/if}
  {/each}

  <div class="flex justify-end">
    <Tooltip
      content={`Clear and reset ${prettyPath} to an empty array`}
      placement="left"
    >
      <button
        class="btn btn-sm btn-ghost btn-square"
        onclick={clearArray}
      >
        <BrushCleaning class="size-4" />
      </button>
    </Tooltip>

    <Tooltip
      content={`Delete ${prettyPath} array from the JSON object`}
      placement="left"
    >
      <button
        class="btn btn-sm btn-ghost btn-square"
        onclick={deleteArray}
      >
        <Trash class="size-4" />
      </button>
    </Tooltip>

    {#if indexesLen > 0}
      <Tooltip
        content={`Remove last item from ${prettyPath} array`}
        placement="left"
      >
        <button
          class="btn btn-sm btn-ghost btn-square"
          onclick={removeItem}
        >
          <Minus class="size-4" />
        </button>
      </Tooltip>
    {/if}

    <Tooltip
      content={`Add item to ${prettyPath} array`}
      placement="left"
    >
      <button
        class="btn btn-sm btn-ghost btn-square"
        onclick={addItem}
      >
        <Plus class="size-4" />
      </button>
    </Tooltip>
  </div>
</Fieldset>
