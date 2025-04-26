import type { Action } from "svelte/action";

type IsscrolledAction = Action<
  HTMLElement,
  undefined,
  { onisscrolled: (e: CustomEvent<boolean>) => void }
>;

/**
 * This action adds a isscrolled event to the node and emits
 * a boolean value indicating if the node is scrolled
 *
 * Usage: <div use:isscrolled onisscrolled={...}></div>
 *
 * @param node The node to add the event to
 */
export const isscrolledAction: IsscrolledAction = (node: HTMLElement) => {
  const handleScroll = () => {
    const isScrolled = node.scrollTop > 0;

    node.dispatchEvent(
      new CustomEvent<boolean>("isscrolled", { detail: isScrolled }),
    );
  };

  $effect(() => {
    node.addEventListener("scroll", handleScroll);
    return () => {
      node.removeEventListener("scroll", handleScroll);
    };
  });
};
