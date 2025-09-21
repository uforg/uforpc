import { createAsyncStore } from "$lib/createAsyncStore.svelte";

type Input = object;
type Output = string;

export interface HistoryItem {
  input: Input;
  output: Output;
  date: string;
}

export interface StoreNode {
  input: Input;
  output: Output;
  history: HistoryItem[];
}

export type StoreNodeInstance = ReturnType<typeof createStoreNode>;

type StoreNodeKey = keyof StoreNode;

const storeNodeDefault: StoreNode = {
  input: {},
  output: "",
  history: [],
};

const storeNodeKeysToPersist: StoreNodeKey[] = ["input", "output", "history"];

export const createStoreNode = (nodeSlug: string) => {
  return createAsyncStore<StoreNode>({
    initialValue: async () => storeNodeDefault,
    keysToPersist: storeNodeKeysToPersist,
    dbName: "storeNode",
    tableName: nodeSlug,
  });
};

/** * Save the current input and output to the history for a given operation ID, initializing it if necessary.
 * Limits the history to the most recent 50 entries.
 *
 * @param operationID the operation ID to save the current input and output for
 */
// export const saveCurrentToHistoryForOperation = (operationID: string) => {
//   const historyLimit = 50;

//   initializeOperation(operationID);
//   const input = getCurrentInputForOperation(operationID);
//   const output = getCurrentOutputForOperation(operationID);
//   if (input && output) {
//     storeHistory.store.operations[operationID].history.unshift({
//       input,
//       output,
//       date: new Date().toISOString(),
//     });
//     // Limit history entries
//     if (
//       storeHistory.store.operations[operationID].history.length > historyLimit
//     ) {
//       storeHistory.store.operations[operationID].history.pop();
//     }
//   }
// };
