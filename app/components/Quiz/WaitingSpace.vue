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
  <div v-if="isAdmin" class="container-fluid mt-2">
    <div class="flex-grow-1 text-center">
      <h1 class="join-page-title">Ready Steady Go</h1>
      <h6>You Can Start Quiz By Pressing Start Quiz button</h6>
    </div>
    <div class="row">
      <div class="col-md-7">
        <form @submit="start_quiz">
          <div class="mb-3 pe-3">
            <div class="divider my-3 text-dark">Link</div>
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
            <div
              class="d-flex align-items-center justify-content-center qr-scale-down"
            >
              <QrCode
                :scan-u-r-l="`${base_url}/join`"
                :quiz-code="code"
                :size="450"
              />
            </div>
          </div>
          <div class="d-flex justify-content-center align-items-center">
            <button type="submit" class="btn btn-primary btn-lg bg-primary">
              Start Quiz
            </button>
          </div>
        </form>
      </div>
      <div class="col-md-5">
        <ListJoinUser></ListJoinUser>
      </div>
    </div>
  </div>
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

.join-page-title {
  color: #663399;
}

@media (max-width: 768px) {
  .qr-scale-down {
    transform: scale(0.8);
  }
}
</style>
