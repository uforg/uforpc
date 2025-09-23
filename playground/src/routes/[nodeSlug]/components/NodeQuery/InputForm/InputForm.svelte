<script lang="ts">
  import { BookText, Braces, Info } from "@lucide/svelte";

  import { storeUi } from "$lib/storeUi.svelte";
  import type { FieldDefinition } from "$lib/urpcTypes";

  import Tabs from "$lib/components/Tabs.svelte";

  import Field from "./Field.svelte";
  import JsonEditor from "./JsonEditor.svelte";

  interface Props {
    fields: FieldDefinition[];
    input: Record<string, any>;
  }

  let { fields, input = $bindable() }: Props = $props();

  let isFormTab = $derived(storeUi.store.inputFormTab === "form");
  let isJsonTab = $derived(storeUi.store.inputFormTab === "json");
</script>

<div
  class="desk:flex-row desk:items-start mt-4 flex flex-col items-end justify-between gap-4"
>
  <div role="alert" class="alert alert-soft alert-info w-fit">
    <Info class="size-4" />
    <span>
      Requests are made from your browser and validations are performed on the
      server side
    </span>
  </div>

  <Tabs
    containerClass="w-auto flex-none"
    buttonClass="btn-xs"
    iconClass="size-3"
    items={[
      { id: "form", label: "Form", icon: BookText },
      { id: "json", label: "JSON", icon: Braces },
    ]}
    bind:active={storeUi.store.inputFormTab}
  />
</div>

{#if isFormTab}
  {#each fields as field}
    <Field {field} path={field.name} bind:input />
  {/each}
{/if}

{#if isJsonTab}
  <JsonEditor bind:input />
{/if}
