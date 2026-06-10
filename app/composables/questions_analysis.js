export function questionsAnalysis(data) {
  let questionAnalysis = {
    rank: 0,
    totalScore: 0,
    totalQuestions: 0,
    correctAnwers: 0,
    wrongAnwers: 0,
    totalSurveyQuestions: 0,
    attemptedSurveyQuestions: 0,
    unAttemptedQuestions: 0,
    accuracy: 0,
  };

  if (!Array.isArray(data)) return questionAnalysis;

  // A leading entry carries the user's rank/total_score; subsequent entries
  // are question rows. We accumulate per-question score on top of the seed
  // total_score so the final value matches the server-side aggregate.
  data.forEach((item) => {
    if (!item) return;
    if (item.rank) {
      questionAnalysis.rank = Number(item.rank) || 0;
      questionAnalysis.totalScore = Number(item.total_score) || 0;
      return;
    }
    questionAnalysis.totalQuestions++;

    if (!item.is_attend) {
      if (item.question_type == "survey") {
        questionAnalysis.totalSurveyQuestions++;
      }
      questionAnalysis.unAttemptedQuestions++;
    } else if (item.question_type != "survey" && item.is_attend) {
      const correctIncorrectFlag = isCorrectAnswer(
        item.selected_answer?.String,
        item.correct_answer
      );

      if (correctIncorrectFlag) {
        questionAnalysis.correctAnwers++;
      } else {
        questionAnalysis.wrongAnwers++;
      }
    } else if (item.question_type == "survey" && item.is_attend) {
      questionAnalysis.totalSurveyQuestions++;
      questionAnalysis.attemptedSurveyQuestions++;
    }

    const score = Number(item.calculated_score);
    if (Number.isFinite(score)) {
      questionAnalysis.totalScore += score;
    }
  });

  questionAnalysis.accuracy = questionAnalysis.totalQuestions
    ? Number(
        (
          ((questionAnalysis.attemptedSurveyQuestions +
            questionAnalysis.correctAnwers) /
            questionAnalysis.totalQuestions) *
          100
        ).toFixed(2)
      )
    : 0;

  return questionAnalysis;
}
