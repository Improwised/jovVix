<script setup>
// core dependencies
import { useNuxtApp } from "nuxt/app";
import { useToast } from "vue-toastification";
import { useMusicStore } from "~~/store/music";
const { kratos_url } = useRuntimeConfig().public;
const musicStore = useMusicStore();
const { getMusic } = musicStore;

const music = computed(() => {
  return getMusic();
});

// define nuxt configs
const app = useNuxtApp();
const toast = useToast();

import { useInvitationCodeStore } from "~/store/invitationcode.js";
import { useListUserstore } from "~/store/userlist";
import { storeToRefs } from "pinia";
import usecopyToClipboard from "~~/composables/copy_to_clipboard";

const invitationCodeStore = useInvitationCodeStore();
const { invitationCode } = storeToRefs(invitationCodeStore);

const listUserStore = useListUserstore();
const { removeAllUsers } = listUserStore;

const startQuiz = ref(false);
const waitingSound = ref(null);

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
  initializeSound();
  copyBtn = document.getElementById("OTP-input-container");
  if (process.client && copyBtn) {
    copyBtn.addEventListener("click", () => {
      usecopyToClipboard(code.value);
    });
  }
});

function initializeSound() {
  if (process.client) {
    waitingSound.value = new Audio("/music/waiting_area_music.mp3");
    if (music.value) {
      waitingSound.value.play();
      waitingSound.value.loop = true;
    }
  }
}

onUnmounted(() => {
  if (waitingSound.value) {
    waitingSound.value.pause();
    waitingSound.value = null;
  }
  if (!startQuiz.value) {
    emits("terminateQuiz");
  }
  if (props.isAdmin) {
    invitationCode.value = null;
    removeAllUsers();
  }
});

watch(
  music,
  (newValue) => {
    if (!newValue && waitingSound.value) {
      waitingSound.value.pause();
    } else if (newValue && waitingSound.value) {
      waitingSound.value.play();
    }
  },
  { deep: true, immediate: true }
);
</script>

<template>
  <Frame
    v-if="isAdmin"
    page-title="Ready Steady Go"
    :page-message="'You Can Start Quiz By Pressing Start Quiz button'"
    :music-component="true"
  >
    <form @submit="start_quiz">
      <div class="mb-3 pe-3">
        <label for="code" class="form-label">Invitation Code</label>
        <div class="d-flex align-items-center justify-content-center gap-2">
          <h2 class="display-4 code">{{ code }}</h2>
          <font-awesome-icon
            id="OTP-input-container"
            icon="fa-solid fa-copy"
            size="xl"
            style="color: #0c6efd"
            class="copy-icon"
            role="button"
          />
        </div>
        <div class="divider mb-5">OR</div>
        <div class="d-flex align-items-center justify-content-center">
          <QrCode :scan-u-r-l="`${kratos_url}/join`" :quiz-code="code" />
        </div>
      </div>
      <button type="submit" class="btn btn-primary btn-lg bg-primary">
        Start Quiz
      </button>
    </form>
  </Frame>
  <Frame
    v-else
    :music-component="true"
    page-title="Ready Steady Go"
    :page-message="data.data"
  >
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
  letter-spacing: 0.5rem;
}
</style>
