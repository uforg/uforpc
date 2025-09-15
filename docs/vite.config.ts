import { sveltekit } from "@sveltejs/kit/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import { run } from "vite-plugin-run";

export default defineConfig({
  plugins: [
    tailwindcss(),
    sveltekit(),
    run([
      {
        name: "gendocs",
        run: ["npm", "run", "gendocs"],
        pattern: ["src/docs/**/*.md", "src/docs/**/*.svx"],
      },
    ]),
  ],
});
