<script lang="ts">
  import { fade } from "svelte/transition";
  import { mergeClasses } from "$lib/helpers/mergeClasses";
  import type { ClassValue } from "$lib/helpers/mergeClasses";

  interface Props {
    children?: any;
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
    if (event.key === "Escape") closeModal();
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
  <div
    class="z-40 w-screen h-screen fixed top-0 left-0"
    transition:fade={{ duration: 200 }}
  >
    <button
      class={mergeClasses(
        "w-full h-full z-10 bg-black/30",
        backdropClassName,
      )}
      onclick={backdropClose ? closeModal : undefined}
      aria-label="Close modal"
    >
    </button>

    <div
      class={mergeClasses(
        "absolute z-20 top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2",
        "bg-base-100 rounded-box p-4 w-[90dvw] max-w-lg max-h-[90dvh] shadow-xl",
        className,
      )}
    >
      {#if children}
        {@render children()}
      {/if}
    </div>
  </div>
{/if}
