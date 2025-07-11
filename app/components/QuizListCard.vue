<script setup>
// define props and emits
const props = defineProps({
  details: {
    type: Object,
    default: () => ({}),
    required: true,
  },
  isPlayedQuiz: {
    type: Boolean,
    required: false,
    default: false,
  },
});
</script>
<template>
  <div class="card">
    <div
      class="row gap-3 gap-md-0 m-0 p-2 justify-content-between justify-content-md-start gap-lg-5"
    >
      <div class="col-3 d-flex align-center m-0 p-0">
        <img
          class="img-fluid bg-primary rounded randomImg h-100"
          src="../assets/images/QuestionLogo.webp"
          alt="QuestionLogo.webp"
        />
      </div>
      <div class="col m-0 p-0">
        <div
          class="card-body d-flex flex-column justify-content-between p-0 h-100"
        >
          <h5 class="card-title mb-0">{{ decodeURI(props.details?.title) }}</h5>
          <p class="card-text mb-0">
            {{ props.details?.description?.String }}
          </p>
          <div class="d-flex gap-3 mb-0">
            <span class="text-muted">
              {{ props.details?.total_questions }} Questions</span
            >
          </div>
          <div
            class="d-flex flex-column flex-md-row align-items-start justify-content-md-between"
          >
            <p class="card-text mb-0">
              <small class="text-muted">
                {{ useGetTime(props.details?.created_at) }}
              </small>
            </p>
            <div>
              <UtilsStartQuiz
                v-if="!isPlayedQuiz"
                :quiz-id="props.details?.id"
              />
              <NuxtLink
                v-if="isPlayedQuiz"
                type="button"
                class="btn text-white btn-primary me-0 mx-2"
                :to="`/admin/played_quiz/${props.details?.id}`"
              >
                View Quiz
              </NuxtLink>
              <NuxtLink
                v-else
                type="button"
                class="btn text-white btn-primary me-0 mx-2"
                :to="`/admin/quiz/list-quiz/${props.details?.id}`"
              >
                View Quiz
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<style scoped>
.randomImg {
  object-fit: cover;
}

img {
  height: auto !important;
  max-height: 120px !important;
}
</style>
