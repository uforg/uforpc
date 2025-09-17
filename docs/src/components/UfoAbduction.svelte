<script lang="ts">
  import { onMount } from "svelte";

  import LogoUfo from "./LogoUfo.svelte";

  interface Props {
    width?: number;
    abductionItems?: string[];
  }

  let {
    width = 150,
    abductionItems = [
      "🐄",
      "👽",
      "💀",
      "🐛",
      "😱",
      "🔥",
      "🤯",
      "😭",
      "⚠️",
      "💩",
      "☠️",
      "complexity",
      "chaos",
      "technical debt",
      "outdated docs",
      "type errors",
      "boilerplate",
      "manual coding",
      "schema drift",
      "version conflicts",
      "any types",
      "runtime errors",
      "config hell",
      "dependency hell",
      "integration pain",
      "sync nightmares",
      "API mismatch",
      "legacy cruft",
      "proto fatigue",
      "GraphQL complexity",
      "REST confusion",
      "maintenance burden",
      "debugging torture",
      "type gymnastics",
      "handwritten clients",
      "TODO: fix later",
      "// don't touch this",
    ],
  }: Props = $props();

  let currentItem = $state("");
  let isAbducting = $state(false);
  let abductionProgress = $state(0);

  let itemOpacity = $derived(isAbducting ? 1 - abductionProgress ** 3 : 0);
  let itemY = $derived(isAbducting ? 30 - abductionProgress * 100 : 30);
  let itemTransform = $derived(`translateX(-50%) translateY(${itemY}px)`);

  onMount(() => {
    const startAbduction = () => {
      if (isAbducting) return;

      currentItem =
        abductionItems[Math.floor(Math.random() * abductionItems.length)];
      isAbducting = true;
      abductionProgress = 0;

      const animate = () => {
        abductionProgress += 0.01;
        if (abductionProgress < 1) {
          requestAnimationFrame(animate);
        } else {
          isAbducting = false;
          currentItem = "";
        }
      };

      animate();
    };

    const interval = setInterval(startAbduction, 300);
    return () => clearInterval(interval);
  });
</script>

<div class="abduction-container">
  <div class="m-auto" style="width: {width}px">
    <LogoUfo class="w-full" />
  </div>

  <div class="beam-container" style="--beam-width: {width}px;">
    <div class="light-beam"></div>

    {#if currentItem}
      <div
        class="abducted-item"
        style="opacity: {itemOpacity}; transform: {itemTransform};"
      >
        {currentItem}
      </div>
    {/if}

    {#each Array(8) as _, i}
      <div
        class="particle"
        style="--delay: {i * 0.3}s; --random: {Math.random()};"
      ></div>
    {/each}
  </div>
</div>

<style>
  .abduction-container {
    position: relative;
    width: 100%;
  }

  .beam-container {
    position: relative;
    height: 120px;
    width: 100%;
    overflow: hidden;
    pointer-events: none;
  }

  .light-beam {
    position: absolute;
    top: 0;
    left: 50%;
    transform: translateX(-50%);
    width: var(--beam-width);
    height: 100%;
    background: linear-gradient(to bottom, #fbbf24 0%, transparent 100%);
    clip-path: polygon(45% 0, 55% 0, 100% 100%, 0% 100%);
    animation: beamShimmer 2s ease-in-out infinite;
  }

  .abducted-item {
    position: absolute;
    top: 60px;
    left: 50%;
    font-size: 1.2rem;
    font-weight: bold;
    pointer-events: none;
    user-select: none;
  }

  .particle {
    position: absolute;
    left: calc(50% + (var(--random) - 0.5) * var(--beam-width) * 0.8);
    top: 70%;
    width: 5px;
    height: 5px;
    background: rgba(251, 191, 36, 0.7);
    border-radius: 50%;
    opacity: 0;
    animation: floatUp 2s linear infinite;
    animation-delay: var(--delay, 0s);
  }

  @keyframes beamShimmer {
    0%,
    100% {
      opacity: 0.6;
    }
    50% {
      opacity: 1;
    }
  }

  @keyframes floatUp {
    0% {
      transform: translateY(0) scale(0.5);
      opacity: 0;
    }
    50% {
      transform: translateY(0) scale(0.8);
      opacity: 0.8;
    }
    100% {
      transform: translateY(-80px) scale(0);
      opacity: 0;
    }
  }
</style>
