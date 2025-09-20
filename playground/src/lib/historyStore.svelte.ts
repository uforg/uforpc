import { createStore } from "./storeHelpers.svelte";

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

export interface HistoryStore {
  operations: Record<string, OperationHistory>;
}

type HistoryStoreKey = keyof HistoryStore;

const defaultHistoryStore: HistoryStore = {
  operations: {},
};

const historyStoreKeysToPersist: HistoryStoreKey[] = ["operations"];

export const historyStore = createStore<HistoryStore>(
  { ...defaultHistoryStore },
  historyStoreKeysToPersist,
);

/**
 * Initialize the operation history for a given operation ID if it doesn't exist.
 *
 * @param operationID the operation ID to initialize
 */
const initializeOperation = (operationID: string) => {
  if (!historyStore.operations[operationID]) {
    historyStore.operations[operationID] = {
      input: {},
      output: "",
      history: [],
    };
  }

  if (!historyStore.operations[operationID].input) {
    historyStore.operations[operationID].input = {};
  }

  if (!historyStore.operations[operationID].output) {
    historyStore.operations[operationID].output = "";
  }

  if (!historyStore.operations[operationID].history) {
    historyStore.operations[operationID].history = [];
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
  return historyStore.operations[operationID].input;
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
  return historyStore.operations[operationID].output;
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
  historyStore.operations[operationID].input = input;
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
  historyStore.operations[operationID].output = output;
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
    historyStore.operations[operationID].history.unshift({
      input,
      output,
      date: new Date().toISOString(),
    });
    // Limit history entries
    if (historyStore.operations[operationID].history.length > historyLimit) {
      historyStore.operations[operationID].history.pop();
    }
  }
};
