<template>
  <div class="container-fluid quiz-container">
    <header class="text-center py-4">
      <h1>Quiz Analysis Report</h1>
    </header>

    <!-- Class Accuracy Bar -->
    <!-- <div class="quiz-header mb-4">
      <div class="quiz-accuracy position-relative w-100">
        <div class="progress">
          <div
            class="progress-bar bg-success"
            role="progressbar"
            :style="{ width: correctWidth + '%' }"
            aria-valuenow="correctWidth"
            aria-valuemin="0"
            aria-valuemax="100"
          ></div>
          <div
            class="progress-bar bg-danger"
            role="progressbar"
            :style="{ width: incorrectWidth + '%' }"
            aria-valuenow="incorrectWidth"
            aria-valuemin="0"
            aria-valuemax="100"
          ></div>
          <div
            class="progress-bar bg-secondary"
            role="progressbar"
            :style="{ width: unattemptedWidth + '%' }"
            aria-valuenow="unattemptedWidth"
            aria-valuemin="0"
            aria-valuemax="100"
          ></div>
        </div>
        <div
          class="accuracy-label position-absolute top-50 start-50 translate-middle"
        >
          71% accuracy
        </div>
      </div>
    </div> -->

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
        <div v-if="selectedTab === 'overview'" class="text-center">
          <QuizUserAnalyticsSpace
            v-for="(oData, index) in rankData"
            :key="index"
            :data="userJson[oData]"
            :survey-questions="surveyQuestions"
            class="user-analytics-item"
            @click="openPopup"
          ></QuizUserAnalyticsSpace>

          <div v-if="popup" class="full-page-popup">
            <QuizPopupUserwiseAnalysis
              v-if="popup"
              @close="closePopup"
            ></QuizPopupUserwiseAnalysis>
          </div>
        </div>
        <div v-if="selectedTab === 'questions'" class="text-center">
          <Frame
            v-for="(qData, index) in questionJson"
            :key="index"
            :page-title="'Q.' + index"
          >
          </Frame>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import lodash from "lodash";

const url = useState("urls");
const headers = useRequestHeaders(["cookie"]);

const analysisJson = ref([]);
// const correctWidth = ref(70); // percentage of correct answers
// const incorrectWidth = ref(20); // percentage of incorrect answers
// const unattemptedWidth = ref(10); // percentage of unattempted answers

const selectedTab = ref("overview");

const selectTab = (tab) => {
  selectedTab.value = tab;
};

const route = useRoute();
const activeQuizId = ref("");

const popup = ref(false);

const userJson = ref({});
const questionJson = ref({});
const rankData = ref([]);

import { useUserScoreboardData } from "~/store/userScoreboardData";
const userScoreboardDataStore = useUserScoreboardData();
const { getUserScoreboardData } = userScoreboardDataStore;

let storedData = {};
const surveyQuestions = ref([]);

const fetchData = () => {
  storedData = getUserScoreboardData();
};

const getAnalysisJson = async (activeQuizId) => {
  try {
    const response = await fetch(
      `${url.value.api_url}/analytics_board/admin?active_quiz_id=${activeQuizId}`,
      {
        method: "GET",
        headers: headers.value,
        mode: "cors",
        credentials: "include",
      }
    );

    const result = await response.json();

    if (response.ok) {
      analysisJson.value = result.data;
      userJson.value = lodash.groupBy(analysisJson.value, "username");

      // Iterate through each item in storedData
      storedData.forEach((data) => {
        rankData.value.push(data.username); //to get usernames rank wise, to pass data from userJson in sorted manner
        let key = data.username; // Get the username (key)

        // Check if the key exists in userJson.value
        if (userJson.value.hasOwnProperty(key)) {
          let totalScore = data.score; // Calculate total_score as score

          // Update userJson.value[key] with rank, total_score, and response_time
          userJson.value[key].push({
            rank: data.rank,
            total_score: totalScore,
            response_time: data.response_time,
          });
        } else {
          console.error(`Key '${key}' not found in userJson.value.`);
        }
      });
      questionJson.value = lodash.groupBy(analysisJson.value, "question");

      let questionNumber = 0;

      for (const key in questionJson.value) {
        questionJson.value[key].forEach((question) => {
          questionNumber++;
          const optionsCount = Object.keys(question.options).length;
          const correctAnswersCount = JSON.parse(
            question.correct_answer
          ).length;

          if (optionsCount === correctAnswersCount) {
            surveyQuestions.value.push(questionNumber);
          }
        });
      }
    } else {
      console.error(result);
    }
  } catch (error) {
    console.error("Failed to fetch data", error);
  }
};

onMounted(() => {
  activeQuizId.value = route.query.active_quiz_id || "";
  getAnalysisJson(activeQuizId.value);
  fetchData();
});

function openPopup() {
  popup.value = true;
}

function closePopup() {
  popup.value = false;
}
</script>

<style scoped>
body,
.quiz-container {
  padding: 0;
  margin: 0;
}

.quiz-container {
  font-family: Arial, sans-serif;
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
  width: 60%;
}

.progress {
  height: 30px;
  border-radius: 15px;
  overflow: hidden;
}

.accuracy-label {
  font-size: 1.25em;
  font-weight: bold;
  color: white;
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
</style>
