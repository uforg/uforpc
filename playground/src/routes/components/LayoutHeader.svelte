<script lang="ts">
  import {
    BookOpenText,
    Ellipsis,
    EllipsisVertical,
    Github,
    Info,
    Menu,
    X,
  } from "@lucide/svelte";

  import { dimensionschangeAction, uiStore } from "$lib/uiStore.svelte";

  import Logo from "$lib/components/Logo.svelte";
  import Offcanvas from "$lib/components/Offcanvas.svelte";
  import Tooltip from "$lib/components/Tooltip.svelte";

  import LayoutHeaderSchema from "./LayoutHeaderSchema.svelte";
  import LayoutHeaderSearch from "./LayoutHeaderSearch.svelte";
  import LayoutHeaderSettings from "./LayoutHeaderSettings.svelte";
  import LayoutHeaderThemeSelect from "./LayoutHeaderThemeSelect.svelte";

  let isOffcanvasOpen = $state(false);
</script>

{#snippet starOnGithub()}
  <a
    href="https://github.com/uforg/uforpc"
    target="_blank"
    class="btn btn-ghost justify-start space-x-1"
  >
    <Github class="size-4" />
    <span>Star on GitHub</span>
    <img
      alt="GitHub Repo stars"
      src="https://img.shields.io/github/stars/uforg/uforpc?style=plastic&label=%20"
    />
  </a>
{/snippet}

{#snippet docsLink()}
  <a
    href="https://uforpc.uforg.dev"
    target="_blank"
    class="btn btn-ghost justify-start space-x-1"
  >
    <BookOpenText class="size-4" />
    <span>Docs</span>
  </a>
{/snippet}

{#snippet offcanvasOpenButton()}
  <Tooltip content="More options" placement="left">
    <button
      class="btn btn-ghost btn-square"
      onclick={() => (isOffcanvasOpen = true)}
    >
      {#if !uiStore.isMobile}
        <EllipsisVertical class="size-4" />
      {/if}
      {#if uiStore.isMobile}
        <Ellipsis class="size-6" />
      {/if}
    </button>
  </Tooltip>
{/snippet}

{#snippet offcanvasMenu()}
  <div class="mt-4 ml-4 flex items-center justify-start space-x-2">
    <button
      class="btn btn-ghost btn-square btn-sm"
      onclick={() => (isOffcanvasOpen = false)}
    >
      <X class="size-6" />
    </button>
    <h2 class="text-lg font-bold">More options</h2>
  </div>
  <div class="flex flex-col items-start p-4 [&>*]:w-full">
    {@render starOnGithub()}
    {@render docsLink()}
    <LayoutHeaderSchema />

    {#if uiStore.isMobile}
      <LayoutHeaderSearch />
      <LayoutHeaderSettings />
      <LayoutHeaderThemeSelect />
    {/if}

    <a href="#/about" class="btn btn-ghost justify-start space-x-1">
      <Info class="size-4" />
      <span>About the project</span>
    </a>
  </div>
{/snippet}

{#if !uiStore.isMobile}
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
      {@render starOnGithub()}
      {@render docsLink()}
      <LayoutHeaderThemeSelect />
      {@render offcanvasOpenButton()}
    </div>
  </header>
{/if}

{#if uiStore.isMobile}
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
      <button
        class="btn btn-ghost btn-square"
        onclick={() => (uiStore.asideOpen = true)}
      >
        <Menu class="size-6" />
      </button>
    </div>

    <Logo class="mx-auto h-[80%]" />

    <div class="flex items-center justify-end space-x-2">
      {@render offcanvasOpenButton()}
    </div>
  </header>
{/if}

<Offcanvas bind:isOpen={isOffcanvasOpen} direction="right">
  {@render offcanvasMenu()}
</Offcanvas>
