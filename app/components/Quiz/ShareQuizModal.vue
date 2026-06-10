<script setup>
import { computed, ref, watch } from "vue";
import { usePush } from "notivue";
import { UserPlus, Users } from "lucide-vue-next";
import { Modal } from "@/components/ui/modal";
import ShareQuizAuthorizeUser from "./ShareQuizAuthorizeUser.vue";
import ShareQuizForm from "./ShareQuizForm.vue";

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  quizId: {
    type: String,
    required: true,
    default: "",
  },
});

const emits = defineEmits(["update:modelValue"]);

const toast = usePush();
const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);

// edit-form state
const id = ref("");
const email = ref("");
const permission = ref("");
const component = ref("");

// Authorized users for this quiz — reactive to quizId so the modal works
// from both the quiz-detail page and the per-quiz list cards.
const {
  refresh: quizAuthorizedUsersDataRefresh,
  data: quizAuthorizedUsersData,
  pending: quizAuthorizedUsersPending,
  error: quizAuthorizedUsersError,
} = useFetch(() => `${url.apiUrl}/shared_quizzes/${props.quizId}`, {
  method: "GET",
  headers,
  mode: "cors",
  credentials: "include",
  immediate: false,
  watch: [() => props.quizId],
});

const authorizedUsers = computed(
  () => quizAuthorizedUsersData.value?.data || []
);

function close() {
  emits("update:modelValue", false);
}

// Reset the form whenever the modal opens, and (re)load the user list.
watch(
  () => props.modelValue,
  (open) => {
    if (open) {
      component.value = "";
      id.value = "";
      email.value = "";
      permission.value = "";
      if (props.quizId) quizAuthorizedUsersDataRefresh();
    }
  }
);

const showEditForm = (idVal, emailVal, permissionVal) => {
  id.value = idVal;
  email.value = emailVal;
  permission.value = permissionVal;
  component.value = "edit-permission";
};

const shareQuiz = async (emailVal, permissionVal) => {
  try {
    await $fetch(`${url.apiUrl}/shared_quizzes/${props.quizId}`, {
      method: "POST",
      headers,
      body: { email: emailVal, permission: permissionVal },
      credentials: "include",
    });
    toast.success("Quiz shared successfully!");
    quizAuthorizedUsersDataRefresh();
    component.value = "";
  } catch (error) {
    console.error("Failed to share the quiz.", error);
    toast.error("Failed to share the quiz.");
  }
};

const updateUserPermission = async (idVal, emailVal, permissionVal) => {
  try {
    await $fetch(
      `${url.apiUrl}/shared_quizzes/${props.quizId}?shared_quiz_id=${idVal}`,
      {
        method: "PUT",
        headers,
        body: { email: emailVal, permission: permissionVal },
        credentials: "include",
      }
    );
    toast.success("User permission updated successfully!");
    quizAuthorizedUsersDataRefresh();
    component.value = "";
  } catch (error) {
    console.error("Failed to update user permission.", error);
    toast.error("Failed to update user permission.");
  }
};

const deleteUserPermission = async (idVal) => {
  try {
    await $fetch(
      `${url.apiUrl}/shared_quizzes/${props.quizId}?shared_quiz_id=${idVal}`,
      {
        method: "DELETE",
        headers,
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

<template>
  <Modal
    :model-value="props.modelValue"
    title="Share Quiz"
    description="Give people access to this quiz and manage their permissions."
    size="md"
    @update:model-value="emits('update:modelValue', $event)"
  >
    <div
      class="mb-4 h-[3px] w-full rounded-full border-b-[3px] border-dashed border-jv-ink/30"
      aria-hidden="true"
    ></div>

    <div v-if="quizAuthorizedUsersPending" class="py-6 text-center">
      <p class="text-[15px] font-semibold text-jv-muted">Loading...</p>
    </div>

    <div
      v-else-if="quizAuthorizedUsersError"
      class="jv-border-rough border-[2px] border-jv-ink bg-jv-white p-4 text-[15px] font-semibold text-jv-coral"
    >
      {{ quizAuthorizedUsersError.message || quizAuthorizedUsersError }}
    </div>

    <template v-else>
      <!-- People with access -->
      <div class="flex items-center justify-between gap-3">
        <h3
          class="flex items-center gap-2 font-body text-[15px] font-black uppercase tracking-[0.14em] text-jv-ink"
        >
          <Users class="size-4" :stroke-width="2.6" />
          People with access
        </h3>
        <button
          type="button"
          title="Add People"
          aria-label="Add people"
          class="inline-flex items-center gap-2 rounded-full border-[2px] border-jv-ink bg-jv-white px-3 py-1.5 text-[13px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-[2deg] active:translate-x-[1px] active:translate-y-[1px] active:shadow-none"
          @click="component = 'give-permission'"
        >
          <UserPlus class="size-4" :stroke-width="2.6" />
          Add
        </button>
      </div>

      <ul class="mt-4 flex max-h-[320px] flex-col gap-3 overflow-y-auto pr-1">
        <li v-for="(user, i) in authorizedUsers" :key="i">
          <ShareQuizAuthorizeUser
            :user="user"
            @show-edit-form="showEditForm"
            @delete-user-permission="deleteUserPermission"
          />
        </li>
        <li
          v-if="authorizedUsers.length === 0"
          class="py-4 text-center text-[14px] font-semibold text-jv-muted"
        >
          No one has access yet. Add someone to get started.
        </li>
      </ul>

      <!-- Add-permission form -->
      <ShareQuizForm
        v-if="component === 'give-permission'"
        form-title="Add People"
        @share-quiz="shareQuiz"
      />

      <!-- Edit-permission form -->
      <ShareQuizForm
        v-if="component === 'edit-permission'"
        :id="id"
        form-title="Edit Permission"
        :email="email"
        :permission="permission"
        @update-user-permission="updateUserPermission"
      />
    </template>

    <template #footer>
      <button
        type="button"
        class="inline-flex h-10 items-center justify-center rounded-full border-[2px] border-jv-ink bg-jv-white px-5 font-body text-[14px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-11 sm:px-6 sm:text-[15px]"
        @click="close"
      >
        Done
      </button>
    </template>
  </Modal>
</template>
