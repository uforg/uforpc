<script lang="ts">
  import { onMount } from "svelte";
  import { dimensionschangeAction, uiStore } from "$lib/uiStore.svelte";
  import { store } from "$lib/store.svelte";
  import AsideSchemaManager from "./AsideSchemaManager.svelte";
  import AsideItem from "./AsideItem.svelte";

  // if has hash anchor navigate to it
  onMount(async () => {
    // wait 500ms to ensure the content is rendered
    await new Promise((resolve) => setTimeout(resolve, 500));

    if (window.location.hash) {
      const element = document.getElementById(
        "navlink-" + window.location.hash.slice(1),
      );
      if (element) {
        element.scrollIntoView({ behavior: "smooth" });
      }
    }
  });
</script>

<aside
  use:dimensionschangeAction
  ondimensionschange={(e) => uiStore.aside = e.detail}
  class={[
    "flex-none w-full max-w-[280px] h-[100dvh] overflow-x-hidden overflow-y-auto scroll-p-[90px]",
  ]}
>
  <a
    class={[
      "flex space-x-2 items-center whitespace-nowrap",
      "h-[72px] w-full sticky top-0 p-4 z-10",
      "bg-base-100/90 backdrop-blur-sm",
      {
        "shadow-xs": uiStore.aside.scroll.isTopScrolled,
      },
    ]}
    href="https://uforpc.uforg.dev"
    target="_blank"
  >
    <img src="/assets/logo.png" alt="UFO RPC Logo" class="h-full">
    <h1 class="font-bold">UFO RPC Playground</h1>
  </a>
  <nav class="p-4">
    <AsideSchemaManager />
    {#each store.jsonSchema.nodes as node}
      <AsideItem {node} />
    {/each}
  </nav>
</aside>
