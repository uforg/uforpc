<script lang="ts">
  import { Moon, Sun } from "@lucide/svelte";

  import { storeUi } from "$lib/storeUi.svelte";

  function toggleTheme() {
    const newTheme = storeUi.store.theme === "dark" ? "light" : "dark";
    storeUi.store.theme = newTheme;
  }
</script>

<button class="btn btn-ghost justify-start space-x-1" onclick={toggleTheme}>
  {#if storeUi.store.theme === "light"}
    <Sun class="size-4" />
  {/if}
  {#if storeUi.store.theme === "dark"}
    <Moon class="size-4" />
  {/if}

  <span class="w-[5ch] text-left">
    {#if storeUi.store.theme === "light"}
      Light
    {/if}
    {#if storeUi.store.theme === "dark"}
      Dark
    {/if}
  </span>
</button>
