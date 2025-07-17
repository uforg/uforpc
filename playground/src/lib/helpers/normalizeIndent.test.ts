// normalizeIndent.test.ts
import { describe, expect, it } from "vitest";

import { normalizeIndent } from "./normalizeIndent.ts";

// Adjust the import path as needed

describe("normalizeIndent", () => {
  it("should normalize a standard code block with tabs", () => {
    const indentedText = `
			// This is a comment.
			function myFunc() {
				// This has extra indentation.
				return true;
			}
		`;

    const expectedText = `
// This is a comment.
function myFunc() {
	// This has extra indentation.
	return true;
}
		`;
    expect(normalizeIndent(indentedText)).toBe(expectedText);
  });

  it("should correctly normalize a Markdown code block", () => {
    const markdownText = `
    Here is some text.

    \`\`\`typescript
      // Some indented code inside a markdown block
      const x = 1;
      if (x > 0) {
        console.log('Positive!');
      }
    \`\`\`
    `;

    const expectedText = `
Here is some text.

\`\`\`typescript
  // Some indented code inside a markdown block
  const x = 1;
  if (x > 0) {
    console.log('Positive!');
  }
\`\`\`
`;
    // The base indentation here is 4 spaces.
    // Notice the lines inside ``` are only dedented by 4 spaces, not fully.
    expect(normalizeIndent(markdownText)).toBe(expectedText);
  });

  it("should not change text that has no leading indentation", () => {
    const unindentedText = `Hello world.
This is a test.
  - With a list item.`;

    expect(normalizeIndent(unindentedText)).toBe(unindentedText);
  });

  it("should handle an empty string as input", () => {
    expect(normalizeIndent("")).toBe("");
  });

  it("should handle a string with only whitespace", () => {
    const whitespaceText = `
      
    
    `;
    expect(normalizeIndent(whitespaceText)).toBe(whitespaceText);
  });

  it("should handle text that starts with empty lines", () => {
    // The function should find the first *non-empty* line to get the indentation.
    const indentedText = `

		// First line of actual content
		const value = 42;
`;

    const expectedText = `

// First line of actual content
const value = 42;
`;
    expect(normalizeIndent(indentedText)).toBe(expectedText);
  });

  it("should handle a mix of spaces and tabs for indentation", () => {
    // The base indentation is `\t  ` (one tab, two spaces)
    const mixedIndentation = `
	  first line
	  	second line with more
	  third line
`;

    const expected = `
first line
	second line with more
third line
`;
    expect(normalizeIndent(mixedIndentation)).toBe(expected);
  });
});
