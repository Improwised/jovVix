<template>
  <PageLayout />
  <div v-if="quizAnalysisPending">loading...</div>
  <div v-else-if="quizAnalysisError">{{ quizAnalysisError }}</div>
  <div v-else class="card m-3 row" v-for="(quiz, index) in quizAnalysis.data">
    <div class="row m-2">
      <div class="col-lg-8 col-sm-12">
        <div>
          <strong class="text-primary">Question: </strong>
          <h4 class="font-bold">{{ quiz.question }}</h4>
        </div>
      </div>
      <div
        class="col-lg-4 col-sm-12 d-flex align-items-center justify-content-around"
      >
        <span class="badge bg-primary"
          >AVG. Response Time:
          {{ (quiz.avg_response_time / 1000).toFixed(2) }}</span
        >
        <span class="badge bg-secondary mx-2"
          >Time Duration: {{ quiz.duration }}</span
        >
        <span v-if="quiz.type === 1" class="badge bg-light-info mx-2 text-dark"
          >Multiple Choice Question</span
        >
        <span v-else class="badge bg-light-info mx-2 text-dark">Survey</span>
      </div>
    </div>
    <div class="border-bottom pb-4 mb-4"></div>
    <div class="row d-flex align-items-stretch">
      <div class="col-sm-12 col-lg-8 part-left">
        <div class="row">
          <div
            v-for="(option, order) in quiz.options"
            class="col-lg-6 col-md-12"
          >
            <div
              v-if="quiz.correct_answer.includes(Number(order))"
              class="bg-light-success option-box d-flex align-items-center m-2 border-rounded justify-content-between position-relative"
            >
              <Option
                class=""
                :order="order"
                :option="option"
                :selected="quiz.selected_answers[order]?.length || 0"
                icon="fa-solid fa-check"
                :isCorrect="true"
              />
              <span
                class="position-absolute top-0 start-50 translate-middle badge rounded-pill bg-success"
              >
                correct
              </span>
            </div>
            <div
              v-else
              class="bg-light-primary option-box d-flex align-items-center m-2 border-rounded justify-content-between"
            >
              <Option
                class=""
                :order="order"
                :option="option"
                :selected="quiz.selected_answers[order]?.length || 0"
              />
            </div>
          </div>
        </div>
      </div>

      <div class="col-sm-12 col-lg-4 part-right">
        <div class="fw-bold">Accuracy:</div>
        <div class="progress">
          <div
            class="progress-bar bg-success"
            role="progressbar"
            :style="{ width: `${quiz.correctPercentage}%` }"
            :aria-valuenow="quiz.correctPercentage"
            aria-valuemin="0"
            aria-valuemax="100"
          >
            {{ quiz.correctPercentage }}%
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
const quizId = route.params.id;

const {
  data: quizAnalysis,
  error: quizAnalysisError,
  pending: quizAnalysisPending,
} = useFetch(`${api_url}/admin/reports/${quizId}/analysis`, {
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
  height: 70px;
  /* width: 70%; */
  border-radius: 30px;
}
.part-left {
  border-right: 1px solid #ccc; /* Vertical divider */
}

.part-right {
  border-top: 1px solid #ccc; /* Horizontal divider */
}

/* Adjust the divider for large screens */
@media (min-width: 992px) {
  .part-left {
    border-right: 1px solid #ccc; /* Keep vertical divider */
    border-bottom: none; /* Remove horizontal divider */
  }

  .part-right {
    border-top: none; /* Remove horizontal divider */
  }
}

/* Adjust the divider for smaller screens */
@media (max-width: 991.98px) {
  .part-left {
    border-right: none; /* Remove vertical divider */
    border-bottom: 1px solid #ccc; /* Add horizontal divider */
  }

  .part-right {
    border-top: none; /* Remove horizontal divider */
  }
}
</style>
