<script lang="ts">
  interface Props {
    children: any;
  }

  const { children }: Props = $props();
</script>

<fieldset
  class={[
    "fieldset-component",
    "fieldset border border-base-content/20 rounded-box w-full p-4 space-y-2",
    "hover:border-base-content/50 has-[.fieldset-component:hover]:hover:border-base-content/20",
  ]}
>
  {@render children()}
</fieldset>
