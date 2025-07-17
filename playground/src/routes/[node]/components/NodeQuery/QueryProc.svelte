<script lang="ts">
  import { Info, Loader, MoveDownLeft, MoveUpRight, Zap } from "@lucide/svelte";
  import { toast } from "svelte-sonner";

  import { joinPath } from "$lib/helpers/joinPath";
  import { getHeadersObject, store } from "$lib/store.svelte";
  import { uiStore } from "$lib/uiStore.svelte";
  import type { ProcedureDefinitionNode } from "$lib/urpcTypes";

  import H2 from "$lib/components/H2.svelte";

  import Field from "./Field.svelte";
  import Output from "./Output.svelte";
  import Snippets from "./Snippets.svelte";

  interface Props {
    proc: ProcedureDefinitionNode;
  }

  const { proc }: Props = $props();

  let value = $state({ root: {} });
  let output: string | null = $state(null);
  let isExecuting = $state(false);
  let cancelRequest = $state<() => void>(() => {});

  async function executeProcedure() {
    if (isExecuting) return;
    isExecuting = true;
    output = null;

    try {
      openOutput(true);
      const controller = new AbortController();
      const signal = controller.signal;

      cancelRequest = () => {
        controller.abort();
        toast.info("Procedure call cancelled");
      };

      const endpoint = joinPath([store.baseUrl, proc.name]);
      const response = await fetch(endpoint, {
        method: "POST",
        body: JSON.stringify(value.root ?? {}),
        headers: getHeadersObject(),
        signal: signal,
      });

      const data = await response.json();
      output = JSON.stringify(data, null, 2);
    } catch (error) {
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
      {#if proc.input}
        <div role="alert" class="alert alert-soft alert-info mt-6 w-fit">
          <Info class="size-4" />
          <span>
            Requests are made from your browser and validations are performed on
            the server side
          </span>
        </div>
        <Field fields={proc.input} path="root" bind:value />
      {:else}
        <div role="alert" class="alert alert-soft alert-warning mt-6 w-fit">
          <Info class="size-4" />
          <span>This procedure does not require any input</span>
        </div>
      {/if}

      <div class="flex w-full justify-end pt-4">
        <button
          class="btn btn-primary"
          disabled={isExecuting}
          onclick={executeProcedure}
        >
          {#if isExecuting}
            <Loader class="animate size-4 animate-spin" />
          {:else}
            <Zap class="size-4" />
          {/if}
          <span>Execute procedure</span>
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
      <Output {cancelRequest} {isExecuting} type="proc" {output} />
    </div>
  </div>

  {#if !uiStore.isMobile}
    <div class="divider divider-horizontal"></div>

    <div class="w-[40%] flex-none">
      <div class="sticky top-[72px] z-10 -mt-4 pt-4">
        <Snippets {value} type="proc" name={proc.name} />
      </div>
    </div>
  {/if}
</div>

{#if uiStore.isMobile}
  <div class="mt-12">
    <Snippets {value} type="proc" name={proc.name} />
  </div>
{/if}
