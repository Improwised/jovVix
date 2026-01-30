<script setup>
import { useToast } from "vue-toastification";
import { useListUserstore } from "~/store/userlist";
import { useSessionStore } from "~~/store/session";
const sessionStore = useSessionStore();
const { setSession } = sessionStore;
const listUserStore = useListUserstore();
const { removeAllUsers } = listUserStore;
let urls = useRuntimeConfig().public;
const router = useRouter();
const toast = useToast();
const requestPending = ref(false);
const activeQuizId = ref(false);
const props = defineProps({
  quizId: {
    default: () => {
      return "";
    },
    type: String,
    required: true,
  },
});

const handleStartDemo = async () => {
  try {
    requestPending.value = true;
    await $fetch(`${urls.apiUrl}/quizzes/${props.quizId}/demo_session`, {
      method: "POST",
      credentials: "include",
      headers: {
        Accept: "application/json",
      },
      onResponse({ response }) {
        if (response.status != 202) {
          requestPending.value = false;
          toast.error("error while start quiz");
          return;
        }
        if (response.status == 202) {
          activeQuizId.value = response._data.data;
          requestPending.value = false;
        }
      },
    });
  } catch (error) {
    toast.error(error.message);
    requestPending.value = false;
    return;
  }

  removeAllUsers();
  setSocketObject(null);
  router.push(`/admin/arrange/${activeQuizId.value}`);

  // add session in store after 1 second
  setTimeout(() => {
    setSession(activeQuizId.value);
  }, 1000);
};
</script>
<template>
  <button
    v-if="requestPending"
    type="button"
    class="btn text-white btn-primary me-0"
  >
    Pending...
  </button>
  <button
    v-else
    type="button"
    class="btn text-white btn-primary me-0"
    @click="handleStartDemo"
  >
    Start Quiz
  </button>
</template>
