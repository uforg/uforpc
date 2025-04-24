<script lang="ts">
  import Header from "./components/Header.svelte";
</script>

<div class="w-[100dvw] h-[100dvh] flex justify-start">
  <aside
    class="
      flex-none w-full max-w-[280px] h-[100dvh] overflow-x-hidden overflow-y-auto
      p-4
    "
  >
    <a
      class="flex space-x-2 items-center whitespace-nowrap"
      href="https://uforpc.uforg.dev"
      target="_blank"
    >
      <img src="/assets/logo.png" alt="UFO RPC Logo" class="h-[38px]">
      <h1 class="font-bold">UFO RPC Playground</h1>
    </a>
  </aside>

  <div class="flex-grow h-[100dvh] overflow-x-hidden overflow-y-auto">
    <Header />
    <main class="w-full p-4">
      <h1>Home</h1>
    </main>
  </div>
</div>
