<script lang="ts">
  import H1 from "$lib/components/H1.svelte";
  import H2 from "$lib/components/H2.svelte";
  import UfoAbduction from "$lib/components/UfoAbduction.svelte";
</script>

<header class="grid w-full grid-cols-1 gap-4 py-8 desk:grid-cols-12">
  <div
    class="flex items-center justify-center desk:col-span-7 desk:justify-end"
  >
    <hgroup>
      <H1>UFO RPC</H1>
      <H2>The RPC framework that puts developer experience first.</H2>
    </hgroup>
  </div>

  <div
    class="flex items-center justify-center desk:col-span-5 desk:justify-start"
  >
    <div>
      <UfoAbduction width={200} />
    </div>
  </div>
</header>
