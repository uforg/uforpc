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
          src: "../urpc/dist/urpc.wasm",
          dest: "app/urpc",
        },
        {
          src: "../urpc/dist/wasm_exec.js",
          dest: "app/urpc",
        },
        {
          src: "node_modules/web-tree-sitter/tree-sitter.wasm",
          dest: "app/cconv",
        },
        {
          src: "node_modules/curlconverter/dist/tree-sitter-bash.wasm",
          dest: "app/cconv",
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
