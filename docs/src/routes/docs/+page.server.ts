import { getDocs } from "$lib/server/docs";

import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async () => {
  const docs = await getDocs();

  return {
    docs,
  };
};
