<script lang="ts">
  import { untrack } from "svelte";

  import Editor from "$lib/components/Editor.svelte";

  interface Props {
    value: Record<string, any>;
  }

  let { value = $bindable() }: Props = $props();

  let valueString = $state(JSON.stringify(value, null, 2));
  $effect(() => {
    const val = valueString;
    untrack(() => {
      try {
        value = JSON.parse(val);
      } catch {
        // Empty object if invalid JSON until the user fixes it
        value = {};
      }
    });
  });
</script>

<Editor
  class="rounded-box h-[600px] w-full overflow-hidden shadow-md"
  lang="json"
  bind:value={valueString}
/>
