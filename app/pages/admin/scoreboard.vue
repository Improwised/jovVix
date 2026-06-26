<script setup>
const url = "/final_score/admin";

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
  <FinalScoreBoard :is-admin="true" :user-u-r-l="url" />
</template>
