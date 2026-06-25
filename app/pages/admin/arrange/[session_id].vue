<script setup>
// core dependencies
import { usePush } from "notivue";

// custom component
import { useSystemEnv } from "~/composables/envs.js";
import { useRouter } from "nuxt/app";
import AdminOperations from "~~/composables/admin_operation";

import { useInvitationCodeStore } from "~/store/invitationcode";
import { useListUserstore } from "~/store/userlist";
import { useUserThatSubmittedAnswer } from "~/store/userSubmittedAnswer";
import { storeToRefs } from "pinia";
import { useSessionStore } from "~~/store/session";
import { useUsersStore } from "~~/store/users";
const sessionStore = useSessionStore();
const {
  setSession,
  setLastComponent,
  getLastComponent,
  setActiveQuizTitle,
  getActiveQuizTitle,
} = sessionStore;
const { activeQuizTitle } = storeToRefs(sessionStore);
const quizTitle = computed(
  () => activeQuizTitle.value || getActiveQuizTitle() || ""
);
const usersStore = useUsersStore();

const invitationCodeStore = useInvitationCodeStore();
const { invitationCode } = storeToRefs(invitationCodeStore);

const listUserStore = useListUserstore();
const { addUser, removeAllUsers } = listUserStore;

const usersThatSubmittedAnswer = useUserThatSubmittedAnswer();
const { addUserSubmittedAnswer, resetUsersSubmittedAnswers } =
  usersThatSubmittedAnswer;

// define nuxt configs
const route = useRoute();
const router = useRouter();
const toast = usePush();
const app = useNuxtApp();
const { apiUrl } = useRuntimeConfig().public;
useSystemEnv();

// define props and emits
const myRef = ref(false);
const data = ref({});
const confirmNeeded = reactive({
  show: false,
  title: "title",
  message: "message",
  positive: "save",
  negative: "cancel",
});
const currentComponent = ref("Loading");
const adminOperationHandler = ref();
const analysisTab = ref("ranking");
const session_id = route.params.session_id;
const runningQuizJoinUser = ref(0);
const isPauseQuiz = ref(false);
const selectedAnswer = ref(0);

// Public-quiz guests land here with session_id="new" and a quiz_id query — they
// need to pick a display name before the guest user + public session are created.
const PENDING_SENTINEL = "new";
const pendingQuizId = computed(() =>
  session_id === PENDING_SENTINEL ? route.query.quiz_id || "" : ""
);
const showHostNameModal = ref(
  session_id === PENDING_SENTINEL && route.query.public === "1"
);
const hostNameSubmitting = ref(false);

// Public-quiz host-also-plays support.
// `public=1` is set by the homepage when a visitor starts a public quiz.
const isPublicPlay = computed(() => route.query.public === "1");
const canPlay = ref(false);
const hostUserPlayedQuiz = ref(null);
const hostPlayedQuizRequested = ref(false);

// When hosting a public quiz, try to register the host as a player too. The API
// allows this only for public quizzes started by someone other than the creator;
// a creator hosting their own quiz gets a 403 and simply stays host-only.
const tryEnableHostPlay = async (code) => {
  if (!isPublicPlay.value || hostPlayedQuizRequested.value || !code) return;
  hostPlayedQuizRequested.value = true;
  try {
    const res = await $fetch(`${apiUrl}/user_played_quizes/${code}`, {
      method: "POST",
      credentials: "include",
    });
    hostUserPlayedQuiz.value = res?.data?.user_played_quiz || null;
    canPlay.value = !!hostUserPlayedQuiz.value;
    if (res?.data?.quiz_title) {
      setActiveQuizTitle(res.data.quiz_title);
    }
  } catch (error) {
    // 403 => host is the quiz creator and may only host, not play. Expected; stay host-only.
    toast.warning(
      "Host is the quiz creator and may only host the public quiz, not play."
    );
    canPlay.value = false;
  }
};
const quizState = computed(() =>
  isPauseQuiz.value ? app.$Pause : app.$Running
);

// event handlers
const handleCustomChange = (isFullScreenEvent) => {
  if (!isFullScreenEvent && myRef.value) {
    toast.error("exit fullscreen mode unexpectedly!!!");
    // handle unexpected behavior
  }
};

