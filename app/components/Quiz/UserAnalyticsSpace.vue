<script setup>
import {
  CheckCircle2,
  XCircle,
  MinusCircle,
  ListChecks,
  X,
} from "lucide-vue-next";
import CodeBlockComponent from "@/components/CodeBlockComponent.vue";

const props = defineProps({
  data: {
    type: Object,
    required: true,
  },
  userName: {
    type: String,
    required: true,
  },
});

const analysis = ref({});
const unattemptedWidth = ref(0);
const incorrectWidth = ref(0);
const isOpen = ref(false);

onMounted(() => {
  analysis.value = questionsAnalysis(props.data);
  unattemptedWidth.value =
    (analysis.value?.unAttemptedQuestions / analysis.value?.totalQuestions) *
    100;
  incorrectWidth.value =
    100 - analysis.value?.accuracy - unattemptedWidth.value;
});

const avatar = computed(() => {
  const rankData = props.data?.filter((item) => item.hasOwnProperty("rank"));

  if (rankData) {
    return getAvatarUrlByName(rankData[0]?.avatar);
  }
  return getAvatarUrlByName("");
});

// Strip the rank entry; only real question rows.
const userQuestions = computed(
  () => props.data?.filter((item) => !item.hasOwnProperty("rank")) || []
);

const parseKeyList = (value) => {
  if (value === null || value === undefined) return [];
  if (Array.isArray(value)) return value.map((n) => Number(n));
  let str = String(value).trim();
  if (str.startsWith("[") && str.endsWith("]")) {
    str = str.slice(1, -1);
  }
  if (!str) return [];
  return str
    .split(",")
    .map((s) => Number(s.trim()))
    .filter((n) => !Number.isNaN(n));
};

const optionLabel = (key) => String.fromCharCode(64 + Number(key));

const isCorrectKey = (quiz, key) =>
  parseKeyList(quiz?.correct_answer).includes(Number(key));

const isSelectedKey = (quiz, key) =>
  parseKeyList(quiz?.selected_answer?.String).includes(Number(key));

const optionRowClass = (quiz, key) => {
  if (isCorrectKey(quiz, key)) {
    return "border-l-4 border-l-jv-accent-green bg-jv-accent-green/25 pl-2";
  }
  if (isSelectedKey(quiz, key)) {
    return "border-l-4 border-l-jv-coral bg-jv-coral/15 pl-2";
  }
  return "border-l-4 border-l-transparent";
};

const statusFor = (quiz) => {
  if (!quiz?.is_attend) {
    return { label: "Skipped", tone: "text-jv-muted" };
  }
  if (quiz?.question_type === "survey") {
    return { label: "Submitted", tone: "text-jv-ink" };
  }
  const correctSet = new Set(parseKeyList(quiz?.correct_answer));
  const selectedSet = new Set(parseKeyList(quiz?.selected_answer?.String));
  if (
    selectedSet.size > 0 &&
    selectedSet.size === correctSet.size &&
    [...selectedSet].every((n) => correctSet.has(n))
  ) {
    return { label: "Correct", tone: "text-jv-accent-green" };
  }
  return { label: "Incorrect", tone: "text-jv-coral" };
};

const responseLabel = (quiz) =>
  quiz?.response_time > 0
    ? `${(quiz.response_time / 1000).toFixed(2)} sec`
    : "-";

const onKeydown = (e) => {
  if (e.key === "Escape") isOpen.value = false;
};

watch(isOpen, (open) => {
  if (typeof document === "undefined") return;
  if (open) {
    document.addEventListener("keydown", onKeydown);
    document.body.style.overflow = "hidden";
  } else {
    document.removeEventListener("keydown", onKeydown);
    document.body.style.overflow = "";
  }
});

onBeforeUnmount(() => {
  if (typeof document !== "undefined") {
    document.removeEventListener("keydown", onKeydown);
    document.body.style.overflow = "";
  }
});
</script>

