<script lang="ts">
  import { Moon, Sun } from "@lucide/svelte";

  import { uiStore } from "$lib/uiStore.svelte";

  function toggleTheme() {
    const newTheme = uiStore.theme === "dark" ? "light" : "dark";
    uiStore.theme = newTheme;
  }
</script>

<button class="btn btn-ghost justify-start space-x-1" onclick={toggleTheme}>
  {#if uiStore.theme === "light"}
    <Sun class="size-4" />
  {/if}
  {#if uiStore.theme === "dark"}
    <Moon class="size-4" />
  {/if}

  <span class="w-[5ch] text-left">
    {#if uiStore.theme === "light"}
      Light
    {/if}
    {#if uiStore.theme === "dark"}
      Dark
    {/if}
  </span>
</button>
