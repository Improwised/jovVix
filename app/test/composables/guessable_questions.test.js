import { describe, it, expect } from "vitest";
import { findGuessablePatterns } from "~/composables/guessable_questions";

// q("stem", ["opt1", "opt2", ...], correctKey) builds a single-answer text question.
const q = (question, opts, correct, extra = {}) => ({
  question,
  question_type_id: 1,
  options_media: "text",
  options: Object.fromEntries(opts.map((text, i) => [String(i + 1), text])),
  correct_answer: JSON.stringify([].concat(correct)),
  ...extra,
});
const codes = (question, opts) =>
  findGuessablePatterns(question, opts).map((p) => p.code);

describe("findGuessablePatterns", () => {
  it("passes a well-built question", () => {
    expect(
      codes(
        q("Which planet is largest?", ["Jupiter", "Saturn", "Mars", "Venus"], 1)
      )
    ).toEqual([]);
  });

  it("flags every option marked correct", () => {
    expect(codes(q("Pick one", ["x", "y", "z"], [1, 2, 3]))).toEqual([
      "all_correct",
    ]);
  });

  it("flags two wrong options with the same text", () => {
    expect(
      codes(q("Capital of France?", ["Paris", "Rome", "rome!", "Oslo"], 1))
    ).toEqual(["duplicate"]);
  });

  it("flags a correct answer much longer than the rest", () => {
    expect(
      codes(
        q(
          "What does a firewall do?",
          ["Cools", "Filters network traffic by rules", "Prints", "Stores"],
          2
        )
      )
    ).toEqual(["longest"]);
  });

  it("does not flag length when the gap is small", () => {
    expect(codes(q("Pick", ["abcdef", "abc", "ab"], 1))).toEqual([]);
  });

  it("flags 'All of the above' as the correct answer", () => {
    expect(
      codes(
        q(
          "Which are fruits?",
          ["Apple trees", "Pear trees", "Plum trees", "All of the above"],
          4
        )
      )
    ).toEqual(["catch_all"]);
  });

  it("flags a hedge word only on the correct answer", () => {
    expect(
      codes(
        q(
          "Antibiotics treat?",
          ["Viral colds", "Usually bacteria", "Broken bones"],
          2
        )
      )
    ).toEqual(["hedge"]);
  });

  it("flags absolute words on every wrong option", () => {
    expect(
      codes(
        q("Sleep at night?", ["Always bad", "Never useful", "Helps memory"], 3)
      )
    ).toEqual(["absolutes"]);
  });

  it("flags a correct answer that alone repeats a question word", () => {
    expect(
      codes(
        q(
          "Which organ pumps blood?",
          ["Liver filter", "Blood pump", "Lung tissue"],
          2
        )
      )
    ).toEqual(["echo"]);
  });

  it("flags an article only the correct answer fits", () => {
    expect(
      codes(q("A tomato is an", ["vegetable", "fruit", "berry", "apple"], 4))
    ).toEqual(["article"]);
    expect(codes(q("Python is a", ["language", "snake", "game"], 1))).toEqual(
      []
    );
  });

  it("skips survey questions", () => {
    expect(
      codes(q("Favourite?", ["x", "x"], [], { question_type_id: 2 }))
    ).toEqual([]);
  });

  it("runs only language-independent checks for other languages", () => {
    const question = q(
      "Elige",
      ["Rojo", "Azul", "Todas las anteriores y además mucho más"],
      3
    );
    expect(codes(question, { english: false })).toEqual(["longest"]);
  });

  it("checks image options only for duplicates", () => {
    const question = q(
      "Which logo?",
      ["a.png", "b.png", "b.png", "a-much-longer-file-name.png"],
      4,
      {
        options_media: "image",
      }
    );
    expect(codes(question)).toEqual(["duplicate"]);
  });

  it("ignores blank option slots and accepts an answers array", () => {
    expect(
      codes(
        q("Which planet is largest?", ["Jupiter", "Saturn", "", ""], 1, {
          correct_answer: [1],
        })
      )
    ).toEqual([]);
  });
});
