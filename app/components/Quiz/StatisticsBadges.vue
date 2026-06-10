<script setup>
import {
  Target,
  Trophy,
  CheckCircle2,
  XCircle,
  CircleSlash,
} from "lucide-vue-next";

const props = defineProps({
  userStatistics: {
    default: () => ({}),
    type: Object,
    required: true,
  },
});

const stats = computed(() => [
  {
    label: "Accuracy",
    value: `${props.userStatistics?.accuracy ?? 0}%`,
    icon: Target,
    bg: "bg-jv-yellow",
  },
  {
    label: "Score",
    value: props.userStatistics?.totalScore ?? 0,
    icon: Trophy,
    bg: "bg-jv-mint",
  },
  {
    label: "Correct",
    value: props.userStatistics?.correctAnwers ?? 0,
    icon: CheckCircle2,
    bg: "bg-jv-mint-2",
  },
  {
    label: "Incorrect",
    value: props.userStatistics?.wrongAnwers ?? 0,
    icon: XCircle,
    bg: "bg-jv-salmon",
  },
  {
    label: "Skipped",
    value: props.userStatistics?.unAttemptedQuestions ?? 0,
    icon: CircleSlash,
    bg: "bg-jv-lavender",
  },
]);
</script>

<template>
  <section class="mx-auto w-full max-w-[1180px]" aria-label="Your statistics">
    <article
      class="relative rotate-[0.3deg] jv-border-rough bg-jv-white p-5 shadow-brutal sm:p-7"
    >
      <span
        class="absolute left-1/2 top-[-12px] z-10 h-4 w-16 -translate-x-1/2 rotate-[1deg] bg-jv-coral"
        aria-hidden="true"
      ></span>

      <header>
        <h2
          class="font-headings text-[22px] leading-none text-jv-ink sm:text-[28px]"
        >
          Your Stats
        </h2>
        <p
          class="mt-1.5 font-body text-[12px] font-bold text-jv-muted sm:text-[14px]"
        >
          How you performed
        </p>
      </header>

      <ul class="mt-5 grid grid-cols-2 gap-2.5 sm:mt-6 sm:gap-4 md:grid-cols-5">
        <li
          v-for="(stat, index) in stats"
          :key="stat.label"
          :class="[
            'flex items-center gap-2.5 border-[2px] border-jv-ink bg-jv-white px-2.5 py-3 shadow-brutal-sm sm:gap-3 sm:px-4 sm:py-4',
            index % 2 === 0 ? 'rotate-[-0.4deg]' : 'rotate-[0.4deg]',
          ]"
        >
          <span
            :class="[
              'grid size-9 shrink-0 place-items-center rounded-[6px] border-[2px] border-jv-ink sm:size-11',
              stat.bg,
            ]"
            aria-hidden="true"
          >
            <component
              :is="stat.icon"
              class="size-4 sm:size-5"
              :stroke-width="2.4"
            />
          </span>
          <div class="min-w-0 flex-1">
            <p
              class="truncate font-body text-[10px] font-bold uppercase tracking-[0.04em] text-jv-muted sm:text-[12px] sm:tracking-[0.08em]"
            >
              {{ stat.label }}
            </p>
            <p
              class="mt-0.5 truncate font-feature text-[18px] font-black tabular-nums text-jv-ink sm:text-[22px]"
            >
              {{ stat.value }}
            </p>
          </div>
        </li>
      </ul>
    </article>
  </section>
</template>
