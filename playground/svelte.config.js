import adapter from "@sveltejs/adapter-static";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";

// https://github.com/sveltejs/kit/issues/9569
const replaceAssets = process.env.REPLACE_ASSETS === "true";

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),

  kit: {
    adapter: adapter({
      pages: "build",
      assets: "build",
      strict: true,
      precompress: false,
    }),

    paths: {
      relative: true,
      assets: replaceAssets ? "http://REPLACEME" : undefined,
    },

    router: {
      type: "hash",
    },
  },
};

export default config;
