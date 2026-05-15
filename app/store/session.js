import { defineStore } from "pinia";
export const useSessionStore = defineStore(
  "session-store",
  () => {
    const session = ref(false);
    const lastComponent = ref("");
    const activeQuizTitle = ref("");

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

    const getActiveQuizTitle = () => activeQuizTitle.value;

    const setActiveQuizTitle = (title) => {
      if (typeof title === "string" && title.trim()) {
        activeQuizTitle.value = title.trim();
      }
    };

    return {
      session,
      getSession,
      setSession,
      lastComponent,
      getLastComponent,
      setLastComponent,
      activeQuizTitle,
      getActiveQuizTitle,
      setActiveQuizTitle,
    };
  },
  {
    persist: true,
  }
);
