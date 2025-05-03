<script lang="ts">
  import { BookOpenText, Github } from "@lucide/svelte";
  import { dimensionschangeAction, uiStore } from "$lib/uiStore.svelte";
  import HeaderSearch from "./HeaderSearch.svelte";
  import HeaderSettings from "./HeaderSettings.svelte";
  import HeaderThemeSelect from "./HeaderThemeSelect.svelte";
</script>

<header
  use:dimensionschangeAction
  ondimensionschange={(e) => uiStore.header = e.detail}
  class={[
    "w-full p-4 flex justify-between items-center space-x-2 h-[72px] sticky top-0 z-30",
    "bg-base-100/90 backdrop-blur-sm",
    {
      "shadow-xs": uiStore.contentWrapper.scroll.isTopScrolled,
    },
  ]}
>
  <div class="flex justify-start items-center space-x-2">
    <HeaderSearch />
    <HeaderSettings />
  </div>
  <div class="flex justify-end items-center space-x-2">
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
    <a
      href="https://uforpc.uforg.dev"
      target="_blank"
      class="btn btn-ghost"
    >
      <BookOpenText class="size-4" /> Docs
    </a>
    <HeaderThemeSelect />
  </div>
</header>
