import { clsx } from "clsx";
import type { ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export type { ClassValue };

/**
 * Merge classes using clsx and tailwind-merge
 * @param inputs - The classes to merge
 * @returns The merged classes
 */
export const mergeClasses = (...inputs: ClassValue[]) => {
  return twMerge(clsx(inputs));
};
