<script lang="ts">
  let { children } = $props();
  import "../app.css";
  import { onMount } from "svelte";
  import { fade } from "svelte/transition";
  import { Loader } from "@lucide/svelte";
  import { Toaster } from "svelte-sonner";
  import { initWasm, waitUntilInitialized } from "$lib/urpc";
  import {
    loadJsonSchemaFromUrpcSchemaUrl,
    loadStore,
  } from "$lib/store.svelte";

  // Initialize the store
  onMount(() => loadStore());

  // Initialize the WebAssembly binary
  let initialized = $state(false);
  onMount(async () => {
    await initWasm();
    await waitUntilInitialized();
    await loadJsonSchemaFromUrpcSchemaUrl("./schema.urpc");
    initialized = true;
  });
</script>

{#if !initialized}
  <main
    out:fade={{ duration: 200 }}
    class="flex flex-col fixed top-0 left-0 h-screen w-screen items-center justify-center"
  >
    <img src="/assets/logo-square.png" alt="UFO RPC Logo" class="size-[150px]">
    <h1 class="text-3xl font-bold mb-2">UFO RPC Playground</h1>
    <h2 class="mb-6">Loading WebAssembly binary...</h2>
    <Loader class="animate animate-spin size-10" />
  </main>
{/if}

{#if initialized}
  <div transition:fade={{ duration: 200 }}>
    {@render children()}
  </div>
{/if}

<Toaster richColors closeButton duration={5000} />
