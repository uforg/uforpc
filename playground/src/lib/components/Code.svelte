<script lang="ts">
  import { toast } from "svelte-sonner";
  import { Copy } from "@lucide/svelte";
  import { darkTheme, getHighlighter, lightTheme } from "$lib/shiki";
  import { transformerColorizedBrackets } from "@shikijs/colorized-brackets";
  import { store } from "$lib/store.svelte";
  import { mergeClasses } from "$lib/helpers/mergeClasses";
  import type { ClassValue } from "$lib/helpers/mergeClasses";

  interface Props {
    code: string;
    lang: "urpc";
    class?: ClassValue;
  }

  const { code, lang, class: className }: Props = $props();

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
        transformers: [transformerColorizedBrackets()],
      });
    })();
  });

  async function copyToClipboard(text: string) {
    try {
      await navigator.clipboard.writeText(text);
      toast.success("Text copied to clipboard", { duration: 1500 });
    } catch (err) {
      console.error("Failed to copy text: ", err);
      toast.error("Failed to copy text", {
        description: `Error: ${err}`,
      });
    }
  }
</script>

{#if urpcSchemaHighlighted !== ""}
  <div class={mergeClasses("relative z-10 group", className)}>
    <button
      class="btn absolute top-2 right-2 hidden group-hover:block"
      onclick={() => copyToClipboard(code)}
    >
      <span class="flex justify-center items-center space-x-2">
        <span>Copy</span>
        <Copy class="size-4" />
      </span>
    </button>
    {@html urpcSchemaHighlighted}
  </div>
{/if}

<style lang="postcss">
  @reference "tailwindcss";
  @plugin "daisyui";

  div {
    :global(pre) {
      @apply p-4 rounded-box shadow-md border border-base-200 bg-base-200!;
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
