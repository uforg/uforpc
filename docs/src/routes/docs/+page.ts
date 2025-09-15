import { getPublishedDocsInfo } from "$lib/docs";

import type { PageLoad } from "./$types";

export const load: PageLoad = async () => {
  const docsInfo = getPublishedDocsInfo();

  return {
    docsInfo,
  };
};
