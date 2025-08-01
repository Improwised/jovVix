<script setup>
const url = "/final_score/admin";

useHead({
  title: "Quiz Scoreboard - Admin Dashboard",
  meta: [
    {
      name: "description",
      content:
        "View quiz results and winners on the admin scoreboard. Switch between table view and winners podium display.",
    },
  ],
});
</script>

<template>
  <div class="w-100">
    <FinalScoreBoard :is-admin="true" :user-u-r-l="url" />
  </div>
</template>
