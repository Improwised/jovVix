<template>
  <Teleport to="body">
    <Transition name="scrim">
      <div
        v-if="visible && isMulti && expanded"
        class="fixed inset-0 z-30 bg-jv-ink/40 sm:hidden"
        @click="expanded = false"
      ></div>
    </Transition>

    <Transition name="dock">
      <div
        v-if="visible"
        class="fixed z-40 left-4 right-4 bottom-[calc(1rem+env(safe-area-inset-bottom))] sm:left-auto sm:right-6 sm:bottom-6 sm:w-[22rem]"
      >
        <!-- Single running quiz: playful hand-drawn card -->
        <div
          v-if="!isMulti"
          class="relative -rotate-2 hover:rotate-0 transition-transform duration-300 ease-out"
        >
          <span
            class="absolute -top-4 -right-3 text-jv-yellow rotate-12 pointer-events-none"
            aria-hidden="true"
          >
            <Sparkle class="size-6 fill-current" :stroke-width="2" />
          </span>
          <span
            class="absolute -bottom-3 -left-3 text-jv-mint -rotate-12 pointer-events-none hidden sm:block"
            aria-hidden="true"
          >
            <svg
              width="44"
              height="18"
              viewBox="0 0 44 18"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M2 10 Q 9 2, 16 8 T 30 9 T 42 5"
                stroke="currentColor"
                stroke-width="3"
                stroke-linecap="round"
                fill="none"
              />
            </svg>
          </span>

          <div
            class="relative bg-jv-white border-[3px] border-jv-ink shadow-brutal p-3.5 sm:p-4"
          >
            <div class="flex items-center gap-2 mb-2.5">
              <LivePulse />
              <span
                class="font-headings text-[11px] tracking-[0.22em] uppercase text-jv-ink/80"
              >
                Live · Quiz running
              </span>
            </div>

            <p
              class="font-headings text-[15px] sm:text-base text-jv-ink leading-snug truncate mb-3"
              :title="sessions[0].title || undefined"
            >
              {{ sessions[0].title || "Your session is still going" }}
            </p>

            <div class="flex items-stretch gap-2">
              <button
                type="button"
                :disabled="resumingId === sessions[0].id"
                class="flex-1 inline-flex items-center justify-center gap-1.5 jv-border-even bg-jv-yellow px-3 py-2 text-[13px] sm:text-sm font-headings text-jv-ink shadow-brutal-sm transition-transform active:translate-x-[2px] active:translate-y-[2px] active:shadow-none hover:-rotate-1 disabled:opacity-60 disabled:cursor-not-allowed"
                @click="resumeQuiz(sessions[0])"
              >
                <ExternalLink class="size-4" :stroke-width="2.4" />
                <span>Resume</span>
              </button>
              <button
                type="button"
                :disabled="stopping"
                class="inline-flex items-center justify-center gap-1.5 jv-border-even bg-jv-coral px-3 py-2 text-[13px] sm:text-sm font-headings text-white shadow-brutal-sm transition-transform active:translate-x-[2px] active:translate-y-[2px] active:shadow-none hover:rotate-1 disabled:opacity-60 disabled:cursor-not-allowed"
                aria-label="Stop quiz"
                @click="openConfirm(sessions[0])"
              >
                <Ban class="size-4" :stroke-width="2.4" />
                <span class="sr-only sm:not-sr-only">Stop</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Many running quizzes: collapsed counter pill -->
        <div v-else-if="!expanded" class="flex justify-end">
          <button
            type="button"
            class="inline-flex items-center gap-2.5 jv-border-even bg-jv-white py-2.5 pl-3 pr-2.5 shadow-brutal transition-transform active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:hover:-rotate-1"
            :aria-expanded="false"
            aria-controls="running-quiz-panel"
            @click="expanded = true"
          >
            <LivePulse />
            <span
              class="font-headings text-[11px] tracking-[0.18em] uppercase text-jv-ink"
            >
              Live · {{ sessions.length }} Quiz running
            </span>
            <span
              class="grid size-6 place-items-center rounded-full bg-jv-yellow border-2 border-jv-ink"
              aria-hidden="true"
            >
              <ChevronUp class="size-3.5" :stroke-width="3" />
            </span>
          </button>
        </div>

        <!-- Many running quizzes: expanded scrollable panel -->
        <div
          v-else
          id="running-quiz-panel"
          class="bg-jv-white border-[3px] border-jv-ink shadow-brutal overflow-hidden"
        >
          <div
            class="flex items-center gap-2 border-b-[3px] border-jv-ink bg-jv-yellow px-3 py-2.5"
          >
            <LivePulse />
            <span
              class="flex-1 font-headings text-[11px] tracking-[0.18em] uppercase text-jv-ink"
            >
              Live · {{ sessions.length }} quizzes
            </span>
            <button
              type="button"
              class="grid size-8 place-items-center border-2 border-jv-ink bg-jv-white transition-transform active:translate-x-[1px] active:translate-y-[1px]"
              aria-label="Collapse running quizzes"
              @click="expanded = false"
            >
              <ChevronDown class="size-4" :stroke-width="3" />
            </button>
          </div>

          <ul
            class="divide-y-2 divide-dashed divide-jv-ink/15 overflow-y-auto overscroll-contain px-3 max-h-[min(58vh,22rem)]"
          >
            <li
              v-for="item in sessions"
              :key="item.id"
              class="flex items-center gap-2 py-2.5"
            >
              <div class="min-w-0 flex-1">
                <p
                  class="truncate font-headings text-[13px] sm:text-sm text-jv-ink leading-snug"
                  :title="item.title || 'Untitled quiz'"
                >
                  {{ item.title || "Untitled quiz" }}
                </p>
                <p
                  v-if="item.invitation_code"
                  class="mt-0.5 font-feature text-[10px] tracking-[0.16em] uppercase text-jv-ink/55"
                >
                  Pin {{ item.invitation_code }}
                </p>
              </div>

              <button
                type="button"
                :disabled="resumingId === item.id"
                class="grid size-10 sm:size-9 shrink-0 place-items-center jv-border-even bg-jv-yellow text-jv-ink shadow-brutal-sm transition-transform active:translate-x-[2px] active:translate-y-[2px] active:shadow-none disabled:opacity-60 disabled:cursor-not-allowed"
                :aria-label="`Resume ${item.title || 'quiz'}`"
                :title="`Resume ${item.title || 'quiz'}`"
                @click="resumeQuiz(item)"
              >
                <ExternalLink class="size-4" :stroke-width="2.4" />
              </button>
              <button
                type="button"
                :disabled="stopping"
                class="grid size-10 sm:size-9 shrink-0 place-items-center jv-border-even bg-jv-coral text-white shadow-brutal-sm transition-transform active:translate-x-[2px] active:translate-y-[2px] active:shadow-none disabled:opacity-60 disabled:cursor-not-allowed"
                :aria-label="`Stop ${item.title || 'quiz'}`"
                :title="`Stop ${item.title || 'quiz'}`"
                @click="openConfirm(item)"
              >
                <Ban class="size-4" :stroke-width="2.4" />
              </button>
            </li>
          </ul>
        </div>
      </div>
    </Transition>
  </Teleport>

  <!-- Confirm dialog -->
  <Teleport to="body">
    <Transition name="overlay">
      <div
        v-if="confirmOpen"
        class="fixed inset-0 z-[100] grid place-items-center bg-jv-ink/40 px-4"
        role="dialog"
        aria-modal="true"
        aria-labelledby="stop-quiz-title"
        @click.self="closeConfirm"
      >
        <div
          class="w-full max-w-md jv-border-uneven bg-jv-white p-6 shadow-brutal"
        >
          <h3
            id="stop-quiz-title"
            class="font-headings text-xl text-jv-ink mb-2"
          >
            Stop the running quiz?
          </h3>
          <p class="text-sm text-jv-ink/75 mb-5">
            <template v-if="pendingStop?.title">
              This will terminate
              <span class="font-semibold">{{ pendingStop.title }}</span> for all
              connected players. This cannot be undone.
            </template>
            <template v-else>
              This will terminate the session for all connected players. This
              cannot be undone.
            </template>
          </p>
          <div class="flex justify-end gap-3">
            <button
              type="button"
              class="inline-flex items-center jv-border-even bg-jv-white px-4 py-2 text-sm font-headings text-jv-ink shadow-brutal-sm transition-transform active:translate-x-[2px] active:translate-y-[2px] active:shadow-none"
              :disabled="stopping"
              @click="closeConfirm"
            >
              Cancel
            </button>
            <button
              type="button"
              class="inline-flex items-center gap-1.5 jv-border-even bg-jv-coral px-4 py-2 text-sm font-headings text-white shadow-brutal-sm transition-transform active:translate-x-[2px] active:translate-y-[2px] active:shadow-none disabled:opacity-60 disabled:cursor-not-allowed"
              :disabled="stopping"
              @click="confirmStop"
            >
              <Ban v-if="!stopping" class="size-4" :stroke-width="2.4" />
              <span>{{ stopping ? "Stopping..." : "Stop quiz" }}</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { computed, h, onBeforeUnmount, onMounted, ref, watch } from "vue";
