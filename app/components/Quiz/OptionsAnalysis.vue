<script setup>
const props = defineProps({
  options: {
    type: [Object, Array],
    required: true,
    default: () => ({}),
  },
  correctAnswer: {
    type: [String, Array],
    required: true,
    default: "",
  },
  selectedAnswer: {
    type: [String, Array],
    required: false,
    default: "",
  },
  selectedAnswers: {
    type: Object,
    required: false,
    default: () => ({}),
  },
  optionsMedia: {
    type: String,
    required: false,
    default: "",
  },
  isAdminAnalysis: {
    type: Boolean,
    required: false,
    default: false,
  },
});

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

const correctSet = computed(() => toIntSet(props.correctAnswer));
const pickedSet = computed(() => toIntSet(props.selectedAnswer));

const isCorrect = (order) => correctSet.value.has(Number(order));
const isPicked = (order) => pickedSet.value.has(Number(order));
</script>

<template>
  <div class="mt-5 grid gap-3 sm:gap-4 md:grid-cols-2">
    <div
      v-for="(option, order) in props.options"
      :key="order"
      :class="[
        'flex w-full items-center gap-3 border-[2px] border-jv-ink px-3 py-3 shadow-brutal-sm sm:gap-4 sm:px-4 sm:py-4',
        isCorrect(order)
          ? 'bg-jv-mint'
          : isPicked(order)
          ? 'bg-jv-salmon/50'
          : 'bg-jv-white',
      ]"
    >
      <Option
        :order="Number(order)"
        :option="option"
        :selected="props.selectedAnswers[order]?.length || 0"
        :is-correct="isCorrect(order)"
        :is-picked="isPicked(order) && !isCorrect(order)"
        :options-media="props.optionsMedia"
        :is-admin-analysis="props.isAdminAnalysis"
      />
    </div>
  </div>
</template>
