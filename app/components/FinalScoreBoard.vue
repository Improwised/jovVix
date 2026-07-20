<script setup>
import { usePush } from "notivue";
import {
  BarChart3,
  Crown,
  ArrowRight,
  Sparkles,
  Home,
  Trophy,
  Medal,
  Award,
} from "lucide-vue-next";
import ScoreBoardTable from "./ScoreBoardTable.vue";
import { useSessionStore } from "~~/store/session";
import { useUsersStore } from "~~/store/users";
import { getAvatarUrlByName } from "~~/composables/avatar";

const url = useRuntimeConfig().public;
const scoreboardData = ref([]);
const route = useRoute();
const router = useRouter();
const activeQuizId = ref("");
const toast = usePush();
const app = useNuxtApp();
useSystemEnv();

const analysisData = ref([]);
const userAnalysisEndpoint = "/analytics_board/user";
const requestPending = ref(false);
const analysisPending = ref(false);
const userStatistics = ref({});
const winnerUI = computed(() => route.query.winner_ui || false);

const usersStore = useUsersStore();
const authPending = ref(true);
const isLoggedIn = computed(
  () => usersStore.getUserData()?.role === "admin-user"
);
const winningSound = ref(null);

const props = defineProps({
  userURL: {
    default: "",
    type: String,
    required: true,
  },
  isAdmin: {
    default: false,
    type: Boolean,
    required: false,
  },
  userName: {
    type: String,
    required: false,
    default: "",
  },
  userPlayedQuiz: {
    type: String,
    required: false,
    default: "",
  },
});

const getFinalScoreboardDetails = async (endpoint) => {
  try {
    requestPending.value = true;
    const response = await $fetch(`${url.apiUrl}${endpoint}`, {
      method: "GET",
      credentials: "include",
      headers: { Accept: "application/json" },
    });
    scoreboardData.value = response?.data || [];
  } catch (error) {
    toast.error(
      error?.data?.message || error?.message || "Failed to load scoreboard"
    );
  } finally {
    requestPending.value = false;
  }
};

const getAnalysisDetails = async () => {
  try {
    analysisPending.value = true;
    const response = await $fetch(
      `${url.apiUrl}${userAnalysisEndpoint}?user_played_quiz=${props.userPlayedQuiz}`,
      {
        method: "GET",
        credentials: "include",
      }
    );
    const data = response?.data || [];
    analysisData.value = data;
    userStatistics.value = questionsAnalysis(data);
  } catch (error) {
    if (error?.statusCode === 401 || error?.response?.status === 401) {
      toast.error(app.$Unauthorized || "Unauthorized");
    } else {
      toast.error(
        error?.data?.message || error?.message || "Failed to load analysis"
      );
    }
  } finally {
    analysisPending.value = false;
  }
};

const loadData = () => {
  if (props.isAdmin) {
    activeQuizId.value = route.query.aqi || "";
    getFinalScoreboardDetails(
      `${props.userURL}?active_quiz_id=${activeQuizId.value}`
    );
  } else {
    getFinalScoreboardDetails(
      `${props.userURL}?user_played_quiz=${props.userPlayedQuiz}`
    );
    getAnalysisDetails();
    setUserDataStore().finally(() => {
      authPending.value = false;
    });
  }
};

if (process.client) {
  loadData();
} else {
  onMounted(() => loadData());
}

const showAnalysis = () => {
  router.push({ path: `/admin/reports/${activeQuizId.value}` });
};

const changeUI = (value) => {
  navigateTo({ path: route.path, query: { ...route.query, winner_ui: value } });
};

// Browsers reject autoplay without a user gesture and surface a noisy promise
// rejection. We can't make a gesture from SSR, so swallow that specific error.
const tryPlay = (audio) => {
  if (!audio) return;
  const result = audio.play();
  if (result && typeof result.catch === "function") {
    result.catch(() => {
      /* autoplay blocked — fine */
    });
  }
};

