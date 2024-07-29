import { defineStore } from "pinia";
export const useUsersStore = defineStore(
  "users-store",
  {
    state: () => ({}),
    getters: {},
    actions: {},
  },
  {
    persist: {
      storage: persistedState.localStorage,
    },
  }
);
