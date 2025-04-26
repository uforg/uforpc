<script lang="ts">
  import { isscrolledAction } from "$lib/actions/isScrolled.svelte";
  import Aside from "./components/Aside.svelte";
  import Header from "./components/Header.svelte";

  let isScrolled = $state(false);
</script>

<div class="w-[100dvw] h-[100dvh] flex justify-start">
  <Aside />
  <div
    use:isscrolledAction
    onisscrolled={(e) => (isScrolled = e.detail)}
    class="flex-grow h-[100dvh] overflow-x-hidden overflow-y-auto"
  >
    <Header {isScrolled} />
    <main class="w-full p-4">
      <h1>Home</h1>
    </main>
  </div>
</div>
