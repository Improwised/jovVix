<template>
  <form
    class="jv-border-rough bg-jv-white p-4 shadow-brutal-sm sm:p-5 lg:p-6"
    :class="mode === 'edit' ? 'rotate-[0.15deg]' : 'rotate-[-0.2deg]'"
    @submit.prevent="submitForm"
  >
    <div
      class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between"
    >
      <div>
        <p
          class="text-[14px] font-black uppercase tracking-[0.12em] text-jv-coral"
        >
          {{ eyebrow }}
        </p>
        <h2
          class="mt-1 font-body text-[24px] font-black leading-tight text-jv-ink sm:text-[28px]"
        >
          {{ title }}
        </h2>
      </div>

      <select
        v-model.number="form.type"
        class="h-11 border-[3px] border-jv-ink bg-jv-white px-3 text-[15px] font-black text-jv-ink shadow-brutal-sm outline-none"
      >
        <option :value="1">Multiple Choice</option>
        <option :value="2">Survey</option>
      </select>
    </div>

    <div class="mt-5 grid gap-4">
      <input
        v-model.trim="form.question"
        type="text"
        required
        placeholder="Enter question..."
        class="h-12 w-full border-[3px] border-jv-ink bg-jv-white px-3 text-[15px] font-semibold text-jv-ink outline-none transition-shadow placeholder:text-jv-ink/40 focus:shadow-brutal-sm sm:text-[16px]"
      />

      <div class="grid gap-3 border-y-2 border-dashed border-jv-ink/15 py-4">
        <div class="flex flex-wrap items-center justify-center gap-2">
          <button
            v-for="choice in mediaChoices"
            :key="`question-${choice.value}`"
            type="button"
            class="inline-flex h-9 items-center gap-1.5 border-[2px] border-jv-ink px-3 text-[13px] font-black shadow-[1px_1px_0_#2D2D2D] transition-transform hover:rotate-[1deg]"
            :class="
              form.question_media === choice.value
                ? 'bg-jv-ink text-jv-white'
                : 'bg-jv-white text-jv-ink'
            "
            @click="form.question_media = choice.value"
          >
            <component :is="choice.icon" class="size-4" :stroke-width="2.3" />
            <span>{{ choice.label }}</span>
          </button>
        </div>

        <textarea
          v-if="form.question_media === 'code'"
          v-model="form.resource"
          rows="5"
          placeholder="Enter code..."
          class="w-full resize-none border-[3px] border-jv-ink bg-jv-canvas px-3 py-2 font-mono text-[14px] text-jv-ink outline-none focus:shadow-brutal-sm"
        ></textarea>

        <label
          v-else-if="form.question_media === 'image'"
          class="flex min-h-14 cursor-pointer items-center justify-center gap-2 border-2 border-dashed border-jv-ink/35 bg-jv-canvas px-3 py-3 text-[14px] font-bold text-jv-muted transition-colors hover:bg-jv-yellow/20"
        >
          <ImageIcon class="size-4" :stroke-width="2.3" />
          <span>{{ questionFileName || "Upload question image" }}</span>
          <input
            type="file"
            class="hidden"
            accept="image/*"
            @change="handleQuestionImage"
          />
        </label>
      </div>
    </div>

    <div class="mt-5 flex flex-wrap justify-end gap-2">
      <button
        v-for="choice in mediaChoices"
        :key="`options-${choice.value}`"
        type="button"
        class="inline-flex h-9 items-center gap-1.5 border-[2px] border-jv-ink px-3 text-[13px] font-black shadow-[1px_1px_0_#2D2D2D] transition-transform hover:rotate-[1deg]"
        :class="
          form.options_media === choice.value
            ? 'bg-jv-ink text-jv-white'
            : 'bg-jv-white text-jv-ink'
        "
        @click="form.options_media = choice.value"
      >
        <component :is="choice.icon" class="size-4" :stroke-width="2.3" />
        <span>{{ choice.label }}</span>
      </button>
    </div>

    <div class="mt-3 grid gap-3">
      <div
        v-for="key in optionKeys"
        :key="key"
        class="grid grid-cols-[minmax(0,1fr)_28px] items-center gap-3 sm:grid-cols-[100px_28px_minmax(0,1fr)]"
      >
        <label class="text-[14px] font-semibold text-jv-muted">
          Option {{ key }}
        </label>
        <input
          v-if="form.type === 1"
          v-model.number="correctAnswer"
          type="radio"
          :value="Number(key)"
          class="size-6 accent-jv-mint"
          :aria-label="`Correct answer option ${key}`"
        />
        <span v-else class="hidden sm:block"></span>

        <input
          v-if="form.options_media === 'text'"
          v-model.trim="form.options[key]"
          type="text"
          :required="Number(key) <= 2"
          :placeholder="`Enter option ${key}`"
          class="col-span-2 h-11 w-full border-[3px] border-jv-ink bg-jv-white px-3 text-[15px] font-semibold text-jv-ink outline-none placeholder:text-jv-ink/40 focus:shadow-brutal-sm sm:col-span-1"
          :class="
            form.type === 1 && Number(key) === Number(correctAnswer)
              ? 'bg-jv-mint/35'
              : ''
          "
        />

        <textarea
          v-else-if="form.options_media === 'code'"
          v-model="form.options[key]"
          rows="2"
          :placeholder="`Enter code for option ${key}`"
          class="col-span-2 min-w-0 resize-none border-[3px] border-jv-ink bg-jv-white px-3 py-2 font-mono text-[14px] text-jv-ink outline-none focus:shadow-brutal-sm sm:col-span-1"
          :class="
            form.type === 1 && Number(key) === Number(correctAnswer)
              ? 'bg-jv-mint/35'
              : ''
          "
        ></textarea>

        <label
          v-else
          class="col-span-2 flex h-11 min-w-0 cursor-pointer items-center gap-2 border-[3px] border-jv-ink bg-jv-white px-3 text-[14px] font-semibold text-jv-muted focus-within:shadow-brutal-sm sm:col-span-1"
          :class="
            form.type === 1 && Number(key) === Number(correctAnswer)
              ? 'bg-jv-mint/35'
              : ''
          "
        >
          <ImageIcon class="size-4 shrink-0" :stroke-width="2.3" />
          <span class="truncate">{{
            optionFileNames[key] || existingImageLabel(key)
          }}</span>
          <input
            type="file"
            class="hidden"
            accept="image/*"
            @change="handleOptionImage($event, key)"
          />
        </label>
      </div>
    </div>

    <div
      class="mt-6 flex flex-col-reverse gap-3 border-t-2 border-dashed border-jv-ink/15 pt-4 sm:flex-row sm:justify-end"
    >
      <button
        v-if="showCancel"
        type="button"
        class="inline-flex h-11 items-center justify-center border-[3px] border-jv-ink bg-jv-white px-5 text-[15px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[-1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none"
        @click="$emit('cancel')"
      >
        Cancel
      </button>
      <button
        type="submit"
        class="inline-flex h-11 items-center justify-center border-[3px] border-jv-ink bg-jv-mint px-5 text-[15px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none disabled:cursor-not-allowed disabled:opacity-60"
        :disabled="saving"
      >
        {{ saving ? "Saving..." : submitLabel }}
      </button>
    </div>
  </form>
