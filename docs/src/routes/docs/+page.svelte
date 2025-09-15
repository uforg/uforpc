<script lang="ts">
  import type { PageProps } from "./$types";

  let { data }: PageProps = $props();
</script>

<h1 class="mb-6 text-2xl font-bold">Docs</h1>

<div class="md:grid-cols-2 lg:grid-cols-3 grid gap-4">
  {#each data.docsInfo as docInfo}
    <div
      class="rounded-lg border border-base-300 bg-base-200 p-4 shadow-sm transition-shadow hover:shadow-md"
    >
      <h2 class="mb-2 text-lg font-semibold">{docInfo.meta.title}</h2>
      <p class="mb-3">{docInfo.meta.description}</p>
      <a
        href={`/docs/${docInfo.slug}`}
        class="text-blue-600 underline hover:text-blue-800">Read more</a
      >
    </div>
  {/each}
</div>
