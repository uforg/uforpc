<script lang="ts">
  import { Info } from "@lucide/svelte";

  import type { ProcedureDefinitionNode } from "$lib/urpcTypes";

  import H2 from "$lib/components/H2.svelte";

  import Field from "./Field.svelte";
  import Snippets from "./Snippets.svelte";

  interface Props {
    proc: ProcedureDefinitionNode;
  }

  const { proc }: Props = $props();

  let value = $state({ root: {} });
  let output = $state({});

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
              "btn  btn-soft btn-primary join-item",
              tab === "input" && "btn-active",
            ]}
            onclick={openInput}
          >
            Input
          </button>
          <button
            class={[
              "btn  btn-soft btn-primary join-item",
              tab === "output" && "btn-active",
            ]}
            onclick={openOutput}
          >
            Output
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
      </div>

      <div
        class={{
          "space-y-2": true,
          hidden: tab === "input",
          block: tab === "output",
        }}
      >
        Output
      </div>
    </div>

    <Snippets {value} procName={proc.name} />
  </div>
{/if}
