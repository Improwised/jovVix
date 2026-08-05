import { defineStore } from "pinia";
import { useUsersStore } from "~~/store/users";

const SYNC_THROTTLE_MS = 30000;

export const useSessionStore = defineStore(
  "session-store",
  () => {
    const session = ref(false);
    const lastComponent = ref("");
    const activeQuizTitle = ref("");
    const activeSessions = ref([]);
    const sessionsSyncedAt = ref(0);
    const sessionsLoading = ref(false);

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

    const getActiveSessions = () => activeSessions.value;

    const removeActiveSession = (id) => {
      activeSessions.value = activeSessions.value.filter((s) => s.id !== id);
      if (session.value === id) {
        session.value = null;
      }
    };

    const syncActiveSessions = async ({ force = false } = {}) => {
      if (!import.meta.client) return;

      if (!useUsersStore().getUserData()) {
        activeSessions.value = [];
        return;
      }

      if (!force && Date.now() - sessionsSyncedAt.value < SYNC_THROTTLE_MS) {
        return;
      }

      if (sessionsLoading.value) return;

      sessionsLoading.value = true;
      const { apiUrl } = useRuntimeConfig().public;

      try {
        const res = await $fetch(`${apiUrl}/quiz/sessions/active`, {
          method: "GET",
          credentials: "include",
        });
        activeSessions.value = Array.isArray(res?.data) ? res.data : [];
        sessionsSyncedAt.value = Date.now();
      } catch (error) {
        const status = error?.response?.status || error?.statusCode;
        if (status === 401 || status === 403) {
          activeSessions.value = [];
          sessionsSyncedAt.value = Date.now();
        }
      } finally {
        sessionsLoading.value = false;
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
      activeSessions,
      sessionsLoading,
      getActiveSessions,
      removeActiveSession,
      syncActiveSessions,
    };
  },
  {
    persist: { pick: ["activeQuizTitle"] },
  }
);
