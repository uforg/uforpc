import { toast } from "svelte-sonner";

export async function copyTextToClipboard(textToCopy: string): Promise<void> {
  try {
    await navigator.clipboard.writeText(textToCopy);
    toast.success("Copied to clipboard", { duration: 1500 });
  } catch (err) {
    console.error("Failed to copy to clipboard: ", err);
    toast.error("Failed to copy to clipboard", {
      description: `Error: ${err}`,
    });
  }
}
