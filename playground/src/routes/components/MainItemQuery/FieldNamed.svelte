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
  <span class="block font-semibold">
    <Label optional={field.optional} {label} />
  </span>

  {#if inputType !== "checkbox"}
    <input
      type={inputType}
      step={inputStep}
      bind:value
      class="input w-full"
      placeholder={`Enter ${label} here...`}
    />
  {/if}

  {#if inputType === "checkbox"}
    <input
      type="checkbox"
      bind:checked={value as boolean}
      class="toggle toggle-lg"
    />
  {/if}
</label>
