import { getDocsMeta } from "$lib/server/docs";

import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async () => {
  const docsMeta = await getDocsMeta();

  return {
    docsMeta,
  };
};
