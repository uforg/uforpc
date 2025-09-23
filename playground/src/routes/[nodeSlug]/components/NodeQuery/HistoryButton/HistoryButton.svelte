<script lang="ts">
  import { History } from "@lucide/svelte";

  import Tooltip from "$lib/components/Tooltip.svelte";

  import type { StoreNodeInstance } from "../../../storeNode.svelte";

  import HistoryOffcanvas from "./HistoryOffcanvas.svelte";

  interface Props {
    nodeName: string;
    storeNode: StoreNodeInstance;
  }

  let { nodeName, storeNode }: Props = $props();

  let isOffcanvasOpen = $state(false);

  function openHistory() {
    isOffcanvasOpen = true;
  }
</script>

<Tooltip content={`View request history for ${nodeName}`}>
  <button class="btn btn-xs border-base-content/20" onclick={openHistory}>
    <History class="size-3" />
    <span>History</span>
  </button>
</Tooltip>

<HistoryOffcanvas bind:isOpen={isOffcanvasOpen} {storeNode} />
