<script lang="ts">
  import Aside from "./components/Aside.svelte";
  import Header from "./components/Header.svelte";
</script>

<div class="w-[100dvw] h-[100dvh] flex justify-start">
  <Aside />
  <div class="flex-grow h-[100dvh] overflow-x-hidden overflow-y-auto">
    <Header />
    <main class="w-full p-4">
      <h1>Home</h1>
    </main>
  </div>
</div>