// main functions
onMounted(() => {
  if (!process.client) return;
  // Guests with a pending session pick a name first; the socket setup runs
  // after the name modal submission redirects to the real session URL.
  if (showHostNameModal.value) return;
  try {
    const lastRenderedComponent = getLastComponent();
    if (socketObject && lastRenderedComponent !== "Waiting") {
      adminOperationHandler.value = new AdminOperations(
        session_id,
        handleQuizEvents,
        handleNetworkEvent,
        confirmSkip
      );
      continueAdmin();
    } else {
      adminOperationHandler.value = new AdminOperations(
        session_id,
        handleQuizEvents,
        handleNetworkEvent,
        confirmSkip
      );
      connectAdmin();
    }
  } catch (err) {
    console.error(err);
    toast.info(app.$ReloadRequired);
  }
});

const handleHostNameSubmit = async ({ name, avatarName }) => {
  if (hostNameSubmitting.value) return;
  if (!pendingQuizId.value) {
    toast.error("Missing quiz to host. Please pick a quiz again.");
    await router.replace("/");
    return;
  }
  hostNameSubmitting.value = true;
  try {
    const userRes = await $fetch(
      `${apiUrl}/user/${encodeURIComponent(
        name
      )}?avatar_name=${encodeURIComponent(avatarName)}`,
      {
        method: "POST",
        credentials: "include",
        headers: { Accept: "application/json" },
      }
    );
    const guest = userRes?.data;
    if (guest) {
      usersStore.setUserData({
        role: "guest-user",
        avatar: guest.img_key || avatarName,
        firstname: guest.first_name || name,
        username: guest.username || name,
      });
    }

    const sessionRes = await $fetch(
      `${apiUrl}/quizzes/${pendingQuizId.value}/public_session`,
      {
        method: "POST",
        credentials: "include",
        headers: { Accept: "application/json" },
      }
    );
    const newSessionId = sessionRes?.data;
    if (!newSessionId) {
      toast.error("Error while starting quiz.");
      return;
    }

    removeAllUsers();
    setSocketObject(null);
    setSession(newSessionId);
    showHostNameModal.value = false;
    await router.replace(`/admin/arrange/${newSessionId}?public=1`);
  } catch (error) {
    toast.error(
      error?.data?.message || error?.message || "Error while starting quiz."
    );
  } finally {
    hostNameSubmitting.value = false;
  }
};

onUnmounted(() => {
  setLastComponent(currentComponent.value);
});

const handleQuizEvents = async (message) => {
  if (message.status == app.$Error) {
    return await router.push(
      "/error?status=" + message.status + "&error=" + message.data
    );
  } else if (message.event == app.$TerminateQuiz) {
    invitationCode.value = undefined;
    removeAllUsers();
    setSession(null);
    // A host who also played sees their own player scoreboard (the admin scoreboard
    // endpoint is Kratos-only and would 401 for guests). Host-only/creators keep the
    // admin scoreboard + analytics they can revisit.
    if (canPlay.value && hostUserPlayedQuiz.value) {
      const playerName = usersStore.getUserData()?.username || "player";
      return await router.push(
        `/join/${encodeURIComponent(playerName)}/scoreboard?user_played_quiz=${
          hostUserPlayedQuiz.value
        }`
      );
    }
    return await router.push(
      "/admin/scoreboard?winner_ui=true&aqi=" + session_id
    );
  } else if (message.event == app.$RedirectToAdmin) {
    return await router.push("/admin/arrange/" + message.data.sessionId);
  } else if (
    message.data == app.$InvitationCodeNotFound ||
    message.data == app.$QuizSessionValidationFailed ||
    message.data == app.$SessionWasCompleted
  ) {
    return await router.push(
      "/admin/arrange?status=" + message.status + "&error=" + message.data
    );
  } else if (
    message.event === app.$EventAnswerSubmittedByUser &&
    message.action === app.$ActionAnserSubmittedByUser
  ) {
    addUserSubmittedAnswer(message.data);
  } else if (message.component === "Running") {
    runningQuizJoinUser.value = message.data;
  } else {
    if (message.component != "Question") {
      resetUsersSubmittedAnswers();
    }
    if (
      message.status == app.$Fail &&
      message.event == app.$InvitationCodeValidation
    ) {
      return await router.push(
        "/join?status=" + message.status + "&error=" + message.data
      );
    }
    // unauthorized ? -> redirect to login page
    if (message.status == app.$Fail && message.data == app.$Unauthorized) {
      router.push(
        "/account/login?error=" + message.data + "&url=" + route.fullPath
      );
      return;
    }
    data.value = message;
    currentComponent.value = message.component;
    if (message.component === "Question") {
      selectedAnswer.value = 0;
    }
    confirmNeeded.value = {
      show: false,
    };

    // Capture the quiz title whenever it shows up so the post-quiz
    // scoreboard view can display it without an extra API call.
    const inboundTitle =
      message?.data?.quizTitle ||
      message?.data?.title ||
      message?.data?.data?.quizTitle;
    if (inboundTitle) {
      setActiveQuizTitle(inboundTitle);
    }

    if (currentComponent.value == "Waiting") {
      if (
        invitationCode.value != undefined &&
        message.data != "no player found"
      ) {
        addUser(message.data);
      }
      if (message.data.code !== undefined) {
        invitationCode.value = message.data.code;
        // Code is now known — register the host as a player if this is a public session.
        tryEnableHostPlay(message.data.code);
      }
    }
  }
};