watch(
  winnerUI,
  (newValue) => {
    const music = newValue == "true";
    if (!music && winningSound.value) {
      winningSound.value.pause();
    } else if (music && winningSound.value) {
      tryPlay(winningSound.value);
    }
  },
  { deep: true, immediate: true }
);

onMounted(() => {
  if (process.client) {
    winningSound.value = new Audio("/music/winning.mp3");
    if (winnerUI.value == "true") {
      tryPlay(winningSound.value);
    }
  }
});

const winners = computed(() => scoreboardData.value.slice(0, 3));

const userRow = computed(() =>
  scoreboardData.value.find((u) => u?.username === props.userName)
);
const userRank = computed(
  () => userRow.value?.rank ?? userStatistics.value?.rank ?? 0
);
const userScore = computed(
  () => userRow.value?.score ?? userStatistics.value?.totalScore ?? 0
);
const playerCount = computed(() => scoreboardData.value.length);

const rankBadge = computed(() => {
  const r = Number(userRank.value);
  if (r === 1) return { label: "1st Place", bg: "bg-jv-yellow", icon: Crown };
  if (r === 2) return { label: "2nd Place", bg: "bg-jv-salmon", icon: Medal };
  if (r === 3) return { label: "3rd Place", bg: "bg-jv-mint", icon: Award };
  if (r > 0) return { label: `Rank #${r}`, bg: "bg-jv-lavender", icon: Trophy };
  return null;
});

// Podium config — drives the 2nd / 1st / 3rd columns in the winner view.
// `heightClass` controls the colored bar height (1st is tallest, 3rd shortest).
const sessionStore = useSessionStore();
const tournamentTitle = computed(() => {
  const fromQuery =
    typeof route.query.title === "string" ? route.query.title.trim() : "";
  if (fromQuery) return fromQuery;
  const stored = sessionStore.getActiveQuizTitle?.();
  if (stored) return stored;
  return "";
});

const podiumOrder = computed(() => {
  const w = scoreboardData.value;
  // Render as 2nd, 1st, 3rd so the visual layout matches the physical podium.
  return [
    {
      slot: "second",
      winner: w[1],
      rank: 2,
      bg: "bg-jv-mint",
      ring: "ring-jv-mint",
      heightClass: "h-[180px] sm:h-[210px]",
      label: "2nd",
    },
    {
      slot: "first",
      winner: w[0],
      rank: 1,
      bg: "bg-jv-yellow",
      ring: "ring-jv-yellow",
      heightClass: "h-[240px] sm:h-[280px]",
      label: "1st",
      highlight: true,
    },
    {
      slot: "third",
      winner: w[2],
      rank: 3,
      bg: "bg-jv-salmon",
      ring: "ring-jv-salmon",
      heightClass: "h-[140px] sm:h-[170px]",
      label: "3rd",
    },
  ];
});

const podiumName = (winner) =>
  winner?.firstname || winner?.username || "Player";
</script>

