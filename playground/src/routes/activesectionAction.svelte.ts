import type { Action } from "svelte/action";

/**
 * This action tracks which section is closest to the top of the viewport
 * and updates the hash accordingly
 *
 * Usage: <div use:activesectionAction onactivesection={...}></div>
 *
 * @param node The scrollable container that contains sections
 */
export const activesectionAction: Action<
  HTMLElement,
  undefined,
  { onactivesection: (e: CustomEvent<string>) => void }
> = (node) => {
  const scrollOffset = 150; // Cherry picked value

  const handleScroll = () => {
    const sections = Array.from(node.querySelectorAll("section[id]"));
    if (!sections.length) return;

    let activeSection: Element | undefined;
    for (const section of sections) {
      const rect = section.getBoundingClientRect();
      const topRel = rect.top;
      const bottomRel = rect.bottom;

      if (topRel <= scrollOffset && bottomRel > scrollOffset) {
        activeSection = section;
        break;
      }
    }

    node.dispatchEvent(
      new CustomEvent("activesection", { detail: activeSection?.id ?? "" }),
    );
  };

  // Debounce scroll events
  let timeout: ReturnType<typeof setTimeout>;
  const debouncedScroll = () => {
    clearTimeout(timeout);
    timeout = setTimeout(handleScroll, 100);
  };

  $effect(() => {
    node.addEventListener("scroll", debouncedScroll);
    return () => {
      clearTimeout(timeout);
      node.removeEventListener("scroll", debouncedScroll);
    };
  });
};
