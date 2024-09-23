<template>
  <PageLayout :current-tab="currentTab" @change-tab="changeTab"/>
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
      <div class="card mb-3 row" v-for="(quiz, index) in quizAnalysis.data">
        <div class="row m-2">
          <div class="col-lg-6 col-sm-12">
            <div>
              <strong class="text-primary">Question: </strong>
              <h3 class="font-bold">{{ quiz.question }}</h3>
            </div>
          </div>
          <div v-if="quiz?.question_media === 'image'" class="d-flex align-items-center justify-content-center">
            <img
              :src="`${quiz?.resource}`"
              :alt="`${quiz?.resource}`"
              class="rounded img-thumbnail"
            />
          </div>
          <CodeBlockComponent v-if="quiz?.question_media === 'code'" :code="quiz?.resource" />
          <div
            class="col-lg-12 d-flex flex-wrap align-items-center justify-content-around"
          >
            <span class="bg-light-primary rounded px-2 text-dark">
              AVG. Response Time:
              {{ Math.abs((quiz.avg_response_time / 1000).toFixed(2)) }}/
              {{ quiz.duration }} seconds
            </span>
            <span v-if="quiz.type === 1" class="badge bg-light-info m-1 text-dark"
              >M.C.Q.</span
            >
            <span v-else class="badge bg-light-info m-1 text-dark">Survey</span>
            <v-progress-circular
              class="mt-2"
              :model-value="quiz.correctPercentage"
              :rotate="360"
              :size="60"
              :width="5"
              :color="quiz.correctPercentage >= 50 ? 'teal' : '#D2042D'"
            >
              {{ quiz.correctPercentage.toFixed(0) }}%
            </v-progress-circular>
          </div>
        </div>

        <div class="row d-flex align-items-stretch m-2">
          <div v-for="(option, order) in quiz.options" class="col-lg-6 col-md-12">
            <div
              v-if="quiz.correct_answer.includes(Number(order))"
              class="bg-light-success option-box"
            >
              <Option
                class="text-success font-bold"
                :order="order"
                :option="option"
                :selected="quiz.selected_answers[order]?.length || 0"
                icon="fa-solid fa-check"
                :isCorrect="true"
                :options-media="quiz?.options_media"
              />
            </div>
            <div v-else class="option-box wrong-option">
              <Option
                class=""
                :order="order"
                :option="option"
                :selected="quiz.selected_answers[order]?.length || 0"
                :options-media="quiz?.options_media"
              />
            </div>
          </div>
        </div>
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
import Option from "~~/components/Option.vue";
import PageLayout from "~~/components/reports/PageLayout.vue";
import lodash from "lodash";
definePageMeta({
  layout: "default",
});

const { api_url } = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const route = useRoute();
const activeQuizId = computed(() => route.params.id );
const currentTab = ref("report")

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

const changeTab = (data)=> {
  currentTab.value = data;
}
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
