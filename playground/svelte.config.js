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
      precompress: false,
    }),

    paths: {
      assets: "http://REPLACEME", // https://github.com/sveltejs/kit/issues/9569
      relative: true,
    },

    router: {
      type: "hash",
    },
  },
};

export default config;
