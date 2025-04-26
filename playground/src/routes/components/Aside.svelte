<script lang="ts">
  import { store } from "$lib/store.svelte";
  import { isscrolledAction } from "$lib/actions/isScrolled.svelte";
  import AsideSchemaManager from "./AsideSchemaManager.svelte";
  import AsideItem from "./AsideItem.svelte";

  let isScrolled = $state(false);
</script>

<aside
  use:isscrolledAction
  onisscrolled={(e) => (isScrolled = e.detail)}
  class={[
    "flex-none w-full max-w-[280px] h-[100dvh] overflow-x-hidden overflow-y-auto",
  ]}
>
  <a
    class={[
      "flex space-x-2 items-center whitespace-nowrap",
      "h-[72px] w-full sticky top-0 p-4 z-10",
      "bg-base-100/90 backdrop-blur-sm",
      {
        "shadow-xs": isScrolled,
      },
    ]}
    href="https://uforpc.uforg.dev"
    target="_blank"
  >
    <img src="/assets/logo.png" alt="UFO RPC Logo" class="h-full">
    <h1 class="font-bold">UFO RPC Playground</h1>
  </a>
  <nav class="p-4 space-y-2">
    <AsideSchemaManager />
    {#each store.jsonSchema.nodes as node}
      <AsideItem {node} />
    {/each}
  </nav>
</aside>
