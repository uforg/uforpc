<script lang="ts">
  import { ChevronDown, ChevronRight, Copy, ScrollText } from "@lucide/svelte";
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
    collapsible?: boolean;
    isOpen?: boolean;
    title?: string;
    rounded?: boolean;
    withBorder?: boolean;
  }

  let {
    code,
    lang,
    class: className,
    collapsible = false,
    isOpen = $bindable(true),
    title = "Code",
    rounded = true,
    withBorder = true,
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

  function toggleCollapse() {
    isOpen = !isOpen;
  }
</script>

{#if urpcSchemaHighlighted !== ""}
  <div class={mergeClasses("group", className)}>
    {#if collapsible}
      <button
        class={[
          "btn w-full justify-start",
          "group/btn",
          {
            "border-base-content/20 border": withBorder,
            "rounded-box": rounded,
            "rounded-b-none": isOpen,
          },
        ]}
        onclick={toggleCollapse}
      >
        <ScrollText class="mr-2 block size-4 group-hover/btn:hidden" />
        {#if isOpen}
          <ChevronDown class="mr-2 hidden size-4 group-hover/btn:block" />
        {:else}
          <ChevronRight class="mr-2 hidden size-4 group-hover/btn:block" />
        {/if}
        {#if title}
          {title}
        {/if}
      </button>
    {/if}

    {#if !collapsible || isOpen}
      <div
        class={[
          "relative z-10 p-4",
          "bg-base-200",
          {
            "border-base-content/20 border": withBorder,
            "rounded-box": rounded,
            "rounded-t-none border-t-0 pt-2": collapsible,
          },
        ]}
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
