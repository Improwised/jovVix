<template>
  <main class="min-h-screen bg-jv-canvas px-4 py-5 sm:px-6 md:px-8 md:py-6">
    <div v-if="quizPending" class="py-8">
      <UtilsQuizListWaiting />
    </div>

    <div
      v-else-if="quizError?.data?.code == 401"
      class="jv-border-rough bg-jv-white p-5 text-[18px] font-semibold text-jv-coral shadow-brutal-sm"
    >
      {{ navigateTo("/account/login") }}
    </div>

    <div
      v-else-if="quizError"
      class="jv-border-rough bg-jv-white p-5 text-[18px] font-semibold text-jv-coral shadow-brutal-sm"
    >
      {{ quizError.message || quizError }}
    </div>

    <section v-else class="mx-auto flex flex-col gap-6">
      <header
        class="flex flex-col gap-4 md:flex-row md:items-start md:justify-between"
      >
        <div class="min-w-0">
          <h1
            class="break-words font-headings text-[38px] leading-none text-jv-ink min-[420px]:text-[44px] sm:text-[52px] md:text-[56px]"
          >
            {{ quizTitle }}
          </h1>
          <div
            class="ml-1 mt-1 h-3 w-24 rounded-full border-b-[3px] border-jv-coral sm:ml-2 sm:w-32"
          ></div>
          <p
            class="mt-3 max-w-2xl text-[16px] leading-[1.5] text-jv-muted sm:text-[18px]"
          >
            {{ quizDescription || "No description added yet." }}
          </p>
        </div>
        <NavigationLink
          url-name="Start Quiz"
          class="rounded-[999px] bg-jv-coral text-white"
          @click="handleStartQuiz"
        >
          <Play class="size-4 fill-current" :stroke-width="2.4" />
        </NavigationLink>
      </header>

      <div
        class="flex flex-col gap-4 jv-border-rough bg-jv-white p-3 shadow-brutal-sm sm:p-4 lg:flex-row lg:items-center lg:justify-between"
      >
        <div class="flex flex-wrap gap-2">
          <span
            class="inline-flex items-center gap-1.5 rounded-[6px] border border-dashed border-jv-ink/35 bg-jv-white px-3 py-2 text-[14px] font-semibold text-jv-muted"
          >
            Total Questions:
            <strong class="text-jv-ink">{{ questions.length }}</strong>
          </span>
          <span
            class="inline-flex items-center gap-1.5 rounded-[6px] border border-dashed border-jv-ink/35 bg-jv-white px-3 py-2 text-[14px] font-semibold text-jv-muted"
          >
            Survey Questions:
            <strong class="text-jv-ink">{{ totalSurveyQuestion }}</strong>
          </span>
          <span
            class="inline-flex items-center gap-1.5 rounded-[6px] border border-dashed border-jv-ink/35 bg-jv-white px-3 py-2 text-[14px] font-semibold text-jv-muted"
          >
            Played Quiz:
            <strong class="text-jv-ink">{{
              quizData?.data?.quiz_played_count || 0
            }}</strong>
          </span>
        </div>

        <div class="flex flex-wrap gap-2 sm:justify-end">
          <NavigationLink
            url-name="Delete Quiz"
            class="bg-jv-white text-jv-coral border-none shadow-none"
            @click="deleteQuiz"
          >
            <Trash2 class="size-4" :stroke-width="2.4" />
          </NavigationLink>
          <NavigationLink
            url-name="Share Quiz"
            class="bg-jv-white"
            @click="handleShareQuiz"
          >
            <Share2 class="size-4" :stroke-width="2.4" />
          </NavigationLink>
        </div>
      </div>

      <div class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_auto] lg:items-center">
        <div class="flex flex-wrap items-center gap-4">
          <label
            class="inline-flex items-center gap-2 text-[15px] font-semibold text-jv-muted"
          >
            <span>Points</span>
            <input
              v-model.number="settings.points"
              type="number"
              min="0"
              max="20"
              :disabled="!canEditQuiz || settingsPending"
              class="h-10 w-14 border-[3px] border-jv-ink bg-jv-white px-2 text-center text-[16px] font-black text-jv-ink shadow-brutal-sm outline-none disabled:opacity-60"
              @change="saveSettings"
            />
          </label>
          <label
            class="inline-flex items-center gap-2 text-[15px] font-semibold text-jv-muted"
          >
            <span>Duration (seconds)</span>
            <input
              v-model.number="settings.duration_in_seconds"
              type="number"
              min="1"
              :disabled="!canEditQuiz || settingsPending"
              class="h-10 w-16 border-[3px] border-jv-ink bg-jv-white px-2 text-center text-[16px] font-black text-jv-ink shadow-brutal-sm outline-none disabled:opacity-60"
              @change="saveSettings"
            />
          </label>
        </div>

        <div
          v-if="canEditQuiz"
          class="flex flex-col gap-3 sm:flex-row sm:justify-end"
        >
          <NavigationLink
            url-name="Import Question by .CSV"
            class="rounded-[999px] bg-jv-white"
            @click="importModalOpen = true"
          >
            <Upload class="size-4" :stroke-width="2.4" />
          </NavigationLink>
          <NavigationLink
            url-name="Add Question"
            class="rounded-[999px] bg-jv-mint"
            @click="openNewQuestionForm"
          >
            <Plus class="size-4" :stroke-width="2.4" />
          </NavigationLink>
        </div>
      </div>

      <QuestionFormCard
        v-if="showNewQuestionForm && canEditQuiz"
        mode="create"
        :saving="savingNewQuestion"
        @save="saveNewQuestion"
        @cancel="showNewQuestionForm = false"
      />

      <section
        v-if="questions.length === 0 && !showNewQuestionForm"
        class="grid min-h-[260px] place-items-center text-center"
      >
        <div class="max-w-md py-8">
          <h2
            class="font-headings text-[30px] leading-tight text-jv-ink sm:text-[36px]"
          >
            No questions yet
          </h2>
          <p class="mt-2 text-[17px] text-jv-muted">
            Create one manually or import a CSV to start building this quiz.
          </p>
        </div>
      </section>

      <section v-else class="grid gap-5">
        <article
          v-for="(question, index) in questions"
          :key="question.question_id"
          class="jv-border-rough bg-jv-white p-4 shadow-brutal-sm transition-transform hover:rotate-[0.25deg] sm:p-5"
          :class="
            editingQuestionId === question.question_id ? 'bg-jv-yellow/10' : ''
          "
          @click="canEditQuiz && startEditQuestion(question)"
        >
          <QuestionFormCard
            v-if="editingQuestionId === question.question_id"
            mode="edit"
            :question="question"
            :saving="savingQuestionId === question.question_id"
            @save="(data) => saveExistingQuestion(question, data)"
            @cancel="editingQuestionId = ''"
            @click.stop
          />

          <template v-else>
            <div class="flex items-start justify-between gap-4">
              <div class="min-w-0">
                <p class="text-[14px] font-black text-jv-coral">
                  Question: {{ index + 1 }}
                </p>
                <h2
                  class="mt-1 break-words text-[20px] font-black leading-snug text-jv-ink sm:text-[22px]"
                >
                  {{ question.question }}
                </h2>
              </div>
              <div v-if="canEditQuiz" class="flex shrink-0 gap-2">
                <NavigationLink
                  class="rounded-full size-10 px-2 py-1 sm:px-2 sm:py-1 md:px-2 md:py-1"
                  aria-label="Edit question"
                  @click.stop="startEditQuestion(question)"
                >
                  <Pencil class="size-4" :stroke-width="2.4" />
                </NavigationLink>
                <NavigationLink
                  class="rounded-full bg-jv-coral text-white size-10 px-2 py-1 sm:px-2 sm:py-1 md:px-2 md:py-1"
                  aria-label="Delete question"
                  @click.stop="deleteQuestion(question.question_id)"
                >
                  <Trash2 class="size-4" :stroke-width="2.4" />
                </NavigationLink>
              </div>
            </div>

            <div
              v-if="question.question_media === 'image' && question.resource"
              class="mt-3 flex justify-center border-[3px] border-jv-ink bg-jv-canvas p-2"
              @click.stop
            >
              <img
                :src="question.resource"
                :alt="question.question"
                class="max-h-64 w-auto max-w-full object-contain"
              />
            </div>

            <div
              v-else-if="
                question.question_media === 'code' && question.resource
              "
              class="mt-3 min-w-0 overflow-x-auto border-none border-jv-ink bg-jv-canvas"
              @click.stop
            >
              <CodeBlockComponent :code="question.resource" />
            </div>

            <div class="mt-4 grid gap-3 sm:grid-cols-2">
              <div
                v-for="(option, key) in question.options"
                :key="key"
                class="flex min-w-0 min-h-12 items-center gap-3 border-[3px] border-jv-ink bg-jv-white px-3 py-2 text-[15px] font-semibold text-jv-ink"
                :class="isCorrectAnswer(question, key) ? 'bg-jv-mint/45' : ''"
              >
                <span
                  class="grid size-7 shrink-0 place-items-center rounded-full border-2 border-jv-ink bg-jv-white text-[13px] font-black"
                >
                  {{ optionLabel(key) }}
                </span>

                <div
                  v-if="question.options_media === 'image' && option"
                  class="flex min-w-0 flex-1 justify-center"
                  @click.stop
                >
                  <img
                    :src="option"
                    :alt="`Option ${optionLabel(key)}`"
                    class="max-h-32 w-auto max-w-full object-contain"
                  />
                </div>

                <div
                  v-else-if="question.options_media === 'code' && option"
                  class="min-w-0 flex-1 overflow-x-auto"
                  @click.stop
                >
                  <CodeBlockComponent :code="option" />
                </div>

                <span v-else class="min-w-0 flex-1 break-words">
                  {{ option }}
                </span>
              </div>
            </div>
          </template>
        </article>
      </section>
    </section>

    <Teleport to="body">
      <div
        v-if="importModalOpen"
        class="fixed inset-0 z-50 grid place-items-center bg-jv-ink/35 px-4 py-6 backdrop-blur-[2px]"
        @click.self="closeImportModal"
      >
        <form
          class="w-full max-w-[680px] rotate-[-0.4deg] border-[4px] border-jv-ink bg-jv-white shadow-brutal-lg"
          @submit.prevent="importCsv"
        >
          <div
            class="flex items-center justify-between gap-4 border-b-[3px] border-jv-ink bg-jv-ink px-5 py-4 text-jv-white sm:px-6"
          >
            <h2
              class="font-body text-[24px] font-black leading-none text-jv-white sm:text-[28px]"
            >
              Add Question to Quiz
            </h2>
            <button
              type="button"
              class="grid size-9 place-items-center text-jv-white transition-transform hover:rotate-[6deg]"
              aria-label="Close import CSV modal"
              @click="closeImportModal"
            >
              <X class="size-6" :stroke-width="2.4" />
            </button>
          </div>

          <div class="grid gap-4 px-5 py-6 sm:px-8">
            <label class="grid gap-2">
              <span
                class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
              >
                Choose CSV File <span class="text-jv-coral">*</span>
              </span>
              <span
                class="flex h-14 cursor-pointer border-[3px] border-jv-ink bg-jv-canvas"
              >
                <span
                  class="inline-flex h-full items-center gap-2 bg-jv-ink px-4 text-[16px] font-black text-jv-white"
                >
                  <Upload class="size-4" :stroke-width="2.4" />
                  Choose File
                </span>
                <span
                  class="flex min-w-0 flex-1 items-center px-4 text-[15px] font-semibold text-jv-muted"
                >
                  {{ csvFileName || "No file chosen" }}
                </span>
                <input
                  type="file"
                  class="hidden"
                  accept=".csv,text/csv"
                  required
                  @change="handleCsvFile"
                />
              </span>
            </label>
            <p class="text-[15px] leading-[1.6] text-jv-muted">
              Upload a CSV file containing additional questions. New rows will
              be appended to this quiz.
            </p>
          </div>

          <div
            class="flex flex-col-reverse gap-3 border-t-[3px] border-jv-ink bg-jv-canvas px-5 py-4 sm:flex-row sm:justify-end sm:px-8"
          >
            <button
              type="button"
              class="inline-flex h-12 items-center justify-center border-[3px] border-jv-ink bg-jv-white px-6 text-[17px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[-1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none"
              @click="closeImportModal"
            >
              Cancel
            </button>
            <button
              type="submit"
              class="inline-flex h-12 items-center justify-center border-[3px] border-jv-ink bg-jv-mint px-6 text-[17px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none disabled:cursor-not-allowed disabled:opacity-60"
              :disabled="importPending"
            >
              {{ importPending ? "Adding..." : "Add Questions" }}
            </button>
          </div>
        </form>
      </div>
    </Teleport>
  </main>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import { Pencil, Play, Plus, Share2, Trash2, Upload, X } from "lucide-vue-next";
