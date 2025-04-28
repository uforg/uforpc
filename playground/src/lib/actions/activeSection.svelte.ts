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
  const handleScroll = () => {
    const sections = Array.from(node.querySelectorAll("section[id]"));
    if (!sections.length) return;

    let closestSection = sections[0];
    let closestDistance = Infinity;
    const scrollOffset = 100;

    for (const section of sections) {
      const distance = Math.abs(
        section.getBoundingClientRect().top - scrollOffset,
      );
      if (distance < closestDistance) {
        closestDistance = distance;
        closestSection = section;
      }
    }

    node.dispatchEvent(
      new CustomEvent("activesection", { detail: closestSection.id }),
    );
  };

  // Debounce scroll events
  let timeout: number;
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
