<script setup>
// core dependencies
import { useToast } from "vue-toastification";

// custom component
import UserOperation from "../../../composables/user_operation.js";
import { useSystemEnv } from "~~/composables/envs";
import constants from "~~/config/constants";

// define nuxt configs
const route = useRoute();
const toast = useToast();
const { session } = await useSession();
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
    userSession.value = new UserOperation(
      route.params.code,
      session.value.user?.username || route.query?.username,
      handleQuizEvents
    );
  }
});

const handleQuizEvents = (message) => {
  if (message.status == constants.Error) {
    navigateTo("/error?status=" + message.status + "&error=" + message.data);
  } else {
    if(message.status == constants.Fail && message.data ==constants.CodeNotFound ){
      navigateTo("/error?status=" + message.status + "&error=" + message.data);
    }
    if (message?.component) {
      const component = message.component;
      data.value = message;
      currentComponent.value = component;
    } else {
      console.log(message);
    }
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
    <LoadingSpace v-if="currentComponent == 'Loading'"></LoadingSpace>
    <WaitingSpace
      v-else-if="currentComponent == 'Waiting'"
      :data="data"
      :is-admin="false"
      @start-quiz="startQuiz"
    ></WaitingSpace>
    <QuestionSpace
      v-else-if="currentComponent == 'Question'"
      :data="data"
      :is-admin="false"
    ></QuestionSpace>
    <ScoreSpace
      v-else-if="currentComponent == 'Score'"
      :data="data"
      :is-admin="false"
    ></ScoreSpace>
  </Playground>
</template>
