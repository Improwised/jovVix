<script setup>
import { format } from "date-fns";
import { Layers } from "lucide-vue-next";
import ScoreBoardTable from "~/components/ScoreBoardTable.vue";
import { combinedRows, overallAccuracy } from "~~/composables/combined_report";

definePageMeta({ layout: "empty" });
useSeoMeta({
  title: "Combined Report - jovVix",
  robots: "noindex, nofollow",
});

const route = useRoute();
const { apiUrl } = useRuntimeConfig().public;
const sessionsQuery = computed(() => route.query.sessions || undefined);
const { data, error, pending } = await useFetch(
  () => `${apiUrl}/admin/reports/quizzes/${route.params.quizId}/combined`,
  {
    query: { sessions: sessionsQuery },
    credentials: "include",
    headers: useRequestHeaders(["cookie"]),
  }
);

const report = computed(() => data.value?.data);
const rows = computed(() => combinedRows(report.value));
const selectedCount = computed(
  () => report.value?.sessions.filter((s) => s.selected).length || 0
);
const stats = computed(() => [
  { label: "Participants", value: report.value?.participants ?? 0 },
  { label: "Sessions", value: selectedCount.value },
  {
    label: "Accuracy",
    value: `${overallAccuracy(report.value?.leaderboard)}%`,
  },
]);

// The checkboxes start from the sessions in the loaded report.
const checked = ref([]);
watch(
  report,
  (value) => {
    checked.value = (value?.sessions || [])
      .filter((s) => s.selected)
      .map((s) => s.id);
  },
  { immediate: true }
);

const generate = () => {
  navigateTo({
    path: route.path,
    query: { ...route.query, sessions: checked.value.join(",") },
  });
};

const when = (time) => format(new Date(time), "dd-MMM-yyyy HH:mm");
</script>

<template>
  <main
    class="flex min-h-screen flex-col gap-8 bg-jv-canvas px-4 py-5 text-jv-ink sm:gap-10 sm:px-6 md:px-8 md:py-6"
  >
    <header class="min-w-0">
      <h1
        class="font-headings text-[38px] leading-none min-[420px]:text-[44px] sm:text-[52px] md:text-[56px]"
      >
        Combined Report
      </h1>
      <div
        class="ml-1 mt-1 h-3 w-40 rounded-full border-b-[3px] border-jv-yellow sm:ml-2 sm:w-48"
        aria-hidden="true"
      ></div>
      <p
        class="mt-4 max-w-2xl text-[15px] font-semibold text-jv-muted sm:text-[16px]"
      >
        Rank everyone who played this quiz in the sessions you choose. A player
        who joined more than one session is ranked by their best play.
      </p>
    </header>

    <NuxtLink
      to="/admin/reports"
      class="inline-flex w-fit items-center gap-2 rounded-full border-2 border-jv-ink bg-jv-white px-5 py-2 font-bold text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-1 active:shadow-none"
      >← Back to reports</NuxtLink
    >

    <p
      v-if="pending && !report"
      role="status"
      class="jv-border-rough bg-jv-white p-5 text-center text-[18px] font-semibold text-jv-muted shadow-brutal-sm"
    >
      Loading...
    </p>
    <p
      v-else-if="error"
      role="alert"
      class="jv-border-rough bg-jv-white p-5 text-[18px] font-semibold text-jv-coral shadow-brutal-sm"
    >
      Could not load the combined report:
      {{ error.data?.data || error.data?.message || error.message }}
    </p>

    <template v-if="report">
      <section
        class="jv-border-rough bg-jv-ink p-5 text-white shadow-brutal sm:p-7 md:p-8"
      >
        <div class="flex flex-col gap-6 lg:flex-row lg:items-center">
          <div class="min-w-0 flex-1">
            <p
              class="flex items-center gap-2 text-[12px] font-black uppercase tracking-[0.16em] text-jv-yellow"
            >
              <Layers class="size-4" :stroke-width="2.6" />Combined leaderboard
            </p>
            <h2
              class="mt-3 font-headings text-[28px] leading-tight text-white sm:text-[32px] md:text-[36px]"
            >
              {{ decodeURI(report.title) }}
            </h2>
            <p
              class="mt-2 text-[14px] font-semibold text-white/70 sm:text-[15px]"
            >
              {{ report.plays }} plays across {{ selectedCount }} of
              {{ report.sessions.length }} sessions
            </p>
          </div>
          <div
            class="grid grid-cols-3 gap-3 jv-border-rough bg-jv-white p-4 text-jv-ink shadow-brutal-sm sm:p-5 lg:w-[340px]"
          >
            <div
              v-for="stat in stats"
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

      <form
        class="jv-border-rough bg-jv-white p-5 shadow-brutal-sm sm:p-6"
        aria-label="Choose sessions"
        @submit.prevent="generate"
      >
        <fieldset>
          <legend class="font-headings text-[24px] leading-none sm:text-[28px]">
            Sessions
          </legend>
          <p class="mt-2 text-[14px] font-semibold text-jv-muted">
            Choose at least two sessions, then generate the report.
          </p>
          <div class="mt-4 flex flex-col gap-2">
            <label
              v-for="(session, index) in report.sessions"
              :key="session.id"
              class="flex cursor-pointer items-center gap-3 rounded-[6px] border-2 border-jv-ink px-3 py-2 font-bold hover:bg-jv-yellow/30"
            >
              <input
                v-model="checked"
                type="checkbox"
                :value="session.id"
                class="size-4 accent-jv-ink"
              />
              <span>Session {{ index + 1 }}</span>
              <span class="text-[14px] font-semibold text-jv-muted">
                {{ when(session.activated_from) }} ·
                {{ session.participants }} player{{
                  session.participants === 1 ? "" : "s"
                }}
              </span>
            </label>
          </div>
        </fieldset>
        <button
          type="submit"
          :disabled="checked.length < 2"
          class="mt-5 cursor-pointer rounded-full border-2 border-jv-ink bg-jv-yellow px-5 py-2 font-bold text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-1 disabled:cursor-not-allowed disabled:opacity-50"
        >
          Generate combined report
        </button>
      </form>

      <ScoreBoardTable
        :scoreboard-data="rows"
        is-admin
        user-name=""
        title="Combined leaderboard"
        :subtitle="`Best play per player across ${selectedCount} sessions`"
      />
    </template>
  </main>
</template>
