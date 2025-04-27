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
</script>

<section {id} class="space-y-2">
  <a href={`#${id}`} class="block">
    <H2 class="flex justify-start items-center space-x-4 group">
      {#if node.kind === "doc"}
        <BookOpenText class="size-6 flex-none group-hover:hidden" />
      {/if}
      {#if node.kind === "rule"}
        <Scale class="size-6 flex-none group-hover:hidden" />
      {/if}
      {#if node.kind === "type"}
        <Type class="size-6 flex-none group-hover:hidden" />
      {/if}
      {#if node.kind === "proc"}
        <ArrowLeftRight class="size-6 flex-none group-hover:hidden" />
      {/if}

      <Hash class="size-6 flex-none hidden group-hover:inline" />

      <span>
        {name}
      </span>
    </H2>
  </a>

  {#if deprecatedMessage !== ""}
    <div role="alert" class="alert alert-warning">
      <TriangleAlert class="size-4" />
      <span>Deprecated: {deprecatedMessage}</span>
    </div>
  {/if}
</section>
