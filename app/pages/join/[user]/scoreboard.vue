<script setup>
const url = "/final_score/user";
const route = useRoute();
const username = computed(() => route.params.user);
</script>

<template>
  <div class="w-100">
    <FinalScoreBoard
      :is-admin="false"
      :user-name="username"
      :user-u-r-l="url"
    />
  </div>
</template>
