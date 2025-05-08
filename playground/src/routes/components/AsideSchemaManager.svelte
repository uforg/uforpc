<script lang="ts">
  import { Save, ScrollText, WandSparkles, X } from "@lucide/svelte";
  import { onMount } from "svelte";
  import { toast } from "svelte-sonner";

  import { loadJsonSchemaFromUrpcSchemaString, store } from "$lib/store.svelte";
  import { cmdFmt } from "$lib/urpc";

  import Editor from "$lib/components/Editor.svelte";
  import Modal from "$lib/components/Modal.svelte";

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
  class="flex h-[90dvh] w-[90dvw] max-w-[1080px] flex-col space-y-6"
>
  <div class="flex w-full items-center justify-between">
    <h3 class="text-xl font-bold">Schema editor</h3>
    <button class="btn btn-circle btn-ghost" onclick={closeModal}>
      <X class="size-4" />
    </button>
  </div>

  <Editor
    bind:value={modifiedSchema}
    class="rounded-box w-full flex-grow overflow-hidden shadow-md"
  />

  <div class="flex w-full items-center justify-end space-x-2">
    <button class="btn btn-ghost" onclick={formatSchema}>
      Format <WandSparkles class="size-4" />
    </button>
    <button class="btn btn-primary" onclick={saveSchema}>
      Save <Save class="size-4" />
    </button>
  </div>
</Modal>
