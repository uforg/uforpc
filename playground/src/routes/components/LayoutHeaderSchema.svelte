<script lang="ts">
  import { ScrollText, X } from "@lucide/svelte";

  import { store } from "$lib/store.svelte";

  import Code from "$lib/components/Code.svelte";
  import Modal from "$lib/components/Modal.svelte";

  let isOpen = $state(false);
  const openModal = () => (isOpen = true);
  const closeModal = () => (isOpen = false);
</script>

<button
  class="btn btn-ghost btn-block justify-start space-x-1"
  onclick={openModal}
>
  <ScrollText class="size-4" />
  <span>Show full schema</span>
</button>

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

  <Code lang="urpc" code={store.urpcSchema} class="w-full flex-grow" />
</Modal>