import {
  Ban,
  ChevronDown,
  ChevronUp,
  ExternalLink,
  Sparkle,
} from "lucide-vue-next";
import { usePush } from "notivue";
import { useSessionStore } from "~~/store/session";
import { useUsersStore } from "~~/store/users";

const POLL_INTERVAL_MS = 60000;

const LivePulse = () =>
  h(
    "span",
    { class: "relative flex size-2.5 shrink-0", "aria-hidden": "true" },
    [
      h("span", {
        class:
          "absolute inline-flex h-full w-full rounded-full bg-jv-coral opacity-75 animate-ping",
      }),
      h("span", {
        class: "relative inline-flex size-2.5 rounded-full bg-jv-coral",
      }),
    ]
  );

const { apiUrl } = useRuntimeConfig().public;
const toast = usePush();
const router = useRouter();

const sessionStore = useSessionStore();
const { syncActiveSessions, removeActiveSession } = sessionStore;

const userDataStore = useUsersStore();
const route = useRoute();

const mounted = ref(false);
const expanded = ref(false);
const confirmOpen = ref(false);
const stopping = ref(false);
const resumingId = ref(null);
const pendingStop = ref(null);

const sessions = computed(() => sessionStore.activeSessions);
const onArrangePage = computed(() => route.path.startsWith("/admin/arrange/"));

