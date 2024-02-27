<script setup>
// core dependencies
import { useToast } from "vue-toastification";

// custom component
import AdminOperation from "../../../composables/admin_operation.js";

// define nuxt configs
const route = useRoute();
const toast = useToast();
const { session } = await useSession();
const cfg = useSystemEnv()

// define props and emits
const url = cfg.value.socket_url
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
  myRef.value = true;
  adminOperationHandler.value.quizStartRequest();
};

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

// main functions
onMounted(() => {
  // core logic
  if (process.client) {
    adminOperationHandler.value = new AdminOperation(
      url,
      route.params.session_id,
      session.value.user?.username,
      handleBackEvent
    );
  }
});

definePageMeta({
  layout: "empty",
});
</script>

<template>
  <Playground :full-screen-enabled="myRef" @is-full-screen="handleCustomChange">
    <LoadingSpace v-if="currentComponent == 'Loading'"></LoadingSpace>
    <WaitingSpace
      v-else-if="currentComponent == 'Waiting'"
      :data="data"
      :is-admin="true"
      @start-quiz="startQuiz"
    ></WaitingSpace>
    <QuestionSpace
      v-else-if="currentComponent == 'Question'"
      :data="data"
      :is-admin="true"
    ></QuestionSpace>
    <ScoreSpace
      v-else-if="currentComponent == 'Score'"
      :data="data"
      :is-admin="true"
    ></ScoreSpace>
  </Playground>
</template>
