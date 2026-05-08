<script setup>
// core dependencies
import { useNuxtApp } from "nuxt/app";
import { useToast } from "vue-toastification";
import { useMusicStore } from "~~/store/music";
import {
  Check,
  Copy,
  Info,
  Keyboard,
  Smile,
  UserRound,
  Users,
} from "lucide-vue-next";
const { baseUrl } = useRuntimeConfig().public;
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
const { listUsers } = storeToRefs(listUserStore);

const startQuiz = ref(false);
const waitingSound = ref(null);
const participantAccentClasses = [
  "bg-jv-yellow",
  "bg-jv-coral text-white",
  "bg-jv-mint",
  "bg-jv-white",
];

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
const participantCount = computed(() => listUsers.value.length);
const joinUrl = computed(() => `${baseUrl}/join`);
const joinUrlWithCode = computed(() => `${joinUrl.value}?code=${code.value}`);

const getParticipantName = (user) =>
  user?.UserName || user?.username || user?.name || "Player";

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
});

const copyToClipBoard = (text) => {
  usecopyToClipboard(text);
};

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
  <main
    v-if="isAdmin"
    class="flex min-h-screen flex-col gap-8 bg-jv-canvas px-4 py-5 text-jv-ink sm:gap-10 sm:px-6 md:px-8 md:py-6"
  >
    <section class="mx-auto flex w-full max-w-[1220px] flex-1 flex-col">
      <header class="mb-7 text-center sm:mb-9">
        <h1
          class="font-headings text-[38px] leading-none text-jv-ink min-[420px]:text-[44px] sm:text-[52px] md:text-[56px]"
        >
          Quiz Lobby
        </h1>
        <p
          class="mx-auto mt-2 max-w-[680px] text-[15px] leading-[1.5] text-jv-muted sm:text-[18px] md:text-[20px]"
        >
          Share the code with players. The quiz will start when you're ready!
        </p>
      </header>

      <div
        class="grid flex-1 gap-6 md:gap-7 xl:grid-cols-[minmax(0,1fr)_minmax(340px,0.95fr)] xl:items-stretch"
      >
        <form
          class="relative flex min-h-0 rotate-[-0.4deg] flex-col overflow-hidden bg-jv-white shadow-brutal-sm jv-border-rough sm:shadow-brutal-lg xl:min-h-[640px]"
          @submit="start_quiz"
        >
          <span
            class="absolute left-1/2 top-[-2px] z-10 h-4 w-12 -translate-x-1/2 rotate-[-1deg] bg-jv-coral"
            aria-hidden="true"
          ></span>

          <div class="bg-jv-yellow px-4 pb-7 pt-6 sm:px-8 sm:pb-9 md:px-10">
            <h2
              class="font-headings text-[30px] leading-tight text-jv-ink min-[420px]:text-[34px] sm:text-[42px] lg:text-[48px]"
            >
              Ready, Steady, Go!
            </h2>

            <div
              class="mt-4 flex flex-col gap-3 jv-border-rough bg-jv-white p-3 shadow-brutal-sm sm:flex-row sm:items-center sm:justify-between sm:p-4"
            >
              <div
                class="flex min-w-0 flex-wrap items-baseline gap-x-2 gap-y-1"
              >
                <p
                  class="shrink-0 font-body text-[12px] font-black uppercase tracking-[0.08em] text-jv-muted"
                >
                  Join at
                </p>
                <p
                  class="min-w-0 break-all font-body text-[16px] font-extrabold leading-[1.35] text-jv-ink min-[420px]:text-[17px] sm:text-[20px] md:text-[22px]"
                >
                  {{ joinUrl.replace(/^https?:\/\//, "") }}
                </p>
              </div>
              <button
                id="URL-input-container"
                type="button"
                class="inline-flex h-11 w-full shrink-0 rotate-[-1deg] items-center justify-center gap-2 rounded-[999px] border-[3px] border-jv-ink bg-jv-coral px-5 font-body text-[16px] font-bold text-white shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-12 sm:w-fit"
                @click="copyToClipBoard(joinUrlWithCode)"
              >
                <Copy class="size-4" :stroke-width="2.5" />
                <span>Copy</span>
              </button>
            </div>
          </div>

          <div class="flex flex-1 flex-col px-4 py-5 sm:px-8 sm:py-6 md:px-10">
            <div
              class="-mt-11 flex flex-col gap-3 jv-border-rough bg-jv-white p-3 shadow-brutal-sm sm:-mt-14 sm:flex-row sm:items-center sm:justify-between sm:p-4 sm:shadow-brutal"
            >
              <div class="flex min-w-0 items-center gap-3 sm:gap-4">
                <span
                  class="grid size-10 shrink-0 place-items-center border-[2px] border-jv-ink bg-jv-mint text-jv-ink sm:size-11"
                >
                  <Keyboard class="size-5" :stroke-width="2.2" />
                </span>
                <div
                  class="flex min-w-0 flex-wrap items-baseline gap-x-2 gap-y-1"
                >
                  <p
                    class="shrink-0 font-body text-[12px] font-black uppercase tracking-[0.08em] text-jv-muted"
                  >
                    Quiz code
                  </p>
                  <h3
                    class="code min-w-0 break-all font-feature text-[22px] font-black leading-none text-jv-ink min-[420px]:text-[24px] sm:text-[26px]"
                  >
                    {{ code }}
                  </h3>
                </div>
              </div>
              <button
                id="OTP-input-container"
                type="button"
                class="inline-flex h-11 w-full shrink-0 rotate-[-1deg] items-center justify-center gap-2 rounded-[999px] border-[3px] border-jv-ink bg-jv-coral px-5 font-body text-[16px] font-bold text-white shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-12 sm:w-fit"
                @click="copyToClipBoard(code)"
              >
                <Copy class="size-4" :stroke-width="2.5" />
                <span>Copy</span>
              </button>
            </div>

            <div
              class="my-5 flex items-center justify-center gap-3 text-center font-body text-[11px] font-black uppercase tracking-[0.08em] text-jv-muted min-[420px]:text-[12px] sm:my-6 sm:gap-4 sm:text-[13px]"
            >
              <span
                class="h-px min-w-8 flex-1 border-t-2 border-dashed border-jv-ink/40 sm:max-w-24"
              ></span>
              <span class="shrink-0">Or scan QR code</span>
              <span
                class="h-px min-w-8 flex-1 border-t-2 border-dashed border-jv-ink/40 sm:max-w-24"
              ></span>
            </div>

            <div class="flex justify-center">
              <div
                class="qr-card grid size-[176px] place-items-center bg-jv-white p-3 shadow-brutal-sm jv-border-rough min-[420px]:size-[196px] sm:size-[220px] sm:p-4 md:size-[240px]"
              >
                <QrCode :scan-u-r-l="joinUrl" :quiz-code="code" :size="200" />
              </div>
            </div>

            <div class="mt-auto pt-5 sm:pt-6">
              <button
                type="submit"
                class="mx-auto flex h-[52px] w-full max-w-[390px] rotate-[-1deg] items-center justify-center rounded-[999px] border-[3px] border-jv-ink bg-jv-mint px-5 font-headings text-[16px] text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-16 sm:px-6 sm:text-[18px]"
              >
                Start Quiz
              </button>
              <p
                class="mt-4 flex items-center justify-center gap-2 text-center font-body text-[12px] leading-[1.4] text-jv-muted sm:text-[13px]"
              >
                <Info class="size-4" :stroke-width="2.3" />
                <span>Host can start the quiz at any time</span>
              </p>
            </div>
          </div>
        </form>

        <aside
          class="relative flex min-h-0 rotate-[0.4deg] flex-col bg-jv-white px-4 py-6 shadow-brutal-sm jv-border-rough sm:px-8 sm:py-8 sm:shadow-brutal-lg md:px-10 xl:min-h-[640px]"
        >
          <span
            class="absolute left-1/2 top-[-10px] h-4 w-12 -translate-x-1/2 rotate-[1deg] bg-jv-salmon"
            aria-hidden="true"
          ></span>

          <div
            class="flex items-start justify-between gap-3 sm:items-center sm:gap-4"
          >
            <div class="flex min-w-0 items-center gap-3 sm:gap-4">
              <span
                class="grid size-10 shrink-0 rotate-[-3deg] place-items-center jv-border-rough bg-jv-mint sm:size-12"
              >
                <Users class="size-5 sm:size-6" :stroke-width="2.4" />
              </span>
              <h2
                class="min-w-0 break-words font-headings text-[28px] leading-tight text-jv-ink min-[420px]:text-[32px] sm:text-[42px] xl:text-[48px]"
              >
                Participants
              </h2>
            </div>
            <span
              class="grid size-11 shrink-0 rotate-[3deg] place-items-center border-[3px] border-jv-ink bg-jv-yellow font-feature text-[24px] font-black shadow-brutal-sm sm:size-14 sm:text-[30px]"
            >
              {{ participantCount }}
            </span>
          </div>

          <div
            class="my-5 border-t-2 border-dashed border-jv-ink/20 sm:my-6"
          ></div>

          <div
            v-if="participantCount"
            class="flex flex-col gap-3 overflow-y-auto pr-1 sm:gap-4 xl:max-h-[430px]"
          >
            <div
              v-for="(user, index) in listUsers"
              :key="user.UserId || user.UserName || index"
              class="flex min-h-14 items-center justify-between gap-3 jv-border-rough bg-jv-white px-3 py-2 shadow-[2px_2px_0_#2D2D2D] sm:min-h-16"
            >
              <div class="flex min-w-0 items-center gap-3">
                <span
                  class="grid size-9 shrink-0 rotate-[-2deg] place-items-center border-[2px] border-jv-ink sm:size-10"
                  :class="
                    participantAccentClasses[
                      index % participantAccentClasses.length
                    ]
                  "
                >
                  <Smile
                    v-if="index % 3 === 0"
                    class="size-5"
                    :stroke-width="2.2"
                  />
                  <UserRound v-else class="size-5" :stroke-width="2.2" />
                </span>
                <span
                  class="truncate font-body text-[16px] font-semibold text-jv-ink sm:text-[20px]"
                >
                  {{ getParticipantName(user) }}
                </span>
              </div>
              <span
                class="grid size-8 shrink-0 rotate-[2deg] place-items-center border-[2px] border-jv-ink bg-jv-mint sm:size-9"
                aria-label="Joined"
              >
                <Check class="size-5" :stroke-width="2.4" />
              </span>
            </div>
          </div>

          <div
            v-else
            class="grid min-h-[220px] flex-1 place-items-center border-2 border-dashed border-jv-ink/20 bg-jv-canvas/50 p-6 text-center sm:min-h-[300px] sm:p-8"
          >
            <div>
              <div
                class="mx-auto grid size-14 place-items-center rounded-full border-[3px] border-jv-ink bg-jv-mint shadow-brutal-sm sm:size-16"
              >
                <Users class="size-7 sm:size-8" :stroke-width="2.4" />
              </div>
              <p
                class="mt-5 font-body text-[18px] font-bold text-jv-muted sm:text-[22px]"
              >
                Waiting for participants...
              </p>
            </div>
          </div>

          <div class="mt-auto pt-6 sm:pt-8">
            <div class="border-t-2 border-jv-ink/15 pt-5 text-center sm:pt-6">
              <div class="mb-3 flex justify-center gap-2 sm:mb-4">
                <span
                  class="size-2.5 rounded-full bg-jv-coral sm:size-3"
                ></span>
                <span
                  class="size-2.5 rounded-full bg-jv-yellow sm:size-3"
                ></span>
                <span class="size-2.5 rounded-full bg-jv-mint sm:size-3"></span>
              </div>
              <p
                class="font-body text-[18px] text-jv-muted sm:text-[22px] md:text-[24px]"
              >
                Waiting for more players...
              </p>
            </div>
          </div>
        </aside>
      </div>

      <footer
        class="mt-8 border-t-2 border-dashed border-jv-ink/15 pt-5 text-center font-body text-[13px] text-jv-muted sm:mt-10 sm:pt-6 sm:text-[16px]"
      >
        &copy; 2026 Jovvix Platform &bull; Let's Play!
      </footer>
    </section>
  </main>
  <Frame
    v-else
    :music-component="true"
    page-title="Ready Steady Go"
    :page-message="data.data"
  >
    <div class="my-24 text-center font-body text-[30px] text-jv-ink">
      {{ data.data }}
    </div>
  </Frame>
</template>
<style scoped>
.code {
  letter-spacing: 0.4rem;
}

@media (max-width: 768px) {
  .code {
    letter-spacing: 0.22rem;
  }
}

.qr-card :deep(svg) {
  width: 100%;
  height: 100%;
}
</style>
