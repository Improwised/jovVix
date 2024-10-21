import { defineStore } from "pinia";
export const useSessionStore = defineStore(
  "session-store",
  () => {
    const session = ref(false);

    const getSession = () => {
      return session.value;
    };

    const setSession = (data) => {
      session.value = data;
    };

    return { session, getSession, setSession };
  },
  {
    persist: {
      storage: persistedState.localStorage,
    },
  }
);
