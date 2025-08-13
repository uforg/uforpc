<script lang="ts">
  import { Download, Loader } from "@lucide/svelte";
  import { downloadZip } from "client-zip";
  import { toast } from "svelte-sonner";

  import { store } from "$lib/store.svelte";
  import { uiStore } from "$lib/uiStore.svelte";
  import {
    cmdCodegen,
    type CmdCodegenOptions,
    type CmdCodegenOutputFile,
  } from "$lib/urpc";

  let isGenerating: boolean = $state(false);

  let downloadFileName = $derived.by(() => {
    if (uiStore.codeSnippetsSdkLang === "typescript-client") {
      return "uforpc-client-sdk.ts";
    }
    if (uiStore.codeSnippetsSdkLang === "golang-client") {
      return "uforpc-client-sdk.go";
    }
    if (uiStore.codeSnippetsSdkLang === "golang-server") {
      return "uforpc-server-sdk.go";
    }
    if (uiStore.codeSnippetsSdkLang === "dart-client") {
      return "uforpc-dart-client-sdk.zip";
    }
    return "unknown";
  });

  function downloadSingleFile(file: CmdCodegenOutputFile) {
    const blob = new Blob([file.content], { type: "text/plain" });
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = downloadFileName;
    link.click();
    link.remove();
  }

  async function downloadMultipleFiles(files: CmdCodegenOutputFile[]) {
    const zipInputFiles = files.map((file) => ({
      name: file.path,
      input: file.content,
      lastModified: new Date(),
    }));

    const blob = await downloadZip(zipInputFiles).blob();
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = downloadFileName;
    link.click();
    link.remove();
  }

  async function generateAndDownload() {
    if (isGenerating) return;
    isGenerating = true;

    try {
      let opts: CmdCodegenOptions = {
        generator: uiStore.codeSnippetsSdkLang,
        schemaInput: store.urpcSchema,
      };
      if (uiStore.codeSnippetsSdkLang === "golang-client") {
        opts.golangPackageName =
          uiStore.codeSnippetsSdkGolangPackageName.trim();
        if (opts.golangPackageName === "") {
          throw new Error("Package name is required");
        }
      }
      if (uiStore.codeSnippetsSdkLang === "dart-client") {
        opts.dartPackageName = uiStore.codeSnippetsSdkDartPackageName.trim();
        if (opts.dartPackageName === "") {
          throw new Error("Package name is required");
        }
      }

      const result = await cmdCodegen(opts);

      if (result.files.length === 1) {
        downloadSingleFile(result.files[0]);
      }
      if (result.files.length > 1) {
        await downloadMultipleFiles(result.files);
      }
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

<div>
  {#if uiStore.codeSnippetsSdkLang === "golang-client"}
    <label class="fieldset">
      <legend class="fieldset-legend">Go package name</legend>
      <input
        id="go-pkg"
        class="input w-full"
        placeholder="Package name..."
        bind:value={uiStore.codeSnippetsSdkGolangPackageName}
      />
      <div class="prose prose-sm text-base-content/50 max-w-none">
        The generated SDK and code examples will use this package name.
      </div>
    </label>
  {/if}

  {#if uiStore.codeSnippetsSdkLang === "dart-client"}
    <label class="fieldset">
      <legend class="fieldset-legend">Dart package name</legend>
      <input
        id="go-pkg"
        class="input w-full"
        placeholder="Package name..."
        bind:value={uiStore.codeSnippetsSdkDartPackageName}
      />
      <div class="prose prose-sm text-base-content/50 max-w-none">
        The generated SDK and code examples will use this package name.
      </div>
    </label>
  {/if}

  <div class="fieldset">
    <legend class="fieldset-legend">Download SDK</legend>
    <button
      class="btn btn-primary btn-block"
      disabled={isGenerating}
      onclick={generateAndDownload}
      type="button"
    >
      {#if isGenerating}
        <Loader class="animate size-4 animate-spin" />
      {/if}
      {#if !isGenerating}
        <Download class="size-4" />
      {/if}
      <span>Download SDK</span>
    </button>
  </div>
</div>