const connectAdmin = () => {
  adminOperationHandler.value.connectAdmin();
};

const continueAdmin = () => {
  adminOperationHandler.value.continueAdmin();
};

function handleNetworkEvent(message) {
  toast.warning(message + ", please reload the page");
}

const startQuiz = () => {
  adminOperationHandler.value.quizStartRequest();
};

const sendAnswer = async (answers) => {
  if (!canPlay.value || !hostUserPlayedQuiz.value) return;
  selectedAnswer.value = 0;
  const { error } = await adminOperationHandler.value.handleSendAnswer(
    answers,
    hostUserPlayedQuiz.value,
    session_id
  );
  if (error) {
    toast.error(error);
    return;
  }
  if (answers.length > 0) {
    selectedAnswer.value = answers[0];
  }
};

const askSkip = () => {
  adminOperationHandler.value.requestSkip(false);
};

// askFor20SecTimerToSkip
const askSkipTimer = () => {
  isPauseQuiz.value = false;
  adminOperationHandler.value.requestSkipTimer();
};

const handlePauseQuiz = () => {
  isPauseQuiz.value = !isPauseQuiz.value;
  adminOperationHandler.value.requestPauseQuiz(isPauseQuiz.value);
};

const confirmSkip = (message) => {
  confirmNeeded.title = "Skip Forcefully !!!";
  confirmNeeded.message = message.data;
  confirmNeeded.positive = "Skip";
  confirmNeeded.show = true;
};

const handleModal = (confirm) => {
  if (confirm) {
    adminOperationHandler.value.requestSkip(true);
  }
  confirmNeeded.show = false;
};

const handleAnalysisTabChange = (tab) => (analysisTab.value = tab);

definePageMeta({
  layout: "empty",
  hideSidebar: true,
  // Public-quiz guests transition from /admin/arrange/new?... to /admin/arrange/<sessionId>?...
  // — keying on fullPath forces a fresh mount so the captured `session_id` const picks up
  // the real session id and the admin socket connects against it.
  key: (route) => route.fullPath,
});

useSeoMeta({
  title: "Quiz Session - jovVix",
  description: "Configure and launch your live quiz session on jovVix.",
  robots: "noindex, nofollow",
});
// custom class to bind component with
</script>

