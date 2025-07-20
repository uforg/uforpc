<script lang="ts">
  import { fade } from "svelte/transition";

  import LogoUfo from "$lib/components/LogoUfo.svelte";
  import SwaggerLogo from "$lib/components/SwaggerLogo.svelte";
  import Tooltip from "$lib/components/Tooltip.svelte";

  import LayoutSwaggerUi from "./LayoutSwaggerUi.svelte";

  const animationDuration = 200;

  let isAnimating = $state(false);
  let isOpen = $state(false);
  let isOpenAnimated = $state(false);
  let opacity = $state(0);

  const toggle = async () => {
    if (isAnimating) return;

    isOpen = !isOpen;
    isAnimating = true;
    if (isOpenAnimated) {
      opacity = 0;
      await new Promise((resolve) => setTimeout(resolve, animationDuration));
      isOpenAnimated = false;
    } else {
      isOpenAnimated = true;
      await new Promise((resolve) => setTimeout(resolve, 10)); // Wait for repaint
      opacity = 1;
    }
    isAnimating = false;
  };

  let tooltipContent = $derived(
    isOpenAnimated ? "Switch to UFO RPC" : "Switch to Swagger UI (OpenAPI)",
  );
</script>

<Tooltip content={tooltipContent} placement="left">
  <button
    class={{
      "group btn btn-circle btn-lg fixed right-4 bottom-4 z-50": true,
      "bg-base-300 border-base-content/20": isOpen,
      "btn-ghost bg-transparent": !isOpen,
    }}
    onclick={toggle}
  >
    {#if !isOpen}
      <span in:fade={{ duration: animationDuration }}>
        <SwaggerLogo class="w-full" />
      </span>
    {/if}
    {#if isOpen}
      <span in:fade={{ duration: animationDuration }}>
        <LogoUfo class="size-10" />
      </span>
    {/if}
  </button>
</Tooltip>

<div
  class={{
    "bg-base-100 fixed top-0 left-0 z-40 h-[100dvh] w-[100dvw] overflow-y-auto": true,
    "transition-opacity": true,
    "translate-x-[300%]": !isOpenAnimated,
  }}
  style="opacity: {opacity}; transition-duration: {animationDuration}ms;"
>
  <LayoutSwaggerUi />
</div>
