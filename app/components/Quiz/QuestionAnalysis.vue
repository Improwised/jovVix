<template>
  <div class="row m-2">
    <div class="col-lg-12">
      <div>
        <div class="d-flex row">
          <strong v-if="props.order !== 0" class="text-primary col-6"
            >Question: {{ props.order }}
          </strong>
          <div class="col-6 d-flex justify-content-end">
            <NuxtLink
              v-if="quizId"
              type="button"
              class="me-5"
              :to="`/admin/quiz/list-quiz/${props.quizId}/${props.question.question_id}`"
            >
              <font-awesome-icon :icon="['fas', 'pen-to-square']" />
            </NuxtLink>
          </div>
        </div>
        <h3 class="font-bold">{{ props.question?.question }}</h3>
      </div>
    </div>
    <div
      v-if="props.question?.question_media === 'image'"
      class="d-flex align-items-center justify-content-center"
    >
      <img
        :src="`${props.question?.resource}`"
        :alt="`${props.question?.resource}`"
        class="rounded img-thumbnail"
      />
    </div>
    <CodeBlockComponent
      v-if="props.question?.question_media === 'code'"
      :code="props.question?.resource"
    />
    <div
      v-if="props.isAdminAnalysis && !props.isForQuiz"
      class="col-lg-12 d-flex flex-wrap align-items-center justify-content-around"
    >
      <span class="bg-light-primary rounded px-2 text-dark">
        AVG. Response Time:
        {{ Math.abs((props.question?.avg_response_time / 1000).toFixed(2)) }}/
        {{ props.question?.duration }} seconds
      </span>
      <span
        v-if="props.question?.type === 1"
        class="badge bg-light-info m-1 text-dark"
        >M.C.Q.</span
      >
      <span v-else class="badge bg-light-info m-1 text-dark">Survey</span>
      <v-progress-circular
        class="mt-2"
        :model-value="props.question?.correctPercentage"
        :rotate="360"
        :size="60"
        :width="5"
        :color="props.question?.correctPercentage >= 50 ? 'teal' : '#D2042D'"
      >
        {{ props.question?.correctPercentage.toFixed(0) }}%
      </v-progress-circular>
    </div>
    <div
      v-else-if="!props.isAdminAnalysis && !props.isForQuiz"
      class="col-lg-12 d-flex flex-wrap align-items-center justify-content-around mt-1"
    >
      <span
        v-if="props.question?.response_time > 0"
        class="bg-light-primary rounded px-2 text-dark"
      >
        Response Time:
        {{ (props.question?.response_time / 1000).toFixed(2) }} seconds
      </span>
      <span v-else class="bg-light-primary rounded px-2 text-dark">
        Response Time: -
      </span>
      <span
        v-if="props.question?.is_attend"
        class="badge bg-success m-1 text-white"
        >Attempted</span
      >
      <span v-else class="badge bg-danger m-1 text-white">Not Attempted</span>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  question: {
    type: Object,
    required: true,
    default: () => {
      return {};
    },
  },
  order: {
    type: Number,
    required: false,
    default: 0,
  },
  isAdminAnalysis: {
    type: Boolean,
    required: true,
    default: false,
  },
  isForQuiz: {
    type: Boolean,
    required: false,
    default: false,
  },
  quizId: {
    type: String,
    required: false,
    default: "",
  },
});
</script>
