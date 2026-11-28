<template>
  <div class="min-h-screen bg-jv-canvas px-4 py-6 sm:px-6 sm:py-10">
    <div class="mx-auto max-w-5xl">
      <header class="mb-6">
        <h1
          class="font-headings text-[38px] leading-none text-jv-ink sm:text-[52px]"
        >
          Edit question
        </h1>
        <p
          class="mt-4 max-w-2xl text-[15px] font-semibold text-jv-muted sm:text-[16px]"
        >
          Update the question and its answer options.
        </p>
      </header>
      <NuxtLink
        :to="`/admin/quiz/list-quiz/${quizId}`"
        class="mb-6 inline-flex w-fit items-center gap-2 rounded-full border-2 border-jv-ink bg-jv-white px-5 py-2 font-bold text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-1 active:shadow-none"
        >← Back to quiz</NuxtLink
      >
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
useSessionExpiredRedirect(questionError);
</script>
