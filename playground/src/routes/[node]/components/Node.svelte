<script lang="ts">
  import { TriangleAlert } from "@lucide/svelte";

  import { deleteMarkdownHeadings } from "$lib/helpers/deleteMarkdownHeadings";
  import { extractNodeFromSchema } from "$lib/helpers/extractNodeFromSchema";
  import { getMarkdownTitle } from "$lib/helpers/getMarkdownTitle";
  import { markdownToHtml } from "$lib/helpers/markdownToHtml";
  import { store } from "$lib/store.svelte";

  import Code from "$lib/components/Code.svelte";
  import H2 from "$lib/components/H2.svelte";

  import NodeQueryProc from "./NodeQuery/QueryProc.svelte";
  import NodeQueryStream from "./NodeQuery/QueryStream.svelte";

  interface Props {
    node: (typeof store.jsonSchema.nodes)[number];
  }

  const { node }: Props = $props();

  let name = $derived.by(() => {
    if (node.kind === "type") return node.name;
    if (node.kind === "proc") return node.name;
    if (node.kind === "stream") return node.name;
    if (node.kind === "doc") {
      return getMarkdownTitle(node.content);
    }

    return "unknown";
  });

  let deprecatedMessage = $derived.by(() => {
    if (node.kind === "doc") return "";
    if (typeof node.deprecated !== "string") return "";
    if (node.deprecated !== "") return node.deprecated;
    return `This ${node.kind} is deprecated and it's use is not recommended`;
  });

  let documentation = $state("");
  $effect(() => {
    (async () => {
      if (node.kind === "doc") {
        documentation = await markdownToHtml(
          deleteMarkdownHeadings(node.content),
        );
      }
      if (
        (node.kind === "stream" ||
          node.kind === "type" ||
          node.kind === "proc") &&
        typeof node.doc === "string" &&
        node.doc !== ""
      ) {
        documentation = await markdownToHtml(deleteMarkdownHeadings(node.doc));
      }
    })();
  });

  let urpcSchema = $state("");
  $effect(() => {
    if (node.kind === "doc") return;
    const extracted = extractNodeFromSchema(store.urpcSchema, node.kind, name);
    if (extracted) urpcSchema = extracted;
  });
</script>

<section class="min-h-[100dvh] space-y-12">
  <div class="prose max-w-5xl">
    <h1>{name}</h1>

    {#if deprecatedMessage !== ""}
      <div
        role="alert"
        class="alert alert-soft alert-error w-fit gap-2 font-bold italic"
      >
        <TriangleAlert class="size-4" />
        <span>Deprecated: {deprecatedMessage}</span>
      </div>
    {/if}

    {#if documentation !== ""}
      {@html documentation}
    {/if}
  </div>

  {#if node.kind === "proc"}
    <div>
      <NodeQueryProc proc={node} />
    </div>
  {/if}

  {#if node.kind === "stream"}
    <div>
      <NodeQueryStream stream={node} />
    </div>
  {/if}

  {#if urpcSchema !== ""}
    <div class="space-y-4">
      <H2>Schema</H2>
      <Code lang="urpc" code={urpcSchema} />
    </div>
  {/if}
</section>
