import { describe, expect, it } from "vitest";
import { detectQuizLanguage } from "../../composables/quizLanguage";

describe("detectQuizLanguage", () => {
  it("detects Japanese text", () => {
    expect(detectQuizLanguage("これはテストです")).toBe("japanese");
  });

  it("detects Chinese text", () => {
    expect(detectQuizLanguage("这是一个测试")).toBe("chinese");
  });

  it("detects Gujarati text", () => {
    expect(detectQuizLanguage("આ એક પરીક્ષણ છે")).toBe("gujarati");
  });

  it("returns english for latin script text", () => {
    expect(detectQuizLanguage("This is a test")).toBe("english");
  });

  it("returns english when input is empty", () => {
    expect(detectQuizLanguage("")).toBe("english");
  });

  it("detects language from an array of text samples", () => {
    expect(detectQuizLanguage(["これは", "テストです"])).toBe("japanese");
  });
});
