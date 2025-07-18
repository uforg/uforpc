<script lang="ts">
  import { joinPath } from "$lib/helpers/joinPath";
  import { getHeadersObject, store } from "$lib/store.svelte";
  import { uiStore } from "$lib/uiStore.svelte";

  import Code from "$lib/components/Code.svelte";
  import H2 from "$lib/components/H2.svelte";

  import SnippetsCode from "./SnippetsCode.svelte";

  interface Props {
    // biome-ignore lint/suspicious/noExplicitAny: it's too dynamic to determine the type
    value: any;
    type: "proc" | "stream";
    name: string;
  }

  const { value, type, name }: Props = $props();

  let curl = $derived.by(() => {
    const endpoint = joinPath([store.baseUrl, name]);
    const payload = value.root ?? {};
    let payloadStr = JSON.stringify(payload, null, 2);
    payloadStr = payloadStr.replace(/'/g, "'\\''");

    let c = `curl -X POST ${endpoint} \\\n`;

    if (type === "stream") {
      c += "-N \\\n";
    }

    let headers = getHeadersObject();
    if (type === "stream") {
      headers.set("Accept", "text/event-stream");
      headers.set("Cache-Control", "no-cache");
    }

    for (const header of headers.entries()) {
      let rawHeader = `${header[0]}: ${header[1]}`;
      c += `-H ${JSON.stringify(rawHeader)} \\\n`;
    }

    c += `-d '${payloadStr}'`;

    return c;
  });

  let maxHeight = $derived.by(() => {
    if (uiStore.isMobile) return "100%";

    const appHeight = uiStore.app.size.offsetHeight;
    const headerHeight = uiStore.header.size.offsetHeight;
    const padding = 16 * 2;

    const mh = appHeight - headerHeight - padding;
    return `${mh}px`;
  });
</script>

<div
  class={{
    "flex h-full flex-col": !uiStore.isMobile,
  }}
  style="max-height: {maxHeight}"
>
  <H2 class="mb-2 flex items-center space-x-2">
    <span>Code snippets</span>
  </H2>

  {#if type === "stream"}
    <p class="pb-4 text-sm">
      Streams use Server-Sent Events. Only curl examples are provided. Build a
      client manually, or generate one with the urpc CLI if your language is
      supported.
      <br />
      <a href="https://uforpc.uforg.dev/r/sse" target="_blank" class="link">
        Learn more here
      </a>
    </p>

    <Code code={curl} lang="bash" />
  {/if}

  {#if type === "proc"}
    <SnippetsCode {curl} />
  {/if}
</div>
