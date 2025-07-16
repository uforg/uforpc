<script lang="ts">
  import { Copy } from "@lucide/svelte";
  import { transformerColorizedBrackets } from "@shikijs/colorized-brackets";
  import { toast } from "svelte-sonner";
  import { slide } from "svelte/transition";

  import { mergeClasses } from "$lib/helpers/mergeClasses";
  import type { ClassValue } from "$lib/helpers/mergeClasses";
  import {
    darkTheme,
    getHighlighter,
    getOrFallbackLanguage,
    lightTheme,
  } from "$lib/shiki";
  import { uiStore } from "$lib/uiStore.svelte";

  interface Props {
    code: string;
    lang: "urpc" | string;
    class?: ClassValue;
    rounded?: boolean;
    withBorder?: boolean;
    scrollY?: boolean;
    scrollX?: boolean;
  }

  let {
    code,
    lang,
    class: className,
    rounded = true,
    withBorder = true,
    scrollY = true,
    scrollX = true,
  }: Props = $props();

  let urpcSchemaHighlighted = $state("");
  $effect(() => {
    const themeMap = {
      dark: darkTheme,
      light: lightTheme,
    };
    let theme = themeMap[uiStore.theme];

    const codeToHighlight = code.trim();

    (async () => {
      const highlighter = await getHighlighter();
      urpcSchemaHighlighted = highlighter.codeToHtml(codeToHighlight, {
        lang: getOrFallbackLanguage(lang),
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
  <div
    class={mergeClasses([
      "group bg-base-200 relative z-10 p-4",
      {
        "overflow-y-auto": scrollY,
        "overflow-x-auto": scrollX,
        "border-base-content/20 border": withBorder,
        "rounded-box": rounded,
      },
      className,
    ])}
    transition:slide={{ duration: 100 }}
  >
    <button
      class="btn absolute top-4 right-4 hidden group-hover:block"
      onclick={() => copyToClipboard(code)}
    >
      <span class="flex items-center justify-center space-x-2">
        <Copy class="size-4" />
        <span>Copy</span>
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
      @apply bg-base-200!;
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
