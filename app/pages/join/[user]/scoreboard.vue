<script setup>
const url = "/final_score/user";
</script>

<template>
  <div class="w-100">
    <FinalScoreBoard :is-admin="false" :user-u-r-l="url" />
  </div>
</template>
