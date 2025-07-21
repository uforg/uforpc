<script lang="ts">
  import { Braces, X } from "@lucide/svelte";

  import Code from "$lib/components/Code.svelte";
  import Modal from "$lib/components/Modal.svelte";

  let isOpen = $state(false);
  const openModal = () => (isOpen = true);
  const closeModal = () => (isOpen = false);

  let openApiSchema = $state("");
  $effect(() => {
    fetch("./openapi.yaml")
      .then((res) => res.text())
      .then((text) => (openApiSchema = text));
  });
</script>

<button
  class="btn btn-ghost btn-block justify-start space-x-1"
  onclick={openModal}
>
  <Braces class="size-4" />
  <span>OpenAPI schema</span>
</button>

<Modal
  bind:isOpen
  class="flex h-[90dvh] w-[90dvw] max-w-[1080px] flex-col space-y-4"
>
  <div class="flex w-full items-center justify-between">
    <h3 class="text-xl font-bold">OpenAPI schema</h3>
    <button class="btn btn-circle btn-ghost" onclick={closeModal}>
      <X class="size-4" />
    </button>
  </div>

  <Code lang="yaml" code={openApiSchema} class="w-full flex-grow" />
</Modal>
