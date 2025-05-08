<script lang="ts">
  import { dimensionschangeAction, uiStore } from "$lib/uiStore.svelte";

  import Aside from "./components/Aside.svelte";
  import Header from "./components/Header.svelte";
  import Main from "./components/Main.svelte";
</script>

<div
  use:dimensionschangeAction
  ondimensionschange={(e) => (uiStore.app = e.detail)}
  class="flex h-[100dvh] w-[100dvw] justify-start"
>
  <Aside />
  <div
    use:dimensionschangeAction
    ondimensionschange={(e) => (uiStore.contentWrapper = e.detail)}
    class="h-[100dvh] flex-grow scroll-p-[90px] overflow-x-hidden overflow-y-auto"
  >
    <Header />
    <Main />
  </div>
</div>
