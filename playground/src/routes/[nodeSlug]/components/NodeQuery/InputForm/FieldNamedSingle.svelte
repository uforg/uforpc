<!--
  This component handles the final single field rendering for a named type that is not an array.

  It handles primitive types only: string, int, float, bool, datetime

  It should handle reactivity and binding of default values correctly.
-->

<script lang="ts">
  import { BrushCleaning, EllipsisVertical, Trash } from "@lucide/svelte";
  import flatpickr from "flatpickr";
  import { get, set, unset } from "lodash-es";
  import { onMount } from "svelte";

  import type { FieldDefinition } from "$lib/urpcTypes";

  import Menu from "$lib/components/Menu.svelte";
  import Tooltip from "$lib/components/Tooltip.svelte";

  import CommonFieldDoc from "./CommonFieldDoc.svelte";
  import CommonLabel from "./CommonLabel.svelte";

  interface Props {
    path: string;
    field: FieldDefinition;
    value: Record<string, any>;
  }

  let { field, value = $bindable(), path }: Props = $props();
  const fieldId = $props.id();

  // We can't bind directly to value[path] because Svelte doesn't support dynamic bindings.
  // So we create a local reactive variable and update value[path] whenever it changes.
  let localValue = $state(get(value, path) ?? undefined);
  $effect(() => {
    if (get(value, path) !== localValue) {
      set(value, path, localValue);
    }
  });

  export const deleteValue = () => {
    if (flatpickrInstance) flatpickrInstance.clear();
    localValue = undefined;
    unset(value, path);
  };

  export const clearValue = () => {
    if (field.typeName === "string") localValue = "";
    if (field.typeName === "int") localValue = 0;
    if (field.typeName === "float") localValue = 0;
    if (field.typeName === "bool") localValue = false;
    if (field.typeName === "datetime") {
      let now = new Date();
      if (flatpickrInstance) flatpickrInstance.setDate(now);
      localValue = now;
    }
  };

  let inputType = $derived.by(() => {
    if (!field.typeName) {
      return "text";
    }

    if (field.typeName === "string") {
      return "text";
    }

    if (["int", "float"].includes(field.typeName)) {
      return "number";
    }

    if (field.typeName === "bool") {
      return "checkbox";
    }

    if (field.typeName === "datetime") {
      return "datetime";
    }
  });

  let inputStep = $derived.by(() => {
    if (field.typeName === "float") {
      return 0.01;
    }

    if (field.typeName === "int") {
      return 1;
    }
  });

  let flatpickrInstance: flatpickr.Instance | null = $state(null);
  onMount(() => {
    if (field.typeName !== "datetime") return;
    let inst = flatpickr(`#${fieldId}`, {
      enableTime: true,
      enableSeconds: true,
      dateFormat: "Z",
      altInput: true,
      altFormat: "F j, Y H:i:S",
    });
    if (Array.isArray(inst)) inst = inst[0];
    flatpickrInstance = inst;
  });
</script>

{#snippet menuContent()}
  <div class="py-1">
    <Tooltip
      content={`Clear and reset ${path} to its default value`}
      placement="left"
    >
      <button
        class="btn btn-ghost btn-block flex items-center justify-start space-x-2"
        onclick={clearValue}
      >
        <BrushCleaning class="size-4" />
        <span>Clear</span>
      </button>
    </Tooltip>

    <Tooltip content={`Delete ${path} from the JSON object`} placement="left">
      <button
        class="btn btn-ghost btn-block flex items-center justify-start space-x-2"
        onclick={deleteValue}
      >
        <Trash class="size-4" />
        <span>Delete</span>
      </button>
    </Tooltip>
  </div>
{/snippet}

{#snippet menu()}
  <Menu content={menuContent} placement="bottom" trigger="mouseenter click">
    <button class="btn btn-ghost btn-square">
      <EllipsisVertical class="size-4" />
    </button>
  </Menu>
{/snippet}

<div>
  <label class="group/field block w-full">
    <span class="mb-1 block font-semibold">
      <CommonLabel optional={field.optional} label={path} />
    </span>

    {#if inputType !== "checkbox" && inputType !== "datetime"}
      <div class="mb-1 flex items-center justify-start">
        <input
          type={inputType}
          step={inputStep}
          bind:value={localValue}
          class="input group-hover/field:border-base-content/50 mr-1 flex-grow"
          placeholder={`Enter ${path} here...`}
        />

        {@render menu()}
      </div>
    {/if}

    {#if inputType === "datetime"}
      <div class="mb-1 flex items-center justify-start">
        <input
          id={fieldId}
          type={inputType}
          step={inputStep}
          bind:value={localValue}
          class="input group-hover/field:border-base-content/50 mr-1 flex-grow"
          placeholder={`Enter ${path} here...`}
        />

        {@render menu()}
      </div>
      <div class="prose prose-sm text-base-content/50 max-w-none font-bold">
        Time is shown in your local timezone and will be sent as UTC
      </div>
    {/if}

    {#if inputType === "checkbox"}
      <div class="flex items-center justify-start space-x-2">
        <input
          type="checkbox"
          bind:checked={localValue}
          class="toggle toggle-lg"
        />

        {@render menu()}
      </div>
    {/if}
  </label>

  <CommonFieldDoc doc={field.doc} />
</div>
