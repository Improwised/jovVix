<script setup>
import {
  Check,
  Flame,
  SkipForward,
  Trophy,
  Volume2,
  VolumeX,
  X,
} from "lucide-vue-next";
import AnswerSubmissionChart from "../AnswerSubmissionChart.vue";
import { useMusicStore } from "~~/store/music";
import { getAvatarUrlByName } from "~~/composables/avatar";

const app = useNuxtApp();
const musicStore = useMusicStore();
const { getMusic, setMusic } = musicStore;
const music = computed(() => getMusic());
const toggleMusic = () => setMusic(!music.value);

const props = defineProps({
  data: {
    default: () => {
      return {};
    },
    type: Object,
    required: true,
  },
  isAdmin: {
    default: false,
    type: Boolean,
    required: false,
  },
  userName: {
    required: false,
    type: String,
    default: "",
  },
  selectedAnswer: {
    required: false,
    type: Number,
    default: 0,
  },
  analysisTab: {
    type: String,
    default: "",
  },
  quizState: {
    type: String,
    required: false,
    default: "",
  },
});

const emits = defineEmits(["askSkipTimer", "changeAnalysisTab"]);
const timer = ref(null);
const time = ref(0);
const isSkip = ref(false);

const question = computed(() => props?.data?.data ?? {});
const options = computed(() => question.value?.options ?? {});
const rankList = computed(() => question.value?.rankList ?? []);
const duration = computed(() => Number(question.value?.duration) || 0);

const progressPercent = computed(() => {
  if (!duration.value) return 0;
  if (time.value < 0) return 100;
  return Math.min(100, Math.max(0, (time.value * 100) / duration.value));
});

const isLastQuestion = computed(() => {
  const q = question.value;
  if (!q) return false;
  return Number(q.question_no) === Number(q.totalQuestions);
});

const avatarBgs = [
  "bg-jv-yellow",
  "bg-jv-salmon",
  "bg-jv-mint",
  "bg-jv-ivory",
  "bg-jv-lavender",
];
const avatarBgFor = (index) => avatarBgs[index % avatarBgs.length];

function handleTimer() {
  clearInterval(timer.value);
  timer.value = setInterval(() => {
    const isPauseQuiz = props.quizState === app.$Pause;
    if (!isPauseQuiz) {
      time.value += 0.1;
      if (time.value == duration.value + 1) {
        clearInterval(timer.value);
        time.value = -1;
        timer.value = null;
      }
    }
  }, 100);
}

handleTimer();

function handleSkipTimer(e) {
  e.preventDefault();
  isSkip.value = true;
  emits("askSkipTimer");
}

const changeAnalysisTab = (tab) => emits("changeAnalysisTab", tab);

onUnmounted(() => {
  if (timer.value) {
    clearInterval(timer.value);
  }
});
</script>

