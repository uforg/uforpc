<script lang="ts">
  import { TriangleAlert } from "@lucide/svelte";
  import { onMount } from "svelte";

  import { deleteMarkdownHeadings } from "$lib/helpers/deleteMarkdownHeadings";
  import { extractNodeFromSchema } from "$lib/helpers/extractNodeFromSchema";
  import { getMarkdownTitle } from "$lib/helpers/getMarkdownTitle";
  import { markdownToHtml } from "$lib/helpers/markdownToHtml";
  import {
    getCurrentInputForOperation,
    getCurrentOutputForOperation,
  } from "$lib/historyStore.svelte";
  import { store } from "$lib/store.svelte";
  import { uiStore } from "$lib/uiStore.svelte";

  import BottomSpace from "$lib/components/BottomSpace.svelte";
  import Code from "$lib/components/Code.svelte";
  import H2 from "$lib/components/H2.svelte";

  import NodeQueryProc from "./NodeQuery/QueryProc.svelte";
  import NodeQueryStream from "./NodeQuery/QueryStream.svelte";
  import Snippets from "./NodeQuery/Snippets/Snippets.svelte";

  interface Props {
    nodeSlug: string;
    node: (typeof store.store.jsonSchema.nodes)[number];
  }

  const { nodeSlug, node }: Props = $props();

  // svelte-ignore non_reactive_update - This is reactive through the store
  let input = getCurrentInputForOperation(nodeSlug);
  // svelte-ignore non_reactive_update - This is reactive through the store
  let output = getCurrentOutputForOperation(nodeSlug);

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
    const extracted = extractNodeFromSchema(
      store.store.urpcSchema,
      node.kind,
      name,
    );
    if (extracted) urpcSchema = extracted;
  });

  let isProcOrStream = $derived(node.kind === "proc" || node.kind === "stream");
</script>

<div
  class={{
    "h-full overflow-hidden": true,
    "grid grid-cols-12": !uiStore.store.isMobile,
  }}
>
  <section
    class={{
      "h-full space-y-12 overflow-y-auto p-4 pt-0": true,
      "col-span-8": !uiStore.store.isMobile && isProcOrStream,
      "col-span-12": !uiStore.store.isMobile && !isProcOrStream,
    }}
  >
    <div
      class={{
        "prose pt-4": true,
        "max-w-none": isProcOrStream,
        "max-w-5xl": !isProcOrStream,
      }}
    >
      <h1 class="break-all">{name}</h1>

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
        <NodeQueryProc proc={node} bind:input bind:output />
      </div>
    {/if}

    {#if node.kind === "stream"}
      <div>
        <NodeQueryStream stream={node} bind:input bind:output />
      </div>
    {/if}

    {#if urpcSchema !== ""}
      <div
        class={{
          "space-y-4": true,
          "max-w-5xl": !isProcOrStream,
        }}
      >
        <H2>Schema</H2>
        <Code lang="urpc" code={urpcSchema} />
      </div>
    {/if}

    {#if uiStore.store.isMobile || !isProcOrStream}
      <BottomSpace />
    {/if}
  </section>

  {#if !uiStore.store.isMobile && (node.kind == "proc" || node.kind == "stream")}
    <div class="col-span-4 overflow-y-auto p-4 pt-0">
      <Snippets {input} type={node.kind} name={node.name} />
      <BottomSpace />
    </div>
  {/if}
</div>
