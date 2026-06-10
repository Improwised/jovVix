<template>
  <section
    class="grid items-start gap-10 sm:gap-12 lg:gap-12 xl:gap-20 lg:grid-cols-[minmax(0,1fr)_300px] xl:grid-cols-[minmax(0,1fr)_400px]"
  >
    <div
      class="flex flex-col gap-3 items-start order-2 lg:order-1 w-full min-w-0"
    >
      <div
        class="rounded-[8px] border-[3px] border-jv-ink bg-jv-yellow/30 px-2.5 py-1 sm:px-3 sm:py-1.5 text-[12px] sm:text-[14px] md:small-text font-semibold uppercase shadow-brutal-sm"
      >
        EXPLORE, LEARN, QUIZ!
      </div>

      <h1 class="main-heading text-jv-ink">
        <span class="block">Open-Source Quizzing</span>
        <span class="block"
          >Built for
          <span
            class="inline-block jv-highlight mt-2 sm:mt-3 px-3 sm:px-5 md:px-6 py-1 sm:py-2 rounded-[4px]"
            >Live Quiz Experiences</span
          >
        </span>
      </h1>

      <!-- Mobile/sm carousel (under md) -->
      <Carousel
        :opts="{ align: 'start', loop: false }"
        class="group mt-5 sm:mt-7 w-full md:hidden"
      >
        <CarouselContent class="-ml-4 py-2">
          <CarouselItem
            v-for="card in actionCards"
            :key="card.title"
            class="pl-4 basis-[80%] min-[420px]:basis-1/2 sm:basis-1/3"
          >
            <HeroActionCard
              :title="card.title"
              :description="card.description"
              :tilt-class="card.tiltClass"
            />
          </CarouselItem>
        </CarouselContent>
        <CarouselPrevious class="jv-carousel-nav !left-2" />
        <CarouselNext class="jv-carousel-nav !right-2" />
      </Carousel>

      <!-- md+ static grid -->
      <div
        class="hidden md:grid mt-5 sm:mt-7 gap-5 sm:gap-6 md:gap-8 lg:gap-10 grid-cols-3 w-full"
      >
        <HeroActionCard
          v-for="card in actionCards"
          :key="card.title"
          :title="card.title"
          :description="card.description"
          :tilt-class="card.tiltClass"
        />
      </div>
    </div>

    <div class="mx-auto max-w-[500px] pt-1 lg:pt-4 order-1 lg:order-2">
      <div class="relative rotate-[3deg]">
        <span
          class="absolute left-1/2 top-[-10px] z-20 h-5 w-14 -translate-x-1/2 rounded-[7px] border-[2px] border-jv-ink bg-jv-slate shadow-brutal-sm"
          aria-hidden="true"
        ></span>
        <div
          class="relative overflow-hidden bg-jv-white p-2 sm:p-3 shadow-brutal-lg jv-border-rough"
        >
          <img
            src="/images/landing/homepage-hero-section.png"
            alt="Live quiz preview"
            class="w-full h-auto"
          />
        </div>
        <div
          class="absolute bottom-[-15px] left-[-14px] grid size-10 sm:size-12 rotate-[-10deg] place-items-center rounded-[8px] border-[3px] border-jv-ink bg-white text-jv-coral shadow-brutal-sm"
        >
          <Users class="size-5 sm:size-6" :stroke-width="2.5" />
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import HeroActionCard from "@/components/landing/HeroActionCard.vue";
import { Users } from "lucide-vue-next";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";

const actionCards = [
  {
    title: "Demo Instance",
    description: "Up to 100 concurrent users",
    tiltClass: "rotate-[-0.7deg]",
  },
  {
    title: "Deploy Your Own",
    description: "Unlimited users & control.",
    tiltClass: "rotate-[0.6deg]",
  },
  {
    title: "Get in touch",
    description: "contact@jovvix.com",
    tiltClass: "rotate-[-0.4deg]",
  },
];
</script>
