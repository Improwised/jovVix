<template>
  <main
    class="flex min-h-screen flex-col gap-8 bg-jv-canvas px-4 py-5 sm:gap-10 sm:px-6 md:px-8 md:py-6"
  >
    <div
      class="flex flex-col gap-5 xl:flex-row xl:items-end xl:justify-between"
    >
      <div class="min-w-0">
        <h1
          class="font-headings text-[38px] leading-none text-jv-ink min-[420px]:text-[44px] sm:text-[52px] md:text-[56px]"
        >
          Reports
        </h1>
        <div
          class="ml-1 mt-1 h-3 w-40 rounded-full border-b-[3px] border-jv-yellow sm:ml-2 sm:w-48"
          aria-hidden="true"
        ></div>
      </div>

      <div
        v-if="reportType === 'quizzes'"
        class="flex flex-col sm:flex-row gap-4 items-center"
      >
        <label
          class="flex h-11 w-full sm:w-auto min-w-[200px] rotate-[1deg] items-center gap-2.5 jv-border-rough bg-jv-white px-3 shadow-brutal-sm sm:h-12 sm:px-4 cursor-pointer hover:bg-jv-white/80 transition-colors"
        >
          <Calendar
            class="size-5 shrink-0 text-jv-ink/40"
            :stroke-width="2.4"
          />
          <input
            v-model="date"
            type="datetime-local"
            class="h-full min-w-0 flex-1 bg-transparent text-[15px] font-semibold text-jv-ink outline-none sm:text-[16px]"
          />
        </label>
        <label
          class="flex h-11 w-full sm:w-auto min-w-[220px] rotate-[-1deg] items-center gap-2.5 jv-border-rough bg-jv-white px-3 shadow-brutal-sm sm:h-12 sm:px-4 hover:bg-jv-white/80 transition-colors"
        >
          <Search class="size-5 shrink-0 text-jv-ink/40" :stroke-width="2.4" />
          <input
            v-model="nameInput"
            type="search"
            class="h-full min-w-0 flex-1 bg-transparent text-[15px] font-semibold text-jv-ink outline-none placeholder:text-jv-ink/45 sm:text-[16px]"
            placeholder="Search quiz"
          />
        </label>
      </div>
      <label
        v-if="reportType === 'rooms'"
        class="flex h-11 w-full items-center gap-2.5 rotate-[1deg] jv-border-rough bg-jv-white px-3 shadow-brutal-sm sm:h-12 sm:w-auto"
      >
        <Calendar class="size-5 shrink-0 text-jv-ink/40" :stroke-width="2.4" />
        <input
          v-model="roomDate"
          type="datetime-local"
          aria-label="Filter rooms by start date"
          class="h-full min-w-0 flex-1 bg-transparent text-[15px] font-semibold text-jv-ink outline-none sm:text-[16px]"
        />
      </label>
    </div>

    <div
      class="grid w-full grid-cols-2 items-stretch gap-1 jv-border-rough bg-jv-white p-1 sm:flex sm:w-fit sm:items-center shadow-brutal-sm"
      role="group"
      aria-label="Report type"
    >
      <button
        v-for="tab in ['quizzes', 'rooms']"
        :key="tab"
        type="button"
        class="cursor-pointer rounded-[6px] px-4 py-2 text-center text-sm font-bold text-jv-ink hover:bg-jv-yellow/40"
        :class="{ 'bg-jv-yellow/60': reportType === tab }"
        :aria-pressed="reportType === tab"
        @click="reportType = tab"
      >
        {{ tab === "quizzes" ? "Quiz reports" : "Room reports" }}
      </button>
    </div>
    <RoomReportList v-if="reportType === 'rooms'" :date="roomDate" />
    <template v-else>
      <!-- State: Pending -->
      <div
        v-if="quizListPending"
        class="jv-border-rough bg-jv-white p-5 text-center text-[18px] font-semibold text-jv-muted shadow-brutal-sm mt-4"
      >
        Loading...
      </div>

      <!-- State: Error -->
      <div
        v-else-if="quizListError"
        class="jv-border-rough bg-jv-white p-5 text-[18px] font-semibold text-jv-coral shadow-brutal-sm mt-4"
      >
        Error while fetching data: {{ quizListError }}
      </div>

      <!-- State: Success -->
      <div v-else class="flex flex-col gap-3">
        <!-- Table Header (Grid) -->
        <div
          class="hidden lg:grid grid-cols-[2fr_1fr_1fr_1fr_1fr_1fr] gap-4 px-5 py-2 text-[12px] font-black uppercase tracking-[0.1em] text-jv-muted"
        >
          <div
            role="button"
            class="flex items-center gap-1 cursor-pointer hover:text-jv-ink transition-colors"
            @click="sortEventHandler('title')"
          >
            QUIZ NAME
            <ChevronUp
              v-if="orderBy === 'title' && order === 'asc'"
              class="size-4"
            />
            <ChevronDown
              v-else-if="orderBy === 'title' && order === 'desc'"
              class="size-4"
            />
            <ChevronsUpDown v-else class="size-4 opacity-40" />
          </div>
          <div
            role="button"
            class="flex items-center justify-center text-center gap-1 cursor-pointer hover:text-jv-ink transition-colors"
            @click="sortEventHandler('participants')"
          >
            TOTAL <br />
            PARTICIPANTS
            <ChevronUp
              v-if="orderBy === 'participants' && order === 'asc'"
              class="size-4"
            />
            <ChevronDown
              v-else-if="orderBy === 'participants' && order === 'desc'"
              class="size-4"
            />
            <ChevronsUpDown v-else class="size-4 opacity-40" />
          </div>
          <div
            role="button"
            class="flex items-center gap-1 cursor-pointer hover:text-jv-ink transition-colors"
            @click="sortEventHandler('activated_from')"
          >
            STARTED AT
            <ChevronUp
              v-if="orderBy === 'activated_from' && order === 'asc'"
              class="size-4"
            />
            <ChevronDown
              v-else-if="orderBy === 'activated_from' && order === 'desc'"
              class="size-4"
            />
            <ChevronsUpDown v-else class="size-4 opacity-40" />
          </div>
          <div
            role="button"
            class="flex items-center gap-1 cursor-pointer hover:text-jv-ink transition-colors"
            @click="sortEventHandler('activated_to')"
          >
            ENDED AT
            <ChevronUp
              v-if="orderBy === 'activated_to' && order === 'asc'"
              class="size-4"
            />
            <ChevronDown
              v-else-if="orderBy === 'activated_to' && order === 'desc'"
              class="size-4"
            />
            <ChevronsUpDown v-else class="size-4 opacity-40" />
          </div>
          <div
            role="button"
            class="flex items-center justify-center text-center gap-1 cursor-pointer hover:text-jv-ink transition-colors"
            @click="sortEventHandler('questions')"
          >
            TOTAL <br />
            QUESTIONS
            <ChevronUp
              v-if="orderBy === 'questions' && order === 'asc'"
              class="size-4"
            />
            <ChevronDown
              v-else-if="orderBy === 'questions' && order === 'desc'"
              class="size-4"
            />
            <ChevronsUpDown v-else class="size-4 opacity-40" />
          </div>
          <div class="flex items-center justify-center text-center">
            ACCURACY
          </div>
        </div>

        <!-- No Data -->
        <div
          v-if="quizList?.data?.count <= 0"
          class="jv-border-rough bg-jv-white p-5 text-center text-[18px] font-semibold text-jv-muted shadow-brutal-sm"
        >
          No quiz found.
        </div>

        <!-- Cards List -->
        <div
          v-for="(quiz, index) in quizList?.data?.data"
          v-else
          :key="index"
          class="group/card relative grid cursor-pointer grid-cols-6 lg:grid-cols-[2fr_1fr_1fr_1fr_1fr_1fr] gap-x-2 gap-y-2 lg:gap-x-4 lg:gap-y-3 items-start lg:items-center jv-border-rough bg-jv-white pl-3 pr-9 py-2.5 lg:py-4 shadow-brutal-sm transition-all hover:-translate-y-0.5 hover:shadow-brutal-md focus-within:shadow-brutal-md sm:pl-5 sm:pr-10"
          @click="viewReport(quiz.id, quiz.title)"
        >
          <ChevronRight
            class="pointer-events-none absolute right-3 top-1/2 size-5 -translate-y-1/2 text-jv-muted transition-colors group-hover/card:text-jv-ink"
            aria-hidden="true"
          />
          <!-- Quiz Name -->
          <div class="contents lg:block lg:min-w-0 lg:pr-4">
            <div class="col-span-full min-w-0">
              <h3
                class="font-headings text-[20px] lg:text-[24px] leading-tight text-jv-ink font-bold truncate"
              >
                <a
                  :href="`/admin/reports/${quiz.id}?title=${encodeURIComponent(decodeURI(quiz.title || ''))}`"
                  class="flex min-w-0 items-center gap-1 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-jv-ink"
                  @click.stop.prevent="viewReport(quiz.id, quiz.title)"
                >
                  <span class="truncate">{{ decodeURI(quiz.title) }}</span>
                </a>
              </h3>
              <p
                v-if="quiz.description.String"
                class="text-[14px] font-semibold text-jv-muted mt-1 truncate"
              >
                {{ quiz.description.String }}
              </p>
            </div>
            <div
              v-if="quiz.sessions > 1"
              class="order-last col-span-full flex items-center lg:mt-2"
            >
              <span class="group relative inline-flex">
                <button
                  type="button"
                  :aria-describedby="`combine-tip-${quiz.id}`"
                  class="inline-flex min-h-8 cursor-pointer whitespace-nowrap items-center gap-1 rounded-full border-2 border-jv-ink bg-jv-yellow px-2.5 py-1 text-xs lg:min-h-9 lg:gap-1.5 lg:px-3 lg:text-[13px] font-bold text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-1"
                  @click.stop="viewCombinedReport(quiz.quiz_id)"
                >
                  <Layers class="size-3.5 lg:size-4" :stroke-width="2.4" />
                  Combine {{ quiz.sessions }} sessions
                </button>
                <span
                  :id="`combine-tip-${quiz.id}`"
                  role="tooltip"
                  class="pointer-events-none invisible absolute left-0 top-full z-20 mt-2 w-64 rounded-[6px] border-2 border-jv-ink bg-jv-ink px-3 py-2 text-[13px] font-semibold leading-snug text-white opacity-0 shadow-brutal-sm transition-opacity group-hover:visible group-hover:opacity-100 group-focus-within:visible group-focus-within:opacity-100"
                >
                  Combine results from multiple sessions of this quiz into one
                  leaderboard.
                </span>
              </span>
            </div>
          </div>

          <!-- Participants -->
          <div
            class="order-1 col-span-2 lg:col-span-1 lg:order-none lg:text-center text-[16px] lg:text-[22px] font-bold lg:font-black text-jv-ink"
          >
            <span
              class="mb-1 block lg:hidden text-[11px] leading-tight text-jv-muted font-bold uppercase"
              >Participants:</span
            >
            {{ quiz.participants }}
          </div>

          <!-- Started At -->
          <div class="hidden lg:block">
            <div class="text-[15px] font-bold leading-snug text-jv-ink">
              <div v-if="quiz.activated_from.Time">
                {{ quiz.activated_from.Time.split(" ")[0] }}
              </div>
              <div
                v-if="quiz.activated_from.Time"
                class="text-[13px] font-semibold text-jv-muted"
              >
                {{ quiz.activated_from.Time.split(" ")[1] }}
              </div>
            </div>
          </div>

          <!-- Ended At -->
          <div class="hidden lg:block">
            <div class="text-[15px] font-bold leading-snug text-jv-ink">
              <div v-if="quiz.activated_to.Time">
                {{ quiz.activated_to.Time.split(" ")[0] }}
              </div>
              <div
                v-if="quiz.activated_to.Time"
                class="text-[13px] font-semibold text-jv-muted"
              >
                {{ quiz.activated_to.Time.split(" ")[1] }}
              </div>
            </div>
          </div>

          <!-- Questions -->
          <div
            class="order-2 col-span-2 lg:col-span-1 lg:order-none lg:text-center text-[16px] lg:text-[22px] font-bold lg:font-black text-jv-ink"
          >
            <span
              class="mb-1 block lg:hidden text-[11px] leading-tight text-jv-muted font-bold uppercase"
              >Questions:</span
            >
            {{ quiz.questions }}
          </div>

          <!-- Accuracy -->
          <div
            class="order-3 col-span-2 lg:col-span-1 lg:order-none lg:flex lg:justify-center"
          >
            <span
              class="mb-1 block lg:hidden text-[11px] leading-tight text-jv-muted font-bold uppercase"
              >Accuracy:</span
            >
            <div
              class="inline-flex items-center justify-center rounded-md px-1.5 py-0.5 text-[14px] font-black leading-5 lg:jv-border-rough lg:min-w-[90px] lg:px-3 lg:py-1.5 lg:text-[16px] lg:leading-normal lg:shadow-brutal-sm"
              :class="
                getAccuracyClass(
                  (quiz.correct_answers /
                    (quiz.participants * quiz.questions)) *
                    100
                )
              "
            >
              {{
                (
                  (quiz.correct_answers /
                    (quiz.participants * quiz.questions)) *
                  100
                ).toFixed(2)
              }}%
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <Pagination :page="page" :num-of-records="quizList?.data?.count / 10" />
    </template>
  </main>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import { format, formatISO } from "date-fns";
