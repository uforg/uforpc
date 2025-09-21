<script lang="ts">
  import { Code, Terminal } from "@lucide/svelte";

  import { storeUi } from "$lib/storeUi.svelte";

  import H2 from "$lib/components/H2.svelte";
  import Tabs from "$lib/components/Tabs.svelte";

  import SnippetsCurl from "./SnippetsCurl.svelte";
  import SnippetsSdk from "./SnippetsSdk.svelte";

  interface Props {
    input: any;
    type: "proc" | "stream";
    name: string;
  }

  let { input, type, name }: Props = $props();
</script>

<div>
  <div
    class={{
      "mb-4": true,
      "bg-base-100 sticky top-0 z-20 pt-4": !storeUi.store.isMobile,
    }}
  >
    <H2 class="mb-4 flex items-center space-x-2">Code snippets</H2>

    <Tabs
      items={[
        { id: "curl", label: "HTTP Snippets", icon: Terminal },
        { id: "sdk", label: "SDK Snippets", icon: Code },
      ]}
      bind:active={storeUi.store.codeSnippetsTab}
    />
  </div>

  <div class="space-y-2">
    {#if storeUi.store.codeSnippetsTab === "sdk"}
      <SnippetsSdk {type} {name} />
    {:else}
      <SnippetsCurl {input} {type} {name} />
    {/if}
  </div>
</div>
