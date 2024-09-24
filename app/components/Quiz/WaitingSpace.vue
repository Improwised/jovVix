<script setup>
// core dependencies
import { useNuxtApp } from "nuxt/app";
import { useToast } from "vue-toastification";

// define nuxt configs
const app = useNuxtApp();
const toast = useToast();

import { useInvitationCodeStore } from "~/store/invitationcode.js";
import { useListUserstore } from "~/store/userlist";
import { storeToRefs } from "pinia";

const invitationCodeStore = useInvitationCodeStore();
const { invitationCode } = storeToRefs(invitationCodeStore);

const listUserStore = useListUserstore();
const { removeAllUsers } = listUserStore;

const startQuiz = ref(false);

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
const emits = defineEmits(["startQuiz", "terminateQuiz"]);

// custom refs
const code = ref(invitationCode.value);

// watchers
watch(
  () => props.data,
  (message) => {
    if (message.status == app.$Fail) {
      toast.error(message.data);
    }
    handleEvent(message);
  },
  { deep: true, immediate: true }
);

// event handlers
function start_quiz(e) {
  e.preventDefault();
  startQuiz.value = true;
  emits("startQuiz");
}

// main function
function handleEvent(message) {
  if (message.event == app.$SentInvitaionCode) {
    code.value = invitationCode.value;
  }
}

let copyBtn;

onMounted(() => {
  copyBtn = document.getElementById("OTP-input-container");
  if (process.client && copyBtn) {
    copyBtn.addEventListener("click", () => {
      copyToClipboard(code.value);
    });
  }
});

onUnmounted(() => {
  if (!startQuiz.value) {
    emits("terminateQuiz");
  }
  if (props.isAdmin) {
    invitationCode.value = null;
    removeAllUsers();
  }
});

function copyToClipboard(text) {
  navigator.clipboard
    .writeText(text)
    .then(() => {
      toast.success("Code copied to clipboard");
    })
    .catch((error) => {
      toast.warning("Error copying to clipboard", error);
    });
}
</script>

<template>
  <Frame
    v-if="isAdmin"
    page-title="Ready Steady Go"
    :page-message="'You Can Start Quiz By Pressing Start Quiz button'"
  >
    <form @submit="start_quiz">
      <div class="mb-3 pe-3">
        <label for="code" class="form-label">Invitation Code</label>
        <div class="d-flex align-items-center justify-content-center gap-2">
          <span id="code" class="code display-2">{{ code }}</span>
          <font-awesome-icon
            id="OTP-input-container"
            icon="fa-solid fa-copy"
            size="xl"
            style="color: #0c6efd"
            class="copy-icon"
            role="button"
          />
        </div>
      </div>
      <button type="submit" class="btn btn-primary btn-lg bg-primary">
        Start Quiz
      </button>
    </form>
  </Frame>
  <Frame v-else page-title="Ready Steady Go" :page-message="data.data">
    <div class="text-center homepage">{{ data.data }}</div>
  </Frame>
</template>
<style scoped>
.homepage {
  margin-top: 100px;
  margin-bottom: 100px;
  font-size: 30px;
}

.code {
  color: #000000;
  padding: 20px;
  letter-spacing: 1rem;
}

@media (max-width: 768px) {
  .code {
    letter-spacing: 0.4rem;
  }
}

@media (max-width: 480px) {
  .code {
    padding: 15px;
    letter-spacing: 0.3rem;
  }
}
</style>
