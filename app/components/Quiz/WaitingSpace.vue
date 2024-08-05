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
        <div
          class="d-flex justify-content-start justify-content-md-center align-items-center position-relative"
        >
          <v-otp-input
            v-model="code"
            max-width="500"
            min-height="90"
            type="number"
            class="large-otp-input"
            disabled
          ></v-otp-input>
          <font-awesome-icon
            id="OTP-input-container"
            icon="fa-solid fa-copy"
            size="xl"
            style="color: #0c6efd"
            class="position-absolute end-0 copy-icon"
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
.copy-icon:hover {
  cursor: pointer;
}

.large-otp-input :deep(input) {
  width: 70px; /* Adjust the width as needed */
  height: 70px; /* Adjust the height as needed */
  font-size: 62px; /* Adjust the font size as needed */

  color: rgb(85, 73, 73) !important;
  text-align: center; /* Center the text horizontally */
  line-height: 70px; /* Vertically center the text */
  padding: 0; /* Ensure no extra padding */
  box-sizing: border-box; /* Ensure padding and border are included in width and height */
  display: flex; /* Flexbox to center the content */
  justify-content: center; /* Center the content horizontally */
  align-items: center; /* Center the content vertically */
  font-family: "Pacifico", cursive;
}

.large-otp-input :deep(input::placeholder) {
  font-size: 72px; /* Match font size with input text */
  color: rgb(12, 11, 11) !important;
  text-align: center; /* Center the placeholder text */
  line-height: 70px; /* Vertically center the placeholder text */
  opacity: 1; /* Ensure placeholder visibility */
  font-family: "Pacifico", cursive;
}

@media only screen and (max-width: 991px) {
  .copy-icon {
    right: -20px !important;
  }
}
</style>
