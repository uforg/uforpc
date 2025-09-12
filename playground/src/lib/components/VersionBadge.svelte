<script>
  import { ExternalLink } from "@lucide/svelte";

  import { versionWithPrefix } from "$lib/version";
</script>

<div class="badge badge-ghost shadow-sm">
  <a
    href="https://github.com/uforg/uforpc/releases/tag/{versionWithPrefix}"
    target="_blank"
    class="flex items-center"
  >
    <ExternalLink class="mr-2 size-4" />
    <span>
      UFO RPC {versionWithPrefix}
    </span>
  </a>
</div>
