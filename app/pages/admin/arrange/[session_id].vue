<script setup>
// core dependencies
import { useToast } from "vue-toastification";

// custom component
import { useSystemEnv } from "~/composables/envs.js";
import { useRouter } from "nuxt/app";
import AdminOperations from "~~/composables/admin_operation";

import { useInvitationCodeStore } from "~/store/invitationcode";
import { useListUserstore } from "~/store/userlist";
import { useUserThatSubmittedAnswer } from "~/store/userSubmittedAnswer";
import { storeToRefs } from "pinia";
import { useSessionStore } from "~~/store/session";
const sessionStore = useSessionStore();
const { setSession, setLastComponent, getLastComponent } = sessionStore;

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
const toast = useToast();
const app = useNuxtApp();
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
  if (process.client) {
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
  }
});

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
    confirmNeeded.value = {
      show: false,
    };

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
  adminOperationHandler.value.quizStartRequest();
};

const sendAnswer = (answers) => {
  adminOperationHandler.value.handleSendAnswer(answers);
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
});
// custom class to bind component with
</script>

<template>
  <div class="bg-image"></div>
  <Playground :full-screen-enabled="myRef" @is-full-screen="handleCustomChange">
    <div
      v-if="currentComponent !== 'Waiting'"
      class="code-display p-3 d-flex align-items-center justify-content-end"
    >
      <div class="d-flex align-items-center gap-2 me-3">
        <span class="text-muted fw-bold fs-2">Code:</span>
        <span class="fw-bold fs-2">{{ invitationCode }}</span>
      </div>
      <div v-if="currentComponent == 'Score'">
        <button
          v-if="isPauseQuiz"
          class="btn btn-danger"
          @click="handlePauseQuiz"
        >
          START
        </button>
        <button v-else class="btn btn-danger" @click="handlePauseQuiz">
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
    <QuizLoadingSpace v-if="currentComponent == 'Loading'"></QuizLoadingSpace>
    <QuizWaitingSpace
      v-else-if="currentComponent == 'Waiting'"
      :data="data"
      :is-admin="true"
      @start-quiz="startQuiz"
    >
    </QuizWaitingSpace>
    <QuizQuestionSpace
      v-else-if="currentComponent == 'Question'"
      :data="data"
      :is-admin="true"
      @send-answer="sendAnswer"
      @ask-skip="askSkip"
    ></QuizQuestionSpace>
    <QuizScoreSpace
      v-else-if="currentComponent == 'Score'"
      :data="data"
      :is-admin="true"
      :analysis-tab="analysisTab"
      :quiz-state="quizState"
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
