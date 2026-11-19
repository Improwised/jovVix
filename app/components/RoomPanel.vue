<script setup>
import { useSessionStore } from "~~/store/session";
import JvSelect from "@/components/ui/select/JvSelect.vue";
import NavigationLink from "@/components/common/NavigationLink.vue";
import RoomStandings from "@/components/RoomStandings.vue";
import { useResizeObserver } from "@vueuse/core";
import {
  ArrowRight,
  Check,
  Flag,
  Home,
  LogOut,
  Radio,
  Users,
  WifiOff,
  BarChart3,
} from "lucide-vue-next";
const props = defineProps({
  roomId: { type: String, default: "" },
  sessionId: { type: String, default: "" },
  participantId: { type: String, default: "" },
  standalone: { type: Boolean, default: false },
  actionOnly: { type: Boolean, default: false },
});
const emit = defineEmits(["room-state", "room-id"]);
const { apiUrl } = useRuntimeConfig().public;
const router = useRouter();
const route = useRoute();
const store = useSessionStore();
const room = ref(null);
const playersList = ref(null);
const playersExpanded = ref(false);
const playersOverflow = ref(false);
const measurePlayers = () => {
  const list = playersList.value;
  playersOverflow.value = Boolean(list && list.scrollHeight > 36);
};
useResizeObserver(playersList, measurePlayers);
const onlineMembers = computed(() =>
  (room.value?.members || []).filter((member) => member.online)
);
const choices = ref([]);
const choicesLoaded = ref(false);
const quizSearch = ref("");
const visibilityFilter = ref("Private");
const quizOptions = computed(() =>
  choices.value
    .filter((quiz) => quiz.is_public === (visibilityFilter.value === "Public"))
    .filter((quiz) =>
      quiz.title.toLowerCase().includes(quizSearch.value.trim().toLowerCase())
    )
    .map((quiz) => ({
      value: quiz.id,
      label: quiz.title,
      badge:
        typeof quiz.question_count === "number"
          ? `${quiz.question_count} ${quiz.question_count === 1 ? "question" : "questions"}`
          : "",
    }))
);
const selected = ref("");
watch(visibilityFilter, () => {
  selected.value = "";
  quizSearch.value = "";
});
watch(onlineMembers, () => nextTick(measurePlayers));
const error = ref("");
const errorStatus = ref(null);
const canRejoin = computed(
  () => !room.value?.is_host && [401, 403, 404, 410].includes(errorStatus.value)
);
const busy = ref(false);
const launching = ref(false);
const movingToQuiz = ref("");
let timer;
let disposed = false;
let navigating = false;
let polling = false;
let leaving = false;
let resolvedRoomId = props.roomId;
const request = (path, options = {}) =>
  $fetch(`${apiUrl}${path}`, { credentials: "include", ...options });
