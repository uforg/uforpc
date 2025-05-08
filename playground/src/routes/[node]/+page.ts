import { error } from "@sveltejs/kit";

import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params }) => {
  if (!params.node) error(404, "Not found");

  const [nodeKind, nodeName] = params.node.split("-");
  if (!nodeKind || !nodeName) error(404, "Not found");

  return {
    node: params.node,
    nodeKind: nodeKind,
    nodeName: nodeName,
  };
};
