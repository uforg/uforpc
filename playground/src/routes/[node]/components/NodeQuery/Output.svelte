<script lang="ts">
  import { CircleX, CloudAlert, Key, Loader, Sparkles } from "@lucide/svelte";

  import { discoverAuthToken } from "$lib/helpers/discoverAuthToken.ts";
  import { setHeader } from "$lib/store.svelte";

  import Editor from "$lib/components/Editor.svelte";
  import H3 from "$lib/components/H3.svelte";
  import Tooltip from "$lib/components/Tooltip.svelte";

  interface Props {
    type: "stream" | "proc";
    cancelRequest: () => void;
    isExecuting: boolean;
    output: string | null;
  }

  const { cancelRequest, isExecuting, output, type }: Props = $props();
  let hasOutput = $derived(!!output);

  // Discover authentication tokens in the response
  let foundTokens = $derived(type === "proc" ? discoverAuthToken(output) : []);
  let hasToken = $derived(foundTokens.length > 0);

  // Function to handle selecting a specific token
  function handleSelectToken(token: (typeof foundTokens)[number]) {
    setHeader("Authorization", `Bearer ${token.value}`);
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

{#if !hasOutput && !isExecuting}
  <div class="mt-12 flex w-full flex-col items-center justify-center space-y-2">
    <CloudAlert class="size-10" />
    <H3 class="flex items-center justify-start space-x-2">No Output</H3>
  </div>

  <p class="pt-4 text-center">
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
  {#if type == "proc" && hasToken}
    <div class="flex flex-wrap items-center justify-start space-x-2">
      <span class="flex flex-none items-center pr-2 text-sm font-bold">
        <Sparkles class="mr-1 size-4" />
        <span>Quick Actions</span>
      </span>

      {#each foundTokens as token}
        <Tooltip content={`Set "Authorization: Bearer <${token.path}>" header`}>
          <button
            class="btn btn-sm btn-ghost flex-none"
            onclick={() => handleSelectToken(token)}
          >
            <Key class="mr-1 size-4" />
            <span>{token.key}</span>
          </button>
        </Tooltip>
      {/each}
    </div>
  {/if}

  {#if type == "stream" && isExecuting}
    <p>
      The data stream is currently active using Server Sent Events (SSE). You
      can stop it by clicking the button below.
    </p>
    {@render CancelButton()}
  {/if}

  <Editor
    class="rounded-box h-[600px] w-full overflow-hidden shadow-md"
    value={output ?? ""}
  />
{/if}
