<script>
  import { Info } from "@lucide/svelte";

  import { uiStore } from "$lib/uiStore.svelte";

  import LayoutHeaderDocsLink from "./LayoutHeaderDocsLink.svelte";
  import LayoutHeaderSchema from "./LayoutHeaderSchema.svelte";
  import LayoutHeaderSearch from "./LayoutHeaderSearch.svelte";
  import LayoutHeaderSettings from "./LayoutHeaderSettings.svelte";
  import LayoutHeaderStarOnGithub from "./LayoutHeaderStarOnGithub.svelte";
  import LayoutHeaderThemeSelect from "./LayoutHeaderThemeSelect.svelte";
</script>

<div class="flex flex-col items-start [&>*]:w-full">
  {#if uiStore.isMobile}
    <LayoutHeaderStarOnGithub />
    <LayoutHeaderDocsLink />
  {/if}

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
