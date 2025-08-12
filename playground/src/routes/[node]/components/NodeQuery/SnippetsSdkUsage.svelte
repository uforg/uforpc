<script lang="ts">
  import { uiStore } from "$lib/uiStore.svelte";

  interface Props {
    type: "proc" | "stream";
    name: string;
  }

  const { type, name }: Props = $props();
</script>

Usage
