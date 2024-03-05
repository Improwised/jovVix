<script setup>
// core dependencies
import { useToast } from "vue-toastification";

// custom component
import UserOperation from "~/composables/user_operation.js";
import { useSystemEnv } from "~/composables/envs.js";
import { useRouter } from "nuxt/app";

// define nuxt configs
const route = useRoute();
const toast = useToast();
const app = useNuxtApp();
const router = useRouter();
useSystemEnv();

// define props and emits
const myRef = ref(false);
const data = ref({});
const currentComponent = ref("Loading");
const userSession = ref();

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
      userSession.value = new UserOperation(
        route.params.code,
        route.query?.username,
        handleQuizEvents,
        handleNetworkEvent
      );
    } catch (err) {
      toast.error(err + ", Please reload the page");
    }
  }
});

const handleQuizEvents = async (message) => {
  if (message.status == app.$Error || message.status == app.$Fail) {
    if (
      message.status == app.$Fail &&
      message.event == app.$InvitationCodeValidation
    ) {
      return await router.push(
        "/join?status=" + message.status + "&error=" + message.data
      );
    }
    return await router.push(
      "/error?status=" + message.status + "&error=" + message.data
    );
  } else {
    if (message?.component) {
      const component = message.component;
      data.value = message;
      currentComponent.value = component;
    } else {
      toast.error(`Error! event:${message.event} action:${message.action}`);
    }
  }
};

function handleNetworkEvent(message) {
  toast.warning(message + ", please reload the page");
}

const startQuiz = () => {
  myRef.value = true;
};

definePageMeta({
  layout: "empty",
});

// custom class to bind component with
</script>

<template>
  <Playground :full-screen-enabled="myRef" @is-full-screen="handleCustomChange">
    <QuizLoadingSpace v-if="currentComponent == 'Loading'"></QuizLoadingSpace>
    <QuizWaitingSpace
      v-else-if="currentComponent == 'Waiting'"
      :data="data"
      :is-admin="false"
      @start-quiz="startQuiz"
    ></QuizWaitingSpace>
    <QuizQuestionSpace
      v-else-if="currentComponent == 'Question'"
      :data="data"
      :is-admin="false"
    ></QuizQuestionSpace>
    <QuizScoreSpace
      v-else-if="currentComponent == 'Score'"
      :data="data"
      :is-admin="false"
    ></QuizScoreSpace>
  </Playground>
</template>
