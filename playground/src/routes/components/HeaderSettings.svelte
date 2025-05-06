<script lang="ts">
  import { Link, Plus, RefreshCcw, Settings, Trash, X } from "@lucide/svelte";
  import { loadDefaultConfig, store } from "$lib/store.svelte";
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
  class="btn btn-ghost flex items-center justify-start space-x-2 text-sm"
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
  <div class="flex w-full items-center justify-between">
    <h3 class="text-xl font-bold">Settings</h3>
    <button class="btn btn-circle btn-ghost" onclick={closeModal}>
      <X class="size-4" />
    </button>
  </div>

  <p>Settings are saved in your browser's local storage.</p>

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
      <p class="label mb-1">Headers to send with requests to the endpoint.</p>

      {#each store.headers as header, index}
        <div class="mb-2 flex gap-2">
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
        <Plus class="mr-1 size-4" />
        Add Header
      </button>
    </fieldset>

    <div class="flex justify-end">
      <button class="btn btn-ghost" onclick={loadDefaultConfig}>
        <RefreshCcw class="mr-1 size-4" />
        Reset default settings
      </button>
    </div>
  </div>
</Modal>
