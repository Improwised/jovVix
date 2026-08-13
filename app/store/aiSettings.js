import { defineStore } from "pinia";

const emptySettings = () => ({
  provider: "",
  baseUrl: "",
  apiKey: "",
  model: "",
});

export const useAiSettingsStore = defineStore(
  "ai-settings-store",
  () => {
    const settings = ref(emptySettings());

    const setSettings = (next) => {
      settings.value = { ...emptySettings(), ...next };
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

    return { settings, setSettings, isConfigured, aiHeaders };
  },
  {
    persist: true,
  }
);
