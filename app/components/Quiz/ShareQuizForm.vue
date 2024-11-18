<template>
  <h5 class="text-subtitle-1">{{ props.formTitle }}</h5>
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
        :disabled="id !== '' ? true : false"
      />
    </div>
    <!-- Permission -->
    <div class="mb-3">
      <label for="permission" class="form-label">Permission</label>
      <select id="permission" v-model="permission" class="form-select" required>
        <option value="" disabled>Select permission level</option>
        <option value="read">Read</option>
        <option value="write">Write</option>
        <option value="share">Share</option>
      </select>
    </div>
    <div>
      <!-- Button -->
      <div class="d-grid">
        <button v-if="id" type="submit" class="btn btn-primary text-white">
          Update Access
        </button>
        <button v-else type="submit" class="btn btn-primary text-white">
          Share Quiz
        </button>
      </div>
    </div>
  </form>
</template>

<script setup>
// define props and emits
const props = defineProps({
  formTitle: {
    type: String,
    required: true,
    default: "",
  },
  id: {
    type: String,
    required: false,
    default: "",
  },
  email: {
    type: String,
    required: false,
    default: "",
  },
  permission: {
    type: String,
    required: false,
    default: "",
  },
});
const emits = defineEmits(["shareQuiz", "updateUserPermission"]);

const email = ref(props.email);
const permission = ref(props.permission);

// Watch for changes in props and update local state
watch(
  () => props.email,
  (newEmail) => {
    email.value = newEmail;
  }
);

watch(
  () => props.permission,
  (newPermission) => {
    permission.value = newPermission;
  }
);

const handleSubmit = () => {
  if (props.id) {
    emits("updateUserPermission", props.id, email.value, permission.value);
  } else {
    emits("shareQuiz", email.value, permission.value);
  }

  email.value = "";
  permission.value = "";
};
</script>
