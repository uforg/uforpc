<script lang="ts">
  import type { PageProps } from "./$types";

  let { data }: PageProps = $props();
</script>

<svelte:head>
  <title>{data.docInfo.meta.title}</title>
  <meta property="og:type" content="article" />
  <meta property="og:title" content={data.docInfo.meta.title} />
</svelte:head>

<article class="prose">
  <h1>{data.docInfo.meta.title}</h1>
  <data.content />
</article>
