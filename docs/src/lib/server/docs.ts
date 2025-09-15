import fg from "fast-glob";

/**
 * Base directory where the documents are stored, must be relative to the root directory.
 */
const baseDirectory = "./src/docs";
const extensions = [".md", ".svx"];

export interface DocMeta {
  slug: string;
  title: string;
  description: string;
  isDraft: boolean;
  publishedAt: string;
  updatedAt: string;
}

interface DocImport {
  metadata: Record<string, any>;
  default: any;
}

interface DocPath {
  absolute: string;
  relativeToBase: string;
  slug: string;
}

/**
 * Recursively retrieves all document paths with specified extensions from the base directory.
 *
 * @returns Array of document paths relative to the root directory
 */
function getDocsPaths(): DocPath[] {
  let paths: string[] = [];

  for (const ext of extensions) {
    const files = fg.sync(`**/*${ext}`, { cwd: baseDirectory });
    paths.push(...files);
  }

  const docPaths: DocPath[] = paths.map((path) => {
    const absolute = `${baseDirectory}/${path}`;
    const relativeToBase = path;
    const slug = getSlugFromPath(path);
    return { absolute, relativeToBase, slug };
  });

  return docPaths;
}

/**
 * Generates a URL-friendly slug from a given file path.
 *
 * @param path - The file path to convert into a slug
 * @returns A URL-friendly slug derived from the file path
 */
function getSlugFromPath(path: string): string {
  let slug = path.replace(/\\/g, "/");
  slug = slug.replaceAll(/ /g, "-");
  slug = slug.replaceAll(/--+/g, "-");

  for (const ext of extensions) {
    if (slug.endsWith(ext)) {
      slug = slug.slice(0, -ext.length);
    }
  }

  slug = slug.replaceAll(/[^a-zA-Z0-9-_/]/g, "");
  return slug;
}

/**
 * Fetches and processes all documents from the base directory, extracting their metadata.
 *
 * @returns Array of document metadata objects
 */
export async function getDocsMeta() {
  const docPaths = getDocsPaths();
  let allImported: { docPath: DocPath; docImport: DocImport }[] = [];
  let docs: DocMeta[] = [];

  for await (const path of docPaths) {
    const relativePath = `../../docs/${path.relativeToBase}`;
    const post = await import(/* @vite-ignore */ relativePath);
    allImported = [
      ...allImported,
      {
        docImport: post,
        docPath: path,
      },
    ];
  }

  for (const imported of allImported) {
    if (typeof imported.docImport !== "object") continue;
    if (!imported.docImport.metadata) continue;
    if (!imported.docPath.slug) continue;

    const metadata = imported.docImport.metadata as Omit<DocMeta, "slug">;
    const post: DocMeta = { ...metadata, slug: imported.docPath.slug };

    if (!post.isDraft) docs.push(post);
  }

  return docs;
}
