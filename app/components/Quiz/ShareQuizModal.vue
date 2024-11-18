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
              <div class="d-flex justify-content-between">
                <h5 class="text-subtitle-1">People with access</h5>
                <button
                  type="buutton"
                  title="Add People"
                  @click="() => (component = 'give-permission')"
                >
                  <font-awesome-icon
                    class="fs-4"
                    :icon="['fas', 'user-plus']"
                  />
                </button>
              </div>
              <div>
                <v-list>
                  <v-list-item
                    v-for="(user, i) in quizAuthorizedUsersData.data"
                    :key="i"
                  >
                    <ShareQuizAuthorizeUser
                      :user="user"
                      @show-edit-form="showEditForm"
                      @delete-user-permission="deleteUserPermission"
                    />
                  </v-list-item>
                </v-list>
              </div>
            </v-card-text>
          </VCard>

          <!-- Form for share quiz permission-->
          <ShareQuizForm
            v-if="component === 'give-permission'"
            form-title="Add People"
            @share-quiz="shareQuiz"
          />

          <!-- Form for update user permission for quiz-->
          <ShareQuizForm
            v-if="component === 'edit-permission'"
            :id="id"
            form-title="Edit Permission"
            :email="email"
            :permission="permission"
            @update-user-permission="updateUserPermission"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useToast } from "vue-toastification";
import ShareQuizAuthorizeUser from "./ShareQuizAuthorizeUser.vue";
import ShareQuizForm from "./ShareQuizForm.vue";
const toast = useToast();
const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);

const id = ref("");
const email = ref("");
const permission = ref("");
const component = ref("");

// define props and emits
const props = defineProps({
  quizId: {
    type: String,
    required: true,
    default: "",
  },
});
const emits = defineEmits(["shareQuiz"]);

// emits the shareQuiz
const shareQuiz = (emailVal, permissionVal) => {
  emits("shareQuiz", emailVal, permissionVal, quizAuthorizedUsersDataRefresh);
};

const showEditForm = (idVal, emailVal, permissionVal) => {
  id.value = idVal;
  email.value = emailVal;
  permission.value = permissionVal;
  component.value = "edit-permission";
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

const updateUserPermission = async (idVal, emailVal, permissionVal) => {
  try {
    const payload = {
      email: emailVal,
      permission: permissionVal,
    };

    await $fetch(
      `${url.api_url}/shared_quizzes/${props.quizId}?shared_quiz_id=${idVal}`,
      {
        method: "PUT",
        headers: headers,
        body: payload,
        credentials: "include",
      }
    );
    toast.success("User permission updated successfully!");
    quizAuthorizedUsersDataRefresh();
    component.value = "give-permission";
  } catch (error) {
    console.error("Failed to update user permission.", error);
    toast.error("Failed to update user permission.");
  }
};

const deleteUserPermission = async (idVal) => {
  try {
    await $fetch(
      `${url.api_url}/shared_quizzes/${props.quizId}?shared_quiz_id=${idVal}`,
      {
        method: "DELETE",
        headers: headers,
        credentials: "include",
      }
    );
    toast.success("User permission deleted successfully!");
    quizAuthorizedUsersDataRefresh();
  } catch (error) {
    console.error("Failed to delete user permission.", error);
    toast.error("Failed to delete user permission.");
  }
};
</script>
