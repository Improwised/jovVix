<script setup>
// core dependencies
import { useNuxtApp } from "nuxt/app";
import { useToast } from "vue-toastification";
import { useMusicStore } from "~~/store/music";
const { base_url } = useRuntimeConfig().public;
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

onMounted(() => {
  initializeSound();
  const copyBtn = document.getElementById("OTP-input-container");
  const urlCopyBtn = document.getElementById("URL-input-container");
  if (process.client && copyBtn && urlCopyBtn) {
    copyBtn.addEventListener("click", () => {
      usecopyToClipboard(code.value);
    });
    urlCopyBtn.addEventListener("click", () => {
      usecopyToClipboard(`${base_url}/join?code=${code.value}`);
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
  <div
    id="waitingspace"
    class="modal fade"
    tabindex="-1"
    aria-labelledby="exampleModalLabel"
    aria-hidden="true"
  >
    <QuizJoinModal :join-u-r-l="`${base_url}/join`" :code="code" />
  </div>
  <Frame
    v-if="isAdmin"
    page-title="Ready Steady Go"
    :page-message="'You Can Start Quiz By Pressing Start Quiz button'"
    :music-component="true"
  >
    <form @submit="start_quiz">
      <div class="mb-3 pe-3">
        <div class="d-flex justify-content-between">
          <label for="code" class="form-label">Link</label>
          <button
            type="button"
            data-bs-toggle="modal"
            data-bs-target="#waitingspace"
          >
            <font-awesome-icon :icon="['fas', 'expand']" />
          </button>
        </div>
        <div class="d-flex align-items-center justify-content-center gap-2">
          <div class="fs-1 text-dark text-decoration-underline">
            quiz.i8d.in/join
          </div>
          <font-awesome-icon
            id="URL-input-container"
            icon="fa-solid fa-copy"
            size="xl"
            style="color: #0c6efd"
            class="copy-icon"
            role="button"
          />
        </div>
        <div class="divider my-3 text-dark">Invitation Code</div>
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
        <div class="divider my-3">QR</div>
        <div class="d-flex align-items-center justify-content-center">
          <QrCode :scan-u-r-l="`${base_url}/join`" :quiz-code="code" />
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
