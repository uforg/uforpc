<script lang="ts">
  import { Loader } from "@lucide/svelte";
  import { toast } from "svelte-sonner";

  import { store } from "$lib/store.svelte";
  import {
    cmdCodegen,
    type CmdCodegenOptions,
    type CmdCodegenOutput,
    type CodegenGenerator,
  } from "$lib/urpc";

  import Code from "$lib/components/Code.svelte";
  import H2 from "$lib/components/H2.svelte";

  interface Props {
    // biome-ignore lint/suspicious/noExplicitAny: consistent with sibling components
    value: any;
    type: "proc" | "stream";
    name: string;
  }

  const { value, type, name }: Props = $props();

  let generator: CodegenGenerator = $state("typescript-client");
  let golangPackageName: string = $state("client");
  let dartPackageName: string = $state("client");

  let isGenerating: boolean = $state(false);
  let output: CmdCodegenOutput | null = $state(null);

  function getLangFromGenerator(gen: CodegenGenerator): string {
    if (gen === "golang-client") return "go";
    if (gen === "dart-client") return "dart";
    return "ts";
  }

  async function generate() {
    if (isGenerating) return;
    isGenerating = true;
    output = null;

    try {
      let opts: CmdCodegenOptions = {
        generator,
        schemaInput: store.urpcSchema,
      };
      if (generator === "golang-client") {
        opts.golangPackageName = golangPackageName.trim() || "client";
      }
      if (generator === "dart-client") {
        opts.dartPackageName = dartPackageName.trim() || "client";
      }

      const result = await cmdCodegen(opts);
      output = result;
    } catch (error) {
      console.error(error);
      toast.error("Failed to generate SDK", {
        description: String(error),
        duration: 5000,
      });
    } finally {
      isGenerating = false;
    }
  }
</script>

<H2 class="mb-2 flex items-center space-x-2">SDK Generator</H2>

<div class="space-y-4">
  <div class="form-control">
    <label class="label" for="sdk-generator-select">
      <span class="label-text">Select generator</span>
    </label>
    <select
      id="sdk-generator-select"
      class="select select-bordered w-full"
      bind:value={generator}
    >
      <option value="typescript-client">TypeScript client</option>
      <option value="golang-client">Go client</option>
      <option value="dart-client">Dart client</option>
    </select>
  </div>

  {#if generator === "golang-client"}
    <div class="form-control">
      <label class="label" for="go-pkg">
        <span class="label-text">Go package name</span>
      </label>
      <input
        id="go-pkg"
        class="input input-bordered w-full"
        placeholder="client"
        bind:value={golangPackageName}
      />
    </div>
  {/if}

  {#if generator === "dart-client"}
    <div class="form-control">
      <label class="label" for="dart-pkg">
        <span class="label-text">Dart package name</span>
      </label>
      <input
        id="dart-pkg"
        class="input input-bordered w-full"
        placeholder="client"
        bind:value={dartPackageName}
      />
    </div>
  {/if}

  <div class="flex w-full justify-end gap-2 pt-2">
    <button
      class="btn btn-primary"
      disabled={isGenerating}
      onclick={generate}
      type="button"
    >
      {#if isGenerating}
        <Loader class="animate size-4 animate-spin" />
      {/if}
      <span>Generate</span>
    </button>
  </div>

  {#if output}
    <div class="divider">Generated files</div>
    {#each output.files as f, i (f.path + i)}
      <div class="card bg-base-100 border-base-content/20 mb-4 border">
        <div class="card-body p-4">
          <div class="pb-2 font-mono text-xs opacity-70">{f.path}</div>
          <Code code={f.content} lang={getLangFromGenerator(generator)} />
        </div>
      </div>
    {/each}
  {/if}
</div>
