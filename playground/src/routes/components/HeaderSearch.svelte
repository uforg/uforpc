<script lang="ts">
  import { Search, X } from "@lucide/svelte";

  const isMac = /mac/.test(navigator.userAgent.toLowerCase());
  const ctrl = isMac ? "⌘" : "CTRL";

  let dialog: HTMLDialogElement | null = null;
  let input: HTMLInputElement | null = null;

  function toggleDialog(isOpen: boolean) {
    if (!dialog) {
      return;
    }
    if (isOpen) {
      dialog.showModal();
      setTimeout(() => {
        input?.focus();
      }, 100);
    }
    if (!isOpen) {
      dialog.close();
    }
  }

  const onKeydown = (e: KeyboardEvent) => {
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === "k") {
      e.preventDefault();
      toggleDialog(true);
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
  class="input input-ghost focus:outline-none"
  onclick={() => toggleDialog(true)}
>
  <Search class="size-4" />
  <span>Search...</span>
  <span class="ml-4">
    <kbd class="kbd kbd-sm">{ctrl}</kbd>
    <kbd class="kbd kbd-sm">K</kbd>
  </span>
</button>

<dialog bind:this={dialog} class="modal">
  <div class="modal-box">
    <div class="flex justify-start items-center space-x-2">
      <label class="input flex-grow">
        <Search class="size-4" />
        <input bind:this={input} type="search" placeholder="Search..." />
      </label>
      <form method="dialog">
        <button class="btn btn-square">
          <X class="size-4" />
        </button>
      </form>
    </div>

    <ul class="list mt-4">
      <li class="list-row">Test 1</li>
      <li class="list-row">Test 2</li>
      <li class="list-row">Test 3</li>
    </ul>
  </div>

  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>
