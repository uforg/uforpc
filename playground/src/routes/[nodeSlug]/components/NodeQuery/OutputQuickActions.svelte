<script lang="ts">
  import { Copy, EarthLock, Key, Sparkles } from "@lucide/svelte";
  import { toast } from "svelte-sonner";

  import {
    discoverAuthToken,
    type TokenInfo,
  } from "$lib/helpers/discoverAuthToken.ts";
  import { setHeader } from "$lib/storeSettings.svelte";

  import Menu from "$lib/components/Menu.svelte";
  import Tooltip from "$lib/components/Tooltip.svelte";

  interface Props {
    output: string | null;
  }

  const { output }: Props = $props();

  // Discover authentication tokens in the response, limit to tokenLimit
  // in order to avoid UI overload
  const tokenLimit = 5;
  let foundTokens = $derived(discoverAuthToken(output));
  let firstTokens = $derived(foundTokens.slice(0, tokenLimit));
  let hasToken = $derived(foundTokens.length > 0);

  function handleSetAuthHeader(token: TokenInfo) {
    setHeader("Authorization", `Bearer ${token.value}`);
  }

  async function handleCopyToClipboard(token: TokenInfo) {
    try {
      await navigator.clipboard.writeText(token.value);
      toast.success("Token copied to clipboard", { duration: 1500 });
    } catch (err) {
      console.error("Failed to copy token: ", err);
      toast.error("Failed to copy token", {
        description: `Error: ${err}`,
      });
    }
  }
</script>

{#if hasToken}
  <div class="flex flex-wrap items-center justify-start space-y-2 space-x-2">
    <span class="flex flex-none items-center pr-2 text-xs font-bold">
      <Sparkles class="mr-1 size-3" />
      <span>Quick Actions</span>
    </span>

    {#each firstTokens as token}
      {#snippet menuContent()}
        <div class="flex flex-col space-y-2 pt-2">
          <button
            class="btn btn-sm btn-ghost w-full justify-start"
            onclick={() => handleCopyToClipboard(token)}
          >
            <Copy class="mr-1 size-4" />
            <span>Copy token</span>
          </button>
          <Tooltip content={`Authorization: Bearer <${token.path}>`}>
            <button
              class="btn btn-sm btn-ghost w-full justify-start"
              onclick={() => handleSetAuthHeader(token)}
            >
              <EarthLock class="mr-1 size-4" />
              <span>Set as header</span>
            </button>
          </Tooltip>
        </div>
      {/snippet}

      <Menu content={menuContent} placement="top" trigger="mouseenter click">
        <button class="btn btn-xs flex-none">
          <Key class="mr-1 size-3" />
          <span>{token.key}</span>
        </button>
      </Menu>
    {/each}
  </div>
{/if}
