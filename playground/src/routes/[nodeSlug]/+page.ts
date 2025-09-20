import { error } from "@sveltejs/kit";

import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params }) => {
  if (!params.nodeSlug) error(404, "Not found");

  const firstDashIndex = params.nodeSlug.indexOf("-");
  if (firstDashIndex === -1) error(404, "Not found");

  const nodeKind = params.nodeSlug.substring(0, firstDashIndex);
  const nodeName = params.nodeSlug.substring(firstDashIndex + 1);

  if (!nodeKind || !nodeName) error(404, "Not found");

  return {
    nodeSlug: params.nodeSlug,
    nodeKind: nodeKind,
    nodeName: nodeName,
  };
};
