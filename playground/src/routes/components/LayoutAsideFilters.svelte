<script lang="ts">
  import {
    ArrowLeftRight,
    BookOpenText,
    CornerRightDown,
    Funnel,
    FunnelX,
    Type,
  } from "@lucide/svelte";

  import { uiStore } from "$lib/uiStore.svelte";

  import Tooltip from "$lib/components/Tooltip.svelte";

  const docsTooltip = $derived(
    uiStore.asideHideDocs ? "Show documentation" : "Hide documentation",
  );
  const typesTooltip = $derived(
    uiStore.asideHideTypes ? "Show data types" : "Hide data types",
  );
  const procsTooltip = $derived(
    uiStore.asideHideProcs ? "Show procedures" : "Hide procedures",
  );
  const streamsTooltip = $derived(
    uiStore.asideHideStreams ? "Show streams" : "Hide streams",
  );

  function toggleDocs() {
    uiStore.asideHideDocs = !uiStore.asideHideDocs;
  }

  function toggleTypes() {
    uiStore.asideHideTypes = !uiStore.asideHideTypes;
  }

  function toggleProcs() {
    uiStore.asideHideProcs = !uiStore.asideHideProcs;
  }

  function toggleStreams() {
    uiStore.asideHideStreams = !uiStore.asideHideStreams;
  }

  function resetFilters() {
    uiStore.asideHideDocs = false;
    uiStore.asideHideTypes = false;
    uiStore.asideHideProcs = false;
    uiStore.asideHideStreams = false;
  }
</script>

<div class="w-full px-4 py-2">
  <div class="join flex">
    <Tooltip content="Reset filters" placement="bottom">
      <button
        class="btn btn-sm join-item group rounded-l-field flex-grow"
        onclick={resetFilters}
      >
        <Funnel class="size-4 group-hover:hidden" />
        <FunnelX class="hidden size-4 group-hover:inline" />
        <span>Filters</span>
      </button>
    </Tooltip>
    <Tooltip content={docsTooltip} placement="bottom">
      <button
        class={[
          "btn btn-sm join-item relative flex-none",
          uiStore.asideHideDocs && "toggle-disabled",
        ]}
        onclick={toggleDocs}
      >
        <BookOpenText class="size-4" />
      </button>
    </Tooltip>
    <Tooltip content={streamsTooltip} placement="bottom">
      <button
        class={[
          "btn btn-sm join-item relative flex-none",
          uiStore.asideHideStreams && "toggle-disabled",
        ]}
        onclick={toggleStreams}
      >
        <CornerRightDown class="size-4" />
      </button>
    </Tooltip>
    <Tooltip content={typesTooltip} placement="bottom">
      <button
        class={[
          "btn btn-sm join-item relative flex-none",
          uiStore.asideHideTypes && "toggle-disabled",
        ]}
        onclick={toggleTypes}
      >
        <Type class="size-4" />
      </button>
    </Tooltip>
    <Tooltip content={procsTooltip} placement="bottom">
      <button
        class={[
          "btn btn-sm join-item rounded-r-field relative flex-none",
          uiStore.asideHideProcs && "toggle-disabled",
        ]}
        onclick={toggleProcs}
      >
        <ArrowLeftRight class="size-4" />
      </button>
    </Tooltip>
  </div>
</div>

<style lang="postcss">
  .toggle-disabled::before {
    content: "";
    position: absolute;
    top: 50%;
    left: 10%;
    width: 80%;
    height: 2px;
    background-color: currentColor;
    transform: translateY(-50%) rotate(-45deg);
    opacity: 0.7;
    z-index: 1;
    pointer-events: none;
  }

  .toggle-disabled {
    opacity: 0.6;
  }
</style>
