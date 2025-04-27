<script lang="ts">
  import { mergeClasses } from "$lib/helpers/mergeClasses";
  import type { ClassValue } from "$lib/helpers/mergeClasses";

  interface Props {
    children: any;
    class?: ClassValue;
    rest?: any;
  }

  let { children, class: className, ...rest }: Props = $props();
</script>

<h2
  class={mergeClasses(
    "text-2xl font-bold",
    className,
  )}
  {...rest}
>
  {@render children()}
</h2>
