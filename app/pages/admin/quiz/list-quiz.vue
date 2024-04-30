<script setup>
const quizList = ref();
const url = useState("urls");
const headers = useRequestHeaders(["cookie"]);
const isLoading = ref(true);

const { data } = await useFetch(url.value.api_url + "/admin/quizzes/list", {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
});
quizList.value = data.value.data;

// remove quiz list loader after 1 sec
setTimeout(() => {
  if (quizList.value) {
    isLoading.value = false;
  }
}, 1000);
</script>
<template>
  <div class="container max-width p-0">
    <div class="d-flex flex-column justify-content-center">
      <!-- list loader -->
      <UtilsQuizListWaiting v-if="isLoading" />

      <!-- quiz details -->
      <div v-else>
        <!-- create quiz if not exists -->
        <div
          v-if="quizList.length < 1"
          class="no-quiz-list d-flex flex-column align-items-center"
        >
          <h1>No Quiz Created By You !</h1>
          <p class="font-italic">Create your first quiz</p>
          <UtilsCreateQuiz />
        </div>

        <!-- show quiz list -->
        <div v-else>
          <!-- Heading -->
          <nav class="navbar pb-4">
            <div class="container-fluid p-0">
              <h1 class="mb-0">Quiz List</h1>
              <UtilsCreateQuiz />
            </div>
          </nav>
          <div class="d-flex flex-column gap-3">
            <div v-for="(details, index) in quizList" :key="index">
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
