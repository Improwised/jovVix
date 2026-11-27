// The quiz report as a PDF (GitHub #38), built in the browser from the data the
// report page already loaded, so it needs no extra API call or access check.
//
// jsPDF's built-in fonts hold Latin-1 only, and jsPDF does no text shaping, so even
// an embedded Gujarati or Devanagari font would break conjuncts and vowel signs.
// Text outside Latin-1 is drawn by the browser on a canvas, with the fonts the page
// itself uses, and placed in the PDF as an image. Latin text stays searchable.

// The analysis API sends the question type as a number (api/constants: Survey = 2).
const SURVEY = 2;
const isSurvey = (quiz) => Number(quiz.type) === SURVEY;

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
        !isSurvey(quiz) && quiz.correct_answer?.includes(Number(key))
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
      // A survey has no right answer: show its votes and skips, ungraded.
      isSurvey(quiz) ? "-" : correctCountFor(quiz),
      isSurvey(quiz) ? "-" : incorrectCountFor(quiz),
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

const LATIN1 = /^[\u0000-ÿ]*$/;
// Canvas pixels per PDF point, so the text images stay sharp when zoomed.
const SCALE = 4;

// Measures and draws text with the browser, which shapes every script it can show.
// Returns null without a 2D canvas (server render, some test DOMs).
export const canvasTextRenderer = () => {
  if (typeof document === "undefined") return null;
  const canvas = document.createElement("canvas");
  const ctx = canvas.getContext?.("2d");
  if (!ctx) return null;
  const family = getComputedStyle(document.body).fontFamily || "sans-serif";
  const font = (size) => `${size * SCALE}px ${family}`;
  const width = (text, size) => {
    ctx.font = font(size);
    return ctx.measureText(text).width / SCALE;
  };
  return {
    width,
    image(text, size, height) {
      canvas.width = Math.max(1, Math.ceil(width(text, size) * SCALE));
      canvas.height = Math.ceil(height * SCALE);
      // Resizing the canvas resets its state, so set the font again.
      ctx.font = font(size);
      ctx.fillStyle = "#141414";
      ctx.textBaseline = "middle";
      ctx.fillText(text, 0, canvas.height / 2);
      return canvas.toDataURL("image/png");
    },
  };
};

// shortcut: breaks lines at spaces only; a single word wider than the cell
// overflows it. Split by grapheme if long unspaced names show up.
const wrapLines = (text, maxWidth, measure) =>
  text.split("\n").flatMap((para) => {
    const lines = [];
    let line = "";
    for (const word of para.split(" ")) {
      const next = line ? `${line} ${word}` : word;
      if (line && measure(next) > maxWidth) {
        lines.push(line);
        line = word;
      } else line = next;
    }
    lines.push(line);
    return lines;
  });

// Returns a jsPDF document; the caller saves it.
export async function buildReportPdf({
  title,
  summary,
  questions,
  participants,
  textRenderer = canvasTextRenderer(),
}) {
  const [{ jsPDF }, { autoTable }] = await Promise.all([
    import("jspdf"),
    import("jspdf-autotable"),
  ]);
  const tables = reportPdfTables({ questions, participants });
  const doc = new jsPDF({ unit: "pt", format: "a4" });
  const left = 40;
  const lineHeight = (size) => size * doc.getLineHeightFactor();

  const drawImageText = (text, x, top, size) => {
    if (!text) return;
    const height = lineHeight(size);
    doc.addImage(
      textRenderer.image(text, size, height),
      "PNG",
      x,
      top,
      textRenderer.width(text, size),
      height
    );
  };
  // y is the baseline, as for doc.text.
  const drawText = (text, x, y, size) => {
    doc.setFontSize(size);
    if (!textRenderer || LATIN1.test(text)) doc.text(text, x, y);
    else drawImageText(text, x, y - size, size);
  };

  drawText(title || "Quiz report", left, 50, 18);
  drawText(summary, left, 70, 10);

  const heading = (text, y) => {
    doc.setFontSize(13);
    doc.text(text, left, y);
    return y + 8;
  };

  // Cells with text outside Latin-1: autoTable lays out blank lines, then the
  // browser-rendered lines are drawn over them.
  const imageCells = new Map();
  const imageTextHooks = textRenderer && {
    didParseCell: ({ cell }) => {
      const text = cell.text.join("\n");
      if (LATIN1.test(text)) return;
      const size = cell.styles.fontSize;
      const lines = wrapLines(
        text,
        cell.styles.cellWidth - cell.padding("horizontal"),
        (t) => textRenderer.width(t, size)
      );
      imageCells.set(cell, lines);
      cell.text = lines.map(() => "");
    },
    didDrawCell: ({ cell }) => {
      const size = cell.styles.fontSize;
      imageCells
        .get(cell)
        ?.forEach((line, i) =>
          drawImageText(
            line,
            cell.x + cell.padding("left"),
            cell.y + cell.padding("top") + i * lineHeight(size),
            size
          )
        );
    },
  };

  const style = {
    styles: { fontSize: 9, cellPadding: 4, overflow: "linebreak" },
    headStyles: { fillColor: [30, 30, 30] },
    margin: { left, right: left },
    ...imageTextHooks,
  };

  // Columns that hold quiz or player text get a fixed width, so non-Latin
  // text can be wrapped before autoTable lays out the rows.
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
    columnStyles: { 1: { cellWidth: 180 } },
  });
  return doc;
}
