<script lang="ts">
  import "../app.css";

  let { children } = $props();
</script>

<main class="h-[100dvh] w-[100dvw] overflow-auto">
  {@render children?.()}
</main>
