<template>
  <main
    class="flex min-h-screen flex-col gap-8 bg-jv-canvas px-4 py-5 sm:gap-10 sm:px-6 md:px-8 md:py-6"
  >
    <!-- Header -->
    <div class="min-w-0">
      <h1
        class="font-headings text-[38px] leading-none text-jv-ink min-[420px]:text-[44px] sm:text-[52px] md:text-[56px]"
      >
        Quiz Analysis
      </h1>
      <div
        class="ml-1 mt-1 h-3 w-40 rounded-full border-b-[3px] border-jv-yellow sm:ml-2 sm:w-48"
        aria-hidden="true"
      ></div>
      <p
        class="mt-4 max-w-2xl text-[15px] font-semibold text-jv-muted sm:text-[16px]"
      >
        A complete performance breakdown for this quiz review per-question
        accuracy and response times, then dive into individual participant
        results, all in one place.
      </p>
    </div>

    <!-- Tabs -->
    <PageLayout :current-tab="currentTab" @change-tab="changeTab" />

    <!-- Class Performance Overview (visible on both tabs once data is loaded) -->
    <section
      v-if="quizAnalysis?.data?.length"
      class="jv-border-rough bg-jv-ink p-5 text-white shadow-brutal sm:p-7 md:p-8"
    >
      <div class="flex flex-col gap-6 lg:flex-row lg:items-center">
        <div class="flex-1 min-w-0">
          <div
            class="flex items-center gap-2 text-[12px] font-black uppercase tracking-[0.16em] text-jv-yellow"
          >
            <Zap class="size-4" :stroke-width="2.6" />
            Class Performance Overview
          </div>
          <h2
            class="mt-3 font-headings text-[28px] leading-tight text-white sm:text-[32px] md:text-[36px]"
          >
            {{ overallMessage }}
          </h2>
          <p
            class="mt-2 max-w-md text-[14px] font-semibold text-white/70 sm:text-[15px]"
          >
            Review question-level metrics and overall accuracy below. Use the
            numbered navigator to jump directly to specific items.
          </p>
        </div>

        <div
          class="jv-border-rough bg-jv-white p-4 text-jv-ink shadow-brutal-sm sm:p-5 lg:w-[340px]"
        >
          <div class="flex items-baseline justify-between">
            <span class="text-[15px] font-bold sm:text-[16px]"
              >Class accuracy</span
            >
            <span class="text-[22px] font-black sm:text-[24px]"
              >{{ classAccuracy.toFixed(0) }}%</span
            >
          </div>
          <div
            class="mt-2 h-2.5 w-full overflow-hidden rounded-full bg-jv-slate"
          >
            <div
              class="h-full rounded-full bg-jv-accent-green"
              :style="{ width: `${Math.min(classAccuracy, 100)}%` }"
            ></div>
          </div>

          <div class="mt-4 grid grid-cols-3 gap-3">
            <div class="min-w-0">
              <div
                class="text-[10px] font-black uppercase tracking-[0.12em] text-jv-muted"
              >
                Questions reviewed
              </div>
              <div class="mt-1 text-[20px] font-black text-jv-ink">
                {{ totalQuestions }}
              </div>
            </div>
            <div class="min-w-0">
              <div
                class="text-[10px] font-black uppercase tracking-[0.12em] text-jv-muted"
              >
                Participants
              </div>
              <div class="mt-1 text-[20px] font-black text-jv-ink">
                {{ totalParticipants }}
              </div>
            </div>
            <div class="min-w-0">
              <div
                class="text-[10px] font-black uppercase tracking-[0.12em] text-jv-muted"
              >
                Avg. completion
              </div>
              <div class="mt-1 text-[20px] font-black text-jv-ink">
                {{ avgCompletion.toFixed(2) }} sec
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Questions tab -->
    <template v-if="currentTab === 'report'">
      <div
        v-if="quizAnalysisPending"
        class="jv-border-rough bg-jv-white p-5 text-center text-[18px] font-semibold text-jv-muted shadow-brutal-sm"
      >
        Loading...
      </div>

      <div
        v-else-if="quizAnalysisError"
        class="jv-border-rough bg-jv-white p-5 text-[18px] font-semibold text-jv-coral shadow-brutal-sm"
      >
        Error while fetching data: {{ quizAnalysisError }}
      </div>

      <template v-else-if="quizAnalysis?.data?.length">
        <!-- Question pager strip -->
        <div class="jv-border-rough bg-jv-white p-4 shadow-brutal-sm sm:p-5">
          <div class="flex flex-wrap items-center gap-3">
            <div
              class="flex items-center gap-2 text-[15px] font-bold text-jv-ink sm:text-[16px]"
            >
              <ListOrdered class="size-5 text-jv-muted" :stroke-width="2.4" />
              {{ totalQuestions }} Questions
            </div>
            <div class="flex flex-wrap items-center gap-2">
              <button
                v-for="n in totalQuestions"
                :key="n"
                type="button"
                class="flex h-8 w-8 items-center justify-center jv-border-rough text-[14px] font-black shadow-brutal-sm transition-transform hover:translate-y-[-1px]"
                :class="
                  activeQuestionIndex === n
                    ? 'bg-jv-coral text-white'
                    : 'bg-jv-white text-jv-ink hover:bg-jv-yellow/50'
                "
                @click="scrollToQuestion(n)"
              >
                {{ n }}
              </button>
            </div>
          </div>
        </div>

        <!-- Questions list -->
        <div class="flex flex-col gap-5">
          <article
            v-for="(quiz, index) in quizAnalysis.data"
            :id="`question-${index + 1}`"
            :key="index"
            role="button"
            tabindex="0"
            class="jv-border-rough bg-jv-white p-4 shadow-brutal-sm sm:p-5 scroll-mt-24 cursor-pointer transition-transform hover:rotate-[0.25deg] focus:outline-none focus:shadow-brutal"
            @click="openQuestionModal(quiz, index)"
            @keydown.enter.prevent="openQuestionModal(quiz, index)"
            @keydown.space.prevent="openQuestionModal(quiz, index)"
          >
            <!-- Question header row -->
            <div
              class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between"
            >
              <div class="min-w-0">
                <p
                  class="text-[12px] font-bold uppercase tracking-[0.14em] text-jv-coral"
                >
                  Question {{ index + 1 }}
                </p>
                <h3
                  class="mt-1 break-words text-[20px] font-bold leading-snug text-jv-ink sm:text-[22px]"
                >
                  {{ quiz.question }}
                </h3>
              </div>

              <div
                class="flex shrink-0 flex-wrap items-start gap-4 text-right sm:gap-5"
              >
                <div>
                  <div
                    class="text-[10px] font-bold uppercase tracking-[0.12em] text-jv-muted"
                  >
                    Avg. time
                  </div>
                  <div class="mt-0.5 text-[14px] font-bold text-jv-ink">
                    {{ Math.abs((quiz.avg_response_time / 1000).toFixed(2)) }}
                    sec
                  </div>
                </div>
                <div>
                  <div
                    class="text-[10px] font-bold uppercase tracking-[0.12em] text-jv-muted"
                  >
                    Correct
                  </div>
                  <div
                    class="mt-0.5 text-[14px] font-bold text-jv-accent-green"
                  >
                    {{ correctCountFor(quiz) }}
                  </div>
                </div>
                <div>
                  <div
                    class="text-[10px] font-bold uppercase tracking-[0.12em] text-jv-muted"
                  >
                    Incorrect
                  </div>
                  <div class="mt-0.5 text-[14px] font-bold text-jv-coral">
                    {{ incorrectCountFor(quiz) }}
                  </div>
                </div>
              </div>
            </div>

            <div
              v-if="quiz.question_media === 'image' && quiz.resource"
              class="mt-3 flex justify-center rounded-md bg-jv-canvas p-2"
              @click.stop
            >
              <img
                :src="quiz.resource"
                :alt="quiz.question"
                class="max-h-64 w-auto max-w-full object-contain"
              />
            </div>
            <div
              v-else-if="quiz.question_media === 'code' && quiz.resource"
              class="mt-3 min-w-0 overflow-x-auto"
              @click.stop
            >
              <CodeBlockComponent :code="quiz.resource" />
            </div>

            <!-- Options -->
            <ul class="mt-4 flex flex-col">
              <li
                v-for="(option, order) in quiz.options"
                :key="order"
                class="flex min-w-0 items-center gap-3 border-b border-jv-ink/10 py-3 pl-3 pr-2 text-[15px] font-medium text-jv-ink last:border-b-0"
                :class="
                  quiz.correct_answer.includes(Number(order))
                    ? 'border-l-4 border-l-jv-accent-green bg-jv-accent-green/25 pl-2'
                    : 'border-l-4 border-l-transparent'
                "
              >
                <span class="w-5 shrink-0 text-[14px] font-bold text-jv-coral">
                  {{ String.fromCharCode(64 + Number(order)) }}.
                </span>

                <div
                  v-if="quiz.options_media === 'image' && option"
                  class="flex min-w-0 flex-1 justify-start"
                  @click.stop
                >
                  <img
                    :src="option"
                    :alt="option"
                    class="max-h-32 w-auto max-w-full object-contain"
                  />
                </div>

                <div
                  v-else-if="quiz.options_media === 'code' && option"
                  class="min-w-0 flex-1 overflow-x-auto"
                  @click.stop
                >
                  <CodeBlockComponent :code="option" />
                </div>

                <span v-else class="min-w-0 flex-1 break-words">
                  {{ option }}
                </span>

                <span
                  class="inline-flex shrink-0 items-center gap-1 text-[13px] font-semibold text-jv-muted"
                >
                  <Users class="size-4" :stroke-width="2.4" />
                  {{ quiz.selected_answers[order]?.length || 0 }}
                </span>
              </li>
            </ul>
          </article>
        </div>
      </template>

      <div
        v-else
        class="jv-border-rough bg-jv-white p-5 text-center text-[18px] font-semibold text-jv-muted shadow-brutal-sm"
      >
        No questions found for this quiz.
      </div>
    </template>

    <!-- Participants tab -->
    <template v-if="currentTab === 'participants'">
      <div
        v-if="quizUserAnalysispending"
        class="jv-border-rough bg-jv-white p-5 text-center text-[18px] font-semibold text-jv-muted shadow-brutal-sm"
      >
        Loading...
      </div>
      <div
        v-else-if="fetchError"
        class="jv-border-rough bg-jv-white p-5 text-[18px] font-semibold text-jv-coral shadow-brutal-sm"
      >
        Error while fetching data: {{ fetchError }}
      </div>
      <div v-else class="flex flex-col gap-4">
        <QuizUserAnalyticsSpace
          v-for="(oData, index) in rankData"
          :key="index"
          :data="userJson[oData]"
          :user-name="oData"
          :survey-questions="surveyQuestions"
        />
      </div>
    </template>

    <QuizQuestionDetailModal
      :is-open="!!selectedQuestion"
      :question="selectedQuestion"
      :index="selectedQuestionIndex"
      :rows="selectedQuestionRows"
      @close="closeQuestionModal"
    />
  </main>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { Zap, ListOrdered, Users } from "lucide-vue-next";
