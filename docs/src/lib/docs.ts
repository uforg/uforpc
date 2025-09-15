import { docsInfo } from "./docs.gen";
import type { DocInfo } from "./docs.gen";

/**
 * Retreives the document information objects for all published documents.
 *
 * @returns Array of document information objects
 */
export function getPublishedDocsInfo() {
  let docs: DocInfo[] = [];

  for (const docInfo of docsInfo) {
    if (!docInfo.meta.isDraft) docs.push(docInfo);
  }

  return docs;
}

/**
 * Fetches a specific published document by its slug.
 *
 * @param slug - The slug of the document to retrieve
 * @returns An object containing the document's info
 */
export function getPublishedDocInfo(slug: string):
  | {
      docInfo: DocInfo;
      exists: true;
    }
  | {
      docInfo: null;
      exists: false;
    } {
  const publishedDocsInfo = getPublishedDocsInfo();
  const matchedPath = publishedDocsInfo.find((p) => p.slug === slug);

  if (!matchedPath) {
    return { docInfo: null, exists: false };
  }

  return { docInfo: matchedPath, exists: true };
}

/**
 * Fetches the content of a specific published document by its slug.
 *
 * @param slug - The slug of the document to retrieve
 * @returns An object containing the document's content and existence status
 */
export async function getPublishedDocContent(slug: string): Promise<
  | {
      docInfo: DocInfo;
      content: any;
      exists: boolean;
    }
  | {
      docInfo: null;
      content: null;
      exists: false;
    }
> {
  const { docInfo, exists } = getPublishedDocInfo(slug);
  if (!exists || !docInfo) {
    return { content: null, docInfo: null, exists: false };
  }

  try {
    const relativePath = `../docs/${docInfo.relative}`;
    const docImport = await import(/* @vite-ignore */ relativePath);
    const content = docImport.default;

    if (!content) {
      throw new Error(`Content or metadata missing for ${slug}`);
    }

    return { content, docInfo, exists: true };
  } catch (e) {
    return { content: null, docInfo: null, exists: false };
  }
}
