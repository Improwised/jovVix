<template>
  <section class="mt-16">
    <div class="mb-6 flex items-center justify-between gap-4">
      <h2
        class="flex items-center gap-2 font-body text-[28px] font-normal leading-tight text-jv-ink sm:text-[32px]"
      >
        <span
          class="grid size-7 place-items-center rounded-full border-[3px] border-jv-coral text-jv-coral"
        >
          <Compass class="size-4" :stroke-width="2.6" />
        </span>
        <span>Explore Public Quizzes</span>
      </h2>
      <NuxtLink
        to="/admin/quiz/list-quiz"
        class="hidden items-center gap-1 text-[18px] font-semibold text-jv-ink underline decoration-1 underline-offset-2 sm:flex"
      >
        View All →
        <ArrowRight class="size-4" aria-hidden="true" />
      </NuxtLink>
    </div>

    <div class="grid gap-6 sm:grid-cols-2 xl:grid-cols-4">
      <PublicQuizCard
        v-for="quiz in quizzes"
        :key="quiz.title"
        :title="quiz.title"
        :description="quiz.description"
        :author="quiz.author"
        :variant="quiz.variant"
        :tilt-class="quiz.tiltClass"
      />
    </div>
  </section>
</template>

<script setup>
import PublicQuizCard from "@/components/landing/PublicQuizCard.vue";
import { ArrowRight, Compass } from "lucide-vue-next";

const quizzes = [
  {
    title: "Pop Hits 2024",
    description:
      "Guess the top charting songs and biggest artists of this year.",
    author: "Alex DJ",
    variant: "pop",
    tiltClass: "rotate-[-1deg]",
  },
  {
    title: "Global Explorer",
    description: "Test your knowledge of world capitals and hidden landmarks.",
    author: "GeoTeacher",
    variant: "global",
    tiltClass: "rotate-[1deg]",
  },
  {
    title: "Movie Buff Trivia",
    description:
      "From classic cinema to modern blockbusters. Can you get a perfect score?",
    author: "Cinephile",
    variant: "movie",
    tiltClass: "rotate-[-1deg]",
  },
  {
    title: "Animal Kingdom",
    description:
      "Identify these wild and wonderful creatures from around the world.",
    author: "WildLife",
    variant: "animal",
    tiltClass: "rotate-[0.6deg]",
  },
];
</script>
