<script setup>
import { useUserThatSubmittedAnswer } from "~/store/userSubmittedAnswer";
import { useListUserstore } from "~~/store/userlist";

const usersThatSubmittedAnswer = useUserThatSubmittedAnswer();
const { usersSubmittedAnswers } = usersThatSubmittedAnswer;

const props = defineProps(["joinedParticipants"]);
</script>

<template>
  <div class="container">
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
            {{ usersSubmittedAnswers.length }}/{{
              props.joinedParticipants
            }}
            People Answered
          </h5>
        </div>
      </div>
    </div>

    <div
      v-if="usersSubmittedAnswers.length"
      class="row justify-content-center mt-5"
    >
      <div class="col-sm-12 col-lg-7 mt-5">
        <div v-for="(user, index) in usersSubmittedAnswers" :key="index">
          <h4
            class="py-3 px-5 border border-1 rounded-pill d-flex justify-content-center"
          >
            {{ user.first_name }} ({{ user.username }})
          </h4>
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
