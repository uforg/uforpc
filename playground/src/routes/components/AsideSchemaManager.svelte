<script lang="ts">
  import { onMount } from "svelte";
  import { Save, ScrollText, WandSparkles, X } from "@lucide/svelte";
  import {
    loadJsonSchemaFromUrpcSchemaString,
    store,
  } from "$lib/store.svelte";
  import { cmdFmt } from "$lib/urpc";
  import Editor from "$lib/components/Editor.svelte";

  let dialog: HTMLDialogElement | null = null;
  const openDialog = () => dialog?.showModal();
  const closeDialog = () => dialog?.close();

  let modifiedSchema = $state("");
  onMount(() => {
    modifiedSchema = store.urpcSchema;
  });

  const saveSchema = async () => {
    try {
      await formatSchema();
      await loadJsonSchemaFromUrpcSchemaString(modifiedSchema);
      closeDialog();
    } catch (error) {
      alert(error);
      console.error(error);
    }
  };

  const formatSchema = async () => {
    try {
      modifiedSchema = await cmdFmt(modifiedSchema);
    } catch (error) {
      alert(error);
      console.error(error);
    }
  };
</script>

<button
  class="btn btn-ghost btn-block justify-start space-x-2"
  title="Show/Edit/Format Schema"
  onclick={openDialog}
>
  <ScrollText class="size-4" />
  <span>Schema</span>
</button>

<dialog bind:this={dialog} class="modal">
  <div class="modal-box space-y-6 w-[90dvw] max-w-[1080px] h-[90dvh] flex flex-col">
    <form method="dialog" class="w-full flex justify-between items-center">
      <h3 class="text-xl font-bold">Schema editor</h3>
      <button class="btn btn-circle btn-ghost">
        <X class="size-4" />
      </button>
    </form>

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
  </div>

  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>
