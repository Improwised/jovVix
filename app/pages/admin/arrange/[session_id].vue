<script setup>
// core dependencies
import { usePush } from "notivue";
import { Ban } from "lucide-vue-next";

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
  action: "skip",
});
const terminating = ref(false);
const lobbyError = ref("");
const currentComponent = ref("Loading");
const adminOperationHandler = ref();
const analysisTab = ref("ranking");
const session_id = route.params.session_id;
const runningQuizJoinUser = ref(0);
const isPauseQuiz = ref(false);
const selectedAnswer = ref(0);
const hostParticipates = ref(false);
const hostPlayedQuiz = ref("");
const hostName = ref("");
const hostUserId = ref("");
const hostUsername = ref("");
const participationPending = ref(false);
const hostNameError = ref("");
const hostNameSuggestions = ref([]);
const restoredAnswer = ref(0);
const questionRenderKey = ref(0);
const restoredSubmissions = ref([]);
const restoredQuestionId = ref("");

const participationUrl = computed(
  () =>
    `${apiUrl}/quiz/host-participation?active_quiz_id=${encodeURIComponent(session_id)}`
);
const loadParticipation = async () => {
  const response = await $fetch(participationUrl.value, {
    credentials: "include",
  });
  const participation = response.data;
  hostParticipates.value = participation.enabled;
  hostPlayedQuiz.value = participation.user_played_quiz || "";
  hostName.value = participation.name || "Host";
  hostUserId.value = participation.user_id || "";
  hostUsername.value = participation.username || "";
  restoredAnswer.value = participation.answer?.[0] || 0;
  restoredSubmissions.value = participation.submitted || [];
  restoredQuestionId.value = participation.question_id || "";
  selectedAnswer.value = restoredAnswer.value;
};
const updateHostParticipation = async (enabled, playerName = "") => {
  participationPending.value = true;
  hostNameError.value = "";
  hostNameSuggestions.value = [];
  try {
    await $fetch(participationUrl.value, {
      method: "PUT",
      credentials: "include",
      body: { enabled, player_name: playerName },
    });
    await loadParticipation();
  } catch (error) {
    const detail = error?.data?.data;
    if (
      enabled &&
      error?.statusCode === 409 &&
      detail &&
      typeof detail === "object"
    ) {
      hostNameError.value =
        detail.message ||
        "This name is already taken. Choose another name for this quiz.";
      hostNameSuggestions.value = detail.suggestions || [];
    } else if (
      enabled &&
      error?.statusCode === 400 &&
      typeof detail === "string"
    ) {
      hostNameError.value = detail;
    } else {
      toast.error(
        typeof detail === "string" ? detail : "Unable to update participation"
      );
    }
  } finally {
    participationPending.value = false;
  }
};
const sendHostAnswer = async (answers) => {
  try {
    const elapsed = Math.max(
      1,
      adminOperationHandler.value.getAnswerResponseTime() || 1
    );
    await $fetch(
      `${apiUrl}/quiz/answer?user_played_quiz=${hostPlayedQuiz.value}&session_id=${session_id}`,
      {
        method: "POST",
        credentials: "include",
        body: {
          id: adminOperationHandler.value.currentQuestion,
          keys: answers,
          response_time: elapsed,
        },
      }
    );
    selectedAnswer.value = answers[0] || 0;
  } catch (error) {
    questionRenderKey.value += 1;
    toast.error(error?.data?.data || "Failed to submit answer");
  }
};

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

const quizState = computed(() =>
  isPauseQuiz.value ? app.$Pause : app.$Running
);

// event handlers
const reloadLobby = () => window.location.reload();

const handleCustomChange = (isFullScreenEvent) => {
  if (!isFullScreenEvent && myRef.value) {
    toast.error("exit fullscreen mode unexpectedly!!!");
    // handle unexpected behavior
  }
};

