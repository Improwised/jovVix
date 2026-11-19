<script setup>
const route = useRoute();
const roomPhase = ref("");
definePageMeta({ layout: "empty", hideSidebar: true });
</script>
<template>
  <main
    class="min-h-screen bg-jv-canvas px-3 py-6 text-jv-ink sm:px-6 md:px-8 md:py-10"
  >
    <header class="mx-auto mb-8 max-w-[1220px] text-center sm:mb-10">
      <p
        class="font-body text-xs font-black uppercase tracking-[0.18em] text-jv-muted sm:text-sm"
      >
        Quiz room
      </p>
      <h1
        class="relative mt-3 inline-block font-headings text-[32px] leading-tight sm:text-[48px] md:text-[56px]"
      >
        <span class="relative z-10">{{
          roomPhase === "closed"
            ? "Thanks for playing!"
            : roomPhase
              ? "Ready for another round?"
              : "Quiz room"
        }}</span>
        <span
          class="absolute bottom-1 left-0 h-2 w-full max-w-[260px] -rotate-1 bg-jv-yellow/70"
          aria-hidden="true"
        ></span>
      </h1>
      <p
        class="mx-auto mt-3 max-w-2xl text-base leading-relaxed text-jv-muted sm:text-lg"
      >
        {{
          roomPhase === "closed"
            ? "The room has ended. Review your final standings and played quizzes."
            : roomPhase
              ? "Stay in the room and play the next quiz together."
              : "Connecting to your room…"
        }}
      </p>
    </header>
    <RoomPanel
      :room-id="String(route.params.id)"
      standalone
      @room-state="roomPhase = $event"
    />
  </main>
</template>
