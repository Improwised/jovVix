import { describe, expect, it } from "vitest";
import { questionsAnalysis } from "../../composables/questions_analysis";

describe("questionsAnalysis", () => {
  it("returns default object for null input", () => {
    const result = questionsAnalysis(null);
    expect(result).toEqual({
      rank: 0,
      totalScore: 0,
      totalQuestions: 0,
      correctAnwers: 0,
      wrongAnwers: 0,
      totalSurveyQuestions: 0,
      attemptedSurveyQuestions: 0,
      unAttemptedQuestions: 0,
      accuracy: 0,
    });
  });

  it("extracts rank and totalScore from the seed entry", () => {
    const result = questionsAnalysis([
      { rank: 3, total_score: 500 },
      {
        question: "Q1",
        is_attend: true,
        selected_answer: { String: "[1]" },
        correct_answer: "[1]",
        question_type: "mcq",
      },
    ]);
    expect(result.rank).toBe(3);
    expect(result.totalScore).toBe(500);
  });

  it("counts correct answers and calculates accuracy", () => {
    const result = questionsAnalysis([
      {
        question: "Q1",
        is_attend: true,
        selected_answer: { String: "[1]" },
        correct_answer: "[1]",
        question_type: "quiz",
      },
      {
        question: "Q2",
        is_attend: true,
        selected_answer: { String: "[2]" },
        correct_answer: "[1]",
        question_type: "quiz",
      },
    ]);
    expect(result.totalQuestions).toBe(2);
    expect(result.correctAnwers).toBe(1);
    expect(result.wrongAnwers).toBe(1);
    expect(result.accuracy).toBe(50);
  });

  it("counts unattempted questions", () => {
    const result = questionsAnalysis([
      { question: "Q1", is_attend: false, question_type: "mcq" },
      {
        question: "Q2",
        is_attend: true,
        selected_answer: { String: "[1]" },
        correct_answer: "[1]",
        question_type: "quiz",
      },
    ]);
    expect(result.unAttemptedQuestions).toBe(1);
    expect(result.correctAnwers).toBe(1);
  });

  it("handles survey questions separately from quiz questions", () => {
    const result = questionsAnalysis([
      { question: "Q1", is_attend: true, question_type: "survey" },
      { question: "Q2", is_attend: false, question_type: "survey" },
    ]);
    expect(result.totalSurveyQuestions).toBe(2);
    expect(result.attemptedSurveyQuestions).toBe(1);
    expect(result.totalQuestions).toBe(2);
  });
});
