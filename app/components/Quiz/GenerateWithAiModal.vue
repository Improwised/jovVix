<script setup>
import { computed, onMounted, ref, watch } from "vue";
import {
  Check,
  Loader2,
  RefreshCw,
  Settings2,
  Sparkles,
  Trash2,
} from "lucide-vue-next";
import { usePush } from "notivue";
import { Modal } from "@/components/ui/modal";
import NavigationLink from "@/components/common/NavigationLink.vue";
import CodeBlockComponent from "@/components/CodeBlockComponent.vue";
import AiSettingsForm from "@/components/Quiz/AiSettingsForm.vue";
import { Skeleton } from "@/components/ui/skeleton";
import JvSelect from "@/components/ui/select/JvSelect.vue";
import { useAiSettingsStore } from "~/store/aiSettings";
import { readApiError } from "@/composables/apiError";

const MIN_QUESTIONS = 1;
const DEFAULT_MAX_QUESTIONS = 20;
const GENERATE_TIMEOUT_MS = 120000;

const DIFFICULTIES = [
  { value: "easy", label: "Easy" },
  { value: "medium", label: "Medium" },
  { value: "hard", label: "Hard" },
];

const LANGUAGES = [
  { value: "english", label: "English" },
  { value: "hindi", label: "Hindi" },
  { value: "gujarati", label: "Gujarati" },
  { value: "marathi", label: "Marathi" },
  { value: "bengali", label: "Bengali" },
  { value: "tamil", label: "Tamil" },
  { value: "telugu", label: "Telugu" },
  { value: "spanish", label: "Spanish" },
  { value: "french", label: "French" },
  { value: "german", label: "German" },
  { value: "portuguese", label: "Portuguese" },
  { value: "arabic", label: "Arabic" },
  { value: "chinese", label: "Chinese" },
  { value: "japanese", label: "Japanese" },
];

const props = defineProps({
  modelValue: { type: Boolean, required: true },
  mode: {
    type: String,
    default: "create",
    validator: (value) => ["create", "append"].includes(value),
  },
  quizId: { type: String, default: "" },
  quizTitle: { type: String, default: "" },
  quizLanguage: { type: String, default: "" },
  points: { type: Number, default: 0 },
  durationInSeconds: { type: Number, default: 0 },
});

const emit = defineEmits(["update:modelValue", "created", "appended"]);

const url = useRuntimeConfig().public;
const toast = usePush();
const aiSettings = useAiSettingsStore();

const form = ref({
  topic: "",
  number_of_questions: 10,
  difficulty: "medium",
  language: "english",
});
const generating = ref(false);
const generated = ref(null);
const generateError = ref("");
const saving = ref(false);
const view = ref("generate");
const aiStatus = ref(null);
const ready = ref(false);

onMounted(() => {
  ready.value = true;
});

const hasResults = computed(() => !!generated.value?.questions?.length);
const isAppend = computed(() => props.mode === "append");
const isSettingsView = computed(() => view.value === "settings");

const maxQuestions = computed(
  () => aiStatus.value?.max_questions || DEFAULT_MAX_QUESTIONS
);

const canGenerate = computed(() => aiSettings.isConfigured);

const saveLabel = computed(() =>
  isAppend.value
    ? `Add ${generated.value?.questions?.length ?? 0} questions`
    : "Add as a quiz"
);

const optionLabel = (index) => String.fromCharCode(65 + index);
const isSurvey = (question) => question.question_type === "survey";
const isCorrectOption = (question, index) =>
  !isSurvey(question) && question.correct_answer === index + 1;

const seedTopic = () => (isAppend.value ? props.quizTitle.trim() : "");

// New questions should match the ones the quiz already has, so append mode starts
// on the language those are written in rather than back at English.
const seedLanguage = () => (isAppend.value && props.quizLanguage) || "english";

// The quiz scoped route sends the questions this quiz already has to the model, so
// it does not generate them again.
const generateUrl = () =>
  isAppend.value
    ? `${url.apiUrl}/quizzes/${props.quizId}/questions/ai/generate`
    : `${url.apiUrl}/ai/questions/generate`;

const resetState = () => {
  form.value = {
    topic: seedTopic(),
    number_of_questions: 10,
    difficulty: "medium",
    language: seedLanguage(),
  };
  generated.value = null;
  generateError.value = "";
  generating.value = false;
  saving.value = false;
};

