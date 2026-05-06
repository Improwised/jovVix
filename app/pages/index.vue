<template>
  <div
    class="min-h-screen overflow-x-hidden bg-jv-canvas text-jv-ink font-body"
  >
    <LandingSidebar />

    <main class="landing-main min-h-screen">
      <div
        class="landing-content mx-auto w-full max-w-[1240px] px-4 pb-10 pt-6 sm:px-6 lg:px-5 lg:pb-9"
      >
        <LandingTopActions />
        <HeroSection />
        <JoinQuizSection />
        <PublicQuizSection />
      </div>
    </main>
  </div>
</template>

<script setup>
import HeroSection from "@/components/landing/HeroSection.vue";
import JoinQuizSection from "@/components/landing/JoinQuizSection.vue";
import LandingSidebar from "@/components/landing/LandingSidebar.vue";
import LandingTopActions from "@/components/landing/LandingTopActions.vue";
import PublicQuizSection from "@/components/landing/PublicQuizSection.vue";

definePageMeta({
  layout: "empty",
});
</script>

<style scoped>
@media (min-width: 1024px) {
  .landing-main {
    margin-left: 14rem;
    position: relative;
  }

  .landing-content {
    padding-top: 24px !important;
    position: relative;
    z-index: 1;
  }

  .landing-main::before {
    background-image: radial-gradient(var(--jv-ink) 1px, transparent 1px);
    background-size: 20px 20px;
    content: "";
    inset: 0;
    opacity: 0.04;
    pointer-events: none;
    position: absolute;
  }
}
</style>
