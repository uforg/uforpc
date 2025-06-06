import fs from "node:fs";
import path from "node:path";

/**
 * Script to replace the REPLACEME placeholder with the correct path to the assets.
 * https://github.com/sveltejs/kit/issues/9569
 */

const readDirRecursive = async (filePath) => {
  const dir = await fs.promises.readdir(filePath);
  const files = await Promise.all(
    dir.map(async (relativePath) => {
      const absolutePath = path.join(filePath, relativePath);
      const stat = await fs.promises.lstat(absolutePath);
      return stat.isDirectory() ? readDirRecursive(absolutePath) : absolutePath;
    }),
  );
  return files.flat();
};

const files = await readDirRecursive("./build");

for (const file of files) {
  if (
    !(
      file.endsWith(".js") ||
      file.endsWith(".html") ||
      file.endsWith(".map") ||
      file.endsWith(".css")
    )
  ) {
    continue;
  }

  const data = await fs.promises.readFile(file, "utf8");
  await fs.promises.writeFile(
    file,
    data.replace(/http:\/\/REPLACEME/g, "."),
    "utf8",
  );
  console.log(`Wrote file '${file}'`);
}
