import { describe, expect, it, vi } from "vitest";

import { debounce } from "./debounce";

describe("debounce", () => {
  vi.useFakeTimers();

  it("should debounce a function", () => {
    const func = vi.fn();
    const debouncedFunc = debounce(func, 1000);

    debouncedFunc();
    debouncedFunc();
    debouncedFunc();

    // The function should not be called immediately
    expect(func).not.toHaveBeenCalled();

    // Fast-forward time by 1000ms
    vi.advanceTimersByTime(1000);

    // Now the function should have been called once
    expect(func).toHaveBeenCalledTimes(1);
  });

  it("should call the function with the last arguments", () => {
    const func = vi.fn();
    const debouncedFunc = debounce(func, 1000);

    debouncedFunc(1);
    debouncedFunc(2);
    debouncedFunc(3);

    vi.advanceTimersByTime(1000);

    expect(func).toHaveBeenCalledWith(3);
  });

  it("should not call the function if the wait time has not passed", () => {
    const func = vi.fn();
    const debouncedFunc = debounce(func, 1000);

    debouncedFunc();

    vi.advanceTimersByTime(500);

    expect(func).not.toHaveBeenCalled();
  });

  it("should call the function again after the wait time", () => {
    const func = vi.fn();
    const debouncedFunc = debounce(func, 1000);

    debouncedFunc();
    vi.advanceTimersByTime(1000);
    expect(func).toHaveBeenCalledTimes(1);

    debouncedFunc();
    vi.advanceTimersByTime(1000);
    expect(func).toHaveBeenCalledTimes(2);
  });

  it("should maintain the correct `this` context", () => {
    const func = vi.fn(function (this: { a: number }) {
      return this.a;
    });
    const context = { a: 1, debouncedFunc: debounce(func, 1000) };

    context.debouncedFunc();
    vi.advanceTimersByTime(1000);

    expect(func).toHaveBeenCalledTimes(1);
    expect(func).toHaveBeenCalledWith();
    expect(func.mock.instances[0]).toEqual({
      a: 1,
      debouncedFunc: expect.any(Function),
    });
  });
});
