<script lang="ts">
  import SwaggerUI from "swagger-ui";
  import "swagger-ui/dist/swagger-ui.css";

  import Logo from "$lib/components/Logo.svelte";

  const elId = $props.id();

  $effect(() => {
    SwaggerUI({
      domNode: document.getElementById(elId),
    });
  });
</script>

<header
  class="mt-12 flex h-[80px] flex-col items-center justify-center space-y-2 px-4 text-center"
>
  <Logo class="flex-grow" />
  <span class="text-base-content/50 text-sm">
    UFO RPC is not affiliated with OpenAPI or Swagger. It just generates OpenAPI
    schemas from UFO RPC code.
  </span>
</header>

<div id={elId}></div>

<style>
  :global([data-theme="dark"] .swagger-ui) {
    filter: invert(88%) hue-rotate(180deg);
  }
  :global([data-theme="dark"] .swagger-ui .microlight) {
    filter: invert(100%) hue-rotate(180deg);
  }
</style>
