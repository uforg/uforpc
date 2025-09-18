<!-- 
  This component overrides the default Starlight SiteTitle component to use a custom logo.
  Default: https://github.com/withastro/starlight/blob/main/packages/starlight/components/SiteTitle.astro
-->

<script lang="ts">
  import Logo from "./Logo.svelte";
</script>

<a href="/" class="flex items-center justify-center">
  <Logo class="custom-logo" />
</a>

<style>
  /* Reset default link styles */
  a,
  a:hover,
  a:focus {
    text-decoration: none;
    background: transparent;
    color: inherit;
    border: none;
    outline: none;
  }

  /* Add focus style */
  a:focus-visible {
    box-shadow: 0 0 0 2px var(--sl-color-gray-2);
    border-radius: 2px;
  }

  /* Same size as the original */
  :global(.custom-logo) {
    height: calc(var(--sl-nav-height) - 2 * var(--sl-nav-pad-y));
    width: auto;
    max-width: 100%;
    object-fit: contain;
    object-position: 0 50%;
  }
</style>
