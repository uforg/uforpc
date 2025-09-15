<script lang="ts">
  import Logo from "$lib/components/Logo.svelte";

  import "../app.css";

  let { children } = $props();
</script>

<header
  class="flex w-full items-center justify-between bg-base-300 p-4 desk:px-12"
>
  <a href="/">
    <Logo class="h-[30px] w-auto" />
  </a>
  <nav>
    <a href="/docs" class="link">Docs</a>
  </nav>
</header>

<main class="h-[100dvh] w-[100dvw] overflow-auto">
  {@render children?.()}
</main>
