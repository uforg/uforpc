<script lang="ts">
  import { Info, Loader, MoveDownLeft, MoveUpRight, Zap } from "@lucide/svelte";
  import { toast } from "svelte-sonner";

  import { ctrlSymbol } from "$lib/helpers/ctrlSymbol";
  import { joinPath } from "$lib/helpers/joinPath";
  import { getHeadersObject, storeSettings } from "$lib/storeSettings.svelte";
  import { storeUi } from "$lib/storeUi.svelte";
  import type { StreamDefinitionNode } from "$lib/urpcTypes";

  import H2 from "$lib/components/H2.svelte";
  import Menu from "$lib/components/Menu.svelte";

  import InputForm from "./InputForm/InputForm.svelte";
  import Output from "./Output.svelte";
  import Snippets from "./Snippets/Snippets.svelte";

  interface Props {
    stream: StreamDefinitionNode;
    input: any;
    output: string | null;
  }

  let { stream, input = $bindable(), output = $bindable() }: Props = $props();

  // biome-ignore lint/suspicious/noExplicitAny: can be any stream response
  let outputArray: any[] = $state([]);
  let isExecuting = $state(false);
  let cancelRequest = $state<() => void>(() => {});

  // let output = $derived(JSON.stringify(outputArray, null, 2));

  async function executeStream() {
    if (isExecuting) return;
    isExecuting = true;
    output = "";

    try {
      openOutput(true);
      const controller = new AbortController();
      const signal = controller.signal;

      cancelRequest = () => {
        controller.abort();
        toast.info("Stream stopped");
      };

      const endpoint = joinPath([storeSettings.store.baseUrl, stream.name]);
      const headers = getHeadersObject();
      headers.set("Accept", "text/event-stream");
      headers.set("Cache-Control", "no-cache");

      const response = await fetch(endpoint, {
        method: "POST",
        body: JSON.stringify(input.root ?? {}),
        headers,
        signal: signal,
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const reader = response.body?.getReader();
      if (!reader) {
        throw new Error("Response body is null");
      }
      const decoder = new TextDecoder();
      let buffer = "";

      while (true) {
        const { done, value } = await reader.read();
        if (done) {
          toast.info("Stream ended by server");
          break;
        }

        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split("\n");

        // Keep the last incomplete line in buffer
        buffer = lines.pop() || "";

        for (const line of lines) {
          if (line.trim() === "") continue;

          if (line.startsWith("data: ")) {
            const eventData = line.slice(6);

            // Skip heartbeat or keep-alive messages
            if (eventData.trim() === "" || eventData.trim() === "heartbeat") {
              continue;
            }

            try {
              const parsedData = JSON.parse(eventData);
              outputArray.unshift(parsedData);
            } catch (parseError) {
              outputArray.unshift(eventData);
            }
          }
        }
      }
    } catch (error: unknown) {
      if (!(error instanceof Error && error.name === "AbortError")) {
        console.error(error);
        toast.error("Failed to send HTTP request", {
          description: `Error: ${error}`,
          duration: 5000,
        });
      }
    } finally {
      isExecuting = false;
      cancelRequest = () => {};
    }
  }

  async function executeStreamFromKbd(event: KeyboardEvent) {
    // CTRL/CMD + ENTER to execute
    if (event.key === "Enter" && (event.ctrlKey || event.metaKey)) {
      event.preventDefault();
      await executeStream();
    }
  }

  let tab: "input" | "output" = $state("input");
  let wrapper: HTMLDivElement | null = $state(null);
  function openInput(scroll = false) {
    if (tab === "input") return;
    tab = "input";
    if (scroll) wrapper?.scrollIntoView({ behavior: "smooth", block: "start" });
  }
  function openOutput(scroll = false) {
    if (tab === "output") return;
    tab = "output";
    if (scroll) wrapper?.scrollIntoView({ behavior: "smooth", block: "start" });
  }
</script>

<div bind:this={wrapper}>
  <div
    class={{
      "bg-base-100 sticky top-0 z-20 pt-4": !storeUi.store.isMobile,
    }}
  >
    <H2 class="mb-4 flex items-center space-x-2">Try it out</H2>

    <div class="join flex w-full">
      <button
        class={[
          "btn join-item border-base-content/20 flex-grow",
          tab === "input" && "btn-primary btn-active",
        ]}
        onclick={() => openInput(false)}
      >
        <MoveUpRight class="size-4" />
        Input
      </button>
      <button
        class={[
          "btn join-item border-base-content/20 flex-grow",
          tab === "output" && "btn-primary btn-active",
        ]}
        onclick={() => openOutput(false)}
      >
        <MoveDownLeft class="size-4" />
        <span>Output</span>
      </button>
    </div>
  </div>

  <div
    class={{
      "space-y-4": true,
      hidden: tab === "output",
      block: tab === "input",
    }}
  >
    {#if stream.input}
      <div
        class="space-y-4"
        onkeydown={executeStreamFromKbd}
        role="button"
        tabindex="0"
      >
        <InputForm fields={stream.input} bind:input />
      </div>
    {:else}
      <div role="alert" class="alert alert-soft alert-warning mt-6 w-fit">
        <Info class="size-4" />
        <span>This stream does not require any input</span>
      </div>
    {/if}

    <div class="flex w-full justify-end gap-2 pt-4">
      {#snippet kbd()}
        <span>
          <kbd class="kbd kbd-sm">{ctrlSymbol()}</kbd>
          <kbd class="kbd kbd-sm">⤶</kbd>
        </span>
      {/snippet}

      <Menu content={kbd} placement="left" trigger="mouseenter">
        <button
          class="btn btn-primary"
          disabled={isExecuting}
          onclick={executeStream}
        >
          {#if isExecuting}
            <Loader class="animate size-4 animate-spin" />
          {:else}
            <Zap class="size-4" />
          {/if}
          <span>Start stream</span>
        </button>
      </Menu>
    </div>
  </div>

  <div
    class={{
      "mt-4 space-y-2": true,
      hidden: tab === "input",
      block: tab === "output",
    }}
  >
    <Output {cancelRequest} {isExecuting} type="stream" {output} />
  </div>
</div>

{#if storeUi.store.isMobile}
  <div class="mt-12">
    <Snippets {input} type="stream" name={stream.name} />
  </div>
{/if}
