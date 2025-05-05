<script lang="ts">
  import { onMount } from "svelte";
  import { dimensionschangeAction, uiStore } from "$lib/uiStore.svelte";
  import { activesectionAction } from "./activesectionAction.svelte";
  import Aside from "./components/Aside.svelte";
  import Header from "./components/Header.svelte";
  import Main from "./components/Main.svelte";

  // Scroll to hash on initial load
  onMount(() => {
    setTimeout(() => {
      const hash = globalThis.location.hash.slice(1);
      if (hash) {
        const element = document.getElementById(hash);
        if (element) {
          element.scrollIntoView({ behavior: "smooth" });
        }
      }
    }, 500);
  });
</script>

<div
  use:dimensionschangeAction
  ondimensionschange={(e) => uiStore.app = e.detail}
  class="w-[100dvw] h-[100dvh] flex justify-start"
>
  <Aside />
  <div
    use:activesectionAction
    use:dimensionschangeAction
    onactivesection={(e) => (uiStore.activeSection = e.detail)}
    ondimensionschange={(e) => uiStore.contentWrapper = e.detail}
    class="flex-grow h-[100dvh] overflow-x-hidden overflow-y-auto scroll-p-[90px]"
  >
    <Header />
    <Main />
  </div>
</div>
