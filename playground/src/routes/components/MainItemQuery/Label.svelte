<script lang="ts">
  import { CircleHelp, SquareAsterisk } from "@lucide/svelte";
  import type { ClassValue } from "$lib/helpers/mergeClasses";
  import { mergeClasses } from "$lib/helpers/mergeClasses";
  import { prettyLabel } from "./prettyLabel";

  interface Props {
    label: string;
    optional: boolean;
    class?: ClassValue;
  }

  const { label, optional, class: className }: Props = $props();

  let plabel = $derived(prettyLabel(label));
  let dataTip = $derived(
    optional ? `${plabel} is optional` : `${plabel} is required`,
  );
</script>

<span
  class={mergeClasses([
    "inline-flex justify-start items-center space-x-1 tooltip tooltip-right",
    className,
  ])}
  data-tip={dataTip}
>
  <span>
    {plabel}
  </span>

  {#if optional}
    <CircleHelp class="size-4 text-info" />
  {:else}
    <SquareAsterisk class="size-4 text-error" />
  {/if}
</span>
