<script setup>
// core dependencies
import { useNuxtApp } from "nuxt/app";
import { usePush } from "notivue";
import { useMusicStore } from "~~/store/music";
import { Check, SkipForward, Trophy, Volume2, VolumeX } from "lucide-vue-next";

const musicStore = useMusicStore();
const { getMusic, setMusic } = musicStore;

const music = computed(() => {
  return getMusic();
});
const toggleMusic = () => setMusic(!music.value);

// define nuxt configs
const app = useNuxtApp();
const toast = usePush();

// define props and emits
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
  // When true, an admin (public-quiz host) is also allowed to answer questions.
  canPlay: {
    default: false,
    type: Boolean,
    required: false,
  },
});
const emits = defineEmits(["sendAnswer", "askSkip"]);

// A participant can answer when they are not the admin, or when they are a
// public-quiz host who is permitted to play alongside the others.
const answerable = computed(() => !props.isAdmin || props.canPlay);

const clockOffset = ref(0); // Calculated once, used for all timers
const isOffsetCalculated = ref(false);

// Calculate offset only ONCE when we first receive server_time
const calculateClockOffset = (serverTimeString) => {
  if (isOffsetCalculated.value || !serverTimeString) return;

  try {
    const clientReceiveTime = Date.now();
    const serverTime = new Date(serverTimeString).getTime();
    clockOffset.value = serverTime - clientReceiveTime;
    isOffsetCalculated.value = true;
  } catch (error) {
    console.error("Error calculating clock offset:", error);
  }
};

// Get current time corrected with offset
const getCorrectedNow = () => {
  return new Date(Date.now() + clockOffset.value);
};

// custom refs
const question = ref();
const selectedKey = ref(null);
const counter = ref(null);
const count = ref(0);
const timer = ref(null);
const time = ref(0);
const questionStartTime = ref(null);
const questionDuration = ref(0);
const isSubmitted = ref(false);

const remainingSeconds = computed(() =>
  Math.max(0, (questionDuration.value || 0) - time.value)
);

// Determine if current question is the last one
const isLastQuestion = computed(() => {
  if (!question?.value) return false;
  return Number(question.value.no) === Number(question.value.totalQuestions);
});

const countdownDisplay = computed(() => {
  if (typeof count.value === "number" && count.value > 0)
    return String(count.value);
  if (typeof count.value === "string") return "Go!";
  return "Get Ready";
});

const isCountdownNumber = computed(
  () => typeof count.value === "number" && count.value > 0
);

// watchers
watch(
  () => props.data,
  (message) => {
    if (message.status == app.$Fail) {
      toast.error(message.data);
      return;
    }
    handleEvent(message);
  },
  { deep: true, immediate: true }
);

// main function
function handleEvent(message) {
  let counterSound = null;
  // Calculate offset ONCE if server sends time
  if (message.data?.server_time) {
    calculateClockOffset(message.data.server_time);
  }

  if (music.value) {
    counterSound = new Audio("/music/clock.mp3");
  }

  if (message.event == app.$GetQuestion) {
    question.value = message.data;
    questionStartTime.value = new Date(message.data.start_time);
    questionDuration.value = Number(message.data.duration);

    time.value = 0;
    count.value = null;
    selectedKey.value = null;
    isSubmitted.value = false;

    handleTimer();
  } else if (message.event == app.$Counter) {
    question.value = null;
    count.value = parseInt(props.data.data.count);
    time.value = 0;
    handleCounter(counterSound);
    if (music.value && counterSound) {
      counterSound.play();
    }
  }
}

function handleTimer() {
  clearInterval(timer.value);

  if (!questionStartTime.value || !questionDuration.value) {
    console.error("Missing question start time or duration");
    return;
  }

  const duration = questionDuration.value;

  timer.value = setInterval(() => {
    const correctedNow = getCorrectedNow();
    const elapsedSeconds = Math.floor(
      (correctedNow.getTime() - questionStartTime.value.getTime()) / 1000
    );

    time.value = Math.min(Math.max(elapsedSeconds, 0), duration);

    if (elapsedSeconds >= duration) {
      clearInterval(timer.value);
      timer.value = null;
      time.value = duration;
    }
  }, 100);
}

function handleCounter(counterSound) {
  clearInterval(counter.value);
  counter.value = setInterval(() => {
    count.value -= 1;
    if (count.value <= 0) {
      clearInterval(counter.value);
      count.value = app.$ReadyMessage;
      counter.value = null;
      if (counterSound) {
        counterSound.pause();
      }
    }
  }, 1000);
}

