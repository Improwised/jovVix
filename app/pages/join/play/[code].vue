<script setup>
// core dependencies
import { useToast } from "vue-toastification";

// custom component
import UserOperation from "../../../composables/user_operation.js";

// define nuxt configs
const route = useRoute();
const toast = useToast();
const { session } = await useSession();
const cfg = useSystemEnv();

// define props and emits
const myRef = ref(false);
const socket_url = ref(cfg.value.socket_url);
const data = ref({});
const currentComponent = ref("Loading");

// event handlers
const handleCustomChange = (isFullScreenEvent) => {
  if (!isFullScreenEvent && myRef.value) {
    toast.error("exit fullscreen mode unexpectedly!!!");
    // handle unexpected behavior
  }
};

console.log(socket_url.value, cfg.value.socket_url);
// main functions
onMounted(() => {
  // core logic
  if (process.server) {
  }
  if (process.client) {
    console.log(route.params.code);
    const userSession = new UserOperation(
      socket_url.value,
      route.params.code,
      session.value.user?.username || route.query?.username,
      handleBackEvent
    );
    console.log(userSession);
  }
});

const handleBackEvent = (message) => {
  if (message.status == "fail") {
    console.log(message);
    navigateTo("/error?status=" + message.status + "&error=" + message.data);
  } else if (message.status == "success") {
    if (message?.component) {
      const component = message.component;
      data.value = message;
      console.log(data.value, component);
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
