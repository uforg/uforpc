import { getPublishedDocsInfo } from "$lib/docs";

import type { LayoutLoad } from "./$types";

export const load: LayoutLoad = async () => {
  const docsInfo = getPublishedDocsInfo();

  return {
    docsInfo,
  };
};
