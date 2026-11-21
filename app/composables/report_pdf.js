// The quiz report as a PDF (GitHub #38), built in the browser from the data the
// report page already loaded, so it needs no extra API call or access check.
// shortcut: jsPDF's built-in fonts cover Latin text only. Embed a Unicode font
// when quizzes in other scripts need a readable PDF.

// Players who picked a correct option of an analysed question.
export const correctCountFor = (quiz) =>
  quiz.correct_answer.reduce(
    (sum, idx) => sum + (quiz.selected_answers[idx]?.length || 0),
    0
  );

export const incorrectCountFor = (quiz) =>
  Math.max(
    (quiz.userParticipants || 0) -
      (quiz.skippedParticipants || 0) -
      correctCountFor(quiz),
    0
  );

const optionLines = (quiz) =>
  Object.entries(quiz.options || {})
    .map(([key, text]) => {
      const label = String.fromCharCode(64 + Number(key));
      const shown = quiz.options_media === "image" ? "[image]" : text;
      const correct =
        quiz.question_type !== "survey" &&
        quiz.correct_answer?.includes(Number(key))
          ? " (correct)"
          : "";
      const players = quiz.selected_answers?.[key]?.length || 0;
      return `${label}. ${shown}${correct}: ${players}`;
    })
    .join("\n");

const seconds = (ms) => `${(Math.abs(ms || 0) / 1000).toFixed(2)} s`;

// questions: the report page's transformed analysis rows.
// participants: [{ name, rank, score, correct, incorrect, skipped }] in rank order.
export const reportPdfTables = ({ questions, participants }) => ({
  questions: {
    head: [
      [
        "#",
        "Question",
        "Options: players",
        "Correct",
        "Incorrect",
        "Skipped",
        "Avg. time",
      ],
    ],
    body: questions.map((quiz, i) => [
      i + 1,
      quiz.question_media === "image"
        ? `${quiz.question} [image]`
        : quiz.question,
      optionLines(quiz),
      correctCountFor(quiz),
      incorrectCountFor(quiz),
      quiz.skippedParticipants || 0,
      seconds(quiz.avg_response_time),
    ]),
  },
  participants: {
    head: [["Rank", "Player", "Score", "Correct", "Incorrect", "Skipped"]],
    body: participants.map((p) => [
      p.rank,
      p.name,
      p.score,
      p.correct,
      p.incorrect,
      p.skipped,
    ]),
  },
});

export const reportPdfFileName = (title, date = new Date()) => {
  const slug =
    (title || "quiz")
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, "-")
      .replace(/^-+|-+$/g, "")
      .slice(0, 60) || "quiz";
  return `${slug}-report-${date.toISOString().slice(0, 10)}.pdf`;
};

// Returns a jsPDF document; the caller saves it.
export async function buildReportPdf({
  title,
  summary,
  questions,
  participants,
}) {
  const [{ jsPDF }, { autoTable }] = await Promise.all([
    import("jspdf"),
    import("jspdf-autotable"),
  ]);
  const tables = reportPdfTables({ questions, participants });
  const doc = new jsPDF({ unit: "pt", format: "a4" });
  const left = 40;

  doc.setFontSize(18);
  doc.text(title || "Quiz report", left, 50);
  doc.setFontSize(10);
  doc.text(summary, left, 70);

  const heading = (text, y) => {
    doc.setFontSize(13);
    doc.text(text, left, y);
    return y + 8;
  };
  const style = {
    styles: { fontSize: 9, cellPadding: 4, overflow: "linebreak" },
    headStyles: { fillColor: [30, 30, 30] },
    margin: { left, right: left },
  };

  autoTable(doc, {
    ...style,
    ...tables.questions,
    startY: heading("Questions", 100),
    columnStyles: { 1: { cellWidth: 150 }, 2: { cellWidth: 150 } },
  });
  autoTable(doc, {
    ...style,
    ...tables.participants,
    startY: heading("Participants", doc.lastAutoTable.finalY + 30),
  });
  return doc;
}
