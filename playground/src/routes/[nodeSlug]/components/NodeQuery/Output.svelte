<script lang="ts">
  import { CircleX, CloudAlert, Copy, Loader, Trash } from "@lucide/svelte";
  import { toast } from "svelte-sonner";

  import Editor from "$lib/components/Editor.svelte";
  import H3 from "$lib/components/H3.svelte";
  import Tooltip from "$lib/components/Tooltip.svelte";

  import type { StoreNodeInstance } from "../../storeNode.svelte";

  import OutputQuickActions from "./OutputQuickActions.svelte";

  interface Props {
    type: "stream" | "proc";
    cancelRequest: () => void;
    isExecuting: boolean;
    storeNode: StoreNodeInstance;
  }

  let { cancelRequest, isExecuting, type, storeNode }: Props = $props();
  let hasOutput = $derived.by(() => {
    if (!storeNode.store.output) return false;
    if (storeNode.store.output === "{}") return false;
    if (storeNode.store.output === "[]") return false;
    return true;
  });

  let prettyOutputDate = $derived.by(() => {
    if (!storeNode.store.outputDate) return "unknown output date";
    let date = new Date(storeNode.store.outputDate);
    return date.toLocaleString();
  });

  async function copyToClipboard() {
    try {
      await navigator.clipboard.writeText(storeNode.store.output);
      toast.success("Output copied to clipboard", { duration: 1500 });
    } catch (err) {
      console.error("Failed to copy output: ", err);
      toast.error("Failed to copy output", {
        description: `Error: ${err}`,
      });
    }
  }
</script>

{#snippet CancelButton()}
  {#if isExecuting}
    <button class="btn btn-error" onclick={cancelRequest}>
      <CircleX class="size-4" />
      {type === "proc" ? "Cancel procedure call" : "Stop stream"}
    </button>
  {/if}
{/snippet}

<div class="w-full space-y-2">
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
    {#if type == "stream" && isExecuting}
      <p class="pb-2 text-sm">
        The data stream is currently active using Server Sent Events (SSE). You
        can stop it by clicking the button below. New messages will be added to
        the top of the output.
      </p>

      <div class="pb-2">
        {@render CancelButton()}
      </div>
    {/if}

    <div class="flex w-full flex-wrap items-start justify-between space-x-2">
      <div>
        {#if type == "proc"}
          <OutputQuickActions output={storeNode.store.output} />
        {/if}
      </div>

      <div class="flex items-center justify-end">
        <Tooltip content="Latest output date" placement="top">
          <button class="btn btn-xs mr-2">
            {prettyOutputDate}
          </button>
        </Tooltip>
        <Tooltip content="Copy output to clipboard" placement="top">
          <button class="btn btn-xs btn-square mr-2" onclick={copyToClipboard}>
            <Copy class="size-3" />
          </button>
        </Tooltip>
        <Tooltip content="Clear output" placement="top">
          <button
            class="btn btn-xs btn-square"
            onclick={storeNode.actions.clearOutput}
          >
            <Trash class="size-3" />
          </button>
        </Tooltip>
      </div>
    </div>

    <Editor
      class="rounded-box h-[600px] w-full overflow-hidden shadow-md"
      lang="json"
      value={storeNode.store.output ?? ""}
    />
  {/if}
</div>
