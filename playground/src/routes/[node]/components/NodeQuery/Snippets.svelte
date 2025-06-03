<script lang="ts">
  import { ChevronLeft, ChevronRight, Code } from "@lucide/svelte";
  import { fade, slide } from "svelte/transition";

  import { getHeadersObject, store } from "$lib/store.svelte";
  import { uiStore } from "$lib/uiStore.svelte";

  import Tooltip from "$lib/components/Tooltip.svelte";

  import SnippetsCode from "./SnippetsCode.svelte";

  interface Props {
    // biome-ignore lint/suspicious/noExplicitAny: it's too dynamic to determine the type
    value: any;
    procName: string;
  }

  const { value, procName }: Props = $props();

  let maxHeight = $derived.by(() => {
    const heightMargin = 16;
    const windowHeight = globalThis.innerHeight;
    const headerHeight = uiStore.header.size.clientHeight;

    const maxHeight = windowHeight - headerHeight - 2 * heightMargin;
    return maxHeight;
  });

  let stickyTop = $derived.by(() => {
    const heightMargin = 16;
    const headerHeight = uiStore.header.size.clientHeight;
    return headerHeight + heightMargin;
  });

  let curl = $derived.by(() => {
    const payload = {
      proc: procName,
      input: value.root ?? {},
    };

    let payloadStr = JSON.stringify(payload, null, 2);
    payloadStr = payloadStr.replace(/'/g, "'\\''");

    let c = `curl -X POST ${store.endpoint} \\\n`;

    for (const header of getHeadersObject().entries()) {
      let rawHeader = `${header[0]}: ${header[1]}`;
      c += `-H ${JSON.stringify(rawHeader)} \\\n`;
    }

    c += `-d '${payloadStr}'`;

    return c;
  });

  let isOpen = $derived(uiStore.codeSnippetsOpen);
  function toggleIsOpen() {
    uiStore.codeSnippetsOpen = !uiStore.codeSnippetsOpen;
  }
</script>

<div
  class={[
    "sticky flex flex-col self-start",
    {
      "w-[40%]": isOpen,
    },
  ]}
  style={`max-height: ${maxHeight}px; top: ${stickyTop}px;`}
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
      onclick={toggleIsOpen}
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
        "rounded-box border-base-content/20 rounded-t-none border border-t-0",
        "overflow-x-auto overflow-y-auto",
      ]}
      in:fade={{ duration: 100 }}
      out:slide={{ duration: 100, axis: "x" }}
    >
      <SnippetsCode {curl} />
    </div>
  {/if}
</div>
