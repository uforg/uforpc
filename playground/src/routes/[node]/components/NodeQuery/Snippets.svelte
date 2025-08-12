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
</script>

<div>
  <div
    class={{
      "mb-4": true,
      "bg-base-100/90 backdrop-blur-sm": !uiStore.isMobile,
      "sticky top-[72px] z-20 -mt-4 pt-4": !uiStore.isMobile,
    }}
  >
    <H2 class="mb-4 flex items-center space-x-2">Code snippets</H2>

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
