<script lang="ts">
  import { store } from "$lib/store.svelte";
  import { darkTheme, getHighlighter, lightTheme } from "$lib/shiki";

  interface Props {
    code: string;
  }

  const { code }: Props = $props();

  let urpcSchemaHighlighted = $state("");
  $effect(() => {
    const theme = store.theme === "dark" || store.theme === "system"
      ? darkTheme
      : lightTheme;
    const codeToHighlight = code.trim();

    (async () => {
      const highlighter = await getHighlighter();
      urpcSchemaHighlighted = highlighter.codeToHtml(codeToHighlight, {
        lang: "urpc",
        theme: theme,
      });
    })();
  });
</script>

{#if urpcSchemaHighlighted !== ""}
  <div class="rounded-box shadow-md p-4 overflow-x-auto">
    {@html urpcSchemaHighlighted}
  </div>
{/if}
