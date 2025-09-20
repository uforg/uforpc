<script lang="ts">
  import Editor from "$lib/components/Editor.svelte";

  interface Props {
    value: Record<string, any>;
  }

  let { value = $bindable() }: Props = $props();

  let valueString = $state(JSON.stringify(value, null, 2));
  $effect(() => {
    if (JSON.stringify(value, null, 2) !== valueString) {
      try {
        value = JSON.parse(valueString);
      } catch {
        // Ignore JSON parse errors
      }
    }
  });

  let tab: "form" | "json" = $state("form");
</script>

<Editor
  class="rounded-box h-[600px] w-full overflow-hidden shadow-md"
  lang="json"
  bind:value={valueString}
/>