import lodash from "lodash";
import PageLayout from "~~/components/reports/PageLayout.vue";
import CodeBlockComponent from "~~/components/CodeBlockComponent.vue";

definePageMeta({
  layout: "empty",
});

useSeoMeta({
  title: "Report Details - jovVix",
  description:
    "Explore detailed participant scores and question analytics for a jovVix quiz.",
  robots: "noindex, nofollow",
});

const { apiUrl } = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const route = useRoute();
const activeQuizId = computed(() => route.params.id);
const currentTab = ref("report");

const {
  data: quizAnalysis,
  error: quizAnalysisError,
  pending: quizAnalysisPending,
} = useFetch(`${apiUrl}/admin/reports/${activeQuizId.value}/analysis`, {
  transform: (quizAnalysis) => {
    quizAnalysis.data?.map((quiz) => {
      quiz.userParticipants = Object.keys(quiz.selected_answers).length;
      const result = {};

      for (const [user, answer] of Object.entries(quiz.selected_answers)) {
        if (!result[answer]) {
          result[answer] = [];
        }
        result[answer].push(user);
      }

      quiz.selected_answers = result;
      quiz.correctPercentage =
        (quiz.correct_answer.reduce(
          (sum, correct_answer) =>
            (sum += quiz.selected_answers[correct_answer]?.length || 0),
          0
        ) /
          quiz.userParticipants) *
        100;

      return quiz;
    });
    return quizAnalysis;
  },
  credentials: "include",
  headers: headers,
});

