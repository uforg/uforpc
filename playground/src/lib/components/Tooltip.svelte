<script lang="ts">
  import "tippy.js/dist/tippy.css";
  import "tippy.js/dist/svg-arrow.css";

  import tippy from "tippy.js";
  import type { Props as TippyProps } from "tippy.js";

  interface Props {
    children: any;
    enabled?: boolean;
    content: TippyProps["content"];
    placement?: TippyProps["placement"];
    interactive?: TippyProps["interactive"];
  }

  let {
    children,
    enabled = true,
    content,
    placement = "top",
    interactive = false,
  }: Props = $props();

  let hiddenEl: HTMLTemplateElement | undefined = $state(undefined);

  const arrow = `
    <svg
      width="16"
      height="6"
      bind:this={arrow}
    >
      <path
        class="svg-arrow"
        d="M0 6s1.796-.013 4.67-3.615C5.851.9 6.93.006 8 0c1.07-.006 2.148.887 3.343 2.385C14.233 6.005 16 6 16 6H0z"
      />
      <path
        class="svg-content"
        d="m0 7s2 0 5-4c1-1 2-2 3-2 1 0 2 1 3 2 3 4 5 4 5 4h-16z"
      />
    </svg>
  `;

  $effect(() => {
    if (!enabled) return;
    if (!hiddenEl) return;

    const el = hiddenEl.nextElementSibling;
    if (!el) return;

    const inst = tippy(el, {
      content,
      placement,
      interactive,
      arrow: arrow,
    });

    return () => {
      inst.destroy();
    };
  });
</script>

<!-- 
  This element does not render anything, it's just used to reference
  the next sibling element as the tooltip target.
-->
<template bind:this={hiddenEl}></template>

{#if children}
  {@render children()}
{/if}

<style lang="postcss">
  @reference "tailwindcss";
  @plugin "daisyui";

  :global(.tippy-box) {
    @apply bg-base-200 text-base-content border border-base-content/20;
  }

  :global(.svg-content) {
    @apply fill-base-200;
  }

  :global(.svg-arrow) {
    @apply fill-base-content/30;
  }
</style>
