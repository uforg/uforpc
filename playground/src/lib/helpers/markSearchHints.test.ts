import { assertEquals } from "@std/assert";
import { markSearchHints, truncateWithMark } from "./markSearchHints.ts";

// Mock SearchResult object for testing
const createMockResult = (terms: string[]) => ({
  id: "test-id",
  terms,
  queryTerms: terms,
  score: 1,
  match: {},
});

Deno.test("markSearchHints - highlights a single term", () => {
  const searchResult = createMockResult(["test"]);
  const text = "This is a test string";
  const result = markSearchHints(searchResult, text);
  assertEquals(result, "This is a <mark>test</mark> string");
});

Deno.test("markSearchHints - highlights multiple occurrences of a term", () => {
  const searchResult = createMockResult(["test"]);
  const text = "test this test string test";
  const result = markSearchHints(searchResult, text);
  assertEquals(
    result,
    "<mark>test</mark> this <mark>test</mark> string <mark>test</mark>",
  );
});

Deno.test("markSearchHints - highlights multiple different terms", () => {
  const searchResult = createMockResult(["test", "string"]);
  const text = "This is a test string";
  const result = markSearchHints(searchResult, text);
  assertEquals(result, "This is a <mark>test</mark> <mark>string</mark>");
});

Deno.test("markSearchHints - case insensitive matching", () => {
  const searchResult = createMockResult(["test"]);
  const text = "This is a TEST string";
  const result = markSearchHints(searchResult, text);
  assertEquals(result, "This is a <mark>TEST</mark> string");
});

Deno.test("markSearchHints - returns original text if no terms match", () => {
  const searchResult = createMockResult(["example"]);
  const text = "This is a test string";
  const result = markSearchHints(searchResult, text);
  assertEquals(result, "This is a test string");
});

Deno.test("truncateWithMark - returns text as is when match is in first 3 words", () => {
  const searchResult = createMockResult(["test"]);
  const text = "This is test string longer than needed";
  const result = truncateWithMark(searchResult, text);
  assertEquals(result, "This is <mark>test</mark> string longer than needed");
});

Deno.test("truncateWithMark - truncates when match is beyond first 3 words", () => {
  const searchResult = createMockResult(["longer"]);
  const text = "This is a test string longer than needed";
  const result = truncateWithMark(searchResult, text);
  assertEquals(result, "... a test string <mark>longer</mark> than needed");
});

Deno.test("truncateWithMark - returns original text when no match is found", () => {
  const searchResult = createMockResult(["example"]);
  const text = "This is a test string longer than needed";
  const result = truncateWithMark(searchResult, text);
  assertEquals(result, "This is a test string longer than needed");
});

Deno.test("truncateWithMark - handles case where match is at the very beginning", () => {
  const searchResult = createMockResult(["this"]);
  const text = "This is a test string";
  const result = truncateWithMark(searchResult, text);
  assertEquals(result, "<mark>This</mark> is a test string");
});

Deno.test("truncateWithMark - handles multiple matches and uses the first one", () => {
  const searchResult = createMockResult(["test"]);
  const text = "One two three four test five six seven test eight";
  const result = truncateWithMark(searchResult, text);
  assertEquals(
    result,
    "... two three four <mark>test</mark> five six seven <mark>test</mark> eight",
  );
});

Deno.test("truncateWithMark - empty string returns empty string", () => {
  const searchResult = createMockResult(["test"]);
  const text = "";
  const result = truncateWithMark(searchResult, text);
  assertEquals(result, "");
});
