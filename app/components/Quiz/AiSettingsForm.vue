<script setup>
import { computed, ref, watch } from "vue";
import { Check, ExternalLink, Eye, EyeOff, Loader2, X } from "lucide-vue-next";
import NavigationLink from "@/components/common/NavigationLink.vue";
import JvSelect from "@/components/ui/select/JvSelect.vue";
import { AI_PROVIDERS, findAiProvider } from "@/config/aiProviders";
import { readApiError } from "@/composables/apiError";

const props = defineProps({
  settings: { type: Object, required: true },
  unlocked: { type: Boolean, default: false },
});

const emit = defineEmits(["save", "unlock", "replace", "cancel"]);

const url = useRuntimeConfig().public;

const withKnownProvider = (settings) => ({
  ...settings,
  provider:
    settings.provider && !findAiProvider(settings.provider)
      ? "custom"
      : settings.provider,
});

const form = ref({
  ...withKnownProvider(props.settings),
  vaultPassword: "",
  vaultPasswordConfirmation: "",
});
const testing = ref(false);
const testResult = ref(null);
const models = ref([]);
const modelsTruncated = ref(false);
const loadingModels = ref(false);
const modelsError = ref("");
const typedModel = ref(false);
const submitAttempted = ref(false);
const showVaultPassword = ref(false);
const showVaultPasswordConfirmation = ref(false);
const showApiKey = ref(false);
let modelsLoadPromise = null;

watch(
  () => props.settings,
  (next) => {
    form.value = {
      ...withKnownProvider(next),
      vaultPassword: "",
      vaultPasswordConfirmation: "",
    };
    testResult.value = null;
    submitAttempted.value = false;
    showVaultPassword.value = false;
    showVaultPasswordConfirmation.value = false;
    showApiKey.value = false;
    resetModels();
  }
);

const selectedProvider = computed(() => findAiProvider(form.value.provider));
const isUnlocking = computed(() => !!form.value.configured && !props.unlocked);
const isSavedAndUnlocked = computed(
  () => !!form.value.configured && props.unlocked
);
const needsApiKey = computed(() => selectedProvider.value?.needsKey !== false);
const hasProviderConnectionDetails = computed(
  () =>
    !!form.value.baseUrl &&
    !!form.value.model &&
    (!needsApiKey.value || !!form.value.apiKey)
);

const canSubmit = computed(() =>
  isUnlocking.value
    ? form.value.vaultPassword.length >= 8
    : hasProviderConnectionDetails.value &&
      form.value.vaultPassword.length >= 8 &&
      form.value.vaultPassword === form.value.vaultPasswordConfirmation
);

const canAttemptSubmit = computed(() =>
  isUnlocking.value ? true : hasProviderConnectionDetails.value
);

const vaultPasswordError = computed(() => {
  if (!submitAttempted.value || form.value.vaultPassword.length >= 8) {
    return "";
  }
  return "Vault Password must be at least 8 characters.";
});

const vaultPasswordConfirmationError = computed(() => {
  if (!submitAttempted.value) return "";
  if (!form.value.vaultPasswordConfirmation) {
    return "Confirm your Vault Password.";
  }
  if (form.value.vaultPassword !== form.value.vaultPasswordConfirmation) {
    return "Vault Passwords do not match.";
  }
  return "";
});

const canLoadModels = computed(() => !!form.value.baseUrl);

const TYPED_MODEL = "__typed__";

const providerOptions = computed(() =>
  AI_PROVIDERS.map((provider) => ({
    value: provider.id,
    label: `${provider.label}${provider.freeTier ? " (free tier)" : ""}`,
  }))
);

const verifiedModels = computed(
  () => selectedProvider.value?.verifiedModels || []
);

// Providers disagree on whether an id carries a namespace: Google lists
// "models/gemini-2.5-flash" where the preset says "gemini-2.5-flash". Compare on
// the last segment. Only ever for matching, never for the value we send.
const modelKey = (model) => model.split("/").pop().toLowerCase();

const isVerifiedModel = (model) =>
  verifiedModels.value.some((known) => modelKey(known) === modelKey(model));

const modelOptions = computed(() => {
  // A model the provider did not list is still shown so a restored or typed one
  // stays selectable, but it never earns a star: the star claims it works.
  const offered = models.value;
  const unlisted =
    form.value.model && !offered.includes(form.value.model)
      ? [form.value.model]
      : [];

  const verified = offered.filter(isVerifiedModel);
  const rest = offered.filter((model) => !isVerifiedModel(model));

  return [
    ...verified.map((model) => ({ value: model, label: `★ ${model}` })),
    ...rest.map((model) => ({ value: model, label: model })),
    ...unlisted.map((model) => ({ value: model, label: model })),
    { value: TYPED_MODEL, label: "Type a model name" },
  ];
});

