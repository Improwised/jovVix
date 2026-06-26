<script setup>
import { Share2 } from "lucide-vue-next";
import NavigationButton from "~~/components/utils/NavigationButton.vue";

definePageMeta({
  layout: "empty",
});

const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);

const {
  data: quizList,
  pending: quizPending,
  error: quizError,
} = useFetch(url.apiUrl + "/shared_quizzes?type=shared_by_me", {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
});
</script>

<template>
  <div class="min-h-screen bg-jv-canvas px-4 py-8 sm:px-6 sm:py-10">
    <div class="mx-auto max-w-[922px]">
      <UtilsQuizListWaiting v-if="quizPending" />

      <div
        v-else-if="quizError"
        class="jv-border-rough border-2 border-jv-coral bg-jv-coral/10 p-4 font-body text-jv-ink"
      >
        {{ quizError.message }}
      </div>

      <template v-else>
        <div
          v-if="quizList?.data.length < 1"
          class="jv-border-rough flex flex-col items-center gap-4 border-2 border-jv-ink bg-jv-white p-8 text-center shadow-brutal sm:p-12"
        >
          <div
            class="grid size-16 place-items-center rounded-[12px] border-2 border-jv-ink bg-jv-mint shadow-brutal-sm"
          >
            <Share2 class="size-8 text-jv-ink" :stroke-width="2.4" />
          </div>
          <h1
            class="font-headings text-[26px] leading-tight text-jv-ink sm:text-[32px]"
          >
            No Quiz Shared By You!
          </h1>
          <p class="font-body text-sm italic text-jv-muted sm:text-base">
            Share your first quiz.
          </p>
          <NavigationButton
            title="Share Quiz"
            navigate-to="/admin/quiz/list-quiz"
          />
        </div>

        <div v-else class="flex flex-col gap-4">
          <header class="mb-2">
            <h1
              class="font-headings text-[28px] leading-tight text-jv-ink sm:text-[36px]"
            >
              Quiz Shared By Me
            </h1>
          </header>
          <QuizListCard
            v-for="(details, index) in quizList?.data"
            :key="index"
            :details="details"
          />
        </div>
      </template>
    </div>
  </div>
</template>
