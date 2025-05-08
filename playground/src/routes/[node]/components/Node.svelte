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

<section {id} class="min-h-[100dvh]">
  <a href={`#${id}`} class="block">
    <H2 class="group flex items-center justify-start text-4xl font-extrabold">
      {#if node.kind === "doc"}
        <BookOpenText class="mr-4 size-8 flex-none group-hover:hidden" />
      {/if}
      {#if node.kind === "rule"}
        <Scale class="mr-4 size-8 flex-none group-hover:hidden" />
      {/if}
      {#if node.kind === "type"}
        <Type class="mr-4 size-8 flex-none group-hover:hidden" />
      {/if}
      {#if node.kind === "proc"}
        <ArrowLeftRight class="mr-4 size-8 flex-none group-hover:hidden" />
      {/if}

      <Hash class="mr-4 hidden size-8 flex-none group-hover:block" />

      {name}
    </H2>
  </a>

  <div class="mt-1 pl-12">
    {#if deprecatedMessage !== ""}
      <div role="alert" class="alert alert-soft alert-error mt-4 w-fit">
        <TriangleAlert class="size-4" />
        <span>Deprecated: {deprecatedMessage}</span>
      </div>
    {/if}

    {#if documentation !== ""}
      <div class="prose prose-headings:mt-0 mt-6 max-w-none">
        {@html documentation}
      </div>
    {/if}

    {#if urpcSchema !== ""}
      <Code
        class="mt-4"
        lang="urpc"
        code={urpcSchema}
        collapsible
        title={`Schema for ${name}`}
        isOpen={node.kind === "rule" || node.kind === "type"}
      />
    {/if}

    {#if node.kind === "proc"}
      <div class="mt-4">
        <NodeQuery proc={node} />
      </div>
    {/if}
  </div>
</section>
