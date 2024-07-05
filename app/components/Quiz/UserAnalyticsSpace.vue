<script setup>
const props = defineProps({
  data: {
    type: Object,
    required: true,
  },
  surveyQuestions: {
    type: Array,
    required: true,
  },
});

let totalScore = 0;
let accuracyArr = [];
let accuracy = 0;
let rank = 0;

const totalSurveyQuestions = props.surveyQuestions.length;

const answeredSurvey = ref(0); // count of the surveys user attempted
let questionCount = 1; // to get the question number

props.data.forEach(function (arrayItem) {
  // if it is the survey question then apply this block's logic
  if (props.surveyQuestions.includes(questionCount)) {
    if (arrayItem.is_attend) {
      answeredSurvey.value++;
    }
  } else {
    // if it is not the survey question then apply this block's logic

    //if it is not the last object which includes the rank, response time and score information then procced to checking if it is attempted or not
    if (!arrayItem.rank) {
      // check if it is attempted or not, if attempted then check correct or incorrect and push accordingly, correct - true, incorrect - false
      if (arrayItem.is_attend) {
        accuracyArr.push(
          isCorrectAnswer(
            arrayItem.selected_answer.String,
            arrayItem.correct_answer
          )
        );
      } else {
        // if it is unattempted then push 0
        accuracyArr.push(0);
      }
    } else {
      //it it is the last object of rank, fetch rank and total score from it
      totalScore = arrayItem.total_score;
      rank = arrayItem.rank;
    }
  }
  questionCount++; // increase the count to use it as question number
});

const countTrue = accuracyArr.filter(Boolean).length;
const notAttempted = accuracyArr.filter((value) => value === 0).length;
const countFalse = accuracyArr.length - countTrue - notAttempted;
accuracy = ((countTrue / accuracyArr.length) * 100).toFixed(2);

//function to check if answer provided by user in all questions are correct or not
function isCorrectAnswer(selectedAnswer, correctAnswer) {
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

const correctWidth = (countTrue / accuracyArr.length) * 100; // percentage of correct answers
const unattemptedWidth = (notAttempted / accuracyArr.length) * 100; // percentage of unattempted answers
const incorrectWidth = 100 - correctWidth - unattemptedWidth; // percentage of incorrect answers

const handleMouseEnter = (event) => {
  event.target.style.transform = "scale(1.05)"; // Scale up on hover
  event.target.style.transition = "transform 0.3s ease"; // Add transition
};

const handleMouseLeave = (event) => {
  event.target.style.transform = "scale(1)"; // Reset scale on leave
};
</script>

<template>
  <div
    class="user-analytics-item"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
  >
    <div class="user-stats-box">
      <div class="header">
        <div class="avatar-container">
          <img
            class="avatar"
            src="../../assets/images/avatar.png"
            alt="Avatar"
          />
          <div class="name">{{ props.data[0].username }}</div>
        </div>
        <div class="stats">
          <div class="stat-item">
            <span class="value">{{ rank }}</span>
            <span class="label">Rank</span>
          </div>
          <div class="stat-item">
            <span class="value">{{ accuracy }}%</span>
            <span class="label">Accuracy</span>
          </div>
          <div class="stat-item">
            <span class="value">{{ totalScore }}</span>
            <span class="label">Score</span>
          </div>
        </div>
      </div>
      <div class="quiz-header mb-4">
        <div class="quiz-accuracy position-relative w-100">
          <div class="progress">
            <div
              class="progress-bar bg-success"
              role="progressbar"
              :style="{ width: correctWidth + '%' }"
              aria-valuenow="70"
              aria-valuemin="0"
              aria-valuemax="100"
            ></div>
            <div
              class="progress-bar bg-danger"
              role="progressbar"
              :style="{ width: incorrectWidth + '%' }"
              aria-valuenow="20"
              aria-valuemin="0"
              aria-valuemax="100"
            ></div>
            <div
              class="progress-bar bg-secondary"
              role="progressbar"
              :style="{ width: unattemptedWidth + '%' }"
              aria-valuenow="10"
              aria-valuemin="0"
              aria-valuemax="100"
            ></div>
          </div>
        </div>
      </div>
      <div>
        &#9989; {{ countTrue }} &ensp; &#10060; {{ countFalse }} &ensp; &#x25CC;
        {{ notAttempted }} &ensp;
        <span v-if="totalSurveyQuestions > 0"
          >&#128203; {{ answeredSurvey }} / {{ totalSurveyQuestions }}</span
        >
      </div>
      <div class="progress-bar">
        <div
          class="correct"
          style="
             {
              width: correctWidth + '%';
            }
          "
        ></div>
        <div
          class="incorrect"
          style="
             {
              width: incorrectWidth + '%';
            }
          "
        ></div>
      </div>
      <div class="progress-bar">
        <div
          class="correct"
          style="
             {
              width: correctWidth + '%';
            }
          "
        ></div>
        <div
          class="incorrect"
          style="
             {
              width: incorrectWidth + '%';
            }
          "
        ></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.user-stats-box {
  border: 1px solid #ddd;
  padding: 10px;
  border-radius: 8px;
  width: 100%;
  max-width: 600px;
  background-color: white;
  margin: 0 auto;
  box-sizing: border-box;
  box-shadow: 0px 2px 4px rgba(0, 0, 0, 0.1); /* Box shadow added */
}
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}
.avatar-container {
  display: flex;
  align-items: center;
}
.avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
}
.name {
  font-size: 18px;
  font-weight: bold;
  margin-left: 10px;
}
.stats {
  display: flex;
  flex-wrap: nowrap;
  justify-content: flex-end;
}
.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-left: 20px;
}
.label {
  font-size: 12px;
  color: #888;
}
.value {
  font-size: 14px;
  font-weight: bold;
}
.divider {
  width: 100%;
  border: none;
  height: 1px;
  background-color: #000;
  margin: 10px 0;
}
.progress-bar {
  display: flex;
  height: 10px;
  border-radius: 5px;
  overflow: hidden;
  width: 100%;
}
.correct {
  background-color: #4caf50;
}
.incorrect {
  background-color: #f44336;
}

@media (max-width: 600px) {
  .header {
    flex-wrap: wrap;
    justify-content: center;
  }
  .avatar-container {
    margin-bottom: 10px;
  }
  .stats {
    justify-content: center;
    flex-wrap: nowrap;
  }
  .stat-item {
    margin-left: 10px;
    margin-top: 0;
  }
}

.user-analytics-item {
  padding: 10px;
  margin-bottom: 10px;
  border: none; /* Remove border */
  border-radius: 5px;
  cursor: pointer;
  transition: transform 0.3s ease; /* Add transition for scale */
}

.user-analytics-item:hover {
  transform: scale(1.05); /* Scale up on hover */
}
</style>
