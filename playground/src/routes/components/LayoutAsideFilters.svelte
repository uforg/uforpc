<script lang="ts">
  import {
    ArrowLeftRight,
    BookOpenText,
    CornerRightDown,
    Funnel,
    FunnelX,
    Search,
    Type,
    X,
  } from "@lucide/svelte";

  import { uiStore } from "$lib/uiStore.svelte";

  import Tooltip from "$lib/components/Tooltip.svelte";

  const searchTooltip = $derived(
    uiStore.asideSearchOpen ? "Close search" : "Open search",
  );
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

  let searchInput: HTMLInputElement | null = $state(null);

  function openSearch() {
    uiStore.asideSearchOpen = true;
    uiStore.asideSearchQuery = "";

    setTimeout(() => {
      searchInput?.focus();
    }, 100);
  }

  function closeSearch() {
    uiStore.asideSearchOpen = false;
    uiStore.asideSearchQuery = "";
  }

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
    uiStore.asideSearchOpen = false;
    uiStore.asideSearchQuery = "";
    uiStore.asideHideDocs = false;
    uiStore.asideHideTypes = true;
    uiStore.asideHideProcs = false;
    uiStore.asideHideStreams = false;
  }
</script>

<div class="flex w-full justify-between px-4 py-2">
  <Tooltip content="Reset filters to default" placement="bottom">
    <button class="btn btn-sm btn-square group" onclick={resetFilters}>
      <Funnel class="size-4 group-hover:hidden" />
      <FunnelX class="hidden size-4 group-hover:inline" />
    </button>
  </Tooltip>

  {#if uiStore.asideSearchOpen}
    <input
      type="text"
      class="input input-sm flex-grow"
      placeholder="Search..."
      bind:this={searchInput}
      bind:value={uiStore.asideSearchQuery}
    />

    <Tooltip content={searchTooltip} placement="bottom">
      <button class={["btn btn-sm btn-square relative"]} onclick={closeSearch}>
        <X class="size-4" />
      </button>
    </Tooltip>
  {/if}

  {#if !uiStore.asideSearchOpen}
    <Tooltip content={searchTooltip} placement="bottom">
      <button class={["btn btn-sm btn-square relative"]} onclick={openSearch}>
        <Search class="size-4" />
      </button>
    </Tooltip>
    <Tooltip content={docsTooltip} placement="bottom">
      <button
        class={[
          "btn btn-sm btn-square relative",
          uiStore.asideHideDocs && "toggle-disabled",
        ]}
        onclick={toggleDocs}
      >
        <BookOpenText class="size-4" />
      </button>
    </Tooltip>
    <Tooltip content={typesTooltip} placement="bottom">
      <button
        class={[
          "btn btn-sm btn-square relative",
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
          "btn btn-sm btn-square relative",
          uiStore.asideHideProcs && "toggle-disabled",
        ]}
        onclick={toggleProcs}
      >
        <ArrowLeftRight class="size-4" />
      </button>
    </Tooltip>
    <Tooltip content={streamsTooltip} placement="bottom">
      <button
        class={[
          "btn btn-sm btn-square relative",
          uiStore.asideHideStreams && "toggle-disabled",
        ]}
        onclick={toggleStreams}
      >
        <CornerRightDown class="size-4" />
      </button>
    </Tooltip>
  {/if}
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
