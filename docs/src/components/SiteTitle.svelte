<!-- 
  This component overrides the default Starlight SiteTitle component to use a custom logo.
  Default: https://github.com/withastro/starlight/blob/main/packages/starlight/components/SiteTitle.astro
-->

<script lang="ts">
  import Logo from "./Logo.svelte";
</script>

<Logo class="custom-logo" />

<style>
  :global(.custom-logo) {
    height: calc(var(--sl-nav-height) - 2 * var(--sl-nav-pad-y));
    width: auto;
    max-width: 100%;
    object-fit: contain;
    object-position: 0 50%;
  }
</style>
