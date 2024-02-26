<script setup>
// core dependencies
import { useToast } from "vue-toastification";

// custom component
import LayoutsPlayground from "../../../components/layouts/playground.vue";
import UserOperation from "../../../composables/user_operation.js";

// define nuxt configs
const route = useRoute();
const toast = useToast();
const { session } = await useSession();
const cfg = useSystemEnv()

// define props and emits
const myRef = ref(false);
const socket_url = cfg.value.api_url

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
    const userSession = initOperation(route, session, socket_url);
  }
});
</script>

<script>
// get userOperation Object
function initOperation(route, session, url) {
  return new UserOperation(
    url,
    route.params.code,
    session.value?.user.username || route.query.username
  );
}

// custom class to bind component with
</script>

<template>
  <LayoutsPlayground
    :full-screen-enabled="myRef"
    @is-full-screen="handleCustomChange"
  >
    <!-- <button @click="() => (myRef = !myRef)">Change to full screen</button> -->
  </LayoutsPlayground>
</template>
