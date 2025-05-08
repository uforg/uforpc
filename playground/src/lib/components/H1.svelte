<script lang="ts">
  import type { Snippet } from "svelte";
  import { mergeClasses } from "$lib/helpers/mergeClasses";
  import type { ClassValue } from "$lib/helpers/mergeClasses";

  interface Props {
    children: Snippet;
    class?: ClassValue;
    rest?: any;
  }

  let { children, class: className, ...rest }: Props = $props();
</script>

<h1 class={mergeClasses("text-4xl font-bold", className)} {...rest}>
  {@render children()}
</h1>
