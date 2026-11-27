<script setup>
import {
  AlertTriangle,
  Award,
  Crown,
  Home,
  Medal,
  Sparkles,
  Trophy,
  Zap,
} from "lucide-vue-next";
import NavigationLink from "@/components/common/NavigationLink.vue";

definePageMeta({ layout: "empty", hideSidebar: true });
useSeoMeta({ title: "Solo results - jovVix", robots: "noindex, nofollow" });

const route = useRoute();
const { apiUrl } = useRuntimeConfig().public;
const { data, error, pending, refresh } = await useFetch(
  `${apiUrl}/solo/results/${route.params.attempt}`,
  { credentials: "include", server: false }
);

const results = computed(() => data.value?.data);
const players = computed(() => results.value?.players || []);
const me = computed(() => players.value.find((p) => p.is_you));

// ScoreBoardTable highlights the row whose username matches this sentinel.
const YOU = "__you__";
const toRow = (player) => ({
  ...player,
  firstname: player.name,
  username: player.is_you ? YOU : player.name,
  img_key: player.img_key,
});

const userStatistics = computed(() => {
  const attempt = results.value?.attempt;
  if (!attempt) return {};
  const totalQuestions = attempt.total_questions ?? attempt.total;
  return {
    accuracy: totalQuestions
      ? Math.round((attempt.correct / totalQuestions) * 100)
      : 0,
    totalScore: attempt.score,
    correctAnwers: attempt.correct,
    wrongAnwers: attempt.total - attempt.correct,
    unAttemptedQuestions: Math.max(0, totalQuestions - attempt.total),
  };
});

const rankBadge = computed(() => {
  const r = me.value?.rank;
  if (r === 1) return { label: "1st Place", bg: "bg-jv-yellow", icon: Crown };
  if (r === 2) return { label: "2nd Place", bg: "bg-jv-salmon", icon: Medal };
  if (r === 3) return { label: "3rd Place", bg: "bg-jv-mint", icon: Award };
  return { label: `Rank #${r}`, bg: "bg-jv-lavender", icon: Trophy };
});

const leaderboardRows = computed(() => {
  const top = players.value.slice(0, 3);
  // Decide by position, not rank: tied scores share a rank, so rank <= 3 can
  // be true even when you sit below the third displayed row.
  const inTop = !me.value || top.some((p) => p.is_you);
  if (inTop) return top.map(toRow);
  return [...top.map(toRow), { divider: true }, toRow(me.value)];
});
</script>

