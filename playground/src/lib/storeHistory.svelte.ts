import { createAsyncStore } from "./createAsyncStore.svelte";

export interface InOut {
  input: object;
  output: string;
  date: string;
}

export interface OperationHistory {
  input: InOut["input"];
  output: InOut["output"];
  history: InOut[];
}

export interface StoreHistory {
  operations: Record<string, OperationHistory>;
}

type StoreHistoryKey = keyof StoreHistory;

const defaultStoreHistory: StoreHistory = {
  operations: {},
};

const storeHistoryKeysToPersist: StoreHistoryKey[] = ["operations"];

export const storeHistory = createAsyncStore<StoreHistory>({
  initialValue: async () => defaultStoreHistory,
  keysToPersist: storeHistoryKeysToPersist,
  storeName: "storeHistory",
});

/**
 * Initialize the operation history for a given operation ID if it doesn't exist.
 *
 * @param operationID the operation ID to initialize
 */
const initializeOperation = (operationID: string) => {
  if (!storeHistory.store.operations[operationID]) {
    storeHistory.store.operations[operationID] = {
      input: {},
      output: "",
      history: [],
    };
  }

  if (!storeHistory.store.operations[operationID].input) {
    storeHistory.store.operations[operationID].input = {};
  }

  if (!storeHistory.store.operations[operationID].output) {
    storeHistory.store.operations[operationID].output = "";
  }

  if (!storeHistory.store.operations[operationID].history) {
    storeHistory.store.operations[operationID].history = [];
  }
};

/**
 * Get the current input for a given operation ID, initializing it if necessary.
 *
 * @param operationID the operation ID to get the current input for
 * @returns The reactive input for the given operation ID
 */
export const getCurrentInputForOperation = (
  operationID: string,
): InOut["input"] => {
  initializeOperation(operationID);
  return storeHistory.store.operations[operationID].input;
};

/** * Get the current output for a given operation ID, initializing it if necessary.
 *
 * @param operationID the operation ID to get the current output for
 * @returns The reactive output for the given operation ID
 */
export const getCurrentOutputForOperation = (
  operationID: string,
): InOut["output"] => {
  initializeOperation(operationID);
  return storeHistory.store.operations[operationID].output;
};

/** * Set the current input for a given operation ID, initializing it if necessary.
 *
 * @param operationID the operation ID to set the current input for
 * @param input the input to set
 */
export const setCurrentInputForOperation = (
  operationID: string,
  input: object,
) => {
  initializeOperation(operationID);
  storeHistory.store.operations[operationID].input = input;
};

/** * Set the current output for a given operation ID, initializing it if necessary.
 *
 * @param operationID the operation ID to set the current output for
 * @param output the output to set
 */
export const setCurrentOutputForOperation = (
  operationID: string,
  output: string,
) => {
  initializeOperation(operationID);
  storeHistory.store.operations[operationID].output = output;
};

/** * Save the current input and output to the history for a given operation ID, initializing it if necessary.
 * Limits the history to the most recent 50 entries.
 *
 * @param operationID the operation ID to save the current input and output for
 */
export const saveCurrentToHistoryForOperation = (operationID: string) => {
  const historyLimit = 50;

  initializeOperation(operationID);
  const input = getCurrentInputForOperation(operationID);
  const output = getCurrentOutputForOperation(operationID);
  if (input && output) {
    storeHistory.store.operations[operationID].history.unshift({
      input,
      output,
      date: new Date().toISOString(),
    });
    // Limit history entries
    if (
      storeHistory.store.operations[operationID].history.length > historyLimit
    ) {
      storeHistory.store.operations[operationID].history.pop();
    }
  }
};
