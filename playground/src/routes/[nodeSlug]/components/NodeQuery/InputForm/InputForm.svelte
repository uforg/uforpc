<script lang="ts">
  import { Info } from "@lucide/svelte";

  import { uiStore } from "$lib/uiStore.svelte";
  import type { FieldDefinition } from "$lib/urpcTypes";

  import Field from "./Field.svelte";
  import JsonEditor from "./JsonEditor.svelte";

  interface Props {
    fields: FieldDefinition[];
    input: Record<string, any>;
  }

  let { fields, input = $bindable() }: Props = $props();

  let isFormTab = $derived(uiStore.store.inputFormTab === "form");
  const switchToForm = () => (uiStore.store.inputFormTab = "form");

  let isJsonTab = $derived(uiStore.store.inputFormTab === "json");
  const switchToJson = () => (uiStore.store.inputFormTab = "json");
</script>

<div
  class="desk:flex-row desk:items-center mt-4 flex flex-col items-end justify-between gap-4"
>
  <div role="alert" class="alert alert-soft alert-info w-fit">
    <Info class="size-4" />
    <span>
      Requests are made from your browser and validations are performed on the
      server side
    </span>
  </div>

  <div class="join">
    <button
      class={[
        "btn btn-xs join-item border-base-content/20 flex-grow",
        isFormTab && "btn-primary btn-active",
      ]}
      onclick={switchToForm}>Form</button
    >
    <button
      class={[
        "btn btn-xs join-item border-base-content/20 flex-grow",
        isJsonTab && "btn-primary btn-active",
      ]}
      onclick={switchToJson}>JSON</button
    >
  </div>
</div>

{#if isFormTab}
  {#each fields as field}
    <Field {field} path={field.name} bind:input />
  {/each}
{/if}

{#if isJsonTab}
  <JsonEditor bind:input />
{/if}