</template>

<script setup>
import { computed, reactive, ref, watch } from "vue";
import { Code2, ImageIcon, Type } from "lucide-vue-next";

const props = defineProps({
  question: {
    type: Object,
    default: null,
  },
  mode: {
    type: String,
    default: "create",
  },
  saving: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["save", "cancel"]);

const mediaChoices = [
  { label: "Text", value: "text", icon: Type },
  { label: "Image", value: "image", icon: ImageIcon },
  { label: "Code", value: "code", icon: Code2 },
];

const form = reactive({
  question: "",
  type: 1,
  question_media: "text",
  options_media: "text",
  resource: "",
  options: {
    1: "",
    2: "",
    3: "",
    4: "",
  },
});

const correctAnswer = ref(1);
const questionFile = ref(null);
const questionFileName = ref("");
const optionFiles = ref({});
const optionFileNames = ref({});

const parseAnswers = (value) => {
  if (Array.isArray(value)) return value;
  if (!value) return [1];
  try {
    const parsed = JSON.parse(value);
    return Array.isArray(parsed) ? parsed : [1];
  } catch {
    return [1];
  }
};

const resetForm = () => {
  const question = props.question;
  form.question = question?.question || "";
  form.type = Number(question?.question_type_id || 1);
  form.question_media = question?.question_media || "text";
  form.options_media = question?.options_media || "text";
  form.resource = question?.resource || "";
  form.options = {
    1: question?.options?.["1"] || "",
    2: question?.options?.["2"] || "",
    3: question?.options?.["3"] || "",
    4: question?.options?.["4"] || "",
  };

  if (question?.options?.["5"]) {
    form.options[5] = question.options["5"];
  }

  const answers = parseAnswers(question?.correct_answer);
  correctAnswer.value = Number(answers[0] || 1);
  questionFile.value = null;
  questionFileName.value = "";
  optionFiles.value = {};
  optionFileNames.value = {};
};

watch(
  () => props.question,
  () => resetForm(),
  { immediate: true, deep: true }
);

const optionKeys = computed(() =>
  Object.keys(form.options)
    .sort((a, b) => Number(a) - Number(b))
    .filter((key) => Number(key) <= 5)
);

const title = computed(() =>
  props.mode === "edit" ? "Edit Question" : "New Question"
);
const eyebrow = computed(() => (props.mode === "edit" ? "Question" : "Create"));
const submitLabel = computed(() =>
  props.mode === "edit" ? "Save Changes" : "Add Question"
);
const showCancel = computed(
  () => props.mode === "edit" || props.mode === "create"
);

const handleQuestionImage = (event) => {
  const file = event.target.files?.[0];
  if (!file) return;
  questionFile.value = file;
  questionFileName.value = file.name;
};

const handleOptionImage = (event, key) => {
  const file = event.target.files?.[0];
  if (!file) return;
  optionFiles.value = {
    ...optionFiles.value,
    [key]: file,
  };
  optionFileNames.value = {
    ...optionFileNames.value,
    [key]: file.name,
  };
};

const existingImageLabel = (key) => {
  if (form.options[key]) return `Option ${key} image`;
  return `Upload image for option ${key}`;
};

const submitForm = () => {
  const nonEmptyOptionKeys = optionKeys.value.filter(
    (key) => form.options[key] || form.options_media === "image"
  );
  const answers =
    form.type === 1
      ? [Number(correctAnswer.value)]
      : nonEmptyOptionKeys.map((key) => Number(key));

  emit("save", {
    payload: {
      question: form.question,
      type: Number(form.type),
      options: { ...form.options },
      answers,
      question_media: form.question_media,
      options_media: form.options_media,
      resource: form.resource,
    },
    files: {
      question: questionFile.value,
      options: { ...optionFiles.value },
    },
  });
};
</script>
