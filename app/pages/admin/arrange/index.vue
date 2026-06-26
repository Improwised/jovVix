<script setup>
import { Home, Library, AlertCircle } from "lucide-vue-next";

definePageMeta({
  layout: "empty",
  hideSidebar: true,
});

const route = useRoute();

// Admins land here from /admin/arrange/[session_id] when the websocket reports
// that the session was completed, the invitation code was missing, or session
// validation failed. The redirect tacks on ?status=...&error=... so we can show
// a context-specific message rather than a generic "temporary redirection".
const presentation = computed(() => {
  const err = String(route.query.error || "").toLowerCase();

  if (err === "session was completed") {
    return {
      title: "Quiz Session Has Ended",
      message:
        "This quiz session is already finished. You can't return to the host controls once a quiz wraps up — head home to start a new one or review results from the dashboard.",
    };
  }
  if (err === "quiz-session-validation-failed") {
    return {
      title: "Quiz Session No Longer Active",
      message:
        "This quiz session can no longer be hosted. It may have ended, been terminated, or expired. Start a new session to host again.",
    };
  }
  if (err === "invitation code not found") {
    return {
      title: "Invitation Code No Longer Valid",
      message:
        "The invitation code for this session is no longer active. The session has likely ended.",
    };
  }
  return {
    title: "No Active Quiz Session",
    message:
      "There's nothing to host here right now. Head back home, or jump into your quizzes to start a new session.",
  };
});

useHead({
  title: "Quiz Session Ended - Admin",
  meta: [
    {
      name: "description",
      content:
        "The quiz session you were hosting has ended. Return home or visit your dashboard to start a new quiz.",
    },
  ],
});
</script>

<template>
  <div
    class="relative min-h-screen overflow-hidden bg-jv-canvas px-4 py-10 sm:px-6 sm:py-14"
  >
    <div
      class="pointer-events-none absolute inset-0 jv-grid"
      aria-hidden="true"
    ></div>

    <div
      class="relative mx-auto flex w-full max-w-[640px] flex-col items-center"
    >
      <div
        class="relative w-full -rotate-[0.6deg] jv-border-rough bg-jv-white p-6 shadow-brutal-lg sm:p-9"
      >
        <span
          class="absolute left-1/2 top-[-12px] z-10 h-4 w-16 -translate-x-1/2 rotate-[2deg] bg-jv-coral"
          aria-hidden="true"
        ></span>

        <div class="flex items-start gap-4">
          <span
            class="grid size-12 shrink-0 rotate-[-3deg] place-items-center rounded-[10px] border-[2px] border-jv-ink bg-jv-yellow text-jv-ink shadow-brutal-sm sm:size-14"
            aria-hidden="true"
          >
            <AlertCircle class="size-6 sm:size-7" :stroke-width="2.4" />
          </span>
          <div class="min-w-0 flex-1">
            <p
              class="font-body text-[11px] font-black uppercase tracking-[0.22em] text-jv-muted sm:text-[12px]"
            >
              Admin
            </p>
            <h1
              class="mt-1 font-headings text-[26px] leading-tight text-jv-ink sm:text-[32px]"
            >
              {{ presentation.title }}
            </h1>
          </div>
        </div>

        <hr class="my-5 border-t-[2px] border-dashed border-jv-ink/25" />

        <p
          class="font-body text-[14px] font-bold leading-relaxed text-jv-muted sm:text-[15px]"
        >
          {{ presentation.message }}
        </p>

        <div
          class="mt-6 flex flex-col items-stretch gap-3 sm:flex-row sm:justify-center"
          role="group"
          aria-label="Recovery navigation"
        >
          <NuxtLink
            to="/"
            class="inline-flex h-11 items-center justify-center gap-2 rounded-full border-[2px] border-jv-ink bg-jv-coral px-6 font-body text-[14px] font-black text-white shadow-brutal-sm transition-transform hover:-rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-12 sm:text-[15px]"
            aria-label="Return to home"
          >
            <Home class="size-4" :stroke-width="2.4" />
            Back to Home
          </NuxtLink>
          <NuxtLink
            to="/admin/quiz/list-quiz"
            class="inline-flex h-11 items-center justify-center gap-2 rounded-full border-[2px] border-jv-ink bg-jv-white px-6 font-body text-[14px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-12 sm:text-[15px]"
            aria-label="Go to my quizzes"
          >
            <Library class="size-4" :stroke-width="2.4" />
            My Quizzes
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.jv-grid {
  background-image: radial-gradient(var(--jv-charcoal) 1px, transparent 1px);
  background-size: 20px 20px;
  opacity: 0.08;
}
</style>
