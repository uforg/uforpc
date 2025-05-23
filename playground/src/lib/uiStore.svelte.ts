import type { Action } from "svelte/action";

export interface UiStoreDimensions {
  element: HTMLElement | null;
  size: {
    clientWidth: number;
    clientHeight: number;
    offsetWidth: number;
    offsetHeight: number;
  };
  scroll: {
    left: number;
    top: number;
    isTopScrolled: boolean;
    isLeftScrolled: boolean;
  };
  parentOffset: {
    top: number;
    right: number;
    bottom: number;
    left: number;
  };
  viewportOffset: {
    top: number;
    right: number;
    bottom: number;
    left: number;
  };
  style: {
    width: number;
    height: number;
    marginTop: number;
    marginRight: number;
    marginBottom: number;
    marginLeft: number;
    paddingTop: number;
    paddingRight: number;
    paddingBottom: number;
    paddingLeft: number;
    borderTop: number;
    borderRight: number;
    borderBottom: number;
    borderLeft: number;
  };
}

const defaultUiStoreDimensions: UiStoreDimensions = {
  element: null,
  size: {
    clientWidth: 0,
    clientHeight: 0,
    offsetWidth: 0,
    offsetHeight: 0,
  },
  scroll: {
    left: 0,
    top: 0,
    isTopScrolled: false,
    isLeftScrolled: false,
  },
  parentOffset: {
    top: 0,
    right: 0,
    bottom: 0,
    left: 0,
  },
  viewportOffset: {
    top: 0,
    right: 0,
    bottom: 0,
    left: 0,
  },
  style: {
    width: 0,
    height: 0,
    marginTop: 0,
    marginRight: 0,
    marginBottom: 0,
    marginLeft: 0,
    paddingTop: 0,
    paddingRight: 0,
    paddingBottom: 0,
    paddingLeft: 0,
    borderTop: 0,
    borderRight: 0,
    borderBottom: 0,
    borderLeft: 0,
  },
};

export type Theme = "system" | "light" | "dark";
export type OsTheme = "light" | "dark";

export interface UiStore {
  loaded: boolean;
  theme: Theme;
  osTheme: OsTheme;
  codeSnippetsOpen: boolean;
  codeSnippetsLang: string;
  app: UiStoreDimensions;
  aside: UiStoreDimensions;
  contentWrapper: UiStoreDimensions;
  header: UiStoreDimensions;
  main: UiStoreDimensions;
}

const localStorageKeys = {
  codeSnippetsOpen: "codeSnippetsOpen",
  codeSnippetsLang: "codeSnippetsLang",
};

export const uiStore = $state<UiStore>({
  loaded: false,
  theme: "system",
  osTheme: "dark",
  codeSnippetsOpen: false,
  codeSnippetsLang: "curl",
  app: { ...defaultUiStoreDimensions },
  aside: { ...defaultUiStoreDimensions },
  contentWrapper: { ...defaultUiStoreDimensions },
  header: { ...defaultUiStoreDimensions },
  main: { ...defaultUiStoreDimensions },
});

$effect.root(() => {
  $effect(() => {
    if (!uiStore.loaded) return;
    saveUiStore();
  });
});

/**
 * Loads the store from the browser's local storage.
 *
 * Should be called only once at the start of the app.
 */
export const loadUiStore = () => {
  if (globalThis.matchMedia?.("(prefers-color-scheme: dark)").matches) {
    uiStore.osTheme = "dark";
  } else {
    uiStore.osTheme = "light";
  }

  // Read more at /static/theme-helper.js
  // biome-ignore lint/suspicious/noExplicitAny: it's a global function
  const theme = (globalThis as any).getTheme();
  uiStore.theme = theme || "system";

  // Load code snippets open state from local storage
  const codeSnippetsOpen = globalThis.localStorage.getItem(
    localStorageKeys.codeSnippetsOpen,
  );
  uiStore.codeSnippetsOpen = codeSnippetsOpen
    ? codeSnippetsOpen === "true"
    : true;

  // Load code snippets lang from local storage
  const codeSnippetsLang = globalThis.localStorage.getItem(
    localStorageKeys.codeSnippetsLang,
  );
  uiStore.codeSnippetsLang = codeSnippetsLang ?? "Curl";

  uiStore.loaded = true;
};

/**
 * Saves the store to the browser's local storage.
 *
 * Should be called when the store is updated.
 */
export const saveUiStore = () => {
  // Read more at /static/theme-helper.js
  // biome-ignore lint/suspicious/noExplicitAny: it's a global function
  (globalThis as any).setTheme(uiStore.theme);

  // Save code snippets open state to local storage
  globalThis.localStorage.setItem(
    localStorageKeys.codeSnippetsOpen,
    uiStore.codeSnippetsOpen.toString(),
  );

  // Save code snippets lang to local storage
  globalThis.localStorage.setItem(
    localStorageKeys.codeSnippetsLang,
    uiStore.codeSnippetsLang,
  );
};

//////////////////////
// Helper functions //
//////////////////////

/**
 * Finds all scrollable ancestor elements of a given HTML element
 *
 * @param {HTMLElement} el - The HTML element to find scrollable ancestors for
 * @returns {(Window | HTMLElement)[]} An array of scrollable ancestors, including the window
 */