import { useDebounceFn } from "@vueuse/core";
import {
  Search,
  Calendar,
  ChevronUp,
  ChevronDown,
  ChevronsUpDown,
  Layers,
  ChevronRight,
} from "lucide-vue-next";
import Pagination from "~/components/Pagination.vue";

definePageMeta({
  layout: "empty",
});

useSeoMeta({
  title: "Reports - jovVix",
  description:
    "View performance reports and analytics for your jovVix quiz sessions.",
  robots: "noindex, nofollow",
});

const reportType = ref(useRoute().query.type === "rooms" ? "rooms" : "quizzes");
watch(reportType, (type) => {
  const current = useRoute();
  navigateTo(
    { path: current.path, query: { ...current.query, type } },
    { replace: true }
  );
});
const roomDate = ref("");
const { apiUrl } = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const route = useRoute();
const page = computed(() => Number(route.query.page) || 1);
const order = computed(() => route.query.order || "desc");
const orderBy = computed(() => route.query.orderBy || "activated_to");
const date = ref(route.query.date);
const name = computed(() => route.query.name || "");
const nameInput = ref(route.query.name);

const {
  data: quizList,
  error: quizListError,
  pending: quizListPending,
} = useFetch(`${apiUrl}/admin/reports/list`, {
  transform: (quizList) => {
    quizList.data.data = quizList.data?.data?.map((quiz) => {
      quiz.activated_from.Time = format(
        formatISO(quiz.activated_from.Time),
        "dd-MMM-yyyy HH:mm:ss"
      );
      quiz.activated_to.Time = format(
        formatISO(quiz.activated_to.Time),
        "dd-MMM-yyyy HH:mm:ss"
      );
      return quiz;
    });
    return quizList;
  },
  credentials: "include",
  headers: headers,
  query: {
    page,
    orderBy,
    order,
    name,
    date,
  },
});
useSessionExpiredRedirect(quizListError);

