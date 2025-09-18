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
      favicon: "/icon.png",
      title: "UFO RPC",
      description: "Modern RPC framework that puts developer experience first.",
      social: [
        {
          icon: "github",
          label: "GitHub",
          href: "https://github.com/uforg/uforpc",
        },
        {
          icon: "discord",
          label: "Discord",
          href: "https://uforpc.uforg.dev/r/discord",
        },
        {
          icon: "reddit",
          label: "Reddit",
          href: "https://uforpc.uforg.dev/r/reddit",
        },
        {
          icon: "x.com",
          label: "X (Twitter)",
          href: "https://uforpc.uforg.dev/r/twitter",
        },
      ],
      sidebar: [
        {
          label: "Guides",
          autogenerate: { directory: "guides" },
        },
        {
          label: "Reference",
          autogenerate: { directory: "reference" },
        },
      ],
      editLink: {
        baseUrl: "https://github.com/uforg/uforpc/tree/main/docs/",
      },
      components: {
        SiteTitle: "./src/components/SiteTitle.svelte",
      },
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