const fetchStatus = async () => {
  try {
    const response = await $fetch(`${url.apiUrl}/ai/status`, {
      headers: { Accept: "application/json" },
      credentials: "include",
    });
    aiStatus.value = response?.data ?? null;
  } catch {
    aiStatus.value = null;
  }
};

watch(
  () => props.modelValue,
  async (open) => {
    if (!open) return;
    resetState();
    await fetchStatus();
    view.value = canGenerate.value ? "generate" : "settings";
  },
  { immediate: true }
);

const saveAiSettings = (next) => {
  aiSettings.setSettings(next);
  view.value = "generate";
};

const discardResults = () => {
  generated.value = null;
  generateError.value = "";
};

const requestClose = (next) => {
  if (next) return;
  if (hasResults.value && !saving.value) {
    const confirmed = window.confirm("Discard these questions and close?");
    if (!confirmed) return;
  }
  emit("update:modelValue", false);
};

const handleGenerate = async () => {
  const topic = form.value.topic.trim();
  const count = Number(form.value.number_of_questions);

  if (topic.length < 3) {
    toast.error("Enter a topic of at least 3 characters.");
    return;
  }

  if (
    !Number.isInteger(count) ||
    count < MIN_QUESTIONS ||
    count > maxQuestions.value
  ) {
    toast.error(
      `Number of questions must be between ${MIN_QUESTIONS} and ${maxQuestions.value}.`
    );
    return;
  }

  try {
    generating.value = true;
    generateError.value = "";

    const response = await $fetch(generateUrl(), {
      method: "POST",
      headers: { Accept: "application/json", ...aiSettings.aiHeaders() },
      body: {
        topic,
        number_of_questions: count,
        difficulty: form.value.difficulty,
        language: form.value.language,
      },
      credentials: "include",
      timeout: GENERATE_TIMEOUT_MS,
    });

    const payload = response?.data;
    if (!payload?.questions?.length) {
      generateError.value =
        "The AI returned no questions. Try a more specific topic.";
      toast.error(generateError.value);
      return;
    }

    generated.value = payload;
    toast.success(`Generated ${payload.questions.length} questions.`);
  } catch (error) {
    generateError.value = readApiError(
      error,
      "Could not generate questions. Try again."
    );
    toast.error(generateError.value);
  } finally {
    generating.value = false;
  }
};

const questionPayload = () =>
  generated.value.questions.map((question) => ({
    question: question.question,
    question_type: question.question_type || "single",
    question_media: question.question_media || "text",
    resource: question.resource || "",
    options: question.options,
    options_media: question.options_media || "text",
    correct_answer: question.correct_answer || 0,
    explanation: question.explanation || "",
  }));

const createQuiz = async () => {
  const response = await $fetch(`${url.apiUrl}/ai/quizzes`, {
    method: "POST",
    headers: { Accept: "application/json" },
    body: {
      title: generated.value.suggested_title,
      description: generated.value.suggested_description || "",
      questions: questionPayload(),
    },
    credentials: "include",
  });

  const quizId = response?.data;
  if (!quizId) {
    toast.error("Error while creating quiz.");
    return;
  }

  toast.success("Quiz created successfully.");
  emit("created", quizId);
};

const appendQuestions = async () => {
  const response = await $fetch(
    `${url.apiUrl}/quizzes/${props.quizId}/questions/ai`,
    {
      method: "POST",
      headers: { Accept: "application/json" },
      body: {
        questions: questionPayload(),
        points: props.points,
        duration_in_seconds: props.durationInSeconds,
      },
      credentials: "include",
    }
  );

  const added = response?.data?.added ?? generated.value.questions.length;
  emit("appended", added);
};

const modalTitle = computed(() => {
  if (isSettingsView.value) return "Connect an AI provider";
  return isAppend.value
    ? "Generate questions with AI"
    : "Generate a quiz with AI";
});

const modalDescription = computed(() => {
  if (isSettingsView.value) return "";
  return isAppend.value
    ? "Enter a topic, check the questions, then add them to this quiz."
    : "Enter a topic, check the questions, then save them as a quiz.";
});

const handleSave = async () => {
  if (!hasResults.value) return;

  try {
    saving.value = true;
    if (isAppend.value) {
      await appendQuestions();
    } else {
      await createQuiz();
    }
  } catch (error) {
    toast.error(
      readApiError(
        error,
        isAppend.value
          ? "Error while adding questions."
          : "Error while creating quiz."
      )
    );
  } finally {
    saving.value = false;
  }
};
</script>

