<script setup>
const url = "/final_score/user";
const route = useRoute();
const username = computed(() => route.params.user);
const userPlayedQuiz = computed(() => route.query.user_played_quiz);

definePageMeta({
  layout: "empty",
  hideSidebar: true,
});

useSeoMeta({
  title: "Quiz Results - jovVix",
  description:
    "View your quiz results, ranking, and per-question analysis after the game.",
});
</script>

<template>
  <FinalScoreBoard
    :is-admin="false"
    :user-name="username"
    :user-u-r-l="url"
    :user-played-quiz="userPlayedQuiz"
  />
</template>