// Derived overview metrics
const totalQuestions = computed(() => quizAnalysis.value?.data?.length || 0);

const totalParticipants = computed(() => {
  const data = quizAnalysis.value?.data;
  if (!data?.length) return 0;
  return data.reduce((m, q) => Math.max(m, q.userParticipants || 0), 0);
});

const classAccuracy = computed(() => {
  const data = quizAnalysis.value?.data;
  if (!data?.length) return 0;
  const valid = data.filter((q) => Number.isFinite(q.correctPercentage));
  if (!valid.length) return 0;
  return valid.reduce((sum, q) => sum + q.correctPercentage, 0) / valid.length;
});

const avgCompletion = computed(() => {
  const data = quizAnalysis.value?.data;
  if (!data?.length) return 0;
  const mean =
    data.reduce((sum, q) => sum + (q.avg_response_time || 0), 0) / data.length;
  return mean / 1000;
});

const overallMessage = computed(() => {
  const acc = classAccuracy.value;
  if (acc >= 70) return "Great overall performance";
  if (acc >= 50) return "Solid performance";
  return "Needs improvement";
});

const correctCountFor = (quiz) =>
  quiz.correct_answer.reduce(
    (sum, idx) => sum + (quiz.selected_answers[idx]?.length || 0),
    0
  );

