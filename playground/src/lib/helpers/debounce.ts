/**
 * Creates a debounced function that delays invoking the provided function until after
 * a specified wait time has elapsed since the last time the debounced function was invoked.
 *
 * @template TArgs - The argument types for the debounced function.
 * @param func - The function to debounce.
 * @param wait - The number of milliseconds to delay.
 * @returns A new debounced function.
 */
export function debounce<TArgs extends unknown[]>(
  func: (...args: TArgs) => void,
  wait: number,
): (...args: TArgs) => void {
  let timeout: ReturnType<typeof setTimeout> | undefined;

  return function (this: unknown, ...args: TArgs): void {
    clearTimeout(timeout);
    timeout = setTimeout(() => {
      func.apply(this, args);
    }, wait);
  };
}
