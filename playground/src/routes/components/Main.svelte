<script lang="ts">
  import { ScrollText } from "@lucide/svelte";
  import { store } from "$lib/store.svelte";
  import { dimensionschangeAction, uiStore } from "$lib/uistore.svelte";
  import MainItem from "./MainItem.svelte";

  let isEmpty = $derived.by(() => {
    return store.jsonSchema.nodes.length === 0;
  });
</script>

<main
  class="w-full p-4 space-y-[80px]"
  use:dimensionschangeAction
  ondimensionschange={(e) => uiStore.main = e.detail}
>
  {#if isEmpty}
    <div class="mt-[200px] flex flex-col justify-center items-center gap-4">
      <ScrollText class="size-[100px]" />
      <h1 class="text-3xl font-bold">
        Add a schema with some content to display
      </h1>
    </div>
  {/if}

  {#each store.jsonSchema.nodes as node}
    <MainItem {node} />
  {/each}
</main>
