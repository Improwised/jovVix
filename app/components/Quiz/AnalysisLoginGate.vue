<script setup>
import { Lock, LogIn, UserPlus } from "lucide-vue-next";

const route = useRoute();
const returnTo = computed(() => encodeURIComponent(route.fullPath));
const loginUrl = computed(() => `/account/login?returnTo=${returnTo.value}`);
const signUpUrl = computed(
  () => `/account/register?returnTo=${returnTo.value}`
);
</script>

<template>
  <section
    class="mx-auto w-full max-w-[1180px]"
    aria-label="Per-question review"
  >
    <header class="mx-auto mb-4 max-w-[820px] sm:mb-6">
      <h2
        class="font-headings text-[22px] leading-none text-jv-ink sm:text-[28px]"
      >
        Question Review
      </h2>
      <p
        class="mt-1.5 font-body text-[12px] font-bold text-jv-muted sm:text-[14px]"
      >
        Log in to see what you got right and where you slipped
      </p>
    </header>

    <div
      class="relative mx-auto max-w-[820px] -rotate-[0.3deg] jv-border-rough bg-jv-white p-5 shadow-brutal sm:p-6"
    >
      <span
        class="absolute left-1/2 top-[-10px] z-10 h-3.5 w-14 -translate-x-1/2 rotate-[1deg] bg-jv-coral"
        aria-hidden="true"
      ></span>

      <div class="flex flex-col items-center gap-4 text-center sm:gap-5">
        <span
          class="flex size-12 items-center justify-center rounded-full border-[2px] border-jv-ink bg-jv-yellow shadow-brutal-sm sm:size-14"
          aria-hidden="true"
        >
          <Lock class="size-5 text-jv-ink sm:size-6" :stroke-width="2.4" />
        </span>

        <p
          class="max-w-[420px] font-body text-[13px] font-bold text-jv-muted sm:text-[15px]"
        >
          Your answer-by-answer breakdown is saved for members. Log in or create
          a free account to unlock it. Your results stay right here.
        </p>

        <div class="flex flex-col items-stretch gap-3 sm:flex-row">
          <NuxtLink
            :to="loginUrl"
            class="inline-flex h-11 items-center justify-center gap-2 rounded-full border-[2px] border-jv-ink bg-jv-coral px-6 font-body text-[14px] font-black text-white shadow-brutal-sm transition-transform hover:-rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-12 sm:text-[15px]"
            aria-label="Log in to view your question review"
          >
            <LogIn class="size-4" :stroke-width="2.4" />
            Log In
          </NuxtLink>
          <NuxtLink
            :to="signUpUrl"
            class="inline-flex h-11 items-center justify-center gap-2 rounded-full border-[2px] border-jv-ink bg-jv-white px-6 font-body text-[14px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-12 sm:text-[15px]"
            aria-label="Create a free account"
          >
            <UserPlus class="size-4" :stroke-width="2.4" />
            Sign Up
          </NuxtLink>
        </div>
      </div>
    </div>
  </section>
</template>
