<template>
  <PageLayout />
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

<script setup>
import Option from "~~/components/Option.vue";
import PageLayout from "~~/components/reports/PageLayout.vue";
definePageMeta({
  layout: "default",
});

const { api_url } = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const route = useRoute();
const activeQuizId = computed(() => route.params.id );

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
