<template>
  <form
    class="jv-border-rough bg-jv-white p-4 shadow-brutal-sm sm:p-5 lg:p-6"
    @submit.prevent="submitForm"
  >
    <div
      class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
    >
      <div>
        <p
          class="text-[12px] font-bold uppercase tracking-[0.14em] text-jv-coral"
        >
          {{ eyebrow }}
        </p>
        <h2
          class="mt-1 font-body text-[22px] font-bold leading-tight text-jv-ink sm:text-[26px]"
        >
          {{ title }}
        </h2>
      </div>

      <select
        v-model.number="form.type"
        class="h-10 rounded-md border border-jv-ink/30 bg-jv-white px-3 text-[14px] font-semibold text-jv-ink outline-none focus:border-jv-ink"
      >
        <option :value="1">Multiple Choice</option>
        <option :value="2">Survey</option>
      </select>
    </div>

    <div class="mt-6 space-y-2">
      <label
        class="text-[13px] font-semibold uppercase tracking-wide text-jv-muted"
      >
        Question
      </label>
      <input
        v-model.trim="form.question"
        type="text"
        required
        placeholder="Enter question..."
        class="h-11 w-full rounded-md border border-jv-ink/30 bg-jv-white px-3 text-[15px] font-medium text-jv-ink outline-none placeholder:text-jv-ink/40 focus:border-jv-ink focus:shadow-brutal-sm"
      />

      <div class="flex flex-wrap items-center gap-1.5 pt-1">
        <NavigationLink
          v-for="choice in mediaChoices"
          :key="`question-${choice.value}`"
          variant="toggle"
          type="button"
          :class="
            form.question_media === choice.value
              ? 'border-jv-ink bg-jv-ink text-jv-white'
              : 'border-jv-ink/25 bg-jv-white text-jv-muted hover:border-jv-ink/60 hover:text-jv-ink'
          "
          @click="form.question_media = choice.value"
        >
          <component :is="choice.icon" class="size-3.5" :stroke-width="2.3" />
          <span>{{ choice.label }}</span>
        </NavigationLink>
      </div>

      <textarea
        v-if="form.question_media === 'code'"
        v-model="form.resource"
        rows="5"
        placeholder="Enter code..."
        class="mt-1 w-full resize-none rounded-md border border-jv-ink/30 bg-jv-canvas px-3 py-2 font-mono text-[14px] text-jv-ink outline-none focus:border-jv-ink focus:shadow-brutal-sm"
      ></textarea>

      <label
        v-else-if="form.question_media === 'image'"
        class="mt-1 flex min-h-12 cursor-pointer items-center justify-center gap-2 rounded-md border border-dashed border-jv-ink/35 bg-jv-canvas px-3 py-3 text-[14px] font-medium text-jv-muted transition-colors hover:bg-jv-yellow/20"
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

    <div class="mt-6 space-y-3">
      <div class="flex flex-wrap items-center justify-between gap-2">
        <label
          class="text-[13px] font-semibold uppercase tracking-wide text-jv-muted"
        >
          Options
        </label>
        <div class="flex flex-wrap items-center gap-1.5">
          <NavigationLink
            v-for="choice in mediaChoices"
            :key="`options-${choice.value}`"
            variant="toggle"
            type="button"
            :class="
              form.options_media === choice.value
                ? 'border-jv-ink bg-jv-ink text-jv-white'
                : 'border-jv-ink/25 bg-jv-white text-jv-muted hover:border-jv-ink/60 hover:text-jv-ink'
            "
            @click="form.options_media = choice.value"
          >
            <component :is="choice.icon" class="size-3.5" :stroke-width="2.3" />
            <span>{{ choice.label }}</span>
          </NavigationLink>
        </div>
      </div>

      <div class="flex flex-col">
        <div
          v-for="key in optionKeys"
          :key="key"
          class="flex min-w-0 items-center gap-3 border-b border-jv-ink/10 py-2.5 pl-2 pr-1 last:border-b-0"
          :class="
            form.type === 1 && Number(key) === Number(correctAnswer)
              ? 'border-l-4 border-l-jv-accent-green bg-jv-accent-green/25 pl-1'
              : 'border-l-4 border-l-transparent'
          "
        >
          <input
            v-if="form.type === 1"
            v-model.number="correctAnswer"
            type="radio"
            :value="Number(key)"
            class="size-5 shrink-0 accent-jv-accent-green"
            :aria-label="`Correct answer option ${key}`"
          />

          <span class="w-6 shrink-0 text-[14px] font-bold text-jv-coral">
            {{ optionLetter(key) }}.
          </span>

          <input
            v-if="form.options_media === 'text'"
            v-model.trim="form.options[key]"
            type="text"
            :required="Number(key) <= 2"
            :placeholder="`Option ${optionLetter(key)}`"
            class="h-10 w-full min-w-0 flex-1 rounded-md border border-jv-ink/25 bg-jv-white px-3 text-[15px] font-medium text-jv-ink outline-none placeholder:text-jv-ink/40 focus:border-jv-ink focus:shadow-brutal-sm"
          />

          <textarea
            v-else-if="form.options_media === 'code'"
            v-model="form.options[key]"
            rows="2"
            :placeholder="`Code for option ${optionLetter(key)}`"
            class="min-w-0 flex-1 resize-none rounded-md border border-jv-ink/25 bg-jv-white px-3 py-2 font-mono text-[13px] text-jv-ink outline-none focus:border-jv-ink focus:shadow-brutal-sm"
          ></textarea>

          <label
            v-else
            class="flex h-10 min-w-0 flex-1 cursor-pointer items-center gap-2 rounded-md border border-dashed border-jv-ink/30 bg-jv-white px-3 text-[13px] font-medium text-jv-muted focus-within:shadow-brutal-sm hover:bg-jv-yellow/10"
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
    </div>

    <div class="mt-6 flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">
      <NavigationLink
        v-if="showCancel"
        class="h-11 bg-jv-white text-jv-ink"
        @click="$emit('cancel')"
      >
        Cancel
      </NavigationLink>
      <NavigationLink class="h-11 bg-jv-accent-green text-white">
        {{ saving ? "Saving..." : submitLabel }}
      </NavigationLink>
    </div>
  </form>