function getScrollableAncestors(el: HTMLElement): (Window | HTMLElement)[] {
  const hosts: (Window | HTMLElement)[] = [window];
  let parent = el.parentElement;

  while (parent) {
    const style = getComputedStyle(parent);
    const overflowY = style.overflowY;
    const overflowX = style.overflowX;
    const canScrollY =
      (overflowY === "auto" || overflowY === "scroll") &&
      parent.scrollHeight > parent.clientHeight;
    const canScrollX =
      (overflowX === "auto" || overflowX === "scroll") &&
      parent.scrollWidth > parent.clientWidth;
    if (canScrollY || canScrollX) hosts.push(parent);
    parent = parent.parentElement;
  }

  return hosts;
}

/**
 * Svelte action that tracks and reports element dimensions and position changes
 *
 * This action monitors an element's size, position, scroll state, and style properties,
 * dispatching a custom event whenever these dimensions change due to resizing, scrolling,
 * or other layout changes.
 *
 * @param {HTMLElement} node - The HTML element to track
 * @returns {object} Action lifecycle methods
 */
export const dimensionschangeAction: Action<
  HTMLElement,
  undefined,
  { ondimensionschange: (e: CustomEvent<UiStoreDimensions>) => void }
> = (node) => {
  let scrollHosts: (Window | HTMLElement)[] = [];
  let ticking = false;

  const dispatchEvent = () => {
    const nodeRect = node.getBoundingClientRect();

    const clientWidth = node.clientWidth;
    const clientHeight = node.clientHeight;
    const offsetWidth = node.offsetWidth;
    const offsetHeight = node.offsetHeight;

    const scrollLeft = node.scrollLeft;
    const scrollTop = node.scrollTop;

    let parentOffset = { top: 0, left: 0, bottom: 0, right: 0 };
    const parent = node.parentElement;
    if (parent) {
      const parentRect = parent.getBoundingClientRect();
      parentOffset = {
        top: nodeRect.top - parentRect.top,
        left: nodeRect.left - parentRect.left,
        bottom: parentRect.bottom - nodeRect.bottom,
        right: parentRect.right - nodeRect.right,
      };
    }

    const viewportOffset = {
      top: nodeRect.top,
      left: nodeRect.left,
      bottom: globalThis.innerHeight - nodeRect.bottom,
      right: globalThis.innerWidth - nodeRect.right,
    };

    const style = globalThis.getComputedStyle(node);
    const width = Number.parseFloat(style.width);
    const height = Number.parseFloat(style.height);
    const marginTop = Number.parseFloat(style.marginTop);
    const marginRight = Number.parseFloat(style.marginRight);
    const marginBottom = Number.parseFloat(style.marginBottom);
    const marginLeft = Number.parseFloat(style.marginLeft);
    const paddingTop = Number.parseFloat(style.paddingTop);
    const paddingRight = Number.parseFloat(style.paddingRight);
    const paddingBottom = Number.parseFloat(style.paddingBottom);
    const paddingLeft = Number.parseFloat(style.paddingLeft);
    const borderTop = Number.parseFloat(style.borderTopWidth);
    const borderRight = Number.parseFloat(style.borderRightWidth);
    const borderBottom = Number.parseFloat(style.borderBottomWidth);
    const borderLeft = Number.parseFloat(style.borderLeftWidth);

    node.dispatchEvent(
      new CustomEvent<UiStoreDimensions>("dimensionschange", {
        detail: {
          element: node,
          size: {
            clientWidth,
            clientHeight,
            offsetWidth,
            offsetHeight,
          },
          scroll: {
            left: scrollLeft,
            top: scrollTop,
            isTopScrolled: scrollTop > 0,
            isLeftScrolled: scrollLeft > 0,
          },
          parentOffset,
          viewportOffset,
          style: {
            width,
            height,
            marginTop,
            marginRight,
            marginBottom,
            marginLeft,
            paddingTop,
            paddingRight,
            paddingBottom,
            paddingLeft,
            borderTop,
            borderRight,
            borderBottom,
            borderLeft,
          },
        },
      }),
    );
  };

  function throttledDispatchEvent() {
    if (!ticking) {
      ticking = true;
      requestAnimationFrame(() => {
        dispatchEvent();
        ticking = false;
      });
    }
  }

  const observer = new ResizeObserver((entries) => {
    if (entries.length !== 1) return;
    throttledDispatchEvent();
  });

  function attachScrollListeners() {
    scrollHosts = getScrollableAncestors(node);
    for (const host of scrollHosts) {
      host.addEventListener("scroll", throttledDispatchEvent, {
        passive: true,
      });
    }

    node.addEventListener("scroll", throttledDispatchEvent);
  }
  function detachScrollListeners() {
    for (const host of scrollHosts) {
      host.removeEventListener("scroll", throttledDispatchEvent);
    }
    scrollHosts = [];

    node.removeEventListener("scroll", throttledDispatchEvent);
  }

  $effect(() => {
    throttledDispatchEvent();

    observer.observe(node);
    globalThis.addEventListener("resize", throttledDispatchEvent);
    attachScrollListeners();

    return () => {
      observer.disconnect();
      globalThis.removeEventListener("resize", throttledDispatchEvent);
      detachScrollListeners();
    };
  });
};