// The title names the downloaded PDF; the report API does not return it.
const viewReport = (activeQuizId, title) => {
  navigateTo({
    path: `/admin/reports/${activeQuizId}`,
    query: { title: decodeURI(title || "") },
  });
};

const viewCombinedReport = (quizId) => {
  navigateTo(`/admin/reports/combined/${quizId}`);
};

const sortEventHandler = (columnName) => {
  let ordercol = order.value;
  if (columnName === orderBy.value) {
    ordercol = ordercol === "asc" ? "desc" : "asc";
  }
  navigateTo({
    path: route.path,
    query: {
      ...route.query,
      orderBy: columnName,
      order: ordercol,
    },
  });
};

const debouncedNavigateTo = useDebounceFn((query) => {
  navigateTo({
    path: route.path,
    query: query,
  });
}, 500);

watch(nameInput, (newName) => {
  debouncedNavigateTo({
    ...route.query,
    name: newName,
  });
});

const getAccuracyClass = (accuracy) => {
  if (isNaN(accuracy)) return "bg-jv-white text-jv-ink";
  if (accuracy >= 60) return "bg-[#d1f4e0] text-[#1e5631]";
  if (accuracy >= 40) return "bg-jv-yellow text-jv-ink";
  return "bg-jv-coral text-white";
};
</script>