<template>
  <main
    class="flex min-h-screen flex-col bg-jv-canvas px-3 py-4 text-jv-ink sm:px-6 sm:py-6 md:px-10 md:py-8"
  >
    <section class="mx-auto w-full max-w-[1180px]">
      <article
        class="relative -rotate-[0.4deg] jv-border-rough bg-jv-white p-5 shadow-brutal sm:p-7 md:p-9"
      >
        <span
          class="absolute left-1/2 top-[-12px] z-10 h-4 w-16 -translate-x-1/2 rotate-[1deg] bg-jv-coral"
          aria-hidden="true"
        ></span>

        <!-- Header: title + music toggle -->
        <header class="flex items-start justify-between gap-4">
          <div class="min-w-0 flex-1">
            <h2
              class="font-headings text-[24px] leading-none text-jv-ink min-[420px]:text-[30px] sm:text-[36px] md:text-[40px]"
            >
              Score Page
            </h2>
            <p
              class="mt-2 font-body text-[12px] font-bold text-jv-muted sm:text-[14px]"
            >
              Rank board
            </p>
          </div>

          <button
            type="button"
            class="grid size-10 shrink-0 place-items-center rounded-[8px] border-[2px] border-jv-ink bg-jv-white text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[2deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:size-11"
            :aria-label="music ? 'Mute music' : 'Unmute music'"
            @click="toggleMusic"
          >
            <Volume2 v-if="music" class="size-4" :stroke-width="2.4" />
            <VolumeX v-else class="size-4" :stroke-width="2.4" />
          </button>
        </header>

        <!-- Coral progress bar -->
        <div
          class="mt-5 h-3 w-full overflow-hidden rounded-full border-[2px] border-jv-ink bg-jv-white sm:mt-6"
          aria-hidden="true"
        >
          <div
            class="h-full bg-jv-coral transition-[width] duration-150"
            :style="{ width: progressPercent + '%' }"
          ></div>
        </div>

        <!-- Question text -->
        <h3
          class="mt-6 font-headings text-[20px] leading-snug text-jv-ink min-[420px]:text-[24px] sm:mt-7 sm:text-[30px] md:text-[36px]"
        >
          {{ question.question }}
        </h3>

        <div
          v-if="question.question_media === 'image'"
          class="mt-5 flex justify-center"
        >
          <img
            :src="question.resource"
            :alt="question.question"
            class="max-h-[260px] w-auto border-[3px] border-jv-ink bg-jv-white object-contain shadow-brutal-sm"
          />
        </div>
        <div v-else-if="question.question_media === 'code'" class="mt-5">
          <CodeBlockComponent :code="question.resource" />
        </div>

        <!-- Options grid: correct answer highlighted -->
        <div class="mt-6 grid gap-3 sm:mt-7 sm:gap-4 md:grid-cols-2 md:gap-5">
          <div
            v-for="(answer, key) in options"
            :key="key"
            :class="[
              'relative flex w-full items-center gap-3 border-[2px] border-jv-ink px-3 py-3 shadow-brutal-sm sm:gap-4 sm:px-4 sm:py-4',
              answer.isAnswer
                ? 'bg-jv-mint'
                : Number(key) === Number(selectedAnswer)
                ? 'bg-jv-salmon/40'
                : 'bg-jv-white',
              selectedAnswer && Number(key) === Number(selectedAnswer)
                ? 'outline-[3px] outline-offset-[2px] outline-jv-ink'
                : '',
              selectedAnswer &&
              !answer.isAnswer &&
              Number(key) !== Number(selectedAnswer)
                ? 'opacity-60'
                : '',
            ]"
          >
            <span
              v-if="selectedAnswer && Number(key) === Number(selectedAnswer)"
              class="absolute -top-3 right-3 inline-flex items-center gap-1 rounded-full border-[2px] border-jv-ink bg-jv-yellow px-2 py-0.5 font-body text-[10px] font-black uppercase tracking-wider text-jv-ink shadow-brutal-sm sm:text-[11px]"
            >
              Your pick
            </span>
            <span
              v-if="answer.isAnswer"
              class="grid size-10 shrink-0 place-items-center rounded-[6px] border-[2px] border-jv-ink bg-jv-white text-jv-ink sm:size-11"
              aria-label="Correct answer"
            >
              <Check class="size-5" :stroke-width="3" />
            </span>
            <span
              v-else-if="
                selectedAnswer && Number(key) === Number(selectedAnswer)
              "
              class="grid size-10 shrink-0 place-items-center rounded-[6px] border-[2px] border-jv-ink bg-jv-white text-jv-coral sm:size-11"
              aria-label="Your incorrect answer"
            >
              <X class="size-5" :stroke-width="3" />
            </span>
            <span
              v-else
              class="grid size-10 shrink-0 place-items-center rounded-[6px] border-[2px] border-jv-ink bg-jv-white font-feature text-[14px] font-black text-jv-ink sm:size-11 sm:text-[16px]"
            >
              {{ String.fromCharCode(64 + Number(key)) }}
            </span>
            <div class="min-w-0 flex-1">
              <div
                v-if="question.options_media === 'image'"
                class="flex items-center justify-center"
              >
                <img
                  :src="answer?.value"
                  :alt="`Option ${key}`"
                  class="max-h-[120px] w-auto object-contain"
                />
              </div>
              <CodeBlockComponent
                v-else-if="question.options_media === 'code'"
                :code="answer?.value"
              />
              <span
                v-else
                class="block break-words font-body text-[15px] font-black leading-snug text-jv-ink sm:text-[18px] md:text-[20px]"
              >
                {{ answer?.value }}
              </span>
            </div>
          </div>
        </div>

        <!-- Admin: skip / finish -->
        <template v-if="isAdmin">
          <div
            class="mt-6 border-t-[2px] border-dashed border-jv-ink/30 sm:mt-7"
          ></div>
          <div class="mt-4 flex justify-end">
            <button
              type="button"
              class="inline-flex h-10 items-center justify-center gap-2 rounded-full border-[2px] border-jv-ink bg-jv-white px-5 font-body text-[14px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none disabled:cursor-not-allowed disabled:opacity-50 sm:h-11 sm:px-6 sm:text-[15px]"
              :disabled="isSkip"
              @click="handleSkipTimer"
            >
              <Trophy
                v-if="isLastQuestion"
                class="size-4"
                :stroke-width="2.4"
              />
              <SkipForward v-else class="size-4" :stroke-width="2.4" />
              <span>{{ isLastQuestion ? "Finish" : "Skip" }}</span>
            </button>
          </div>
        </template>

        <!-- Tabs (admin only) -->
        <div
          v-if="isAdmin"
          class="mt-6 border-b-[2px] border-jv-ink/15 sm:mt-7"
        >
          <div class="flex">
            <button
              type="button"
              class="relative px-5 pb-3 font-headings text-[16px] transition-colors sm:px-7 sm:text-[20px]"
              :class="
                analysisTab === 'ranking'
                  ? 'text-jv-ink'
                  : 'text-jv-muted hover:text-jv-ink'
              "
              @click="changeAnalysisTab('ranking')"
            >
              Rankings
              <span
                v-if="analysisTab === 'ranking'"
                class="absolute -bottom-[2px] left-2 right-2 h-[3px] bg-jv-coral"
                aria-hidden="true"
              ></span>
            </button>
            <button
              type="button"
              class="relative px-5 pb-3 font-headings text-[16px] transition-colors sm:px-7 sm:text-[20px]"
              :class="
                analysisTab === 'chart'
                  ? 'text-jv-ink'
                  : 'text-jv-muted hover:text-jv-ink'
              "
              @click="changeAnalysisTab('chart')"
            >
              Chart
              <span
                v-if="analysisTab === 'chart'"
                class="absolute -bottom-[2px] left-2 right-2 h-[3px] bg-jv-coral"
                aria-hidden="true"
              ></span>
            </button>
          </div>
        </div>

        <!-- Rankings list -->
        <div
          v-if="!isAdmin || analysisTab === 'ranking'"
          class="mt-5 flex flex-col gap-3 sm:mt-6 sm:gap-3.5"
        >
          <div
            v-for="(user, index) in rankList"
            :key="user.username || index"
            :class="[
              'flex items-center gap-3 border-[2px] border-jv-ink px-3 py-2.5 shadow-brutal-sm sm:gap-4 sm:px-4 sm:py-3',
              user.username === userName ? 'bg-jv-yellow-soft' : 'bg-jv-white',
            ]"
          >
            <span
              class="min-w-[24px] text-center font-feature text-[18px] font-black text-jv-ink sm:min-w-[32px] sm:text-[22px]"
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
                class="truncate font-body text-[15px] font-black leading-tight text-jv-ink sm:text-[17px]"
              >
                {{ user.firstname }}
              </p>
              <p
                v-if="user.username === userName && user.username"
                class="truncate font-body text-[11px] font-bold text-jv-muted sm:text-[12px]"
              >
                @{{ user.username }}
              </p>
            </div>
            <span
              class="font-feature text-[16px] font-black text-jv-ink sm:text-[20px]"
            >
              {{ user.score }}
            </span>
            <span
              class="ml-1 inline-flex min-w-[36px] items-center justify-end gap-1 font-feature text-[14px] font-black text-jv-coral sm:ml-2 sm:text-[16px]"
            >
              <Flame class="size-4 sm:size-[18px]" :stroke-width="2.4" />
              {{ user.streak_count }}
            </span>
          </div>

          <div
            v-if="rankList.length === 0"
            class="mx-auto mt-2 flex w-fit items-center gap-2 rounded-full border-[2px] border-dashed border-jv-ink/30 px-4 py-2 font-body text-[13px] font-bold text-jv-muted sm:text-[14px]"
          >
            <span
              class="size-2 rounded-full bg-jv-ink/30"
              aria-hidden="true"
            ></span>
            No rankings yet
          </div>
        </div>

        <!-- Chart -->
        <div
          v-if="isAdmin && analysisTab === 'chart'"
          class="mt-5 border-[2px] border-jv-ink bg-jv-white p-4 shadow-brutal-sm sm:mt-6 sm:p-5"
        >
          <AnswerSubmissionChart
            :options="options"
            :options-media="question?.options_media || ''"
            :responses="question?.userResponses || []"
          />
        </div>
      </article>
    </section>
  </main>
</template>
