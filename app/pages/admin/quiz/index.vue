<template>
  <div
    class="d-flex justify-content-center align-items-center flex-wrap quiz-container"
  >
    <ActionCard
      action-title="Create New Quiz"
      @card-click="navigateTo('/admin/quiz/create-quiz')"
    >
      <template #image>
        <img src="@/assets/images/add-file.gif" class="card-image" />
      </template>
    </ActionCard>
    <ActionCard
      action-title="My Quizzes"
      @card-click="navigateTo('/admin/quiz/list-quiz')"
    >
      <template #image>
        <img src="@/assets/images/list.gif" class="card-image" />
      </template>
    </ActionCard>
  </div>
</template>
<script setup>
import ActionCard from "~~/components/Quiz/ActionCard.vue";
</script>

<style scoped>
.card-image {
  height: 250px;
  width: 200px;
  border: none;
}

.quiz-container {
  min-height: 90vh;
}
</style>
