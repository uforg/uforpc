<!--
  This uses the following helper to handle theme changes and is loaded
  in the head of the document to prevent any flash of unstyled content.

  /static/theme-helper.js
-->

<script lang="ts">
  import { Moon, MoonStar, Palette, Sun, SunMoon } from "@lucide/svelte";
  import { onMount } from "svelte";

  const themesArr = ["system", "corporate", "dark", "dracula"];

  let currentTheme = $state("");
  onMount(() => {
    const theme = (window as any).getTheme();
    currentTheme = theme || "system";
  });

  function setTheme(theme: string) {
    currentTheme = theme;
    (window as any).setTheme(currentTheme);
    (document.activeElement as HTMLElement)?.blur();
  }
</script>

{#snippet themeName(showIcon: boolean, tname: string)}
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
    {@render themeName(false, currentTheme)}
  </div>
  <ul
    tabindex="-1"
    class="dropdown-content menu bg-base-100 rounded-box z-1 w-[120px] p-2 shadow-md"
  >
    {#each themesArr as themeItem}
      <li>
        <button onclick={() => setTheme(themeItem)}>
          {@render themeName(true, themeItem)}
        </button>
      </li>
    {/each}
  </ul>
</div>
