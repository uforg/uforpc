<script lang="ts">
  import "tippy.js/dist/tippy.css";
  import "tippy.js/dist/svg-arrow.css";

  import tippy from "tippy.js";
  import type { Props as TippyProps } from "tippy.js";

  // the following props are set by the component and should
  // not be passed in from the parent
  type customTippyProps = Omit<
    TippyProps,
    "arrow" | "appendTo" | "triggerTarget"
  >;

  export interface Props extends Partial<customTippyProps> {
    children: any;
    enabled?: boolean;
  }

  let {
    children,
    enabled = true,
    ...tippyProps
  }: Props = $props();

  let hiddenEl: HTMLTemplateElement | undefined = $state(undefined);
  let arrow: SVGElement | undefined = $state(undefined);

  $effect(() => {
    if (!enabled) return;
    if (!hiddenEl) return;
    if (!arrow) return;

    const el = hiddenEl.nextElementSibling;
    if (!el) return;

    const inst = tippy(el, {
      ...tippyProps,
      arrow,
      appendTo: el,
      triggerTarget: el,
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

<template>
  <svg
    width="16"
    height="6"
    bind:this={arrow}
  >
    <path
      class="svg-arrow-border fill-base-content/30"
      d="M0 6s1.796-.013 4.67-3.615C5.851.9 6.93.006 8 0c1.07-.006 2.148.887 3.343 2.385C14.233 6.005 16 6 16 6H0z"
    />
    <path
      class="svg-arrow fill-base-200"
      d="m0 7s2 0 5-4c1-1 2-2 3-2 1 0 2 1 3 2 3 4 5 4 5 4h-16z"
    />
  </svg>
</template>

<style lang="postcss">
  @reference "tailwindcss";
  @plugin "daisyui";

  :global(.tippy-box) {
    @apply bg-base-200 text-base-content border border-base-content/20;
  }
</style>
