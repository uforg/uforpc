<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    children: Snippet;
  }

  const { children }: Props = $props();
</script>

<fieldset
  class={[
    "fieldset-component",
    "fieldset border-base-content/20 rounded-box w-full space-y-2 border p-4",
    "hover:border-base-content/50 has-[.fieldset-component:hover]:hover:border-base-content/20",
  ]}
>
  {@render children()}
</fieldset>
