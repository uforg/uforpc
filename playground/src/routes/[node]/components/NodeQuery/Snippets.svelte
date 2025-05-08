<script lang="ts">
  import { ChevronLeft, ChevronRight, Code } from "@lucide/svelte";
  import { fade, slide } from "svelte/transition";

  import { store } from "$lib/store.svelte";
  import { uiStore, type UiStoreDimensions } from "$lib/uiStore.svelte";

  import Tooltip from "$lib/components/Tooltip.svelte";

  import SnippetsCode from "./SnippetsCode.svelte";

  interface Props {
    // biome-ignore lint/suspicious/noExplicitAny: it's too dynamic to determine the type
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

    const maxHeight = windowHeight - headerHeight - 2 * heightMargin;
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
  <Tooltip
    content={isOpen ? "Hide code snippets" : "Show code snippets"}
    placement="left"
  >
    <button
      class={[
        "btn rounded-box group border-base-content/20 flex w-full items-center",
        {
          "justify-start rounded-b-none px-4": isOpen,
          "px-3": !isOpen,
        },
      ]}
      onclick={() => (isOpen = !isOpen)}
    >
      <span
        class={{
          "mr-2": isOpen,
        }}
      >
        <Code class="size-4 group-hover:hidden" />
        {#if isOpen}
          <ChevronRight class="hidden size-4 group-hover:block" />
        {:else}
          <ChevronLeft class="hidden size-4 group-hover:block" />
        {/if}
      </span>

      {#if isOpen}
        <span>Code snippets for {procName}</span>
      {/if}
    </button>
  </Tooltip>

  {#if isOpen}
    <div
      class={[
        "rounded-box border-base-content/20 rounded-t-none border border-t-0 p-4",
        "overflow-x-auto overflow-y-auto",
      ]}
      in:fade={{ duration: 100 }}
      out:slide={{ duration: 100, axis: "x" }}
    >
      <SnippetsCode {curl} />
    </div>
  {/if}
</div>
