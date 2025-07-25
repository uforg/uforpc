<script lang="ts">
  import {
    BookOpenText,
    Brain,
    Cable,
    Github,
    NotebookText,
    ShieldCheck,
  } from "@lucide/svelte";

  import { uiStore } from "$lib/uiStore.svelte";

  import Logo from "$lib/components/Logo.svelte";

  import UfoAbduction from "../lib/components/UfoAbduction.svelte";
</script>

<svelte:head>
  <title>UFO RPC Playground</title>
</svelte:head>

<div
  class={{
    "container mx-auto": true,
    "py-16": !uiStore.isMobile,
  }}
>
  <div class="mb-16 text-center">
    {#if !uiStore.isMobile}
      <div class="mb-8">
        <Logo class="mx-auto h-full max-w-[600px]" />
      </div>
    {/if}

    <p
      class={{
        "mx-auto mb-8 max-w-2xl text-xl": true,
        "text-base-content/70": !uiStore.isMobile,
      }}
    >
      Test and experiment with UFO RPC in this interactive playground using the
      left sidebar to navigate the schema.
    </p>

    <div class="mb-8 flex items-center justify-center">
      <UfoAbduction width={200} />
    </div>

    <div class="flex items-center justify-center space-x-2">
      <a
        href="https://github.com/uforg/uforpc"
        target="_blank"
        class={{
          "btn btn-primary": true,
          "btn-lg": !uiStore.isMobile,
        }}
      >
        <Github class="size-5" />
        View on GitHub
      </a>
      <a
        href="https://uforpc.uforg.dev"
        target="_blank"
        class={{
          "btn btn-outline": true,
          "btn-lg": !uiStore.isMobile,
        }}
      >
        <BookOpenText class="size-5" />
        Documentation
      </a>
    </div>
  </div>

  <div class="desk:grid-cols-2 mb-16 grid gap-8">
    <div class="card bg-base-200 shadow-sm">
      <div class="card-body">
        <h2 class="card-title mb-4 text-2xl">What is UFO RPC?</h2>
        <p class="text-base-content/80 leading-relaxed">
          UFO RPC is a DX-focused RPC system that provides type safety and
          cross-language compatibility. It includes its own DSL (Domain Specific
          Language) that makes it easy to define and validate RPC services with
          a clean, intuitive syntax.
        </p>
      </div>
    </div>

    <div class="card bg-base-200 shadow-sm">
      <div class="card-body">
        <h2 class="card-title mb-4 text-2xl">About this Playground</h2>
        <p class="text-base-content/80 leading-relaxed">
          This playground provides an interactive environment where you can test
          and experiment with a provided UFO RPC schema.
        </p>
      </div>
    </div>
  </div>

  <div class="text-center">
    <h2 class="mb-8 text-3xl font-bold">Key Features</h2>
    <div class="desk:grid-cols-4 grid gap-6">
      <div class="card bg-base-200 shadow-sm">
        <div class="card-body text-center">
          <div
            class="bg-primary/10 mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full"
          >
            <span class="text-2xl">
              <ShieldCheck class="size-6" />
            </span>
          </div>
          <h3 class="mb-2 font-semibold">Type Safety</h3>
          <p class="text-base-content/70 text-sm">
            Strong type validation and business rules.
          </p>
        </div>
      </div>

      <div class="card bg-base-200 shadow-sm">
        <div class="card-body text-center">
          <div
            class="bg-secondary/10 mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full"
          >
            <span class="text-2xl">
              <NotebookText class="size-6" />
            </span>
          </div>
          <h3 class="mb-2 font-semibold">Declarative Syntax</h3>
          <p class="text-base-content/70 text-sm">
            Intuitive and human-readable format.
          </p>
        </div>
      </div>

      <div class="card bg-base-200 shadow-sm">
        <div class="card-body text-center">
          <div
            class="bg-accent/10 mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full"
          >
            <span class="text-2xl">
              <Cable class="size-6" />
            </span>
          </div>
          <h3 class="mb-2 font-semibold">Cross-Language Compatibility</h3>
          <p class="text-base-content/70 text-sm">
            Generate efficient type-safe server and client code from your UFO
            RPC schema.
          </p>
        </div>
      </div>

      <div class="card bg-base-200 shadow-sm">
        <div class="card-body text-center">
          <div
            class="bg-info/10 mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full"
          >
            <span class="text-2xl">
              <Brain class="size-6" />
            </span>
          </div>
          <h3 class="mb-2 font-semibold">Developer Experience</h3>
          <p class="text-base-content/70 text-sm">
            Best possible DX with strict data integrity.
          </p>
        </div>
      </div>
    </div>
  </div>
</div>
