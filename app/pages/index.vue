<template>
  <main class="flex flex-col gap-8 sm:gap-10 px-4 sm:px-6 md:px-8 py-5 md:py-6">
    <LandingTopActions />
    <HeroSection />
    <JoinLiveQuiz />
    <PublicQuizSection />
  </main>
</template>

<script setup>
import HeroSection from "@/components/landing/HeroSection.vue";
import JoinLiveQuiz from "@/components/Quiz/JoinLiveQuiz.vue";
import LandingTopActions from "@/components/landing/LandingTopActions.vue";
import PublicQuizSection from "@/components/landing/PublicQuizSection.vue";

definePageMeta({
  layout: "empty",
});
</script>
