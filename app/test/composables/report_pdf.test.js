import { describe, it, expect } from "vitest";
import {
  buildReportPdf,
  reportPdfFileName,
  reportPdfTables,
} from "~/composables/report_pdf";

// Rows as the report page holds them after its transform.
const questions = [
  {
    question: "Capital of France?",
    question_type: "single answer",
    options: { 1: "Paris", 2: "Lyon" },
    correct_answer: [1],
    selected_answers: { 1: ["ann", "bob"], 2: ["cy"] },
    userParticipants: 4,
    skippedParticipants: 1,
    avg_response_time: 2500,
  },
  {
    question: "Favourite colour?",
    question_type: "survey",
    options: { 1: "Red", 2: "Blue" },
    correct_answer: [1, 2],
    selected_answers: { 2: ["ann"] },
    userParticipants: 1,
    skippedParticipants: 0,
    avg_response_time: 1000,
  },
];
const participants = [
  { name: "Ann", rank: 1, score: 950, correct: 1, incorrect: 0, skipped: 0 },
  { name: "Bob", rank: 2, score: 700, correct: 1, incorrect: 0, skipped: 0 },
];

describe("report PDF", () => {
  it("tabulates each question's options, counts and time", () => {
    const { questions: q, participants: p } = reportPdfTables({
      questions,
      participants,
    });
    expect(q.body[0]).toEqual([
      1,
      "Capital of France?",
      "A. Paris (correct): 2\nB. Lyon: 1",
      2,
      1,
      1,
      "2.50 s",
    ]);
    // A survey has no correct option to mark.
    expect(q.body[1][2]).toBe("A. Red: 0\nB. Blue: 1");
    expect(p.body).toEqual([
      [1, "Ann", 950, 1, 0, 0],
      [2, "Bob", 700, 1, 0, 0],
    ]);
  });

  it("names the file after the quiz title and date", () => {
    const date = new Date("2026-11-21T10:00:00Z");
    expect(reportPdfFileName("Week 3: Maps & Rivers!", date)).toBe(
      "week-3-maps-rivers-report-2026-11-21.pdf"
    );
    expect(reportPdfFileName("", date)).toBe("quiz-report-2026-11-21.pdf");
  });

  it("writes the title, questions and participants into the PDF", async () => {
    const doc = await buildReportPdf({
      title: "Geography",
      summary: "2 questions, 4 players",
      questions,
      participants,
    });
    const pdf = doc.output();
    expect(pdf.startsWith("%PDF-")).toBe(true);
    for (const text of [
      "Geography",
      "Capital of France?",
      "A. Paris \\(correct\\): 2",
      "Participants",
      "Ann",
      "Bob",
    ]) {
      expect(pdf).toContain(text);
    }
  });
});
