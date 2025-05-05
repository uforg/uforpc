<script lang="ts">
  import Menu from "$lib/components/Menu.svelte";
  import Tooltip from "$lib/components/Tooltip.svelte";
  import { uiStore } from "$lib/uiStore.svelte";
  import type { Theme } from "$lib/uiStore.svelte";
  import { Moon, Palette, Sun, SunMoon } from "@lucide/svelte";

  const themesArr: Theme[] = ["system", "light", "dark"];

  function setTheme(theme: Theme) {
    uiStore.theme = theme;
    (document.activeElement as HTMLElement)?.blur();
  }
</script>

{#snippet themeName(showIcon: boolean, tname: Theme)}
  {#if showIcon && tname === "system"}
    <SunMoon class="size-4" />
  {/if}
  {#if showIcon && tname === "light"}
    <Sun class="size-4" />
  {/if}
  {#if showIcon && tname === "dark"}
    <Moon class="size-4" />
  {/if}

  {#if tname === "system"}
    System
  {/if}
  {#if tname === "light"}
    Light
  {/if}
  {#if tname === "dark"}
    Dark
  {/if}
{/snippet}

{#snippet content()}
  <div class="py-1 space-y-2">
    {#each themesArr as themeItem}
      <button
        class="flex justify-start items-center space-x-2 btn btn-ghost btn-block"
        onclick={() => setTheme(themeItem)}
      >
        {@render themeName(true, themeItem)}
      </button>
    {/each}
  </div>
{/snippet}

<Menu {content}>
  <div>
    <Tooltip content="Theme" placement="left">
      <button class="btn btn-ghost">
        <Palette class="size-4" />
        {@render themeName(false, uiStore.theme)}
      </button>
    </Tooltip>
  </div>
</Menu>
