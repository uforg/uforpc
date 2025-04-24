import type { Action } from "svelte/action";

type ClickOutsideAction = Action<
  HTMLElement,
  undefined,
  { onclickoutside: (e: CustomEvent) => void }
>;

/**
 * This action adds a clickoutside event to the node
 *
 * Usage: <div use:clickOutside on:clickoutside={...}></div>
 *
 * @param node The node to add the event to
 */
export const clickOutside: ClickOutsideAction = (node: HTMLElement) => {
  function handleClick(e: MouseEvent) {
    if (!e.target) return;
    if (!node.contains(e.target as Node)) {
      node.dispatchEvent(new CustomEvent("clickoutside"));
    }
  }

  $effect(() => {
    document.addEventListener("click", handleClick);

    return () => {
      document.removeEventListener("click", handleClick);
    };
  });
};
