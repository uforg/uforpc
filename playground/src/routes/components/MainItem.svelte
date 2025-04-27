<script lang="ts">
  import { onMount } from "svelte";
  import {
    ArrowLeftRight,
    BookOpenText,
    Hash,
    Scale,
    TriangleAlert,
    Type,
  } from "@lucide/svelte";
  import { getMarkdownTitle } from "$lib/helpers/getMarkdownTitle";
  import { deleteMarkdownHeadings } from "$lib/helpers/deleteMarkdownHeadings";
  import { markdownToHtml } from "$lib/helpers/markdownToHtml";
  import type { store } from "$lib/store.svelte";
  import H2 from "$lib/components/H2.svelte";

  interface Props {
    node: typeof store.jsonSchema.nodes[number];
  }

  const { node }: Props = $props();

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
    if (node.kind === "rule") return `rule-${name}`;
    if (node.kind === "type") return `type-${name}`;
    if (node.kind === "proc") return `proc-${name}`;
    if (node.kind === "doc") return `doc-${name}`;
    return "#";
  });

  let deprecatedMessage = $derived.by(() => {
    if (node.kind === "doc") return "";
    if (typeof node.deprecated !== "string") return "";
    if (node.deprecated !== "") return node.deprecated;
    return `This ${node.kind} is deprecated and it's use is not recommended.`;
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
        (
          node.kind === "rule" ||
          node.kind === "type" ||
          node.kind === "proc"
        ) &&
        typeof node.doc === "string" &&
        node.doc != ""
      ) {
        documentation = await markdownToHtml(
          deleteMarkdownHeadings(node.doc),
        );
      }
    })();
  });
</script>

<section {id}>
  <a href={`#${id}`} class="block">
    <H2 class="flex justify-start items-center group text-4xl font-extrabold">
      {#if node.kind === "doc"}
        <BookOpenText class="size-8 mr-4 flex-none group-hover:hidden" />
      {/if}
      {#if node.kind === "rule"}
        <Scale class="size-8 mr-4 flex-none group-hover:hidden" />
      {/if}
      {#if node.kind === "type"}
        <Type class="size-8 mr-4 flex-none group-hover:hidden" />
      {/if}
      {#if node.kind === "proc"}
        <ArrowLeftRight class="size-8 mr-4 flex-none group-hover:hidden" />
      {/if}

      <Hash class="size-8 mr-4 flex-none hidden group-hover:block" />

      {name}
    </H2>
  </a>

  <div class="pl-12 mt-1">
    {#if deprecatedMessage !== ""}
      <div role="alert" class="alert alert-warning mt-6 mb-4">
        <TriangleAlert class="size-4" />
        <span>Deprecated: {deprecatedMessage}</span>
      </div>
    {/if}

    {#if documentation !== ""}
      <div class="prose prose-headings:mb-0 prose-headings:mt-0 max-w-none">
        {@html documentation}
      </div>
    {/if}
  </div>
</section>
