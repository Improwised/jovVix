import { defineStore } from "pinia";

const KEY_STORAGE_ID = "ai-vault-password";

const emptySettings = () => ({
  provider: "",
  baseUrl: "",
  apiKey: "",
  configured: false,
  model: "",
  maskedApiKey: "",
});

const readStoredKey = () => {
  if (import.meta.server) return "";
  try {
    return sessionStorage.getItem(KEY_STORAGE_ID) ?? "";
  } catch {
    return "";
  }
};

const writeStoredKey = (apiKey) => {
  if (import.meta.server) return;
  try {
    if (apiKey) sessionStorage.setItem(KEY_STORAGE_ID, apiKey);
    else sessionStorage.removeItem(KEY_STORAGE_ID);
  } catch {
    // a browser with storage disabled keeps the key in memory for this page only
  }
};

export const useAiSettingsStore = defineStore(
  "ai-settings-store",
  () => {
    const settings = ref(emptySettings());

    // Only the Vault Password lives for this browser session; the provider key
    // is encrypted in the database and never persisted by this store.
    const restoreApiKey = () => {
      settings.value.apiKey = readStoredKey();
    };

    const fetchSettings = async () => {
      const storedKey = readStoredKey();
      const headers = storedKey ? { "X-AI-Vault-Password": storedKey } : {};
      const response = await $fetch(
        `${useRuntimeConfig().public.apiUrl}/ai/settings`,
        { credentials: "include", headers }
      );
      settings.value = { ...emptySettings(), ...(response?.data || {}) };
    };
    const setSettings = async (next) => {
      const headers = {
        "X-AI-Api-Key": next.apiKey,
        "X-AI-Vault-Password": next.vaultPassword,
      };
      const response = await $fetch(
        `${useRuntimeConfig().public.apiUrl}/ai/settings`,
        {
          method: "PUT",
          credentials: "include",
          headers,
          body: {
            provider: next.provider,
            baseUrl: next.baseUrl,
            model: next.model,
          },
        }
      );
      settings.value = { ...emptySettings(), ...(response?.data || {}) };
      writeStoredKey(next.vaultPassword);
      settings.value.apiKey = next.vaultPassword;
    };
    const unlock = async (password) => {
      await $fetch(`${useRuntimeConfig().public.apiUrl}/ai/settings/unlock`, {
        method: "POST",
        credentials: "include",
        headers: { "X-AI-Vault-Password": password },
      });
      settings.value.apiKey = password;
      writeStoredKey(password);
    };
    const deleteSettings = async () => {
      await $fetch(`${useRuntimeConfig().public.apiUrl}/ai/settings`, {
        method: "DELETE",
        credentials: "include",
      });
      settings.value = emptySettings();
      writeStoredKey("");
    };

    const clear = () => {
      settings.value = emptySettings();
      writeStoredKey("");
    };

    const isConfigured = computed(
      () =>
        !!settings.value.configured &&
        !!settings.value.baseUrl &&
        !!settings.value.model
    );

    const aiHeaders = () => {
      if (!isConfigured.value) return {};

      const headers = {
        "X-AI-Vault-Password": settings.value.apiKey,
      };
      return headers;
    };

    return {
      settings,
      setSettings,
      unlock,
      deleteSettings,
      clear,
      restoreApiKey,
      fetchSettings,
      isConfigured,
      aiHeaders,
    };
  },
  {
    persist: {
      pick: ["settings.provider", "settings.baseUrl", "settings.model"],
    },
  }
);
