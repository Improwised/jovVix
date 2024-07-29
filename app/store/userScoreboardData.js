import { defineStore } from "pinia";

export const useUserScoreboardData = defineStore(
  "userscoreboard-store",
  () => {
    const userScoreboardData = ref([]);

    const addData = (data) => {
      userScoreboardData.value.push(...data);
    };

    const getUserScoreboardData = () => {
      return userScoreboardData.value;
    };

    const resetStore = () => {
      userScoreboardData.value = []; // Reset to null or empty object
    };

    return { userScoreboardData, addData, getUserScoreboardData, resetStore };
  },
  {
    persist: {
      storage: persistedState.localStorage,
    },
  }
);
