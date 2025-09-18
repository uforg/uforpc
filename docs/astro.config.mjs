// @ts-check
import starlight from "@astrojs/starlight";
import svelte from "@astrojs/svelte";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "astro/config";
import { run } from "vite-plugin-run";

const syntaxUrl =
  "https://cdn.jsdelivr.net/gh/uforg/uforpc-vscode@0.1.7/syntaxes/urpc.tmLanguage.json";
const syntaxOutputFile = "urpc.tmLanguage.json";

// https://astro.build/config
export default defineConfig({
  integrations: [
    svelte(),
    starlight({
      title: "Docs with Tailwind",
      social: [
        {
          icon: "github",
          label: "GitHub",
          href: "https://github.com/withastro/starlight",
        },
      ],
      sidebar: [
        {
          label: "Guides",
          items: [
            // Each item here is one entry in the navigation menu.
            { label: "Example Guide", slug: "guides/example" },
          ],
        },
        {
          label: "Reference",
          autogenerate: { directory: "reference" },
        },
      ],
      customCss: ["./src/styles/global.css"],
    }),
  ],

  vite: {
    plugins: [
      tailwindcss(),
      run([
        {
          name: "download urpc syntax",
          run: ["wget", "-q", syntaxUrl, "-O", syntaxOutputFile],
          pattern: ["astro.config.mjs"],
          build: true,
          startup: true,
        },
      ]),
    ],
  },
});
