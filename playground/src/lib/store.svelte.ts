import { transpileUrpcToJson } from "../lib/urpc.ts";
import type { Schema } from "../lib/urpcTypes.ts";

export type Theme = "system" | "light" | "dark";

export interface Store {
  loaded: boolean;
  theme: Theme;
  endpoint: string;
  headers: Record<string, string>;
  urpcSchema: string;
  jsonSchema: Schema;
}

export const store: Store = $state({
  loaded: false,
  theme: "system",
  endpoint: "",
  headers: {},
  urpcSchema: `version 1`,
  jsonSchema: { version: 1, nodes: [] },
});

$effect.root(() => {
  $effect(() => {
    if (!store.loaded) return;
    saveStore();
  });
});

/**
 * Loads the store from the browser's local storage.
 *
 * Should be called only once at the start of the app.
 */
export const loadStore = () => {
  // Read more at /static/theme-helper.js
  // deno-lint-ignore no-explicit-any
  const theme = (globalThis as any).getTheme();
  store.theme = theme || "system";

  const endpoint = localStorage.getItem("endpoint");
  if (endpoint) {
    store.endpoint = endpoint;
  }

  const headers = localStorage.getItem("headers");
  if (headers) {
    store.headers = JSON.parse(headers);
  }

  store.loaded = true;
};

/**
 * Saves the store to the browser's local storage.
 *
 * Should be called when the store is updated.
 */
export const saveStore = () => {
  // Read more at /static/theme-helper.js
  // deno-lint-ignore no-explicit-any
  (globalThis as any).setTheme(store.theme);
  localStorage.setItem("endpoint", store.endpoint);
  localStorage.setItem("headers", JSON.stringify(store.headers));
};

/**
 * Fetches and loads an URPC schema from a specified URL.
 *
 * This function attempts to retrieve a schema from the given URL and, if successful,
 * updates the `urpcSchema` store with the fetched content. If the fetch fails,
 * an error is logged to the console.
 *
 * @param url The URL from which to fetch the URPC schema.
 * @throws Logs an error to the console if the fetch operation fails.
 */
export const loadUrpcSchemaFromUrl = async (url: string) => {
  const response = await fetch(url);
  if (!response.ok) {
    console.error(`Failed to fetch schema from ${url}`);
    return;
  }

  const sch = await response.text();
  store.urpcSchema = sch;
};

/**
 * Updates the `urpcSchema` store with a provided URPC schema string.
 *
 * This function directly sets the `urpcSchema` store to the provided schema string,
 * allowing for immediate updates to the schema without fetching from a URL.
 *
 * @param sch The URPC schema string to be loaded into the store.
 */
export const loadUrpcSchemaFromString = (sch: string) => {
  store.urpcSchema = sch;
};

/**
 * Transpiles the current URPC schema to JSON format and updates the `jsonSchema` store.
 *
 * This asynchronous function takes the current value of `urpcSchema`, transpiles it to JSON
 * using the `transpileUrpcToJson` utility, and then updates the `jsonSchema` store with the result.
 */
export const loadJsonSchemaFromCurrentUrpcSchema = async () => {
  store.jsonSchema = await transpileUrpcToJson(store.urpcSchema);
};

/**
 * Fetches an URPC schema from a URL, loads it, and then transpiles it to JSON.
 *
 * This function combines the operations of `loadUrpcSchemaFromUrl` and
 * `loadJsonSchemaFromCurrentUrpcSchema`. It first fetches and loads the URPC schema
 * from the specified URL, then transpiles that schema to JSON, updating both
 * the `urpcSchema` and `jsonSchema` stores in the process.
 *
 * @param url The URL from which to fetch the URPC schema.
 */
export const loadJsonSchemaFromUrpcSchemaUrl = async (url: string) => {
  await loadUrpcSchemaFromUrl(url);
  await loadJsonSchemaFromCurrentUrpcSchema();
};

/**
 * Loads an URPC schema from a string and transpiles it to JSON.
 *
 * This function takes an URPC schema as a string, loads it into the `urpcSchema` store,
 * and then transpiles it to JSON, updating both the `urpcSchema` and `jsonSchema` stores.
 * It's useful for processing schemas that are already available as strings without needing
 * to fetch from a URL.
 *
 * @param sch The URPC schema string to load and transpile.
 */
export const loadJsonSchemaFromUrpcSchemaString = async (sch: string) => {
  loadUrpcSchemaFromString(sch);
  await loadJsonSchemaFromCurrentUrpcSchema();
};
