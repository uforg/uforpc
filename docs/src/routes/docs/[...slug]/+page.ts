import { error } from "@sveltejs/kit";

import { getPublishedDocContent } from "$lib/docs";

import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params }) => {
  try {
    const { content, docInfo, exists } = await getPublishedDocContent(
      params.slug,
    );

    if (!exists || !docInfo || !content) {
      return error(404, `Could not find ${params.slug}`);
    }

    if (!content) {
      return error(404, `Missing content for ${params.slug}`);
    }

    if (!docInfo) {
      return error(404, `Missing metadata for ${params.slug}`);
    }

    return {
      docInfo,
      content,
    };
  } catch (e) {
    const errMsg = `Could not load ${params.slug} (${(e as Error).message})`;
    return error(404, errMsg);
  }
};
