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
const question = ref({
  question: "",
  options: {},
});
const answer = ref([]);

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
  console.log(message);
  if (message.event == app.$GetQuestion) {
    question.value.question = message.data.question;
    question.value.options = message.data.options;
    answer.value = null;
  }
}

function handleSubmit(e) {
  e.preventDefault();
  emits('sendAnswer')
  console.log(answer.value);
}
</script>

<template>
  <Frame page-title="Question" page-message="let's play">
    <div>{{ question.question }}</div>
    <div class="d-flex">
      <div
        v-for="(value, key) in question.options"
        :key="key"
        class="flex-grow-1 border m-1 rounded p-1"
      >
        <div class="form-check form-check-inline">
          <input
            id="answer"
            v-model="answer"
            class="form-check-input"
            type="radio"
            name="answer"
            :value="{ key }"
          />
          <label class="form-check-label" for="answer">{{ value }}</label>
        </div>
      </div>
    </div>
    <button type="button" class="btn btn-primary mt-3" @click="handleSubmit">
      submit
    </button>
  </Frame>
</template>
