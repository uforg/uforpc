<script lang="ts">
  import { setAtPath } from "$lib/helpers/setAtPath";
  import type { FieldDefinitionWithLabel } from "./types";
  import { PackageOpen, Plus, Trash } from "@lucide/svelte";
  import FieldNamed from "./FieldNamed.svelte";
  import FieldInline from "./FieldInline.svelte";

  interface Props {
    field: FieldDefinitionWithLabel;
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

  function addItem() {
    indexes.push(indexesLen);
  }

  function removeItem() {
    if (indexes.length <= 0) return;

    value = setAtPath(value, `${path}.${lastIndex}`, null, {
      removeNullOrUndefined: true,
    });

    indexes.pop();
  }
</script>

<fieldset class="fieldset border border-base-300 rounded-box w-full p-4 space-y-2">
  <legend class="fieldset-legend">{field.name}</legend>

  {#if indexesLen == 0}
    <PackageOpen class="size-6 mx-auto" />
    <p class="text-sm text-center italic">
      No items, add one using the button below
    </p>
  {/if}

  {#each indexes as index}
    {#if field.typeName}
      <FieldNamed
        field={{ ...field, label: `${field.name}[${index}]` }}
        path={`${path}.${index}`}
        bind:value
      />
    {/if}

    {#if field.typeInline}
      <fieldset class="fieldset border border-base-300 rounded-box w-full p-4 space-y-2">
        <legend class="fieldset-legend">
          {`${field.name}[${index}]`}
        </legend>
        <FieldInline
          fields={field.typeInline.fields}
          path={`${path}.${index}`}
          bind:value
        />
      </fieldset>
    {/if}
  {/each}

  <div class="flex justify-end space-x-2">
    {#if indexesLen > 0}
      <button class="btn btn-sm btn-soft btn-error" onclick={removeItem}>
        <Trash class="size-4" />
        Remove last item
      </button>
    {/if}

    <button class="btn btn-sm" onclick={addItem}>
      <Plus class="size-4" />
      Add item to {field.name}
    </button>
  </div>
</fieldset>
