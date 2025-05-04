<script lang="ts">
  import { fade, slide } from "svelte/transition";
  import { ChevronLeft, ChevronRight, Code } from "@lucide/svelte";
  import { uiStore, type UiStoreDimensions } from "$lib/uiStore.svelte";
  import { store } from "$lib/store.svelte";
  import SnippetsCode from "./SnippetsCode.svelte";

  interface Props {
    value: any;
    procName: string;
    parentDimensions: UiStoreDimensions | undefined;
  }

  const { value, procName, parentDimensions }: Props = $props();

  let height = $state(0);
  let isOpen = $state(true);

  let maxHeight = $derived.by(() => {
    if (!parentDimensions) return 0;

    const heightMargin = 16;
    const windowHeight = globalThis.innerHeight;
    const headerHeight = uiStore.header.size.clientHeight;

    const maxHeight = windowHeight - headerHeight - (2 * heightMargin);
    return maxHeight;
  });

  let marginTop = $derived.by(() => {
    if (!parentDimensions) return 0;

    const heightMargin = 16;
    const headerHeight = uiStore.header.size.clientHeight;
    const topThreshold = headerHeight + heightMargin;

    const parentVpTop = parentDimensions.viewportOffset.top;
    const parentHeight = parentDimensions.size.clientHeight;

    const maxMarginTop = Math.max(0, parentHeight - height);

    let marginTop = 0;
    if (parentVpTop <= topThreshold) {
      const desiredMarginTop = Math.max(0, topThreshold - parentVpTop);
      marginTop = Math.min(desiredMarginTop, maxMarginTop);
    }

    return marginTop;
  });

  let curl = $derived.by(() => {
    const payload = {
      proc: procName,
      input: value.root ?? {},
    };

    let payloadStr = JSON.stringify(payload, null, 2);
    payloadStr = payloadStr.replace(/'/g, "'\\''");

    let c = `curl -X POST ${store.endpoint} \\\n`;
    c += `-H "Content-Type: application/json" \\\n`;
    c += `-d '${payloadStr}'`;

    return c;
  });
</script>

<div
  class={[
    "flex flex-col self-start",
    {
      "w-[40%]": isOpen,
    },
  ]}
  style={`max-height: ${maxHeight}px; margin-top: ${marginTop}px;`}
  bind:clientHeight={height}
>
  <button
    class={[
      "btn rounded-box w-full flex items-center group border-base-content/20",
      "tooltip tooltip-left",
      {
        "px-4 justify-start rounded-b-none": isOpen,
        "px-3": !isOpen,
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
      <Code class="size-4 group-hover:hidden" />
      {#if isOpen}
        <ChevronRight class="size-4 hidden group-hover:block" />
      {:else}
        <ChevronLeft class="size-4 hidden group-hover:block" />
      {/if}
    </span>

    {#if isOpen}
      <span>Code snippets for {procName}</span>
    {/if}
  </button>

  {#if isOpen}
    <div
      class={[
        "p-4 rounded-box rounded-t-none border border-t-0 border-base-content/20",
        "overflow-x-auto overflow-y-auto",
      ]}
      in:fade={{ duration: 100 }}
      out:slide={{ duration: 100, axis: "x" }}
    >
      <SnippetsCode {curl} />
    </div>
  {/if}
</div>
