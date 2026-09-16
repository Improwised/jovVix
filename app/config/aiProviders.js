// verifiedModels are the ones known to hold the generator's JSON contract. They
// are sorted to the top of the model list and marked, and everything else the
// provider offers still appears below them: this is a shortcut, not a whitelist.
// These ids go stale as providers retire models. The settings form only stars a
// model the provider actually lists, and replaces a preset it does not, so a
// stale entry here degrades the first pick rather than breaking generation.
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
    verifiedModels: [
      "deepseek/deepseek-chat",
      "google/gemini-2.5-flash",
      "openai/gpt-4o-mini",
      "meta-llama/llama-3.3-70b-instruct",
    ],
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
    verifiedModels: [
      "llama-3.3-70b-versatile",
      "llama-3.1-8b-instant",
      "openai/gpt-oss-120b",
      "qwen/qwen3-32b",
    ],
  },
  {
    id: "google",
    label: "Google AI Studio",
    baseUrl: "https://generativelanguage.googleapis.com/v1beta/openai",
    // Google lists its ids namespaced, and the id is sent exactly as listed.
    model: "models/gemini-3.5-flash",
    keyUrl: "https://aistudio.google.com/apikey",
    freeTier: true,
    needsKey: true,
    hint: "Free tier with daily limits.",
    verifiedModels: ["models/gemini-3.5-flash", "models/gemini-3.5-flash-lite"],
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
    verifiedModels: ["gpt-4o-mini", "gpt-4.1-mini", "gpt-4o", "gpt-4.1"],
  },
  {
    id: "custom",
    label: "Custom / LiteLLM",
    freeTier: false,
    needsKey: false,
    hint: "Any endpoint that speaks the OpenAI chat-completions format.",
    verifiedModels: [],
  },
];

export const findAiProvider = (id) =>
  AI_PROVIDERS.find((provider) => provider.id === id) || null;
