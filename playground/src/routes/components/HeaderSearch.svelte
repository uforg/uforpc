<script lang="ts">
  import { Search, X } from "@lucide/svelte";
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
</script>

<button
  class="input input-ghost focus:outline-none ps-0"
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
  <div class="flex justify-start items-center space-x-2">
    <label class="input flex-grow">
      <Search class="size-4" />
      <input bind:this={input} type="search" placeholder="Search..." />
    </label>
    <button class="btn btn-square" onclick={closeModal}>
      <X class="size-4" />
    </button>
  </div>

  <ul class="list mt-4">
    <li class="list-row">Test 1</li>
    <li class="list-row">Test 2</li>
    <li class="list-row">Test 3</li>
  </ul>
</Modal>
