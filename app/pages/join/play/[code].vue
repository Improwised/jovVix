<script setup>
// core dependencies
import { useToast } from "vue-toastification";
import { useRouter } from "nuxt/app";

// custom component
import UserOperation from "~/composables/user_operation.js";
import { useSystemEnv } from "~/composables/envs.js";

// define nuxt configs
const route = useRoute();
const router = useRouter();
const toast = useToast();
const app = useNuxtApp();
useSystemEnv();

// define props and emits
const myRef = ref(false);
const data = ref({});
const currentComponent = ref("Loading");
const userOperationHandler = ref();

const monitorTerminateQuiz = ref(false);

// for notification bars
const showConnectingBar = ref(false);
const showReconnectedBar = ref(false);

// get query params
const username = computed(() => route.query.username);
const firstname = computed(() => route.query.firstname);
const userPlayedQuiz = computed(() => route.query.user_played_quiz);
const sessionId = computed(() => route.query.session_id);

const selectedAnswer = ref(0);

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
      userOperationHandler.value = new UserOperation(
        route.params.code,
        username.value,
        handleQuizEvents,
        handleNetworkEvent,
        handleNetworkEstablished
      );
    } catch (err) {
      toast.info(app.$ReloadRequired);
      console.error(err);
    }
  }
});

const handleQuizEvents = async (message) => {
  if (message.status == app.$Error) {
    console.error(message);
    return await router.push(
      "/error?status=" + message.status + "&error=" + message.data
    );
  } else if (message.event == app.$TerminateQuiz) {
    monitorTerminateQuiz.value = true;
    return await router.push(
      `/join/${username.value}/scoreboard?user_played_quiz=${userPlayedQuiz.value}`
    );
  } else if (message.event == app.$RedirectToAdmin) {
    return await router.push("/admin/arrange/" + message.data.sessionId);
  } else if (
    message.data == app.$InvitationCodeNotFound ||
    message.data == app.$QuizSessionValidationFailed
  ) {
    return await router.push(
      "/join?status=" + message.status + "&error=" + message.data
    );
  } else if (message.data == app.$AdminDisconnected) {
    toast.warning(app.$AdminDisconnectedMessage);
  } else if (
    message.status == app.$Fail &&
    message.event == app.$InvitationCodeValidation
  ) {
    console.error(message);
    return await router.push(
      "/join?status=" + message.status + "&error=" + message.data
    );
  }
  // unauthorized ? -> redirect to login page
  else if (message.status == app.$Fail && message.data == app.$Unauthorized) {
    console.error(message);
    router.push(
      "/account/login?error=" + message.data + "&url=" + route.fullPath
    );
    return;
  } else if (message.data === app.$Unauthorized) {
    toast.error(
      "You are unauthorized to access the resource or Your JWT token is expired"
    );
  } else {
    if (message?.component === "Question") {
      selectedAnswer.value = 0;
    }
    data.value = message;
    currentComponent.value = message.component;
  }
};

function handleNetworkEvent() {
  showConnectingBar.value = true;
}

function hideConnectingBar() {
  showConnectingBar.value = false;
}

function handleNetworkEstablished() {
  showConnectingBar.value = false;
  showReconnectedBar.value = true;

  setTimeout(() => {
    showReconnectedBar.value = false;
  }, 2000);
}

const startQuiz = () => {
  myRef.value = true;
};

const sendAnswer = async (answers) => {
  selectedAnswer.value = 0;
  try {
    const { error } = await userOperationHandler.value.handleSendAnswer(
      answers,
      userPlayedQuiz.value,
      sessionId.value
    );

    if (error) {
      toast.error(error);
      return;
    }
    toast.success(app.$AnswerSubmitted);
    if (answers.length > 0) {
      selectedAnswer.value = answers[0];
    }
  } catch (err) {
    toast.error("An error occurred while submitting the answer.");
  }
};

definePageMeta({
  layout: "empty",
});

onBeforeUnmount(() => {
  if (!monitorTerminateQuiz.value) {
    userOperationHandler.value.endQuiz();
  }
});
</script>

<template>
  <div class="bg-image"></div>
  <div>
    <div v-if="showConnectingBar" class="top-bar-red">
      <div class="doodle">&#128641;</div>
      <div class="text-inside-bar">
        Problem connecting with server, reconnecting...
      </div>
      <button class="close-button" @click="hideConnectingBar">×</button>
    </div>

    <div v-if="showReconnectedBar" class="top-bar-green">
      <div class="text-inside-bar">Reconnected &#128515;</div>
    </div>

    <Playground
      :full-screen-enabled="myRef"
      @is-full-screen="handleCustomChange"
    >
      <UserName :user-name="firstname"></UserName>

      <QuizLoadingSpace
        v-if="currentComponent === 'Loading'"
      ></QuizLoadingSpace>
      <QuizWaitingSpace
        v-else-if="currentComponent === 'Waiting'"
        :data="data"
        :is-admin="false"
        @start-quiz="startQuiz"
      >
      </QuizWaitingSpace>
      <QuizQuestionSpace
        v-else-if="currentComponent === 'Question'"
        :data="data"
        :is-admin="false"
        @send-answer="sendAnswer"
      ></QuizQuestionSpace>
      <QuizScoreSpace
        v-else-if="currentComponent === 'Score'"
        :data="data"
        :user-name="username"
        :is-admin="false"
        :selected-answer="selectedAnswer"
      ></QuizScoreSpace>
    </Playground>
  </div>
</template>

<style scoped>
.top-bar-red,
.top-bar-green {
  background-color: #e41d3b;
  color: #000000;
  padding: 10px 0;
  text-align: center;
  opacity: 1;
  transition: opacity 0.3s ease-in-out;
  position: relative;
}

.top-bar-green {
  background-color: green;
}

.text-inside-bar {
  display: inline-block;
  font-size: 18px;
}

@keyframes doodle-animation {
  0% {
    left: calc(100% + 50px);
  }

  100% {
    left: -50px;
  }
}

.doodle {
  position: absolute;
  animation: doodle-animation 5s linear infinite;
  font-size: 24px;
}

.close-button {
  position: absolute;
  top: 5px;
  right: 10px;
  background: none;
  border: none;
  cursor: pointer;
  color: white;
  font-size: 20px;
}

.close-button:hover {
  color: #ccc;
}

.bg-image {
  background-image: url("@/assets/images/que-web-bg.png");
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
    background-image: url("@/assets/images/Que-mob-bg.png");
  }
}
</style>
