// Checks for multiple-choice questions whose correct answer can be spotted without knowing
// the subject (GitHub #342). docs/guessable-questions.md lists each pattern and what is not
// checked. Option order is shuffled at play time, so no check looks at positions.

const SURVEY_TYPE = 2;

const CATCH_ALL =
  /^(all|none|both|neither) (of )?(the )?(above|these|them|options|answers)$|^both [a-e] and [a-e]$|^all (of )?the above$/;
const HEDGES = new Set([
  "usually",
  "typically",
  "generally",
  "often",
  "sometimes",
  "may",
  "might",
  "can",
  "could",
  "probably",
  "likely",
  "mostly",
  "frequently",
  "normally",
]);
const ABSOLUTES = new Set([
  "always",
  "never",
  "all",
  "none",
  "only",
  "every",
  "must",
  "completely",
]);
const STOPWORDS = new Set([
  "what",
  "which",
  "when",
  "where",
  "who",
  "whom",
  "whose",
  "why",
  "how",
  "this",
  "that",
  "these",
  "those",
  "there",
  "their",
  "they",
  "them",
  "with",
  "from",
  "into",
  "about",
  "does",
  "have",
  "been",
  "were",
  "will",
  "would",
  "should",
  "could",
  "following",
  "best",
  "most",
  "called",
  "known",
]);

const normalize = (text) =>
  String(text)
    .toLowerCase()
    .replace(/[^\p{L}\p{N}\s]/gu, " ")
    .replace(/\s+/g, " ")
    .trim();
const words = (text) => normalize(text).split(" ").filter(Boolean);
// shortcut: strips a plural "s" only; upgrade to a real stemmer if echo misses obvious matches.
const stem = (word) => (word.length > 4 ? word.replace(/s$/, "") : word);
const contentWords = (text) =>
  new Set(
    words(text)
      .filter((w) => w.length >= 4 && !STOPWORDS.has(w))
      .map(stem)
  );
const hasAny = (text, set) => words(text).some((w) => set.has(w));
// shortcut: spelling, not sound ("an hour", "a university" misjudge); fine for a warning.
const startsWithVowel = (text) => /^[aeiou]/.test(normalize(text));

const parseAnswers = (value) => {
  if (Array.isArray(value)) return value.map(Number);
  try {
    const parsed = JSON.parse(value);
    return Array.isArray(parsed) ? parsed.map(Number) : [];
  } catch {
    return [];
  }
};

// Returns [{ code, reason }]. Each reason completes "This question might be too easy to
// guess because ...". `english` turns on the checks that read English words.
export const findGuessablePatterns = (question, { english = true } = {}) => {
  if (!question || Number(question.question_type_id) === SURVEY_TYPE) return [];

  const answers = new Set(parseAnswers(question.correct_answer));
  const options = Object.entries(question.options || {})
    .filter(([, text]) => String(text ?? "").trim() !== "")
    .map(([key, text]) => ({ key: Number(key), text: String(text).trim() }));
  const correct = options.filter((o) => answers.has(o.key));
  const wrong = options.filter((o) => !answers.has(o.key));
  const found = [];
  const add = (code, reason) => found.push({ code, reason });

  if (options.length > 1 && wrong.length === 0) {
    add("all_correct", "every option is marked correct");
    return found;
  }
  if (correct.length !== 1 || wrong.length === 0) return found;

  const seen = new Map();
  for (const o of wrong) {
    const n = normalize(o.text);
    if (seen.has(n)) {
      add(
        "duplicate",
        "two wrong options have the same text, so players can rule both out"
      );
      break;
    }
    seen.set(n, o.key);
  }

  // Images and code are compared as text only for duplicates.
  if (question.options_media && question.options_media !== "text") return found;

  const answer = correct[0].text;
  const longestWrong = Math.max(...wrong.map((o) => o.text.length));
  if (
    answer.length >= longestWrong * 1.5 &&
    answer.length - longestWrong >= 10
  ) {
    add("longest", "the correct answer is much longer than the other options");
  }

  if (!english) return found;

  if (CATCH_ALL.test(normalize(answer))) {
    add(
      "catch_all",
      `the correct answer is "${answer}", which players often pick by habit`
    );
  }

  if (wrong.length >= 2) {
    if (hasAny(answer, HEDGES) && !wrong.some((o) => hasAny(o.text, HEDGES))) {
      add(
        "hedge",
        'only the correct answer is qualified with a word like "usually" or "often"'
      );
    }
    if (
      !hasAny(answer, ABSOLUTES) &&
      wrong.every((o) => hasAny(o.text, ABSOLUTES))
    ) {
      add(
        "absolutes",
        'every wrong option says "always", "never" or similar, and the correct one does not'
      );
    }
    const stemWords = contentWords(question.question);
    const echoes = (text) =>
      [...contentWords(text)].some((w) => stemWords.has(w));
    if (echoes(answer) && !wrong.some((o) => echoes(o.text))) {
      add("echo", "only the correct answer repeats a word from the question");
    }
  }

  const article = words(question.question).at(-1);
  if (article === "a" || article === "an") {
    const fits = (text) => startsWithVowel(text) === (article === "an");
    if (fits(answer) && !wrong.some((o) => fits(o.text))) {
      add(
        "article",
        `the question ends with "${article}", which only the correct answer fits`
      );
    }
  }

  return found;
};
