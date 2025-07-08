<script lang="ts">
  import { ScrollText, X } from "@lucide/svelte";

  import { store } from "$lib/store.svelte";

  import Code from "$lib/components/Code.svelte";
  import Modal from "$lib/components/Modal.svelte";
  import Tooltip from "$lib/components/Tooltip.svelte";

  let isOpen = $state(false);
  const openModal = () => (isOpen = true);
  const closeModal = () => (isOpen = false);
</script>

<Tooltip content="Show full schema">
  <button
    class="btn btn-ghost btn-block justify-start space-x-2 border-transparent"
    onclick={openModal}
  >
    <ScrollText class="size-4" />
    <span>Schema</span>
  </button>
</Tooltip>

<Modal
  bind:isOpen
  class="flex h-[90dvh] w-[90dvw] max-w-[1080px] flex-col space-y-4"
>
  <div class="flex w-full items-center justify-between">
    <h3 class="text-xl font-bold">Full schema</h3>
    <button class="btn btn-circle btn-ghost" onclick={closeModal}>
      <X class="size-4" />
    </button>
  </div>

  <Code
    lang="urpc"
    code={store.urpcSchema}
    class="w-full flex-grow overflow-auto"
  />
</Modal>
