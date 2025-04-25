<script lang="ts">
  import { onMount } from "svelte";
  import loader from "@monaco-editor/loader";
  import type * as Monaco from "monaco-editor/esm/vs/editor/editor.api";

  let { value = $bindable(), ...others }: { [key: string]: any } = $props();
  let editorContainer: HTMLElement;
  let editor: Monaco.editor.IStandaloneCodeEditor;

  onMount(async () => {
    const monaco = await loader.init();
    editor = monaco.editor.create(
      editorContainer,
      {
        value: value,
        language: "urpc",
      },
    );

    editor.onDidChangeModelContent(() => {
      value = editor.getValue();
    });
  });

  $effect(() => {
    // Variable reassign to let svelte know that the value has changed
    // IDK why it's needed, but without it, the editor doesn't update
    let newValue = value;

    if (!editor) return;
    if (newValue === editor.getValue()) return;
    editor?.setValue(newValue);
  });
</script>

<div {...others} bind:this={editorContainer}></div>