const modelsNotice = computed(() => {
  if (!models.value.length) return "";

  const parts = [];
  if (verifiedModels.value.length) parts.push("★ = known to work well");
  if (modelsTruncated.value) parts.push("list shortened");

  return parts.join(" · ");
});

const showModelList = computed(
  () => models.value.length > 0 && !typedModel.value
);

const resetModels = () => {
  models.value = [];
  modelsTruncated.value = false;
  modelsError.value = "";
  typedModel.value = false;
};

const selectModel = (value) => {
  typedModel.value = value === TYPED_MODEL;
  if (typedModel.value) return;
  form.value.model = value;
  testResult.value = null;
};

// Every field belongs to the provider that was selected, so switching wipes the
// lot and refills from the new preset. Custom has no preset and starts empty.
const applyProvider = (id) => {
  const provider = findAiProvider(id);

  form.value.provider = id;
  form.value.baseUrl = provider?.baseUrl || "";
  form.value.model = provider?.model || "";
  form.value.apiKey = "";

  testResult.value = null;
  resetModels();

  if (provider && id !== "custom") maybeLoadModels();
};

const loadModels = () => {
  if (!canLoadModels.value) return Promise.resolve();
  if (modelsLoadPromise) return modelsLoadPromise;

  modelsLoadPromise = (async () => {
    try {
      loadingModels.value = true;
      modelsError.value = "";

      const headers = {
        Accept: "application/json",
        "X-AI-Base-Url": form.value.baseUrl,
      };
      if (form.value.apiKey) headers["X-AI-Api-Key"] = form.value.apiKey;

      const response = await $fetch(`${url.apiUrl}/ai/models`, {
        headers,
        credentials: "include",
      });

      models.value = response?.data?.models ?? [];
      modelsTruncated.value = response?.data?.truncated ?? false;
      typedModel.value = false;
      if (!models.value.length) {
        modelsError.value = "That provider returned no usable models.";
        return;
      }

      // A preset goes stale when the provider retires a model or names it
      // differently, so anything the provider did not list is replaced with one it
      // did rather than left to fail at generate time.
      if (!form.value.model || !models.value.includes(form.value.model)) {
        const preferred = models.value.find(isVerifiedModel);
        selectModel(preferred || models.value[0]);
      }
    } catch (error) {
      models.value = [];
      modelsTruncated.value = false;
      modelsError.value = readApiError(error, "Could not load models.");
    } finally {
      loadingModels.value = false;
    }
  })();

  return modelsLoadPromise.finally(() => {
    modelsLoadPromise = null;
  });
};

const maybeLoadModels = async () => {
  if (models.value.length) return;
  if (selectedProvider.value?.needsKey && !form.value.apiKey) return;
  await loadModels();
};

const handleBaseUrlInput = () => {
  markCustom();
  resetModels();
};

// A provider is its endpoint, so only the base URL can make the choice custom.
// Switching to another of that provider's models keeps the hints and the ★ marks.
const markCustom = () => {
  const provider = selectedProvider.value;
  if (!provider) return;
  if (form.value.baseUrl !== provider.baseUrl) {
    form.value.provider = "custom";
  }
};

const handleTest = async () => {
  if (!hasProviderConnectionDetails.value) return;

  // Clicking Test immediately after entering a key also blurs the key field.
  // Wait for that model fetch so this request uses the provider's current model.
  await maybeLoadModels();
  if (!hasProviderConnectionDetails.value) return;

  try {
    testing.value = true;
    testResult.value = null;

    const headers = {
      Accept: "application/json",
      "X-AI-Base-Url": form.value.baseUrl,
      "X-AI-Model": form.value.model,
    };
    if (form.value.apiKey) headers["X-AI-Api-Key"] = form.value.apiKey;

    const response = await $fetch(`${url.apiUrl}/ai/test`, {
      method: "POST",
      headers,
      credentials: "include",
    });

    const latency = response?.data?.latency_ms ?? 0;
    const sample = response?.data?.sample;

    testResult.value = {
      ok: true,
      message: sample
        ? `Wrote a question in ${latency} ms: "${sample}"`
        : `Connected in ${latency} ms.`,
    };
  } catch (error) {
    testResult.value = {
      ok: false,
      message: readApiError(error, "Could not reach that provider."),
    };
  } finally {
    testing.value = false;
  }
};