const visible = computed(
  () =>
    mounted.value &&
    sessions.value.length > 0 &&
    Boolean(userDataStore.getUserData()) &&
    !onArrangePage.value
);

const isMulti = computed(() => sessions.value.length > 1);

const sync = (force = false) => syncActiveSessions({ force });

const onFocus = () => sync();

const onVisibilityChange = () => {
  if (document.visibilityState === "visible") sync();
};

let pollTimer = null;

const startPolling = () => {
  if (pollTimer) return;
  pollTimer = setInterval(() => {
    if (document.visibilityState === "visible") sync(true);
  }, POLL_INTERVAL_MS);
};

const stopPolling = () => {
  if (!pollTimer) return;
  clearInterval(pollTimer);
  pollTimer = null;
};

onMounted(() => {
  mounted.value = true;
  sync();
  window.addEventListener("focus", onFocus);
  document.addEventListener("visibilitychange", onVisibilityChange);
  window.addEventListener("keydown", onKeydown);
});

onBeforeUnmount(() => {
  window.removeEventListener("focus", onFocus);
  document.removeEventListener("visibilitychange", onVisibilityChange);
  window.removeEventListener("keydown", onKeydown);
  stopPolling();
});

watch(visible, (isVisible) => {
  if (isVisible) {
    startPolling();
  } else {
    stopPolling();
    expanded.value = false;
  }
});

watch(
  () => userDataStore.userData,
  (value) => {
    if (value) sync(true);
  }
);

watch(onArrangePage, (isOnArrange, wasOnArrange) => {
  if (wasOnArrange && !isOnArrange) sync(true);
});

watch(sessions, (list) => {
  if (list.length <= 1) expanded.value = false;
});

const onKeydown = (event) => {
  if (event.key !== "Escape") return;
  if (confirmOpen.value) closeConfirm();
  else expanded.value = false;
};

const resumeQuiz = async (target) => {
  if (!target?.id || resumingId.value) return;
  resumingId.value = target.id;
  try {
    await sync(true);
    if (!sessions.value.some((s) => s.id === target.id)) {
      toast.error("That session has already ended.");
      return;
    }
    router.push(`/admin/arrange/${target.id}`);
  } finally {
    resumingId.value = null;
  }
};

const openConfirm = (target) => {
  pendingStop.value = target;
  confirmOpen.value = true;
};

const closeConfirm = () => {
  if (stopping.value) return;
  confirmOpen.value = false;
  pendingStop.value = null;
};

const confirmStop = async () => {
  const target = pendingStop.value;
  if (!target?.id || stopping.value) return;
  stopping.value = true;
  try {
    await $fetch(`${apiUrl}/quiz/terminate?session_id=${target.id}`, {
      method: "GET",
      credentials: "include",
    });
    removeActiveSession(target.id);
    toast.success("Quiz stopped successfully.");
    confirmOpen.value = false;
    pendingStop.value = null;
    await sync(true);
  } catch (error) {
    console.error("failed to stop quiz from banner", error);
    toast.error("Failed to stop running quiz.");
  } finally {
    stopping.value = false;
  }
};
</script>

<style scoped>
.dock-enter-active {
  transition: transform 0.35s cubic-bezier(0.2, 0.9, 0.3, 1.4),
    opacity 0.25s ease;
}
.dock-leave-active {
  transition: transform 0.2s ease, opacity 0.2s ease;
}
.dock-enter-from {
  transform: translateY(120%) rotate(-8deg);
  opacity: 0;
}
.dock-leave-to {
  transform: translateY(60%) rotate(-4deg);
  opacity: 0;
}

.overlay-enter-active,
.overlay-leave-active,
.scrim-enter-active,
.scrim-leave-active {
  transition: opacity 0.15s ease;
}
.overlay-enter-from,
.overlay-leave-to,
.scrim-enter-from,
.scrim-leave-to {
  opacity: 0;
}

@media (prefers-reduced-motion: reduce) {
  .dock-enter-active,
  .dock-leave-active {
    transition: opacity 0.2s ease;
  }
  .dock-enter-from,
  .dock-leave-to {
    transform: none;
  }
}
</style>
