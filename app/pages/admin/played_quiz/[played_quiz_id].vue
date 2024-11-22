<script setup>
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
    <div>
      <div v-if="quizError" class="alert alert-danger" role="alert">
        {{ quizError.data }}
      </div>
      <div v-else-if="quizPending">Pending...</div>
      <div v-else>
        <QuizStatisticsBadges :user-statistics="userStatistics" />
        <QuizAnalysis :data="quizList.data" />
      </div>
    </div>
  </ClientOnly>
</template>
