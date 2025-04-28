<script lang="ts">
  import { store } from "$lib/store.svelte";
  import { darkTheme, getHighlighter, lightTheme } from "$lib/shiki";

  interface Props {
    code: string;
    lang: string;
  }

  const { code, lang }: Props = $props();

  let urpcSchemaHighlighted = $state("");
  $effect(() => {
    const theme = store.theme === "dark" || store.theme === "system"
      ? darkTheme
      : lightTheme;
    const codeToHighlight = code.trim();

    (async () => {
      const highlighter = await getHighlighter();
      urpcSchemaHighlighted = highlighter.codeToHtml(codeToHighlight, {
        lang: lang,
        theme: theme,
      });
    })();
  });
</script>

{#if urpcSchemaHighlighted !== ""}
  <div>
    {@html urpcSchemaHighlighted}
  </div>
{/if}

<style lang="postcss">
  @reference "tailwindcss";
  @plugin "daisyui";

  div {
    :global(pre) {
      @apply p-4 rounded-box shadow-md border border-base-200 bg-base-100;
      @apply overflow-x-auto;
    }

    /* 
      Classes to handle line numbers
      https://github.com/shikijs/shiki/issues/3
    */

    :global(code) {
      counter-reset: step;
      counter-increment: step 0;
    }

    :global(code .line::before) {
      content: counter(step);
      counter-increment: step;
      width: 1rem;
      margin-right: 1.5rem;
      display: inline-block;
      text-align: right;
      @apply text-base-content/40;
    }
  }
</style>
