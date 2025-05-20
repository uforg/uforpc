<script lang="ts">
  import { onNavigate } from "$app/navigation";
  import { Loader } from "@lucide/svelte";
  import { onMount } from "svelte";
  import { toast, Toaster } from "svelte-sonner";
  import { fade } from "svelte/transition";

  import { initializeShiki } from "$lib/shiki";
  import {
    loadJsonSchemaFromUrpcSchemaUrl,
    loadStore,
  } from "$lib/store.svelte";
  import {
    dimensionschangeAction,
    loadUiStore,
    uiStore,
  } from "$lib/uiStore.svelte";
  import { initWasm, waitUntilInitialized } from "$lib/urpc";

  import "../app.css";

  import LayoutAside from "./components/LayoutAside.svelte";
  import LayoutHeader from "./components/LayoutHeader.svelte";

  let { children } = $props();

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

  // Handle view transitions
  onNavigate((navigation) => {
    if (!document.startViewTransition) return;

    return new Promise((resolve) => {
      document.startViewTransition(async () => {
        resolve();
        await navigation.complete;
      });
    });
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
  <div
    transition:fade={{ duration: 200 }}
    use:dimensionschangeAction
    ondimensionschange={(e) => (uiStore.app = e.detail)}
    class="flex h-[100dvh] w-[100dvw] justify-start"
  >
    <LayoutAside />
    <div
      use:dimensionschangeAction
      ondimensionschange={(e) => (uiStore.contentWrapper = e.detail)}
      class="h-[100dvh] flex-grow scroll-p-[90px] overflow-x-hidden overflow-y-auto"
    >
      <LayoutHeader />
      <main
        class="w-full p-4"
        use:dimensionschangeAction
        ondimensionschange={(e) => (uiStore.main = e.detail)}
      >
        {@render children()}
      </main>
    </div>
  </div>
{/if}

<Toaster richColors closeButton duration={5000} />
