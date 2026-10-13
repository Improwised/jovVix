import { describe, expect, it } from "vitest";
import { isCorrectAnswer } from "../../composables/check_is_correct";

describe("isCorrectAnswer", () => {
  it("returns true when a single correct answer is selected", () => {
    expect(isCorrectAnswer("[1]", "[1]")).toBe(true);
  });

  it("returns false when a wrong answer is selected", () => {
    expect(isCorrectAnswer("[2]", "[1]")).toBe(false);
  });

  it("returns true when all multiple correct answers are selected in order", () => {
    expect(isCorrectAnswer("[1,2,3]", "[1,2,3]")).toBe(true);
  });

  it("returns true when a subset of correct answers is selected", () => {
    expect(isCorrectAnswer("[1]", "[1,2,3]")).toBe(true);
  });

  it("returns false when selected includes a wrong answer", () => {
    expect(isCorrectAnswer("[1,4]", "[1,2,3]")).toBe(false);
  });

  it("returns false when selected answer is empty", () => {
    expect(isCorrectAnswer("[]", "[1]")).toBe(false);
  });

  it("returns true regardless of answer order", () => {
    expect(isCorrectAnswer("[3,1,2]", "[1,2,3]")).toBe(true);
  });

  it("returns false when selected is a bare number string", () => {
    expect(isCorrectAnswer("1", "[1]")).toBe(false);
  });
});
