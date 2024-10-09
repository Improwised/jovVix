<script setup>
import { useUserThatSubmittedAnswer } from "~/store/userSubmittedAnswer";
import { useNuxtApp } from "nuxt/app";
import { useToast } from "vue-toastification";
const app = useNuxtApp();
const toast = useToast();

const usersThatSubmittedAnswer = useUserThatSubmittedAnswer();
const { usersSubmittedAnswers } = usersThatSubmittedAnswer;
const totalUser = ref(0);

const props = defineProps({
  data: {
    default: () => {
      return {};
    },
    type: Object,
    required: true,
  },
});

watch(
  () => props.data,
  (message) => {
    if (message.status == app.$Fail) {
      toast.error(message.data);
      return;
    }
    handleCountUser(message);
  },
  { deep: true, immediate: true }
);

// main function
function handleCountUser(message) {
  if (message.event == app.$GetQuestion) {
    totalUser.value = message.data.totalJoinUser;
  }
}
</script>

<template>
  <div class="container" style="max-width: 800px">
    <div class="row justify-content-center">
      <div v-if="usersSubmittedAnswers.length == 0" class="col-7 col-md-4 mt-5">
        <div
          class="d-flex border border-1 justify-content-center align-items-center px-3 py-2 py-md-4 gap-3 border-radius"
        >
          <font-awesome-icon icon="fa-solid fa-users" size="xl" />
          <h5 class="text-center mb-0">No One Answered Till Now..</h5>
        </div>
      </div>

      <div v-else class="col-6 col-md-4 mt-5">
        <div
          class="d-flex border border-1 justify-content-center align-items-center px-3 py-2 py-md-4 gap-3 border-radius"
        >
          <font-awesome-icon icon="fa-solid fa-users" size="xl" />
          <h5 class="text-center text-sm fs-5 mb-0">
            {{ usersSubmittedAnswers.length }}/{{ totalUser }}
            People Answered
          </h5>
        </div>
      </div>
    </div>

    <v-card
      v-if="usersSubmittedAnswers.length"
      :flat="true"
      class="mb-2 d-flex flex-wrap justify-content-center"
      color="#00000000"
    >
      <div
        v-for="user in usersSubmittedAnswers"
        :key="user.UserId"
        class="chip m-2"
      >
        <img
          :src="getAvatarUrlByName(user?.img_key)"
          alt="Person"
          width="96"
          height="96"
        />
        {{ user.first_name }} ({{ user.username }})
      </div>
    </v-card>
  </div>
</template>

<style scoped>
.border-radius {
  border-radius: 2rem !important;
}

.chip {
  display: inline-block;
  padding: 0 25px;
  height: 50px;
  font-size: 16px;
  line-height: 50px;
  border-radius: 25px;
  max-width: 600px;
  background-color: #f1f1f1;
}

.chip img {
  float: left;
  margin: 0 10px 0 -25px;
  height: 50px;
  width: 50px;
  border-radius: 50%;
}
</style>
