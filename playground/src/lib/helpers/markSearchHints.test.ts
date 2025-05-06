import { expect, describe, it } from "vitest";
import { markSearchHints, truncateWithMark } from "./markSearchHints.ts";

// Mock SearchResult object for testing
const createMockResult = (terms: string[]) => ({
  id: "test-id",
  terms,
  queryTerms: terms,
  score: 1,
  match: {},
});

describe("markSearchHints", () => {
  it("highlights a single term", () => {
    const searchResult = createMockResult(["test"]);
    const text = "This is a test string";
    const result = markSearchHints(searchResult, text);
    expect(result).toBe("This is a <mark>test</mark> string");
  });

  it("highlights multiple occurrences of a term", () => {
    const searchResult = createMockResult(["test"]);
    const text = "test this test string test";
    const result = markSearchHints(searchResult, text);
    expect(result).toBe(
      "<mark>test</mark> this <mark>test</mark> string <mark>test</mark>",
    );
  });

  it("highlights multiple different terms", () => {
    const searchResult = createMockResult(["test", "string"]);
    const text = "This is a test string";
    const result = markSearchHints(searchResult, text);
    expect(result).toBe("This is a <mark>test</mark> <mark>string</mark>");
  });

  it("case insensitive matching", () => {
    const searchResult = createMockResult(["test"]);
    const text = "This is a TEST string";
    const result = markSearchHints(searchResult, text);
    expect(result).toBe("This is a <mark>TEST</mark> string");
  });

  it("returns original text if no terms match", () => {
    const searchResult = createMockResult(["example"]);
    const text = "This is a test string";
    const result = markSearchHints(searchResult, text);
    expect(result).toBe("This is a test string");
  });
});

describe("truncateWithMark", () => {
  it("returns text as is when match is in first 3 words", () => {
    const searchResult = createMockResult(["test"]);
    const text = "This is test string longer than needed";
    const result = truncateWithMark(searchResult, text);
    expect(result).toBe("This is <mark>test</mark> string longer than needed");
  });

  it("truncates when match is beyond first 3 words", () => {
    const searchResult = createMockResult(["longer"]);
    const text = "This is a test string longer than needed";
    const result = truncateWithMark(searchResult, text);
    expect(result).toBe("... a test string <mark>longer</mark> than needed");
  });

  it("returns original text when no match is found", () => {
    const searchResult = createMockResult(["example"]);
    const text = "This is a test string longer than needed";
    const result = truncateWithMark(searchResult, text);
    expect(result).toBe("This is a test string longer than needed");
  });

  it("handles case where match is at the very beginning", () => {
    const searchResult = createMockResult(["this"]);
    const text = "This is a test string";
    const result = truncateWithMark(searchResult, text);
    expect(result).toBe("<mark>This</mark> is a test string");
  });

  it("handles multiple matches and uses the first one", () => {
    const searchResult = createMockResult(["test"]);
    const text = "One two three four test five six seven test eight";
    const result = truncateWithMark(searchResult, text);
    expect(result).toBe(
      "... two three four <mark>test</mark> five six seven <mark>test</mark> eight",
    );
  });

  it("empty string returns empty string", () => {
    const searchResult = createMockResult(["test"]);
    const text = "";
    const result = truncateWithMark(searchResult, text);
    expect(result).toBe("");
  });
});
