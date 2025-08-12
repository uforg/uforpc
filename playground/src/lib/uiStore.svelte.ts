import type { Action } from "svelte/action";

import { debounce } from "./helpers/debounce";
import type { CodegenGenerator } from "./urpc";

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

const matchMediaColor = globalThis.matchMedia?.("(prefers-color-scheme: dark)");

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

export type Theme = "light" | "dark";

export interface UiStore {
  loaded: boolean;
  isMobile: boolean;
  theme: Theme;
  codeSnippetsTab: "sdk" | "curl";
  codeSnippetsCurlLang: string;
  codeSnippetsSdkLang: CodegenGenerator;
  codeSnippetsSdkStep: "download" | "setup" | "usage" | "";
  codeSnippetsSdkDartPackageName: string;
  codeSnippetsSdkGolangPackageName: string;
  asideOpen: boolean;
  asideSearchOpen: boolean;
  asideSearchQuery: string;
  asideHideDocs: boolean;
  asideHideTypes: boolean;
  asideHideProcs: boolean;
  asideHideStreams: boolean;
  app: UiStoreDimensions;
  aside: UiStoreDimensions;
  contentWrapper: UiStoreDimensions;
  header: UiStoreDimensions;
  main: UiStoreDimensions;
}

const localStorageKeys = {
  theme: "theme",
  codeSnippetsTab: "codeSnippetsTab",
  codeSnippetsCurlLang: "codeSnippetsCurlLang",
  codeSnippetsSdkLang: "codeSnippetsSdkLang",
  codeSnippetsSdkStep: "codeSnippetsSdkStep",
  codeSnippetsSdkDartPackageName: "codeSnippetsSdkDartPackageName",
  codeSnippetsSdkGolangPackageName: "codeSnippetsSdkGolangPackageName",
  asideSearchOpen: "asideSearchOpen",
  asideSearchQuery: "asideSearchQuery",
  asideHideDocs: "asideHideDocs",
  asideHideTypes: "asideHideTypes",
  asideHideProcs: "asideHideProcs",
  asideHideStreams: "asideHideStreams",
};

export const uiStore = $state<UiStore>({
  loaded: false,
  isMobile: false,
  theme: "dark",
  codeSnippetsTab: "sdk",
  codeSnippetsCurlLang: "Curl",
  codeSnippetsSdkLang: "typescript-client",
  codeSnippetsSdkStep: "download",
  codeSnippetsSdkDartPackageName: "uforpc",
  codeSnippetsSdkGolangPackageName: "uforpc",
  asideOpen: false,
  asideSearchOpen: false,
  asideSearchQuery: "",
  asideHideDocs: false,
  asideHideTypes: true,
  asideHideProcs: false,
  asideHideStreams: false,
  app: { ...defaultUiStoreDimensions },
  aside: { ...defaultUiStoreDimensions },
  contentWrapper: { ...defaultUiStoreDimensions },
  header: { ...defaultUiStoreDimensions },
  main: { ...defaultUiStoreDimensions },
});

