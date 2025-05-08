import adapter from "@sveltejs/adapter-static";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),

  kit: {
    adapter: adapter({
      pages: "build",
      assets: "build",
      strict: true,
      precompress: true,
    }),

    paths: {
      relative: true,
    },

    router: {
      type: "hash",
    },
  },
};

export default config;
