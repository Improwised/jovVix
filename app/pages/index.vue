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

const route = useRoute();
const { baseUrl } = useRuntimeConfig().public;
const canonicalUrl = new URL(route.path, baseUrl).href;
const shareImageUrl = new URL("/og-image.png", baseUrl).href;

useHead({
  link: [{ rel: "canonical", href: canonicalUrl }],
});

useSeoMeta({
  title: "jovVix - Real-time Quiz Platform",
  description:
    "Create and host real-time interactive quizzes with jovVix. Engage your audience with live multiplayer games, instant scoring, and shareable results.",
  ogTitle: "jovVix - Real-time Quiz Platform",
  ogDescription:
    "Create and host real-time interactive quizzes with jovVix. Engage your audience with live multiplayer games, instant scoring, and shareable results.",
  ogType: "website",
  ogUrl: canonicalUrl,
  ogImage: shareImageUrl,
});
</script>