const incorrectCountFor = (quiz) =>
  Math.max((quiz.userParticipants || 0) - correctCountFor(quiz), 0);

// Question pager: scroll behavior + active tracking
const activeQuestionIndex = ref(1);

const scrollToQuestion = (n) => {
  activeQuestionIndex.value = n;
  if (typeof document === "undefined") return;
  const el = document.getElementById(`question-${n}`);
  el?.scrollIntoView({ behavior: "smooth", block: "start" });
};

let observer = null;

onMounted(() => {
  // Wait for questions to render, then observe them for active-chip tracking.
  watchEffect(() => {
    if (!quizAnalysis.value?.data?.length) return;
    nextTick(() => {
      observer?.disconnect();
      observer = new IntersectionObserver(
        (entries) => {
          const visible = entries
            .filter((e) => e.isIntersecting)
            .sort(
              (a, b) =>
                a.target.getBoundingClientRect().top -
                b.target.getBoundingClientRect().top
            )[0];
          if (visible) {
            const id = visible.target.id;
            const n = Number(id.replace("question-", ""));
            if (n) activeQuestionIndex.value = n;
          }
        },
        { rootMargin: "-30% 0px -60% 0px", threshold: 0 }
      );
      for (let i = 1; i <= quizAnalysis.value.data.length; i++) {
        const el = document.getElementById(`question-${i}`);
        if (el) observer.observe(el);
      }
    });
  });
});

onBeforeUnmount(() => {
  observer?.disconnect();
});

// Participants tab data (unchanged behavior)
const analysisJson = ref([]);
const userJson = ref({});
const questionJson = ref({});
const rankData = ref([]);
const surveyQuestions = ref(0);
const ranks = ref();
const quizUserAnalysispending = ref(false);
const fetchError = ref("");

const getAnalysisJson = async () => {
  try {
    quizUserAnalysispending.value = true;
    const response = await fetch(
      `${apiUrl}/analytics_board/admin?active_quiz_id=${activeQuizId.value}`,
      {
        method: "GET",
        headers: headers.value,
        mode: "cors",
        credentials: "include",
      }
    );

    const ranksResponse = await fetch(
      `${apiUrl}/final_score/admin?active_quiz_id=${activeQuizId.value}`,
      {
        method: "GET",
        headers: headers.value,
        mode: "cors",
        credentials: "include",
      }
    );

    const result = await response.json();
    ranks.value = await ranksResponse.json();

    if (response.ok && ranksResponse.ok) {
      analysisJson.value = result.data;

      userJson.value = lodash.groupBy(analysisJson.value, "username");

      ranks.value?.data.forEach((data) => {
        rankData.value.push(data.username);
        let key = data.username;

        if (userJson.value.hasOwnProperty(key)) {
          let totalScore = data.score;

          userJson.value[key].push({
            rank: data.rank,
            total_score: totalScore,
            response_time: data.response_time,
            avatar: data.img_key,
          });
        } else {
          console.error(`Key '${key}' not found in userJson.value.`);
        }
      });

      questionJson.value = lodash.groupBy(analysisJson.value, "question");

      for (const key in userJson.value) {
        userJson.value[key].forEach((question) => {
          if (!question.rank) {
            if (question.question_type == "survey") {
              surveyQuestions.value++;
            }
          }
        });
        break;
      }
    } else {
      console.error(result);
    }

    quizUserAnalysispending.value = false;
  } catch (error) {
    quizUserAnalysispending.value = false;
    fetchError.value = error;
    console.error("Failed to fetch data", error);
  }
};

onMounted(() => {
  getAnalysisJson();
});

const changeTab = (data) => {
  currentTab.value = data;
};

// Per-question detail modal
const selectedQuestion = ref(null);
const selectedQuestionIndex = ref(0);

const rowsByQuestion = computed(() =>
  lodash.groupBy(analysisJson.value || [], "question")
);

const selectedQuestionRows = computed(() => {
  if (!selectedQuestion.value) return [];
  return rowsByQuestion.value[selectedQuestion.value.question] || [];
});

const openQuestionModal = (quiz, index) => {
  selectedQuestion.value = quiz;
  selectedQuestionIndex.value = index + 1;
};

const closeQuestionModal = () => {
  selectedQuestion.value = null;
};
</script>
