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

const session_id = route.params.session_id;

// event handlers
const handleCustomChange = (isFullScreenEvent) => {
  if (!isFullScreenEvent && myRef.value) {
    toast.error("exit fullscreen mode unexpectedly!!!");
    // handle unexpected behavior
  }
};

// main functions
onMounted(() => {
  // core logic
  if (process.client) {
    try {
      if (socketObject) {
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

const handleQuizEvents = async (message) => {
  if (message.status == app.$Error) {
    return await router.push(
      "/error?status=" + message.status + "&error=" + message.data
    );
  } else if (message.event == app.$TerminateQuiz) {
    invitationCode.value = undefined;
    removeAllUsers();
    return await router.push("/admin/scoreboard?aqi=" + session_id);
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
}

const continueAdmin = () => {
  adminOperationHandler.value.continueAdmin();
}

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
  adminOperationHandler.value.requestSkipTimer();
};

const confirmSkip = (message) => {
  confirmNeeded.title = "Skip Forcefully !!!";
  confirmNeeded.message = message.data;
  confirmNeeded.positive = "skip";
  confirmNeeded.show = true;
};

const handleModal = (confirm) => {
  if (confirm) {
    adminOperationHandler.value.requestSkip(true);
  }
  confirmNeeded.show = false;
};

definePageMeta({
  layout: "empty",
});
// custom class to bind component with
</script>

<template>
  <Playground :full-screen-enabled="myRef" @is-full-screen="handleCustomChange">
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
      @terminate-quiz="terminateQuizHandler"
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
      @ask-skip-timer="askSkipTimer"
    ></QuizScoreSpace>
    <ListJoinUser v-if="currentComponent == 'Waiting'"></ListJoinUser>
    <QuizListUserAnswered
      :data="data"
      v-if="currentComponent == 'Question'"
    ></QuizListUserAnswered>
  </Playground>
</template>
