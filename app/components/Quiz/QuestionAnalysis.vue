<script setup>
import {
  Pencil,
  Trash2,
  Clock,
  BarChart3,
  Check,
  X,
  CircleSlash,
} from "lucide-vue-next";

const props = defineProps({
  question: {
    type: Object,
    required: true,
    default: () => ({}),
  },
  order: {
    type: Number,
    required: false,
    default: 0,
  },
  isAdminAnalysis: {
    type: Boolean,
    required: false,
    default: false,
  },
  isForQuiz: {
    type: Boolean,
    required: false,
    default: false,
  },
  isEditable: {
    type: String,
    required: false,
    default: "",
  },
});

const emits = defineEmits(["deleteQuestion", "editQuestion"]);

const deleteQuestion = (questionId) => emits("deleteQuestion", questionId);
const editQuestion = (questionId) => emits("editQuestion", questionId);

const correctPercent = computed(() => {
  const v = Number(props.question?.correctPercentage);
  return Number.isFinite(v) ? v : 0;
});

const correctPercentColor = computed(() =>
  correctPercent.value >= 50 ? "bg-jv-mint" : "bg-jv-salmon"
);

const toIntSet = (val) => {
  if (val == null) return new Set();
  if (Array.isArray(val)) {
    return new Set(val.map((v) => Number(v)).filter((n) => Number.isFinite(n)));
  }
  if (typeof val === "number") {
    return Number.isFinite(val) ? new Set([val]) : new Set();
  }
  if (typeof val === "string") {
    return new Set(
      val
        .replace(/[[\]\s]/g, "")
        .split(",")
        .map((s) => Number(s))
        .filter((n) => Number.isFinite(n))
    );
  }
  return new Set();
};

const isFetchableUrl = (val) =>
  typeof val === "string" && /^(https?:|data:|blob:)/i.test(val.trim());
const renderQuestionImage = computed(
  () =>
    props.question?.question_media === "image" &&
    isFetchableUrl(props.question?.resource)
);

const userResultChip = computed(() => {
  if (!props.question?.is_attend) {
    return {
      label: "Not Attempted",
      bg: "bg-jv-salmon",
      icon: CircleSlash,
    };
  }
  if (props.question?.question_type === "survey") {
    return {
      label: "Attempted",
      bg: "bg-jv-mint",
      icon: Check,
    };
  }
  const correct = toIntSet(props.question?.correct_answer);
  const picked = toIntSet(props.question?.selected_answer?.String);
  if (picked.size === 0) {
    return {
      label: "Not Attempted",
      bg: "bg-jv-salmon",
      icon: CircleSlash,
    };
  }
  const isCorrect =
    picked.size === correct.size && [...picked].every((p) => correct.has(p));
  return isCorrect
    ? { label: "Correct", bg: "bg-jv-mint", icon: Check }
    : { label: "Incorrect", bg: "bg-jv-salmon", icon: X };
});
</script>

<template>
  <div class="w-full">
    <div
      v-if="props.order !== 0 || props.isEditable"
      class="flex items-center justify-between gap-3"
    >
      <span
        v-if="props.order !== 0"
        class="inline-flex items-center gap-2 rounded-full border-[2px] border-jv-ink bg-jv-yellow px-3 py-1 font-feature text-[12px] font-black uppercase tracking-[0.08em] text-jv-ink shadow-brutal-sm sm:text-[13px]"
      >
        Question {{ props.order }}
      </span>

      <div v-if="props.isEditable" class="flex items-center gap-2">
        <button
          type="button"
          class="grid size-9 place-items-center rounded-[8px] border-[2px] border-jv-ink bg-jv-yellow text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[2deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none"
          title="Edit question"
          @click="editQuestion(props.question?.question_id)"
        >
          <Pencil class="size-4" :stroke-width="2.4" />
        </button>
        <button
          type="button"
          class="grid size-9 place-items-center rounded-[8px] border-[2px] border-jv-ink bg-jv-coral text-white shadow-brutal-sm transition-transform hover:rotate-[-2deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none"
          title="Delete question"
          data-bs-toggle="modal"
          data-bs-target="#deleteQuestion"
        >
          <Trash2 class="size-4" :stroke-width="2.4" />
        </button>
        <DeleteDialog
          id="deleteQuestion"
          @confirm-delete="deleteQuestion(props.question?.question_id)"
        />
      </div>
    </div>

    <h3
      class="mt-3 font-headings text-[20px] leading-snug text-jv-ink sm:text-[24px] md:text-[28px]"
    >
      {{ props.question?.question }}
    </h3>

    <div v-if="renderQuestionImage" class="mt-4 flex justify-center">
      <img
        :src="props.question?.resource"
        alt="Question image"
        class="max-h-[220px] w-auto border-[3px] border-jv-ink bg-jv-white object-contain shadow-brutal-sm"
      />
    </div>
    <div v-if="props.question?.question_media === 'code'" class="mt-4">
      <CodeBlockComponent :code="props.question?.resource" />
    </div>

    <div
      v-if="props.isAdminAnalysis && !props.isForQuiz"
      class="mt-4 flex flex-wrap items-center gap-2 sm:gap-3"
    >
      <span
        class="inline-flex items-center gap-2 rounded-full border-[2px] border-jv-ink bg-jv-white px-3 py-1 font-feature text-[12px] font-black text-jv-ink shadow-brutal-sm sm:text-[13px]"
      >
        <Clock class="size-3.5" :stroke-width="2.4" />
        AVG.
        {{ Math.abs((props.question?.avg_response_time / 1000).toFixed(2)) }}s
        <span class="text-jv-muted">/ {{ props.question?.duration }}s</span>
      </span>
      <span
        class="inline-flex items-center rounded-full border-[2px] border-jv-ink bg-jv-lavender px-3 py-1 font-feature text-[12px] font-black text-jv-ink shadow-brutal-sm sm:text-[13px]"
      >
        {{ props.question?.type === 1 ? "M.C.Q." : "Survey" }}
      </span>
      <span
        :class="[
          'inline-flex items-center gap-2 rounded-full border-[2px] border-jv-ink px-3 py-1 font-feature text-[12px] font-black text-jv-ink shadow-brutal-sm sm:text-[13px]',
          correctPercentColor,
        ]"
      >
        <BarChart3 class="size-3.5" :stroke-width="2.4" />
        {{ correctPercent.toFixed(0) }}% correct
      </span>
    </div>

    <div
      v-else-if="!props.isAdminAnalysis && !props.isForQuiz"
      class="mt-4 flex flex-wrap items-center gap-2 sm:gap-3"
    >
      <span
        class="inline-flex items-center gap-2 rounded-full border-[2px] border-jv-ink bg-jv-white px-3 py-1 font-feature text-[12px] font-black text-jv-ink shadow-brutal-sm sm:text-[13px]"
      >
        <Clock class="size-3.5" :stroke-width="2.4" />
        <template v-if="props.question?.response_time > 0">
          {{ (props.question?.response_time / 1000).toFixed(2) }}s
        </template>
        <template v-else>—</template>
      </span>
      <span
        :class="[
          'inline-flex items-center gap-1.5 rounded-full border-[2px] border-jv-ink px-3 py-1 font-feature text-[12px] font-black text-jv-ink shadow-brutal-sm sm:text-[13px]',
          userResultChip.bg,
        ]"
      >
        <component
          :is="userResultChip.icon"
          class="size-3.5"
          :stroke-width="2.6"
        />
        {{ userResultChip.label }}
      </span>
    </div>
  </div>
</template>
