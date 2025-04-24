/**
 * This action adds a clickoutside event to the node
 *
 * Usage: <div use:clickOutside on:clickoutside={...}></div>
 *
 * @param node The node to add the event to
 */
export function clickOutside(node: HTMLElement) {
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
}