const refresh = async () => {
  if (disposed || polling || navigating) return;
  clearTimeout(timer);
  polling = true;
  try {
    const path = resolvedRoomId
      ? `/rooms/${resolvedRoomId}`
      : `/rooms/lookup?session_id=${encodeURIComponent(props.sessionId)}&participant_id=${encodeURIComponent(props.participantId)}`;
    const response = await request(path);
    if (disposed || leaving) return;
    room.value = response.data;
    emit("room-state", room.value.phase);
    emit("room-id", room.value.id);
    resolvedRoomId = room.value.id;
    error.value = "";
    errorStatus.value = null;
    const r = room.value;
    if (
      route.query.view !== "quiz" &&
      ["lobby", "playing"].includes(r.phase) &&
      (props.standalone ||
        (props.sessionId && r.session_id !== props.sessionId) ||
        (props.participantId && r.user_played_quiz !== props.participantId))
    ) {
      navigating = true;
      if (r.is_host) {
        store.setSession(r.session_id);
        store.setActiveQuizTitle(r.title);
        setSocketObject(null);
        await router.push(`/admin/arrange/${r.session_id}`);
      } else {
        const played = (
          await request(`/user_played_quizes/${r.code}`, { method: "POST" })
        ).data;
        if (disposed || leaving) return;
        store.setPlaySession({
          code: String(r.code),
          username: r.username,
          firstname: played.player_name,
          sessionId: played.session_id,
          userPlayedQuiz: played.user_played_quiz,
        });
        store.setActiveQuizTitle(r.title);
        movingToQuiz.value = r.title || "the next quiz";
        await new Promise((resolve) => setTimeout(resolve, 1800));
        if (disposed || leaving) return;
        await router.push(`/join/play/${r.code}`);
      }
    }
    if (
      !props.actionOnly &&
      r.is_host &&
      r.phase === "results" &&
      !choicesLoaded.value
    ) {
      const quizzes = (await request(`/rooms/${r.id}/quizzes`)).data;
      choices.value = quizzes;
      choicesLoaded.value = true;
    }
  } catch (err) {
    errorStatus.value = err?.statusCode || err?.response?.status;
    navigating = false;
    movingToQuiz.value = "";
    if (err?.statusCode === 404 && !props.standalone) {
      emit("room-state", "unavailable");
      return;
    }
    error.value =
      typeof err?.data?.data === "string"
        ? err.data.data
        : "Connection interrupted. Retrying…";
  } finally {
    polling = false;
    if (
      !disposed &&
      !leaving &&
      !navigating &&
      room.value?.phase !== "closed" &&
      errorStatus.value !== 404
    )
      timer = setTimeout(refresh, 3000);
  }
};
const next = async () => {
  if (busy.value || !selected.value) return;
  busy.value = true;
  launching.value = true;
  try {
    await request(`/rooms/${room.value.id}/next`, {
      method: "POST",
      body: {
        quiz_id: selected.value,
        expected_session_id: room.value.session_id,
      },
    });
    clearTimeout(timer);
    await refresh();
  } catch (err) {
    error.value = err?.data?.data || "Could not select this quiz.";
  } finally {
    busy.value = false;
    launching.value = false;
  }
};
const leave = async () => {
  if (busy.value) return;
  leaving = true;
  busy.value = true;
  try {
    await request(`/rooms/${room.value.id}/leave`, { method: "POST" });
    store.clearPlaySession();
    await router.push("/join");
  } catch {
    leaving = false;
    error.value = "Could not leave the room. Try again.";
    clearTimeout(timer);
    timer = setTimeout(refresh, 3000);
  } finally {
    busy.value = false;
  }
};
const closeConfirmOpen = ref(false);
const close = () => {
  if (!busy.value) closeConfirmOpen.value = true;
};
const confirmClose = async (confirmed) => {
  closeConfirmOpen.value = false;
  if (!confirmed || busy.value) return;
  busy.value = true;
  try {
    await request(`/rooms/${room.value.id}/close`, { method: "POST" });
    clearTimeout(timer);
    await refresh();
  } catch (err) {
    error.value = err?.data?.data || "Could not end the room.";
  } finally {
    busy.value = false;
  }
};
onMounted(refresh);
onBeforeUnmount(() => {
  disposed = true;
  clearTimeout(timer);
});
</script>
<template>
  <div
    v-if="actionOnly && (room || error || $slots.actions)"
    class="flex flex-wrap items-center justify-between gap-3"
  >
    <UtilsConfirmModal
      v-if="closeConfirmOpen"
      modal-title="End this room?"
      modal-message="This will end the room for everyone and disable its join code. Completed quiz results will remain available."
      model-positive-message="End room"
      @confirm-message="confirmClose"
    />
    <p v-if="error" role="alert" class="text-sm font-bold text-jv-muted">
      {{ error }}
    </p>
    <p
      v-if="movingToQuiz"
      role="status"
      aria-live="polite"
      class="text-sm font-bold text-jv-muted"
    >
      The host selected {{ movingToQuiz }}. Moving to the lobby…
    </p>
    <p
      v-else-if="room?.phase === 'closed'"
      role="status"
      class="text-sm font-bold text-jv-muted"
    >
      This room has ended. You can join another quiz.
    </p>
    <p
      v-else-if="room && !room.is_host"
      class="text-sm font-bold text-jv-muted"
      role="status"
    >
      {{
        room.host_online
          ? "Waiting for the host to choose the next quiz"
          : "Waiting for the host to return"
      }}
      <span
        v-if="room.host_online"
        class="waiting-dots ml-1 inline-flex gap-1"
        aria-hidden="true"
        ><i></i><i></i><i></i
      ></span>
    </p>
    <div
      class="ml-auto flex w-full min-w-0 max-w-full flex-wrap items-center justify-center gap-3 sm:w-auto sm:justify-end"
    >
      <NavigationLink
        v-if="room && room.phase !== 'closed' && room.is_host"
        variant="result"
        :url="`/room/${room.id}`"
        class="w-full bg-jv-mint sm:w-auto"
        >Choose next quiz<ArrowRight :stroke-width="2.4"
      /></NavigationLink>
      <NavigationLink
        v-if="room && room.phase !== 'closed' && !room.is_host"
        variant="result"
        type="button"
        :disabled="busy"
        @click="leave"
        ><LogOut :stroke-width="2.4" />Leave room</NavigationLink
      >
      <slot name="actions" />
      <NavigationLink
        v-if="room && room.phase !== 'closed' && room.is_host"
        variant="result"
        type="button"
        :disabled="busy"
        @click="close"
        ><Flag :stroke-width="2.4" />End room</NavigationLink
      >
    </div>
  </div>
  <div v-else-if="!actionOnly && (room || error || standalone)">
    <UtilsConfirmModal
      v-if="closeConfirmOpen"
      modal-title="End this room?"
      modal-message="This will end the room for everyone and disable its join code. Completed quiz results will remain available."
      model-positive-message="End room"
      @confirm-message="confirmClose"
    />

    <section
      v-if="!(standalone && room?.phase === 'closed')"
      class="room-panel relative z-20 mx-auto my-3 w-[calc(100%-2rem)] max-w-7xl overflow-hidden bg-jv-white font-body text-jv-ink shadow-brutal jv-border-rough"
      :class="{ 'room-page': standalone }"
      aria-label="Quiz room"
    >
      <div
        v-if="error"
        role="alert"
        class="border-b-2 border-jv-ink bg-jv-yellow-soft px-5 py-3 text-sm"
      >
        {{ error }}
        <button
          type="button"
          class="ml-2 cursor-pointer font-bold underline"
          @click="refresh"
        >
          Try again
        </button>
        <NuxtLink v-if="canRejoin" to="/join" class="ml-2 font-bold underline"
          >Join a room</NuxtLink
        >
      </div>
      <template v-if="room">
        <header
          class="flex flex-wrap items-center justify-between gap-3 border-b-2 border-jv-ink bg-jv-white px-4 py-2 sm:px-5"
        >
          <p
            class="flex items-center gap-2 text-xs font-black uppercase tracking-widest"
          >
            <span
              class="grid size-5 place-items-center rounded-full border-2 border-jv-ink bg-jv-mint"
              ><Check class="size-3" :stroke-width="3" /></span
            >{{
              room.phase === "closed"
                ? "Thanks for playing"
                : "Stay for the next quiz"
            }}
          </p>
          <p v-if="room.phase !== 'closed'" class="flex items-center gap-2">
            <span class="text-sm font-bold text-jv-muted">Code:</span>
            <strong class="font-feature text-xl font-black text-jv-coral">{{
              room.code
            }}</strong>
          </p>
        </header>
        <div class="room-main grid items-center gap-4 px-4 py-3 sm:px-5">
          <div class="room-art" aria-hidden="true">
            <NuxtImg
              src="/images/jovvix-logo.png"
              alt="jovVix"
              width="300"
              height="122"
              class="room-brand"
            />
          </div>
          <div class="room-copy min-w-0 self-center">
            <p
              class="mb-1 text-[10px] font-black uppercase tracking-[.12em] text-jv-muted"
            >
              {{ room.phase === "results" ? "Quiz complete" : "Quiz ready" }}
            </p>
            <h2
              class="font-headings text-xl font-bold leading-tight sm:text-[22px]"
            >
              {{
                room.is_host
                  ? "One more round?"
                  : room.host_online
                    ? "A new challenge is on its way!"
                    : "Saving your spot."
              }}
            </h2>
            <p
              v-if="room.is_host"
              class="mt-1 text-xs leading-relaxed text-jv-muted"
            >
              Pick the next quiz. Your crew and join code come with you.
            </p>
            <div v-else class="mt-2" role="status" aria-live="polite">
              <p class="flex items-center gap-2 text-xs font-extrabold">
                <template v-if="movingToQuiz">
                  The host selected {{ movingToQuiz }}. Moving to the lobby…
                </template>
                <template v-else>
                  <Radio
                    v-if="room.host_online"
                    class="size-4 shrink-0 text-jv-accent-green"
                  /><WifiOff v-else class="size-4 shrink-0" /><span
                    >{{
                      room.host_online
                        ? "Waiting for the host to choose the next quiz"
                        : "The host is disconnected. Waiting for them to return."
                    }}<span
                      v-if="room.host_online"
                      class="waiting-dots ml-1 inline-flex gap-1"
                      aria-hidden="true"
                      ><i></i><i></i><i></i></span
                  ></span>
                </template>
              </p>
            </div>
          </div>
          <div
            v-if="room.is_host && room.phase === 'results'"
            class="room-picker"
          >
            <label class="min-w-0">
              <span
                class="mb-1 flex items-center gap-2 text-[11px] font-black uppercase tracking-[0.08em]"
                >Choose the next quiz</span
              >
              <JvSelect
                v-model="selected"
                :options="quizOptions"
                :placeholder="`Choose a ${visibilityFilter.toLowerCase()} quiz…`"
                aria-label="Choose the next quiz"
                :disabled="busy"
                class="jv-card !h-10 !px-3 !text-sm"
              >
                <template #header>
                  <div
                    class="flex gap-1 border-b border-jv-ink/20 p-2"
                    role="group"
                    aria-label="Filter by visibility"
                  >
                    <button
                      v-for="option in ['Private', 'Public']"
                      :key="option"
                      type="button"
                      class="cursor-pointer rounded-[6px] px-3 py-1.5 text-sm font-semibold text-jv-ink transition-colors hover:bg-jv-yellow/40"
                      :class="
                        option === visibilityFilter ? 'bg-jv-yellow/60' : ''
                      "
                      :aria-pressed="option === visibilityFilter"
                      @click.stop="visibilityFilter = option"
                    >
                      {{ option }}
                    </button>
                  </div>
                  <div class="border-b border-jv-ink/20 p-2">
                    <input
                      v-model="quizSearch"
                      type="search"
                      aria-label="Search quizzes"
                      placeholder="Search quizzes…"
                      class="w-full rounded border border-jv-ink/30 bg-jv-white px-3 py-2 text-sm text-jv-ink"
                      @keydown.stop
                    />
                  </div>
                </template>
                <template #empty>
                  <p class="px-3 py-3 text-sm text-jv-muted">
                    {{
                      choicesLoaded
                        ? quizSearch.trim()
                          ? "No quizzes match your search."
                          : `No ${visibilityFilter.toLowerCase()} quizzes are available to host.`
                        : "Loading quizzes..."
                    }}
                  </p>
                  <button
                    v-if="
                      choicesLoaded &&
                      visibilityFilter === 'Private' &&
                      !quizSearch.trim()
                    "
                    type="button"
                    class="mx-3 mb-3 cursor-pointer text-sm font-bold underline underline-offset-4"
                    @click.stop="visibilityFilter = 'Public'"
                  >
                    Browse public quizzes
                  </button>
                </template>
              </JvSelect>
            </label>
            <NavigationLink
              type="button"
              variant="result"
              class="w-full bg-jv-mint"
              :disabled="busy || !selected"
              @click="next"
            >
              <span>{{
                launching ? "Getting the lobby ready…" : "Continue to lobby"
              }}</span>
              <ArrowRight class="size-4" :stroke-width="2.4" />
            </NavigationLink>
          </div>
        </div>
        <footer
          class="border-t border-jv-ink/20 bg-jv-canvas px-4 py-2 sm:px-5"
        >
          <div class="flex flex-wrap items-center justify-between gap-3">
            <p class="flex items-center gap-2 text-[11px] font-extrabold">
              <Users class="size-4" />{{ onlineMembers.length }}
              {{ onlineMembers.length === 1 ? "player" : "players" }}
              connected
            </p>
            <NavigationLink
              v-if="!room.is_host"
              type="button"
              variant="result"
              :disabled="busy"
              @click="leave"
              ><LogOut class="size-4" :stroke-width="2.4" />Leave
              room</NavigationLink
            >
            <NavigationLink
              v-else-if="room.phase === 'results'"
              type="button"
              variant="result"
              :disabled="busy"
              @click="close"
              ><Flag class="size-4" :stroke-width="2.4" />End
              room</NavigationLink
            >
          </div>
          <ul
            v-if="standalone && onlineMembers.length"
            id="room-players"
            ref="playersList"
            class="mt-3 flex flex-wrap gap-2 overflow-hidden"
            :class="{ 'max-h-9': !playersExpanded }"
            aria-label="Players in this room"
          >
            <li
              v-for="member in onlineMembers"
              :key="member.name"
              class="inline-flex h-9 max-w-full shrink-0 items-center gap-2 rounded-full border-2 border-jv-ink bg-jv-white px-3 text-sm font-bold"
            >
              <span class="truncate">{{ member.name }}</span>
              <span
                class="size-2 shrink-0 rounded-full bg-jv-accent-green"
                aria-label="Online"
              ></span>
            </li>
          </ul>
          <button
            v-if="standalone && playersOverflow"
            type="button"
            class="mt-3 cursor-pointer text-xs font-bold underline underline-offset-4"
            :aria-expanded="playersExpanded"
            aria-controls="room-players"
            @click="playersExpanded = !playersExpanded"
          >
            {{ playersExpanded ? "Show less" : "Show all players" }}
          </button>
        </footer>
      </template>
      <p v-else-if="!error" class="p-7 text-sm font-bold" role="status">
        Connecting to room…
      </p>
    </section>
    <RoomStandings
      v-if="standalone && room?.phase === 'closed'"
      :room-id="room.id"
      class="mx-auto mt-4 w-[calc(100%-2rem)] max-w-7xl"
    />
    <div
      v-if="standalone && room"
      class="mb-2 mt-8 flex flex-wrap items-center justify-center gap-3"
      role="group"
      aria-label="Room actions"
    >
      <NavigationLink url="/" variant="result" aria-label="Return to home"
        ><Home :stroke-width="2.4" />Home</NavigationLink
      >
      <NavigationLink
        v-if="room.phase === 'closed'"
        :url="`/room/report/${room.id}`"
        variant="result"
        class="bg-jv-yellow"
        ><BarChart3 :stroke-width="2.4" />View final room
        standings</NavigationLink
      >
      <NavigationLink
        v-else-if="
          room.phase === 'results' && (room.is_host || room.user_played_quiz)
        "
        :url="
          room.is_host
            ? `/admin/scoreboard?aqi=${room.session_id}`
            : `/join/${encodeURIComponent(room.username)}/scoreboard?user_played_quiz=${room.user_played_quiz}`
        "
        variant="result"
        ><BarChart3 class="size-4" :stroke-width="2.4" />View quiz
        results</NavigationLink
      >
    </div>
  </div>
