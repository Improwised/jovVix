<script setup>
const props = defineProps({
  data: {
    type: [Array, Object],
    required: true,
    default: () => [],
  },
});

const questionsAnalysis = computed(() => {
  if (!Array.isArray(props.data)) return [];
  return props.data.filter(
    (item) => item && !Object.prototype.hasOwnProperty.call(item, "rank")
  );
});
</script>

<template>
  <section
    class="mx-auto w-full max-w-[1180px]"
    aria-label="Per-question review"
  >
    <header class="mx-auto mb-4 max-w-[820px] sm:mb-6">
      <h2
        class="font-headings text-[22px] leading-none text-jv-ink sm:text-[28px]"
      >
        Question Review
      </h2>
      <p
        class="mt-1.5 font-body text-[12px] font-bold text-jv-muted sm:text-[14px]"
      >
        See what you got right and where you slipped
      </p>
    </header>

    <div
      v-if="questionsAnalysis.length > 0"
      class="mx-auto flex max-w-[820px] flex-col gap-5 sm:gap-6"
    >
      <article
        v-for="(quiz, index) in questionsAnalysis"
        :key="quiz?.question_id || index"
        :class="[
          'relative jv-border-rough bg-jv-white p-5 shadow-brutal sm:p-6',
          index % 2 === 0 ? '-rotate-[0.3deg]' : 'rotate-[0.3deg]',
        ]"
      >
        <span
          class="absolute left-1/2 top-[-10px] z-10 h-3.5 w-14 -translate-x-1/2 rotate-[1deg] bg-jv-coral"
          aria-hidden="true"
        ></span>

        <QuizQuestionAnalysis :question="quiz" :order="index + 1" />
        <QuizOptionsAnalysis
          :options="quiz?.options"
          :correct-answer="quiz?.correct_answer"
          :selected-answer="quiz?.selected_answer?.String"
          :options-media="quiz?.options_media"
        />
      </article>
    </div>

    <div
      v-else
      class="mx-auto mt-3 flex w-fit items-center gap-2 rounded-full border-[2px] border-dashed border-jv-ink/30 px-4 py-2 font-body text-[13px] font-bold text-jv-muted sm:text-[14px]"
    >
      <span class="size-2 rounded-full bg-jv-ink/30" aria-hidden="true"></span>
      No question data
    </div>
  </section>
</template>
