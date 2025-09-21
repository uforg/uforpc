<script lang="ts">
  import { page } from "$app/state";
  import {
    ArrowLeftRight,
    BookOpenText,
    CornerRightDown,
    TriangleAlert,
    Type,
  } from "@lucide/svelte";

  import { getMarkdownTitle } from "$lib/helpers/getMarkdownTitle";
  import { markSearchHints } from "$lib/helpers/markSearchHints";
  import { slugify } from "$lib/helpers/slugify";
  import type { storeSettings } from "$lib/storeSettings.svelte";
  import { uiStore } from "$lib/uiStore.svelte";

  import Tooltip from "$lib/components/Tooltip.svelte";

  interface Props {
    node: (typeof storeSettings.store.jsonSchema.nodes)[number];
  }

  const { node }: Props = $props();

  let name = $derived.by(() => {
    if (node.kind === "type") return node.name;
    if (node.kind === "proc") return node.name;
    if (node.kind === "stream") return node.name;
    if (node.kind === "doc") return getMarkdownTitle(node.content);
    return "unknown";
  });

  let nameHtml = $derived.by(() => {
    if (!uiStore.store.asideSearchOpen || !uiStore.store.asideSearchQuery) {
      return name;
    }

    return markSearchHints([uiStore.store.asideSearchQuery], name);
  });

  let title = $derived.by(() => {
    const deprecated = isDeprecated ? " (Deprecated)" : "";

    if (node.kind === "type") return `${name} data type${deprecated}`;
    if (node.kind === "proc") return `${name} procedure${deprecated}`;
    if (node.kind === "stream") return `${name} stream${deprecated}`;
    if (node.kind === "doc") return `${name} documentation${deprecated}`;
    return "Unknown";
  });

  let contentId = $derived.by(() => {
    if (node.kind === "type") return slugify(`type-${name}`);
    if (node.kind === "proc") return slugify(`proc-${name}`);
    if (node.kind === "stream") return slugify(`stream-${name}`);
    if (node.kind === "doc") return slugify(`doc-${name}`);
    return "";
  });

  let id = $derived(`navlink-${contentId}`);

  let href = $derived(`#/${contentId}`);

  let isDeprecated = $derived.by(() => {
    if (node.kind === "doc") return false;
    if (typeof node.deprecated === "string") return true;
    return false;
  });

  let isActive = $derived.by(() => {
    const paramsNode = page.params.node;
    if (!paramsNode) return false;

    return paramsNode === contentId;
  });
</script>

<Tooltip content={title}>
  <a
    {id}
    {href}
    onclick={() => (uiStore.store.asideOpen = false)}
    class={[
      "btn btn-ghost btn-block justify-start space-x-2 border-transparent",
      {
        "hover:bg-blue-500/20": node.kind === "doc",
        "hover:bg-purple-500/20": node.kind === "type",
        "hover:bg-green-500/20": node.kind === "proc",
        "hover:bg-yellow-500/20": node.kind === "stream",
        "bg-blue-500/20": isActive && node.kind === "doc",
        "bg-purple-500/20": isActive && node.kind === "type",
        "bg-green-500/20": isActive && node.kind === "proc",
        "bg-yellow-500/20": isActive && node.kind === "stream",
      },
    ]}
  >
    {#if node.kind === "doc"}
      <BookOpenText class="size-4 flex-none" />
    {/if}
    {#if node.kind === "type"}
      <Type class="size-4 flex-none" />
    {/if}
    {#if node.kind === "proc"}
      <ArrowLeftRight class="size-4 flex-none" />
    {/if}
    {#if node.kind === "stream"}
      <CornerRightDown class="size-4 flex-none" />
    {/if}

    <span
      class={[
        "overflow-hidden overflow-ellipsis whitespace-nowrap",
        {
          "line-through": isDeprecated,
        },
      ]}
    >
      {@html nameHtml}
    </span>

    {#if isDeprecated}
      <TriangleAlert class="text-warning size-4 flex-none" />
    {/if}
  </a>
</Tooltip>
