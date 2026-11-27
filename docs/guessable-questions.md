# Guessable question warnings

The quiz detail page (`/admin/quiz/list-quiz/[quiz_id]`) warns a creator when a question's correct
answer can be picked without knowing the subject (#342). The warning is advice. It does not block
saving or hosting, and only people who can edit the quiz see it.

The checks are in `app/composables/guessable_questions.js` (`findGuessablePatterns`; `guessWarningsFor`
applies the editors-only and English-only rules), with one test
per pattern in `app/test/composables/guessable_questions.test.js`. They run in the browser each time
the page shows the questions, so a saved edit is checked again at once.

## Patterns checked

| Code | The warning fires when | Why it gives the answer away |
| --- | --- | --- |
| `duplicate` | The correct answer has the same text as a wrong option, or two wrong options have the same text (case and punctuation ignored). | A repeated correct answer stands out; repeated wrong options can't both be right, so players drop both. |
| `longest` | The correct answer is at least 1.5 times as long as the longest wrong option and at least 10 characters longer. | Writers add detail to the true answer; players know to pick the longest, most detailed one. |
| `shortest` | The correct answer is at most two thirds as long as the shortest wrong option and at least 10 characters shorter. | A short, plain answer among long, padded ones stands out the same way. |
| `catch_all` | The correct answer is "All of the above", "None of these", "Both A and B" and similar. | Players pick these by habit, and they make no sense once options are shuffled. |
| `hedge` | Only the correct answer has a qualifier such as "usually", "often", "may" or "typically". | True statements tend to be hedged; invented ones are not. |
| `echo` | Only the correct answer repeats a content word (4+ letters, not a question word) from the question. | Players match the words in the question. |

The checks run on questions with exactly one correct answer, which the editor, CSV import and AI
generation all require. `longest` and `shortest` need at least one wrong option. `hedge` and `echo`
need at least two, so a true/false question is not flagged by them.

## When checks are skipped

- Survey questions have no correct answer, so no check runs.
- Image and code options run only `duplicate`.
- `catch_all`, `hedge` and `echo` read English words. They run only when
  the quiz language is detected as English. Latin-script quizzes in other languages are detected as
  English too, so these checks can miss or misfire there.

## Not checked

- Option position, such as "the third option is always correct". jovVix shuffles options at play
  time.
- Patterns across a whole quiz, such as the correct answer being the longest in most questions.
- Singular and plural agreement between the question and the options. It needs a parser to be
  reliable.
- Absolute words ("always", "never") on the wrong options only. A useful word list would be too long
  to keep up, and a short one misses most cases.
- A question ending in "a" or "an" that only the correct answer fits. It needs the sound of the next
  word, not its spelling ("an hour", "a university"), and such questions are rare.
- Options that are opposites of each other, one of them correct. Detecting opposites needs a word
  list per topic.
