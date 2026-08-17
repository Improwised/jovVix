// Script ranges for the non-Latin languages the generator offers. Latin script
// covers English, Spanish, French, German and Portuguese alike, so those cannot
// be told apart by characters and fall back to English.
const SCRIPTS = [
  { language: "japanese", pattern: /[぀-ヿ]/ },
  { language: "chinese", pattern: /[一-鿿]/ },
  { language: "gujarati", pattern: /[઀-૿]/ },
  { language: "bengali", pattern: /[ঀ-৿]/ },
  { language: "tamil", pattern: /[஀-௿]/ },
  { language: "telugu", pattern: /[ఀ-౿]/ },
  { language: "arabic", pattern: /[؀-ۿ]/ },
  // Hindi and Marathi share Devanagari; Hindi is the likelier of the two.
  { language: "hindi", pattern: /[ऀ-ॿ]/ },
];

const DEFAULT_LANGUAGE = "english";

export const detectQuizLanguage = (texts) => {
  const sample = (Array.isArray(texts) ? texts : [texts])
    .filter(Boolean)
    .slice(0, 5)
    .join(" ");

  if (!sample) return DEFAULT_LANGUAGE;

  const counts = SCRIPTS.map(({ language, pattern }) => {
    const global = new RegExp(pattern.source, "g");
    return { language, count: (sample.match(global) || []).length };
  });

  const winner = counts.reduce((best, entry) =>
    entry.count > best.count ? entry : best
  );

  return winner.count > 0 ? winner.language : DEFAULT_LANGUAGE;
};
