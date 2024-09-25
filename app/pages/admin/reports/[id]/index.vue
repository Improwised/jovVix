<template>
  <PageLayout :current-tab="currentTab" @change-tab="changeTab" />
  <template v-if="currentTab == 'report'">
    <div v-if="quizAnalysisPending">Loading...</div>
    <div
      v-else-if="quizAnalysisError"
      class="text-danger alert alert-danger mt-3"
    >
      Error while fetching data:
      <span>
        {{ quizAnalysisError }}
      </span>
    </div>
    <div v-else class="container mt-3">
      <div
        v-for="(quiz, index) in quizAnalysis.data"
        :key="index"
        class="card mb-3 row"
      >
        <QuizQuestionAnalysis
          :question="quiz"
          :order="index + 1"
          :is-admin-analysis="true"
        />
        <QuizOptionsAnalysis
          :options="quiz?.options"
          :correct-answer="quiz?.correct_answer"
          :selected-answer="quiz?.selected_answer?.String"
          :selected-answers="quiz.selected_answers"
          :options-media="quiz?.options_media"
          :is-admin-analysis="true"
        />
      </div>
    </div>
  </template>

  <template v-if="currentTab == 'participants'">
    <div v-if="quizUserAnalysispending">Loading...</div>
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
</template>

<script setup>
import PageLayout from "~~/components/reports/PageLayout.vue";
import lodash from "lodash";
definePageMeta({
  layout: "default",
});

const { api_url } = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const route = useRoute();
const activeQuizId = computed(() => route.params.id);
const currentTab = ref("report");

const {
  data: quizAnalysis,
  error: quizAnalysisError,
  pending: quizAnalysisPending,
} = useFetch(`${api_url}/admin/reports/${activeQuizId.value}/analysis`, {
  transform: (quizAnalysis) => {
    quizAnalysis.data?.map((quiz) => {
      quiz.userParticipants = Object.keys(quiz.selected_answers).length;
      const result = {};

      for (const [user, answer] of Object.entries(quiz.selected_answers)) {
        if (!result[answer]) {
          result[answer] = [];
        }
        result[answer].push(user);
      }

      quiz.selected_answers = result;
      quiz.correctPercentage =
        (quiz.correct_answer.reduce(
          (sum, correct_answer) =>
            (sum += quiz.selected_answers[correct_answer]?.length || 0),
          0
        ) /
          quiz.userParticipants) *
        100;

      return quiz;
    });
    return quizAnalysis;
  },
  credentials: "include",
  headers: headers,
});

const analysisJson = ref([]);
const userJson = ref({});
const questionJson = ref({});
const rankData = ref([]);
const surveyQuestions = ref(0);
const ranks = ref();
const quizUserAnalysispending = ref(false);
const fetchError = ref("");

const getAnalysisJson = async () => {
  try {
    quizUserAnalysispending.value = true;
    const response = await fetch(
      `${api_url}/analytics_board/admin?active_quiz_id=${activeQuizId.value}`,
      {
        method: "GET",
        headers: headers.value,
        mode: "cors",
        credentials: "include",
      }
    );

    const ranksResponse = await fetch(
      `${api_url}/final_score/admin?active_quiz_id=${activeQuizId.value}`,
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

    quizUserAnalysispending.value = false;
  } catch (error) {
    quizUserAnalysispending.value = false;
    fetchError.value = error;
    console.error("Failed to fetch data", error);
  }
};

onMounted(() => {
  getAnalysisJson();
});

const changeTab = (data) => {
  currentTab.value = data;
};
</script>

<style scoped>
.option-box {
  min-height: 70px;
  padding-top: 3px;
  border-radius: 30px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 5px;
  margin-top: 3px;
}

.wrong-option {
  border: 1px solid var(--bs-light-primary);
}
/* Adjust the divider for large screens */
</style>
