<script lang="ts">
  import { onMount } from "svelte";
  import { Save, ScrollText, WandSparkles, X } from "@lucide/svelte";
  import { toast } from "svelte-sonner";
  import {
    loadJsonSchemaFromUrpcSchemaString,
    store,
  } from "$lib/store.svelte";
  import { cmdFmt } from "$lib/urpc";
  import Modal from "$lib/components/Modal.svelte";
  import Editor from "$lib/components/Editor.svelte";

  let isOpen = $state(false);
  const openModal = () => (isOpen = true);
  const closeModal = () => (isOpen = false);

  let modifiedSchema = $state("");
  onMount(() => {
    modifiedSchema = store.urpcSchema;
  });

  const saveSchema = async () => {
    try {
      await formatSchema();
      await loadJsonSchemaFromUrpcSchemaString(modifiedSchema);
      closeModal();
    } catch (error: unknown) {
      toast.error("Failed to save schema", {
        description: `${error}`,
      });
      console.error(error);
    }
  };

  const formatSchema = async () => {
    try {
      modifiedSchema = await cmdFmt(modifiedSchema);
    } catch (error: unknown) {
      toast.error("Failed to format schema", {
        description: `${error}`,
      });
      console.error(error);
    }
  };
</script>

<button
  class="btn btn-ghost btn-block justify-start space-x-2"
  title="Show/Edit/Format Schema"
  onclick={openModal}
>
  <ScrollText class="size-4" />
  <span>Schema</span>
</button>

<Modal
  bind:isOpen
  class="space-y-6 w-[90dvw] max-w-[1080px] h-[90dvh] flex flex-col"
>
  <div class="w-full flex justify-between items-center">
    <h3 class="text-xl font-bold">Schema editor</h3>
    <button class="btn btn-circle btn-ghost" onclick={closeModal}>
      <X class="size-4" />
    </button>
  </div>

  <Editor
    bind:value={modifiedSchema}
    class="w-full flex-grow rounded-box shadow-md overflow-hidden"
  />

  <div class="w-full flex justify-end items-center space-x-2">
    <button class="btn btn-ghost" onclick={formatSchema}>
      Format <WandSparkles class="size-4" />
    </button>
    <button class="btn btn-primary" onclick={saveSchema}>
      Save <Save class="size-4" />
    </button>
  </div>
</Modal>
