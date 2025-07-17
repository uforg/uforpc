<script lang="ts">
  import { BrushCleaning, EllipsisVertical, Trash } from "@lucide/svelte";
  import flatpickr from "flatpickr";
  import { onMount, untrack } from "svelte";

  import { setAtPath } from "$lib/helpers/setAtPath";
  import type { FieldDefinition } from "$lib/urpcTypes";

  import Menu from "$lib/components/Menu.svelte";
  import Tooltip from "$lib/components/Tooltip.svelte";

  import FieldDoc from "./FieldDoc.svelte";
  import Label from "./Label.svelte";
  import { prettyLabel } from "./prettyLabel";

  interface Props {
    field: FieldDefinition;
    path: string;
    // biome-ignore lint/suspicious/noExplicitAny: it's too dynamic to determine the type
    value: any;
  }

  let { field, path, value: globalValue = $bindable() }: Props = $props();
  const fieldId = $props.id();

  // biome-ignore lint/suspicious/noExplicitAny: it's too dynamic to determine the type
  let value: any = $state(null);

  // Listen to changes and update the global value
  // Use untrack to avoid infinite loop
  // https://svelte.dev/docs/svelte/svelte#untrack
  $effect(() => {
    const val = value;
    untrack(() => {
      globalValue = setAtPath(globalValue, path, val);
    });
  });

  export const deleteValue = () => {
    if (flatpickrInstance) flatpickrInstance.clear();
    value = null;
  };

  export const clearValue = () => {
    if (field.typeName === "string") value = "";
    if (field.typeName === "int") value = 0;
    if (field.typeName === "float") value = 0;
    if (field.typeName === "bool") value = false;
    if (field.typeName === "datetime") {
      let now = new Date();
      if (flatpickrInstance) flatpickrInstance.setDate(now);
      value = now;
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

  let label = $derived(prettyLabel(path));

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
      content={`Clear and reset ${label} to its default value`}
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

    <Tooltip content={`Delete ${label} from the JSON object`} placement="left">
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
      <Label optional={field.optional} {label} />
    </span>

    {#if inputType !== "checkbox" && inputType !== "datetime"}
      <div class="mb-1 flex items-center justify-start">
        <input
          type={inputType}
          step={inputStep}
          bind:value
          class="input group-hover/field:border-base-content/50 mr-1 flex-grow"
          placeholder={`Enter ${label} here...`}
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
          bind:value
          class="input group-hover/field:border-base-content/50 mr-1 flex-grow"
          placeholder={`Enter ${label} here...`}
        />

        {@render menu()}
      </div>
    {/if}

    {#if inputType === "checkbox"}
      <div class="flex items-center justify-start space-x-2">
        <input
          type="checkbox"
          bind:checked={value as boolean}
          class="toggle toggle-lg"
        />

        {@render menu()}
      </div>
    {/if}
  </label>

  <FieldDoc doc={field.doc} />
</div>
