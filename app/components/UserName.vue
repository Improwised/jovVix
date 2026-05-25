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
  return user?.avatar ? getAvatarUrlByName(user?.avatar) : "";
});
</script>

<template>
  <div class="mb-8 flex justify-center">
    <div
      class="jv-card flex items-center gap-3 border-2 border-jv-ink bg-jv-white py-1.5 pl-1.5 pr-5 shadow-brutal-sm"
    >
      <img
        :src="avatar || 'https://api.dicebear.com/9.x/bottts/svg?seed=Eden'"
        :alt="props.userName"
        class="size-[60px] shrink-0 rounded-full border-2 border-jv-ink object-cover"
      />
      <h5 class="m-0 font-headings text-base text-jv-ink sm:text-lg">
        {{ props.userName }}
      </h5>
    </div>
  </div>
</template>
