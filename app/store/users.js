import { defineStore } from "pinia";
export const useUsersStore = defineStore(
  "users-store",
  () => {
    const userData = ref(null);

    const setUserData = (data) => {
      userData.value = data;
    };

    const getUserData = () => {
      return userData.value;
    };

    return { userData, setUserData, getUserData };
  },
  {
    persist: {
      storage: persistedState.localStorage,
    },
  }
);
