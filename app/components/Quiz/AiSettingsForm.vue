<script setup>
import { computed, ref, watch } from "vue";
import { Check, ExternalLink, Loader2, X } from "lucide-vue-next";
import NavigationLink from "@/components/common/NavigationLink.vue";
import JvSelect from "@/components/ui/select/JvSelect.vue";
import { AI_PROVIDERS, findAiProvider } from "@/config/aiProviders";
import { readApiError } from "@/composables/apiError";

const props = defineProps({
  settings: { type: Object, required: true },
  canCancel: { type: Boolean, default: false },
});

const emit = defineEmits(["save", "cancel"]);

const url = useRuntimeConfig().public;

const withKnownProvider = (settings) => ({
  ...settings,
  provider:
    settings.provider && !findAiProvider(settings.provider)
      ? "custom"
      : settings.provider,
});

const form = ref(withKnownProvider(props.settings));
const testing = ref(false);
const testResult = ref(null);
const models = ref([]);
const modelsTruncated = ref(false);
const loadingModels = ref(false);
const modelsError = ref("");
const typedModel = ref(false);

watch(
  () => props.settings,
  (next) => {
    form.value = withKnownProvider(next);
    testResult.value = null;
    resetModels();
  }
);

const selectedProvider = computed(() => findAiProvider(form.value.provider));

const canSubmit = computed(() => !!form.value.baseUrl && !!form.value.model);

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

const loadModels = async () => {
  if (!canLoadModels.value || loadingModels.value) return;

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
};

const maybeLoadModels = () => {
  if (models.value.length || loadingModels.value) return;
  if (selectedProvider.value?.needsKey && !form.value.apiKey) return;
  loadModels();
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
  if (!canSubmit.value) return;

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
  emit("save", { ...form.value });
};

const handleSave = () => {
  if (!canSubmit.value) return;
  emit("save", { ...form.value });
};
</script>

<template>
  <div class="flex flex-col gap-5">
    <p
      class="jv-border-rough bg-jv-mint/25 p-3 font-body text-[14px] font-semibold text-jv-ink"
    >
      Pick a provider to generate quizzes. jovVix passes your key to your
      provider to generate questions. We never save it on our servers or write
      it to our logs, and it is cleared when you log out.
    </p>

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
      >
        API key
        <span v-if="selectedProvider?.needsKey" class="text-jv-coral">*</span>
      </span>
      <input
        v-model.trim="form.apiKey"
        type="password"
        autocomplete="off"
        placeholder="sk-…"
        class="h-14 border-[3px] border-jv-ink bg-jv-canvas px-4 text-[17px] font-semibold text-jv-ink caret-jv-ink outline-none transition-shadow focus:shadow-brutal-sm"
        @blur="maybeLoadModels"
      />
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
        :disabled="testing || !canSubmit"
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
        v-if="canCancel"
        url-name="Cancel"
        class="bg-jv-white font-[500]"
        @click="emit('cancel')"
      />
      <NavigationLink
        url-name="Save and continue"
        class="bg-jv-mint font-[500]"
        :disabled="!canSubmit"
        @click="handleSave"
      />
    </div>
  </div>
</template>
