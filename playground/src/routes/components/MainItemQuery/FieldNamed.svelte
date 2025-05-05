<script lang="ts">
  import { untrack } from "svelte";
  import { BrushCleaning, Trash } from "@lucide/svelte";
  import type { FieldDefinition } from "$lib/urpcTypes";
  import { setAtPath } from "$lib/helpers/setAtPath";
  import Tooltip from "$lib/components/Tooltip.svelte";
  import { prettyLabel } from "./prettyLabel";
  import Label from "./Label.svelte";

  interface Props {
    field: FieldDefinition;
    path: string;
    value: any;
  }

  let {
    field,
    path,
    value: globalValue = $bindable(),
  }: Props = $props();

  let value: any = $state(null);

  // Listen to changes and update the global value
  // Use untrack to avoid infinite loop
  // https://svelte.dev/docs/svelte/svelte#untrack
  $effect(() => {
    let val = value;
    untrack(() => {
      globalValue = setAtPath(globalValue, path, val);
    });
  });

  export const deleteValue = () => (value = null);
  export const clearValue = () => {
    if (field.typeName === "string") value = "";
    if (field.typeName === "int") value = 0;
    if (field.typeName === "float") value = 0;
    if (field.typeName === "boolean") value = false;
    if (field.typeName === "datetime") value = new Date();
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

    if (field.typeName === "boolean") {
      return "checkbox";
    }

    if (field.typeName === "datetime") {
      return "datetime-local";
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
</script>

<label class="block space-y-1 w-full group/field">
  <span class="block font-semibold">
    <Label optional={field.optional} {label} />
  </span>

  {#if inputType !== "checkbox"}
    <div class="flex justify-start items-center">
      <input
        type={inputType}
        step={inputStep}
        bind:value
        class="input flex-grow mr-1 group-hover/field:border-base-content/50"
        placeholder={`Enter ${label} here...`}
      />

      <Tooltip
        content={`Clear and reset ${label} to its default value`}
        placement="left"
      >
        <button
          class="w-8 btn btn-ghost btn-square"
          onclick={clearValue}
        >
          <BrushCleaning class="size-4" />
        </button>
      </Tooltip>

      <Tooltip
        content={`Delete ${label} from the JSON object`}
        placement="left"
      >
        <button
          class="w-8 btn btn-ghost btn-square"
          onclick={deleteValue}
        >
          <Trash class="size-4" />
        </button>
      </Tooltip>
    </div>
  {/if}

  {#if inputType === "checkbox"}
    <div class="flex justify-start items-center space-x-2">
      <input
        type="checkbox"
        bind:checked={value as boolean}
        class="toggle toggle-lg"
      />

      <Tooltip
        content={`Delete ${label} from the JSON object`}
        placement="right"
      >
        <button
          class="btn btn-ghost btn-square"
          onclick={deleteValue}
        >
          <Trash class="size-4" />
        </button>
      </Tooltip>
    </div>
  {/if}
</label>
