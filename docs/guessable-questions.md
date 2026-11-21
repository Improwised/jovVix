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
| `all_correct` | Every option is marked correct. | Any choice scores. |
| `duplicate` | Two wrong options have the same text (case and punctuation ignored). | Only one option can be right, so players drop both. |
| `longest` | The correct answer is at least 1.5 times as long as the longest wrong option and at least 10 characters longer. | Writers add detail to the true answer; players know to pick the longest, most detailed one. |
| `catch_all` | The correct answer is "All of the above", "None of these", "Both A and B" and similar. | Players pick these by habit, and they make no sense once options are shuffled. |
| `hedge` | Only the correct answer has a qualifier such as "usually", "often", "may" or "typically". | True statements tend to be hedged; invented ones are not. |
| `absolutes` | Every wrong option says "always", "never", "all", "only" or similar, and the correct one does not. | Absolute statements are usually false, so players rule them out. |
| `echo` | Only the correct answer repeats a content word (4+ letters, not a question word) from the question. | Players match the words in the question. |
| `article` | The question ends in "a" or "an" and only the correct answer fits that article. | Grammar points to the one option that reads correctly. |

`longest` needs at least one wrong option. `hedge`, `absolutes` and `echo` need at least two, so a
true/false question is not flagged by them. Only `all_correct` applies when the question has more
than one correct answer.

## When checks are skipped

- Survey questions have no correct answer, so no check runs.
- Image and code options run only `all_correct` and `duplicate`.
- `catch_all`, `hedge`, `absolutes`, `echo` and `article` read English words. They run only when
  the quiz language is detected as English. Latin-script quizzes in other languages are detected as
  English too, so these checks can miss or misfire there.
- `article` judges by the first letter, not the sound: "an hour" and "a university" are misjudged.

## Not checked

- Option position, such as "the third option is always correct". jovVix shuffles options at play
  time.
- Patterns across a whole quiz, such as the correct answer being the longest in most questions.
- Singular and plural agreement between the question and the options. It needs a parser to be
  reliable.
- Options that are opposites of each other, one of them correct. Detecting opposites needs a word
  list per topic.
