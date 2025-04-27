<script lang="ts">
  import { isscrolledAction } from "$lib/actions/isScrolled.svelte";
  import { onMount } from "svelte";
  import Aside from "./components/Aside.svelte";
  import Header from "./components/Header.svelte";
  import Main from "./components/Main.svelte";

  let isScrolled = $state(false);

  // if has hash anchor navigate to it
  onMount(async () => {
    // wait 500ms to ensure the content is rendered
    await new Promise((resolve) => setTimeout(resolve, 500));

    if (window.location.hash) {
      const element = document.getElementById(
        window.location.hash.slice(1),
      );
      if (element) {
        element.scrollIntoView({ behavior: "smooth" });
      }
    }
  });
</script>

<div class="w-[100dvw] h-[100dvh] flex justify-start">
  <Aside />
  <div
    use:isscrolledAction
    onisscrolled={(e) => (isScrolled = e.detail)}
    class="flex-grow h-[100dvh] overflow-x-hidden overflow-y-auto scroll-p-[90px]"
  >
    <Header {isScrolled} />
    <Main />
  </div>
</div>
