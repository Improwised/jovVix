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
const emits = defineEmits(["startQuiz"]);

// custom refs
const code = ref(app.$InvitationCode);

// watchers
watch(
  () => props.data,
  (message) => {
    if (message.status == app.$Fail) {
      toast.error(message);
    }
    handleEvent(message);
  },
  { deep: true, immediate: true }
);

// event handlers
function start_quiz(e) {
  e.preventDefault();
  emits("startQuiz");
}

// main function
function handleEvent(message) {
  if (message.event == app.$SentInvitaionCode) {
    code.value = message.data.code;
  }
}
</script>

<template>
  <Frame
    v-if="isAdmin"
    page-title="Ready-shady-go"
    :page-message="'you can start quiz by pressing start button'"
  >
    <form @submit="start_quiz">
      <div class="mb-3 pe-3">
        <label for="code" class="form-label">Invitation code</label>
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
  </Frame>
  <Frame v-else :page-title="'Ready-shady-go'" :page-message="data.data">
    <div class="text-center homepage">{{ data.data }}</div>
  </Frame>
</template>
<style scoped>
.homepage {
  margin-top: 100px;
  margin-bottom: 100px;
  font-size: 30px;
}
</style>
