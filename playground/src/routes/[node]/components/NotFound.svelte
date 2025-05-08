<script>
  import { FileWarning } from "@lucide/svelte";

  import H1 from "$lib/components/H1.svelte";
</script>

<H1 class="flex items-center justify-start space-x-2">
  <FileWarning class="size-8" />
  <span>Not Found</span>
</H1>

<p class="pt-4">
  What you're looking for doesn't exist in the current schema or has been
  renamed. Use the menu on the left to navigate between available pages.
</p>
