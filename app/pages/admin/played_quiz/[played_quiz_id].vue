<script setup>
const headers = useRequestHeaders(["cookie"]);
const url = useRuntimeConfig().public;
const route = useRoute();

const userAccuracy = ref(0)
const userTotalScore = ref(0)
const played_quiz_id = computed(() => route.params.played_quiz_id);

const {
  data: quizList,
  pending: quizPending,
  error: quizError,
} = useFetch(`${url.api_url}/user_played_quizes/${played_quiz_id.value}`, {
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
    const analysis = questionsAnalysis(quizList.value?.data)
    userAccuracy.value = analysis.accuracy;
    userTotalScore.value = analysis.totalScore;
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
      <div v-else-if="quizPending">
        Pending...
      </div>
      <div v-else>
        <h3 class="text-center">Accuracy: {{ userAccuracy }}%</h3>
        <h3 class="text-center">Total Score: {{ userTotalScore }}</h3>
        <QuizQuestionAnalysis :data="quizList.data" />
      </div>
    </div>
  </ClientOnly>
</template>
