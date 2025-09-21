<script lang="ts">
  export interface TabItem {
    id: string;
    label: string;
  }

  interface Props {
    items: TabItem[];
    activeId: string;
    onSelect?: (id: string) => void;
  }

  let {
    items,
    activeId = $bindable(),
    onSelect: parentOnSelect,
  }: Props = $props();

  const onSelect = (id: string) => {
    activeId = id;
    if (parentOnSelect) parentOnSelect(id);
  };
</script>

<div class="join bg-base-100 flex w-full overflow-x-auto overflow-y-hidden">
  {#each items as tab}
    <button
      class={[
        "btn join-item border-base-content/20 flex-grow",
        activeId === tab.id && "btn-primary",
      ]}
      onclick={() => onSelect(tab.id)}
      type="button"
    >
      {tab.label}
    </button>
  {/each}
</div>
