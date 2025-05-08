<script lang="ts">
  import { getMarkdownTitle } from "$lib/helpers/getMarkdownTitle";
  import { slugify } from "$lib/helpers/slugify";
  import { store } from "$lib/store.svelte";

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
</script>

{#if !nodeExists}
  <NotFound />
{/if}

{#if nodeExists}
  <Node node={store.jsonSchema.nodes[nodeIndex]} />
{/if}