$effect.root(() => {
  // Effect to check if the screen is mobile (even on resize) with debounce
  $effect(() => {
    const calcIsMobile = debounce(() => {
      const mobileThreshold = 1200;
      uiStore.isMobile = globalThis.innerWidth < mobileThreshold;
    }, 100);

    calcIsMobile();
    globalThis.addEventListener("resize", calcIsMobile);
    return () => {
      globalThis.removeEventListener("resize", calcIsMobile);
    };
  });

  // Effect to save the store to the browser's local storage
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
  /**
   * IMPORTANT:
   * The theme should be loaded before anything else in
   * the app.html file.
   */

  // Load code snippets tab from local storage
  const codeSnippetsTab = globalThis.localStorage.getItem(
    localStorageKeys.codeSnippetsTab,
  );
  uiStore.codeSnippetsTab = codeSnippetsTab === "curl" ? "curl" : "sdk";

  // Load code snippets curl lang from local storage
  const codeSnippetsCurlLang = globalThis.localStorage.getItem(
    localStorageKeys.codeSnippetsCurlLang,
  );
  uiStore.codeSnippetsCurlLang = codeSnippetsCurlLang ?? "Curl";

  // Load code snippets sdk lang from local storage
  const codeSnippetsSdkLang = globalThis.localStorage.getItem(
    localStorageKeys.codeSnippetsSdkLang,
  );
  uiStore.codeSnippetsSdkLang = (codeSnippetsSdkLang ??
    "typescript-client") as CodegenGenerator;

  // Load code snippets sdk step from local storage
  const codeSnippetsSdkStep = globalThis.localStorage.getItem(
    localStorageKeys.codeSnippetsSdkStep,
  );
  uiStore.codeSnippetsSdkStep = (codeSnippetsSdkStep ?? "download") as
    | "download"
    | "setup"
    | "usage";

  // Load code snippets sdk dart package name from local storage
  const codeSnippetsSdkDartPackageName = globalThis.localStorage.getItem(
    localStorageKeys.codeSnippetsSdkDartPackageName,
  );
  uiStore.codeSnippetsSdkDartPackageName =
    codeSnippetsSdkDartPackageName ?? "uforpc";

  // Load code snippets sdk golang package name from local storage
  const codeSnippetsSdkGolangPackageName = globalThis.localStorage.getItem(
    localStorageKeys.codeSnippetsSdkGolangPackageName,
  );
  uiStore.codeSnippetsSdkGolangPackageName =
    codeSnippetsSdkGolangPackageName ?? "uforpc";

  // Load aside search open state from local storage
  const asideSearchOpen = globalThis.localStorage.getItem(
    localStorageKeys.asideSearchOpen,
  );
  uiStore.asideSearchOpen = asideSearchOpen === "true";

  // Load aside search query from local storage
  const asideSearchQuery = globalThis.localStorage.getItem(
    localStorageKeys.asideSearchQuery,
  );
  uiStore.asideSearchQuery = asideSearchQuery ?? "";

  // Load aside hide docs from local storage
  const asideHideDocs = globalThis.localStorage.getItem(
    localStorageKeys.asideHideDocs,
  );
  uiStore.asideHideDocs = asideHideDocs === "true";

  // Load aside hide types from local storage
  const asideHideTypes = globalThis.localStorage.getItem(
    localStorageKeys.asideHideTypes,
  );
  uiStore.asideHideTypes = asideHideTypes ? asideHideTypes === "true" : true;

  // Load aside hide procs from local storage
  const asideHideProcs = globalThis.localStorage.getItem(
    localStorageKeys.asideHideProcs,
  );
  uiStore.asideHideProcs = asideHideProcs === "true";

  // Load aside hide streams from local storage
  const asideHideStreams = globalThis.localStorage.getItem(
    localStorageKeys.asideHideStreams,
  );
  uiStore.asideHideStreams = asideHideStreams === "true";

  uiStore.loaded = true;
};

/**
 * Saves the store to the browser's local storage.
 *
 * Should be called when the store is updated.
 */
export const saveUiStore = () => {
  // Save theme to local storage
  globalThis.localStorage.setItem(localStorageKeys.theme, uiStore.theme);
  setThemeAttribute(uiStore.theme);

  // Save code snippets curl lang to local storage
  globalThis.localStorage.setItem(
    localStorageKeys.codeSnippetsTab,
    uiStore.codeSnippetsTab,
  );

  // Save code snippets curl lang to local storage
  globalThis.localStorage.setItem(
    localStorageKeys.codeSnippetsCurlLang,
    uiStore.codeSnippetsCurlLang,
  );

  // Save code snippets sdk lang to local storage
  globalThis.localStorage.setItem(
    localStorageKeys.codeSnippetsSdkLang,
    uiStore.codeSnippetsSdkLang,
  );

  // Save code snippets sdk step to local storage
  globalThis.localStorage.setItem(
    localStorageKeys.codeSnippetsSdkStep,
    uiStore.codeSnippetsSdkStep,
  );

  // Save code snippets sdk dart package name to local storage
  globalThis.localStorage.setItem(
    localStorageKeys.codeSnippetsSdkDartPackageName,
    uiStore.codeSnippetsSdkDartPackageName,
  );

  // Save code snippets sdk golang package name to local storage
  globalThis.localStorage.setItem(
    localStorageKeys.codeSnippetsSdkGolangPackageName,
    uiStore.codeSnippetsSdkGolangPackageName,
  );

  // Save aside search open state to local storage
  globalThis.localStorage.setItem(
    localStorageKeys.asideSearchOpen,
    uiStore.asideSearchOpen.toString(),
  );

  // Save aside search query to local storage
  globalThis.localStorage.setItem(
    localStorageKeys.asideSearchQuery,
    uiStore.asideSearchQuery,
  );

  // Save aside hide docs to local storage
  globalThis.localStorage.setItem(
    localStorageKeys.asideHideDocs,
    uiStore.asideHideDocs.toString(),
  );

  // Save aside hide types to local storage
  globalThis.localStorage.setItem(
    localStorageKeys.asideHideTypes,
    uiStore.asideHideTypes.toString(),
  );

  // Save aside hide procs to local storage
  globalThis.localStorage.setItem(
    localStorageKeys.asideHideProcs,
    uiStore.asideHideProcs.toString(),
  );

  // Save aside hide streams to local storage
  globalThis.localStorage.setItem(
    localStorageKeys.asideHideStreams,
    uiStore.asideHideStreams.toString(),
  );
};

/////////////////////
// THEME UTILITIES //
/////////////////////
/**
 * Returns the system theme based on the css color scheme
 */
function getSystemTheme(): Theme {
  if (!matchMediaColor) return "dark";
  return matchMediaColor?.matches ? "dark" : "light";
}

/**
 * Sets the theme attribute of the document element, used by DaisyUI
 */
function setThemeAttribute(theme: Theme) {
  document.documentElement.setAttribute("data-theme", theme);
}

/**
 * Sets the theme stored in the local storage, it falls back to the
 * system theme
 */
export function initTheme() {
  const theme = localStorage.getItem(localStorageKeys.theme);
  if (theme === "light" || theme === "dark") {
    uiStore.theme = theme;
  } else {
    uiStore.theme = getSystemTheme();
  }
  setThemeAttribute(uiStore.theme);

  // Listen for changes in the color scheme to change the theme dinamically
  matchMediaColor?.addEventListener("change", () => {
    uiStore.theme = getSystemTheme();
    setThemeAttribute(uiStore.theme);
  });
}

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
