<template>
  <section class="mt-2 sm:mt-4">
    <div class="mb-5 sm:mb-6 flex items-center justify-between gap-3 sm:gap-4">
      <h3
        class="flex items-center gap-2 font-headings text-[22px] min-[420px]:text-[26px] sm:text-[30px] md:text-[34px] leading-[1.2] text-jv-ink"
      >
        <span
          class="grid size-6 sm:size-7 shrink-0 place-items-center rounded-full border-[3px] border-jv-coral text-jv-coral"
        >
          <Compass class="size-3.5 sm:size-4" :stroke-width="2.6" />
        </span>
        <span>Explore Public Quizzes</span>
      </h3>
      <NuxtLink
        to="/admin/quiz/list-quiz"
        class="flex items-center gap-1 text-[14px] min-[420px]:text-[16px] sm:text-[18px] font-semibold text-jv-ink underline decoration-1 underline-offset-2 shrink-0"
      >
        View All
        <ArrowRight class="size-3.5 sm:size-4" aria-hidden="true" />
      </NuxtLink>
    </div>

    <!-- Carousel: under lg -->
    <Carousel
      :opts="{ align: 'start', loop: false }"
      class="group w-full lg:hidden"
    >
      <CarouselContent class="-ml-4 py-2">
        <CarouselItem
          v-for="quiz in quizzes"
          :key="quiz.title"
          class="pl-4 basis-[85%] min-[420px]:basis-[60%] sm:basis-1/2 md:basis-[40%]"
        >
          <PublicQuizCard
            :title="quiz.title"
            :description="quiz.description"
            :author="quiz.author"
            :tilt-class="quiz.tiltClass"
            :image="quiz.image"
          />
        </CarouselItem>
      </CarouselContent>
      <CarouselPrevious class="jv-carousel-nav !left-3" />
      <CarouselNext class="jv-carousel-nav !right-3" />
    </Carousel>

    <!-- Grid: lg+ -->
    <div class="hidden lg:grid gap-5 sm:gap-6 lg:grid-cols-3 xl:grid-cols-4">
      <PublicQuizCard
        v-for="quiz in quizzes"
        :key="quiz.title"
        :title="quiz.title"
        :description="quiz.description"
        :author="quiz.author"
        :tilt-class="quiz.tiltClass"
        :image="quiz.image"
      />
    </div>
  </section>
</template>

<script setup>
import PublicQuizCard from "@/components/landing/PublicQuizCard.vue";
import { ArrowRight, Compass } from "lucide-vue-next";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";

const quizzes = [
  {
    title: "Pop Hits 2024",
    description:
      "Guess the top charting songs and biggest artists of this year.",
    author: "Alex DJ",
    tiltClass: "rotate-[-1deg]",
    image: "/images/landing/homepage-public-quiz-1.png",
  },
  {
    title: "Global Explorer",
    description: "Test your knowledge of world capitals and hidden landmarks.",
    author: "GeoTeacher",
    tiltClass: "rotate-[1deg]",
    image: "/images/landing/homepage-public-quiz-2.png",
  },
  {
    title: "Movie Buff Trivia",
    description:
      "From classic cinema to modern blockbusters. Can you get a perfect score?",
    author: "Cinephile",
    tiltClass: "rotate-[-1deg]",
    image: "/images/landing/homepage-public-quiz-3.png",
  },
  {
    title: "Animal Kingdom",
    description:
      "Identify these wild and wonderful creatures from around the world.",
    author: "WildLife",
    tiltClass: "rotate-[0.6deg]",
    image: "/images/landing/homepage-public-quiz-4.png",
  },
];
</script>
