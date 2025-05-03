<script lang="ts">
  import { untrack } from "svelte";
  import type { FieldDefinition } from "$lib/urpcTypes";
  import { setAtPath } from "$lib/helpers/setAtPath";
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

  export const validate = () => {
    isTouched = true;
    return isValid;
  };

  let isTouched = $state(false);
  let value = $state();

  // Listen to changes and update the global value
  // Use untrack to avoid infinite loop
  // https://svelte.dev/docs/svelte/svelte#untrack
  $effect(() => {
    let val = value;
    untrack(() => {
      globalValue = setAtPath(globalValue, path, val);
    });
  });

  let isValid = $derived.by(() => {
    if (!field.typeName) return true;

    if (field.optional && !value) return true;

    if (field.typeName === "string") {
      return typeof value === "string" && value.length > 0;
    }

    if (["int", "float"].includes(field.typeName)) {
      return typeof value === "number";
    }

    if (field.typeName === "boolean") {
      return typeof value === "boolean";
    }

    if (field.typeName === "datetime") {
      return typeof value === "string";
    }

    return false;
  });

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

<label class="block space-y-1 w-full">
  <span
    class={{
      "block font-semibold": true,
      "text-error": isTouched && !isValid,
    }}
  >
    <Label optional={field.optional} {label} />
  </span>

  {#if inputType !== "checkbox"}
    <input
      type={inputType}
      step={inputStep}
      bind:value
      class={{
        "input w-full": true,
        "input-error placeholder:text-error": isTouched && !isValid,
      }}
      placeholder={`Enter ${label} here...`}
      onblur={() => (isTouched = true)}
    />
  {/if}

  {#if inputType === "checkbox"}
    <input
      type="checkbox"
      bind:checked={value as boolean}
      class={{
        "toggle toggle-lg": true,
        "toggle-error": isTouched && !isValid,
      }}
      onblur={() => (isTouched = true)}
    />
  {/if}

  {#if isTouched && !isValid}
    <p class="block label text-xs text-error">
      {label} is required
    </p>
  {/if}
</label>
