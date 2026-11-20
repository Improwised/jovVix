<script setup>
import { usePush } from "notivue";
import { useListUserstore } from "~/store/userlist";
import { useSessionStore } from "~~/store/session";

const sessionStore = useSessionStore();
const { setSession } = sessionStore;
const listUserStore = useListUserstore();
const { removeAllUsers } = listUserStore;
let urls = useRuntimeConfig().public;
const router = useRouter();
const toast = usePush();
const requestPending = ref(false);
const activeQuizId = ref(false);
const props = defineProps({
  isPublic: { type: Boolean, default: false },
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
  setSession(activeQuizId.value);
  router.push(`/admin/arrange/${activeQuizId.value}`);
};
</script>
<template>
  <QuizModeAction
    :label="requestPending ? 'Pending...' : 'Host Live'"
    secondary
    :disabled="requestPending"
    @click="handleStartDemo"
  />
</template>
