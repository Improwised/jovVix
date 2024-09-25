<script setup>
const props = defineProps({
  data: {
    type: Object,
    required: true,
  },
});

const questionsAnalysis = computed(() => {
  // to remove rank object from data
  const filteredData = props.data?.filter(
    (item) => !item.hasOwnProperty("rank")
  );
  return filteredData;
});
</script>

<template>
  <div class="container mt-3">
    <div
      v-for="(quiz, index) in questionsAnalysis"
      :key="index"
      class="card mb-3 row"
    >
      <QuizQuestionAnalysis :question="quiz" :order="index + 1" />
      <QuizOptionsAnalysis
        :options="quiz?.options"
        :correct-answer="quiz?.correct_answer"
        :selected-answer="quiz?.selected_answer?.String"
        :options-media="quiz?.options_media"
      />
    </div>
  </div>
</template>
