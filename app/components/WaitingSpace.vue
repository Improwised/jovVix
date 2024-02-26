<script setup>
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
const emits = defineEmits(["startQuiz"]);

// custom refs
const code = ref(123455);

// watchers
watch(
  () => props.data,
  (new_val) => {
    handleEvent(new_val);
  },
  { deep: true }
);

// event handlers
function start_quiz(e) {
  e.preventDefault();
  emits("startQuiz");
}

// main function
function handleEvent(message) {
  if (message.event == "send code to admin") {
    code.value = message.data.code;
  }
}
</script>

<template>
  <FrameLayout
    v-if="isAdmin"
    :page-title="'Ready-shady-go'"
    :page-welcome-message="'you can start quiz by pressing start button'"
  >
    <form @submit="start_quiz">
      <div class="mb-3 pe-3">
        <label for="code" class="form-label">Code</label>
        <v-otp-input
          v-model="code"
          max-width="500"
          min-height="20"
          type="number"
          disabled
        ></v-otp-input>
      </div>
      <button type="submit" class="btn btn-primary btn-lg bg-primary">
        Start
      </button>
    </form>
  </FrameLayout>
  <FrameLayout
    v-else
    :page-title="'Ready-shady-go'"
    :page-welcome-message="'quiz is about to start'"
  >
    <h2>Rules</h2>
    <div></div>
  </FrameLayout>
</template>
