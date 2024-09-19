<template>
  <div class="container-fluid quiz-container">
    <header class="text-center py-4">
      <h1>Quiz Analysis</h1>
    </header>

    <div class="quiz-content">
      <ul class="nav nav-tabs justify-content-center mb-4">
        <li class="nav-item">
          <button
            class="nav-link"
            :class="{ active: selectedTab === 'overview' }"
            @click="selectTab('overview')"
          >
            Overview
          </button>
        </li>
        <li class="nav-item">
          <button
            class="nav-link"
            :class="{ active: selectedTab === 'questions' }"
            @click="selectTab('questions')"
          >
            Questions
          </button>
        </li>
      </ul>

      <div class="tab-content">
        <div
          v-if="selectedTab === 'overview'"
          class="text-center user-analytics-item"
        >
          <QuizUserAnalyticsSpace
            v-for="(oData, index) in rankData"
            :key="index"
            :data="userJson[oData]"
            :user-name="oData"
          ></QuizUserAnalyticsSpace>
        </div>
        <div v-if="selectedTab === 'questions'" class="text-center">
          <!-- Display Total Questions in Card -->
          <div class="card total-questions-card mx-auto mb-4 text-center">
            <div class="card-body">
              <h5 class="card-title">Total Questions</h5>
              <p class="card-text">{{ totalQuestions }}</p>
            </div>
          </div>

          <Frame
            v-for="(qData, index) in questionJson"
            :key="index"
            :page-title="`Q. ${qData[0] ? qData[0].order_no : ''} ${index}`"
            class="mb-2"
          >
            <div v-if="qData[0]?.question_media === 'image'" class="d-flex align-items-center justify-content-center">
              <img
                :src="`${qData[0]?.resource}`"
                :alt="`${qData[0]?.resource}`"
                class="rounded img-thumbnail"
              />
            </div>
            <div class="row m-2">
              <div
                class="d-flex flex-wrap align-items-center justify-content-between gap-2"
              >
                <span class="badge bg-primary">
                  AVG. Response Time:
                  {{ (qData[0].response_time / 1000).toFixed(2) }} seconds
                </span>
                <span
                  v-if="qData[0].question_type === 'mcq'"
                  class="badge bg-light-info mx-2 text-dark"
                >
                  Multiple Choice Question
                </span>
                <span v-else class="badge bg-light-info mx-2 text-dark">
                  Survey Question
                </span>
              </div>
            </div>

            <ul class="options-list">
              <li
                v-for="(option, key) in qData[0].options"
                :key="key"
                :class="
                  qData[0].correct_answer.includes(key) ? 'correct-answer' : ''
                "
              >
                <span v-if="qData[0]?.options_media === 'text'">{{ key }}: {{ option }}</span>
                <div v-if="qData[0]?.options_media === 'image'" class="d-flex align-items-center justify-content-center">
                  <span>{{ key }}:</span>
                  <img
                    :src="`${option}`"
                    :alt="`${option}`"
                    class="rounded img-thumbnail"
                  />
                </div>
              </li>
            </ul>
          </Frame>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRoute } from "vue-router";
import lodash from "lodash";

const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const route = useRoute();

const selectedTab = ref("overview");

const userJson = ref({});
const questionJson = ref({});
const rankData = ref([]);
const totalQuestions = ref(0);
const activeQuizId = computed(() => route.query.active_quiz_id);

const {
  data: adminFinalScore,
  pending: adminFinalScorePending,
  error: adminFinalScoreError,
} = useFetch(
  `${url.api_url}/final_score/admin?active_quiz_id=${activeQuizId.value}`,
  {
    headers: headers,
    mode: "cors",
    credentials: "include",
  }
);

const {
  data: adminAnalyticsBoard,
  pending: adminAnalyticsBoardPending,
  error: adminAnalyticsBoardError,
} = useFetch(
  `${url.api_url}/analytics_board/admin?active_quiz_id=${activeQuizId.value}`,
  {
    headers: headers,
    mode: "cors",
    credentials: "include",
  }
);

