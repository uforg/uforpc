<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { isscrolledAction } from "$lib/actions/isScrolled.svelte.ts";
  import { activesectionAction } from "$lib/actions/activeSection.svelte.ts";
  import { store } from "$lib/store.svelte";
  import Aside from "./components/Aside.svelte";
  import Header from "./components/Header.svelte";
  import Main from "./components/Main.svelte";

  let isScrolled = $state(false);
  let activeSection = $state("");

  // Update URL hash when active section changes
  $effect(() => {
    if (activeSection) {
      store.activeSection = activeSection;
      // history.replaceState(null, "", `#${activeSection}`);
      goto(`#${activeSection}`, { noScroll: true });
    }
  });

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

<div class="w-[100dvw] h-[100dvh] flex justify-start">
  <Aside />
  <div
    use:isscrolledAction
    use:activesectionAction
    onisscrolled={(e) => (isScrolled = e.detail)}
    onactivesection={(e) => (activeSection = e.detail)}
    class="flex-grow h-[100dvh] overflow-x-hidden overflow-y-auto scroll-p-[90px]"
  >
    <Header {isScrolled} />
    <Main />

    <div class="h-[calc(100dvh-100px)]"></div>
  </div>
</div>
