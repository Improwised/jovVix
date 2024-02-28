<script setup>
// core dependencies
import { useToast } from "vue-toastification";

// custom component
import AdminOperation from "../../../composables/admin_operation.js";

// define nuxt configs
const route = useRoute();
const toast = useToast();
const app = useNuxtApp();
const { session } = await useSession();
useSystemEnv();

// define props and emits
const myRef = ref(false);
const data = ref({});
const adminOperationHandler = ref();
const currentComponent = ref("Loading");

// event handlers
const handleCustomChange = (isFullScreenEvent) => {
  if (!isFullScreenEvent && myRef.value) {
    toast.error("exit fullscreen mode unexpectedly!!!");
    // handle unexpected behavior
  }
};

const startQuiz = () => {
  // myRef.value = true;
  adminOperationHandler.value.quizStartRequest();
};

const handleQuizEvents = (message) => {
  // error ? -> redirect to error page
  if (message.status == app.$Error) {
    adminOperationHandler.value.printLog();
    navigateTo("/error?status=" + message.status + "&error=" + message.data);
  } else {
    // unauthorized ? -> redirect to login page
    if (message.status == app.$Fail && message.data == app.$Unauthorized) {
      navigateTo(
        "/account/login?error=" + message.data + "&url=" + route.fullPath
      );
    }
    data.value = message;
    currentComponent.value = message.component;
  }
};

// main functions
onMounted(() => {
  // core logic
  if (process.client) {
    adminOperationHandler.value = new AdminOperation(
      route.params.session_id,
      session.value.user?.username,
      handleQuizEvents
    );
  }
});

definePageMeta({
  layout: "empty",
});
</script>

<template>
  <Playground :full-screen-enabled="myRef" @is-full-screen="handleCustomChange">
    <QuizLoadingSpace v-if="currentComponent == 'Loading'"></QuizLoadingSpace>
    <QuizWaitingSpace
      v-else-if="currentComponent == 'Waiting'"
      :data="data"
      :is-admin="true"
      @start-quiz="startQuiz"
    ></QuizWaitingSpace>
    <QuizQuestionSpace
      v-else-if="currentComponent == 'Question'"
      :data="data"
      :is-admin="true"
    ></QuizQuestionSpace>
    <QuizScoreSpace
      v-else-if="currentComponent == 'Score'"
      :data="data"
      :is-admin="true"
    ></QuizScoreSpace>
  </Playground>
</template>
