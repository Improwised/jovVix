<template>
  <div v-if="quizPending">Pending...</div>
  <div v-else-if="quizError?.data?.code == 401">
    {{ navigateTo("/account/login") }}
  </div>
  <div v-else-if="quizError">{{ quizError }}</div>
  <div v-else class="container mt-3">
    <div class="card mb-3 row">
      <div
        class="d-flex flex-wrap align-items-center justify-content-around bg-white rounded p-5"
      >
        <span
          class="badge rounded-pill bg-light-primary text-dark m-1 px-2 fs-4"
        >
          Played Quiz : {{ quizData?.data?.quiz_played_count }}
        </span>
        <span
          class="badge rounded-pill bg-light-primary text-dark m-1 px-2 fs-4"
        >
          Total Questions: {{ quizData?.data?.data.length }}
        </span>
        <span
          class="badge rounded-pill bg-light-primary text-dark m-1 px-2 fs-4"
        >
          Survey Questions: {{ totalSurveyQuestion }}</span
        >

        <!-- Delete quiz button -->
        <button
          v-if="canEditQuiz"
          type="button"
          class="btn btn-outline-danger"
          data-bs-toggle="modal"
          data-bs-target="#deleteQuiz"
        >
          <font-awesome-icon :icon="['fas', 'trash-can']" class="pr-2" />Delete
          Quiz
        </button>
        <DeleteDialog id="deleteQuiz" @confirm-delete="deleteQuiz" />

        <!-- Share quiz buttion -->
        <button
          v-if="quizData?.data?.permission === 'share'"
          class="badge rounded-pill bg-light-info text-dark m-1 px-2 fs-4"
          data-bs-toggle="modal"
          data-bs-target="#shareQuizModal"
          title="Share Quiz"
        >
          <font-awesome-icon :icon="['fas', 'share-from-square']" />
        </button>
        <ShareQuizModal :quiz-id="quizId" @share-quiz="shareQuiz" />
      </div>
    </div>
    <div
      v-for="(quiz, index) in quizData?.data?.data"
      :key="index"
      class="card mb-3 row"
    >
      <QuizQuestionAnalysis
        :question="quiz"
        :order="index + 1"
        :is-admin-analysis="true"
        :is-for-quiz="true"
        :is-editable="canEditQuiz"
        @delete-question="deleteQuestion"
        @edit-question="navagateToEditQuestion"
      />
      <QuizOptionsAnalysis
        :options="quiz?.options"
        :correct-answer="quiz?.correct_answer"
        :options-media="quiz?.options_media"
      />
    </div>
  </div>
</template>

<script setup>
import { useToast } from "vue-toastification";
import DeleteDialog from "~~/components/DeleteDialog.vue";
import ShareQuizModal from "~~/components/Quiz/ShareQuizModal.vue";
const toast = useToast();
const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const route = useRoute();
const quizId = computed(() => route.params.quiz_id || "");
const totalSurveyQuestion = computed(() => {
  const data = quizData.value?.data?.data;
  if (!data) return 0;

  return data.reduce((count, item) => {
    return item.question_type === "survey" ? count + 1 : count;
  }, 0);
});

const canEditQuiz = computed(() => {
  const permission = quizData.value?.data?.permission;
  const isEditable = quizData.value?.data?.is_quiz_editable;
  return (permission === "write" || permission === "share") && isEditable;
});

const {
  refresh,
  data: quizData,
  pending: quizPending,
  error: quizError,
} = useFetch(`${url.apiUrl}/quizzes/${quizId.value}/questions`, {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
});

const navagateToEditQuestion = (questionId) => {
  navigateTo(`/admin/quiz/list-quiz/${quizId.value}/${questionId}`);
};

const deleteQuiz = async () => {
  try {
    await $fetch(`${url.apiUrl}/quizzes/${quizId.value}`, {
      method: "DELETE",
      headers: headers,
      credentials: "include",
    });
    toast.success("Question delete successfully!");
    navigateTo("/admin/quiz/list-quiz");
  } catch (error) {
    console.error("Failed to update the question", error);
    toast.error("Failed to update the question.");
  }
};

const deleteQuestion = async (questionId) => {
  try {
    await $fetch(
      `${url.apiUrl}/quizzes/${quizId.value}/questions/${questionId}`,
      {
        method: "DELETE",
        headers: headers,
        credentials: "include",
      }
    );
    toast.success("Question delete successfully!");
    refresh();
  } catch (error) {
    console.error("Failed to update the question", error);
    toast.error("Failed to update the question.");
  }
};

const shareQuiz = async (email, permission, quizAuthorizedUsersDataRefresh) => {
  try {
    const payload = {
      email: email,
      permission: permission,
    };

    await $fetch(`${url.apiUrl}/shared_quizzes/${quizId.value}`, {
      method: "POST",
      headers: headers,
      body: payload,
      credentials: "include",
    });
    toast.success("Quiz shared successfully!");
    quizAuthorizedUsersDataRefresh();
  } catch (error) {
    console.error("Failed to share the quiz.", error);
    toast.error("Failed to share the quiz.");
  }
};
</script>
