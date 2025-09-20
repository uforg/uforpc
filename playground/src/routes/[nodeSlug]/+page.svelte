<script lang="ts">
  import { untrack } from "svelte";

  import { getMarkdownTitle } from "$lib/helpers/getMarkdownTitle";
  import { slugify } from "$lib/helpers/slugify";
  import { store } from "$lib/store.svelte";
  import { uiStore } from "$lib/uiStore.svelte";

  import type { PageProps } from "./$types";
  import Node from "./components/Node.svelte";
  import NotFound from "./components/NotFound.svelte";

  let { data }: PageProps = $props();

  let nodeIndex = $derived.by(() => {
    for (const [index, node] of store.jsonSchema.nodes.entries()) {
      if (node.kind !== data.nodeKind) continue;

      const isDoc = node.kind === "doc";
      let nodeName = isDoc ? getMarkdownTitle(node.content) : node.name;
      nodeName = slugify(nodeName);

      if (data.nodeName === nodeName) return index;
    }

    return -1; // Node not found in store
  });

  let nodeExists = $derived(nodeIndex !== -1);

  let node = $derived(store.jsonSchema.nodes[nodeIndex]);

  let name = $derived.by(() => {
    if (node.kind === "type") return node.name;
    if (node.kind === "proc") return node.name;
    if (node.kind === "stream") return node.name;
    if (node.kind === "doc") {
      return getMarkdownTitle(node.content);
    }

    return "unknown";
  });

  let humanKind = $derived.by(() => {
    if (node.kind === "type") return "type";
    if (node.kind === "proc") return "procedure";
    if (node.kind === "stream") return "stream";
    if (node.kind === "doc") return "documentation";
    return "unknown";
  });

  let title = $derived.by(() => {
    if (!nodeExists) return "UFO RPC Playground";

    return `${name} ${humanKind} - UFO RPC Playground`;
  });

  // Scroll to top of page when node changes
  $effect(() => {
    nodeIndex; // Just to add a dependency to trigger the effect
    untrack(() => {
      // Untrack the uiStore.contentWrapper.element to avoid infinite loop
      uiStore.contentWrapper.element?.scrollTo({ top: 0, behavior: "smooth" });
    });
  });
</script>

<svelte:head>
  <title>{title}</title>
</svelte:head>

{#if !nodeExists}
  <NotFound />
{/if}

{#if nodeExists}
  {#key nodeIndex}
    <Node nodeSlug={data.nodeSlug} {node} />
  {/key}
{/if}
