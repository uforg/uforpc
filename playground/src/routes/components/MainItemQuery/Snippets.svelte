<script lang="ts">
  import { fade, slide } from "svelte/transition";
  import H3 from "$lib/components/H3.svelte";
  import { ChevronLeft, ChevronRight, Code } from "@lucide/svelte";

  interface Props {
    value: any;
  }

  const { value }: Props = $props();

  let isOpen = $state(true);
</script>

<div
  class={[
    {
      "w-[400px] max-w-[400px]": isOpen,
      "flex-none": !isOpen,
    },
  ]}
>
  <button
    class={[
      "btn rounded-box w-full flex items-center group border-base-content/20",
      "tooltip tooltip-left",
      {
        "px-4 justify-start rounded-b-none": isOpen,
        "px-2 btn-square": !isOpen,
      },
    ]}
    data-tip={isOpen ? "Hide code snippets" : "Show code snippets"}
    onclick={() => (isOpen = !isOpen)}
  >
    <span
      class={{
        "mr-2": isOpen,
      }}
    >
      <Code class="size-6 group-hover:hidden" />

      {#if isOpen}
        <ChevronRight class="size-6 hidden group-hover:block" />
      {:else}
        <ChevronLeft class="size-6 hidden group-hover:block" />
      {/if}
    </span>

    {#if isOpen}
      <H3>Code snippets</H3>
    {/if}
  </button>

  {#if isOpen}
    <div
      class={[
        "p-4 rounded-box rounded-t-none border border-t-0 border-base-content/20",
        "overflow-x-auto",
      ]}
      in:fade={{ duration: 100 }}
      out:slide={{ duration: 100, axis: "x" }}
    >
      <pre>{JSON.stringify(value, null, 2)}</pre>
    </div>
  {/if}
</div>