<template>
  <Modal
    :model-value="modelValue"
    size="xl"
    :close-on-backdrop="false"
    :title="modalTitle"
    :description="modalDescription"
    @update:model-value="requestClose"
  >
    <div class="max-h-[62vh] overflow-y-auto pb-3 pr-1">
      <Skeleton v-if="!ready" class="h-64 w-full" />

      <AiSettingsForm
        v-else-if="isSettingsView"
        :settings="aiSettings.settings"
        :can-cancel="canGenerate"
        @save="saveAiSettings"
        @cancel="view = 'generate'"
      />

      <template v-else>
        <form
          class="grid gap-5 sm:grid-cols-2 md:grid-cols-3"
          @submit.prevent="handleGenerate"
        >
          <label class="grid gap-2 sm:col-span-2 md:col-span-3">
            <span
              class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
            >
              Topic <span v-if="!isAppend" class="text-jv-coral">*</span>
            </span>
            <input
              v-model.trim="form.topic"
              type="text"
              required
              maxlength="200"
              placeholder="JavaScript closures"
              class="h-14 border-[3px] border-jv-ink bg-jv-canvas px-4 text-[17px] font-semibold text-jv-ink caret-jv-ink outline-none transition-shadow focus:shadow-brutal-sm"
            />
            <span class="font-body text-[13px] font-semibold text-jv-muted">
              <template v-if="isAppend">
                Taken from this quiz. Change it to cover a different area.
              </template>
              <template v-else>
                Be specific. "World War II treaties" works better than
                "history". Coding topics get code snippets automatically.
              </template>
            </span>
          </label>

          <label class="grid gap-2">
            <span
              class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
            >
              Number of questions <span class="text-jv-coral">*</span>
            </span>
            <input
              v-model.number="form.number_of_questions"
              type="number"
              required
              :min="MIN_QUESTIONS"
              :max="maxQuestions"
              step="1"
              class="h-14 border-[3px] border-jv-ink bg-jv-canvas px-4 text-[17px] font-semibold text-jv-ink caret-jv-ink outline-none transition-shadow focus:shadow-brutal-sm"
            />
          </label>

          <label class="grid gap-2">
            <span
              class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
            >
              Difficulty <span class="text-jv-coral">*</span>
            </span>
            <JvSelect v-model="form.difficulty" :options="DIFFICULTIES" />
          </label>

          <label class="grid gap-2">
            <span
              class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
            >
              Language <span class="text-jv-coral">*</span>
            </span>
            <JvSelect v-model="form.language" :options="LANGUAGES" />
          </label>

          <div
            class="flex flex-wrap items-center gap-3 sm:col-span-2 md:col-span-3"
          >
            <NavigationLink
              type="submit"
              url-name="Generate Questions"
              class="w-full bg-jv-coral py-2 font-[500] text-white sm:w-fit"
              :disabled="generating"
            >
              <Loader2
                v-if="generating"
                class="size-[18px] animate-spin"
                :stroke-width="2.4"
              />
              <Sparkles v-else class="size-[18px]" :stroke-width="2.4" />
            </NavigationLink>

            <button
              type="button"
              class="inline-flex items-center gap-1.5 font-body text-[13px] font-bold text-jv-muted underline"
              @click="view = 'settings'"
            >
              <Settings2 class="size-4" :stroke-width="2.4" />
              {{
                aiSettings.isConfigured ? "AI settings" : "Use my own AI key"
              }}
            </button>
          </div>
        </form>

        <section v-if="generating" class="mt-6 flex flex-col gap-4">
          <p class="font-body text-sm font-bold text-jv-muted">
            Asking the model. This can take up to a minute.
          </p>
          <Skeleton v-for="n in 3" :key="n" class="h-40 w-full" />
        </section>

        <section
          v-if="generateError && !generating"
          class="jv-border-rough mt-6 bg-jv-coral/15 p-4 sm:p-5"
        >
          <p class="text-[16px] font-bold text-jv-ink">{{ generateError }}</p>
          <div class="mt-4 flex flex-wrap gap-3">
            <NavigationLink
              url-name="Try again"
              class="bg-jv-white py-2 font-[500]"
              :disabled="generating"
              @click="handleGenerate"
            />
            <NavigationLink
              url-name="Fix AI settings"
              class="bg-jv-yellow py-2 font-[500]"
              @click="view = 'settings'"
            >
              <Settings2 class="size-[18px]" :stroke-width="2.4" />
            </NavigationLink>
          </div>
        </section>

        <section
          v-if="hasResults && !generating"
          class="mt-6 flex flex-col gap-4"
        >
          <p class="text-[16px] font-bold text-jv-ink">
            {{ generated.generated }} questions · {{ generated.topic }} ·
            {{ generated.difficulty }}
            <template
              v-if="generated.language && generated.language !== 'english'"
            >
              · {{ generated.language }}
            </template>
          </p>

          <p
            v-if="generated.notice"
            class="jv-border-rough bg-jv-yellow/25 p-3 text-[14px] font-semibold text-jv-ink"
          >
            {{ generated.notice }}
          </p>

          <p
            v-if="!isAppend && generated.suggested_title"
            class="font-body text-[14px] font-semibold text-jv-muted"
          >
            Saved as "{{ generated.suggested_title }}".
          </p>

          <article
            v-for="(question, index) in generated.questions"
            :key="index"
            class="jv-border-rough bg-jv-white p-4 shadow-brutal-sm sm:p-5"
          >
            <div class="flex flex-wrap items-center gap-2">
              <p
                class="text-[12px] font-bold uppercase tracking-[0.14em] text-jv-coral"
              >
                Question {{ index + 1 }}
              </p>
              <span
                v-if="isSurvey(question)"
                class="border-[2px] border-jv-ink bg-jv-yellow px-2 py-0.5 text-[11px] font-black uppercase tracking-[0.12em] text-jv-ink"
              >
                Survey: everyone scores
              </span>
            </div>
            <h3
              class="mt-1 break-words text-[18px] font-bold leading-snug text-jv-ink sm:text-[20px]"
            >
              {{ question.question }}
            </h3>

            <div
              v-if="question.question_media === 'code' && question.resource"
              class="mt-3 min-w-0 overflow-x-auto"
            >
              <CodeBlockComponent :code="question.resource" />
            </div>

            <ul class="mt-4 flex flex-col">
              <li
                v-for="(option, optionIndex) in question.options"
                :key="optionIndex"
                class="flex min-w-0 items-center gap-3 border-b border-jv-ink/10 py-3 pl-3 pr-2 text-[15px] font-medium text-jv-ink last:border-b-0"
                :class="
                  isCorrectOption(question, optionIndex)
                    ? 'border-l-4 border-l-jv-accent-green bg-jv-accent-green/25 pl-2'
                    : 'border-l-4 border-l-transparent'
                "
              >
                <span class="w-5 shrink-0 text-[14px] font-bold text-jv-coral">
                  {{ optionLabel(optionIndex) }}.
                </span>
                <div
                  v-if="question.options_media === 'code' && option"
                  class="min-w-0 flex-1 overflow-x-auto"
                >
                  <CodeBlockComponent :code="option" />
                </div>
                <span v-else class="min-w-0 flex-1 break-words">{{
                  option
                }}</span>
                <Check
                  v-if="isCorrectOption(question, optionIndex)"
                  class="size-5 shrink-0 text-jv-accent-green"
                  :stroke-width="3"
                />
              </li>
            </ul>

            <p
              v-if="question.explanation"
              class="mt-3 font-body text-[14px] italic text-jv-muted"
            >
              {{ question.explanation }}
            </p>
          </article>
        </section>
      </template>
    </div>

    <template v-if="hasResults && !isSettingsView" #footer>
      <NavigationLink
        url-name="Discard"
        class="bg-jv-white font-[500]"
        :disabled="saving"
        @click="discardResults"
      >
        <Trash2 class="size-[18px]" :stroke-width="2.4" />
      </NavigationLink>
      <NavigationLink
        url-name="Regenerate"
        class="bg-jv-yellow font-[500]"
        :disabled="generating || saving"
        @click="handleGenerate"
      >
        <RefreshCw class="size-[18px]" :stroke-width="2.4" />
      </NavigationLink>
      <NavigationLink
        :url-name="saveLabel"
        class="bg-jv-mint font-[500]"
        :disabled="saving"
        @click="handleSave"
      >
        <Loader2
          v-if="saving"
          class="size-[18px] animate-spin"
          :stroke-width="2.4"
        />
      </NavigationLink>
    </template>
  </Modal>
</template>
