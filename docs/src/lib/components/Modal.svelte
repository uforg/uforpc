<script lang="ts">
  import type { Snippet } from "svelte";
  import Portal from "svelte-portal";
  import { fade } from "svelte/transition";

  import { mergeClasses } from "$lib/helpers/mergeClasses";
  import type { ClassValue } from "$lib/helpers/mergeClasses";

  interface Props {
    children?: Snippet;
    isOpen?: boolean;
    class?: ClassValue;
    backdropClass?: ClassValue;
    backdropClose?: boolean;
    escapeClose?: boolean;
  }

  let {
    children,
    isOpen = $bindable(false),
    class: className,
    backdropClass: backdropClassName,
    backdropClose = true,
    escapeClose = true,
  }: Props = $props();

  const closeModal = () => (isOpen = false);

  const handleEscapeKey = (event: KeyboardEvent) => {
    if (!escapeClose) return;
    if (event.key === "Escape") {
      event.preventDefault();
      closeModal();
    }
  };
</script>

{#if isOpen}
  <Portal target="body">
    <div
      class="fixed top-0 left-0 z-40 h-screen w-screen"
      transition:fade={{ duration: 100 }}
      onkeydown={handleEscapeKey}
      role="button"
      tabindex="0"
    >
      <button
        class={mergeClasses(
          "z-10 h-full w-full bg-black/30",
          backdropClassName,
        )}
        onclick={backdropClose ? closeModal : undefined}
        aria-label="Close modal"
      >
      </button>

      <div
        class={mergeClasses(
          "absolute top-1/2 left-1/2 z-20 -translate-x-1/2 -translate-y-1/2",
          "max-h-[90dvh] w-[90dvw] max-w-lg rounded-box bg-base-100 p-4 shadow-xl",
          "overflow-hidden",
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
