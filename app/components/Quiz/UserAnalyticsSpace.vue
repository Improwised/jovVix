<script setup>
import {
  CheckCircle2,
  XCircle,
  MinusCircle,
  ListChecks,
  X,
} from "lucide-vue-next";

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
        class="relative z-10 flex max-h-[90vh] w-full max-w-3xl flex-col jv-border-rough bg-jv-white shadow-brutal-lg"
      >
        <div
          class="flex items-center justify-between gap-4 border-b-2 border-jv-ink/15 px-5 py-4 sm:px-6"
        >
          <h2
            :id="`${props.userName}-title`"
            class="font-headings text-[22px] leading-tight text-jv-ink sm:text-[26px]"
          >
            Questions Analysis
          </h2>
          <button
            type="button"
            class="flex size-9 shrink-0 items-center justify-center jv-border-rough bg-jv-white text-jv-ink shadow-brutal-sm hover:bg-jv-yellow/50 transition-colors"
            aria-label="Close"
            @click="isOpen = false"
          >
            <X class="size-5" :stroke-width="2.6" />
          </button>
        </div>
        <div class="overflow-y-auto px-5 py-5 sm:px-6 sm:py-6">
          <QuizAnalysis :data="props.data" />
        </div>
      </div>
    </div>
  </Teleport>
</template>
