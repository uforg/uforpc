<script lang="ts">
  import { CloudAlert, Key, Sparkles } from "@lucide/svelte";

  import { discoverAuthToken } from "$lib/helpers/discoverAuthToken.ts";
  import { setHeader } from "$lib/store.svelte";

  import Editor from "$lib/components/Editor.svelte";
  import H3 from "$lib/components/H3.svelte";
  import Tooltip from "$lib/components/Tooltip.svelte";

  interface Props {
    output: object | string | null;
  }

  const { output }: Props = $props();

  let hasOutput = $derived(!!output);
  let outputString = $derived(
    typeof output === "string" ? output : JSON.stringify(output, null, 2),
  );

  // Discover authentication tokens in the response
  let foundTokens = $derived(discoverAuthToken(output));
  let hasToken = $derived(foundTokens.length > 0);

  // Function to handle selecting a specific token
  function handleSelectToken(token: (typeof foundTokens)[number]) {
    setHeader("Authorization", `Bearer ${token.value}`);
  }
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
  {#if hasToken}
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

  <Editor
    class="rounded-box h-[600px] w-full overflow-hidden shadow-md"
    value={outputString}
  />
{/if}
