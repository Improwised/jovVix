<script setup>
const url = "/final_score/admin";
const route = useRoute();

definePageMeta({
  layout: "empty",
  hideSidebar: true,
});

useSeoMeta({
  title: "Quiz Scoreboard - jovVix",
  description:
    "View quiz results and winners on the admin scoreboard. Switch between table view and winners podium display.",
  robots: "noindex, nofollow",
});
</script>

<template>
  <FinalScoreBoard
    :key="String(route.query.aqi || '')"
    :is-admin="true"
    :user-u-r-l="url"
  />
</template>
