<script lang="ts">
  import type { FieldDefinition } from "$lib/urpcTypes";

  interface Props {
    field: FieldDefinition;
    value?: string | number | boolean | undefined;
  }

  let {
    field,
    value = $bindable(),
  }: Props = $props();

  export const validate = () => {
    isTouched = true;
    return isValid;
  };

  let isTouched = $state(false);

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
</script>

<label class="block space-y-1 w-full">
  <span
    class={{
      "block font-semibold": true,
      "text-error": isTouched && !isValid,
    }}
  >
    {field.name}
    {#if !field.optional}
      <span class="text-error">*</span>
    {/if}
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
      placeholder={`Enter "${field.name}" here...`}
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

  {#if field.optional}
    <p class="block label text-xs">Optional</p>
  {/if}

  {#if isTouched && !isValid}
    <p class="block label text-xs text-error">
      {field.name} is required
    </p>
  {/if}
</label>
