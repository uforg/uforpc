<script lang="ts">
  import "tippy.js/dist/tippy.css";
  import "tippy.js/themes/light-border.css";
  import "tippy.js/themes/translucent.css";

  import tippy from "tippy.js";
  import type { Props as TippyProps } from "tippy.js";
  import { uiStore } from "$lib/uiStore.svelte";
  import type { Theme } from "$lib/uiStore.svelte";

  interface Props {
    children: any;
    content: TippyProps["content"];
    placement?: TippyProps["placement"];
    interactive?: TippyProps["interactive"];
    followCursor?: TippyProps["followCursor"];
  }

  let {
    children,
    content,
    placement = "top",
    interactive = false,
    followCursor = false,
  }: Props = $props();

  const lightTheme = "light-border";
  const darkTheme = "translucent";

  let hiddenEl: HTMLTemplateElement | undefined = $state(undefined);

  $effect(() => {
    if (!hiddenEl) return;

    const el = hiddenEl.nextElementSibling;
    if (!el) return;

    const themeMap: Record<Theme, string> = {
      light: lightTheme,
      dark: darkTheme,
      system: uiStore.osTheme === "dark" ? darkTheme : lightTheme,
    };

    const theme = themeMap[uiStore.theme];

    const inst = tippy(el, {
      content,
      placement,
      interactive,
      followCursor,
      theme,
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