<template>
  <!-- WINNER PODIUM (admin only, when winner_ui=true) -->
  <ClientOnly v-if="winnerUI == 'true' && props.isAdmin">
    <div
      class="relative flex min-h-screen flex-col overflow-hidden bg-jv-canvas"
    >
      <!-- Continuous victory-stage confetti, blasted inward from both edges.
           Mounts only with this view, so unmounting handles cleanup. -->
      <WinnerConfetti />

      <!-- TOP BAR: logo + tournament title + next button -->
      <header
        class="relative z-10 flex items-center justify-between gap-3 px-4 py-4 sm:px-8 sm:py-5 md:px-12"
        role="banner"
      >
        <div class="flex min-w-0 items-center gap-3 sm:gap-4">
          <div class="hidden min-w-0 flex-col sm:flex">
            <p
              class="font-body text-[10px] font-bold uppercase tracking-[0.16em] text-jv-muted sm:text-[11px]"
            >
              Final Results
            </p>
            <!-- Title can come from the persisted session store; render only on
               the client to avoid SSR/CSR hydration mismatches. -->
            <ClientOnly>
              <h2
                v-if="tournamentTitle"
                class="relative truncate font-headings text-[20px] leading-tight sm:text-[26px]"
              >
                <span class="relative z-10">{{ tournamentTitle }}</span>
                <span
                  class="absolute bottom-[2px] left-0 z-0 h-[6px] w-[80%] max-w-[180px] rotate-[-1deg] bg-jv-yellow/70"
                  aria-hidden="true"
                ></span>
              </h2>
            </ClientOnly>
          </div>
        </div>

        <div class="flex items-center gap-2 sm:gap-3">
          <NuxtLink
            to="/"
            class="inline-flex h-10 items-center justify-center gap-1.5 rounded-[8px] border-[2px] border-jv-ink bg-jv-white px-4 font-body text-[13px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[-1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-11 sm:px-5 sm:text-[14px]"
            aria-label="Return to home"
          >
            <Home class="size-4" :stroke-width="2.4" />
            Home
          </NuxtLink>
          <button
            type="button"
            class="inline-flex h-10 items-center justify-center gap-1.5 rounded-[8px] border-[2px] border-jv-ink bg-jv-white px-4 font-body text-[13px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[-1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-11 sm:px-5 sm:text-[14px]"
            aria-label="Continue to scoreboard table view"
            @click="changeUI(false)"
          >
            Next
            <ArrowRight class="size-4" :stroke-width="2.4" />
          </button>
        </div>
      </header>

      <main
        id="main-content"
        class="relative z-10 mx-auto flex w-full max-w-[1180px] flex-1 flex-col px-3 pb-6 pt-2 sm:px-6 sm:pb-10 md:px-10"
        role="main"
        aria-label="Quiz winners podium"
      >
        <div class="relative flex w-full flex-1 flex-col justify-evenly">
          <!-- Title with confetti-like decorations -->
          <div class="relative flex items-center justify-center">
            <span
              class="absolute left-[8%] top-[10%] hidden h-5 w-1.5 -rotate-[30deg] rounded-full bg-jv-coral sm:block"
              aria-hidden="true"
            ></span>
            <span
              class="absolute left-[20%] top-[60%] hidden size-2.5 rounded-full bg-jv-yellow sm:block"
              aria-hidden="true"
            ></span>
            <span
              class="absolute right-[8%] top-[15%] hidden h-5 w-1.5 rotate-[20deg] rounded-full bg-jv-mint sm:block"
              aria-hidden="true"
            ></span>
            <h2
              class="text-center text-[40px] leading-none text-jv-ink sm:text-[44px] md:text-[56px]"
            >
              Quiz Winners!
            </h2>
          </div>

          <!-- Loading -->
          <div
            v-if="requestPending"
            class="mt-12 flex justify-center"
            role="status"
            aria-live="polite"
          >
            <span
              class="inline-flex items-center gap-2 rounded-full border-[2px] border-dashed border-jv-ink/30 px-5 py-2 font-body text-[14px] font-bold text-jv-muted"
            >
              <Sparkles class="size-4" :stroke-width="2.4" />
              Loading winners…
            </span>
          </div>

          <!-- Podium -->
          <section
            v-else-if="winners.length"
            class="relative mx-auto mt-10 flex w-full max-w-[640px] items-end justify-center gap-3 sm:mt-14 sm:gap-5"
            aria-label="Top three winners podium"
          >
            <div
              v-for="step in podiumOrder"
              :key="step.slot"
              class="flex flex-1 flex-col items-center"
            >
              <template v-if="step.winner">
                <!-- Name tag -->
                <span
                  class="relative z-10 inline-flex max-w-full rotate-[-1deg] items-center justify-center rounded-[6px] border-[2px] border-jv-ink bg-jv-ink px-2.5 py-1 font-feature text-[11px] font-black uppercase tracking-[0.06em] text-white shadow-brutal-sm sm:text-[12px]"
                >
                  <span class="block max-w-[88px] truncate sm:max-w-[120px]">
                    {{ podiumName(step.winner) }}
                  </span>
                </span>

                <!-- Crown for 1st -->
                <Crown
                  v-if="step.rank === 1"
                  class="relative z-10 mt-1 size-6 text-jv-yellow-2 drop-shadow sm:size-7"
                  :stroke-width="2.4"
                  fill="currentColor"
                />

                <!-- Avatar: colored ring matches the podium color of this rank. -->
                <span
                  :class="[
                    'relative z-10 mt-1 grid place-items-center overflow-hidden rounded-full border-[3px] border-jv-ink bg-jv-white shadow-brutal-sm ring-4 ring-offset-[2px] ring-offset-jv-canvas',
                    step.ring,
                    step.rank === 1
                      ? 'size-[72px] sm:size-[88px]'
                      : 'size-[56px] sm:size-[68px]',
                  ]"
                >
                  <img
                    :src="getAvatarUrlByName(step.winner?.img_key)"
                    :alt="`Avatar for ${podiumName(step.winner)}`"
                    class="size-full object-cover"
                  />
                </span>

                <!-- Step / pedestal -->
                <div
                  :class="[
                    'relative mt-2 flex w-full items-center justify-center border-[2px] border-jv-ink shadow-brutal-sm',
                    step.bg,
                    step.heightClass,
                  ]"
                >
                  <!-- Spotlight glow behind 1st place -->
                  <span
                    v-if="step.highlight"
                    class="pointer-events-none absolute -top-[76%] left-1/2 h-[140%] w-[160%] -translate-x-1/2 bg-gradient-to-b from-jv-yellow/45 via-jv-yellow/25 to-transparent [clip-path:polygon(30%_0%,70%_0%,100%_100%,0%_100%)]"
                    aria-hidden="true"
                  ></span>
                  <span
                    class="relative font-feature font-black text-jv-ink"
                    :class="
                      step.rank === 1
                        ? 'text-[40px] sm:text-[56px]'
                        : 'text-[28px] sm:text-[36px]'
                    "
                  >
                    {{ step.label }}
                  </span>
                </div>
              </template>

              <!-- Placeholder when fewer than 3 winners exist -->
              <template v-else>
                <span
                  class="inline-flex items-center gap-1 rounded-full border-[2px] border-dashed border-jv-ink/25 px-3 py-1 font-body text-[11px] font-bold text-jv-muted"
                >
                  —
                </span>
                <span
                  class="mt-1 grid place-items-center overflow-hidden rounded-full border-[3px] border-dashed border-jv-ink/25 bg-jv-white"
                  :class="
                    step.rank === 1
                      ? 'size-[72px] sm:size-[88px]'
                      : 'size-[56px] sm:size-[68px]'
                  "
                  aria-hidden="true"
                ></span>
                <div
                  :class="[
                    'relative mt-2 flex w-full items-center justify-center border-[2px] border-dashed border-jv-ink/25 bg-jv-white',
                    step.heightClass,
                  ]"
                >
                  <span
                    class="font-feature font-black text-jv-ink/25"
                    :class="
                      step.rank === 1
                        ? 'text-[40px] sm:text-[56px]'
                        : 'text-[28px] sm:text-[36px]'
                    "
                  >
                    {{ step.label }}
                  </span>
                </div>
              </template>
            </div>
          </section>

          <!-- Empty state -->
          <div v-else class="mt-12 flex justify-center" role="status">
            <span
              class="inline-flex items-center gap-2 rounded-full border-[2px] border-dashed border-jv-ink/30 px-4 py-2 font-body text-[13px] font-bold text-jv-muted"
            >
              <span
                class="size-2 rounded-full bg-jv-ink/30"
                aria-hidden="true"
              ></span>
              No winners yet
            </span>
          </div>
        </div>

        <!-- Status pill + copyright -->
        <div
          v-if="!requestPending"
          class="mt-6 flex justify-center"
          aria-live="polite"
        >
          <span
            class="inline-flex items-center gap-2 rounded-full border-[2px] border-dashed border-jv-ink/30 bg-jv-white px-4 py-2 font-body text-[12px] font-bold text-jv-muted sm:text-[13px]"
          >
            <span
              class="size-2 rounded-full border border-jv-ink bg-jv-mint"
              aria-hidden="true"
            ></span>
            Tournament completed successfully
          </span>
        </div>
      </main>
    </div>
  </ClientOnly>

  <!-- TABLE + ANALYSIS VIEW -->
  <ClientOnly v-else>
    <div
      v-if="requestPending"
      class="flex min-h-screen items-center justify-center bg-jv-canvas px-4"
      role="status"
      aria-live="polite"
    >
      <div
        class="inline-flex items-center gap-2 rounded-full border-[2px] border-dashed border-jv-ink/30 bg-jv-white px-5 py-2.5 font-body text-[14px] font-bold text-jv-muted shadow-brutal-sm"
      >
        <Sparkles class="size-4" :stroke-width="2.4" />
        Loading scoreboard…
      </div>
    </div>

    <main
      v-else
      id="main-content"
      class="min-h-screen bg-jv-canvas px-3 py-6 text-jv-ink sm:px-6 sm:py-8 md:px-10"
      role="main"
    >
      <div class="mx-auto flex w-full max-w-[1180px] flex-col gap-6 sm:gap-8">
        <!-- HERO: celebrate the user / show admin title -->
        <header
          class="relative -rotate-[0.4deg] jv-border-rough bg-jv-white p-5 shadow-brutal sm:p-7"
        >
          <span
            class="absolute left-1/2 top-[-12px] z-10 h-4 w-16 -translate-x-1/2 rotate-[2deg] bg-jv-coral"
            aria-hidden="true"
          ></span>

          <p
            class="font-body text-[11px] font-black uppercase tracking-[0.22em] text-jv-muted sm:text-[12px]"
          >
            Quiz Complete
          </p>

          <div
            class="mt-1 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between sm:gap-5"
          >
            <div class="min-w-0 flex-1">
              <h1
                class="font-headings text-[28px] leading-tight text-jv-ink sm:text-[40px]"
              >
                {{ props.isAdmin ? "Final Scoreboard" : "Great game!" }}
              </h1>
              <p
                v-if="!props.isAdmin"
                class="mt-2 font-body text-[14px] font-bold text-jv-muted sm:text-[16px]"
              >
                Here's how everyone stacked up.
              </p>
              <p
                v-else
                class="mt-2 font-body text-[14px] font-bold text-jv-muted sm:text-[16px]"
              >
                {{ playerCount }} player{{ playerCount === 1 ? "" : "s" }}
                completed the quiz.
              </p>
            </div>

            <div
              v-if="!props.isAdmin && rankBadge"
              class="flex flex-wrap items-center gap-2 sm:gap-3"
            >
              <span
                :class="[
                  'inline-flex items-center gap-2 rounded-full border-[2px] border-jv-ink px-3 py-1.5 font-feature text-[13px] font-black text-jv-ink shadow-brutal-sm sm:text-[14px]',
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
                class="inline-flex items-baseline gap-1.5 rounded-full border-[2px] border-jv-ink bg-jv-white px-3 py-1.5 shadow-brutal-sm"
              >
                <span
                  class="font-body text-[10px] font-bold uppercase tracking-[0.12em] text-jv-muted"
                  >Score</span
                >
                <span
                  class="font-feature text-[16px] font-black tabular-nums text-jv-ink sm:text-[18px]"
                >
                  {{ userScore }}
                </span>
              </span>
            </div>
          </div>
        </header>

        <!-- USER STATS (non-admin only) -->
        <section
          v-if="!props.isAdmin"
          aria-label="Your quiz statistics"
          class="contents"
        >
          <div
            v-if="analysisPending"
            class="mx-auto flex w-fit items-center gap-2 rounded-full border-[2px] border-dashed border-jv-ink/30 bg-jv-white px-5 py-2.5 font-body text-[13px] font-bold text-jv-muted shadow-brutal-sm sm:text-[14px]"
            role="status"
            aria-live="polite"
          >
            <Sparkles class="size-4" :stroke-width="2.4" />
            Crunching your stats…
          </div>
          <QuizStatisticsBadges v-else :user-statistics="userStatistics" />
        </section>

        <!-- RANKINGS -->
        <ScoreBoardTable
          v-if="scoreboardData"
          :scoreboard-data="scoreboardData"
          :is-admin="props.isAdmin"
          :user-name="props.userName"
        />

        <!-- ADMIN CONTROLS -->
        <div
          v-if="props.isAdmin"
          class="flex flex-col items-stretch gap-3 sm:flex-row sm:justify-center"
          role="group"
          aria-label="Admin controls"
        >
          <NuxtLink
            to="/"
            class="inline-flex h-11 items-center justify-center gap-2 rounded-full border-[2px] border-jv-ink bg-jv-white px-6 font-body text-[14px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-12 sm:text-[15px]"
            aria-label="Return to home"
          >
            <Home class="size-4" :stroke-width="2.4" />
            Home
          </NuxtLink>
          <NuxtLink
            class="inline-flex h-11 items-center justify-center gap-2 rounded-full border-[2px] border-jv-ink bg-jv-white px-6 font-body text-[14px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-12 sm:text-[15px]"
            aria-label="View detailed quiz analysis"
            @click="showAnalysis"
          >
            <BarChart3 class="size-4" :stroke-width="2.4" />
            Show Analysis
          </NuxtLink>
          <NuxtLink
            class="inline-flex h-11 items-center justify-center gap-2 rounded-full border-[2px] border-jv-ink bg-jv-coral px-6 font-body text-[14px] font-black text-white shadow-brutal-sm transition-transform hover:-rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-12 sm:text-[15px]"
            aria-label="Switch to winners podium view"
            @click="changeUI(true)"
          >
            <Crown class="size-4" :stroke-width="2.4" />
            Show Winners
          </NuxtLink>
        </div>

        <!-- USER ANALYSIS -->
        <template v-if="!props.isAdmin">
          <div
            v-if="analysisPending || authPending"
            class="mx-auto flex w-fit items-center gap-2 rounded-full border-[2px] border-dashed border-jv-ink/30 bg-jv-white px-5 py-2.5 font-body text-[13px] font-bold text-jv-muted shadow-brutal-sm sm:text-[14px]"
            role="status"
            aria-live="polite"
          >
            <Sparkles class="size-4" :stroke-width="2.4" />
            Loading question review…
          </div>
          <QuizAnalysis v-else-if="isLoggedIn" :data="analysisData" />
          <QuizAnalysisLoginGate v-else />

          <div
            class="mb-2 flex flex-col items-stretch gap-3 sm:flex-row sm:justify-center"
            role="group"
            aria-label="Quiz end actions"
          >
            <NuxtLink
              to="/"
              class="inline-flex h-11 items-center justify-center gap-2 rounded-full border-[2px] border-jv-ink bg-jv-white px-6 font-body text-[14px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-12 sm:text-[15px]"
              aria-label="Return to home"
            >
              <Home class="size-4" :stroke-width="2.4" />
              Home
            </NuxtLink>
            <NuxtLink
              to="/join"
              class="inline-flex h-11 items-center justify-center gap-2 rounded-full border-[2px] border-jv-ink bg-jv-yellow px-6 font-body text-[14px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-12 sm:text-[15px]"
              aria-label="Join another quiz"
            >
              <ArrowRight class="size-4" :stroke-width="2.4" />
              Play another
            </NuxtLink>
          </div>
        </template>
      </div>
    </main>
  </ClientOnly>
</template>

<style scoped>
.jv-grid {
  background-image: radial-gradient(var(--jv-charcoal) 1px, transparent 1px);
  background-size: 20px 20px;
  opacity: 0.08;
}
</style>
