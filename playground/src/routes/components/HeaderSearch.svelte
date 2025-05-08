<script lang="ts">
  import {
    ArrowLeftRight,
    BookOpenText,
    FileX,
    Scale,
    Search,
    Type,
    X,
  } from "@lucide/svelte";

  import {
    markSearchHints,
    truncateWithMark,
  } from "$lib/helpers/markSearchHints";
  import { miniSearch } from "$lib/store.svelte";

  import H2 from "$lib/components/H2.svelte";
  import Modal from "$lib/components/Modal.svelte";

  const isMac = /mac/.test(navigator.userAgent.toLowerCase());
  const ctrl = isMac ? "⌘" : "CTRL";

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
  class="btn btn-ghost flex items-center justify-start space-x-2 text-sm"
  onclick={openModal}
>
  <Search class="size-4" />
  <span>Search...</span>
  <span class="ml-4">
    <kbd class="kbd kbd-sm">{ctrl}</kbd>
    <kbd class="kbd kbd-sm">K</kbd>
  </span>
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
        <li class="list-row hover:bg-base-200">
          <a href={`#${result.slug}`} onclick={closeModal}>
            <span class="flex items-center justify-start text-lg font-bold">
              {#if result.kind === "doc"}
                <BookOpenText class="mr-2 size-4 flex-none" />
              {/if}
              {#if result.kind === "rule"}
                <Scale class="mr-2 size-4 flex-none" />
              {/if}
              {#if result.kind === "type"}
                <Type class="mr-2 size-4 flex-none" />
              {/if}
              {#if result.kind === "proc"}
                <ArrowLeftRight class="mr-2 size-4 flex-none" />
              {/if}
              {@html markSearchHints(result, result.name)}
            </span>
            <p class="truncate text-sm">
              {@html truncateWithMark(result, result.doc)}
            </p>
          </a>
        </li>
      {/each}
    </ul>
  {/if}
</Modal>
