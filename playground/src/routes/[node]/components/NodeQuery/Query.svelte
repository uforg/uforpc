<script lang="ts">
  import { Info, Loader, MoveDownLeft, MoveUpRight, Zap } from "@lucide/svelte";
  import { toast } from "svelte-sonner";

  import { getHeadersObject, store } from "$lib/store.svelte";
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
  let output: object | null = $state(null);

  let isExecuting = $state(false);
  async function executeProcedure() {
    if (isExecuting) return;
    isExecuting = true;
    output = null;

    try {
      const response = await fetch(store.endpoint, {
        method: "POST",
        body: JSON.stringify({
          proc: proc.name,
          input: value.root,
        }),
        headers: getHeadersObject(),
      });

      const data = await response.json();
      output = data;

      openOutput();
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
  function openInput() {
    if (tab === "input") return;
    tab = "input";
    wrapper?.scrollIntoView({ behavior: "smooth", block: "start" });
  }
  function openOutput() {
    if (tab === "output") return;
    tab = "output";
    wrapper?.scrollIntoView({ behavior: "smooth", block: "start" });
  }
</script>

{#if proc.input}
  <div class="flex space-x-4" bind:this={wrapper}>
    <div class="flex-grow space-y-2 rounded-t-none">
      <div
        class={[
          "sticky top-[72px] z-10 flex w-full items-center justify-between",
          "bg-base-100 -mt-4 pt-4 pb-2",
        ]}
      >
        <H2>Try {proc.name}</H2>
        <div class="join">
          <button
            class={[
              "btn join-item",
              tab === "input" && "btn-primary btn-active",
            ]}
            onclick={openInput}
          >
            <MoveUpRight class="size-4" />
            Input
          </button>
          <button
            class={[
              "btn join-item",
              tab === "output" && "btn-primary btn-active",
            ]}
            onclick={openOutput}
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

        <Field fields={proc.input} path="root" bind:value />

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
        <Output {output} />
      </div>
    </div>

    <Snippets {value} procName={proc.name} />
  </div>
{/if}
