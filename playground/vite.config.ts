import { sveltekit } from "@sveltejs/kit/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import { viteStaticCopy } from "vite-plugin-static-copy";

export default defineConfig({
  plugins: [
    tailwindcss(),
    sveltekit(),
    viteStaticCopy({
      targets: [
        {
          src: "node_modules/web-tree-sitter/tree-sitter.wasm",
          dest: "_app/_cconv",
        },
        {
          src: "node_modules/curlconverter/dist/tree-sitter-bash.wasm",
          dest: "_app/_cconv",
        },
      ],
    }),
  ],
  server: {
    host: "0.0.0.0",
  },
  optimizeDeps: {
    esbuildOptions: {
      target: "esnext",
    },
  },
  build: {
    target: "ES2022",
  },
});
