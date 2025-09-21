<script>
  import { Info } from "@lucide/svelte";

  import { storeUi } from "$lib/storeUi.svelte";

  import LayoutHeaderDocsLink from "./LayoutHeaderDocsLink.svelte";
  import LayoutHeaderOpenApiSchema from "./LayoutHeaderOpenApiSchema.svelte";
  import LayoutHeaderSchema from "./LayoutHeaderSchema.svelte";
  import LayoutHeaderSearch from "./LayoutHeaderSearch.svelte";
  import LayoutHeaderSettings from "./LayoutHeaderSettings.svelte";
  import LayoutHeaderStarOnGithub from "./LayoutHeaderStarOnGithub.svelte";
  import LayoutHeaderThemeSelect from "./LayoutHeaderThemeSelect.svelte";
</script>

<div class="flex flex-col items-start [&>*]:w-full">
  {#if storeUi.store.isMobile}
    <LayoutHeaderStarOnGithub />
    <LayoutHeaderDocsLink />
  {/if}

  <LayoutHeaderSchema />
  <LayoutHeaderOpenApiSchema />

  {#if storeUi.store.isMobile}
    <LayoutHeaderSearch />
    <LayoutHeaderSettings />
    <LayoutHeaderThemeSelect />
  {/if}

  <a href="#/about" class="btn btn-ghost justify-start space-x-1">
    <Info class="size-4" />
    <span>About the project</span>
  </a>
</div>