<template>
  <div class="bg-image"></div>
  <QuizHostNameModal
    v-if="showHostNameModal"
    :submitting="hostNameSubmitting"
    @submit="handleHostNameSubmit"
  />
  <Playground :full-screen-enabled="myRef" @is-full-screen="handleCustomChange">
    <div
      v-if="currentComponent !== 'Waiting' && currentComponent !== 'Loading'"
      class="code-display flex flex-col gap-3 px-4 py-3 sm:flex-row sm:items-center sm:justify-between sm:px-6 md:px-8"
    >
      <div
        v-if="quizTitle"
        class="flex min-w-0 max-w-full flex-col gap-0.5 jv-border-rough bg-jv-white px-3 py-2 shadow-brutal-sm sm:px-4"
        :title="quizTitle"
      >
        <span
          class="font-body text-[10px] font-black uppercase tracking-[0.14em] text-jv-muted sm:text-[11px]"
        >
          Now hosting
        </span>
        <span
          class="min-w-0 truncate font-headings text-[18px] leading-tight text-jv-ink sm:text-[22px]"
        >
          {{ quizTitle }}
        </span>
      </div>
      <div
        class="flex min-w-0 items-center justify-between gap-2 jv-border-rough bg-jv-white px-3 py-2 shadow-brutal-sm sm:justify-start"
      >
        <span class="text-[18px] font-bold text-jv-muted sm:text-[22px]">
          Code:
        </span>
        <span
          class="min-w-0 break-all font-feature text-[22px] font-black text-jv-coral sm:text-[28px]"
        >
          {{ invitationCode }}
        </span>
      </div>
      <div v-if="currentComponent == 'Score'" class="sm:ml-2">
        <button
          v-if="isPauseQuiz"
          class="inline-flex h-11 w-full items-center justify-center rounded-[8px] border-[3px] border-jv-ink bg-jv-mint px-5 text-[15px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:w-fit sm:text-[16px]"
          @click="handlePauseQuiz"
        >
          START
        </button>
        <button
          v-else
          class="inline-flex h-11 w-full items-center justify-center rounded-[8px] border-[3px] border-jv-ink bg-jv-coral px-5 text-[15px] font-black text-white shadow-brutal-sm transition-transform hover:rotate-[-1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:w-fit sm:text-[16px]"
          @click="handlePauseQuiz"
        >
          PAUSE
        </button>
      </div>
    </div>
    <UtilsConfirmModal
      v-if="confirmNeeded.show"
      :modal-title="confirmNeeded.title"
      :modal-message="confirmNeeded.message"
      :model-positive-message="confirmNeeded.positive"
      @confirm-message="(c) => handleModal(c)"
    ></UtilsConfirmModal>
    <QuizWaitingSpaceSkeleton
      v-if="currentComponent == 'Loading'"
    ></QuizWaitingSpaceSkeleton>
    <QuizWaitingSpace
      v-else-if="currentComponent == 'Waiting'"
      :data="data"
      :is-admin="true"
      :quiz-title-override="quizTitle"
      @start-quiz="startQuiz"
    >
    </QuizWaitingSpace>
    <QuizQuestionSpace
      v-else-if="currentComponent == 'Question'"
      :data="data"
      :is-admin="true"
      :can-play="canPlay"
      :quiz-title="quizTitle"
      @send-answer="sendAnswer"
      @ask-skip="askSkip"
    ></QuizQuestionSpace>
    <QuizScoreSpace
      v-else-if="currentComponent == 'Score'"
      :data="data"
      :is-admin="true"
      :selected-answer="selectedAnswer"
      :analysis-tab="analysisTab"
      :quiz-state="quizState"
      :quiz-title="quizTitle"
      @change-analysis-tab="handleAnalysisTabChange"
      @ask-skip-timer="askSkipTimer"
    ></QuizScoreSpace>
    <QuizListUserAnswered
      v-if="currentComponent == 'Question' && data?.event !== '5_sec_counter'"
      :data="data"
      :running-quiz-join-user="runningQuizJoinUser"
      @auto-skip="askSkip"
    ></QuizListUserAnswered>
  </Playground>
</template>

<style scoped>
.bg-image {
  background-image: url("@/assets/images/que-web-bg.webp");
  position: fixed;
  right: 0;
  bottom: 0;
  min-width: 100%;
  min-height: 100%;
  width: 100%;
  height: auto;
  z-index: -1;
  opacity: 0.2;
}

@media (max-width: 576px) {
  .bg-image {
    background-image: url("@/assets/images/Que-mob-bg.webp");
  }
}
</style>
