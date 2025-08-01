import { defineStore } from "pinia";

export const useUserThatSubmittedAnswer = defineStore(
  "users-that-submitted-answers-store",
  () => {
    const usersSubmittedAnswers = ref([]);

    const addUserSubmittedAnswer = (user) => {
      if (usersSubmittedAnswers.value.includes(user)) {
        return;
      } else {
        usersSubmittedAnswers.value.push(user);
      }
    };

    const resetUsersSubmittedAnswers = () => {
      usersSubmittedAnswers.value = [];
    };

    return {
      usersSubmittedAnswers,
      addUserSubmittedAnswer,
      resetUsersSubmittedAnswers,
    };
  },
  {
    persist: true,
  }
);
