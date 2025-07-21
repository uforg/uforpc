<script lang="ts">
  import { Copy, Download, EllipsisVertical } from "@lucide/svelte";
  import { transformerColorizedBrackets } from "@shikijs/colorized-brackets";
  import { toast } from "svelte-sonner";

  import { getLangExtension } from "$lib/helpers/getLangExtension";
  import { mergeClasses } from "$lib/helpers/mergeClasses";
  import type { ClassValue } from "$lib/helpers/mergeClasses";
  import {
    darkTheme,
    getHighlighter,
    getOrFallbackLanguage,
    lightTheme,
  } from "$lib/shiki";
  import { uiStore } from "$lib/uiStore.svelte";

  import Menu from "./Menu.svelte";

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

  let codeHighlighted = $state("");
  $effect(() => {
    const themeMap = {
      dark: darkTheme,
      light: lightTheme,
    };
    let theme = themeMap[uiStore.theme];

    const codeToHighlight = code.trim();

    (async () => {
      const highlighter = await getHighlighter();
      codeHighlighted = highlighter.codeToHtml(codeToHighlight, {
        lang: getOrFallbackLanguage(lang),
        theme: theme,
        transformers: [transformerColorizedBrackets()],
      });
    })();
  });

  async function copyToClipboard() {
    try {
      await navigator.clipboard.writeText(code);
      toast.success("Code copied to clipboard", { duration: 1500 });
    } catch (err) {
      console.error("Failed to copy code: ", err);
      toast.error("Failed to copy code", {
        description: `Error: ${err}`,
      });
    }
  }

  const downloadCode = () => {
    try {
      // Create a blob from the code
      const blob = new Blob([code], { type: "text/plain" });
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");

      // Find extension from lang
      const extension = getLangExtension(lang);
      const fileName = `code.${extension}`;

      // Download the file
      a.href = url;
      a.download = fileName;
      a.click();

      toast.success("Code downloaded", { duration: 1500 });
    } catch (error) {
      console.error("Failed to download code: ", error);
      toast.error("Failed to download code", {
        description: `Error: ${error}`,
      });
    }
  };
</script>

{#snippet menuContent()}
  <div class="flex flex-col space-y-1">
    <button
      class="btn btn-ghost btn-sm justify-start"
      onclick={() => copyToClipboard()}
    >
      <Copy class="size-4" />
      <span>Copy to clipboard</span>
    </button>
    <button
      class="btn btn-ghost btn-sm justify-start"
      onclick={() => downloadCode()}
    >
      <Download class="size-4" />
      <span>Download</span>
    </button>
  </div>
{/snippet}

{#if codeHighlighted !== ""}
  <div
    class={mergeClasses([
      "bg-base-200 relative z-10 p-4",
      {
        "overflow-y-auto": scrollY,
        "overflow-x-auto": scrollX,
        "border-base-content/20 border": withBorder,
        "rounded-box": rounded,
      },
      className,
    ])}
  >
    <div
      class={["sticky top-0 left-0 z-10 flex justify-end", "-mb-6 h-6 w-full"]}
    >
      <Menu
        content={menuContent}
        trigger="mouseenter focus"
        placement="left-start"
      >
        <button class="btn btn-sm btn-square">
          <EllipsisVertical class="size-4" />
        </button>
      </Menu>
    </div>
    {@html codeHighlighted}
  </div>
{/if}

<style lang="postcss">
  @reference "tailwindcss";
  @plugin "daisyui";

  div {
    :global(pre) {
      @apply bg-base-200!;
    }

    :global(pre:focus-visible) {
      @apply outline-none;
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
