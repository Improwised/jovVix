//function to check if answer provided by user in all questions are correct or not
export function isCorrectAnswer(selectedAnswer, correctAnswer) {
  // Function to parse and sort the answer string into an array
  function parseAndSort(answer) {
    return answer.length <= 2
      ? []
      : answer
          .slice(1, -1)
          .split(",")
          .map(Number)
          .sort((a, b) => a - b);
  }

  // Parse and sort both answers
  const selectedArray = parseAndSort(selectedAnswer);
  const correctArray = parseAndSort(correctAnswer);

  // Check if selectedArray is not empty and every element in selectedArray is in correctArray
  return (
    selectedArray.length > 0 &&
    selectedArray.every((value) => correctArray.includes(value))
  );
}