</template>

<script setup>
import { computed, reactive, ref, watch } from "vue";
import { Code2, ImageIcon, Type } from "lucide-vue-next";
import { usePush } from "notivue";
import NavigationLink from "../common/NavigationLink.vue";

const app = useNuxtApp();
const toast = usePush();
const { maxImageFileSize } = useRuntimeConfig().public;

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
const questionFileName = ref("");
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
  questionFileName.value = "";
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

const optionLetter = (key) => String.fromCharCode(64 + Number(key));

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

const validateImageFile = (file) => {
  if (!app.$validImageTypes.includes(file.type)) {
    toast.error(
      "Please upload a valid image file (JPEG, PNG, GIF, WEBP, HEIC, HEIF)."
    );
    return false;
  }
  if (file.size > maxImageFileSize) {
    const limitKb = Math.round(maxImageFileSize / 1024);
    toast.error(`Please upload an image less than ${limitKb} KB.`);
    return false;
  }
  return true;
};

const readAsBase64 = (file) =>
  new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (e) => resolve(e.target.result);
    reader.onerror = () => reject(reader.error);
    reader.readAsDataURL(file);
  });

const handleQuestionImage = async (event) => {
  const file = event.target.files?.[0];
  if (!file) return;
  if (!validateImageFile(file)) {
    event.target.value = "";
    return;
  }
  form.resource = await readAsBase64(file);
  questionFileName.value = file.name;
};

const handleOptionImage = async (event, key) => {
  const file = event.target.files?.[0];
  if (!file) return;
  if (!validateImageFile(file)) {
    event.target.value = "";
    return;
  }
  form.options[key] = await readAsBase64(file);
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
  });
};
</script>
