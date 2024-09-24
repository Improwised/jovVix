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
  data.filter((item) => {
    if (item.rank) {
      questionAnalysis.rank = item.rank;
      questionAnalysis.totalScore = item.total_score;
      return;
    }
    questionAnalysis.totalQuestions++;
    let correctIncorrectFlag = false;

    if (!item.is_attend) {
      if (item.question_type == "survey") {
        questionAnalysis.totalSurveyQuestions++;
      }
      questionAnalysis.unAttemptedQuestions++;
    } else if (item.question_type != "survey" && item.is_attend) {
      //check if the answer is correct or not
      correctIncorrectFlag = isCorrectAnswer(
        item.selected_answer.String,
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
    questionAnalysis.totalScore += item.calculated_score;
  });

  questionAnalysis.accuracy = Number(
    (
      ((questionAnalysis.attemptedSurveyQuestions +
        questionAnalysis.correctAnwers) /
        questionAnalysis.totalQuestions) *
      100
    ).toFixed(2)
  );

  return questionAnalysis;
}
