<script lang="ts">
  import { slide } from "svelte/transition";
  import { ChevronDown, ChevronRight, Info, Zap } from "@lucide/svelte";
  import type { ProcedureDefinitionNode } from "$lib/urpcTypes";
  import { dimensionschangeAction } from "$lib/uiStore.svelte";
  import type { UiStoreDimensions } from "$lib/uiStore.svelte";
  import H2 from "$lib/components/H2.svelte";
  import Field from "./Field.svelte";
  import Snippets from "./Snippets.svelte";

  interface Props {
    proc: ProcedureDefinitionNode;
  }

  const { proc }: Props = $props();

  let value = $state({ root: {} });
  let isOpen = $state(true);
  const toggleIsOpen = () => (isOpen = !isOpen);

  let dimensions: UiStoreDimensions | undefined = $state(undefined);
</script>

{#if proc.input}
  <div
    use:dimensionschangeAction
    ondimensionschange={(e) => (dimensions = e.detail)}
    class="flex space-x-4"
  >
    <div class="flex-grow">
      <button
        class={[
          "btn justify-start border-base-content/20 rounded-box w-full",
          "group/btn",
          {
            "rounded-b-none": isOpen,
          },
        ]}
        onclick={toggleIsOpen}
      >
        <Zap class="size-4 block group-hover/btn:hidden" />
        {#if isOpen}
          <ChevronDown class="size-4 hidden group-hover/btn:block" />
        {/if}
        {#if !isOpen}
          <ChevronRight class="size-4 hidden group-hover/btn:block" />
        {/if}
        <span class="ml-2">Try it out</span>
      </button>
      {#if isOpen}
        <div
          class="p-4 rounded-box rounded-t-none border border-t-0 border-base-content/20 space-y-2 fieldset"
          transition:slide={{ duration: 100 }}
        >
          <H2>Input</H2>

          <div role="alert" class="alert alert-soft alert-info w-fit">
            <Info class="size-4" />
            <span>
              All validations are performed on the server side
            </span>
          </div>

          <Field fields={proc.input} path="root" bind:value />
        </div>
      {/if}
    </div>

    <Snippets {value} procName={proc.name} parentDimensions={dimensions} />
  </div>
{/if}
