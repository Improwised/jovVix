<script setup>
import { AlertTriangle } from "lucide-vue-next";

definePageMeta({
  layout: "empty",
});

useSeoMeta({
  title: "Quiz Results - jovVix",
  description:
    "Review detailed results and analytics for a completed quiz session on jovVix.",
  robots: "noindex, nofollow",
});

const headers = useRequestHeaders(["cookie"]);
const url = useRuntimeConfig().public;
const route = useRoute();

const played_quiz_id = computed(() => route.params.played_quiz_id);
const userStatistics = ref({});

const {
  data: quizList,
  pending: quizPending,
  error: quizError,
} = useFetch(`${url.apiUrl}/user_played_quizes/${played_quiz_id.value}`, {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
});

watch(
  [quizPending],
  () => {
    if (quizPending.value || quizError.value) {
      return;
    }
    userStatistics.value = questionsAnalysis(quizList.value?.data);
  },
  { immediate: true, deep: true }
);
</script>

<template>
  <ClientOnly>
    <div class="min-h-screen bg-jv-canvas px-4 py-6 sm:px-6 sm:py-10">
      <div class="mx-auto max-w-6xl">
        <header class="mb-6">
          <h1
            class="font-headings text-[38px] leading-none text-jv-ink sm:text-[52px]"
          >
            Quiz results
          </h1>
          <p
            class="mt-4 max-w-2xl text-[15px] font-semibold text-jv-muted sm:text-[16px]"
          >
            Review your results and question-by-question analysis.
          </p>
        </header>
        <NuxtLink
          to="/admin"
          class="mb-6 inline-flex w-fit items-center gap-2 rounded-full border-2 border-jv-ink bg-jv-white px-5 py-2 font-bold text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-1 active:shadow-none"
          >← Back to dashboard</NuxtLink
        >
        <div
          v-if="quizError"
          role="alert"
          class="jv-border-rough flex items-start gap-3 border-2 border-jv-coral bg-jv-coral/10 p-4 font-body text-jv-ink"
        >
          <AlertTriangle
            class="mt-0.5 size-5 shrink-0 text-jv-coral"
            :stroke-width="2.4"
          />
          <span>{{ quizError.data }}</span>
        </div>

        <div
          v-else-if="quizPending"
          class="jv-border-rough flex items-center justify-center gap-3 border-2 border-jv-ink bg-jv-white p-8 font-body text-jv-muted shadow-brutal-sm"
        >
          <span
            class="size-3 animate-pulse rounded-full bg-jv-coral"
            aria-hidden="true"
          ></span>
          Loading quiz results...
        </div>

        <div v-else class="flex flex-col gap-6">
          <QuizStatisticsBadges :user-statistics="userStatistics" />
          <QuizAnalysis :data="quizList.data" />
        </div>
      </div>
    </div>
  </ClientOnly>
</template>
