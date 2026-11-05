import { defineStore } from "pinia";

export const useUserThatSubmittedAnswer = defineStore(
  "users-that-submitted-answers-store",
  () => {
    const usersSubmittedAnswers = ref([]);

    const addUserSubmittedAnswer = (user) => {
      const id = user?.id || user?.UserId || user?.user_id;
      if (
        !id ||
        usersSubmittedAnswers.value.some(
          (entry) => (entry?.id || entry?.UserId || entry?.user_id) === id
        )
      )
        return;
      usersSubmittedAnswers.value.push(user);
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