function handleOptionClick(key) {
  if (!answerable.value || isSubmitted.value) return;
  if (props.data?.data?.options === undefined) {
    toast.warning(app.$NoAnswerFound);
    return;
  }
  selectedKey.value = Number(key);
  emits("sendAnswer", [Number(key)]);
  isSubmitted.value = true;
}

function handleSkip(e) {
  e.preventDefault();
  emits("askSkip");
}

// Cleanup on unmount
onUnmounted(() => {
  if (timer.value) {
    clearInterval(timer.value);
  }
  if (counter.value) {
    clearInterval(counter.value);
  }
});
</script>

<template>
  <!-- Countdown view -->
  <main
    v-if="question == null"
    class="flex min-h-[75vh] flex-col items-center justify-center bg-jv-canvas px-4 py-10 text-jv-ink sm:py-14"
  >
    <article
      class="relative w-full max-w-[460px] -rotate-[0.4deg] jv-border-rough bg-jv-white px-6 py-9 shadow-brutal sm:px-10 sm:py-12"
    >
      <span
        class="absolute left-1/2 top-[-12px] h-4 w-16 -translate-x-1/2 rotate-[1deg] bg-jv-coral"
        aria-hidden="true"
      ></span>

      <p
        class="text-center font-body text-[12px] font-black uppercase tracking-[0.22em] text-jv-muted sm:text-[13px]"
      >
        Get Ready
      </p>

      <transition name="pop" mode="out-in">
        <div
          v-if="isCountdownNumber"
          :key="String(count)"
          class="mx-auto mt-5 grid size-[180px] place-items-center rounded-full border-[3px] border-dashed border-jv-ink bg-jv-yellow font-headings text-[110px] leading-none text-jv-ink sm:size-[220px] sm:text-[140px]"
        >
          {{ countdownDisplay }}
        </div>
        <h1
          v-else
          :key="'msg'"
          class="mt-6 text-center font-headings text-[48px] leading-none text-jv-ink sm:text-[64px]"
        >
          {{ countdownDisplay }}
        </h1>
      </transition>

      <div class="mt-6 flex justify-center gap-2">
        <span class="size-2 rounded-full bg-jv-coral"></span>
        <span class="size-2 rounded-full bg-jv-yellow"></span>
        <span class="size-2 rounded-full bg-jv-mint"></span>
      </div>
    </article>
  </main>

  <!-- Question view -->
  <main
    v-else
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

        <!-- Header: question count + timer -->
        <header class="flex items-start justify-between gap-4">
          <div class="min-w-0 flex-1">
            <h2
              class="font-headings text-[24px] leading-none text-jv-ink min-[420px]:text-[30px] sm:text-[36px] md:text-[40px]"
            >
              Question {{ question.no }}
              <span class="text-jv-muted">/ {{ question.totalQuestions }}</span>
            </h2>
            <p
              class="mt-2 font-body text-[12px] font-bold text-jv-muted sm:text-[14px]"
            >
              Let's Play
            </p>
          </div>

          <div class="flex shrink-0 items-center gap-2 sm:gap-3">
            <button
              type="button"
              class="grid size-10 place-items-center rounded-[8px] border-[2px] border-jv-ink bg-jv-white text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[2deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:size-11"
              :aria-label="music ? 'Mute music' : 'Unmute music'"
              @click="toggleMusic"
            >
              <Volume2 v-if="music" class="size-4" :stroke-width="2.4" />
              <VolumeX v-else class="size-4" :stroke-width="2.4" />
            </button>

            <div
              class="grid size-12 place-items-center rounded-full border-[3px] border-dashed border-jv-ink bg-jv-yellow font-feature text-[20px] font-black text-jv-ink sm:size-[60px] sm:text-[24px]"
              :aria-label="`${remainingSeconds} seconds remaining`"
            >
              {{ remainingSeconds }}
            </div>
          </div>
        </header>

        <div
          class="my-5 border-t-[2px] border-dashed border-jv-ink/30 sm:my-6"
        ></div>

        <!-- Question text -->
        <h3
          class="font-headings text-[20px] leading-snug text-jv-ink min-[420px]:text-[24px] sm:text-[30px] md:text-[36px]"
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

        <!-- Options grid -->
        <div class="mt-6 grid gap-3 sm:mt-7 sm:gap-4 md:grid-cols-2 md:gap-5">
          <button
            v-for="(value, key) in question.options"
            :key="key"
            type="button"
            :disabled="!answerable || isSubmitted"
            :class="[
              'group flex w-full items-center gap-3 border-[2px] border-jv-ink bg-jv-white px-3 py-3 text-left transition-all sm:gap-4 sm:px-4 sm:py-4',
              answerable && !isSubmitted
                ? 'cursor-pointer shadow-brutal-sm hover:translate-x-[2px] hover:translate-y-[2px] hover:shadow-none'
                : 'shadow-brutal-sm',
              isSubmitted && selectedKey === Number(key)
                ? 'bg-jv-yellow outline-[3px] outline-offset-[2px] outline-jv-ink'
                : '',
              isSubmitted && selectedKey !== Number(key) ? 'opacity-50' : '',
              !answerable ? 'cursor-default' : '',
            ]"
            @click="handleOptionClick(key)"
          >
            <span
              class="grid size-9 shrink-0 place-items-center rounded-[6px] border-[2px] border-jv-ink bg-jv-white font-feature text-[14px] font-black text-jv-ink sm:size-11 sm:text-[16px]"
            >
              {{ String.fromCharCode(64 + Number(key)) }}
            </span>
            <div class="min-w-0 flex-1">
              <div
                v-if="question.options_media === 'image'"
                class="flex items-center justify-center"
              >
                <img
                  :src="value"
                  :alt="`Option ${key}`"
                  class="max-h-[120px] w-auto object-contain"
                />
              </div>
              <CodeBlockComponent
                v-else-if="question.options_media === 'code'"
                :code="value"
              />
              <span
                v-else
                class="block break-words font-body text-[15px] font-black leading-snug text-jv-ink sm:text-[18px] md:text-[20px]"
              >
                {{ value }}
              </span>
            </div>
            <Check
              v-if="isSubmitted && selectedKey === Number(key)"
              class="size-5 shrink-0 text-jv-ink"
              :stroke-width="3"
            />
          </button>
        </div>

        <!-- Admin footer: waiting status + skip -->
        <template v-if="isAdmin">
          <div
            class="mt-6 border-t-[2px] border-dashed border-jv-ink/30 sm:mt-7"
          ></div>
          <footer
            class="mt-4 flex flex-col items-stretch justify-between gap-3 sm:flex-row sm:items-center sm:gap-4"
          >
            <div class="flex items-center gap-3">
              <div class="flex gap-1">
                <span
                  class="size-2 rounded-full bg-jv-coral sm:size-2.5"
                ></span>
                <span
                  class="size-2 rounded-full bg-jv-yellow sm:size-2.5"
                ></span>
                <span class="size-2 rounded-full bg-jv-mint sm:size-2.5"></span>
              </div>
              <p
                class="font-body text-[13px] font-bold text-jv-muted sm:text-[14px]"
              >
                Waiting for participants to answer
              </p>
            </div>
            <button
              type="button"
              class="inline-flex h-10 items-center justify-center gap-2 self-end rounded-full border-[2px] border-jv-ink bg-jv-white px-5 font-body text-[14px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-11 sm:self-auto sm:px-6 sm:text-[15px]"
              @click="handleSkip"
            >
              <Trophy
                v-if="isLastQuestion"
                class="size-4"
                :stroke-width="2.4"
              />
              <SkipForward v-else class="size-4" :stroke-width="2.4" />
              <span>{{ isLastQuestion ? "Finish Quiz" : "Skip" }}</span>
            </button>
          </footer>
        </template>

        <!-- User submission acknowledgement -->
        <div
          v-if="answerable && isSubmitted"
          class="mt-6 flex items-center justify-center gap-2 rounded-full border-[2px] border-dashed border-jv-ink/40 px-4 py-2 font-body text-[13px] font-bold text-jv-muted sm:text-[14px]"
        >
          <Check class="size-4" :stroke-width="2.6" />
          Answer locked — waiting for the next question
        </div>
      </article>
    </section>
  </main>
</template>

<style scoped>
.pop-enter-from {
  opacity: 0;
  transform: scale(0.4) rotate(-12deg);
}
.pop-leave-to {
  opacity: 0;
  transform: scale(1.3) rotate(8deg);
}
.pop-enter-active,
.pop-leave-active {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
</style>
