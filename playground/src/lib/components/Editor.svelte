<script lang="ts">
  import { onMount } from "svelte";
  import loader from "@monaco-editor/loader";
  import type * as Monaco from "monaco-editor/esm/vs/editor/editor.api";
  import { shikiToMonaco } from "@shikijs/monaco";
  import { uiStore } from "$lib/uiStore.svelte";
  import { darkTheme, getHighlighter, lightTheme } from "$lib/shiki";

  let { value = $bindable(), ...others }: { [key: string]: any } = $props();
  let editorContainer: HTMLElement;
  let monaco: typeof Monaco | null = $state(null);
  let editor: Monaco.editor.IStandaloneCodeEditor | null = $state(null);

  onMount(async () => {
    const highlighter = await getHighlighter();

    monaco = await loader.init();
    monaco.languages.register({ id: "urpc" });
    shikiToMonaco(highlighter, monaco);

    editor = monaco.editor.create(editorContainer, {
      value: value,
      language: "urpc",
      tabSize: 2,
      insertSpaces: true,
      padding: { top: 30, bottom: 30 },
    });

    editor.onDidChangeModelContent(() => {
      value = editor?.getValue();
    });
  });

  // Effect that manages the editor's value
  $effect(() => {
    if (!editor) return;
    if (value === editor.getValue()) return;
    editor.setValue(value);
  });

  // Effect that manages the editor's theme
  $effect(() => {
    if (!monaco) return;

    const themeMap = {
      system: uiStore.osTheme === "dark" ? darkTheme : lightTheme,
      light: lightTheme,
      dark: darkTheme,
    };

    monaco.editor.setTheme(themeMap[uiStore.theme]);
  });
</script>

<div {...others} bind:this={editorContainer}></div>
