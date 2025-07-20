<script lang="ts">
  import { Ellipsis, EllipsisVertical, Menu, X } from "@lucide/svelte";

  import { dimensionschangeAction, uiStore } from "$lib/uiStore.svelte";

  import H4 from "$lib/components/H3.svelte";
  import Logo from "$lib/components/Logo.svelte";
  import MenuComponent from "$lib/components/Menu.svelte";
  import Offcanvas from "$lib/components/Offcanvas.svelte";

  import LayoutHeaderDocsLink from "./LayoutHeaderDocsLink.svelte";
  import LayoutHeaderMoreOptions from "./LayoutHeaderMoreOptions.svelte";
  import LayoutHeaderSearch from "./LayoutHeaderSearch.svelte";
  import LayoutHeaderSettings from "./LayoutHeaderSettings.svelte";
  import LayoutHeaderStarOnGithub from "./LayoutHeaderStarOnGithub.svelte";
  import LayoutHeaderThemeSelect from "./LayoutHeaderThemeSelect.svelte";

  let isMobileOffcanvasOpen = $state(false);
</script>

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
      <LayoutHeaderStarOnGithub />
      <LayoutHeaderDocsLink />
      <LayoutHeaderThemeSelect />

      {#snippet menuContent()}
        <LayoutHeaderMoreOptions />
      {/snippet}

      <MenuComponent
        content={menuContent}
        trigger="mouseenter focus"
        placement="bottom-start"
      >
        <button class="btn btn-ghost btn-square">
          <EllipsisVertical class="size-4" />
        </button>
      </MenuComponent>
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
      <button
        class="btn btn-ghost btn-square"
        onclick={() => (isMobileOffcanvasOpen = true)}
      >
        <Ellipsis class="size-6" />
      </button>
    </div>
  </header>

  <Offcanvas bind:isOpen={isMobileOffcanvasOpen} direction="right">
    <div class="mt-4 ml-4 flex items-center justify-start space-x-2">
      <button
        class="btn btn-ghost btn-square btn-sm"
        onclick={() => (isMobileOffcanvasOpen = false)}
      >
        <X class="size-6" />
      </button>
      <h2 class="text-lg font-bold">More options</h2>
    </div>
    <div class="p-4">
      <LayoutHeaderMoreOptions />
    </div>
  </Offcanvas>
{/if}
