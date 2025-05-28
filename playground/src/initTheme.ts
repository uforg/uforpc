import { browser } from "$app/environment";

import { initTheme } from "$lib/uiStore.svelte";

/**
 * This file is used to initialize the theme of the app. It should be loaded
 * before anything else in the app.html file to prevent flickering.
 */

(() => {
  if (!browser) return;
  initTheme();
})();
