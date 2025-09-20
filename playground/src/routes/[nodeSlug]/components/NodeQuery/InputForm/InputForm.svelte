<script lang="ts">
  import type { FieldDefinition } from "$lib/urpcTypes";

  import Field from "./Field.svelte";
  import JsonEditor from "./JsonEditor.svelte";

  interface Props {
    fields: FieldDefinition[];
    input: Record<string, any>;
  }

  let { fields, input = $bindable() }: Props = $props();

  let tab: "form" | "json" = $state("form");
</script>

<div class="flex justify-end">
  <div class="join">
    <button
      class={[
        "btn btn-xs join-item border-base-content/20 flex-grow",
        tab === "form" && "btn-primary btn-active",
      ]}
      onclick={() => (tab = "form")}>Form</button
    >
    <button
      class={[
        "btn btn-xs join-item border-base-content/20 flex-grow",
        tab === "json" && "btn-primary btn-active",
      ]}
      onclick={() => (tab = "json")}>JSON</button
    >
  </div>
</div>

{#if tab === "form"}
  {#each fields as field}
    <Field {field} path={field.name} bind:input />
  {/each}
{/if}

{#if tab === "json"}
  <JsonEditor bind:input />
{/if}
