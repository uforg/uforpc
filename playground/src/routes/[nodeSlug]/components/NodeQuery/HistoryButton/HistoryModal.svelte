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

<Modal bind:isOpen class="h-full w-full max-w-4xl">
  {#if historyItem}
    <div class="flex h-full flex-col space-y-4 overflow-hidden">
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-semibold">
          {formatISODate(historyItem.date)}
        </h2>
        <button
          class="btn btn-ghost btn-sm btn-circle"
          onclick={() => (isOpen = false)}
          aria-label="Close"
        >
          <X class="size-4" />
        </button>
      </div>

      <div>
        <Tabs
          items={[
            { id: "input", label: "Input", icon: MoveUpRight },
            { id: "output", label: "Output", icon: MoveDownLeft },
          ]}
          bind:active={activeTab}
        />
      </div>

      <div class="flex-1 flex-grow overflow-hidden">
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
