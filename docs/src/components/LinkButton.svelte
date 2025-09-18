<script lang="ts">
  interface Props {
    href: string;
    class?: string;
    children?: any;
    rest?: any;
  }

  let { href, class: className, children, ...rest }: Props = $props();
</script>

<a
  {href}
  class={[
    "inline-flex w-full items-center justify-center rounded-full border border-transparent bg-indigo-600 px-8 py-4 text-base leading-6 font-bold text-white no-underline hover:bg-indigo-500 md:w-auto",
    className ?? "",
  ]}
  {...rest}
>
  {#if children}
    {@render children()}
  {/if}
</a>
