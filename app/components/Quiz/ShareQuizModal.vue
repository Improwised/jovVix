<template>
  <div
    id="shareQuizModal"
    class="modal fade"
    tabindex="-1"
    aria-labelledby="exampleModalLabel"
    aria-hidden="true"
  >
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header z-2">
          <h1 id="exampleModalLabel" class="modal-title fs-3">Share Quiz</h1>
          <button
            id="share-quiz-btn-close"
            type="button"
            class="btn-close"
            data-bs-dismiss="modal"
            aria-label="Close"
          ></button>
        </div>
        <div class="modal-body">
          <div v-if="quizAuthorizedUsersPending">Pending...</div>
          <div v-else-if="quizAuthorizedUsersError">
            {{ quizAuthorizedUsersError }}
          </div>

          <!-- Authorized Users with Permission -->
          <VCard v-else elevation="0" class="overflow-hidden">
            <v-card-text class="pa-0">
              <h5 class="text-subtitle-1">People with access</h5>
              <div>
                <v-list>
                  <v-list-item
                    v-for="(user, i) in quizAuthorizedUsersData.data"
                    :key="i"
                  >
                    <ShareQuizAuthorizeUser :user="user" />
                  </v-list-item>
                </v-list>
              </div>
            </v-card-text>
          </VCard>

          <!-- Form for share quiz permission-->
          <h5 class="text-subtitle-1">Add People</h5>
          <form @submit.prevent="handleSubmit">
            <!-- Email -->
            <div class="mb-3">
              <label for="email" class="form-label">Email</label>
              <input
                id="email"
                v-model="email"
                type="email"
                name="identifier"
                class="form-control"
                required
              />
            </div>
            <!-- Permission -->
            <div class="mb-3">
              <label for="permission" class="form-label">Permission</label>
              <select
                id="permission"
                v-model="permission"
                class="form-select"
                required
              >
                <option value="" disabled>Select permission level</option>
                <option value="read">Read</option>
                <option value="write">Write</option>
                <option value="share">Share</option>
              </select>
            </div>
            <div>
              <!-- Button -->
              <div class="d-grid">
                <button type="submit" class="btn btn-primary text-white">
                  Submit
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import ShareQuizAuthorizeUser from "./ShareQuizAuthorizeUser.vue";
const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);

const email = ref("");
const permission = ref("");

// define props and emits
const props = defineProps({
  quizId: {
    type: String,
    required: true,
    default: "",
  },
});
const emits = defineEmits(["shareQuiz"]);

// emits the shareQuiz and close the modal
const handleSubmit = () => {
  emits(
    "shareQuiz",
    email.value,
    permission.value,
    quizAuthorizedUsersDataRefresh
  );
  email.value = "";
  permission.value = "";

  // close the modal
  var closeModalButton = document.getElementById("share-quiz-btn-close");
  closeModalButton.click();
};

// Get authorized users data for perticular quiz
const {
  refresh: quizAuthorizedUsersDataRefresh,
  data: quizAuthorizedUsersData,
  pending: quizAuthorizedUsersPending,
  error: quizAuthorizedUsersError,
} = useFetch(`${url.api_url}/shared_quizzes/${props.quizId}`, {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
});
</script>
