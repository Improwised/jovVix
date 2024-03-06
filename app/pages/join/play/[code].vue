<script setup>
// core dependencies
import { useToast } from "vue-toastification";

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
    userOperationHandler.value = new UserOperation(
      route.params.code,
      route.query?.username,
      handleQuizEvents
    );
  }
});

const handleQuizEvents = (message) => {
  // error ? -> redirect to error page
  if (message.status == app.$Error) {
    userOperationHandler.value.printLog();
    router.push("/error?status=" + message.status + "&error=" + message.data);
  } else if (message.event == app.$TerminateQuiz) {
    router.push("/join/scoreboard");
  } else {
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
