<script lang="ts">
  import ThemeSelect from "$lib/components/ThemeSelect.svelte";
  import { BookOpenText, Github } from "@lucide/svelte";
  import Search from "./Search.svelte";
</script>

<header class="w-full p-4 flex justify-between items-center space-x-2">
  <div class="flex-grow">
    <Search />
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
    <ThemeSelect />
  </div>
</header>
