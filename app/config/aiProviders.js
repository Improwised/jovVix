export const AI_PROVIDERS = [
  {
    id: "openrouter",
    label: "OpenRouter",
    baseUrl: "https://openrouter.ai/api/v1",
    model: "deepseek/deepseek-chat",
    keyUrl: "https://openrouter.ai/keys",
    freeTier: true,
    needsKey: true,
    hint: "Models with :free in the name cost nothing.",
  },
  {
    id: "groq",
    label: "Groq",
    baseUrl: "https://api.groq.com/openai/v1",
    model: "llama-3.3-70b-versatile",
    keyUrl: "https://console.groq.com/keys",
    freeTier: true,
    needsKey: true,
    hint: "Free tier, no card needed.",
  },
  {
    id: "google",
    label: "Google AI Studio",
    baseUrl: "https://generativelanguage.googleapis.com/v1beta/openai",
    model: "gemini-2.5-flash",
    keyUrl: "https://aistudio.google.com/apikey",
    freeTier: true,
    needsKey: true,
    hint: "Free tier with daily limits.",
  },
  {
    id: "openai",
    label: "OpenAI",
    baseUrl: "https://api.openai.com/v1",
    model: "gpt-4o-mini",
    keyUrl: "https://platform.openai.com/api-keys",
    freeTier: false,
    needsKey: true,
    hint: "Paid. You are billed for every quiz you generate.",
  },
  {
    id: "custom",
    label: "Custom / LiteLLM",
    freeTier: false,
    needsKey: false,
    hint: "Any endpoint that speaks the OpenAI chat-completions format.",
  },
];

export const findAiProvider = (id) =>
  AI_PROVIDERS.find((provider) => provider.id === id) || null;
