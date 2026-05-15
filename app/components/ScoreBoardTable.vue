<script setup>
import { Crown, Medal, Award, Star } from "lucide-vue-next";
import { getAvatarUrlByName } from "~~/composables/avatar";

const props = defineProps({
  scoreboardData: {
    type: Array,
    required: true,
    default: () => [],
  },
  isAdmin: {
    type: Boolean,
    required: false,
    default: false,
  },
  userName: {
    type: String,
    required: true,
    default: "",
  },
});

const avatarAccents = [
  "bg-jv-yellow",
  "bg-jv-salmon",
  "bg-jv-mint",
  "bg-jv-ivory",
  "bg-jv-lavender",
];
const avatarBgFor = (index) => avatarAccents[index % avatarAccents.length];

const rankBadgeFor = (rank) => {
  const r = Number(rank);
  if (r === 1) return { icon: Crown, bg: "bg-jv-yellow", aria: "First place" };
  if (r === 2) return { icon: Medal, bg: "bg-jv-salmon", aria: "Second place" };
  if (r === 3) return { icon: Award, bg: "bg-jv-mint", aria: "Third place" };
  return null;
};
</script>

<template>
  <section aria-label="Final rankings">
    <article
      class="relative -rotate-[0.3deg] jv-border-rough bg-jv-white p-5 shadow-brutal sm:p-7 md:p-8"
    >
      <span
        class="absolute left-1/2 top-[-12px] z-10 h-4 w-16 -translate-x-1/2 rotate-[1deg] bg-jv-coral"
        aria-hidden="true"
      ></span>

      <header class="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h2
            class="font-headings text-[26px] leading-none text-jv-ink sm:text-[32px]"
          >
            Rankings
          </h2>
          <p
            class="mt-1.5 font-body text-[12px] font-bold text-jv-muted sm:text-[14px]"
          >
            Final scoreboard
          </p>
        </div>
        <p
          class="inline-flex items-center gap-2 rounded-full border-[2px] border-dashed border-jv-ink/30 px-3 py-1 font-feature text-[11px] font-black uppercase tracking-[0.16em] text-jv-muted sm:text-[12px]"
        >
          <span
            class="size-2 rounded-full bg-jv-coral"
            aria-hidden="true"
          ></span>
          {{ props.scoreboardData.length }} player{{
            props.scoreboardData.length === 1 ? "" : "s"
          }}
        </p>
      </header>

      <div
        class="my-4 border-t-[2px] border-dashed border-jv-ink/25 sm:my-5"
        aria-hidden="true"
      ></div>

      <ul
        v-if="props.scoreboardData.length > 0"
        class="flex flex-col gap-3 sm:gap-3.5"
      >
        <li
          v-for="(user, index) in props.scoreboardData"
          :key="`${user.username}-${user.rank}-${index}`"
          :class="[
            'flex items-center gap-3 border-[2px] border-jv-ink px-3 py-2.5 shadow-brutal-sm transition-transform sm:gap-4 sm:px-4 sm:py-3',
            user.username === userName && !isAdmin
              ? 'bg-jv-yellow-soft -rotate-[0.3deg]'
              : 'bg-jv-white',
          ]"
        >
          <!-- Rank: medal icon for top 3, number otherwise -->
          <span
            v-if="rankBadgeFor(user.rank)"
            :class="[
              'grid size-8 shrink-0 place-items-center rounded-[6px] border-[2px] border-jv-ink shadow-brutal-sm sm:size-10',
              rankBadgeFor(user.rank).bg,
            ]"
            :aria-label="rankBadgeFor(user.rank).aria"
          >
            <component
              :is="rankBadgeFor(user.rank).icon"
              class="size-4 sm:size-5"
              :stroke-width="2.4"
            />
          </span>
          <span
            v-else
            class="w-8 shrink-0 text-center font-feature text-[16px] font-black text-jv-ink tabular-nums sm:w-10 sm:text-[20px]"
            :aria-label="`Rank ${user.rank}`"
          >
            {{ user.rank }}
          </span>

          <span
            :class="[
              'grid size-10 shrink-0 place-items-center overflow-hidden rounded-[6px] border-[2px] border-jv-ink sm:size-11',
              avatarBgFor(index),
            ]"
          >
            <img
              :src="getAvatarUrlByName(user?.img_key)"
              :alt="`${user.firstname || user.username || 'Player'} avatar`"
              class="size-full object-cover"
            />
          </span>

          <div class="min-w-0 flex-1">
            <p
              class="flex flex-wrap items-baseline gap-x-2 gap-y-0.5 font-body text-[15px] font-black leading-tight text-jv-ink sm:text-[16px]"
            >
              <span class="truncate">{{ user.firstname }}</span>
              <span
                v-if="user.username === userName && !isAdmin"
                class="inline-flex shrink-0 items-center gap-1 rounded-full border-[2px] border-jv-ink bg-jv-white px-2 py-0.5 font-feature text-[10px] font-black uppercase tracking-[0.12em] text-jv-ink"
              >
                <Star class="size-3" :stroke-width="2.6" />
                You
              </span>
              <span
                v-else-if="
                  (isAdmin || user.username === userName) && user.username
                "
                class="truncate font-body text-[12px] font-bold text-jv-muted sm:text-[13px]"
              >
                @{{ user.username }}
              </span>
            </p>
          </div>

          <span
            class="shrink-0 font-feature text-[16px] font-black tabular-nums text-jv-ink sm:text-[20px]"
            :aria-label="`Score ${user.score}`"
          >
            {{ user.score }}
          </span>
        </li>
      </ul>

      <div
        v-else
        class="mx-auto mt-2 flex w-fit items-center gap-2 rounded-full border-[2px] border-dashed border-jv-ink/30 px-4 py-2 font-body text-[13px] font-bold text-jv-muted sm:text-[14px]"
      >
        <span
          class="size-2 rounded-full bg-jv-ink/30"
          aria-hidden="true"
        ></span>
        No rankings yet
      </div>
    </article>
  </section>
</template>
