<template>
  <PageLayout
    v-model:question-type-filter="filter"
    v-model:user-filter="userRankFilter"
    :current-tab="currentTab"
    :total-question="filterQuestionAnalysis.length"
    :all-question="allQuestionsLen"
    @change-tab="changeTab"
  />
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
        v-for="(quiz, index) in filterQuestionAnalysis"
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
          v-for="oData in rankData"
          :key="oData"
          :data="userJson[oData]"
          :user-name="oData"
          :survey-questions="surveyQuestions"
          class="user-analytics-item"
        ></QuizUserAnalyticsSpace>
      </div>
    </div>
    <div
      v-if="!userRankFilter.showTop10"
      class="d-flex align-items-center justify-content-center"
    >
      <Pagination
        :page="Math.floor(userRankFilter.offset / userRankFilter.limit) + 1"
        :num-of-records="Math.ceil(totalUsers / userRankFilter.limit)"
      />
    </div>
  </template>
</template>

<script setup>
import PageLayout from "~~/components/reports/PageLayout.vue";
import lodash from "lodash";
import { ref, computed } from "vue";

definePageMeta({
  layout: "default",
});

const { apiUrl } = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const route = useRoute();
const activeQuizId = computed(() => route.params.id);
const currentTab = ref("report");

const {
  data: quizAnalysis,
  error: quizAnalysisError,
  pending: quizAnalysisPending,
} = useFetch(`${apiUrl}/admin/reports/${activeQuizId.value}/analysis`, {
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
const quizUserAnalysispending = ref(false);
const fetchError = ref("");
const filter = ref("all");

const userRankFilter = ref({
  isAsc: false,
  limit: 10,
  offset: 0,
  showTop10: false,
});

const getAnalysisJson = async () => {
  try {
    quizUserAnalysispending.value = true;

    const response = await fetch(
      `${apiUrl}/analytics_board/admin?active_quiz_id=${activeQuizId.value}`,
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
      questionJson.value = lodash.groupBy(analysisJson.value, "question");
    } else {
      console.error(result);
    }

    quizUserAnalysispending.value = false;
  } catch (error) {
    quizUserAnalysispending.value = false;
    fetchError.value = error;
    console.error("Failed to fetch analysis data", error);
  }
};

const totalUsers = ref(0);

const getRankData = async () => {
  try {
    let url = `${apiUrl}/final_score/admin?active_quiz_id=${activeQuizId.value}`;

    if (userRankFilter.value.showTop10) {
      url += `&user_limit=10&starting_at=0&order_by=desc`;
    } else {
      const orderBy = userRankFilter.value.isAsc ? "asc" : "desc";
      const userLimit = userRankFilter.value.limit;
      const offset = userRankFilter.value.offset;

      url += `&user_limit=${userLimit}&starting_at=${offset}&order_by=${orderBy}`;
    }

    const ranksResponse = await fetch(url, {
      method: "GET",
      headers: headers.value,
      mode: "cors",
      credentials: "include",
    });

    const ranksResult = await ranksResponse.json();

    if (ranksResponse.ok) {
      rankData.value = [];
      totalUsers.value = ranksResult.count; // 👈 store total count

      ranksResult.data.forEach((data) => {
        rankData.value.push(data.username);

        let key = data.username;
        if (userJson.value.hasOwnProperty(key)) {
          userJson.value[key].push({
            rank: data.rank,
            total_score: data.score,
            response_time: data.response_time,
            avatar: data.img_key,
          });
        }
      });
    } else {
      console.error(ranksResult);
    }
  } catch (error) {
    console.error("Failed to fetch rank data", error);
  }
};

onMounted(() => {
  getAnalysisJson();
  getRankData();
});

watch(
  userRankFilter,
  () => {
    getRankData();
  },
  { deep: true }
);

const changeTab = (data) => {
  currentTab.value = data;
};

const filterQuestionAnalysis = computed(() => {
  if (!quizAnalysis.value?.data) return [];

  if (filter.value === "all") {
    return quizAnalysis.value.data;
  }
  return quizAnalysis.value.data.filter((q) => q.type == filter.value);
});

const allQuestionsLen = computed(() => quizAnalysis.value?.data.length);
const page = computed(() => Number(route.query.page) || 1);

watch(page, (newPage) => {
  userRankFilter.value.offset = (newPage - 1) * userRankFilter.value.limit;
});
</script>
