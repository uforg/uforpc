<script lang="ts">
  import { Settings, X } from "@lucide/svelte";
  import Modal from "$lib/components/Modal.svelte";

  const isMac = /mac/.test(navigator.userAgent.toLowerCase());
  const ctrl = isMac ? "⌘" : "CTRL";

  let isOpen = $state(false);
  const openModal = () => (isOpen = true);
  const closeModal = () => (isOpen = false);

  const onKeydown = (e: KeyboardEvent) => {
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === ",") {
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
  class="btn btn-ghost flex justify-start items-center space-x-2 text-sm"
  onclick={openModal}
>
  <Settings class="size-4" />
  <span>Settings</span>
  <span class="ml-4">
    <kbd class="kbd kbd-sm">{ctrl}</kbd>
    <kbd class="kbd kbd-sm">,</kbd>
  </span>
</button>

<Modal bind:isOpen>
  <div class="w-full flex justify-between items-center">
    <h3 class="text-xl font-bold">Settings</h3>
    <button class="btn btn-circle btn-ghost" onclick={closeModal}>
      <X class="size-4" />
    </button>
  </div>

  <ul class="list mt-4">
    <li class="list-row">Test 1</li>
    <li class="list-row">Test 2</li>
    <li class="list-row">Test 3</li>
  </ul>
</Modal>
