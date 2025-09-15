<script lang="ts">
  import Logo from "$lib/components/Logo.svelte";

  import "../app.css";

  let { children } = $props();
</script>

<div class="flex h-[100dvh] w-[100dvw] flex-col overflow-hidden">
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

  <main class="h-full w-full flex-1 overflow-auto">
    {@render children?.()}
  </main>
</div>
