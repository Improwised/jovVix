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

//for username
const headers = useRequestHeaders(["cookie"]);
const url = useState("urls");
const endpoint = "/user/who";
const userName = ref("");

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
        route.query?.username,
        handleQuizEvents,
        handleNetworkEvent
      );
    } catch (err) {
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
    monitorTerminateQuiz.value = true;
    return await router.push(`/join/${route.query.username}/scoreboard`);
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
  } else {
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
  }
};

function handleNetworkEvent(message) {
  toast.warning(message + ", please reload the page");
}

const startQuiz = () => {
  myRef.value = true;
};

const sendAnswer = async (answers) => {
  const response = await userOperationHandler.value.handleSendAnswer(answers);

  if (response?.error) {
    toast.error(response.error);
    return;
  }
  toast.success(app.$AnswerSubmitted);
};

definePageMeta({
  layout: "empty",
});

// Listen for beforeunload event to close WebSocket connection
onBeforeUnmount(() => {
  if (!monitorTerminateQuiz.value) {
    userOperationHandler.value.endQuiz();
  }
});

async function getUserNameData() {
  const response = await $fetch(url.value.api_url + endpoint, {
    method: "GET",
    headers,
    credentials: "include",
    mode: "cors",
  });
  userName.value = response.data.username;
}

setTimeout(async () => {
  await getUserNameData();
}, 1000);

// custom class to bind component with
</script>

<template>
  <Playground :full-screen-enabled="myRef" @is-full-screen="handleCustomChange">
    <UserName :user-name="userName"></UserName>

    <QuizLoadingSpace v-if="currentComponent == 'Loading'"></QuizLoadingSpace>
    <QuizWaitingSpace
      v-else-if="currentComponent == 'Waiting'"
      :data="data"
      :is-admin="false"
      @start-quiz="startQuiz"
    >
    </QuizWaitingSpace>
    <QuizQuestionSpace
      v-else-if="currentComponent == 'Question'"
      :data="data"
      :is-admin="false"
      @send-answer="sendAnswer"
    ></QuizQuestionSpace>
    <QuizScoreSpace
      v-else-if="currentComponent == 'Score'"
      :data="data"
      :is-admin="false"
    ></QuizScoreSpace>
  </Playground>
</template>