<template>
  <!-- Participant card -->
  <button
    type="button"
    class="w-full text-left jv-border-rough bg-jv-white p-5 sm:p-6 shadow-brutal-sm transition-transform hover:translate-y-[-2px] hover:shadow-brutal cursor-pointer"
    @click="isOpen = true"
  >
    <!-- Top row: avatar + name | stats -->
    <div
      class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
    >
      <div class="flex items-center gap-4 min-w-0">
        <img
          :src="avatar"
          :alt="props.data[0].firstname"
          class="size-14 shrink-0 rounded-full border-[3px] border-jv-ink bg-jv-slate object-cover shadow-brutal-sm"
        />
        <div class="min-w-0">
          <div
            class="font-headings text-[22px] leading-tight text-jv-ink sm:text-[24px] truncate"
          >
            {{ props.data[0].firstname }}
          </div>
          <div class="text-[14px] font-semibold text-jv-muted truncate">
            ({{ props.data[0].username }})
          </div>
        </div>
      </div>

      <div class="flex items-end gap-6 sm:gap-8">
        <div class="text-center">
          <div class="text-[22px] font-black text-jv-ink sm:text-[24px]">
            {{ analysis?.rank }}
          </div>
          <div
            class="mt-0.5 text-[11px] font-black uppercase tracking-[0.14em] text-jv-muted"
          >
            Rank
          </div>
        </div>
        <div class="text-center">
          <div class="text-[22px] font-black text-jv-ink sm:text-[24px]">
            {{ analysis?.accuracy }}%
          </div>
          <div
            class="mt-0.5 text-[11px] font-black uppercase tracking-[0.14em] text-jv-muted"
          >
            Accuracy
          </div>
        </div>
        <div class="text-center">
          <div class="text-[22px] font-black text-jv-ink sm:text-[24px]">
            {{ analysis?.totalScore }}
          </div>
          <div
            class="mt-0.5 text-[11px] font-black uppercase tracking-[0.14em] text-jv-muted"
          >
            Score
          </div>
        </div>
      </div>
    </div>

    <!-- Multi-segment progress bar -->
    <div
      class="mt-5 flex h-2.5 w-full overflow-hidden rounded-full bg-jv-slate"
    >
      <div
        class="h-full bg-jv-accent-green"
        :style="{ width: `${analysis?.accuracy || 0}%` }"
      ></div>
      <div
        class="h-full bg-jv-coral"
        :style="{ width: `${Math.max(incorrectWidth, 0)}%` }"
      ></div>
      <div
        class="h-full bg-jv-ink/20"
        :style="{ width: `${Math.max(unattemptedWidth, 0)}%` }"
      ></div>
    </div>

    <!-- Bottom counts row -->
    <div
      class="mt-4 flex flex-wrap items-center justify-center gap-x-6 gap-y-2 text-[14px] font-bold text-jv-ink"
    >
      <span class="inline-flex items-center gap-1.5">
        <CheckCircle2 class="size-4 text-jv-accent-green" :stroke-width="2.6" />
        {{ analysis?.correctAnwers }}
      </span>
      <span class="inline-flex items-center gap-1.5">
        <XCircle class="size-4 text-jv-coral" :stroke-width="2.6" />
        {{ analysis?.wrongAnwers }}
      </span>
      <span class="inline-flex items-center gap-1.5">
        <MinusCircle class="size-4 text-jv-muted" :stroke-width="2.6" />
        {{ analysis?.unAttemptedQuestions }}
      </span>
      <span
        v-if="analysis?.totalSurveyQuestions > 0"
        class="inline-flex items-center gap-1.5"
      >
        <ListChecks class="size-4 text-jv-muted" :stroke-width="2.6" />
        {{ analysis?.attemptedSurveyQuestions }} /
        {{ analysis?.totalSurveyQuestions }}
      </span>
    </div>
  </button>

  <!-- Per-user Questions Analysis Modal (Vue-controlled, JovVix style) -->
  <Teleport to="body">
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
      role="dialog"
      aria-modal="true"
      :aria-labelledby="`${props.userName}-title`"
    >
      <div class="absolute inset-0 bg-jv-ink/60" @click="isOpen = false"></div>

      <div
        class="relative z-10 flex max-h-[92vh] w-full max-w-[1000px] rotate-[-0.4deg] flex-col border-[4px] border-jv-ink bg-jv-white shadow-brutal-lg"
      >
        <!-- Dark header -->
        <div
          class="flex flex-col gap-4 border-b-[3px] border-jv-ink bg-jv-ink px-5 py-4 text-jv-white sm:flex-row sm:items-center sm:justify-between sm:px-6"
        >
          <div class="flex items-center gap-3 min-w-0">
            <img
              :src="avatar"
              :alt="props.data[0].firstname"
              class="size-12 shrink-0 rounded-full border-[2px] border-jv-white/30 bg-jv-slate object-cover"
            />
            <div class="min-w-0">
              <h2
                :id="`${props.userName}-title`"
                class="font-headings text-[22px] leading-tight text-jv-white sm:text-[26px] truncate"
              >
                {{ props.data[0].firstname }}
              </h2>
              <div
                class="text-[12px] font-bold uppercase tracking-[0.14em] text-jv-yellow"
              >
                Questions Analysis
              </div>
            </div>
          </div>

          <div
            class="grid grid-cols-3 gap-x-5 gap-y-2 border-[3px] border-jv-white/20 bg-jv-ink/60 px-4 py-2 sm:gap-x-7"
          >
            <div>
              <div
                class="text-[10px] font-black uppercase tracking-[0.14em] text-jv-white/60"
              >
                Rank
              </div>
              <div class="mt-0.5 text-[18px] font-black text-jv-white">
                {{ analysis?.rank }}
              </div>
            </div>
            <div>
              <div
                class="text-[10px] font-black uppercase tracking-[0.14em] text-jv-white/60"
              >
                Accuracy
              </div>
              <div class="mt-0.5 text-[18px] font-black text-jv-accent-green">
                {{ analysis?.accuracy }}%
              </div>
            </div>
            <div>
              <div
                class="text-[10px] font-black uppercase tracking-[0.14em] text-jv-white/60"
              >
                Score
              </div>
              <div class="mt-0.5 text-[18px] font-black text-jv-white">
                {{ analysis?.totalScore }}
              </div>
            </div>
          </div>

          <button
            type="button"
            class="absolute right-3 top-3 grid size-9 place-items-center text-jv-white transition-transform hover:rotate-[6deg] sm:hidden"
            aria-label="Close"
            @click="isOpen = false"
          >
            <X class="size-6" :stroke-width="2.4" />
          </button>
        </div>

        <!-- Body: per-question list -->
        <div
          class="flex-1 overflow-y-auto bg-jv-canvas px-4 py-5 sm:px-6 sm:py-6"
        >
          <div class="flex flex-col gap-5">
            <article
              v-for="(quiz, index) in userQuestions"
              :key="index"
              class="jv-border-rough bg-jv-white p-4 shadow-brutal-sm sm:p-5"
            >
              <!-- Header -->
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
                      Response
                    </div>
                    <div class="mt-0.5 text-[14px] font-bold text-jv-ink">
                      {{ responseLabel(quiz) }}
                    </div>
                  </div>
                  <div>
                    <div
                      class="text-[10px] font-bold uppercase tracking-[0.12em] text-jv-muted"
                    >
                      Status
                    </div>
                    <div
                      class="mt-0.5 text-[14px] font-bold"
                      :class="statusFor(quiz).tone"
                    >
                      {{ statusFor(quiz).label }}
                    </div>
                  </div>
                </div>
              </div>

              <!-- Resource -->
              <div
                v-if="quiz.question_media === 'image' && quiz.resource"
                class="mt-3 flex justify-center rounded-md bg-jv-canvas p-2"
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
              >
                <CodeBlockComponent :code="quiz.resource" />
              </div>

              <!-- Options -->
              <ul class="mt-4 flex flex-col">
                <li
                  v-for="(option, key) in quiz.options"
                  :key="key"
                  class="flex min-w-0 items-center gap-3 border-b border-jv-ink/10 py-3 pl-3 pr-2 text-[15px] font-medium text-jv-ink last:border-b-0"
                  :class="optionRowClass(quiz, key)"
                >
                  <span
                    class="w-5 shrink-0 text-[14px] font-bold text-jv-coral"
                  >
                    {{ optionLabel(key) }}.
                  </span>

                  <div
                    v-if="quiz.options_media === 'image' && option"
                    class="flex min-w-0 flex-1 justify-start"
                  >
                    <img
                      :src="option"
                      :alt="`Option ${optionLabel(key)}`"
                      class="max-h-32 w-auto max-w-full object-contain"
                    />
                  </div>

                  <div
                    v-else-if="quiz.options_media === 'code' && option"
                    class="min-w-0 flex-1 overflow-x-auto"
                  >
                    <CodeBlockComponent :code="option" />
                  </div>

                  <span v-else class="min-w-0 flex-1 break-words">
                    {{ option }}
                  </span>
                </li>
              </ul>
            </article>

            <div
              v-if="!userQuestions.length"
              class="jv-border-rough bg-jv-white p-5 text-center text-[16px] font-semibold text-jv-muted shadow-brutal-sm"
            >
              No questions found for this participant.
            </div>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>
