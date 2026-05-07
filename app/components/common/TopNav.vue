<template>
  <header
    class="sticky top-0 z-40 border-b-[3px] border-jv-ink bg-jv-canvas/95 backdrop-blur"
  >
    <div
      class="flex items-center justify-between gap-4 px-4 sm:px-6 md:px-8 py-3 md:py-3.5"
    >
      <NuxtLink
        to="/"
        class="flex items-center gap-2 text-[20px] sm:text-[22px] md:text-[24px] font-black tracking-[-0.5px] text-jv-ink no-underline"
      >
        <span
          class="grid size-8 sm:size-9 rotate-[-6deg] place-items-center rounded-[8px] border-2 border-jv-ink bg-jv-coral text-white shadow-brutal-sm"
        >
          <Zap class="size-4 sm:size-[18px]" :stroke-width="2.4" />
        </span>
        <span>JovVix</span>
      </NuxtLink>

      <nav class="hidden md:flex items-center gap-6 lg:gap-7">
        <NuxtLink
          to="/#features"
          class="text-jv-ink text-sm lg:text-[15px] font-body hover:text-jv-coral transition-colors no-underline"
        >
          Features
        </NuxtLink>
        <NuxtLink
          to="/#docs"
          class="text-jv-ink text-sm lg:text-[15px] font-body hover:text-jv-coral transition-colors no-underline"
        >
          Docs
        </NuxtLink>
        <NuxtLink
          to="/#community"
          class="text-jv-ink text-sm lg:text-[15px] font-body hover:text-jv-coral transition-colors no-underline"
        >
          Community
        </NuxtLink>
      </nav>

      <NuxtLink
        to="/join"
        class="inline-flex items-center justify-center rotate-[-1deg] rounded-[8px] border-2 border-jv-ink bg-jv-yellow-2 text-jv-ink shadow-brutal-sm px-4 sm:px-5 py-2 sm:py-2.5 text-xs sm:text-[13px] font-body no-underline transition-transform hover:rotate-[2deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none"
      >
        Play a Quiz
      </NuxtLink>
    </div>
  </header>
</template>

<script setup>
import { Zap } from "lucide-vue-next";
</script>
