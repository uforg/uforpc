<script lang="ts">
  import { Info, Loader, MoveDownLeft, MoveUpRight, Zap } from "@lucide/svelte";
  import { toast } from "svelte-sonner";

  import { getHeadersObject, store } from "$lib/store.svelte";
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
  let output: object | null = $state(null);

  let isExecuting = $state(false);
  async function executeStream() {
    if (isExecuting) return;
    isExecuting = true;
    output = null;

    try {
      const response = await fetch(store.endpoint, {
        method: "POST",
        body: JSON.stringify({
          type: "stream",
          name: stream.name,
          input: value.root,
        }),
        headers: getHeadersObject(),
      });

      const data = await response.json();
      output = data;

      openOutput(true);
    } catch (error) {
      console.error(error);
      toast.error("Failed to send HTTP request", {
        description: `Error: ${error}`,
        duration: 5000,
      });
    } finally {
      isExecuting = false;
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

{#if stream.input}
  <div class="flex space-x-4" bind:this={wrapper}>
    <div class="flex-grow space-y-2 rounded-t-none">
      <div
        class={[
          "sticky top-[72px] z-10 flex w-full items-center justify-between",
          "bg-base-100 -mt-4 space-x-2 pt-4 pb-2",
        ]}
      >
        <H2 class="flex items-center space-x-2 break-all">
          <Zap class="size-6 flex-none" />
          <span>Try {stream.name}</span>
        </H2>
        <div class="join">
          <button
            class={[
              "btn join-item border-base-content/20",
              tab === "input" && "btn-primary btn-active",
            ]}
            onclick={() => openInput(false)}
          >
            <MoveUpRight class="size-4" />
            Input
          </button>
          <button
            class={[
              "btn join-item border-base-content/20",
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
          "space-y-2": true,
          hidden: tab === "output",
          block: tab === "input",
        }}
      >
        <div role="alert" class="alert alert-soft alert-info mt-6 w-fit">
          <Info class="size-4" />
          <span> All validations are performed on the server side </span>
        </div>

        <Field fields={stream.input} path="root" bind:value />

        <div class="flex w-full justify-end pt-4">
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
          "space-y-2": true,
          hidden: tab === "input",
          block: tab === "output",
        }}
      >
        <Output {output} />
      </div>
    </div>

    <Snippets {value} type="stream" name={stream.name} />
  </div>
{/if}
