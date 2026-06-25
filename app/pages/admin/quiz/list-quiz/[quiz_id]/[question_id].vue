<template>
  <div class="min-h-screen bg-jv-canvas px-4 py-6 sm:px-6 sm:py-10">
    <div class="mx-auto max-w-5xl">
      <div
        v-if="questionPending"
        class="jv-border-rough flex items-center justify-center gap-3 border-2 border-jv-ink bg-jv-white p-8 font-body text-jv-muted shadow-brutal-sm"
      >
        <span
          class="size-3 animate-pulse rounded-full bg-jv-coral"
          aria-hidden="true"
        ></span>
        Loading question...
      </div>

      <div
        v-else-if="questionError"
        role="alert"
        class="jv-border-rough border-2 border-jv-coral bg-jv-coral/10 p-4 font-body text-jv-ink"
      >
        {{ questionError }}
      </div>

      <div
        v-else
        class="jv-border-rough border-2 border-jv-ink bg-jv-white p-4 shadow-brutal sm:p-6"
      >
        <QuizEditQuestion
          :question="questionData?.data"
          :quiz-id="quizId"
          :question-id="questionId"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: "empty",
});

useSeoMeta({
  title: "Edit Question - jovVix",
  description: "Edit and configure an individual quiz question on jovVix.",
  robots: "noindex, nofollow",
});

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
  `${url.apiUrl}/quizzes/${quizId.value}/questions/${questionId.value}`,
  {
    method: "GET",
    headers: headers,
    mode: "cors",
    credentials: "include",
  }
);
</script>
