<template>
  <PageLayout />
  <div v-if="pending">Loading...</div>
  <div v-else-if="fetchError" class="text-danger alert alert-danger mt-3">
    Error while fetching data:
    <span>
      {{ fetchError }}
    </span>
  </div>
  <div v-else class="quiz-content">
    <div class="tab-content">
      <QuizUserAnalyticsSpace
        v-for="(oData, index) in rankData"
        :key="index"
        :data="userJson[oData]"
        :user-name="oData"
        :survey-questions="surveyQuestions"
        class="user-analytics-item"
      ></QuizUserAnalyticsSpace>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import lodash from "lodash";

const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);

const analysisJson = ref([]);

const route = useRoute();
const activeQuizId = ref("");

const userJson = ref({});
const questionJson = ref({});
const rankData = ref([]);
const ranks = ref();
const pending = ref(false);
const fetchError = ref("");
import PageLayout from "~~/components/reports/PageLayout.vue";

const surveyQuestions = ref(0);

const getAnalysisJson = async (activeQuizId) => {
  try {
    pending.value = true;
    const response = await fetch(
      `${url.api_url}/analytics_board/admin?active_quiz_id=${activeQuizId}`,
      {
        method: "GET",
        headers: headers.value,
        mode: "cors",
        credentials: "include",
      }
    );

    const ranksResponse = await fetch(
      `${url.api_url}/final_score/admin?active_quiz_id=${activeQuizId}`,
      {
        method: "GET",
        headers: headers.value,
        mode: "cors",
        credentials: "include",
      }
    );

    const result = await response.json();
    ranks.value = await ranksResponse.json();

    if (response.ok && ranksResponse.ok) {
      analysisJson.value = result.data;

      userJson.value = lodash.groupBy(analysisJson.value, "username");

      ranks.value?.data.forEach((data) => {
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

      // from userJson, count total points of all questions and count of total survey questions
      for (const key in userJson.value) {
        userJson.value[key].forEach((question) => {
          if (!question.rank) {
            if (question.question_type == "survey") {
              surveyQuestions.value++;
            }
          }
        });
        break;
      }
    } else {
      console.error(result);
    }

    pending.value = false;
  } catch (error) {
    pending.value = false;
    fetchError.value = error;
    console.error("Failed to fetch data", error);
  }
};

onMounted(() => {
  activeQuizId.value = route.params.id || "";
  getAnalysisJson(activeQuizId.value);
});
</script>

<style scoped>
body,
.quiz-container {
  padding: 0;
  margin: 0;
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