// main functions
onMounted(async () => {
  if (!import.meta.client) return;
  // Guests with a pending session pick a name first; the socket setup runs
  // after the name modal submission redirects to the real session URL.
  if (showHostNameModal.value) return;
  try {
    await loadParticipation();
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
    if (err?.statusCode === 401) {
      lobbyError.value =
        "Your session has expired. Please sign in and try again.";
    } else if (err?.statusCode === 403 || err?.statusCode === 404) {
      lobbyError.value =
        "This quiz session is unavailable, or you do not have access to host it.";
    } else {
      lobbyError.value =
        "The quiz lobby could not connect. Please check your connection and try again.";
    }
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
    return await goToScoreboardAfterTerminate();
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
      if (
        message.event === app.$GetQuestion &&
        message.data?.id !== restoredQuestionId.value
      ) {
        restoredAnswer.value = 0;
        restoredSubmissions.value = [];
      }
      if (
        message.event === app.$GetQuestion &&
        message.data?.id === restoredQuestionId.value &&
        restoredSubmissions.value.length
      ) {
        resetUsersSubmittedAnswers();
        restoredSubmissions.value.forEach(addUserSubmittedAnswer);
        restoredSubmissions.value = [];
      }
      if (message.event === app.$Counter) {
        selectedAnswer.value = 0;
        restoredAnswer.value = 0;
        restoredSubmissions.value = [];
        restoredQuestionId.value = "";
      }
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
  if (participationPending.value) return;
  adminOperationHandler.value.quizStartRequest();
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
  confirmNeeded.action = "skip";
  confirmNeeded.show = true;
};

// Clean up local session state and route the host to the results view. Shared by
// both the inbound terminate_quiz handler and the host-initiated "End Quiz" button.
const goToScoreboardAfterTerminate = async () => {
  invitationCode.value = undefined;
  removeAllUsers();
  setSession(null);
  // Guest hosts must sign in before accessing the Kratos-protected host scoreboard.
  const isGuestHost = usersStore.getUserData()?.role === "guest-user";
  const scoreboardPath = `/admin/scoreboard?winner_ui=true&aqi=${session_id}`;
  if (isGuestHost && hostParticipates.value && hostPlayedQuiz.value) {
    return await router.push(
      `/join/${encodeURIComponent(hostUsername.value)}/scoreboard?user_played_quiz=${hostPlayedQuiz.value}&host=1`
    );
  }
  if (isGuestHost) {
    return await router.push(
      `/account/login?returnTo=${encodeURIComponent(scoreboardPath)}`
    );
  }
  return await router.push(scoreboardPath);
};

const askEndQuiz = () => {
  confirmNeeded.title = "End this quiz?";
  confirmNeeded.message =
    "This will end the quiz for all connected players and send them to the results. This cannot be undone.";
  confirmNeeded.positive = "End Quiz";
  confirmNeeded.action = "endQuiz";
  confirmNeeded.show = true;
};

const endQuiz = async () => {
  if (terminating.value) return;
  terminating.value = true;
  try {
    await $fetch(`${apiUrl}/quiz/terminate?session_id=${session_id}`, {
      method: "GET",
      credentials: "include",
    });
    await goToScoreboardAfterTerminate();
  } catch (error) {
    console.error("failed to end quiz", error);
    toast.error("Failed to end the quiz. Please try again.");
  } finally {
    terminating.value = false;
  }
};

const handleModal = (confirm) => {
  const action = confirmNeeded.action;
  confirmNeeded.show = false;
  if (!confirm) {
    return;
  }
  if (action === "endQuiz") {
    endQuiz();
  } else {
    adminOperationHandler.value.requestSkip(true);
  }
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
  <main v-if="lobbyError" class="min-h-screen bg-jv-canvas p-6">
    <section
      role="alert"
      class="mx-auto max-w-xl jv-border-rough bg-jv-white p-6 shadow-brutal-sm"
    >
      <h1 class="font-headings text-3xl">Unable to load quiz lobby</h1>
      <p class="mt-3">{{ lobbyError }}</p>
      <button
        type="button"
        class="mt-5 rounded-full bg-jv-coral px-5 py-2 font-bold text-white"
        @click="reloadLobby"
      >
        Try again
      </button>
      <NuxtLink to="/" class="ml-4 font-bold underline">Back to home</NuxtLink>
    </section>
  </main>
  <Playground
    v-else
    :full-screen-enabled="myRef"
    @is-full-screen="handleCustomChange"
  >
    <div
      v-if="currentComponent !== 'Waiting' && currentComponent !== 'Loading'"
      class="code-display grid grid-cols-1 gap-3 px-4 py-3 sm:grid-cols-[1fr_auto_1fr] sm:items-center sm:px-6 md:px-8"
    >
      <div
        v-if="quizTitle"
        class="flex min-w-0 max-w-full flex-col gap-0.5 jv-border-rough bg-jv-white px-3 py-2 shadow-brutal-sm sm:col-start-1 sm:justify-self-start sm:px-4"
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
        class="flex min-w-0 items-center justify-between gap-2 jv-border-rough bg-jv-white px-3 py-2 shadow-brutal-sm sm:col-start-2 sm:justify-self-center sm:justify-start"
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
      <div
        class="flex flex-col gap-3 sm:col-start-3 sm:flex-row sm:items-center sm:justify-self-end"
      >
        <button
          v-if="currentComponent == 'Score'"
          :class="[
            'inline-flex h-11 w-full items-center justify-center rounded-[8px] border-[3px] border-jv-ink px-5 text-[15px] font-black shadow-brutal-sm transition-transform active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:w-fit sm:text-[16px]',
            isPauseQuiz
              ? 'bg-jv-mint text-jv-ink hover:rotate-[1deg]'
              : 'bg-jv-coral text-white hover:rotate-[-1deg]',
          ]"
          @click="handlePauseQuiz"
        >
          {{ isPauseQuiz ? "RESUME" : "PAUSE" }}
        </button>
        <button
          v-if="currentComponent == 'Question' || currentComponent == 'Score'"
          type="button"
          :disabled="terminating"
          class="inline-flex h-11 w-full items-center justify-center gap-2 rounded-[8px] border-[3px] border-jv-ink bg-jv-coral px-5 text-[15px] font-black text-white shadow-brutal-sm transition-transform hover:rotate-[-1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none disabled:cursor-not-allowed disabled:opacity-60 sm:w-fit sm:text-[16px]"
          aria-label="End quiz for all players"
          @click="askEndQuiz"
        >
          <Ban class="size-4" :stroke-width="2.4" />
          <span>{{ terminating ? "Ending..." : "End Quiz" }}</span>
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
      :host-participates="hostParticipates"
      :host-name="hostName"
      :host-user-id="hostUserId"
      :participation-pending="participationPending"
      :host-name-error="hostNameError"
      :host-name-suggestions="hostNameSuggestions"
      @update-host-participation="updateHostParticipation"
      @retry-host-participation="(name) => updateHostParticipation(true, name)"
      @cancel-host-name="
        hostNameError = '';
        hostNameSuggestions = [];
      "
      @start-quiz="startQuiz"
    >
    </QuizWaitingSpace>
    <QuizQuestionSpace
      v-else-if="currentComponent == 'Question'"
      :key="questionRenderKey"
      :data="data"
      :is-admin="true"
      :host-participates="hostParticipates"
      :restored-answer="restoredAnswer"
      :quiz-title="quizTitle"
      @send-answer="sendHostAnswer"
      @ask-skip="askSkip"
    ></QuizQuestionSpace>
    <QuizScoreSpace
      v-else-if="currentComponent == 'Score'"
      :data="data"
      :is-admin="true"
      :selected-answer="selectedAnswer"
      :user-name="hostUsername"
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
