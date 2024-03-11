<script setup>
// core dependencies
import { useNuxtApp } from "nuxt/app";
import { useToast } from "vue-toastification";

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
const emits = defineEmits(["sendAnswer"]);

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

  if (final_val != null) {
    emits("sendAnswer", [parseInt(answer.value.key)]);
    isSubmitted.value = true;
  }
}
</script>

<template>
  <Frame
    v-if="question != null"
    page-title="Question"
    page-message="let's play"
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
    <div>{{ question.no }}</div>
    <div>
      {{ question.question }}
    </div>
    <div>
      {{ question.duration }}
    </div>
    <div class="d-flex">
      <div
        v-for="(value, key) in question.options"
        :key="key"
        class="flex-grow-1 border m-1 rounded p-1"
      >
        <div class="form-check form-check-inline">
          <input
            v-if="!isAdmin"
            v-model="answer"
            class="form-check-input"
            type="radio"
            name="answer"
            :value="{ key }"
            :disabled="isSubmitted"
          />
          <label class="form-check-label">{{ value }}</label>
        </div>
      </div>
    </div>
    <button
      v-if="!isAdmin"
      type="button"
      class="btn btn-primary mt-3"
      :disabled="isSubmitted"
      @click="handleSubmit"
    >
      submit
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
