# UFO RPC Docs

This is the documentation website of UFO RPC. Read this entire document and only pay attention to the sections related to your tasks.

## Commands

- `npm run dev`: Starts a local development server with live reloading (this blocks terminal until the process is killed)
- `npm run build`: Builds the app for production
- `npm run lint`: Lints the codebase and run typescript checks
- `npm run fmt`: Formats the codebase using prettier

Always execute the `fmt`, `lint` and `build` commands when you finish your changes, and ensure there are no errors. If there are errors, fix them until there are none.

## Stack

### Svelte and SvelteKit v5

The project is built using SvelteKit and configured to generate a static site. You can read the SvelteKit configuration in `svelte.config.js` and `vite.config.ts` files.

The Svelte version used is 5, and the syntax is different from previous versions. You can read this documentation for more information: https://svelte.dev/llms-small.txt

Always use typescript and svelte runes.

### MDsveX

The project includes support for `.svx` files using MDsveX, you can write your docs using markdown and svelte components.

You can read the configuration in `svelte.config.js`.

And you can read more about `MDsveX` here: https://context7.com/websites/mdsvex_pngwn_io/llms.txt

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

## Writing Docs

All the docs should be written on a file based routing system inside `src/docs` directory, if you need to create a new doc, just create a new `.md` or `.svx` file inside that directory, and it will be automatically available in the sidebar.

Before writing a new doc, read 3 existing docs to understand the writing style and format.

### How docs are rendered?

Every time you edit a doc, the script `scripts/gendocs.ts` generates the file `src/lib/docs.gen.ts` where the docs metadata is stored, then the SvelteKit files inside `src/routes/docs/[...slug]` directory use that metadata to import and render the docs.

The files inside `src/routes/docs` also uses the generated docs metadata to show a navigation sidebar with the docs structure.

All the helpers used to work with the `src/lib/docs.gen.ts` file are stored in `src/lib/docs.ts` file.
