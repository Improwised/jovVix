<script setup>
// core dependencies
import { useNuxtApp } from "nuxt/app";
import { useToast } from "vue-toastification";
import { useMusicStore } from "~~/store/music";
const musicStore = useMusicStore();
const { getMusic } = musicStore;

const music = computed(() => {
  return getMusic();
});

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
  let counterSound = null;

  if (music.value) {
    counterSound = new Audio("/music/clock.mp3");
  }
  if (message.event == app.$GetQuestion) {
    question.value = message.data;
    count.value = null;
    answer.value = [];
    handleTimer();
  } else if (message.event == app.$Counter) {
    question.value = null;
    count.value = parseInt(props.data.data.count);
    time.value = 0;
    handleCounter(counterSound);
    if (music.value && counterSound) {
      counterSound.play();
    }
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

function handleCounter(counterSound) {
  clearInterval(counter.value);
  counter.value = setInterval(() => {
    count.value -= 1;
    if (count.value <= 0) {
      clearInterval(counter.value);
      count.value = app.$ReadyMessage;
      counter.value = null;
      if (counterSound) {
        counterSound.pause();
      }
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
    :music-component="true"
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

    <!-- Question -->
    <QuizQuestionAnalysis
      :question="props.data?.data"
      :is-for-quiz="true"
      :order="question.no"
    />

    <!-- Options -->
    <div class="row d-flex align-items-stretch m-2">
      <div
        v-for="(value, key) in question.options"
        :key="key"
        class="col-lg-6 col-md-12"
      >
        <input
          v-if="!isAdmin"
          :id="`${key}`"
          v-model="answer"
          type="radio"
          class="option"
          name="answer"
          :value="{ key }"
          :disabled="isSubmitted"
          @change="handleSubmit"
        />
        <label :for="`${key}`" class="option-box wrong-option">
          <Option
            :order="Number(key)"
            :option="value"
            :options-media="question.options_media"
          />
        </label>
      </div>
    </div>
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
.option-box {
  min-height: 70px;
  padding-top: 3px;
  border-radius: 30px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 5px;
  margin-top: 3px;
  transition: all 0.3s ease;
}

.wrong-option {
  border: 2px solid var(--bs-light-primary);
}

input[type="radio"]:checked + .option-box {
  border-color: #3c3535eb;
  transform: scale(1.05);
  opacity: 1;
}

input[type="radio"]:not(checked) + .option-box {
  opacity: 0.7;
}

.option {
  display: none;
}
</style>
