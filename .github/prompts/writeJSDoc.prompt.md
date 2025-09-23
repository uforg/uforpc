---
mode: agent
---

Your primary task is to write JSDOC documentation for JavaScript and TypeScript code blocks. You must adhere to the following rules with absolute precision.

**VITAL CONSTRAINTS:**

1.  **DO NOT MODIFY EXISTING CODE.** Your only responsibility is to add JSDOC comment blocks (`/** ... */`) directly above the code they document. You must not change, add, or remove any part of the source code itself.
2.  **LANGUAGE:** All documentation must be written in **English** unless a different language is explicitly requested.
3.  **CLARITY AND BREVITY:** The documentation must be clear, concise, and short. Explain the purpose of the code, its parameters, and what it returns in a way that is easy for any developer to understand quickly, ensuring that brevity does not compromise understandability.

**DOCUMENTATION RULES:**

- **For TypeScript files (`.ts`, `.tsx`):** Since TypeScript uses its own type system, your JSDOC **must not** include type definitions. Focus solely on the description.
  - Example:
    ```typescript
    /**
     * Greets a user by their name.
     * @param name - The name of the user to greet.
     * @returns A greeting message.
     */
    ```
- **For JavaScript files (`.js`, `.jsx`):** Your JSDOC **must include** type definitions for all parameters and return values using standard JSDOC syntax.
  - Example:
    ```javascript
    /**
     * Greets a user by their name.
     * @param {string} name - The name of the user to greet.
     * @returns {string} A greeting message.
     */
    ```

**SUCCESS CRITERIA:**

- The output perfectly preserves the original code.
- High-quality, clear, and concise JSDOC documentation has been added.
- The correct JSDOC style (with or without types) has been applied based on the file type (JavaScript or TypeScript).
