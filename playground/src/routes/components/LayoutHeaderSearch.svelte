<script lang="ts">
  import {
    ArrowLeftRight,
    BookOpenText,
    CornerRightDown,
    FileX,
    Search,
    Type,
    X,
  } from "@lucide/svelte";

  import { ctrlSymbol } from "$lib/helpers/ctrlSymbol";
  import {
    markSearchHintsMinisearch,
    truncateWithMarkMinisearch,
  } from "$lib/helpers/markSearchHints";
  import { miniSearch } from "$lib/store.svelte";
  import { uiStore } from "$lib/uiStore.svelte";

  import H2 from "$lib/components/H2.svelte";
  import Modal from "$lib/components/Modal.svelte";

  let input: HTMLInputElement | null = null;
  let isOpen = $state(false);
  const openModal = () => {
    isOpen = true;
    setTimeout(() => {
      input?.focus();
    }, 100);
  };
  const closeModal = () => (isOpen = false);

  const onKeydown = (e: KeyboardEvent) => {
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === "k") {
      e.preventDefault();
      openModal();
    }
  };

  $effect(() => {
    window.addEventListener("keydown", onKeydown);
    return () => {
      window.removeEventListener("keydown", onKeydown);
    };
  });

  let searchQuery = $state("");
  let searchResults = $derived(miniSearch.search(searchQuery));
</script>

<button
  class="btn btn-ghost flex items-center justify-start space-x-1 text-sm"
  onclick={openModal}
>
  <Search class="size-4" />
  <span>Search...</span>
  {#if !uiStore.isMobile}
    <span class="ml-4">
      <kbd class="kbd kbd-sm">{ctrlSymbol()}</kbd>
      <kbd class="kbd kbd-sm">K</kbd>
    </span>
  {/if}
</button>

<Modal bind:isOpen>
  <div class="flex items-center justify-start space-x-2">
    <label class="input flex-grow">
      <Search class="size-4" />
      <input
        bind:this={input}
        bind:value={searchQuery}
        type="search"
        placeholder="Search..."
      />
    </label>
    <button class="btn btn-square" onclick={closeModal}>
      <X class="size-4" />
    </button>
  </div>

  {#if searchResults.length === 0}
    <div
      class="my-8 flex flex-col items-center justify-center space-y-2 text-center"
    >
      <FileX class="size-12" />
      <H2>No results found</H2>
    </div>
  {/if}

  {#if searchResults.length > 0}
    <ul class="list mt-4">
      {#each searchResults as result}
        <a
          href={`#/${result.slug}`}
          onclick={closeModal}
          class="list-row hover:bg-base-200 block"
        >
          <span class="flex items-center justify-start text-lg font-bold">
            {#if result.kind === "doc"}
              <BookOpenText class="mr-2 size-4 flex-none" />
            {/if}
            {#if result.kind === "type"}
              <Type class="mr-2 size-4 flex-none" />
            {/if}
            {#if result.kind === "proc"}
              <ArrowLeftRight class="mr-2 size-4 flex-none" />
            {/if}
            {#if result.kind === "stream"}
              <CornerRightDown class="mr-2 size-4 flex-none" />
            {/if}
            {@html markSearchHintsMinisearch(result, result.name)}
          </span>
          <p class="truncate text-sm">
            {@html truncateWithMarkMinisearch(result, result.doc)}
          </p>
        </a>
      {/each}
    </ul>
  {/if}
</Modal>
