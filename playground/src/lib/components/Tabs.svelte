<script lang="ts" generics="T extends string">
  import type { IconProps } from "@lucide/svelte";
  import type { Component } from "svelte";

  import { type ClassValue, mergeClasses } from "$lib/helpers/mergeClasses";

  export interface TabItem<T extends string = string> {
    id: T;
    label?: string;
    icon?: Component<IconProps, {}, "">;
  }

  interface Props<T extends string> {
    items: TabItem<T>[];
    active: T;
    onSelect?: (id: T) => void;
    containerClass?: ClassValue;
    buttonClass?: ClassValue;
    activeButtonClass?: ClassValue;
    inactiveButtonClass?: ClassValue;
  }

  let {
    items,
    active = $bindable(),
    onSelect,
    containerClass,
    buttonClass,
    activeButtonClass,
    inactiveButtonClass,
  }: Props<T> = $props();

  const handleSelect = (id: T) => {
    active = id;
    if (onSelect) onSelect(id);
  };
</script>

<div
  class={mergeClasses(
    "join bg-base-100 flex w-full overflow-x-auto overflow-y-hidden",
    containerClass,
  )}
>
  {#each items as tab}
    <button
      class={[
        "btn join-item border-base-content/20 flex-grow",
        buttonClass,
        active === tab.id && "btn-primary",
        active === tab.id && activeButtonClass,
        active !== tab.id && inactiveButtonClass,
      ]}
      onclick={() => handleSelect(tab.id)}
      type="button"
    >
      {#if tab.icon}
        <tab.icon class="size-4" />
      {/if}
      {#if tab.label}
        <span>
          {tab.label}
        </span>
      {/if}
    </button>
  {/each}
</div>
