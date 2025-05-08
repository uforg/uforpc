<script lang="ts">
  import { goto } from "$app/navigation";
  import {
    ArrowLeftRight,
    BookOpenText,
    Hash,
    Scale,
    TriangleAlert,
    Type,
  } from "@lucide/svelte";

  import { deleteMarkdownHeadings } from "$lib/helpers/deleteMarkdownHeadings";
  import { extractNodeFromSchema } from "$lib/helpers/extractNodeFromSchema";
  import { getMarkdownTitle } from "$lib/helpers/getMarkdownTitle";
  import { markdownToHtml } from "$lib/helpers/markdownToHtml";
  import { slugify } from "$lib/helpers/slugify";
  import { store } from "$lib/store.svelte";
  import { uiStore } from "$lib/uiStore.svelte";

  import Code from "$lib/components/Code.svelte";
  import H2 from "$lib/components/H2.svelte";

  import NodeQuery from "./NodeQuery/Query.svelte";

  interface Props {
    node: (typeof store.jsonSchema.nodes)[number];
  }

  const { node }: Props = $props();

  let humanKind = $derived.by(() => {
    if (node.kind === "rule") return "validation rule";
    if (node.kind === "type") return "type";
    if (node.kind === "proc") return "procedure";
    if (node.kind === "doc") return "documentation";
    return "unknown";
  });

  let name = $derived.by(() => {
    if (node.kind === "rule") return node.name;
    if (node.kind === "type") return node.name;
    if (node.kind === "proc") return node.name;
    if (node.kind === "doc") {
      return getMarkdownTitle(node.content);
    }

    return "unknown";
  });

  let id = $derived.by(() => {
    if (node.kind === "rule") return slugify(`rule-${name}`);
    if (node.kind === "type") return slugify(`type-${name}`);
    if (node.kind === "proc") return slugify(`proc-${name}`);
    if (node.kind === "doc") return slugify(`doc-${name}`);
    return "#";
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
        (node.kind === "rule" ||
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

  // TODO: Add this to layout
  // document.title = `${name} ${humanKind} - UFO RPC Playground`;
</script>

<section {id} class="mb-[200px] min-h-[100dvh] space-y-12">
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
    <div class="mt-4">
      <NodeQuery proc={node} />
    </div>
  {/if}

  {#if urpcSchema !== ""}
    <div class="space-y-2">
      <H2>URPC Schema</H2>
      <Code lang="urpc" code={urpcSchema} collapsible={false} isOpen />
    </div>
  {/if}
</section>
