<script lang="ts">
  import { CircleHelp, SquareAsterisk } from "@lucide/svelte";
  import type { ClassValue } from "$lib/helpers/mergeClasses";
  import { mergeClasses } from "$lib/helpers/mergeClasses";
  import Tooltip from "$lib/components/Tooltip.svelte";
  import { prettyLabel } from "./prettyLabel";

  interface Props {
    label: string;
    optional: boolean;
    class?: ClassValue;
  }

  const { label, optional, class: className }: Props = $props();

  let plabel = $derived(prettyLabel(label));
  let dataTip = $derived(
    `${plabel} is marked as ${optional ? "optional" : "required"}`,
  );
</script>

<Tooltip content={dataTip} placement="right">
  <span
    class={mergeClasses([
      "inline-flex items-center justify-start space-x-1",
      className,
    ])}
  >
    <span>
      {plabel}
    </span>

    {#if optional}
      <CircleHelp class="text-info size-4" />
    {:else}
      <SquareAsterisk class="text-error size-4" />
    {/if}
  </span>
</Tooltip>
