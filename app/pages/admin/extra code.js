// Initialize data structures to collect analysis
let questionAnalysis = {};
let questionCount = {};
let totalResponseTime = {};

// Iterate through each response
data.data.forEach((response) => {
  let question = response.question;
  let options = response.options;
  let correctAnswer = response.correct_answer;
  let selectedAnswer = response.selected_answer.String;

  // Initialize question analysis if not already initialized
  if (!questionAnalysis[question]) {
    questionAnalysis[question] = {
      options: [],
    };
    questionCount[question] = 0;
    totalResponseTime[question] = 0;
  }

  // Calculate average response time per question
  if (response.response_time !== -1) {
    totalResponseTime[question] += response.response_time;
    questionCount[question]++;
  }

  // Initialize option statistics for the current question
  let optionStats = {};

  // Populate option statistics
  Object.keys(options).forEach((key) => {
    let optionText = options[key];
    let isCorrect = key === correctAnswer.slice(1, -1);
    let count = selectedAnswer === `[${key}]` ? 1 : 0;

    optionStats[key] = {
      text: optionText,
      correct: isCorrect,
      count: count,
    };
  });

  // Add option statistics to the current question analysis
  questionAnalysis[question].options.push({
    stats: optionStats,
  });
});

// Calculate average response time per question
Object.keys(questionAnalysis).forEach((question) => {
  if (questionCount[question] > 0) {
    questionAnalysis[question].average_response_time =
      totalResponseTime[question] / questionCount[question];
  } else {
    questionAnalysis[question].average_response_time = 0;
  }
});

// Prepare final analysis object
let analysisResult = {
  total_correct_answers: 0,
  total_incorrect_answers: 0,
  question_analysis: [],
};

// Populate total correct and incorrect answers
data.data.forEach((response) => {
  if (response.selected_answer.String === response.correct_answer) {
    analysisResult.total_correct_answers++;
  } else {
    analysisResult.total_incorrect_answers++;
  }
});

// Format question analysis into the final structure
Object.keys(questionAnalysis).forEach((question) => {
  let options = [];
  questionAnalysis[question].options.forEach((option) => {
    let stats = [];
    Object.keys(option.stats).forEach((key) => {
      stats.push({
        text: option.stats[key].text,
        correct: option.stats[key].correct,
        count: option.stats[key].count,
      });
    });
    options.push(stats);
  });
  analysisResult.question_analysis.push({
    question: question,
    options: options,
    average_response_time: questionAnalysis[question].average_response_time,
  });
});

// Convert the analysis result to JSON string
let jsonResult = JSON.stringify(analysisResult, null, 2);
console.log(jsonResult);
