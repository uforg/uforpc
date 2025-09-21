import { browser } from "$app/environment";
import localforage from "localforage";
import { debounce } from "lodash-es";
import { toast } from "svelte-sonner";

interface CreateAsyncStoreOptions<T extends Record<string, unknown>> {
  initialValue: () => Promise<T>;
  keysToPersist: (keyof T)[];
  storeName?: string;
}

interface AsyncStoreStatus {
  loading: boolean;
  saving: boolean;
}

interface AsyncStoreResult<T extends Record<string, unknown>> {
  store: T;
  status: AsyncStoreStatus;
  ready: () => Promise<void>;
}

/**
 * Creates an asynchronous Svelte store that initializes its value from a
 * provided async function and persists specified keys to IndexedDB using
 * localforage with fallback to localStorage if IndexedDB is unavailable.
 *
 * Important: this function works only with JSON-serializable values.
 *
 * @template T - The type of the store object, which should be a record with string keys and unknown values.
 * @param opts - Configuration options for creating the async store.
 * @param opts.initialValue - An async function that returns the initial value of the store.
 * @param opts.keysToPersist - An array of keys from the store that should be persisted to IndexedDB.
 * @param opts.storeName - An optional name for the store, used to create a unique isolated database instead of the global one.
 * @returns An object containing the Svelte store, its lifecycle (status) state and a promise to wait for initialization.
 *
 * @example
 * ```ts
 * const { store, status, ready } = createAsyncStore({
 *   initialValue: async () => ({ theme: 'light', fontSize: 14 }),
 *   keysToPersist: ['theme'],
 *   storeName: 'userPreferences',
 * });
 * ```
 */

// biome-ignore lint/suspicious/noExplicitAny: the values are dynamic and varied between different stores
export function createAsyncStore<T extends Record<string, any>>(
  opts: CreateAsyncStoreOptions<T>,
): AsyncStoreResult<T> {
  // Initialize Svelte stores
  let store = $state<T>({} as T);
  const status = $state({
    loading: true,
    saving: false,
  });

  // Promise for waiting for initialization
  let resolveReady: () => void;
  const readyPromise = new Promise<void>((resolve) => {
    resolveReady = resolve;
  });

  // Asynchronously manage the store lifecycle
  (async () => {
    // Browser-only check
    if (!browser) {
      // biome-ignore lint/style/noNonNullAssertion: the function is always defined above
      resolveReady!();
      status.loading = false;
      return;
    }

    // Create the localforage database name, it' will be used to isolate
    // different stores between themselves
    let dbName = createGlobalDbNamePrefix();
    if (opts.storeName && opts.storeName.trim() !== "") {
      dbName += `-${opts.storeName.trim()}`;
    }

    // Create localforage database instance
    // https://localforage.github.io/localForage/#multiple-instances-createinstance
    // https://localforage.github.io/localForage/#settings-api-config
    const db = localforage.createInstance({
      name: dbName,
      driver: localforage.INDEXEDDB,
    });

    // Load the initial store value
    try {
      const initialValue = await opts.initialValue();
      for (const key in initialValue) {
        (store[key] as unknown) = initialValue[key];
      }
    } catch (error) {
      toast.error("Failed to load initial store value", {
        description: `Error: ${error}`,
      });
    }

    // Load persisted values from the database
    try {
      const promises = opts.keysToPersist.map(async (keyToPersist) => {
        const value = await db.getItem(keyToPersist as string);
        if (value === null) return;
        (store[keyToPersist] as unknown) = value;
      });
      await Promise.all(promises);
    } catch (error) {
      toast.error(`Failed to load persisted store values from ${dbName}`, {
        description: `Error: ${error}`,
      });
    }

    // Create map with the debounced persist functions for each key
    const persistDebouncedMap = new Map<string, (value: unknown) => void>();
    for (const keyToPersist of opts.keysToPersist) {
      const persistFn = async (value: unknown) => {
        status.saving = true;

        try {
          // Delete null or undefined values from the database
          if (value === null || value === undefined) {
            await db.removeItem(keyToPersist as string);
          } else {
            // localforage can only store JSON-serializable values and svelte
            // wraps stores in proxies, so we need to convert the value to a
            // plain object before persisting it
            const plainValue = JSON.parse(JSON.stringify(value));
            await db.setItem(keyToPersist as string, plainValue);
          }
        } catch (error) {
          // On error, remove the item to avoid stale or corrupted data
          db.removeItem(keyToPersist as string);

          // Show error toast and log to console for debugging
          toast.error(
            `Failed to persist ${keyToPersist as string} value to the database ${dbName}`,
            {
              description: `Error: ${error}`,
            },
          );
          console.error(
            `Failed to persist ${keyToPersist as string} value to the database ${dbName}`,
            error,
          );
          console.log({ value });
        } finally {
          status.saving = false;
        }
      };

      const delayMs = 300;
      const persistFnDebounced = debounce(persistFn, delayMs);
      persistDebouncedMap.set(keyToPersist as string, persistFnDebounced);
    }

    // Create an $effect to persist changes to the database
    $effect.root(() => {
      for (const keyToPersist of opts.keysToPersist) {
        $effect(() => {
          const value = store[keyToPersist];
          const persistFn = persistDebouncedMap.get(keyToPersist as string);
          if (persistFn) persistFn(value);
        });
      }
    });

    // biome-ignore lint/style/noNonNullAssertion: the function is always defined above
    resolveReady!();
    status.loading = false;
  })();

  return {
    store,
    status,
    ready: () => readyPromise,
  };
}

/**
 * Creates a prefix string based on the current URL path. This prefix allows
 * to use the same database names across different deployments under the
 * same domain, avoiding collisions.
 *
 * @returns A prefix string based on the current URL path, suitable for use in database names.
 */
function createGlobalDbNamePrefix(): string {
  if (!browser) return "";

  const prefix = globalThis.location.pathname
    .replace(/[^a-z0-9]/gi, "-")
    .toLowerCase();

  return prefix;
}
