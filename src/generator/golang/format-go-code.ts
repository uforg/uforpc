import { promises as fs } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { exec } from "node:child_process";

/**
 * Formats Go code using `go fmt`.
 *
 * Creates a temporary file, writes the code to it, runs `go fmt` on the file,
 * reads the formatted code, deletes the temporary file, and returns the
 * formatted code.
 *
 * @param {string} code - The Go code to format.
 * @returns {Promise<string>} The formatted Go code.
 *
 * @requires Go installed on the system.
 */
export async function formatGoCode(code: string): Promise<string> {
  const tempDir = tmpdir();
  const tempFileName = `temp-${crypto.randomUUID()}.go`;
  const tempFilePath = join(tempDir, tempFileName);

  try {
    await fs.writeFile(tempFilePath, code);

    await new Promise<void>((resolve, reject) => {
      exec(`go fmt ${tempFilePath}`, (error) => {
        if (error) {
          reject(error);
        } else {
          resolve();
        }
      });
    });

    const formattedCode = await fs.readFile(tempFilePath, "utf8");

    return formattedCode;
  } finally {
    await fs.unlink(tempFilePath);
  }
}
