<script setup>
import { useSessionStore } from "~~/store/session";
const props = defineProps({
  sessionId: { type: String, default: "" },
  isHost: { type: Boolean, default: false },
});
const { apiUrl } = useRuntimeConfig().public;
const router = useRouter();
const sessionStore = useSessionStore();
let timer;
let disposed = false;
const ping = async () => {
  try {
    if (props.sessionId) {
      const { data: room } = await $fetch(
        `${apiUrl}/rooms/lookup?session_id=${encodeURIComponent(props.sessionId)}`,
        { credentials: "include" }
      );
      if (
        !disposed &&
        (["results", "closed"].includes(room.phase) ||
          room.session_id !== props.sessionId)
      ) {
        const saved = sessionStore.getPlaySessionFor(room.code);
        const participant =
          saved?.sessionId === props.sessionId
            ? saved.userPlayedQuiz
            : room.user_played_quiz;
        const username =
          saved?.sessionId === props.sessionId ? saved.username : room.username;
        if (
          !props.isHost &&
          ["results", "closed"].includes(room.phase) &&
          participant &&
          username
        ) {
          sessionStore.clearPlaySession();
          await router.replace(
            `/join/${encodeURIComponent(username)}/scoreboard?user_played_quiz=${encodeURIComponent(participant)}`
          );
        } else {
          await router.replace(`/room/${room.id}`);
        }
      }
    }
  } catch (err) {
    if (
      !disposed &&
      !props.isHost &&
      (err?.statusCode === 401 || err?.statusCode === 403)
    )
      await router.replace("/join");
  } finally {
    if (!disposed) timer = setTimeout(ping, 30000);
  }
};
onMounted(ping);
onBeforeUnmount(() => {
  disposed = true;
  clearTimeout(timer);
});
</script>
<template><span hidden aria-hidden="true" /></template>
