<script lang="ts">
  import { markdownToHtml } from "$lib/helpers/markdownToHtml";
  import { mergeClasses } from "$lib/helpers/mergeClasses";

  interface Props {
    doc?: string;
    class?: string;
  }

  let { doc, class: className }: Props = $props();

  let docHtml = $state("");

  async function renderDoc() {
    if (!doc) return;
    docHtml = await markdownToHtml(doc);
  }

  $effect(() => {
    renderDoc();
  });
</script>

{#if docHtml !== ""}
  <div
    class={mergeClasses([
      "prose prose-sm text-base-content/50 max-w-none",
      className,
    ])}
  >
    {@html docHtml}
  </div>
{/if}