const forgetApiKey = () => {
  form.value.apiKey = "";
  testResult.value = null;
};

const handleSave = () => {
  submitAttempted.value = true;
  if (!canSubmit.value) return;
  if (isUnlocking.value) {
    emit("unlock", form.value.vaultPassword);
    return;
  }
  emit("save", { ...form.value });
};

watch(
  () => [
    form.value.provider,
    form.value.baseUrl,
    form.value.model,
    form.value.apiKey,
  ],
  () => {
    testResult.value = null;
  }
);
</script>

<template>
  <div class="flex flex-col gap-5">
    <p
      class="jv-border-rough bg-jv-mint/25 p-3 font-body text-[14px] font-semibold text-jv-ink"
    >
      jovVix stores your API key only as an encrypted vault value. Your Vault
      Password is never stored and stays only in this browser session.
    </p>

    <template v-if="isSavedAndUnlocked">
      <p class="font-body text-[14px] font-semibold text-jv-muted">
        Your saved {{ form.provider || "AI" }} settings are unlocked for this
        browser session.
      </p>
      <div class="grid gap-1 font-body text-[14px] text-jv-muted">
        <span>Provider: {{ form.provider }}</span
        ><span>Base URL: {{ form.baseUrl }}</span
        ><span>Model: {{ form.model }}</span>
      </div>
      <div class="flex justify-end gap-3">
        <NavigationLink
          url-name="Back to generation"
          class="bg-jv-white font-[500]"
          @click="emit('cancel')"
        />
        <NavigationLink
          url-name="Clear AI settings"
          class="bg-jv-white font-[500]"
          @click="emit('replace')"
        />
      </div>
    </template>

    <form
      v-else-if="isUnlocking"
      class="flex flex-col gap-5"
      @submit.prevent="handleSave"
    >
      <p class="font-body text-[14px] font-semibold text-jv-muted">
        Your saved {{ form.provider || "AI" }} settings are ready. Enter only
        your Vault Password to unlock them for this session.
      </p>
      <label class="grid gap-2">
        <span
          class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
          >Vault Password</span
        >
        <div class="relative">
          <input
            v-model="form.vaultPassword"
            :type="showVaultPassword ? 'text' : 'password'"
            autocomplete="current-password"
            placeholder="Enter your Vault Password"
            :aria-invalid="!!vaultPasswordError"
            class="h-14 w-full border-[3px] bg-jv-canvas py-0 pl-4 pr-14 text-[17px] font-semibold text-jv-ink outline-none transition-shadow focus:shadow-brutal-sm"
            :class="vaultPasswordError ? 'border-jv-coral' : 'border-jv-ink'"
          />
          <button
            type="button"
            class="absolute inset-y-0 right-0 grid w-14 place-items-center text-jv-muted transition-colors hover:text-jv-ink focus:outline-none focus-visible:text-jv-ink"
            :aria-label="
              showVaultPassword ? 'Hide Vault Password' : 'Show Vault Password'
            "
            :aria-pressed="showVaultPassword"
            @click="showVaultPassword = !showVaultPassword"
          >
            <EyeOff
              v-if="showVaultPassword"
              class="size-5"
              :stroke-width="2.4"
            />
            <Eye v-else class="size-5" :stroke-width="2.4" />
          </button>
        </div>
        <span
          v-if="vaultPasswordError"
          class="font-body text-[13px] font-semibold text-jv-coral"
        >
          {{ vaultPasswordError }}
        </span>
      </label>
      <div class="flex justify-end gap-3">
        <NavigationLink
          url-name="Clear AI settings"
          class="bg-jv-white font-[500]"
          type="button"
          @click="emit('replace')"
        />
        <NavigationLink
          url-name="Unlock"
          class="bg-jv-mint font-[500]"
          :disabled="!canAttemptSubmit"
          type="submit"
        />
      </div>
    </form>

    <form
      v-else
      class="flex min-w-0 flex-col gap-5"
      @submit.prevent="handleSave"
    >
      <label class="grid gap-2">
        <span
          class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
        >
          Provider
        </span>
        <JvSelect
          :model-value="form.provider"
          :options="providerOptions"
          placeholder="Choose a provider"
          aria-label="Provider"
          @update:model-value="applyProvider"
        />
        <span
          v-if="selectedProvider?.hint"
          class="font-body text-[13px] font-semibold text-jv-muted"
        >
          {{ selectedProvider.hint }}
        </span>
      </label>

      <label class="grid gap-2">
        <span
          class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
        >
          Base URL <span class="text-jv-coral">*</span>
        </span>
        <input
          v-model.trim="form.baseUrl"
          type="text"
          required
          placeholder="https://openrouter.ai/api/v1"
          class="h-14 border-[3px] border-jv-ink bg-jv-canvas px-4 text-[17px] font-semibold text-jv-ink caret-jv-ink outline-none transition-shadow focus:shadow-brutal-sm"
          @input="handleBaseUrlInput"
        />
      </label>

      <label class="grid gap-2">
        <span
          class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
          >Vault Password <span class="text-jv-coral">*</span></span
        >
        <div class="relative">
          <input
            v-model="form.vaultPassword"
            :type="showVaultPassword ? 'text' : 'password'"
            autocomplete="new-password"
            placeholder="At least 8 characters"
            :aria-invalid="!!vaultPasswordError"
            class="h-14 w-full border-[3px] bg-jv-canvas py-0 pl-4 pr-14 text-[17px] font-semibold text-jv-ink outline-none transition-shadow focus:shadow-brutal-sm"
            :class="vaultPasswordError ? 'border-jv-coral' : 'border-jv-ink'"
          />
          <button
            type="button"
            class="absolute inset-y-0 right-0 grid w-14 place-items-center text-jv-muted transition-colors hover:text-jv-ink focus:outline-none focus-visible:text-jv-ink"
            :aria-label="
              showVaultPassword ? 'Hide Vault Password' : 'Show Vault Password'
            "
            :aria-pressed="showVaultPassword"
            @click="showVaultPassword = !showVaultPassword"
          >
            <EyeOff
              v-if="showVaultPassword"
              class="size-5"
              :stroke-width="2.4"
            />
            <Eye v-else class="size-5" :stroke-width="2.4" />
          </button>
        </div>
        <span
          v-if="vaultPasswordError"
          class="font-body text-[13px] font-semibold text-jv-coral"
        >
          {{ vaultPasswordError }}
        </span>
      </label>

      <label class="grid gap-2">
        <span
          class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
          >Confirm Vault Password <span class="text-jv-coral">*</span></span
        >
        <div class="relative">
          <input
            v-model="form.vaultPasswordConfirmation"
            :type="showVaultPasswordConfirmation ? 'text' : 'password'"
            autocomplete="new-password"
            placeholder="Enter it again to confirm"
            :aria-invalid="!!vaultPasswordConfirmationError"
            class="h-14 w-full border-[3px] bg-jv-canvas py-0 pl-4 pr-14 text-[17px] font-semibold text-jv-ink outline-none transition-shadow focus:shadow-brutal-sm"
            :class="
              vaultPasswordConfirmationError
                ? 'border-jv-coral'
                : 'border-jv-ink'
            "
          />
          <button
            type="button"
            class="absolute inset-y-0 right-0 grid w-14 place-items-center text-jv-muted transition-colors hover:text-jv-ink focus:outline-none focus-visible:text-jv-ink"
            :aria-label="
              showVaultPasswordConfirmation
                ? 'Hide confirmation Vault Password'
                : 'Show confirmation Vault Password'
            "
            :aria-pressed="showVaultPasswordConfirmation"
            @click="
              showVaultPasswordConfirmation = !showVaultPasswordConfirmation
            "
          >
            <EyeOff
              v-if="showVaultPasswordConfirmation"
              class="size-5"
              :stroke-width="2.4"
            />
            <Eye v-else class="size-5" :stroke-width="2.4" />
          </button>
        </div>
        <span
          v-if="vaultPasswordConfirmationError"
          class="font-body text-[13px] font-semibold text-jv-coral"
        >
          {{ vaultPasswordConfirmationError }}
        </span>
      </label>

      <label class="grid gap-2">
        <span
          class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
        >
          API key
          <span v-if="selectedProvider?.needsKey" class="text-jv-coral">*</span>
        </span>
        <div class="relative">
          <input
            v-model.trim="form.apiKey"
            :type="showApiKey ? 'text' : 'password'"
            autocomplete="off"
            placeholder="sk-…"
            class="h-14 w-full border-[3px] border-jv-ink bg-jv-canvas py-0 pl-4 pr-14 text-[17px] font-semibold text-jv-ink caret-jv-ink outline-none transition-shadow focus:shadow-brutal-sm"
            @blur="maybeLoadModels"
          />
          <button
            type="button"
            class="absolute inset-y-0 right-0 grid w-14 place-items-center text-jv-muted transition-colors hover:text-jv-ink focus:outline-none focus-visible:text-jv-ink"
            :aria-label="showApiKey ? 'Hide API key' : 'Show API key'"
            :aria-pressed="showApiKey"
            @click="showApiKey = !showApiKey"
          >
            <EyeOff v-if="showApiKey" class="size-5" :stroke-width="2.4" />
            <Eye v-else class="size-5" :stroke-width="2.4" />
          </button>
        </div>
        <span
          class="flex flex-wrap items-center gap-2 font-body text-[13px] font-semibold text-jv-muted"
        >
          <a
            v-if="selectedProvider?.keyUrl"
            :href="selectedProvider.keyUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="inline-flex items-center gap-1 underline"
          >
            Get a key
            <ExternalLink class="size-3" :stroke-width="2.6" />
          </a>
          <span v-if="!selectedProvider?.needsKey">
            Not needed for a local model.
          </span>
          <button
            v-if="form.apiKey"
            type="button"
            class="inline-flex items-center gap-1 underline"
            @click="forgetApiKey"
          >
            Forget my key
          </button>
        </span>
      </label>

      <label class="grid gap-2">
        <span
          class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
        >
          Model <span class="text-jv-coral">*</span>
        </span>
        <JvSelect
          v-if="showModelList"
          :model-value="form.model"
          :options="modelOptions"
          aria-label="Model"
          required
          @update:model-value="selectModel"
        />
        <input
          v-else
          v-model.trim="form.model"
          type="text"
          required
          placeholder="deepseek/deepseek-chat"
          class="h-14 border-[3px] border-jv-ink bg-jv-canvas px-4 text-[17px] font-semibold text-jv-ink caret-jv-ink outline-none transition-shadow focus:shadow-brutal-sm"
          @focus="maybeLoadModels"
        />
      </label>

      <div class="-mt-2 flex flex-wrap items-center gap-3">
        <button
          type="button"
          class="inline-flex items-center gap-1.5 font-body text-[13px] font-bold text-jv-muted underline disabled:no-underline disabled:opacity-60"
          :disabled="loadingModels || !canLoadModels"
          @click="loadModels"
        >
          <Loader2
            v-if="loadingModels"
            class="size-4 animate-spin"
            :stroke-width="2.4"
          />
          {{ models.length ? "Refresh models" : "Load models" }}
        </button>
        <span
          v-if="modelsError"
          class="font-body text-[13px] font-semibold text-jv-coral"
        >
          {{ modelsError }}
        </span>
        <span
          v-else-if="modelsNotice"
          class="font-body text-[13px] font-semibold text-jv-muted"
        >
          {{ modelsNotice }}
        </span>
      </div>

      <div class="flex flex-wrap items-center gap-3">
        <NavigationLink
          url-name="Test connection"
          class="bg-jv-white py-2 font-[500]"
          :disabled="testing || !hasProviderConnectionDetails"
          type="button"
          @click="handleTest"
        >
          <Loader2
            v-if="testing"
            class="size-[18px] animate-spin"
            :stroke-width="2.4"
          />
        </NavigationLink>

        <span
          v-if="testing"
          class="font-body text-[14px] font-semibold text-jv-muted"
        >
          Asking this model for one question…
        </span>
        <p
          v-else-if="testResult"
          class="flex min-w-0 items-start gap-2 font-body text-[14px] font-bold"
          :class="testResult.ok ? 'text-jv-accent-green' : 'text-jv-coral'"
        >
          <Check
            v-if="testResult.ok"
            class="mt-0.5 size-4 shrink-0"
            :stroke-width="3"
          />
          <X v-else class="mt-0.5 size-4 shrink-0" :stroke-width="3" />
          <span class="min-w-0 break-words">{{ testResult.message }}</span>
        </p>
      </div>

      <div class="flex flex-wrap justify-end gap-3">
        <NavigationLink
          url-name="Cancel"
          class="bg-jv-white font-[500]"
          type="button"
          @click="emit('cancel')"
        />
        <NavigationLink
          url-name="Save and continue"
          class="bg-jv-mint font-[500]"
          :disabled="!canAttemptSubmit"
          type="submit"
        />
      </div>
    </form>
  </div>
</template>
