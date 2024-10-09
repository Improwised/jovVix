<script setup>
import { useUsersStore } from "~~/store/users";
import { getAvatarUrlByName } from "~~/composables/avatar";
const userData = useUsersStore();
const { getUserData } = userData;

const props = defineProps({
  userName: {
    type: String,
    required: false,
    default: "",
  },
});

const avatar = computed(() => {
  const user = getUserData();
  return getAvatarUrlByName(user?.avatar);
});
</script>

<template>
  <div class="container-fluid mb-5">
    <div class="row justify-content-center">
      <div class="col-6 col-md-4 mt-5 d-flex justify-content-center">
        <div
          class="border border-1 px-7 py-2 py-md-4 border-radius d-flex align-items-center"
        >
          <img :src="avatar" height="70" width="70" class="me-3" />
          <h5 class="text-center text-sm fs-5 mb-0">{{ props.userName }}</h5>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.border-radius {
  border-radius: 2rem !important;
}
</style>
