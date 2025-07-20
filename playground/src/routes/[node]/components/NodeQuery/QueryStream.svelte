<script lang="ts">
  import { Info, Loader, MoveDownLeft, MoveUpRight, Zap } from "@lucide/svelte";
  import { toast } from "svelte-sonner";

  import { joinPath } from "$lib/helpers/joinPath";
  import { getHeadersObject, store } from "$lib/store.svelte";
  import { uiStore } from "$lib/uiStore.svelte";
  import type { StreamDefinitionNode } from "$lib/urpcTypes";

  import H2 from "$lib/components/H2.svelte";

  import Field from "./Field.svelte";
  import Output from "./Output.svelte";
  import Snippets from "./Snippets.svelte";

  interface Props {
    stream: StreamDefinitionNode;
  }

  const { stream }: Props = $props();

  let value = $state({ root: {} });
  let output: string | null = $state(null);
  let isExecuting = $state(false);
  let cancelRequest = $state<() => void>(() => {});

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

      const endpoint = joinPath([store.baseUrl, stream.name]);
      const response = await fetch(endpoint, {
        method: "POST",
        body: JSON.stringify(value.root ?? {}),
        headers: {
          ...getHeadersObject(),
          Accept: "text/event-stream",
          "Cache-Control": "no-cache",
        },
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

              // Add message to output at the beginning
              if (output) {
                output = `${JSON.stringify(parsedData, null, 2)}\n\n${output}`;
              } else {
                output = JSON.stringify(parsedData, null, 2);
              }
            } catch (parseError) {
              // If not JSON, treat as plain text
              if (output) {
                output = `${eventData}\n${"─".repeat(50)}\n${output}`;
              } else {
                output = eventData;
              }
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

<div class="flex" bind:this={wrapper}>
  <div class="flex-grow">
    <H2 class="mb-4 flex items-center space-x-2">Try it out</H2>

    <div
      class={{
        "join bg-base-100 flex w-full": true,
        "sticky top-[72px] z-10 -mt-4 pt-4": !uiStore.isMobile,
      }}
    >
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

    <div
      class={{
        "space-y-4": true,
        hidden: tab === "output",
        block: tab === "input",
      }}
    >
      {#if stream.input}
        <div role="alert" class="alert alert-soft alert-info mt-6 w-fit">
          <Info class="size-4" />
          <span>
            Requests are made from your browser and validations are performed on
            the server side
          </span>
        </div>
        <Field fields={stream.input} path="root" bind:value />
      {:else}
        <div role="alert" class="alert alert-soft alert-warning mt-6 w-fit">
          <Info class="size-4" />
          <span>This stream does not require any input</span>
        </div>
      {/if}

      <div class="flex w-full justify-end gap-2 pt-4">
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

  {#if !uiStore.isMobile}
    <div class="divider divider-horizontal"></div>

    <div class="w-[40%] flex-none">
      <div class="sticky top-[72px] z-10 -mt-4 pt-4">
        <Snippets {value} type="stream" name={stream.name} />
      </div>
    </div>
  {/if}
</div>

{#if uiStore.isMobile}
  <div class="mt-12">
    <Snippets {value} type="stream" name={stream.name} />
  </div>
{/if}
