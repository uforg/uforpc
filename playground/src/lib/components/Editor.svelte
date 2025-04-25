<script lang="ts">
  import { onMount } from "svelte";
  import loader from "@monaco-editor/loader";
  import type * as Monaco from "monaco-editor/esm/vs/editor/editor.api";
  import { createHighlighter } from "shiki";
  import { shikiToMonaco } from "@shikijs/monaco";

  let { value = $bindable(), ...others }: { [key: string]: any } = $props();
  let editorContainer: HTMLElement;
  let editor: Monaco.editor.IStandaloneCodeEditor;

  onMount(async () => {
    const urpcSyntaxUrl =
      "https://cdn.jsdelivr.net/gh/uforg/uforpc-vscode/syntaxes/urpc.tmLanguage.json";
    const urpcSyntax = await fetch(urpcSyntaxUrl);
    const urpcSyntaxJson = await urpcSyntax.json();
    urpcSyntaxJson.name = "urpc";

    const highlighter = await createHighlighter({
      langs: [urpcSyntaxJson],
      themes: ["dracula", "andromeeda", "github-light"],
    });

    const monaco = await loader.init();
    monaco.languages.register({ id: "urpc" });
    shikiToMonaco(highlighter, monaco);

    editor = monaco.editor.create(
      editorContainer,
      {
        value: value,
        language: "urpc",
        theme: "dracula",
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
