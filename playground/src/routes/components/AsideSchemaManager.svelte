<script lang="ts">
  import { onMount } from "svelte";
  import { Save, ScrollText, X } from "@lucide/svelte";
  import {
    loadJsonSchemaFromUrpcSchemaString,
    store,
  } from "../store.svelte";
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
      await loadJsonSchemaFromUrpcSchemaString(modifiedSchema);
      closeDialog();
    } catch (error) {
      alert(error);
      console.error(error);
    }
  };
</script>

<button
  class="btn btn-ghost btn-block justify-start space-x-2"
  title="Show/Edit Schema"
  onclick={openDialog}
>
  <ScrollText class="size-4" />
  <span>Schema</span>
</button>

<dialog bind:this={dialog} class="modal">
  <div class="modal-box space-y-6 w-[90dvw] max-w-[1080px]">
    <form method="dialog" class="w-full flex justify-between items-center">
      <h3 class="text-xl font-bold">Schema editor</h3>
      <button class="btn btn-circle btn-ghost">
        <X class="size-4" />
      </button>
    </form>

    <Editor bind:value={modifiedSchema} class="w-full h-[600px]" />

    <div class="w-full flex justify-end items-center">
      <button class="btn btn-primary" onclick={saveSchema}>
        Save <Save class="size-4" />
      </button>
    </div>
  </div>

  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>
