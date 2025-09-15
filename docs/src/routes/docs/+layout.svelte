<script lang="ts">
  import H4 from "$lib/components/H4.svelte";

  import type { LayoutProps } from "./$types";

  let { children, data }: LayoutProps = $props();
</script>

<div class="flex h-full w-full">
  <nav class="overflow-auto">
    <ul>
      {#each data.docsInfo as docInfo}
        <li
          class="rounded-lg border border-base-300 bg-base-200 p-4 shadow-sm transition-shadow hover:shadow-md"
        >
          <a href={`/docs/${docInfo.slug}`}>
            <H4 class="mb-2 text-lg font-semibold">{docInfo.meta.title}</H4>
          </a>
        </li>
      {/each}
    </ul>
  </nav>

  <section class="mx-auto w-full overflow-auto">
    {@render children?.()}
  </section>
</div>
