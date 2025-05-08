<script lang="ts">
  import { BookOpenText, Github } from "@lucide/svelte";

  import { dimensionschangeAction, uiStore } from "$lib/uiStore.svelte";

  import LayoutHeaderSearch from "./LayoutHeaderSearch.svelte";
  import LayoutHeaderSettings from "./LayoutHeaderSettings.svelte";
  import LayoutHeaderThemeSelect from "./LayoutHeaderThemeSelect.svelte";
</script>

<header
  use:dimensionschangeAction
  ondimensionschange={(e) => (uiStore.header = e.detail)}
  class={[
    "sticky top-0 z-30 flex h-[72px] w-full items-center justify-between space-x-2 p-4",
    "bg-base-100/90 backdrop-blur-sm",
    {
      "shadow-xs": uiStore.contentWrapper.scroll.isTopScrolled,
    },
  ]}
>
  <div class="flex items-center justify-start space-x-2">
    <LayoutHeaderSearch />
    <LayoutHeaderSettings />
  </div>
  <div class="flex items-center justify-end space-x-2">
    <a
      href="https://github.com/uforg/uforpc"
      target="_blank"
      class="btn btn-ghost"
    >
      <Github class="size-4" />
      <span>Star on GitHub</span>
      <img
        alt="GitHub Repo stars"
        src="https://img.shields.io/github/stars/uforg/uforpc?style=plastic&label=%20"
      />
    </a>
    <a href="https://uforpc.uforg.dev" target="_blank" class="btn btn-ghost">
      <BookOpenText class="size-4" /> Docs
    </a>
    <LayoutHeaderThemeSelect />
  </div>
</header>