</template>
<style scoped>
.room-main {
  grid-template-columns: 90px minmax(180px, 1fr) minmax(380px, 1.4fr);
}
.room-main:has(.room-picker) {
  column-gap: 20px;
}
.room-art {
  display: grid;
  place-items: center;
  width: 90px;
}
.room-brand {
  width: 90px;
  height: auto;
  object-fit: contain;
}
.room-picker {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: end;
  gap: 12px;
  min-width: 0;
}
.room-picker > p {
  grid-column: 1 / -1;
}
.room-main:not(:has(.room-picker)) {
  grid-template-columns: 90px minmax(0, 1fr);
}
.waiting-dots i {
  width: 4px;
  height: 4px;
  background: currentColor;
  border-radius: 50%;
  animation: waiting-dot 1.5s ease-in-out infinite;
}
.waiting-dots i:nth-child(2) {
  animation-delay: 0.2s;
}
.waiting-dots i:nth-child(3) {
  animation-delay: 0.4s;
}
@keyframes waiting-dot {
  0%,
  60%,
  100% {
    transform: translateY(0);
    opacity: 0.35;
  }
  30% {
    transform: translateY(-3px);
    opacity: 1;
  }
}
@media (max-width: 1023px) {
  .room-main {
    grid-template-columns: 70px minmax(0, 1fr);
  }
  .room-art,
  .room-brand {
    width: 70px;
  }
  .room-picker {
    grid-column: 1 / -1;
  }
  .room-main:not(:has(.room-picker)) {
    grid-template-columns: 70px minmax(0, 1fr);
  }
}
@media (max-width: 540px) {
  .room-picker {
    grid-template-columns: minmax(0, 1fr);
    gap: 8px;
  }
  .room-main {
    gap: 12px;
  }
}

