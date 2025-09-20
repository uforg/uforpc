<script lang="ts">
  import { CircleX, CloudAlert, Loader } from "@lucide/svelte";

  import Editor from "$lib/components/Editor.svelte";
  import H3 from "$lib/components/H3.svelte";

  import OutputQuickActions from "./OutputQuickActions.svelte";

  interface Props {
    type: "stream" | "proc";
    cancelRequest: () => void;
    isExecuting: boolean;
    output: string | null;
  }

  const { cancelRequest, isExecuting, output, type }: Props = $props();
  let hasOutput = $derived.by(() => {
    if (!output) return false;
    if (output === "{}") return false;
    if (output === "[]") return false;
    return true;
  });
</script>

{#snippet CancelButton()}
  {#if isExecuting}
    <button class="btn btn-error" onclick={cancelRequest}>
      <CircleX class="size-4" />
      {type === "proc" ? "Cancel procedure call" : "Stop stream"}
    </button>
  {/if}
{/snippet}

{#if !hasOutput && !isExecuting}
  <div
    class="mt-[100px] flex w-full flex-col items-center justify-center space-y-2"
  >
    <CloudAlert class="size-10" />
    <H3 class="flex items-center justify-start space-x-2">No Output</H3>
  </div>

  <p class="mb-[100px] pt-4 text-center">
    Please execute the {type === "proc" ? "procedure" : "stream"} from the input
    tab to see the output.
  </p>
{/if}

{#if !hasOutput && isExecuting}
  <div
    class="mt-12 mb-4 flex w-full flex-col items-center justify-center space-y-2"
  >
    <Loader class="animate size-10 animate-spin" />
    <H3 class="flex items-center justify-start space-x-2">
      {type === "proc" ? "Executing procedure" : "Starting data stream"}
    </H3>
  </div>

  <div class="flex justify-center">
    {@render CancelButton()}
  </div>
{/if}

{#if hasOutput}
  {#if type == "proc"}
    <OutputQuickActions {output} />
  {/if}

  {#if type == "stream" && isExecuting}
    <p>
      The data stream is currently active using Server Sent Events (SSE). You
      can stop it by clicking the button below. New messages will be added to
      the top of the output.
    </p>
    {@render CancelButton()}
  {/if}

  <Editor
    class="rounded-box h-[600px] w-full overflow-hidden shadow-md"
    lang="json"
    value={output ?? ""}
  />
{/if}
