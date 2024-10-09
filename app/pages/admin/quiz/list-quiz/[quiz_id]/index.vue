<template>
  <div v-if="quizPending">Pending...</div>
  <div v-else-if="quizError">{{ quizError }}</div>
  <div v-else class="container mt-3">
    <div
      v-for="(quiz, index) in quizData.data"
      :key="index"
      class="card mb-3 row"
    >
      <QuizQuestionAnalysis
        :question="quiz"
        :order="index + 1"
        :is-admin-analysis="true"
        :is-for-quiz="true"
        :quiz-id="quizId"
      />
      <QuizOptionsAnalysis
        :options="quiz?.options"
        :correct-answer="quiz?.correct_answer"
        :options-media="quiz?.options_media"
      />
    </div>
  </div>
</template>

<script setup>
const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const route = useRoute();
const quizId = computed(() => route.params.quiz_id || "");
const {
  data: quizData,
  pending: quizPending,
  error: quizError,
} = useFetch(`${url.api_url}/admin/quizzes/question/${quizId.value}`, {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
});
</script>
