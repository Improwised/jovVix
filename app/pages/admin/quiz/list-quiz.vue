<script setup>
const quizList = ref(false);
</script>
<template>
  <div class="container">
    <div class="d-flex flex-column justify-content-center">
      <!-- Heading -->
      <nav class="navbar pb-4">
        <div class="container-fluid">
          <h1>Quiz List</h1>
          <UtilsCreateQuiz />
        </div>
      </nav>
      <!-- list loader -->
      <UtilsQuizListWaiting />
      <!-- create quiz if not list exists -->
      <div
        v-if="quizList"
        class="no-quiz-list d-flex flex-column align-items-center"
      >
        <p class="fw-semibold">Create your first quiz</p>
        <UtilsCreateQuiz />
      </div>
    </div>
  </div>
</template>
<style scoped>
.max-width {
  max-width: 1320px;
}
</style>
