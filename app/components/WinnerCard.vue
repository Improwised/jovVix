<script setup>
import { Crown, Medal, Award } from "lucide-vue-next";
import { getAvatarUrlByName } from "~~/composables/avatar";

const props = defineProps({
  winner: {
    type: Object,
    required: true,
    default: () => ({}),
  },
});

const avatar = computed(() => getAvatarUrlByName(props.winner?.img_key));

const rankConfig = computed(() => {
  const rank = Number(props.winner?.rank);
  if (rank === 1) {
    return {
      bg: "bg-jv-yellow",
      icon: Crown,
      label: "1st",
      rotate: "rotate-[-1deg]",
      scale: "sm:scale-110",
    };
  }
  if (rank === 2) {
    return {
      bg: "bg-jv-salmon",
      icon: Medal,
      label: "2nd",
      rotate: "rotate-[1deg]",
      scale: "",
    };
  }
  return {
    bg: "bg-jv-mint",
    icon: Award,
    label: "3rd",
    rotate: "rotate-[-0.5deg]",
    scale: "",
  };
});

const getOrdinal = (rank) => {
  const suffixes = ["th", "st", "nd", "rd"];
  const v = rank % 100;
  return rank + (suffixes[(v - 20) % 10] || suffixes[v] || suffixes[0]);
};
</script>

<template>
  <article
    :class="[
      'relative flex w-full max-w-[280px] flex-col items-center gap-3 jv-border-rough p-5 shadow-brutal-lg sm:gap-4 sm:p-6',
      rankConfig.bg,
      rankConfig.rotate,
      rankConfig.scale,
    ]"
    :aria-label="`${getOrdinal(props.winner.rank)} place: ${
      props.winner.firstname
    } with score ${props.winner.score}`"
    role="article"
  >
    <span
      class="absolute left-1/2 top-[-12px] z-10 h-4 w-16 -translate-x-1/2 rotate-[2deg] bg-jv-coral"
      aria-hidden="true"
    ></span>

    <div
      class="inline-flex items-center gap-2 rounded-full border-[2px] border-jv-ink bg-jv-white px-3 py-1 font-feature text-[14px] font-black text-jv-ink shadow-brutal-sm sm:text-[15px]"
    >
      <component
        :is="rankConfig.icon"
        class="size-4 text-jv-ink"
        :stroke-width="2.4"
      />
      {{ rankConfig.label }} Place
    </div>

    <div
      class="grid size-24 place-items-center overflow-hidden rounded-full border-[3px] border-jv-ink bg-jv-white shadow-brutal-sm sm:size-28"
    >
      <img
        :src="avatar"
        :alt="`Avatar for ${props.winner.firstname}`"
        class="size-full object-cover"
      />
    </div>

    <div class="flex w-full flex-col items-center gap-1 text-center">
      <h2
        :id="`winner-name-${props.winner.rank}`"
        class="break-words font-headings text-[22px] uppercase leading-tight text-jv-ink sm:text-[26px]"
      >
        {{ props.winner.firstname }}
      </h2>
      <p
        v-if="props.winner.username"
        class="font-body text-[12px] font-bold text-jv-muted sm:text-[13px]"
      >
        @{{ props.winner.username }}
      </p>
    </div>

    <div
      class="mt-1 inline-flex items-baseline gap-2 rounded-full border-[2px] border-jv-ink bg-jv-white px-4 py-1.5 shadow-brutal-sm"
    >
      <span class="font-body text-[11px] font-bold text-jv-muted">SCORE</span>
      <span
        class="font-feature text-[22px] font-black tabular-nums text-jv-ink sm:text-[26px]"
      >
        {{ props.winner.score }}
      </span>
    </div>
  </article>
</template>
