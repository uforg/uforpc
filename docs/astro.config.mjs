// @ts-check
import starlight from "@astrojs/starlight";
import svelte from "@astrojs/svelte";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "astro/config";
import { execSync } from "child_process";
import fse from "fs-extra";
import { bundledLanguages } from "shiki";

const syntaxUrl =
  "https://cdn.jsdelivr.net/gh/uforg/uforpc-vscode@0.1.7/syntaxes/urpc.tmLanguage.json";
const syntaxOutputFile = "urpc.tmLanguage.json";

// https://docs.astro.build/en/reference/configuration-reference/#markdown-options
// https://shiki.style/guide/load-lang
function getShikiLangs() {
  if (!fse.existsSync(syntaxOutputFile)) {
    execSync(`wget -q -O ${syntaxOutputFile} ${syntaxUrl}`);
  }

  const urpcLang = fse.readJSONSync(syntaxOutputFile);
  urpcLang.name = "urpc";

  return [...Object.keys(bundledLanguages), urpcLang];
}

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

  markdown: {
    syntaxHighlight: "shiki",
    shikiConfig: {
      langs: getShikiLangs(),
    },
  },

  vite: {
    plugins: [tailwindcss()],
  },
});
