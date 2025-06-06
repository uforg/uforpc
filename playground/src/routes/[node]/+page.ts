import { error } from "@sveltejs/kit";

import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params }) => {
  if (!params.node) error(404, "Not found");

  const firstDashIndex = params.node.indexOf("-");
  if (firstDashIndex === -1) error(404, "Not found");

  const nodeKind = params.node.substring(0, firstDashIndex);
  const nodeName = params.node.substring(firstDashIndex + 1);

  if (!nodeKind || !nodeName) error(404, "Not found");

  return {
    node: params.node,
    nodeKind: nodeKind,
    nodeName: nodeName,
  };
};
