<script lang="ts">
  import { browser } from "$app/environment";
  import { onNavigate } from "$app/navigation";
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
    initTheme,
    uiStore,
  } from "$lib/uiStore.svelte";
  import { initWasm, waitUntilInitialized } from "$lib/urpc";

  import Logo from "$lib/components/Logo.svelte";

  import "../app.css";

  import LayoutAside from "./components/LayoutAside.svelte";
  import LayoutHeader from "./components/LayoutHeader.svelte";
  import LayoutSwaggerSwitch from "./components/LayoutSwaggerSwitch.svelte";

  let { children } = $props();

  // Initialize theme
  onMount(() => {
    if (!browser) return;
    initTheme();
  });

  // Initialize the stores
  onMount(() => {
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

  let mainWidth = $derived.by(() => {
    if (uiStore.isMobile) return uiStore.app.size.offsetWidth;
    return uiStore.app.size.offsetWidth - uiStore.aside.size.offsetWidth;
  });

  let mainHeight = $derived.by(() => {
    return uiStore.app.size.offsetHeight - uiStore.header.size.offsetHeight;
  });

  let mainStyle = $derived.by(() => {
    return `width: ${mainWidth}px; height: ${mainHeight}px;`;
  });
</script>

{#if !initialized}
  <main
    out:fade={{ duration: 200 }}
    class="fixed top-0 left-0 flex h-screen w-screen flex-col items-center justify-center"
  >
    <Logo class="mb-6 h-auto w-[90dvw] max-w-[600px]" />
    <h2>{message}...</h2>
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
      class="h-[100dvh] flex-grow scroll-p-[90px]"
    >
      <LayoutHeader />
      <main
        class="overflow-hidden"
        style={mainStyle}
        use:dimensionschangeAction
        ondimensionschange={(e) => (uiStore.main = e.detail)}
      >
        {@render children()}
      </main>
    </div>
  </div>

  <!-- Requires a space at bottom of the page to fit the switch button without covering other content -->
  <LayoutSwaggerSwitch />
{/if}

<Toaster richColors closeButton duration={5000} />
