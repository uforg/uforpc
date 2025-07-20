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
  class="flex h-[90px] flex-col items-center justify-center space-y-2 px-4 pt-4 text-center"
>
  <Logo class="flex-grow" />
  <span class="text-base-content/50 text-sm">
    UFO RPC is not affiliated with OpenAPI or Swagger. It just generates OpenAPI
    schemas from UFO RPC code.
  </span>
</header>

<div id={elId}></div>
