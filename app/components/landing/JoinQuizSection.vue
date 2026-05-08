<template>
  <section
    class="relative mt-10 flex min-h-[300px] items-center justify-center border-[3px] border-jv-ink bg-jv-yellow px-5 py-14 text-center shadow-brutal-lg jv-card lg:mt-9 lg:min-h-[360px]"
  >
    <div
      class="absolute left-1/2 top-8 grid size-16 -translate-x-1/2 rotate-[3deg] place-items-center rounded-[8px] border-[3px] border-jv-ink bg-jv-white shadow-brutal"
      aria-hidden="true"
    >
      <Gamepad2 class="size-8" :stroke-width="2.6" />
    </div>

    <div class="mx-auto mt-10 max-w-[620px]">
      <h2
        class="font-body text-[34px] font-black leading-tight text-jv-ink sm:text-[40px]"
      >
        Join a Live Quiz
      </h2>
      <p
        class="mx-auto mt-4 max-w-[620px] text-[18px] leading-[1.35] text-jv-muted sm:text-[20px]"
      >
        Enter your event invitation code below to jump straight into the action
        and compete with friends!
      </p>
      <LandingPillButton to="/join" class="mt-9" width="148px">
        Join Game
      </LandingPillButton>
    </div>
  </section>
</template>

<script setup>
import LandingPillButton from "@/components/landing/LandingPillButton.vue";
import { Gamepad2 } from "lucide-vue-next";
</script>
