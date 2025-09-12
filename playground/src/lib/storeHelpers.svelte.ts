/**
 * Creates a Svelte store with persistence to local storage.
 *
 * @param initialStoreValue The initial value of the store.
 * @param keysToPersist The store keys of the store to persist.
 * @returns The created store.
 */
import { browser } from "$app/environment";

// biome-ignore lint/suspicious/noExplicitAny: since it's generic, it needs to be any
export function createStore<T extends Record<string, any>>(
  initialStoreValue: T,
  keysToPersist: (keyof T)[],
): T {
  // Make a shallow copy of the initial store value to avoid mutating it
  const initialValue = { ...initialStoreValue };

  // Load persisted values from local storage
  for (const keyToPersist of keysToPersist as string[]) {
    const localStorageKey = prefixLocalStorageKey(keyToPersist);

    // Search for the key in the local storage and if the value is not
    // found, do nothing, the default value is already set
    const value = globalThis.localStorage.getItem(localStorageKey);
    if (value === null) continue;

    // Deserialize and update the value based on the javascript type of it's default value
    const defaultValue = initialValue[keyToPersist];
    const defaultValueType = typeof defaultValue;

    switch (defaultValueType) {
      case "string":
        (initialValue[keyToPersist] as string) = value;
        break;
      case "boolean":
        (initialValue[keyToPersist] as boolean) = value === "true";
        break;
      case "number": {
        const numberValue = Number(value);
        if (!Number.isNaN(numberValue)) {
          (initialValue[keyToPersist] as unknown as number) = numberValue;
        }
        break;
      }
      case "object":
        try {
          const parsedValue = JSON.parse(value);
          (initialValue[keyToPersist] as object) = parsedValue;
        } catch {
          // Ignore invalid persisted object
        }
        break;
    }
  }

  // Create a Svelte store and add an $effect to persist changes to local storage
  const store = $state(initialValue);
  $effect.root(() => {
    $effect(() => {
      for (const keyToPersist of keysToPersist as string[]) {
        const localStorageKey = prefixLocalStorageKey(keyToPersist);
        const value = store[keyToPersist];
        const valueType = typeof value;

        // Delete null or undefined values from local storage
        if (value === null || value === undefined) {
          globalThis.localStorage.removeItem(localStorageKey);
          continue;
        }

        switch (valueType) {
          case "string":
            globalThis.localStorage.setItem(
              localStorageKey,
              value as unknown as string,
            );
            break;
          case "boolean":
            globalThis.localStorage.setItem(
              localStorageKey,
              (value as unknown as boolean).toString(),
            );
            break;
          case "number":
            globalThis.localStorage.setItem(
              localStorageKey,
              (value as unknown as number).toString(),
            );
            break;
          case "object":
            try {
              const stringifiedValue = JSON.stringify(value);
              globalThis.localStorage.setItem(
                localStorageKey,
                stringifiedValue,
              );
            } catch {
              // Ignore invalid object
            }
            break;
        }
      }
    });
  });

  return store;
}

function createLocalStoragePrefix(): string {
  if (!browser) return "";

  const prefix = globalThis.location.pathname
    .replace(/[^a-z0-9]/gi, "-")
    .toLowerCase();

  return prefix;
}

export function prefixLocalStorageKey(key: string): string {
  return `${createLocalStoragePrefix()}-${key}`;
}
