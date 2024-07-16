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
  if (props.isAdmin) {
    invitationCode.value = null
    removeAllUsers()
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
            min-height="20"
            type="number"
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

@media only screen and (max-width: 991px) {
  .copy-icon {
    right: -20px !important;
  }
}
</style>
