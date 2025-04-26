<script lang="ts">
  import { Link, Plus, Settings, Trash, X } from "@lucide/svelte";
  import { store } from "$lib/store.svelte";
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

  const addHeader = () => {
    store.headers = [...store.headers, { key: "", value: "" }];
  };

  const removeHeader = (index: number) => {
    store.headers = store.headers.filter((_, i) => i !== index);
  };
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

  <p>
    Settings are saved in your browser's local storage.
  </p>

  <div class="mt-4 space-y-4">
    <fieldset class="fieldset">
      <legend class="fieldset-legend">Endpoint</legend>
      <label class="input w-full">
        <Link class="size-4" />
        <input
          type="url"
          class="grow"
          spellcheck="false"
          placeholder="https://example.com/api/v1/urpc"
          bind:value={store.endpoint}
        />
      </label>
      <p class="label">The endpoint where the UFO RPC server is running.</p>
    </fieldset>

    <fieldset class="fieldset">
      <legend class="fieldset-legend">Headers</legend>

      {#each store.headers as header, index}
        <div class="flex gap-2 mb-2">
          <input
            type="text"
            class="input flex-1"
            spellcheck="false"
            placeholder="Key"
            bind:value={header.key}
          />
          <input
            type="text"
            class="input flex-1"
            spellcheck="false"
            placeholder="Value"
            bind:value={header.value}
          />
          <button
            class="btn btn-square btn-ghost btn-error"
            onclick={() => removeHeader(index)}
            title="Remove header"
          >
            <Trash class="size-4" />
          </button>
        </div>
      {/each}

      <button class="btn btn-outline mt-2" onclick={addHeader}>
        <Plus class="size-4 mr-1" />
        Add Header
      </button>
      <p class="label mt-2">Headers to send with requests to the endpoint.</p>
    </fieldset>
  </div>
</Modal>
