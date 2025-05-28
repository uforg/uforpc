<script lang="ts">
  import { Moon, Sun } from "@lucide/svelte";
  import { fade } from "svelte/transition";

  import { uiStore } from "$lib/uiStore.svelte";

  import Tooltip from "$lib/components/Tooltip.svelte";

  function toggleTheme() {
    const newTheme = uiStore.theme === "dark" ? "light" : "dark";
    uiStore.theme = newTheme;
  }

  let tooltipContent = $derived(
    uiStore.theme === "dark" ? "Set light theme" : "Set dark theme",
  );
</script>

<Tooltip content={tooltipContent} placement="left">
  <button class="btn btn-ghost" onclick={toggleTheme}>
    <span class="relative size-4">
      {#if uiStore.theme === "light"}
        <span transition:fade={{ duration: 100 }} class="absolute inset-0">
          <Sun class="size-4" />
        </span>
      {/if}
      {#if uiStore.theme === "dark"}
        <span transition:fade={{ duration: 100 }} class="absolute inset-0">
          <Moon class="size-4" />
        </span>
      {/if}
    </span>

    <span class="w-[5ch]">
      {#if uiStore.theme === "light"}
        Light
      {/if}
      {#if uiStore.theme === "dark"}
        Dark
      {/if}
    </span>
  </button>
</Tooltip>
