<template>
  <div class="row d-flex align-items-stretch m-2">
    <div
      v-for="(option, order) in props.options"
      :key="order"
      class="col-lg-6 col-md-12"
    >
      <div
        v-if="props.correctAnswer.includes(Number(order))"
        class="bg-light-success option-box"
      >
        <Option
          :order="Number(order)"
          :option="option"
          :selected="props.selectedAnswers[order]?.length || 0"
          :is-correct="true"
          :options-media="props.optionsMedia"
          :is-admin-analysis="props.isAdminAnalysis"
        />
      </div>
      <div
        v-else-if="
          props.selectedAnswer.includes(order) &&
          !props.correctAnswer.includes(order)
        "
        class="bg-light-danger option-box wrong-option"
      >
        <Option
          :order="Number(order)"
          :option="option"
          :selected="props.selectedAnswers[order]?.length || 0"
          :options-media="props.optionsMedia"
          :is-admin-analysis="props.isAdminAnalysis"
        />
      </div>
      <div v-else class="option-box wrong-option">
        <Option
          :order="Number(order)"
          :option="option"
          :selected="props.selectedAnswers[order]?.length || 0"
          :options-media="props.optionsMedia"
          :is-admin-analysis="props.isAdminAnalysis"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  options: {
    type: Object,
    required: true,
    default: () => {
      return {};
    },
  },
  correctAnswer: {
    type: String,
    required: true,
    default: () => {
      return [];
    },
  },
  selectedAnswer: {
    type: String,
    required: false,
    default: "",
  },
  selectedAnswers: {
    type: Object,
    required: false,
    default: () => {
      return {};
    },
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
</script>

<style scoped>
.option-box {
  min-height: 70px;
  padding-top: 3px;
  border-radius: 30px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 5px;
  margin-top: 3px;
}

.wrong-option {
  border: 1px solid var(--bs-light-primary);
}
</style>