import { useToast } from "vue-toastification";
import QuestionFormCard from "@/components/quiz-manage/QuestionFormCard.vue";
import CodeBlockComponent from "@/components/CodeBlockComponent.vue";
import usecopyToClipboard from "@/composables/copy_to_clipboard";
import { useListUserstore } from "~/store/userlist";
import { useSessionStore } from "~~/store/session";
import NavigationLink from "~/components/common/NavigationLink.vue";

definePageMeta({
  layout: "empty",
});

const app = useNuxtApp();
const toast = useToast();
const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const route = useRoute();
const router = useRouter();
const sessionStore = useSessionStore();
const listUserStore = useListUserstore();

const quizId = computed(() => route.params.quiz_id || "");
const importModalOpen = ref(false);
const importPending = ref(false);
const csvFile = ref(null);
const csvFileName = ref("");
const showNewQuestionForm = ref(false);
const showAddAnother = ref(false);
const savingNewQuestion = ref(false);
const editingQuestionId = ref("");
const savingQuestionId = ref("");
const settingsPending = ref(false);
const startingQuiz = ref(false);
const settings = ref({
  points: 10,
  duration_in_seconds: 10,
});

const {
  refresh,
  data: quizData,
  pending: quizPending,
  error: quizError,
} = useFetch(`${url.apiUrl}/quizzes/${quizId.value}/questions`, {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
});

