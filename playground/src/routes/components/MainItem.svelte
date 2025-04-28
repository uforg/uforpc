<script lang="ts">
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
  import { slugify } from "$lib/helpers/slugify";
  import { store } from "$lib/store.svelte";
  import { extractNodeFromSchema } from "$lib/helpers/extractNodeFromSchema";
  import H2 from "$lib/components/H2.svelte";
  import UrpcCode from "$lib/components/UrpcCode.svelte";

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

  let urpcSchema = $state("");
  $effect(() => {
    if (node.kind === "doc") return;
    const extracted = extractNodeFromSchema(
      store.urpcSchema,
      node.kind,
      name,
    );
    if (extracted) urpcSchema = extracted;
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
      <div role="alert" class="alert alert-warning mt-4">
        <TriangleAlert class="size-4" />
        <span>Deprecated: {deprecatedMessage}</span>
      </div>
    {/if}

    {#if documentation !== ""}
      <div class="prose prose-headings:mb-0 prose-headings:mt-0 max-w-none mt-2">
        {@html documentation}
      </div>
    {/if}

    {#if urpcSchema !== ""}
      <UrpcCode code={urpcSchema} />
    {/if}
  </div>
</section>
