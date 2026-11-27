<template>
  <section
    v-if="show"
    class="jv-border-rough bg-jv-yellow/30 p-4 sm:p-6 shadow-brutal-sm"
    aria-labelledby="first-quiz-heading"
  >
    <h2
      id="first-quiz-heading"
      class="font-headings text-[24px] sm:text-[30px] leading-tight text-jv-ink"
    >
      You're in! Start your first quiz
    </h2>
    <p class="mt-1 text-[15px] sm:text-[17px] text-jv-muted">
      Players join on their phones with a code or a QR scan.
    </p>
    <div class="mt-4 flex flex-col gap-3 sm:flex-row">
      <a
        href="#public-quiz"
        class="inline-flex items-center justify-center gap-2 font-headings jv-border-uneven shadow-brutal transition-transform hover:rotate-[2deg] active:shadow-none bg-jv-coral text-white px-5 py-2.5 text-base"
      >
        <Play class="size-5" :stroke-width="2.5" />
        Host a ready-made quiz now
      </a>
      <NuxtLink
        to="/admin/quiz/list-quiz?create=1"
        class="inline-flex items-center justify-center gap-2 font-headings jv-border-uneven shadow-brutal transition-transform hover:rotate-[2deg] active:shadow-none bg-jv-white text-jv-ink px-5 py-2.5 text-base"
      >
        <Plus class="size-5" :stroke-width="2.5" />
        Write your own questions
      </NuxtLink>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { Play, Plus } from "lucide-vue-next";
import { takeNewSignup } from "@/composables/first_quiz_welcome";
import { useUsersStore } from "~~/store/users";

const usersStore = useUsersStore();
const justSignedUp = ref(false);

// The sidebar fills the user store after mount, so wait for it reactively.
const show = computed(() => justSignedUp.value && !!usersStore.getUserData());

onMounted(() => {
  justSignedUp.value = takeNewSignup();
});
</script>