const questions = computed(() => quizData.value?.data?.data || []);
const quizTitle = computed(() => quizData.value?.data?.quiz_title || "Quiz");
const quizDescription = computed(
  () => quizData.value?.data?.quiz_description?.String || "",
);
const totalSurveyQuestion = computed(() =>
  questions.value.reduce(
    (count, item) => (item.question_type_id === 2 ? count + 1 : count),
    0,
  ),
);
const canEditQuiz = computed(() => {
  const permission = quizData.value?.data?.permission;
  const isEditable = quizData.value?.data?.is_quiz_editable;
  return (permission === "write" || permission === "share") && isEditable;
});

watch(
  () => quizData.value?.data,
  (data) => {
    if (!data) return;
    settings.value = {
      points: Number(data.points ?? 10),
      duration_in_seconds: Number(data.duration_in_seconds ?? 10),
    };
  },
  { immediate: true },
);

const parseAnswers = (value) => {
  if (Array.isArray(value)) return value;
  if (!value) return [];
  try {
    const parsed = JSON.parse(value);
    return Array.isArray(parsed) ? parsed : [];
  } catch {
    return [];
  }
};

const optionLabel = (key) => String.fromCharCode(64 + Number(key));

const isCorrectAnswer = (question, key) =>
  parseAnswers(question.correct_answer).includes(Number(key));

