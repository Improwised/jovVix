<template>
  <div
    class="d-flex flex-wrap align-items-center justify-content-around container bg-white rounded p-5"
  >
    <span class="badge rounded-pill bg-light-primary text-dark m-2 px-2 fs-5">
      Accuracy: {{ props.userStatistics?.accuracy }}%
    </span>
    <span class="badge rounded-pill bg-light-primary text-dark m-2 px-2 fs-5">
      Total Score: {{ props.userStatistics?.totalScore }}
    </span>
    <span class="badge rounded-pill bg-light-primary text-dark m-2 px-2 fs-5">
      Total Correct: {{ props.userStatistics?.correctAnwers }}</span
    >
    <span class="badge rounded-pill bg-light-secondary text-dark m-2 px-2 fs-5">
      Total Incorrect: {{ props.userStatistics?.wrongAnwers }}
    </span>
    <span class="badge rounded-pill bg-light-secondary text-dark m-2 px-2 fs-5">
      Total Un-attmpted: {{ props.userStatistics?.unAttemptedQuestions }}
    </span>
  </div>
</template>

<script setup>
const props = defineProps({
  userStatistics: {
    default: () => {
      return {};
    },
    type: Object,
    required: true,
  },
});
</script>
