<script lang="ts">
  import { page } from "$app/state";
  import { Home } from "@lucide/svelte";
  import { onMount } from "svelte";

  import { getMarkdownTitle } from "$lib/helpers/getMarkdownTitle";
  import { store } from "$lib/store.svelte";
  import { dimensionschangeAction, uiStore } from "$lib/uiStore.svelte";
  import type { Schema } from "$lib/urpcTypes";

  import Logo from "$lib/components/Logo.svelte";
  import Tooltip from "$lib/components/Tooltip.svelte";

  import LayoutAsideFilters from "./LayoutAsideFilters.svelte";
  import LayoutAsideItem from "./LayoutAsideItem.svelte";
  import LayoutAsideSchema from "./LayoutAsideSchema.svelte";

  // if has hash anchor navigate to it
  onMount(async () => {
    // wait 500ms to ensure the content is rendered
    await new Promise((resolve) => setTimeout(resolve, 500));

    if (window.location.hash) {
      const element = document.getElementById(
        `navlink-${window.location.hash.slice(1)}`,
      );
      if (element) {
        element.scrollIntoView({ behavior: "smooth" });
      }
    }
  });

  let isHome = $derived(page.url.hash === "" || page.url.hash === "#/");

  function getNodeName(node: Schema["nodes"][number]) {
    if (node.kind === "type") return node.name;
    if (node.kind === "proc") return node.name;
    if (node.kind === "stream") return node.name;
    if (node.kind === "doc") return getMarkdownTitle(node.content);
    return "unknown";
  }

  function shouldShowNode(
    kind: string,
    node: Schema["nodes"][number],
  ): boolean {
    if (node.kind !== kind) return false;

    if (uiStore.asideSearchOpen) {
      if (uiStore.asideSearchQuery === "") return true;

      // Do the search
      const name = getNodeName(node).toLowerCase();
      const query = uiStore.asideSearchQuery.toLowerCase();
      return name.includes(query);
    }

    if (node.kind === "doc") return !uiStore.asideHideDocs;
    if (node.kind === "type") return !uiStore.asideHideTypes;
    if (node.kind === "proc") return !uiStore.asideHideProcs;
    if (node.kind === "stream") return !uiStore.asideHideStreams;

    return false;
  }
</script>

<aside
  use:dimensionschangeAction
  ondimensionschange={(e) => (uiStore.aside = e.detail)}
  class={[
    "h-[100dvh] w-full max-w-[280px] flex-none scroll-p-[90px]",
    "overflow-x-hidden overflow-y-auto",
  ]}
>
  <header
    class={[
      "bg-base-100/90 sticky top-0 z-10 w-full backdrop-blur-sm",
      {
        "shadow-xs": uiStore.aside.scroll.isTopScrolled,
      },
    ]}
  >
    <a
      class={[
        "sticky top-0 z-10 flex h-[72px] w-full items-end p-4",
        {
          "shadow-xs": uiStore.aside.scroll.isTopScrolled,
        },
      ]}
      href="https://uforpc.uforg.dev"
      target="_blank"
    >
      <Logo
        class="mx-auto h-full"
        animateAuto
        animateAutoSpeed={2}
        animateHover
        animateHoverSpeed={0.5}
      />
    </a>

    <LayoutAsideFilters />
  </header>

  <nav class="p-4">
    <LayoutAsideSchema />

    <Tooltip content="RPC Home">
      <a
        href="#/"
        class={[
          "btn btn-ghost btn-block justify-start space-x-2 border-transparent",
          "hover:bg-blue-500/20",
          { "bg-blue-500/20": isHome },
        ]}
      >
        <Home class="size-4" />
        <span>Home</span>
      </a>
    </Tooltip>

    {#each store.jsonSchema.nodes as node}
      {#if shouldShowNode("doc", node)}
        <LayoutAsideItem {node} />
      {/if}
      {#if shouldShowNode("type", node)}
        <LayoutAsideItem {node} />
      {/if}
      {#if shouldShowNode("proc", node)}
        <LayoutAsideItem {node} />
      {/if}
      {#if shouldShowNode("stream", node)}
        <LayoutAsideItem {node} />
      {/if}
    {/each}
  </nav>
</aside>
