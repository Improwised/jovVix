<template>
  <div class="flex flex-col gap-5">
    <div class="flex flex-col gap-2">
      <label
        for="edit-question-text"
        class="font-body text-xs font-bold uppercase tracking-wide text-jv-ink"
      >
        Question
      </label>
      <input
        id="edit-question-text"
        v-model="editableQuestion.question"
        placeholder="Edit Question"
        class="jv-card w-full border-2 border-jv-ink bg-jv-white px-4 py-3 font-body text-base font-bold text-jv-ink shadow-brutal-sm outline-none transition-all focus:translate-x-[1px] focus:translate-y-[1px] focus:shadow-none"
      />
    </div>

    <div
      v-if="editableQuestion.question_media === 'image'"
      class="flex flex-col gap-3"
    >
      <label
        class="jv-card flex w-fit cursor-pointer items-center gap-2 border-2 border-jv-ink bg-jv-yellow px-4 py-2 font-headings text-sm text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none"
      >
        <Camera class="size-4" :stroke-width="2.4" />
        Replace Image
        <input
          id="image-attachment-question"
          type="file"
          accept="image/*"
          class="hidden"
          :name="editableQuestion.question_id"
          @change="onImageChange"
        />
      </label>
      <div v-if="editableQuestion.resource" class="flex justify-center">
        <img
          :src="editableQuestion.resource"
          :alt="editableQuestion.question"
          class="max-h-[240px] w-auto border-2 border-jv-ink bg-jv-white object-contain shadow-brutal-sm"
        />
      </div>
    </div>

    <CodeBlockComponent
      v-if="editableQuestion.question_media === 'code'"
      :code="editableQuestion.resource"
      :read-only="false"
      @code-change="changeCode"
    />

    <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
      <div
        v-for="(option, order) in props.question.options"
        :key="order"
        class="jv-card flex flex-col gap-3 border-2 border-jv-ink bg-jv-white p-4 shadow-brutal-sm"
      >
        <div
          v-if="props.question.options_media === 'image'"
          class="flex flex-col gap-2"
        >
          <label
            class="jv-card flex w-fit cursor-pointer items-center gap-2 border-2 border-jv-ink bg-jv-yellow-soft px-3 py-1.5 font-headings text-xs text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none"
          >
            <Camera class="size-3.5" :stroke-width="2.4" />
            Option {{ order }}
            <input
              :id="`image-attachment-option-${order}`"
              type="file"
              accept="image/*"
              class="hidden"
              :name="`${order}_${props.question.question_id}`"
              @change="onImageChange"
            />
          </label>
          <div v-if="editableOptions[order]" class="flex justify-center">
            <img
              :src="editableOptions[order]"
              :alt="option"
              class="max-h-[160px] w-auto border-2 border-jv-ink bg-jv-white object-contain shadow-brutal-sm"
            />
          </div>
        </div>

        <CodeBlockComponent
          v-if="editableQuestion.options_media === 'code'"
          :code="option"
          :read-only="false"
          :option-order="Number(order)"
          @code-change="changeCode"
        />

        <input
          v-if="editableQuestion.options_media === 'text'"
          v-model="editableOptions[order]"
          placeholder="Edit Option"
          class="jv-card w-full border-2 border-jv-ink bg-jv-white px-3 py-2 font-body text-sm text-jv-ink shadow-brutal-sm outline-none transition-all focus:translate-x-[1px] focus:translate-y-[1px] focus:shadow-none"
        />

        <label
          v-if="props.question.question_type_id == 1"
          class="inline-flex cursor-pointer items-center gap-2 font-body text-sm font-bold text-jv-ink"
        >
          <input
            :id="`correct-${order}`"
            v-model="picked"
            type="radio"
            :value="order"
            class="size-4 accent-jv-coral"
          />
          Correct answer
        </label>
      </div>
    </div>

    <div class="flex justify-end">
      <button
        type="button"
        class="jv-card inline-flex h-11 items-center justify-center gap-2 border-2 border-jv-ink bg-jv-coral px-6 font-headings text-base text-white shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none"
        @click="updateQuestion"
      >
        <Save class="size-4" :stroke-width="2.4" />
        Save Changes
      </button>
    </div>
  </div>
</template>

<script setup>
import { usePush } from "notivue";
import { Camera, Save } from "lucide-vue-next";

const app = useNuxtApp();
const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const toast = usePush();

const props = defineProps({
  question: {
    type: Object,
    required: true,
    default: () => ({}),
  },
  quizId: {
    type: String,
    required: false,
    default: "",
  },
  questionId: {
    type: String,
    required: false,
    default: "",
  },
});

const editableQuestion = ref({ ...props.question });
const editableOptions = ref({ ...props.question.options });
const correctAnswers = ref(props.question.correct_answer);
const checkedNames = ref(JSON.parse(correctAnswers.value));
const picked = ref(checkedNames.value[0]);

const changeCode = (data, order) => {
  if (order) {
    editableOptions.value[order] = data;
    return;
  }
  editableQuestion.value.resource = data;
};

const onImageChange = (e) => {
  if (e.target.files.length === 0) {
    toast.error("Please select a file to upload.");
    return;
  }

  const file = e.target.files[0];

  if (!app.$validImageTypes.includes(file.type)) {
    toast.error(
      "Please upload a valid image file (JPEG, PNG, GIF, WEBP, HEIC, HEIF)."
    );
    return;
  }

  if (file.size > url.maxImageFileSize) {
    const limitKb = Math.round(url.maxImageFileSize / 1024);
    toast.error(`Please upload an image less than ${limitKb} KB.`);
    return;
  }

  const reader = new FileReader();
  reader.onload = (ev) => {
    if (e.target.id.startsWith("image-attachment-option")) {
      const order = e.target.name[0];
      editableOptions.value[order] = ev.target.result;
    } else {
      editableQuestion.value.resource = ev.target.result;
    }
  };
  reader.readAsDataURL(file);
};

const updateQuestion = async () => {
  try {
    const payload = {
      question: editableQuestion.value.question,
      type: props.question.question_type_id,
      options: editableOptions.value,
      answers:
        props.question.question_type_id === 1
          ? [Number(picked.value)]
          : checkedNames.value.map(Number),
      points: props.question.points,
      duration_in_seconds: props.question.duration_in_seconds,
      question_media: props.question.question_media,
      options_media: props.question.options_media,
      resource: editableQuestion.value.resource,
    };

    await $fetch(
      `${url.apiUrl}/quizzes/${props.quizId}/questions/${props.questionId}`,
      {
        method: "PUT",
        headers: headers,
        body: payload,
        credentials: "include",
      }
    );
    toast.success("Question updated successfully!");
  } catch (error) {
    console.error("Failed to update the question", error);
    toast.error("Failed to update the question.");
  }
};
</script>
