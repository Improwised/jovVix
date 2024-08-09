<script setup>
const url = "/final_score/user";
const route = useRoute();
const username = computed(() => route.params.user);
const userPlayedQuiz = computed(() => route.query.user_played_quiz);
</script>

<template>
  <div class="w-100">
    <FinalScoreBoard
      :is-admin="false"
      :user-name="username"
      :user-u-r-l="url"
      :user-played-quiz="userPlayedQuiz"
    />
  </div>
</template>
