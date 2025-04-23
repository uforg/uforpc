<script lang="ts">
  import ThemeSelect from "$lib/components/ThemeSelect.svelte";
  import { BookOpenText, Github, Search } from "@lucide/svelte";
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
    <header class="w-full p-4 flex justify-between items-center space-x-2">
      <label class="input">
        <Search />
        <input type="search" class="grow" placeholder="Search..." />
        <kbd class="kbd kbd-sm">CTRL</kbd>
        <kbd class="kbd kbd-sm">K</kbd>
      </label>

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
        <ThemeSelect />
      </div>
    </header>
    <main class="w-full p-4">
      <h1>Home</h1>
    </main>
  </div>
</div>
