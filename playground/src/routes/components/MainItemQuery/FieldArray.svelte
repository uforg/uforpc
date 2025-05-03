<script lang="ts">
  import { PackageOpen, Plus, Trash } from "@lucide/svelte";
  import { setAtPath } from "$lib/helpers/setAtPath";
  import type { FieldDefinition } from "$lib/urpcTypes";
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

  function addItem() {
    indexes.push(indexesLen);
  }

  function removeItem() {
    if (indexes.length <= 0) return;
    value = setAtPath(value, `${path}.${lastIndex}`, null);
    indexes.pop();
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
</Fieldset>
