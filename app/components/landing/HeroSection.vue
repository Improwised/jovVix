<template>
  <section
    class="grid items-start gap-9 lg:grid-cols-[minmax(0,1fr)_390px] xl:grid-cols-[minmax(0,1fr)_430px]"
  >
    <div>
      <div
        class="inline-flex rotate-[-1deg] items-center rounded-[8px] border-[3px] border-jv-ink bg-[#fff8d8] px-3 py-1.5 text-[17px] font-bold uppercase tracking-[0.08em] shadow-brutal-sm"
      >
        EXPLORE, LEARN, QUIZ!
      </div>

      <h1
        class="mt-4 max-w-[760px] font-body text-[48px] font-normal leading-[0.98] tracking-normal text-jv-ink sm:text-[64px] lg:text-[58px] xl:text-[68px]"
      >
        <span class="block">Open-Source Quizzing</span>
        <span class="block">Built for</span>
        <span class="jv-highlight mt-2 inline-block px-3 pb-1 leading-[1.05]"
          >Live Engagement</span
        >
      </h1>

      <div class="mt-7 grid gap-5 sm:grid-cols-3 lg:max-w-[770px]">
        <HeroActionCard
          title="Demo Instance"
          description="Up to 100 concurrent"
          tilt-class="rotate-[-0.7deg]"
        />
        <HeroActionCard
          title="Deploy Your Own"
          description="Unlimited users & control."
          tilt-class="rotate-[0.6deg]"
        />
        <HeroActionCard
          title="Get in touch"
          description="contact@jovvix.com"
          tilt-class="rotate-[-0.4deg]"
        />
      </div>
    </div>

    <div class="mx-auto w-full max-w-[430px] pt-1 lg:pt-4">
      <div class="relative rotate-[3deg]">
        <span
          class="absolute left-1/2 top-[-10px] z-20 h-5 w-14 -translate-x-1/2 rounded-[7px] border-[2px] border-jv-ink bg-jv-slate shadow-brutal-sm"
          aria-hidden="true"
        ></span>
        <div
          class="relative min-h-[245px] overflow-hidden bg-jv-white p-3 shadow-brutal-lg jv-border-rough sm:min-h-[285px]"
        >
          <!-- TODO: Replace this CSS sketch with the final classroom illustration asset when available. -->
          <div
            class="relative h-full min-h-[220px] border border-jv-ink/20 bg-[#fbfdff] p-4"
          >
            <div
              class="absolute inset-x-8 top-4 h-16 rounded-full bg-sky-100/70 blur-sm"
            ></div>
            <div
              class="relative mx-auto w-[160px] rotate-[1deg] rounded-[2px] border-2 border-jv-ink bg-white p-2 text-center shadow-brutal-sm"
            >
              <p class="text-[9px] font-black uppercase leading-tight">
                What is the capital of France?
              </p>
              <div class="mt-1 grid grid-cols-2 gap-1 text-[8px]">
                <span class="rounded border border-jv-ink/40 bg-jv-yellow-soft"
                  >Paris</span
                >
                <span class="rounded border border-jv-ink/40">Rome</span>
                <span class="rounded border border-jv-ink/40">Berlin</span>
                <span class="rounded border border-jv-ink/40">Delhi</span>
              </div>
            </div>

            <div
              class="absolute bottom-7 left-1/2 h-20 w-16 -translate-x-1/2 rounded-t-full border-[3px] border-jv-ink bg-white"
            >
              <div
                class="mx-auto mt-5 size-9 rounded-full border-2 border-jv-ink bg-[#ffe0c2]"
              ></div>
              <div
                class="mx-auto mt-2 h-5 w-12 rounded-t-full border-2 border-jv-ink bg-jv-mint"
              ></div>
            </div>

            <div class="absolute bottom-5 left-5 grid grid-cols-4 gap-2">
              <span
                v-for="item in 12"
                :key="`left-${item}`"
                class="size-6 rounded-full border border-jv-ink/50 bg-white"
              ></span>
            </div>
            <div class="absolute bottom-5 right-5 grid grid-cols-4 gap-2">
              <span
                v-for="item in 12"
                :key="`right-${item}`"
                class="size-6 rounded-full border border-jv-ink/50 bg-white"
              ></span>
            </div>
            <span class="absolute left-6 top-12 text-2xl">?</span>
            <span class="absolute right-8 top-16 text-2xl">!</span>
          </div>
        </div>
        <div
          class="absolute bottom-[-15px] left-[-14px] grid size-12 rotate-[-10deg] place-items-center rounded-[8px] border-[3px] border-jv-ink bg-white text-jv-coral shadow-brutal-sm"
        >
          <Users class="size-6" :stroke-width="2.5" />
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import HeroActionCard from "@/components/landing/HeroActionCard.vue";
import { Users } from "lucide-vue-next";
</script>
