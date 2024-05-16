<script setup>
const url = "/final_score/admin";
</script>

<template>
  <div class="w-100">
    <FinalScoreBoard :is-admin="true" :user-u-r-l="url" />
  </div>
</template>
