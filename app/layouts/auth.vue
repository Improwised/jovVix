<template>
  <div class="min-h-screen bg-jv-canvas flex flex-col">
    <TopNav />

    <!-- Main content -->
    <main
      class="relative flex-1 flex items-center justify-center px-4 sm:px-6 py-6 sm:py-8 overflow-hidden"
    >
      <!-- Decorative background grid -->
      <div
        class="absolute inset-0 jv-grid pointer-events-none"
        aria-hidden="true"
      ></div>

      <!-- Decorative sparkles & squiggle -->
      <div
        class="absolute hidden sm:block left-[10%] top-[18%] rotate-12 text-jv-coral pointer-events-none"
        aria-hidden="true"
      >
        <Sparkle class="size-7 md:size-9 fill-current" :stroke-width="2" />
      </div>
      <div
        class="absolute hidden md:block right-[10%] top-[14%] -rotate-12 text-jv-yellow pointer-events-none"
        aria-hidden="true"
      >
        <svg
          width="120"
          height="48"
          viewBox="0 0 120 48"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M2 28 Q 22 4, 42 22 T 82 26 T 118 14"
            stroke="currentColor"
            stroke-width="3.5"
            stroke-linecap="round"
            fill="none"
          />
        </svg>
      </div>
      <div
        class="absolute hidden md:block right-[14%] bottom-[18%] -rotate-12 text-jv-mint pointer-events-none"
        aria-hidden="true"
      >
        <Sparkle class="size-6 fill-current" :stroke-width="2" />
      </div>

      <slot />
    </main>
  </div>
</template>

<script setup>
import { Sparkle } from "lucide-vue-next";
import TopNav from "@/components/common/TopNav.vue";
</script>
