<script lang="ts">
  import { uiStore } from "$lib/uiStore.svelte";

  import H2 from "$lib/components/H2.svelte";
  import Tabs from "$lib/components/Tabs.svelte";

  import SnippetsCurl from "./SnippetsCurl.svelte";
  import SnippetsSdk from "./SnippetsSdk.svelte";

  interface Props {
    // biome-ignore lint/suspicious/noExplicitAny: consistent with sibling components
    value: any;
    type: "proc" | "stream";
    name: string;
  }

  const { value, type, name }: Props = $props();

  let activeTab: "sdk" | "curl" = $state("sdk");

  let maxHeight = $derived.by(() => {
    if (uiStore.isMobile) return "100%";

    const appHeight = uiStore.app.size.offsetHeight;
    const headerHeight = uiStore.header.size.offsetHeight;
    const padding = 16 * 2;

    const mh = appHeight - headerHeight - padding;
    return `${mh}px`;
  });
</script>

<div
  class={{
    "flex h-full flex-col": !uiStore.isMobile,
  }}
  style="max-height: {maxHeight}"
>
  <H2 class="mb-4 flex items-center space-x-2">Code snippets</H2>

  <div class="mb-4">
    <Tabs
      items={[
        { id: "sdk", label: "Client SDK Snippets" },
        { id: "curl", label: "Curl Snippets" },
      ]}
      activeId={activeTab}
      onSelect={(id) => (activeTab = id as "sdk" | "curl")}
    />
  </div>

  <div class="space-y-2 overflow-y-auto">
    {#if activeTab === "sdk"}
      <SnippetsSdk {value} {type} {name} />
    {:else}
      <SnippetsCurl {value} {type} {name} />
    {/if}
  </div>
</div>