const openNewQuestionForm = () => {
  showAddAnother.value = false;
  showNewQuestionForm.value = true;
  editingQuestionId.value = "";
};

const startEditQuestion = (question) => {
  showNewQuestionForm.value = false;
  showAddAnother.value = false;
  editingQuestionId.value = question.question_id;
};

const saveSettings = async () => {
  if (!canEditQuiz.value) return;

  try {
    settingsPending.value = true;
    await $fetch(`${url.apiUrl}/quizzes/${quizId.value}/settings`, {
      method: "PUT",
      headers,
      credentials: "include",
      body: settings.value,
    });
    toast.success("Quiz settings updated.");
    refresh();
  } catch (error) {
    toast.error(
      error?.data?.message ||
        error?.message ||
        "Failed to update quiz settings.",
    );
  } finally {
    settingsPending.value = false;
  }
};

const validateImage = (file) => {
  if (!file) return true;

  if (!app.$validImageTypes.includes(file.type)) {
    toast.error("Please upload a valid image file.");
    return false;
  }
  if (file.size > url.maxImageFileSize) {
    toast.error(
      `Please upload an image less than ${
        url.maxImageFileSize / 1024 / 1024
      } MB.`,
    );
    return false;
  }
  return true;
};

const uploadImages = async (questionId, files) => {
  const imageForm = new FormData();
  let hasImages = false;

  if (files?.question) {
    if (!validateImage(files.question)) return false;
    imageForm.append("image-attachment", files.question, questionId);
    hasImages = true;
  }

  Object.entries(files?.options || {}).forEach(([key, file]) => {
    if (!validateImage(file)) return;
    imageForm.append("image-attachment", file, `${key}_${questionId}`);
    hasImages = true;
  });

  if (!hasImages) return true;

  await $fetch(`${url.apiUrl}/images?quiz_id=${quizId.value}`, {
    method: "POST",
    headers: {
      Accept: "application/json",
    },
    body: imageForm,
    mode: "cors",
    credentials: "include",
  });

  return true;
};

const saveNewQuestion = async ({ payload, files }) => {
  try {
    savingNewQuestion.value = true;
    const payloadWithDefaults = {
      ...payload,
      points: Number(settings.value.points),
      duration_in_seconds: Number(settings.value.duration_in_seconds),
    };
    const response = await $fetch(
      `${url.apiUrl}/quizzes/${quizId.value}/questions`,
      {
        method: "POST",
        headers,
        credentials: "include",
        body: payloadWithDefaults,
      },
    );

    const questionId = response?.data;
    if (questionId) {
      await uploadImages(questionId, files);
    }

    toast.success("Question added successfully.");
    showNewQuestionForm.value = false;
    showAddAnother.value = true;
    refresh();
  } catch (error) {
    toast.error(
      error?.data?.message || error?.message || "Failed to add question.",
    );
  } finally {
    savingNewQuestion.value = false;
  }
};

