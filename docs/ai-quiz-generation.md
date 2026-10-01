# AI Quiz Generation Guide

jovVix lets you generate entire quizzes - or add more questions to an existing quiz - using AI. You bring your own API key from any OpenAI-compatible provider, and the server takes care of the rest. Your key is encrypted with a Vault Password that only you know; the server never stores or sees it in plaintext.

## Supported Providers

jovVix ships with presets for these providers:

| Provider | Base URL | Free Tier | Where to Get a Key |
|----------|----------|-----------|--------------------|
| OpenRouter | `https://openrouter.ai/api/v1` | Yes | [openrouter.ai/keys](https://openrouter.ai/keys) |
| Groq | `https://api.groq.com/openai/v1` | Yes | [console.groq.com/keys](https://console.groq.com/keys) |
| Google AI Studio | `https://generativelanguage.googleapis.com/v1beta/openai` | Yes | [aistudio.google.com/apikey](https://aistudio.google.com/apikey) |
| OpenAI | `https://api.openai.com/v1` | No | [platform.openai.com/api-keys](https://platform.openai.com/api-keys) |
| Custom / LiteLLM | _you specify_ | - | Any endpoint that speaks the OpenAI chat-completions format |

Each preset also lists **verified models** known to reliably follow the generation contract. These are starred in the model dropdown as a shortcut, but you can use any chat-capable model your provider offers.

## How the Vault Works

API keys are sensitive. jovVix encrypts yours before storing it:

1. You choose a **Vault Password** (8-512 characters) when saving your AI settings.
2. Keys are encrypted with **AES-256-GCM**, with the encryption key derived from your password using **Argon2id** (a hardened key-derivation function).
3. The Vault Password itself is **never stored** - not in the database, not in logs, not on disk.

When you need to generate or append questions, the UI prompts you to unlock the vault with your password. The server decrypts the key in memory just long enough to make the AI call, then discards it.

You can change providers or reset your API key without changing the Vault Password, or delete the vault entirely to remove your saved credentials.

## Configuring AI Settings

1. Navigate to the **Quiz** section from the left sidebar, then click **AI Settings**.
2. Choose a provider from the dropdown. The base URL fills in automatically; change it only if you are using a custom endpoint.
3. Enter your **API key** from the provider.
4. Pick a **model**. You can either fetch the provider's model list or type a model name directly.
5. Click **Test Connection** to verify the provider responds with a valid quiz question.
6. Set a **Vault Password** (and confirm it), then click **Save Settings**.

Your saved settings persist across sessions. The sidebar shows a masked version of your key so you can confirm which provider is configured.

## Generating a Quiz

You can generate questions in two modes:

### Create a New Quiz

Click **Generate with AI** from the quiz list page, then fill in:

- **Topic** - the subject of the quiz (e.g. "world capitals", "Python async patterns")
- **Language** - 14 languages supported (see below)
- **Difficulty** - easy, medium, or hard
- **Number of questions** - 1 to 20 (configurable via `AI_MAX_QUESTIONS`, capped at 50)

Click **Generate**. The server calls your provider and returns questions matching the JSON contract. Once you are happy with the generated set, you can create the quiz with a single click.

### Append to an Existing Quiz

From a quiz's detail page, click **Generate with AI** in append mode. The server sends the text of your existing questions so the model avoids asking the same facts again. Otherwise the flow is identical: choose topic, language, difficulty, and count, then generate and append.

## Supported Languages and Difficulties

### Languages

| Code | Language |
|------|----------|
| `english` | English |
| `hindi` | Hindi |
| `gujarati` | Gujarati |
| `marathi` | Marathi |
| `bengali` | Bengali |
| `tamil` | Tamil |
| `telugu` | Telugu |
| `spanish` | Spanish |
| `french` | French |
| `german` | German |
| `portuguese` | Portuguese |
| `arabic` | Arabic |
| `chinese` | Simplified Chinese |
| `japanese` | Japanese |

All questions, options, and the quiz title and description are written in the chosen language. Only the JSON keys and the values of `question_type`, `question_media`, and `options_media` remain in English.

### Difficulties

| Level | Guidance |
|-------|----------|
| Easy | Common knowledge that a beginner would recognise; single-step recall |
| Medium | Requires real familiarity with the topic and one step of reasoning |
| Hard | Requires specialist knowledge, precise detail, or multi-step reasoning; distractors should be near-misses |

## Question Types and Media

The AI generates two question types:

- **Single** - exactly one option is factually correct. This is the default.
- **Survey** - a genuine opinion question where every option counts as correct. Used sparingly (at most one in every five questions).

Questions can include code snippets by setting `question_media` or `options_media` to `code`. Snippets are stored as raw source (no markdown fences), and the generator mixes text and code questions naturally rather than making every question code-based.

## Environment Variables

| Variable | File | Default | Description |
|----------|------|---------|-------------|
| `AI_TEMPERATURE` | `api/.env` | `0.4` | Controls generation creativity (0.0-2.0). Lower values produce more deterministic output. |
| `AI_TIMEOUT_SECONDS` | `api/.env` | `90` | AI request timeout in seconds. Must stay below the frontend's 120-second fetch timeout. |
| `AI_MAX_QUESTIONS` | `api/.env` | `20` | Maximum questions per generation request. Capped at 50. |
| `AI_JSON_MODE` | `api/.env` | `true` | Enforce `json_object` response format on the provider. Disable if your provider does not support it. |

## API Endpoints

All AI endpoints require Kratos authentication and live under `/api/v1/ai`:

| Method | Endpoint | Purpose |
|--------|----------|---------|
| `GET` | `/ai/status` | Check whether AI is configured for the current user |
| `GET` | `/ai/settings` | Retrieve saved provider and model (key is masked) |
| `PUT` | `/ai/settings` | Save or update provider, model, API key, and vault |
| `POST` | `/ai/settings/unlock` | Unlock the vault for a session |
| `DELETE` | `/ai/settings` | Delete the vault and all saved credentials |
| `GET` | `/ai/models` | List chat-capable models from the configured provider |
| `POST` | `/ai/test` | Test the connection: returns a sample question and response time |
| `POST` | `/ai/questions/generate` | Generate questions for a new quiz |
| `POST` | `/ai/quizzes` | Create a quiz from generated questions |
| `POST` | `/quizzes/:quiz_id/questions/ai/generate` | Generate questions for an existing quiz |
| `POST` | `/quizzes/:quiz_id/questions/ai` | Append generated questions to an existing quiz |

## How Provider Adaptations Work

Not every OpenAI-compatible provider supports the same set of features. The server sends a best-effort request with JSON mode, `max_tokens`, and `temperature`. If the provider rejects a parameter with a 400 or 422 status, the server retries without that parameter:

1. **JSON mode** - dropped if the provider rejects `response_format`
2. **max_tokens** - switched to `max_completion_tokens` if the provider rejects `max_tokens`
3. **Temperature** - dropped if the provider rejects the `temperature` field

Up to 3 adaptations are attempted before giving up.

## Troubleshooting

**"AI quiz generation is not configured"** - you have not saved your AI settings yet. Go to AI Settings from the quiz sidebar and configure a provider.

**"Could not connect to that base URL"** - the jovVix server makes this call, not your browser. If the server cannot reach the provider (network restrictions, DNS, TLS mismatch), try another provider or use a custom endpoint accessible from your server.

**"That host name does not resolve"** or **"TLS handshake failed"** - double-check the base URL. Many providers need `/v1` at the end. Verify http vs https.

**"The provider rejected that API key"** - the key may be invalid, expired, or has insufficient permissions for the chosen model.

**"Endpoint or model not found"** - the model name is likely wrong. Fetch the model list from the provider or type the exact model ID from the provider's documentation.

**"Connected, but this model did not return a usable quiz question"** - the model is reachable but cannot produce valid JSON in the required shape. Pick a different model from the verified list for your provider.

**"The AI service returned an empty response"** - the provider accepted the request but returned no content. This can happen with very restrictive free tiers or when the topic is outside the model's knowledge. Try a different topic or model.

**Timeouts** - if generation consistently times out, increase `AI_TIMEOUT_SECONDS` (max 120s to stay under the frontend fetch timeout). For free-tier providers, smaller question counts (5-10) reduce the chance of hitting rate limits.

**Vault Password forgotten** - there is no recovery mechanism. Delete the vault and configure your provider again with a new password.