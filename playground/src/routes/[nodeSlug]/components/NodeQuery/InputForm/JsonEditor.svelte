<script lang="ts">
  import { onMount, untrack } from "svelte";

  import Editor from "$lib/components/Editor.svelte";

  interface Props {
    input: Record<string, any>;
  }

  let { input = $bindable() }: Props = $props();

  let initialValue = $state("");
  onMount(() => {
    initialValue = JSON.stringify(input, null, 2);
  });

  const handleChange = (newValue: string) => {
    let shouldUpdate = true;
    if (newValue === initialValue) shouldUpdate = false;
    if (newValue.trim() === "") shouldUpdate = false;
    if (!shouldUpdate) return;

    try {
      const parsed = JSON.parse(newValue) as Record<string, any>;

      // Mutate the existing input object in-place to preserve its reactive proxy
      // and any references other components may be holding (historyStore, snippets, etc.)

      // 1) remove keys not present in parsed
      for (const key of Object.keys(input)) {
        if (!(key in parsed)) delete (input as Record<string, any>)[key];
      }

      // 2) assign/overwrite keys from parsed
      for (const [key, value] of Object.entries(parsed)) {
        (input as Record<string, any>)[key] = value;
      }
    } catch {
      // Empty object if invalid JSON until the user fixes it
      for (const key of Object.keys(input)) {
        delete (input as Record<string, any>)[key];
      }
    }
  };
</script>

<Editor
  class="rounded-box h-[600px] w-full overflow-hidden shadow-md"
  lang="json"
  value={initialValue}
  onChange={handleChange}
/>
