<script lang="ts">
  import { untrack } from "svelte";

  import Editor from "$lib/components/Editor.svelte";

  interface Props {
    input: Record<string, any>;
  }

  let { input = $bindable() }: Props = $props();

  let valueString = $state(JSON.stringify(input, null, 2));
  $effect(() => {
    const val = valueString;
    untrack(() => {
      try {
        input = JSON.parse(val);
      } catch {
        // Empty object if invalid JSON until the user fixes it
        input = {};
      }
    });
  });
</script>

<Editor
  class="rounded-box h-[600px] w-full overflow-hidden shadow-md"
  lang="json"
  bind:value={valueString}
/>
