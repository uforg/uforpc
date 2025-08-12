<script lang="ts">
  import { Link, Plus, RefreshCcw, Settings, Trash, X } from "@lucide/svelte";

  import { loadDefaultConfig, store } from "$lib/store.svelte";
  import { uiStore } from "$lib/uiStore.svelte";

  import Modal from "$lib/components/Modal.svelte";
  import Tooltip from "$lib/components/Tooltip.svelte";

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
    store.headers = [...store.headers, { key: "", value: "", enabled: true }];
  };

  const removeHeader = (index: number) => {
    store.headers = store.headers.filter((_, i) => i !== index);
  };

  const loadDefaultConfigConfirm = () => {
    if (confirm("Are you sure you want to reset the default settings?")) {
      loadDefaultConfig();
    }
  };
</script>

<button
  class="btn btn-ghost flex items-center justify-start space-x-1 text-sm"
  onclick={openModal}
>
  <Settings class="size-4" />
  <span>Settings</span>
  {#if !uiStore.isMobile}
    <span class="ml-4">
      <kbd class="kbd kbd-sm">{ctrl}</kbd>
      <kbd class="kbd kbd-sm">,</kbd>
    </span>
  {/if}
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
      <legend class="fieldset-legend">Base URL</legend>
      <label class="input w-full">
        <Link class="size-4" />
        <input
          type="url"
          class="grow"
          spellcheck="false"
          placeholder="https://example.com/api/v1/urpc"
          bind:value={store.baseUrl}
        />
      </label>
      <p class="label text-wrap">
        This is the base URL where the UFO RPC server is running, all requests
        will be sent to {`<base-url>/{operationName}`}
        where {`{operationName}`} is the name of the procedure or stream you want
        to call.
      </p>
    </fieldset>

    <fieldset class="fieldset">
      <legend class="fieldset-legend">Headers</legend>
      <p class="label mb-1">Headers to send with requests to the endpoint.</p>

      {#each store.headers as header, index}
        <div class="mb-2 flex items-center gap-2">
          <Tooltip
            content={header.enabled ? "Disable header" : "Enable header"}
          >
            <input
              type="checkbox"
              class="toggle"
              bind:checked={header.enabled}
            />
          </Tooltip>
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
          <Tooltip content="Remove header">
            <button
              class="btn btn-square btn-ghost btn-error"
              onclick={() => removeHeader(index)}
            >
              <Trash class="size-4" />
            </button>
          </Tooltip>
        </div>
      {/each}

      <button class="btn btn-outline mt-2" onclick={addHeader}>
        <Plus class="mr-1 size-4" />
        Add Header
      </button>
    </fieldset>

    <div class="flex justify-end">
      <button class="btn btn-ghost" onclick={loadDefaultConfigConfirm}>
        <RefreshCcw class="mr-1 size-4" />
        Reset default settings
      </button>
    </div>
  </div>
</Modal>
