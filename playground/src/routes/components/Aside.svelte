<script lang="ts">
  import AsideSchemaManager from "./AsideSchemaManager.svelte";
</script>

<aside
  class={[
    "flex-none w-full max-w-[280px] h-[100dvh] overflow-x-hidden overflow-y-auto",
  ]}
>
  <a
    class={[
      "flex space-x-2 items-center whitespace-nowrap",
      "h-[72px] w-full sticky top-0 p-4",
    ]}
    href="https://uforpc.uforg.dev"
    target="_blank"
  >
    <img src="/assets/logo.png" alt="UFO RPC Logo" class="h-full">
    <h1 class="font-bold">UFO RPC Playground</h1>
  </a>
  <nav class="p-4">
    <AsideSchemaManager />
  </nav>
</aside>