<template>
  <ClientOnly>
    <main
      class="min-h-screen bg-jv-canvas px-3 py-6 text-jv-ink sm:px-6 sm:py-8 md:px-10"
    >
      <div class="mx-auto flex w-full max-w-[1180px] flex-col gap-6 sm:gap-8">
        <header>
          <h1
            class="font-headings text-[38px] leading-none text-jv-ink sm:text-[52px]"
          >
            Solo results
          </h1>
          <p
            class="mt-4 max-w-2xl text-[15px] font-semibold text-jv-muted sm:text-[16px]"
          >
            Review your score and see how you rank.
          </p>
        </header>
        <NuxtLink
          to="/#public-quiz"
          class="inline-flex w-fit items-center gap-2 rounded-full border-2 border-jv-ink bg-jv-white px-5 py-2 font-bold text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-1 active:shadow-none"
          >← Back to public quizzes</NuxtLink
        >
        <div
          v-if="pending"
          class="flex min-h-[40vh] items-center justify-center"
          role="status"
          aria-live="polite"
        >
          <div
            class="inline-flex items-center gap-2 rounded-full border-[2px] border-dashed border-jv-ink/30 bg-jv-white px-5 py-2.5 font-body text-[14px] font-bold text-jv-muted shadow-brutal-sm"
          >
            <Sparkles class="size-4" :stroke-width="2.4" />
            Loading solo results…
          </div>
        </div>

        <div
          v-else-if="error"
          role="alert"
          class="jv-border-rough flex items-start gap-3 border-2 border-jv-coral bg-jv-coral/10 p-4 font-body text-jv-ink"
        >
          <AlertTriangle
            class="mt-0.5 size-5 shrink-0 text-jv-coral"
            :stroke-width="2.4"
          />
          <div>
            <p>Unable to load your completed solo results.</p>
            <NavigationLink variant="result" class="mt-3" @click="refresh()">
              Try again
            </NavigationLink>
          </div>
        </div>

        <template v-else-if="results">
          <header
            data-testid="solo-results-hero"
            class="relative -rotate-[0.4deg] jv-border-rough bg-jv-white p-5 shadow-brutal sm:p-7"
          >
            <span
              class="absolute left-1/2 top-[-12px] z-10 h-4 w-16 -translate-x-1/2 rotate-[2deg] bg-jv-coral"
              aria-hidden="true"
            ></span>

            <div class="flex flex-wrap items-start justify-between gap-3">
              <p
                class="font-body text-[11px] font-black uppercase tracking-[0.22em] text-jv-muted sm:text-[12px]"
              >
                Solo complete · {{ results.attempt.title }}
              </p>
              <div v-if="me" class="flex flex-wrap items-start gap-2 sm:gap-3">
                <span
                  :class="[
                    'inline-flex h-9 items-center gap-2 rounded-full border-[2px] border-jv-ink px-3 font-feature text-[13px] font-black text-jv-ink shadow-brutal-sm sm:text-[14px]',
                    rankBadge.bg,
                  ]"
                >
                  <component
                    :is="rankBadge.icon"
                    class="size-4"
                    :stroke-width="2.4"
                  />
                  {{ rankBadge.label }}
                </span>
                <span
                  class="inline-flex h-9 items-center gap-1.5 rounded-full border-[2px] border-jv-ink bg-jv-white px-3 shadow-brutal-sm"
                >
                  <span
                    class="font-body text-[10px] font-bold uppercase leading-none tracking-[0.12em] text-jv-muted"
                    >Score</span
                  >
                  <span
                    class="font-feature text-[16px] font-black leading-none tabular-nums text-jv-ink sm:text-[18px]"
                  >
                    {{ results.attempt.score }}
                  </span>
                </span>
              </div>
            </div>

            <div class="mt-1 min-w-0">
              <h2
                class="font-headings text-[28px] leading-tight text-jv-ink sm:text-[40px]"
              >
                Great game!
              </h2>
              <p
                class="mt-2 font-body text-[14px] font-bold text-jv-muted sm:text-[16px]"
              >
                {{ results.attempt.correct }} /
                {{ results.attempt.total_questions ?? results.attempt.total }}
                correct<template v-if="me">
                  · Personal best: {{ me.score }} points</template
                >
              </p>
              <p
                v-if="players.length === 1"
                class="mt-1 font-body text-[14px] font-bold text-jv-muted sm:text-[16px]"
              >
                You're the first to finish this quiz!
              </p>
            </div>
          </header>

          <QuizStatisticsBadges :user-statistics="userStatistics" />

          <ScoreBoardTable
            :scoreboard-data="leaderboardRows"
            :user-name="YOU"
            title="Overall Rankings"
            :rank-icons="false"
            :total-players="players.length"
            subtitle="Best attempt per player · Ties share a rank"
          />

          <div
            class="flex flex-wrap items-center justify-center gap-3"
            role="group"
            aria-label="Solo result actions"
          >
            <NavigationLink
              url="/"
              variant="result"
              aria-label="Return to home"
            >
              <Home class="size-4" :stroke-width="2.4" />
              Home
            </NavigationLink>
            <NavigationLink
              :url="`/admin/arrange/new?public=1&mode=solo&quiz_id=${results.attempt.quiz_id}`"
              variant="result"
              class="bg-jv-yellow"
              aria-label="Play this quiz again"
            >
              <Zap class="size-4" :stroke-width="2.4" />
              Beat your score
            </NavigationLink>
          </div>
        </template>
      </div>
    </main>
  </ClientOnly>
</template>
