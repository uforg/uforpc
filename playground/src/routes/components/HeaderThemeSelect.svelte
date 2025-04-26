<script lang="ts">
  import { setTheme, store } from "$lib/store.svelte";
  import type { Theme } from "$lib/store.svelte";
  import { Moon, MoonStar, Palette, Sun, SunMoon } from "@lucide/svelte";

  const themesArr: Theme[] = ["system", "corporate", "dark", "dracula"];

  function setThemeAndBlur(theme: Theme) {
    setTheme(theme);
    (document.activeElement as HTMLElement)?.blur();
  }
</script>

{#snippet themeName(showIcon: boolean, tname: Theme)}
  {#if showIcon && tname === "system"}
    <SunMoon class="size-4" />
  {/if}
  {#if showIcon && tname === "corporate"}
    <Sun class="size-4" />
  {/if}
  {#if showIcon && tname === "dark"}
    <Moon class="size-4" />
  {/if}
  {#if showIcon && tname === "dracula"}
    <MoonStar class="size-4" />
  {/if}

  {#if tname === "system"}
    System
  {/if}
  {#if tname === "corporate"}
    Light
  {/if}
  {#if tname === "dark"}
    Dark
  {/if}
  {#if tname === "dracula"}
    Dracula
  {/if}
{/snippet}

<div class="dropdown dropdown-end">
  <div
    tabindex="-1"
    role="button"
    class="btn btn-ghost tooltip tooltip-left"
    data-tip="Theme"
  >
    <Palette class="size-4" />
    {@render themeName(false, store.theme)}
  </div>
  <ul
    tabindex="-1"
    class="dropdown-content menu bg-base-100 rounded-box z-1 w-[120px] p-2 shadow-md"
  >
    {#each themesArr as themeItem}
      <li>
        <button onclick={() => setThemeAndBlur(themeItem)}>
          {@render themeName(true, themeItem)}
        </button>
      </li>
    {/each}
  </ul>
</div>
