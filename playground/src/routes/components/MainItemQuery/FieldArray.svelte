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
    <button
      class="btn btn-sm btn-ghost btn-square tooltip tooltip-left"
      data-tip={`Clear and reset ${prettyPath} to an empty array`}
      onclick={clearArray}
    >
      <BrushCleaning class="size-4" />
    </button>

    <button
      class="btn btn-sm btn-ghost btn-square tooltip tooltip-left"
      data-tip={`Delete ${prettyPath} array from the JSON object`}
      onclick={deleteArray}
    >
      <Trash class="size-4" />
    </button>

    {#if indexesLen > 0}
      <button
        class="btn btn-sm btn-ghost btn-square tooltip tooltip-left"
        data-tip={`Remove last item from ${prettyPath} array`}
        onclick={removeItem}
      >
        <Minus class="size-4" />
      </button>
    {/if}

    <button
      class="btn btn-sm btn-ghost btn-square tooltip tooltip-left"
      data-tip={`Add item to ${prettyPath} array`}
      onclick={addItem}
    >
      <Plus class="size-4" />
    </button>
  </div>
</Fieldset>
