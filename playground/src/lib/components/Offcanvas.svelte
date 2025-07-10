<script lang="ts">
  import type { Snippet } from "svelte";
  import Portal from "svelte-portal";
  import { fade } from "svelte/transition";

  import { mergeClasses } from "$lib/helpers/mergeClasses";
  import type { ClassValue } from "$lib/helpers/mergeClasses";

  interface Props {
    children?: Snippet;
    direction?: "left" | "right";
    isOpen?: boolean;
    class?: ClassValue;
    backdropClass?: ClassValue;
    backdropClose?: boolean;
    escapeClose?: boolean;
  }

  let {
    children,
    direction = "left",
    isOpen = $bindable(false),
    class: className,
    backdropClass: backdropClassName,
    backdropClose = true,
    escapeClose = true,
  }: Props = $props();

  const closeOffcanvas = () => (isOpen = false);

  const handleEscapeKey = (event: KeyboardEvent) => {
    if (event.key === "Escape") closeOffcanvas();
  };
  $effect(() => {
    if (escapeClose) {
      document.addEventListener("keydown", handleEscapeKey);
      return () => {
        document.removeEventListener("keydown", handleEscapeKey);
      };
    }
  });
</script>

{#if isOpen}
  <Portal target="body">
    <div
      class="fixed top-0 left-0 z-40 h-screen w-screen"
      transition:fade={{ duration: 100 }}
    >
      <button
        class={mergeClasses(
          "z-10 h-full w-full bg-black/30",
          backdropClassName,
        )}
        onclick={backdropClose ? closeOffcanvas : undefined}
        aria-label="Close modal"
      >
      </button>

      <div
        class={mergeClasses(
          "absolute top-0 z-20",
          "bg-base-100 h-[100dvh] w-[280px]",
          "overflow-x-hidden overflow-y-auto",
          {
            "left-0": direction === "left",
            "right-0": direction === "right",
          },
          className,
        )}
      >
        {#if children}
          {@render children()}
        {/if}
      </div>
    </div>
  </Portal>
{/if}
