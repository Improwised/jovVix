<script setup>
// core dependencies
import { usePush } from "notivue";
import { useRouter } from "nuxt/app";
import { storeToRefs } from "pinia";

// custom component
import UserOperation from "~/composables/user_operation.js";
import { useSystemEnv } from "~/composables/envs.js";
import { useSessionStore } from "~~/store/session";

// define nuxt configs
const route = useRoute();
const router = useRouter();
const toast = usePush();
const app = useNuxtApp();
useSystemEnv();

const sessionStore = useSessionStore();
const { setActiveQuizTitle } = sessionStore;
const { activeQuizTitle } = storeToRefs(sessionStore);

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
const quizState = ref(app.$Running);

const incomingQuizTitle = (route.query.quiz_title || "").toString().trim();
if (incomingQuizTitle) {
  setActiveQuizTitle(incomingQuizTitle);
}
const quizTitle = computed(
  () => activeQuizTitle.value || incomingQuizTitle || ""
);
const showQuizTitleBar = computed(
  () => !!quizTitle.value && currentComponent.value !== "Waiting"
);

// event handlers
const handleCustomChange = (isFullScreenEvent) => {
  if (!isFullScreenEvent && myRef.value) {
    toast.error("exit fullscreen mode unexpectedly!!!");
    // handle unexpected behavior
  }
};

// Safety net: server closed the socket before we entered the lobby (typically
// a back-nav into a session that has already wrapped up). Bounce to /join so
// the banner there can explain what happened instead of leaving the player
// stuck on the "Connecting to lobby…" skeleton.
const handleSessionUnavailable = () => {
  if (monitorTerminateQuiz.value) return;
  if (currentComponent.value !== "Loading") return;
  router.replace(
    "/join?status=fail&error=" + encodeURIComponent(app.$SessionWasCompleted)
  );
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
        handleNetworkEstablished,
        handleSessionUnavailable
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
    toast.info(app.$HostEndedQuizMessage);
    return await router.push(
      `/join/${username.value}/scoreboard?user_played_quiz=${userPlayedQuiz.value}`
    );
  } else if (message.event == app.$RedirectToAdmin) {
    return await router.push("/admin/arrange/" + message.data.sessionId);
  } else if (
    message.data == app.$InvitationCodeNotFound ||
    message.data == app.$QuizSessionValidationFailed ||
    message.data == app.$SessionWasCompleted
  ) {
    // Players hit this when they back-navigate from the scoreboard into a quiz
    // whose session already terminated — bounce them to /join so the banner can
    // explain what happened instead of leaving them on a stalled play screen.
    return await router.push(
      "/join?status=" + message.status + "&error=" + message.data
    );
  } else if (message.data == app.$AdminDisconnected) {
    toast.warning(app.$AdminDisconnectedMessage);
  } else if (message.data == app.$PauseQuiz) {
    quizState.value = app.$Pause;
    toast.info(app.$PauseQuizMessage);
  } else if (message.data == app.$ResumeQuiz) {
    quizState.value = app.$Running;
    toast.success(app.$ResumeQuizMessage);
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
      quizState.value = app.$Running;
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
  hideSidebar: true,
});

useSeoMeta({
  title: "Play Live Quiz - jovVix",
  description:
    "Answer questions in real time and climb the live scoreboard in this jovVix multiplayer quiz game.",
  robots: "noindex, nofollow",
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
      <div
        v-if="showQuizTitleBar"
        class="flex justify-center px-3 pt-3 sm:px-6 sm:pt-4 md:px-10"
      >
        <div
          class="inline-flex max-w-full items-center gap-3 rotate-[-0.4deg] jv-border-rough bg-jv-white px-3 py-2 shadow-brutal-sm sm:gap-4 sm:px-5 sm:py-2.5"
          :title="quizTitle"
        >
          <span
            class="grid size-9 shrink-0 rotate-[2deg] place-items-center rounded-full border-[2px] border-jv-ink bg-jv-yellow text-jv-ink sm:size-10"
            aria-hidden="true"
          >
            <span class="text-[16px] font-black sm:text-[18px]">Q</span>
          </span>
          <div class="flex min-w-0 flex-col">
            <p
              class="font-body text-[10px] font-bold uppercase tracking-[0.12em] text-jv-muted sm:text-[12px]"
            >
              Playing
            </p>
            <p
              class="min-w-0 truncate font-headings text-[18px] leading-tight text-jv-ink sm:text-[22px]"
            >
              {{ quizTitle }}
            </p>
          </div>
        </div>
      </div>

      <!-- <UserName :user-name="firstname"></UserName> -->

      <QuizWaitingSpacePlayerSkeleton
        v-if="currentComponent === 'Loading'"
      ></QuizWaitingSpacePlayerSkeleton>
      <QuizWaitingSpace
        v-else-if="currentComponent === 'Waiting'"
        :data="data"
        :is-admin="false"
        :user-name="firstname"
        :quiz-title-override="quizTitle"
        @start-quiz="startQuiz"
      >
      </QuizWaitingSpace>
      <QuizQuestionSpace
        v-else-if="currentComponent === 'Question'"
        :data="data"
        :is-admin="false"
        :quiz-title="quizTitle"
        @send-answer="sendAnswer"
      ></QuizQuestionSpace>
      <QuizScoreSpace
        v-else-if="currentComponent === 'Score'"
        :data="data"
        :user-name="username"
        :is-admin="false"
        :selected-answer="selectedAnswer"
        :quiz-state="quizState"
        :quiz-title="quizTitle"
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
