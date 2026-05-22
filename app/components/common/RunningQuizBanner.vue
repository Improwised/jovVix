<template>
  <!-- Floating live card -->
  <Teleport to="body">
    <Transition name="dock">
      <div
        v-if="visible"
        class="fixed z-40 bottom-4 left-4 right-4 sm:bottom-6 sm:right-6 sm:left-auto sm:max-w-[22rem]"
      >
        <div
          class="group relative -rotate-2 hover:rotate-0 transition-transform duration-300 ease-out"
        >
          <!-- Decorative sparkle peeking out top-right -->
          <span
            class="absolute -top-4 -right-3 text-jv-yellow rotate-12 pointer-events-none"
            aria-hidden="true"
          >
            <Sparkle class="size-6 fill-current" :stroke-width="2" />
          </span>
          <!-- Decorative squiggle peeking out bottom-left -->
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

          <!-- Card body -->
          <div
            class="relative bg-jv-white border-[3px] border-jv-ink shadow-brutal p-3.5 sm:p-4"
          >
            <!-- Header strip: live pulse + label + rocket -->
            <div class="flex items-center gap-2 mb-2.5">
              <span class="relative flex size-2.5" aria-hidden="true">
                <span
                  class="absolute inline-flex h-full w-full rounded-full bg-jv-coral opacity-75 animate-ping"
                ></span>
                <span
                  class="relative inline-flex size-2.5 rounded-full bg-jv-coral"
                ></span>
              </span>
              <span
                class="font-headings text-[11px] tracking-[0.22em] uppercase text-jv-ink/80"
              >
                Live · Quiz running
              </span>
            </div>

            <!-- Title -->
            <p
              v-if="quizTitle"
              class="font-headings text-[15px] sm:text-base text-jv-ink leading-snug truncate mb-3"
              :title="quizTitle"
            >
              {{ quizTitle }}
            </p>
            <p
              v-else
              class="font-headings text-[15px] sm:text-base text-jv-ink leading-snug mb-3"
            >
              Your session is still going
            </p>

            <!-- Actions -->
            <div class="flex items-stretch gap-2">
              <button
                type="button"
                class="flex-1 inline-flex items-center justify-center gap-1.5 jv-border-even bg-jv-yellow px-3 py-2 text-[13px] sm:text-sm font-headings text-jv-ink shadow-brutal-sm transition-transform active:translate-x-[2px] active:translate-y-[2px] active:shadow-none hover:-rotate-1"
                @click="resumeQuiz"
              >
                <ExternalLink class="size-4" :stroke-width="2.4" />
                <span>Resume</span>
              </button>
              <button
                type="button"
                :disabled="stopping"
                class="inline-flex items-center justify-center gap-1.5 jv-border-even bg-jv-coral px-3 py-2 text-[13px] sm:text-sm font-headings text-white shadow-brutal-sm transition-transform active:translate-x-[2px] active:translate-y-[2px] active:shadow-none hover:rotate-1 disabled:opacity-60 disabled:cursor-not-allowed"
                aria-label="Stop quiz"
                @click="openConfirm"
              >
                <Ban class="size-4" :stroke-width="2.4" />
                <span class="sr-only sm:not-sr-only">Stop</span>
              </button>
            </div>
          </div>
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
            <template v-if="quizTitle">
              This will terminate
              <span class="font-semibold">{{ quizTitle }}</span> for all
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
import { computed, ref } from "vue";
import { Ban, ExternalLink, Sparkle } from "lucide-vue-next";
import { useToast } from "vue-toastification";
import { useSessionStore } from "~~/store/session";
import { useUsersStore } from "~~/store/users";

const { apiUrl } = useRuntimeConfig().public;
const toast = useToast();
const router = useRouter();

const sessionStore = useSessionStore();
const { getSession, setSession, getActiveQuizTitle } = sessionStore;

const userDataStore = useUsersStore();
const route = useRoute();

const activeSession = computed(() => getSession());
const quizTitle = computed(() => getActiveQuizTitle());
const isAdminUser = computed(
  () => userDataStore.getUserData()?.role === "admin-user"
);
const onArrangePage = computed(() => route.path.startsWith("/admin/arrange/"));

const visible = computed(
  () =>
    Boolean(activeSession.value) && isAdminUser.value && !onArrangePage.value
);

const confirmOpen = ref(false);
const stopping = ref(false);

const resumeQuiz = () => {
  if (!activeSession.value) return;
  router.push(`/admin/arrange/${activeSession.value}`);
};

const openConfirm = () => {
  confirmOpen.value = true;
};

const closeConfirm = () => {
  if (stopping.value) return;
  confirmOpen.value = false;
};

const confirmStop = async () => {
  if (!activeSession.value || stopping.value) return;
  stopping.value = true;
  try {
    await $fetch(`${apiUrl}/quiz/terminate?session_id=${activeSession.value}`, {
      method: "GET",
      credentials: "include",
    });
    setSession(null);
    toast.success("Quiz stopped successfully.");
    confirmOpen.value = false;
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
.overlay-leave-active {
  transition: opacity 0.15s ease;
}
.overlay-enter-from,
.overlay-leave-to {
  opacity: 0;
}
</style>
