<template>
  <div class="mb-8 hidden h-[50px] items-start justify-end gap-6 lg:flex">
    <LandingPillButton to="/join" width="149px"> Enter Code </LandingPillButton>
    <LandingPillButton to="/admin/quiz/create-quiz" width="161px">
      Create Quiz
    </LandingPillButton>
  </div>
</template>

<script setup>
import LandingPillButton from "@/components/landing/LandingPillButton.vue";
</script>