watch(
  [adminAnalyticsBoardPending, adminFinalScorePending],
  () => {
    if (adminAnalyticsBoardPending.value || adminFinalScorePending.value) {
      return;
    } else if (adminFinalScoreError.value || adminAnalyticsBoardError.value) {
      toast.error("error while get analysis");
      return;
    }

    let analysisJson = adminAnalyticsBoard.value.data;
    userJson.value = lodash.groupBy(analysisJson, "username");

    adminFinalScore.value?.data.forEach((data) => {
      rankData.value.push(data.username);
      let key = data.username;

      if (userJson.value.hasOwnProperty(key)) {
        let totalScore = data.score;

        userJson.value[key].push({
          rank: data.rank,
          total_score: totalScore,
          response_time: data.response_time,
        });
      } else {
        console.error(`Key '${key}' not found in userJson.value.`);
      }
    });

    questionJson.value = lodash.groupBy(analysisJson, "question");
    totalQuestions.value = Object.keys(questionJson.value).length;
    for (const key in questionJson.value) {
      questionJson.value[key].forEach((question) => {
        try {
          const correctAnswers = JSON.parse(question.correct_answer);
          question.question_type = correctAnswers.length > 1 ? "survey" : "mcq";
        } catch (error) {
          console.error(
            `Error parsing correct_answer for question: ${question.question}`,
            error
          );
          question.question_type = "mcq"; // defaulting to MCQ if error occurs
        }
      });
    }
  },
  { immediate: true, deep: true }
);

const selectTab = (tab) => {
  selectedTab.value = tab;
};
</script>

<style scoped>
.total-questions-card {
  max-width: 300px; /* Set a max width for larger screens */
  width: 100%; /* Ensure it scales down on smaller screens */
}

@media (max-width: 576px) {
  .total-questions-card {
    padding: 10px; /* Adjust padding on smaller screens */
  }
}

.options-list {
  padding: 0;
  list-style: none;
  margin: 0;
}

.options-list li {
  background-color: #f2f4f8;
  border-radius: 8px;
  padding: 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px; /* Add space between options */
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: background-color 0.3s ease;
}

.options-list li.correct-answer {
  background-color: #d4f5d4;
}

.option-text {
  font-size: 1rem;
  font-weight: 500;
}

.divider {
  width: 100%;
  height: 1px;
  background-color: #ccc;
  margin-top: 20px;
}

.total-questions-card {
  width: 210px; /* Fixed width */
  height: 70px; /* Fixed height to make it more square-like */
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 10px; /* Adjust padding as needed */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  display: flex; /* Use flexbox to center content */
  flex-direction: column; /* Stack elements vertically */
  justify-content: center; /* Center content vertically */
  align-items: center; /* Center content horizontally */
  text-align: center; /* Center text */
}

.card-title {
  font-size: 1rem; /* Adjusted font size */
  margin-bottom: 8px; /* Reduced margin */
  font-weight: bold;
  color: #4a2c77;
}

.card-text {
  font-size: 1.25rem; /* Adjusted font size */
  font-weight: bold;
  color: #4a2c77;
}

.user-row {
  background-color: #8968cd !important;
  box-shadow: none;
}

.quiz-container {
  font-family: Arial, sans-serif;
  padding: 20px;
}

header {
  margin-bottom: 20px;
}

.quiz-header {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 20px;
}

.quiz-accuracy {
  position: relative;
  width: 80%;
  margin-bottom: 30px; /* Adds space between the bar and the tabs */
}

.progress {
  height: 35px;
  border-radius: 20px;
  overflow: hidden;
}

.accuracy-label {
  font-size: 1.25em;
  font-weight: bold;
  color: white;
}

.nav-tabs {
  margin-top: 20px; /* Ensures there's space above the tabs */
}

.nav-tabs .nav-link {
  cursor: pointer;
}

.nav-tabs .nav-link.active {
  background-color: #007bff;
  color: white;
}

.tab-content {
  text-align: center;
  margin-top: 20px;
}

.user-analytics-item {
  margin-bottom: 20px; /* Adjust spacing between each user analytics item */
}

.progress-bar {
  display: flex;
  height: 30px;
  border-radius: 20px;
  overflow: hidden;
  width: 100%;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: white;
}
.correct {
  background-color: #4caf50;
}
.incorrect {
  background-color: #f44336;
}

.quiz-accuracy {
  position: relative;
  width: 80%;
  margin-bottom: 30px; /* Adds space between the bar and the tabs */
}

.progress {
  height: 40px;
  border-radius: 20px;
  overflow: hidden;
  position: relative; /* Ensure the circle is positioned relative to the progress bar */
}

.progress-bar {
  display: flex;
  height: 40px;
  border-radius: 20px;
  overflow: hidden;
  width: 100%;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: white;
}

.correct {
  background-color: #4caf50;
}

.incorrect {
  background-color: #f44336;
}

.progress-circle {
  position: absolute;
  top: -15px; /* Adjust this value to move the circle vertically */
  width: 70px;
  height: 70px;
  background-color: #fff;
  border: 2px solid #4caf50;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
  transition: left 1s ease-in-out; /* Animation for the left property */
}
</style>
