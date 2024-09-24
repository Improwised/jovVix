<script setup>
// core dependencies
import { useNuxtApp } from "nuxt/app";
import { useToast } from "vue-toastification";
import CodeBlockComponent from "../CodeBlockComponent.vue";

// define nuxt configs
const app = useNuxtApp();
const toast = useToast();

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
});
const emits = defineEmits(["sendAnswer", "askSkip"]);

// custom refs
const question = ref();
const answer = ref([]);
const counter = ref(null);
const count = ref(0);
const timer = ref(null);
const time = ref(0);
const progressValue = computed(() => {
  return (time.value * 100) / question.value.duration;
});
const progressQuiz = computed(() =>
  ((question.value.no / question.value.totalQuestions) * 100).toFixed(0)
);
const isSubmitted = ref(false);

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
  if (message.event == app.$GetQuestion) {
    question.value = message.data;
    count.value = null;
    answer.value = [];
    handleTimer();
  } else if (message.event == app.$Counter) {
    question.value = null;
    count.value = 1;
    time.value = 0;
    handleCounter();
  }
}

function handleTimer() {
  clearInterval(timer.value);
  timer.value = setInterval(() => {
    time.value += 1;
    if (time.value == question.value?.duration + 1) {
      clearInterval(timer.value);
      time.value = -1;
      timer.value = null;
    }
  }, 1000);
}

function handleCounter() {
  clearInterval(counter.value);
  counter.value = setInterval(() => {
    count.value += 1;
    if (parseInt(props.data.data.count) <= count.value) {
      clearInterval(counter.value);
      count.value = app.$ReadyMessage;
      counter.value = null;
    }
  }, 1000);
}

function handleSubmit(e) {
  e.preventDefault();
  let final_val = answer.value;

  if (final_val.length != 0) {
    emits("sendAnswer", [parseInt(answer.value.key)]);
    isSubmitted.value = true;
  } else {
    toast.warning(app.$NoAnswerFound);
  }
}

function handleSkip(e) {
  e.preventDefault();
  emits("askSkip");
}
</script>

<template>
  <Frame
    v-if="question != null"
    :page-title="`Question ${question.no} / ${question.totalQuestions}`"
    page-message="Let's Play"
  >
    <template #sub-title>
      <v-progress-circular
        :model-value="progressValue"
        :rotate="0"
        :size="80"
        :width="13"
        color="primary"
      >
        {{ question.duration - time }}
      </v-progress-circular>
    </template>
    <div>
      <div class="progress">
        <div
          class="progress-bar"
          role="progressbar"
          :aria-valuenow="progressQuiz"
          aria-valuemin="0"
          aria-valuemax="100"
          :style="{ width: ` ${progressQuiz}%` }"
        ></div>
      </div>
    </div>
    <div>
      <span>{{ question.no }}. </span>
      <span>{{ question.question }}</span>
      <img
        v-if="question?.question_media === 'image'"
        :src="`${question?.resource}`"
        :alt="`${question?.resource}`"
        class="rounded img-thumbnail"
      />
      <CodeBlockComponent
        v-if="question?.question_media === 'code'"
        :code="question?.resource"
      />
    </div>
    <div class="d-flex flex-column">
      <div
        v-for="(value, key) in question.options"
        :key="key"
        class="border m-1 rounded p-1"
      >
        <label class="form-check-label d-flex align-items-center">
          <div v-if="!isAdmin" class="form-check form-check-inline me-0">
            <input
              :id="`${key}`"
              v-model="answer"
              class="form-check-input"
              type="radio"
              name="answer"
              :value="{ key }"
              :disabled="isSubmitted"
            />
          </div>
          <p
            v-if="question?.options_media === 'text'"
            class="mb-0 ms-2"
            :for="`${key}`"
          >
            {{ value }}
          </p>
          <img
            v-if="question?.options_media === 'image'"
            :src="`${value}`"
            :alt="`${value}`"
            class="rounded img-thumbnail"
            :for="`${key}`"
          />
          <CodeBlockComponent
            v-if="question?.options_media === 'code'"
            :code="value"
          />
        </label>
      </div>
    </div>
    <button
      v-if="!isAdmin"
      type="button"
      class="btn text-white btn-primary mt-3"
      :disabled="isSubmitted"
      @click="handleSubmit"
    >
      submit
    </button>
    <button
      v-if="isAdmin"
      type="button"
      class="btn text-white btn-primary mt-3"
      @click="handleSkip"
    >
      skip
    </button>
  </Frame>
  <div
    v-else
    class="d-flex align-center justify-content-center align-items-center fs-1"
    style="height: 100vh"
  >
    <div>
      {{ count }}
    </div>
  </div>
</template>

<style scoped>
@media (max-width: 768px) {
  .d-flex.flex-column.flex-md-row {
    flex-direction: column !important;
  }
  .d-flex.flex-column.flex-md-row .border {
    width: 100% !important;
  }
  .form-check-label {
    justify-content: flex-start !important;
  }
  .form-check-input + .mb-0 {
    margin-left: 0.5rem; /* Adjust the spacing between the radio button and the text */
  }
}
</style>
