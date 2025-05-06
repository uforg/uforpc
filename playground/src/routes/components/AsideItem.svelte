<script lang="ts">
  import {
    ArrowLeftRight,
    BookOpenText,
    Scale,
    TriangleAlert,
    Type,
  } from "@lucide/svelte";
  import { store } from "$lib/store.svelte";
  import { uiStore } from "$lib/uiStore.svelte";
  import { getMarkdownTitle } from "$lib/helpers/getMarkdownTitle";
  import { slugify } from "$lib/helpers/slugify";
  import Tooltip from "$lib/components/Tooltip.svelte";

  interface Props {
    node: (typeof store.jsonSchema.nodes)[number];
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
    const deprecated = isDeprecated ? " (Deprecated)" : "";

    if (node.kind === "rule") return `${name} validation rule${deprecated}`;
    if (node.kind === "type") return `${name} type${deprecated}`;
    if (node.kind === "proc") return `${name} procedure${deprecated}`;
    if (node.kind === "doc") return `${name} documentation${deprecated}`;
    return "Unknown";
  });

  let contentId = $derived.by(() => {
    if (node.kind === "rule") return slugify(`rule-${name}`);
    if (node.kind === "type") return slugify(`type-${name}`);
    if (node.kind === "proc") return slugify(`proc-${name}`);
    if (node.kind === "doc") return slugify(`doc-${name}`);
    return "";
  });

  let id = $derived(`navlink-${contentId}`);

  let href = $derived(`#${contentId}`);

  let isDeprecated = $derived.by(() => {
    if (node.kind === "doc") return false;
    if (typeof node.deprecated === "string") return true;
    return false;
  });

  const navigate = (contentId: string) => {
    uiStore.activeSection = contentId;
  };

  let isActive = $derived.by(() => {
    return uiStore.activeSection === contentId;
  });
</script>

<Tooltip content={title}>
  <a
    {id}
    {href}
    onclick={() => navigate(contentId)}
    class={[
      "btn btn-ghost btn-block justify-start space-x-2 border-transparent",
      {
        "hover:bg-blue-500/20": node.kind === "doc",
        "hover:bg-yellow-500/20": node.kind === "rule",
        "hover:bg-purple-500/20": node.kind === "type",
        "hover:bg-green-500/20": node.kind === "proc",
        "bg-blue-500/20": isActive && node.kind === "doc",
        "bg-yellow-500/20": isActive && node.kind === "rule",
        "bg-purple-500/20": isActive && node.kind === "type",
        "bg-green-500/20": isActive && node.kind === "proc",
      },
    ]}
  >
    {#if node.kind === "doc"}
      <BookOpenText class="size-4 flex-none" />
    {/if}
    {#if node.kind === "rule"}
      <Scale class="size-4 flex-none" />
    {/if}
    {#if node.kind === "type"}
      <Type class="size-4 flex-none" />
    {/if}
    {#if node.kind === "proc"}
      <ArrowLeftRight class="size-4 flex-none" />
    {/if}

    <span
      class={[
        "overflow-hidden overflow-ellipsis whitespace-nowrap",
        {
          "line-through": isDeprecated,
        },
      ]}
    >
      {name}
    </span>

    {#if isDeprecated}
      <TriangleAlert class="text-warning size-4 flex-none" />
    {/if}
  </a>
</Tooltip>
