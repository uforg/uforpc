import MiniSearch from "minisearch";

import { getCurrentHost } from "./helpers/getCurrentHost.ts";
import { getMarkdownTitle } from "./helpers/getMarkdownTitle.ts";
import { markdownToText } from "./helpers/markdownToText.ts";
import { slugify } from "./helpers/slugify.ts";
import { transpileUrpcToJson } from "./urpc.ts";
import type { Schema } from "./urpcTypes.ts";

type SearchItem = {
  id: number;
  kind: "doc" | "type" | "proc" | "stream";
  name: string;
  slug: string;
  doc: string;
};

export const miniSearch = new MiniSearch({
  fields: ["kind", "name", "doc"],
  storeFields: ["kind", "name", "slug", "doc"],
  searchOptions: {
    boost: { title: 2 },
    fuzzy: 0.2,
    prefix: true,
  },
  tokenize: (text: string, _?: string): string[] => {
    const tokens: string[] = [];

    // First split by spaces
    const spaceTokens = text.split(" ");
    tokens.push(...spaceTokens);

    // Then split each space token by uppercase letters
    for (const token of spaceTokens) {
      const upperCaseTokens = token.split(/(?=[A-Z])/);
      tokens.push(...upperCaseTokens);
    }

    return tokens;
  },
});

export interface Header {
  key: string;
  value: string;
}

export interface Store {
  loaded: boolean;
  baseUrl: string;
  headers: Header[];
  urpcSchema: string;
  jsonSchema: Schema;
}

export const store: Store = $state({
  loaded: false,
  baseUrl: "",
  headers: [],
  urpcSchema: "version 1",
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
export const loadStore = async () => {
  // Prioritize the config stored in the browser's local storage
  await loadDefaultConfig();

  const baseUrl = localStorage.getItem("baseUrl");
  if (baseUrl) {
    store.baseUrl = baseUrl;
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
  localStorage.setItem("baseUrl", store.baseUrl);
  localStorage.setItem("headers", JSON.stringify(store.headers));
};

/**
 * Add or update a header in the store.headers array.
 *
 * @param key The key of the header to add or update.
 * @param value The value of the header to add or update.
 */
export const setHeader = (key: string, value: string) => {
  const currHeaders = getHeadersObject();
  currHeaders.set(key, value);

  const newHeaders: Header[] = [];
  for (const header of currHeaders.entries()) {
    newHeaders.push({ key: header[0], value: header[1] });
  }

  store.headers = newHeaders;
};

/**
 * Converts the headers array to a Headers object for use in fetch requests
 * it will set the "Content-Type" header to "application/json" by default
 * if that header is not present in the store.headers array.
 *
 * @returns A Headers object
 */
export const getHeadersObject = (): Headers => {
  const headers = new Headers();
  headers.set("Content-Type", "application/json");

  for (const header of store.headers) {
    if (header.key.trim()) headers.set(header.key, header.value);
  }

  return headers;
};

/**
 * Loads the default configuration from the static/config.json file.
 */
export const loadDefaultConfig = async () => {
  const response = await fetch("./config.json");
  if (!response.ok) {
    console.error("Failed to fetch default config");
    return;
  }
  const config = await response.json();

  if (typeof config.baseUrl === "string" && config.baseUrl.trim() !== "") {
    store.baseUrl = config.baseUrl;
  } else {
    store.baseUrl = `${getCurrentHost()}/api/v1/urpc`;
  }

  if (Array.isArray(config.headers)) {
    store.headers = config.headers;
  }
};

/**
 * Transpiles the current URPC schema to JSON format and updates the `jsonSchema` store.
 *
 * This asynchronous function takes the current value of `urpcSchema`, transpiles it to JSON
 * using the `transpileUrpcToJson` utility, and then updates the `jsonSchema` store with the result.
 */
export const loadJsonSchemaFromCurrentUrpcSchema = async () => {
  store.jsonSchema = await transpileUrpcToJson(store.urpcSchema);
  await indexSearchItems();
};

/**
 * Indexes the search items for the current URPC JSON schema.
 */
const indexSearchItems = async () => {
  const searchItems = await Promise.all(
    store.jsonSchema.nodes.map(async (node, index) => {
      let name = "";
      let doc = "";

      if (node.kind === "doc") {
        name = getMarkdownTitle(node.content);
        doc = node.content;
      } else {
        name = node.name;
        doc = node.doc ?? "";
      }

      const item: SearchItem = {
        id: index,
        kind: node.kind,
        name,
        doc,
        slug: slugify(`${node.kind}-${name}`),
      };

      item.doc = await markdownToText(item.doc);

      return item;
    }),
  );

  miniSearch.removeAll();
  miniSearch.addAll(searchItems);
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
