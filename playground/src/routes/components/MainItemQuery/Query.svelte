<script lang="ts">
  import { fade } from "svelte/transition";
  import { ChevronDown, ChevronRight, ChevronUp } from "@lucide/svelte";
  import type { ProcedureDefinitionNode } from "$lib/urpcTypes";
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
</script>

{#if proc.input}
  <div class="flex space-x-4">
    <div class="flex-grow">
      <button
        class={[
          "btn justify-start border-base-content/20 rounded-box w-full",
          {
            "rounded-b-none": isOpen,
          },
        ]}
        onclick={toggleIsOpen}
      >
        {#if isOpen}
          <ChevronDown class="size-5" />
        {/if}
        {#if !isOpen}
          <ChevronRight class="size-5" />
        {/if}
        <span class="ml-2">Try it out</span>
      </button>
      {#if isOpen}
        <div
          class="p-4 rounded-box rounded-t-none border border-t-0 border-base-content/20 space-y-2"
          transition:fade={{ duration: 100 }}
        >
          <H2>Input</H2>
          <Field fields={proc.input} path="root" bind:value />
        </div>
      {/if}
    </div>

    <Snippets {value} />
  </div>
{/if}
