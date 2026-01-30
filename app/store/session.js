import { defineStore } from "pinia";
export const useSessionStore = defineStore(
  "session-store",
  () => {
    const session = ref(false);
    const lastComponent = ref("");

    const getSession = () => {
      return session.value;
    };

    const setSession = (data) => {
      session.value = data;
    };

    const getLastComponent = () => {
      return lastComponent.value;
    };

    const setLastComponent = (data) => {
      lastComponent.value = data;
    };

    return {
      session,
      getSession,
      setSession,
      lastComponent,
      getLastComponent,
      setLastComponent,
    };
  },
  {
    persist: true,
  }
);
