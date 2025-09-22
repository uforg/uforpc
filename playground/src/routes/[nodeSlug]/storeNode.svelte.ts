import { createAsyncStore } from "$lib/createAsyncStore.svelte";

type Input = object;
type Output = string;
type Date = string;

export interface HistoryItem {
  input: Input;
  output: Output;
  date: Date;
}

export interface StoreNode {
  input: Input;
  output: Output;
  date: Date;
  history: HistoryItem[];
}

export type StoreNodeInstance = ReturnType<typeof createStoreNode>;

type StoreNodeKey = keyof StoreNode;

const storeNodeDefault: StoreNode = {
  input: {},
  output: "",
  date: "",
  history: [],
};

const storeNodeKeysToPersist: StoreNodeKey[] = [
  "input",
  "output",
  "date",
  "history",
];

export const createStoreNode = (nodeSlug: string) => {
  return createAsyncStore<StoreNode>({
    initialValue: async () => storeNodeDefault,
    keysToPersist: storeNodeKeysToPersist,
    dbName: "storeNode",
    tableName: nodeSlug,
  });
};

/**
 * Save the current input and output to the history for a given operation ID, initializing it if necessary.
 * Limits the history to the most recent 150 entries.
 *
 * @param operationID the operation ID to save the current input and output for
 */
export const saveCurrentToHistory = (
  storeNode: ReturnType<typeof createStoreNode>,
) => {
  const historyLimit = 150;
  const input = storeNode.store.input;
  const output = storeNode.store.output;

  if (storeNode.status.loading) return;
  if (!input && !output) return;

  storeNode.store.history.unshift({
    input,
    output,
    date: new Date().toISOString(),
  });

  if (storeNode.store.history.length > historyLimit) {
    storeNode.store.history.pop();
  }
};
