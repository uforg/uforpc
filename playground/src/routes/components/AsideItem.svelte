<script lang="ts">
  import { getMarkdownTitle } from "$lib/helpers/getMarkdownTitle";
  import type { store } from "$lib/store.svelte";
  import {
    ArrowLeftRight,
    BookOpenText,
    Scale,
    TriangleAlert,
    Type,
  } from "@lucide/svelte";

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

  let title = $derived.by(() => {
    if (node.kind === "rule") return `Rule: ${name}`;
    if (node.kind === "type") return `Type: ${name}`;
    if (node.kind === "proc") return `Procedure: ${name}`;
    if (node.kind === "doc") return `Documentation: ${name}`;
    return "Unknown";
  });

  let href = $derived.by(() => {
    if (node.kind === "rule") return `#rule-${name}`;
    if (node.kind === "type") return `#type-${name}`;
    if (node.kind === "proc") return `#proc-${name}`;
    if (node.kind === "doc") return `#doc-${name}`;
    return "#";
  });

  let isDeprecated = $derived.by(() => {
    if (node.kind === "doc") return false;
    if (typeof node.deprecated === "string") return true;
    return false;
  });
</script>

<a
  {href}
  {title}
  class={[
    "btn btn-ghost btn-block justify-start space-x-2",
    {
      "tooltip tooltip-top": isDeprecated,
      "hover:bg-blue-500/20": node.kind === "doc",
      "hover:bg-yellow-500/20": node.kind === "rule",
      "hover:bg-purple-500/20": node.kind === "type",
      "hover:bg-green-500/20": node.kind === "proc",
    },
  ]}
  data-tip={isDeprecated ? "Deprecated" : ""}
>
  {#if node.kind === "doc"}
    <BookOpenText class="flex-none size-4" />
  {/if}
  {#if node.kind === "rule"}
    <Scale class="flex-none size-4" />
  {/if}
  {#if node.kind === "type"}
    <Type class="flex-none size-4" />
  {/if}
  {#if node.kind === "proc"}
    <ArrowLeftRight class="flex-none size-4" />
  {/if}

  <span
    class={[
      "overflow-hidden overflow-ellipsis",
      {
        "line-through": isDeprecated,
      },
    ]}
  >
    {name}
  </span>

  {#if isDeprecated}
    <TriangleAlert class="flex-none size-4 text-warning" />
  {/if}
</a>
