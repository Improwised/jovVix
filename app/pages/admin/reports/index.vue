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
          Quiz Reports
        </h1>
        <div
          class="ml-1 mt-1 h-3 w-40 rounded-full border-b-[3px] border-jv-yellow sm:ml-2 sm:w-48"
          aria-hidden="true"
        ></div>
      </div>

      <div class="flex flex-col sm:flex-row gap-4 items-center">
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
    </div>

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
        <div class="flex items-center justify-center text-center">ACCURACY</div>
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
        role="button"
        class="grid grid-cols-1 lg:grid-cols-[2fr_1fr_1fr_1fr_1fr_1fr] gap-4 items-center jv-border-rough bg-jv-white px-5 py-4 shadow-brutal-sm hover:translate-y-[-2px] hover:shadow-brutal-md transition-all cursor-pointer"
        @click="viewReport(quiz.id)"
      >
        <!-- Quiz Name -->
        <div class="min-w-0 pr-4">
          <h3
            class="font-headings text-[22px] sm:text-[24px] leading-tight text-jv-ink font-bold truncate"
          >
            {{ decodeURI(quiz.title) }}
          </h3>
          <p class="text-[14px] font-semibold text-jv-muted mt-1 truncate">
            {{ quiz.description.String }}
          </p>
        </div>

        <!-- Participants -->
        <div class="lg:text-center text-[22px] font-black text-jv-ink">
          <span
            class="lg:hidden text-[14px] text-jv-muted font-bold mr-2 uppercase tracking-wide"
            >Participants:</span
          >
          {{ quiz.participants }}
        </div>

        <!-- Started At -->
        <div>
          <span
            class="lg:hidden text-[14px] text-jv-muted font-bold mr-2 uppercase tracking-wide"
            >Started At:</span
          >
          <div
            class="inline-block lg:block text-[15px] font-bold text-jv-ink leading-snug"
          >
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
        <div>
          <span
            class="lg:hidden text-[14px] text-jv-muted font-bold mr-2 uppercase tracking-wide"
            >Ended At:</span
          >
          <div
            class="inline-block lg:block text-[15px] font-bold text-jv-ink leading-snug"
          >
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
        <div class="lg:text-center text-[22px] font-black text-jv-ink">
          <span
            class="lg:hidden text-[14px] text-jv-muted font-bold mr-2 uppercase tracking-wide"
            >Questions:</span
          >
          {{ quiz.questions }}
        </div>

        <!-- Accuracy -->
        <div class="lg:flex lg:justify-center">
          <span
            class="lg:hidden text-[14px] text-jv-muted font-bold mr-2 uppercase tracking-wide"
            >Accuracy:</span
          >
          <div
            class="inline-flex items-center justify-center px-3 py-1.5 jv-border-rough shadow-brutal-sm text-[16px] font-black min-w-[90px]"
            :class="
              getAccuracyClass(
                (quiz.correct_answers / (quiz.participants * quiz.questions)) *
                  100
              )
            "
          >
            {{
              (
                (quiz.correct_answers / (quiz.participants * quiz.questions)) *
                100
              ).toFixed(2)
            }}%
          </div>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <Pagination :page="page" :num-of-records="quizList?.data?.count / 10" />
  </main>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import { format, formatISO } from "date-fns";
import debounce from "lodash/debounce";
import {
  Search,
  Calendar,
  ChevronUp,
  ChevronDown,
  ChevronsUpDown,
} from "lucide-vue-next";
import Pagination from "~/components/Pagination.vue";

definePageMeta({
  layout: "empty",
});

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

const viewReport = (activeQuizId) => {
  navigateTo(`/admin/reports/${activeQuizId}`);
};

const sortEventHandler = (columnName) => {
  let ordercol = order.value;
  if (columnName === orderBy.value) {
    ordercol === "asc" ? (ordercol = "desc") : (ordercol = "asc");
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

const debouncedNavigateTo = debounce((query) => {
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
