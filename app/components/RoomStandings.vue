<script setup>
import ReportTabs from "@/components/reports/PageLayout.vue";
import { Trophy, Zap, ArrowRight } from "lucide-vue-next";
const props = defineProps({
  roomId: { type: String, required: true },
  detailed: { type: Boolean, default: false },
});
const emit = defineEmits(["report-loaded"]);
const route = useRoute();
const currentTab = computed(() =>
  route.query.tab === "quizzes" ? "quizzes" : "standings"
);
const changeReportTab = (tab) => {
  router.push({ path: route.path, query: { ...route.query, tab } });
};
const { apiUrl } = useRuntimeConfig().public;
const router = useRouter();
const { data, error, status, refresh } = await useFetch(
  () => `${apiUrl}/rooms/${props.roomId}/report`,
  { credentials: "include", headers: useRequestHeaders(["cookie"]) }
);
const report = computed(() => data.value?.data);
watch(
  report,
  (value) => {
    if (value) emit("report-loaded", value);
  },
  { immediate: true }
);
const refreshOnFocus = () => {
  if (!report.value?.ended) refresh();
};
onMounted(() => window.addEventListener("focus", refreshOnFocus));
onBeforeUnmount(() => window.removeEventListener("focus", refreshOnFocus));
</script>
<template>
  <section
    class="font-body text-jv-ink"
    :class="
      detailed
        ? 'space-y-8 sm:space-y-10'
        : 'space-y-5 rounded-[12px] border-[3px] border-jv-ink bg-jv-white p-5 shadow-brutal sm:p-7'
    "
    aria-label="Room results"
  >
    <ReportTabs
      v-if="detailed"
      :current-tab="currentTab"
      first-value="standings"
      second-value="quizzes"
      first-label="Overall standings"
      second-label="Played quizzes"
      @change-tab="changeReportTab"
    />
    <section
      v-if="detailed && report"
      class="jv-border-rough bg-jv-ink p-5 text-white shadow-brutal sm:p-7 md:p-8"
    >
      <div class="flex flex-col gap-6 lg:flex-row lg:items-center">
        <div class="min-w-0 flex-1">
          <p
            class="flex items-center gap-2 text-[12px] font-black uppercase tracking-[0.16em] text-jv-yellow"
          >
            <Zap class="size-4" :stroke-width="2.6" />Room recap
          </p>
          <h2
            class="mt-3 font-headings text-[28px] leading-tight text-white sm:text-[32px] md:text-[36px]"
          >
            {{
              report.ended
                ? "The final scores are in."
                : "The competition continues."
            }}
          </h2>
          <p
            class="mt-2 max-w-md text-[14px] font-semibold text-white/70 sm:text-[15px]"
          >
            {{
              report.ended
                ? "Review the final standings or revisit a completed quiz."
                : "Scores update after each quiz. Choose Played quizzes to revisit individual results."
            }}
          </p>
        </div>
        <div
          class="grid grid-cols-3 gap-3 jv-border-rough bg-jv-white p-4 text-jv-ink shadow-brutal-sm sm:p-5 lg:w-[340px]"
        >
          <div
            v-for="stat in [
              { label: 'Quizzes played', value: report.quizzes.length },
              { label: 'Total players', value: report.standings.length },
              { label: 'Room status', value: report.ended ? 'Ended' : 'Open' },
            ]"
            :key="stat.label"
            class="flex min-w-0 flex-col"
          >
            <p
              class="text-[10px] font-black uppercase tracking-[0.12em] text-jv-muted"
            >
              {{ stat.label }}
            </p>
            <p class="mt-auto pt-1 text-[20px] font-black text-jv-ink">
              {{ stat.value }}
            </p>
          </div>
        </div>
      </div>
    </section>
    <div v-if="!detailed">
      <p class="text-xs font-black uppercase tracking-widest text-jv-muted">
        {{ report?.ended ? "Room complete" : "Room results" }}
      </p>
      <div
        class="mt-2 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between sm:gap-5"
      >
        <div class="min-w-0 flex-1">
          <h2 class="font-headings text-3xl sm:text-4xl">
            {{ report?.ended ? "Final room standings" : "Overall standings" }}
          </h2>
          <p v-if="report" class="mt-2 text-sm text-jv-muted">
            Total points across {{ report.quizzes.length }} completed
            {{ report.quizzes.length === 1 ? "quiz" : "quizzes" }}. Missed
            quizzes add zero points.
          </p>
        </div>
        <p
          v-if="report?.standings.length"
          class="inline-flex w-fit items-center gap-2 rounded-full border-[2px] border-dashed border-jv-ink/30 px-3 py-1 font-feature text-[11px] font-black uppercase tracking-[0.16em] text-jv-muted sm:text-[12px]"
        >
          <span
            class="size-2 rounded-full bg-jv-coral"
            aria-hidden="true"
          ></span>
          {{ report.standings.length }} player{{
            report.standings.length === 1 ? "" : "s"
          }}
        </p>
      </div>
    </div>
    <p v-if="status === 'pending'" role="status">Loading room results…</p>
    <div
      v-else-if="error"
      role="alert"
      class="rounded border-2 border-jv-ink bg-jv-yellow-soft p-4"
    >
      Could not load room results.
      <button
        type="button"
        class="cursor-pointer font-bold underline"
        @click="refresh"
      >
        Try again
      </button>
    </div>
    <template v-else-if="report">
      <template v-if="!detailed || currentTab === 'standings'">
        <div
          v-if="!detailed && report.ended && report.standings.length"
          class="flex flex-wrap justify-center gap-3 py-4"
          aria-label="Final room winners"
        >
          <div
            v-for="(player, index) in report.standings.filter(
              (p) => p.rank <= 3
            )"
            :key="index"
            class="rounded-xl border-2 border-jv-ink px-5 py-4 text-center shadow-brutal-sm"
            :class="
              player.rank === 1
                ? 'bg-jv-yellow'
                : player.rank === 2
                  ? 'bg-jv-mint'
                  : 'bg-jv-salmon'
            "
          >
            <Trophy class="mx-auto mb-2 size-6" />
            <p class="font-headings text-2xl">{{ player.name }}</p>
            <p class="mt-1 font-bold">
              #{{ player.rank }} · {{ player.score }} points
            </p>
          </div>
        </div>
        <ScoreBoardTable
          v-if="report.standings.length"
          :flat="!detailed"
          :scoreboard-data="
            report.standings.map((p) => ({ ...p, firstname: p.name }))
          "
          :is-admin="report.is_host"
          :user-name="report.username"
          subtitle="Combined points from completed quizzes"
        />
        <p v-else class="py-6 text-center text-jv-muted">
          {{
            report.ended
              ? "No player scores were recorded in this room."
              : "Overall scores will appear after a quiz finishes."
          }}
        </p>
        <p
          v-if="report.quizzes.some((q) => !q.participant) && !report.is_host"
          class="text-sm text-jv-muted"
        >
          Missed quizzes count as zero in your overall score.
        </p>
      </template>
      <section v-else class="space-y-4" aria-label="Played quizzes">
        <p
          v-if="!report.quizzes.length"
          class="jv-border-rough bg-jv-white p-5 text-center text-jv-muted shadow-brutal-sm"
        >
          Completed quizzes will appear here.
        </p>
        <article
          v-for="(quiz, index) in report.quizzes"
          :key="quiz.id"
          class="flex flex-wrap items-center justify-between gap-4 jv-border-rough bg-jv-white p-5 shadow-brutal-sm sm:p-7"
        >
          <div class="min-w-0">
            <p
              class="text-xs font-black uppercase tracking-widest text-jv-coral"
            >
              Quiz {{ index + 1 }}
            </p>
            <h3 class="mt-2 break-words font-headings text-2xl font-bold">
              {{ quiz.title }}
            </h3>
          </div>
          <NuxtLink
            v-if="report.is_host || quiz.participant"
            :to="
              report.is_host
                ? `/admin/reports/${quiz.id}?room=${roomId}&title=${encodeURIComponent(quiz.title || '')}`
                : `/join/${encodeURIComponent(report.username)}/scoreboard?user_played_quiz=${quiz.participant}&view=quiz`
            "
            class="inline-flex h-11 items-center gap-2 rounded-full border-2 border-jv-ink bg-jv-white px-5 font-bold shadow-brutal-sm transition-transform hover:-rotate-1 active:shadow-none"
            >{{ report.is_host ? "View analysis" : "View results"
            }}<ArrowRight class="size-4" /></NuxtLink
          ><span v-else class="text-sm font-bold text-jv-muted"
            >Did not participate</span
          >
        </article>
      </section>
    </template>
  </section>
</template>
