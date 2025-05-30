<script lang="ts">
  import { CloudAlert } from "@lucide/svelte";

  import Editor from "$lib/components/Editor.svelte";
  import H3 from "$lib/components/H3.svelte";

  interface Props {
    output: object | null;
  }

  const { output }: Props = $props();

  let hasOutput = $derived(!!output);
  let outputString = $derived(JSON.stringify(output, null, 2));
</script>

{#if !hasOutput}
  <div class="mt-12 flex w-full flex-col items-center justify-center space-y-2">
    <CloudAlert class="size-10" />
    <H3 class="flex items-center justify-start space-x-2">No Output</H3>
  </div>

  <p class="pt-4 text-center">
    Please execute the procedure from the input tab to see the output.
  </p>
{/if}

{#if hasOutput}
  <Editor
    class="rounded-box h-[600px] w-full overflow-hidden shadow-md"
    value={outputString}
  />
{/if}
