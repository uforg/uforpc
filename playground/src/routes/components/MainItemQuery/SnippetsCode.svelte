<script lang="ts">
  import Code from "$lib/components/Code.svelte";
  import * as curlconverter from "curlconverter";

  interface Props {
    curl: string;
  }

  const { curl }: Props = $props();

  interface Lang {
    group: string;
    langCode: string;
    label: string;
    func: (code: string) => string;
  }

  interface LangGroup {
    group: string;
    langs: Lang[];
  }

  const langs: Lang[] = [
    {
      group: "Bash",
      langCode: "bash",
      label: "Curl",
      func: (code: string) => code,
    },
    {
      group: "JavaScript",
      langCode: "javascript",
      label: "JavaScript Fetch",
      func: curlconverter.toJavaScript,
    },
    {
      group: "JavaScript",
      langCode: "javascript",
      label: "JavaScript JQuery",
      func: curlconverter.toJavaScriptJquery,
    },
    {
      group: "JavaScript",
      langCode: "javascript",
      label: "JavaScript XHR",
      func: curlconverter.toJavaScriptXHR,
    },
    {
      group: "PHP",
      langCode: "php",
      label: "PHP Curl",
      func: curlconverter.toPhp,
    },
    {
      group: "PHP",
      langCode: "php",
      label: "PHP Guzzle",
      func: curlconverter.toPhpGuzzle,
    },
  ];

  // This takes every lang and puts it into it's group
  const langGroups = $derived.by(() => {
    const groups: LangGroup[] = [];

    for (const lang of langs) {
      const group = groups.find((group) => group.group === lang.group);
      if (group) {
        group.langs.push(lang);
      } else {
        groups.push({ group: lang.group, langs: [lang] });
      }
    }

    return groups;
  });

  let pickedLabel = $state(langs[0].label);

  let pickedLang = $derived.by(() => {
    const lang = langs.find((lang) => lang.label === pickedLabel);
    if (!lang) return langs[0].langCode;
    return lang.langCode;
  });

  let pickedCode = $derived.by(() => {
    const lang = langs.find((lang) => lang.label === pickedLabel);
    if (!lang) return langs[0].func(curl);
    return lang.func(curl);
  });
</script>

<fieldset class="fieldset">
  <legend class="fieldset-legend">Language</legend>
  <select class="select w-full" bind:value={pickedLabel}>
    {#each langGroups as langGroup}
      {#if langGroup.langs.length > 1}
        <optgroup label={langGroup.group}>
          {#each langGroup.langs as lang}
            <option value={lang.label}>{lang.label}</option>
          {/each}
        </optgroup>
      {:else}
        <option value={langGroup.langs[0].label}>
          {langGroup.langs[0].label}
        </option>
      {/if}
    {/each}
  </select>
</fieldset>

<div class="prose mt-2">
  <Code code={pickedCode} lang={pickedLang} />
</div>
