import { defineStore } from "pinia";

const KEY_STORAGE_ID = "ai-settings-key";

const emptySettings = () => ({
  provider: "",
  baseUrl: "",
  apiKey: "",
  model: "",
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

    // The key lives in sessionStorage so it does not outlive the browser session,
    // and is left out of the persisted state so it never reaches localStorage.
    const restoreApiKey = () => {
      settings.value.apiKey = readStoredKey();
    };

    const setSettings = (next) => {
      settings.value = { ...emptySettings(), ...next };
      writeStoredKey(settings.value.apiKey);
    };

    const clear = () => {
      settings.value = emptySettings();
      writeStoredKey("");
    };

    const isConfigured = computed(
      () => !!settings.value.baseUrl && !!settings.value.model
    );

    const aiHeaders = () => {
      if (!isConfigured.value) return {};

      const headers = {
        "X-AI-Base-Url": settings.value.baseUrl,
        "X-AI-Model": settings.value.model,
      };
      if (settings.value.apiKey) {
        headers["X-AI-Api-Key"] = settings.value.apiKey;
      }
      return headers;
    };

    return {
      settings,
      setSettings,
      clear,
      restoreApiKey,
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
