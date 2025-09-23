# UFO RPC Playground

This is the playground webapp of UFO RPC. Read this entire document and only pay attention to the sections related to your tasks.

## Commands

Read the `package.json` file to see all available commands, the most important ones are:

- `npm run dev`: Starts a local development server with live reloading (this blocks terminal until the process is killed)
- `npm run build`: Builds the app for production
- `npm run lint`: Lints the codebase and run typescript checks
- `npm run fmt`: Formats the codebase using prettier

Always execute the `fmt`, `lint` and `build` commands when you finish your changes, and ensure there are no errors. If there are errors, fix them until there are none, the only exception is when you are only writing documentation, in that case only `fmt` is required.

Before running a command like `npm run {command}`, make sure you are inside `/workspaces/uforpc/playground` directory.

## Stack

### Svelte and SvelteKit v5

The project is built using SvelteKit and configured to generate a static site. You can read the SvelteKit configuration in `svelte.config.js` and `vite.config.ts` files.

The Svelte version used is 5, and the syntax is different from previous versions. You can read this documentation for more information: https://svelte.dev/llms-small.txt

Always use typescript and svelte runes.

Don't use svelte <style></style> tags, instead use TailwindCSS and DaisyUI classes to style the HTML elements.

### TailwindCSS v4

The project styling is done using TailwindCSS v4. You can find the configuration in `src/app.css` file.

There is only one breakpoint configured named `desk`, so, when you want to create styles for mobile dont use any prefix, and for desktop use the `desk:` prefix, forget about all other breakpoints.

### DaisyUI v5

In addition to TailwindCSS the project is using DaisyUI v5, always prefer DaisyUI classes, you can read the documentation here: https://daisyui.com/llms.txt

### Lucide Icons

This project uses the Lucide icons from `@lucide/svelte` package, this is how you use the icons:

```svelte
<script lang="ts">
  import { Code, Github, Zap } from "@lucide/svelte";
</script>

<Zap class="size-4" />
```
