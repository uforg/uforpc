// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";

// https://astro.build/config
export default defineConfig({
  site: "https://uforpc.uforg.dev",

  integrations: [
    starlight({
      title: "UFO RPC Docs",
      social: [
        {
          icon: "github",
          label: "GitHub",
          href: "https://github.com/uforg/uforpc",
        },
      ],
      sidebar: [
        {
          label: "Docs",
          autogenerate: { directory: "docs" },
        },
      ],
    }),
  ],
});
