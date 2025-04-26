<script lang="ts">
  import { onMount } from "svelte";
  import loader from "@monaco-editor/loader";
  import type * as Monaco from "monaco-editor/esm/vs/editor/editor.api";
  import { createHighlighter } from "shiki";
  import type { BundledTheme } from "shiki";
  import { shikiToMonaco } from "@shikijs/monaco";
  import { store } from "$lib/store.svelte";

  const lightTheme: BundledTheme = "github-light";
  const darkTheme: BundledTheme = "github-dark";

  let { value = $bindable(), ...others }: { [key: string]: any } = $props();
  let editorContainer: HTMLElement;
  let monaco: typeof Monaco | null = $state(null);
  let editor: Monaco.editor.IStandaloneCodeEditor | null = $state(null);

  onMount(async () => {
    const urpcSyntaxUrl =
      "https://cdn.jsdelivr.net/gh/uforg/uforpc-vscode/syntaxes/urpc.tmLanguage.json";
    const urpcSyntax = await fetch(urpcSyntaxUrl);
    const urpcSyntaxJson = await urpcSyntax.json();
    urpcSyntaxJson.name = "urpc";

    const highlighter = await createHighlighter({
      langs: [urpcSyntaxJson],
      themes: [lightTheme, darkTheme],
    });

    monaco = await loader.init();
    monaco.languages.register({ id: "urpc" });
    shikiToMonaco(highlighter, monaco);

    editor = monaco.editor.create(
      editorContainer,
      {
        value: value,
        language: "urpc",
        tabSize: 2,
        insertSpaces: true,
        padding: { top: 30, bottom: 30 },
      },
    );

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

    let defaultTheme: BundledTheme = lightTheme;
    if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
      defaultTheme = darkTheme;
    } else {
      defaultTheme = lightTheme;
    }

    const themeMap: Record<string, BundledTheme> = {
      system: defaultTheme,
      light: lightTheme,
      dark: darkTheme,
    };

    monaco.editor.setTheme(themeMap[store.theme]);
  });
</script>

<div {...others} bind:this={editorContainer}></div>
