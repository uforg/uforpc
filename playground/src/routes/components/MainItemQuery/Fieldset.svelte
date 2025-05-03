<script lang="ts">
  interface Props {
    children: any;
  }

  const { children }: Props = $props();
</script>

<fieldset
  class={[
    "fieldset border border-base-content/20 rounded-box w-full p-4 space-y-2",
  ]}
>
  {@render children()}
</fieldset>