const buildUpdatePayload = (question, payload) => {
  const questionId = question.question_id;
  const options = { ...payload.options };

  if (payload.options_media === "image") {
    Object.keys(options).forEach((key) => {
      options[key] = `${quizId.value}/${key}_${questionId}`;
    });
  }

  return {
    ...payload,
    options,
    points: Number(settings.value.points),
    duration_in_seconds: Number(settings.value.duration_in_seconds),
    resource:
      payload.question_media === "image"
        ? `${quizId.value}/${questionId}`
        : payload.resource,
  };
};

const saveExistingQuestion = async (question, { payload, files }) => {
  const questionId = question.question_id;

  try {
    savingQuestionId.value = questionId;
    await uploadImages(questionId, files);
    await $fetch(
      `${url.apiUrl}/quizzes/${quizId.value}/questions/${questionId}`,
      {
        method: "PUT",
        headers,
        credentials: "include",
        body: buildUpdatePayload(question, payload),
      },
    );

    toast.success("Question updated successfully.");
    editingQuestionId.value = "";
    refresh();
  } catch (error) {
    toast.error(
      error?.data?.message || error?.message || "Failed to update question.",
    );
  } finally {
    savingQuestionId.value = "";
  }
};

const deleteQuestion = async (questionId) => {
  if (!window.confirm("Delete this question?")) return;

  try {
    await $fetch(
      `${url.apiUrl}/quizzes/${quizId.value}/questions/${questionId}`,
      {
        method: "DELETE",
        headers,
        credentials: "include",
      },
    );
    toast.success("Question deleted successfully.");
    refresh();
  } catch (error) {
    toast.error(
      error?.data?.message || error?.message || "Failed to delete question.",
    );
  }
};

const deleteQuiz = async () => {
  if (!window.confirm("Delete this quiz?")) return;

  try {
    await $fetch(`${url.apiUrl}/quizzes/${quizId.value}`, {
      method: "DELETE",
      headers,
      credentials: "include",
    });
    toast.success("Quiz deleted successfully.");
    router.push("/admin/quiz/list-quiz");
  } catch (error) {
    toast.error(
      error?.data?.message || error?.message || "Failed to delete quiz.",
    );
  }
};

const handleCsvFile = (event) => {
  const file = event.target.files?.[0];
  csvFile.value = file || null;
  csvFileName.value = file?.name || "";
};

const closeImportModal = () => {
  importModalOpen.value = false;
  csvFile.value = null;
  csvFileName.value = "";
};

const importCsv = async () => {
  if (!csvFile.value) {
    toast.error("Please select a CSV file.");
    return;
  }

  try {
    importPending.value = true;
    const formData = new FormData();
    formData.append("attachment", csvFile.value);

    await $fetch(`${url.apiUrl}/quizzes/${quizId.value}/questions/upload`, {
      method: "POST",
      headers: {
        Accept: "application/json",
      },
      body: formData,
      credentials: "include",
    });

    toast.success("Questions imported successfully.");
    closeImportModal();
    refresh();
  } catch (error) {
    toast.error(
      error?.data?.data ||
        error?.data?.message ||
        error?.message ||
        "Failed to import CSV.",
    );
  } finally {
    importPending.value = false;
  }
};

const handleShareQuiz = () => {
  if (!import.meta.client) return;
  usecopyToClipboard(window.location.href);
};

const handleStartQuiz = async () => {
  if (questions.value.length === 0) {
    toast.error("Add at least one question before starting the quiz.");
    return;
  }

  try {
    startingQuiz.value = true;
    const response = await $fetch(
      `${url.apiUrl}/quizzes/${quizId.value}/demo_session`,
      {
        method: "POST",
        credentials: "include",
        headers: {
          Accept: "application/json",
        },
      },
    );

    const activeQuizId = response?.data;
    if (!activeQuizId) {
      toast.error("Error while starting quiz.");
      return;
    }

    listUserStore.removeAllUsers();
    setSocketObject(null);
    sessionStore.setSession(activeQuizId);
    router.push(`/admin/arrange/${activeQuizId}`);
  } catch (error) {
    toast.error(error?.message || "Error while starting quiz.");
  } finally {
    startingQuiz.value = false;
  }
};
</script>
