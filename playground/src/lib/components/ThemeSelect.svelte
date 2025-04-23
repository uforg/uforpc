<!--
  This uses the following helper to handle theme changes and is loaded
  in the head of the document to prevent any flash of unstyled content.

  /static/theme-helper.js
-->

<script lang="ts">
  import { onMount } from "svelte";

  interface Props {
    class?: string;
  }
  let { class: className = "" }: Props = $props();

  let currentTheme = $state("");
  onMount(() => {
    const theme = (window as any).getTheme();
    currentTheme = theme || "system";
  });

  $effect(() => {
    (window as any).setTheme(currentTheme);
  });
</script>

<label class="select {className}">
  <span class="label">Theme</span>
  <select bind:value={currentTheme}>
    <option value="system">System 🖥️</option>
    <option value="corporate">Light ☀️</option>
    <option value="dark">Dark 🌑</option>
    <option value="dracula">Dracula 🧛</option>
  </select>
</label>
