import { describe, expect, it } from "vitest";
import {
  accuracy,
  combinedRows,
  overallAccuracy,
} from "../../composables/combined_report";

const report = {
  sessions: [{ id: "s1" }, { id: "s2" }],
  leaderboard: [
    { firstname: "Chen", session_id: "s2", correct: 1, questions: 2, plays: 1 },
    { firstname: "Asha", session_id: "s2", correct: 2, questions: 2, plays: 2 },
    { firstname: "Bo", session_id: "s1", correct: 0, questions: 3, plays: 1 },
  ],
};

describe("combined report helpers", () => {
  it("rounds accuracy and treats zero questions as 0%", () => {
    expect(accuracy(2, 3)).toBe(67);
    expect(accuracy(0, 0)).toBe(0);
  });

  it("weights overall accuracy by answers, not by player", () => {
    // 3 of 7 answers are correct; the mean of 50%, 100% and 0% would be 50%.
    expect(overallAccuracy(report.leaderboard)).toBe(43);
    expect(overallAccuracy()).toBe(0);
  });

  it("labels each row with its session, accuracy and repeat plays", () => {
    expect(combinedRows(report).map((r) => r.detail)).toEqual([
      "Session 2 · 1/2 correct (50%)",
      "Session 2 · 2/2 correct (100%) · best of 2 plays",
      "Session 1 · 0/3 correct (0%)",
    ]);
    expect(combinedRows(report)[0].firstname).toBe("Chen");
    expect(combinedRows(null)).toEqual([]);
  });
});
