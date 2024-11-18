<script setup>
import NavigationButton from "~~/components/utils/NavigationButton.vue";
const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);

const {
  data: quizList,
  pending: quizPending,
  error: quizError,
} = useFetch(url.api_url + "/shared_quizzes?type=shared_by_me", {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
});
</script>
<template>
  <div class="container max-width p-0">
    <div class="d-flex flex-column justify-content-center">
      <!-- list loader -->
      <UtilsQuizListWaiting v-if="quizPending" />

      <div v-else-if="quizError">{{ quizError.message }}</div>

      <!-- quiz details -->
      <div v-else>
        <div
          v-if="quizList?.data.length < 1"
          class="no-quiz-list d-flex flex-column align-items-center"
        >
          <h1>No Quiz Shared By You !</h1>
          <p class="font-italic">Share Your First Quiz</p>
          <NavigationButton
            :title="'Share Quiz'"
            :navigate-to="'/admin/quiz/list-quiz'"
          />
        </div>

        <!-- show quiz list -->
        <div v-else>
          <!-- Heading -->
          <nav class="navbar pb-4">
            <div class="container-fluid p-0">
              <h1 class="mb-0">Quiz List</h1>
            </div>
          </nav>
          <div class="d-flex flex-column gap-3">
            <div v-for="(details, index) in quizList?.data" :key="index">
              <QuizListCard :details="details" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<style scoped>
.max-width {
  max-width: 922px;
}
</style>
