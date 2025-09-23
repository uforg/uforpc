<script lang="ts">
  import { MoveDownLeft, MoveUpRight, X } from "@lucide/svelte";

  import { formatISODate } from "$lib/helpers/formatISODate";

  import Code from "$lib/components/Code.svelte";
  import Modal from "$lib/components/Modal.svelte";
  import Tabs from "$lib/components/Tabs.svelte";

  import type { HistoryItem } from "../../../storeNode.svelte";

  interface Props {
    isOpen: boolean;
    historyItem: HistoryItem | null;
  }

  let { isOpen = $bindable(), historyItem }: Props = $props();

  let activeTab: "input" | "output" = $state("input");

  // Reset tab when modal opens
  $effect(() => {
    if (isOpen) activeTab = "input";
  });
</script>

<Modal bind:isOpen class="w-[95vw] max-w-4xl">
  {#if historyItem}
    <div class="flex h-full max-h-[80vh] flex-col">
      <div
        class="border-base-content/20 flex items-center justify-between border-b pb-4"
      >
        <h2 class="text-lg font-semibold">
          {formatISODate(historyItem.date)}
        </h2>
        <button
          class="btn btn-ghost btn-sm"
          onclick={() => (isOpen = false)}
          aria-label="Close"
        >
          <X class="size-4" />
        </button>
      </div>

      <div class="mt-4">
        <Tabs
          items={[
            { id: "input", label: "Input", icon: MoveUpRight },
            { id: "output", label: "Output", icon: MoveDownLeft },
          ]}
          bind:active={activeTab}
        />
      </div>

      <div class="mt-4 flex-1 overflow-hidden">
        {#if activeTab === "input"}
          <Code lang="json" code={historyItem.input} />
        {/if}

        {#if activeTab === "output"}
          <Code lang="json" code={historyItem.output} />
        {/if}
      </div>
    </div>
  {/if}
</Modal>