.room-page {
  max-width: 1220px;
}
.room-page > header {
  padding: 20px 32px;
}
.room-page .room-main,
.room-page .room-main:not(:has(.room-picker)) {
  grid-template-columns: minmax(160px, 0.75fr) minmax(0, 2fr);
  gap: 16px 48px;
  padding: 48px;
  min-height: 290px;
}
.room-page .room-art {
  grid-row: 1 / 3;
  width: 100%;
}
.room-page .room-brand {
  width: min(100%, 220px);
}
.room-page .room-copy h2 {
  font-size: clamp(30px, 3vw, 44px);
}
.room-page .room-copy > p:first-child {
  font-size: 12px;
}
.room-page .room-copy > p,
.room-page .room-copy [role="status"] p {
  font-size: 16px;
}
.room-page .room-picker {
  grid-column: 2;
  width: 100%;
}
.room-page .room-picker label > span {
  font-size: 13px;
}
.room-page > footer {
  padding: 20px 32px;
}
.room-page footer p {
  font-size: 14px;
}
@media (max-width: 640px) {
  .room-page > header {
    padding: 16px;
  }
  .room-page .room-main,
  .room-page .room-main:not(:has(.room-picker)) {
    grid-template-columns: 1fr;
    gap: 20px;
    padding: 28px 20px;
    min-height: 0;
  }
  .room-page .room-art {
    grid-row: auto;
    justify-self: start;
  }
  .room-page .room-brand {
    width: 160px;
  }
  .room-page .room-picker {
    grid-column: 1;
  }
  .room-page > footer {
    align-items: flex-start;
    flex-direction: column;
    padding: 16px 20px;
  }
  .room-page > footer > div {
    flex: none;
    width: 100%;
  }
}
@media (prefers-reduced-motion: reduce) {
  .waiting-dots i,
  .room-panel *,
  .room-panel *::before,
  .room-panel *::after {
    animation: none !important;
    transition: none !important;
  }
}
</style>
