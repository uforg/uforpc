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
</script>

<div>
  <div
    class={{
      "mb-4": true,
      "bg-base-100 sticky top-0 z-20 pt-4": !uiStore.isMobile,
    }}
  >
    <H2 class="mb-4 flex items-center space-x-2">Code snippets</H2>

    <Tabs
      items={[
        { id: "sdk", label: "Client SDK Snippets" },
        { id: "curl", label: "Curl Snippets" },
      ]}
      activeId={uiStore.codeSnippetsTab}
      onSelect={(id) => (uiStore.codeSnippetsTab = id as "sdk" | "curl")}
    />
  </div>

  <div class="space-y-2">
    {#if uiStore.codeSnippetsTab === "sdk"}
      <SnippetsSdk {type} {name} />
    {:else}
      <SnippetsCurl {value} {type} {name} />
    {/if}
  </div>
</div>
