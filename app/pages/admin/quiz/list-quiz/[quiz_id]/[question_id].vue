<template>
  <div v-if="questionPending">Pending...</div>
  <div v-else-if="questionError">{{ questionError }}</div>
  <div v-else class="container mt-3">
    <div class="card mb-3 row pb-2">
      <QuizEditQuestion
        :question="questionData?.data"
        :quiz-id="quizId"
        :question-id="questionId"
      />
    </div>
  </div>
</template>

<script setup>
const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const route = useRoute();
const quizId = computed(() => route.params.quiz_id || "");
const questionId = computed(() => route.params.question_id || "");
const {
  data: questionData,
  pending: questionPending,
  error: questionError,
} = useFetch(
  `${url.api_url}/quizzes/${quizId.value}/questions/${questionId.value}`,
  {
    method: "GET",
    headers: headers,
    mode: "cors",
    credentials: "include",
  }
);
</script>
