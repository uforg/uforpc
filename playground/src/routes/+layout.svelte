<script lang="ts">
  let { children } = $props();
  import "../app.css";
  import { onMount } from "svelte";
  import { fade } from "svelte/transition";
  import { Loader } from "@lucide/svelte";
  import { toast, Toaster } from "svelte-sonner";
  import { initWasm, waitUntilInitialized } from "$lib/urpc";
  import { initializeShiki } from "$lib/shiki";
  import { loadUiStore } from "$lib/uiStore.svelte";
  import {
    loadJsonSchemaFromUrpcSchemaUrl,
    loadStore,
  } from "$lib/store.svelte";

  // Initialize the stores
  onMount(() => {
    loadUiStore();
    loadStore();
  });

  // Initialize the WebAssembly binary
  let initialized = $state(false);
  let message = $state("Starting playground");
  onMount(async () => {
    const handleError = (error: unknown) => {
      console.error(error);
      toast.error("Failed to initialize the Playground", {
        description: `Please try again or contact the developers if the problem persists. Error: ${error}`,
        duration: 15000,
      });
    };

    message = "Loading code highlighter";
    try {
      await initializeShiki();
    } catch (error) {
      handleError(error);
      return;
    }

    message = "Loading URPC WebAssembly binary";
    try {
      await initWasm();
      await waitUntilInitialized();
    } catch (error) {
      handleError(error);
      return;
    }

    message = "Loading URPC Schema";
    try {
      await loadJsonSchemaFromUrpcSchemaUrl("./schema.urpc");
    } catch (error) {
      handleError(error);
    }

    initialized = true;
  });
</script>

{#if !initialized}
  <main
    out:fade={{ duration: 200 }}
    class="fixed top-0 left-0 flex h-screen w-screen flex-col items-center justify-center"
  >
    <img
      src="/assets/logo-square.png"
      alt="UFO RPC Logo"
      class="size-[150px]"
    />
    <h1 class="mb-2 text-3xl font-bold">UFO RPC Playground</h1>
    <h2 class="mb-6">{message}...</h2>
    <Loader class="animate size-10 animate-spin" />
  </main>
{/if}

{#if initialized}
  <div transition:fade={{ duration: 200 }}>
    {@render children()}
  </div>
{/if}

<Toaster richColors closeButton duration={5000} />
