<template>
  <main class="min-h-screen bg-jv-canvas px-4 py-5 sm:px-6 md:px-8 md:py-6">
    <div
      v-if="quizError?.data?.code == 401"
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
          <div class="flex flex-wrap items-center gap-3">
            <h1
              class="break-words font-headings text-[38px] leading-none text-jv-ink min-[420px]:text-[44px] sm:text-[52px] md:text-[56px]"
            >
              {{ quizTitle }}
            </h1>
            <span
              v-if="isPublicQuiz"
              class="inline-flex rotate-[-2deg] items-center gap-1.5 border-[3px] border-jv-ink bg-jv-mint px-3 py-1.5 text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink shadow-brutal-sm sm:text-[14px]"
              aria-label="Public quiz"
            >
              <Globe class="size-4" :stroke-width="2.6" />
              Public
            </span>
          </div>
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

      <div
        v-if="canEditPublicMeta && canEditQuiz"
        class="grid gap-5 sm:grid-cols-2 sm:gap-x-8"
      >
        <label class="grid content-start gap-2">
          <span
            class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
          >
            Category
          </span>
          <select
            v-model="settings.category_id"
            :disabled="settingsPending"
            class="h-14 border-[3px] border-jv-ink bg-jv-canvas px-4 text-[16px] font-semibold text-jv-ink outline-none transition-shadow focus:shadow-brutal-sm disabled:opacity-60"
          >
            <option value="">No category (shows under "Other")</option>
            <option
              v-for="category in categories"
              :key="category.id"
              :value="category.id"
            >
              {{ category.name }}
            </option>
          </select>
        </label>

        <div class="grid content-start gap-2">
          <span
            class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
          >
            Cover Image
          </span>
          <label
            v-if="!settings.cover_image"
            class="flex h-14 cursor-pointer items-center justify-center gap-2 border-[3px] border-dashed border-jv-ink/40 bg-jv-canvas px-4 text-[15px] font-semibold text-jv-muted transition-colors hover:bg-jv-yellow/20"
          >
            <ImageIcon class="size-4 shrink-0" :stroke-width="2.3" />
            <span class="truncate">Upload cover image</span>
            <input
              type="file"
              class="hidden"
              accept="image/*"
              :disabled="settingsPending"
              @change="handleCoverImage"
            />
          </label>
          <div
            v-else
            class="flex h-14 items-center gap-3 border-[3px] border-jv-ink bg-jv-canvas px-3"
          >
            <img
              :src="settings.cover_image"
              alt="Cover image preview"
              class="h-10 w-14 shrink-0 border-2 border-jv-ink object-cover"
            />
            <span
              class="min-w-0 flex-1 truncate text-[14px] font-semibold text-jv-ink"
            >
              {{ coverImageName || "Current cover image" }}
            </span>
            <button
              type="button"
              class="grid size-8 shrink-0 place-items-center text-jv-ink transition-colors hover:text-jv-coral"
              aria-label="Remove cover image"
              :disabled="settingsPending"
              @click="removeCoverImage"
            >
              <X class="size-4" :stroke-width="2.4" />
            </button>
          </div>
        </div>
      </div>

      <div class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_auto] lg:items-center">
        <div
          v-if="questions.length > 0"
          class="flex flex-wrap items-center gap-4"
        >
          <TooltipProvider>
            <label
              class="inline-flex items-center gap-2 text-[15px] font-semibold text-jv-muted"
            >
              <span class="inline-flex items-center gap-1.5">
                Points
                <Tooltip>
                  <TooltipTrigger as-child>
                    <button
                      type="button"
                      class="inline-flex size-5 items-center justify-center rounded-full border-[2px] border-jv-ink bg-jv-white text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-[2deg] active:translate-x-[1px] active:translate-y-[1px] active:shadow-none"
                      aria-label="What does Points mean?"
                    >
                      <Info class="size-3" :stroke-width="2.6" />
                    </button>
                  </TooltipTrigger>
                  <TooltipContent side="bottom" align="start">
                    Maximum points awarded for answering this quiz's questions
                    correctly. Players who answer faster earn closer to the full
                    amount; slower correct answers earn proportionally less.
                  </TooltipContent>
                </Tooltip>
              </span>
              <input
                v-model.number="settings.points"
                type="number"
                min="0"
                max="20"
                :disabled="!canEditQuiz || settingsPending"
                class="h-10 w-14 border-[3px] border-jv-ink bg-jv-white px-2 text-center text-[16px] font-black text-jv-ink shadow-brutal-sm outline-none disabled:opacity-60"
              />
            </label>
            <label
              class="inline-flex items-center gap-2 text-[15px] font-semibold text-jv-muted"
            >
              <span class="inline-flex items-center gap-1.5">
                Duration (seconds)
                <Tooltip>
                  <TooltipTrigger as-child>
                    <button
                      type="button"
                      class="inline-flex size-5 items-center justify-center rounded-full border-[2px] border-jv-ink bg-jv-white text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-[2deg] active:translate-x-[1px] active:translate-y-[1px] active:shadow-none"
                      aria-label="What does Duration mean?"
                    >
                      <Info class="size-3" :stroke-width="2.6" />
                    </button>
                  </TooltipTrigger>
                  <TooltipContent side="bottom" align="start">
                    Time in seconds each player has to answer a question before
                    it auto-skips. Shorter durations make the quiz faster-paced
                    but harder.
                  </TooltipContent>
                </Tooltip>
              </span>
              <input
                v-model.number="settings.duration_in_seconds"
                type="number"
                min="1"
                :disabled="!canEditQuiz || settingsPending"
                class="h-10 w-16 border-[3px] border-jv-ink bg-jv-white px-2 text-center text-[16px] font-black text-jv-ink shadow-brutal-sm outline-none disabled:opacity-60"
              />
            </label>
          </TooltipProvider>
        </div>

        <div
          v-if="canEditQuiz"
          class="flex flex-col gap-3 sm:flex-row sm:justify-end"
        >
          <NavigationLink
            v-if="canEditQuiz && hasUnsavedSettings"
            url-name="Save Changes"
            class="rounded-[999px] bg-jv-white"
            :disabled="settingsPending"
            @click="saveSettings"
          />
          <NavigationLink
            url-name="Import Question by .CSV"
            class="rounded-[999px] bg-jv-white"
            @click="importModalOpen = true"
          >
            <Upload class="size-4" :stroke-width="2.4" />
          </NavigationLink>
          <NavigationLink
            url-name="Add Question"
            class="rounded-[999px] bg-jv-accent-green text-white"
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

      <draggable
        v-else
        v-model="orderedQuestions"
        :component-data="{
          tag: 'section',
          name: 'flip-list',
          type: 'transition-group',
        }"
        class="grid gap-5"
        item-key="question_id"
        handle=".drag-handle"
        ghost-class="question-ghost"
        chosen-class="question-chosen"
        drag-class="question-drag"
        :disabled="!canEditQuiz"
        :animation="180"
        :scroll="scrollContainer"
        :bubble-scroll="true"
        :force-autoscroll-fallback="true"
        :scroll-sensitivity="80"
        :scroll-speed="14"
      >
        <template #item="{ element: question, index }">
          <article
            :class="
              editingQuestionId === question.question_id
                ? ''
                : 'jv-border-rough bg-jv-white p-4 shadow-brutal-sm transition-transform hover:rotate-[0.25deg] sm:p-5'
            "
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
                <div class="flex min-w-0 items-start gap-2">
                  <button
                    v-if="canEditQuiz"
                    type="button"
                    class="drag-handle mt-1 grid size-7 shrink-0 cursor-grab place-items-center rounded text-jv-muted hover:bg-jv-canvas hover:text-jv-ink active:cursor-grabbing"
                    aria-label="Drag to reorder"
                    @click.stop
                  >
                    <GripVertical class="size-4" :stroke-width="2.4" />
                  </button>
                  <div class="min-w-0">
                    <p
                      class="text-[12px] font-bold uppercase tracking-[0.14em] text-jv-coral"
                    >
                      Question {{ index + 1 }}
                    </p>
                    <h2
                      class="mt-1 break-words text-[20px] font-bold leading-snug text-jv-ink sm:text-[22px]"
                    >
                      {{ question.question }}
                    </h2>
                  </div>
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
                class="mt-3 flex justify-center rounded-md bg-jv-canvas p-2"
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
                class="mt-3 min-w-0 overflow-x-auto"
                @click.stop
              >
                <CodeBlockComponent :code="question.resource" />
              </div>

              <ul class="mt-4 flex flex-col">
                <li
                  v-for="(option, key) in question.options"
                  :key="key"
                  class="flex min-w-0 items-center gap-3 border-b border-jv-ink/10 py-3 pl-3 pr-2 text-[15px] font-medium text-jv-ink last:border-b-0"
                  :class="
                    isCorrectAnswer(question, key)
                      ? 'border-l-4 border-l-jv-accent-green bg-jv-accent-green/25 pl-2'
                      : 'border-l-4 border-l-transparent'
                  "
                >
                  <span
                    class="w-5 shrink-0 text-[14px] font-bold text-jv-coral"
                  >
                    {{ optionLabel(key) }}.
                  </span>

                  <div
                    v-if="question.options_media === 'image' && option"
                    class="flex min-w-0 flex-1 justify-start"
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

                  <Check
                    v-if="isCorrectAnswer(question, key)"
                    class="size-5 shrink-0 text-jv-accent-green"
                    :stroke-width="3"
                  />
                </li>
              </ul>
            </template>
          </article>
        </template>
      </draggable>
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
            <NavigationLink
              url="/files/demo.csv"
              url-name="Download Sample"
              external
              download="demo.csv"
              class="w-fit rounded-[999px] bg-jv-white"
            >
              <Download class="size-4" :stroke-width="2.4" />
            </NavigationLink>
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
              class="inline-flex h-12 items-center justify-center border-[3px] border-jv-ink bg-jv-accent-green px-6 text-[17px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none disabled:cursor-not-allowed disabled:opacity-60"
              :disabled="importPending"
            >
              {{ importPending ? "Adding..." : "Add Questions" }}
            </button>
          </div>
        </form>
      </div>
    </Teleport>

    <ShareQuizModal v-model="shareModalOpen" :quiz-id="quizId" />
  </main>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import {
  Check,
  Download,
  Globe,
  GripVertical,
  Image as ImageIcon,
  Info,
  Pencil,
  Play,
  Plus,
  Share2,
  Trash2,
  Upload,
  X,
} from "lucide-vue-next";
import { usePush } from "notivue";
import draggable from "vuedraggable";
import QuestionFormCard from "@/components/quiz-manage/QuestionFormCard.vue";
import CodeBlockComponent from "@/components/CodeBlockComponent.vue";
import ShareQuizModal from "@/components/Quiz/ShareQuizModal.vue";
import { useListUserstore } from "~/store/userlist";
import { useSessionStore } from "~~/store/session";
import NavigationLink from "~/components/common/NavigationLink.vue";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

definePageMeta({
  layout: "empty",
});

useSeoMeta({
  title: "Edit Quiz - jovVix",
  description: "Edit quiz details and manage questions for your jovVix quiz.",
  robots: "noindex, nofollow",
});

const toast = usePush();
const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const route = useRoute();
const router = useRouter();
const sessionStore = useSessionStore();
const listUserStore = useListUserstore();

const quizId = computed(() => route.params.quiz_id || "");
const importModalOpen = ref(false);
const shareModalOpen = ref(false);
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

const {
  refresh,
  data: quizData,
  error: quizError,
} = await useFetch(`${url.apiUrl}/quizzes/${quizId.value}/questions`, {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
});

const settings = ref({
  points: Number(quizData.value?.data?.points ?? 10),
  duration_in_seconds: Number(quizData.value?.data?.duration_in_seconds ?? 10),
  category_id: nullableString(quizData.value?.data?.category_id),
  cover_image: nullableString(quizData.value?.data?.cover_image),
});
const coverImageName = ref("");
const { pickCoverImage } = useCoverImage();

const { data: categoriesData } = await useFetch(`${url.apiUrl}/categories`, {
  method: "GET",
  headers: headers,
  credentials: "include",
});
const categories = computed(() => categoriesData.value?.data || []);

const orderedQuestions = ref([...(quizData.value?.data?.data || [])]);
const originalOrderIds = ref(
  (quizData.value?.data?.data || []).map((q) => q.question_id)
);
const scrollContainer = ref(null);

onMounted(() => {
  scrollContainer.value =
    document.querySelector(".lg\\:overflow-y-auto") ||
    document.scrollingElement;
});

const questions = computed(() => quizData.value?.data?.data || []);
const quizTitle = computed(() => quizData.value?.data?.quiz_title || "Quiz");
const quizDescription = computed(
  () => quizData.value?.data?.quiz_description?.String || ""
);
const isPublicQuiz = computed(() => !!quizData.value?.data?.is_public);
// The server decides this: public quiz + the configured admin allowlist.
const canEditPublicMeta = computed(
  () => !!quizData.value?.data?.can_edit_public_meta
);
const totalSurveyQuestion = computed(() =>
  questions.value.reduce(
    (count, item) => (item.question_type_id === 2 ? count + 1 : count),
    0
  )
);
const canEditQuiz = computed(() => {
  const permission = quizData.value?.data?.permission;
  const isEditable = quizData.value?.data?.is_quiz_editable;
  return (permission === "write" || permission === "share") && isEditable;
});

const hasUnsavedSettings = computed(() => {
  const data = quizData.value?.data;
  if (!data) return false;
  const currentIds = orderedQuestions.value.map((q) => q.question_id);
  const orderChanged =
    currentIds.length !== originalOrderIds.value.length ||
    currentIds.some((id, i) => id !== originalOrderIds.value[i]);
  return (
    Number(settings.value.points) !== Number(data.points ?? 10) ||
    Number(settings.value.duration_in_seconds) !==
      Number(data.duration_in_seconds ?? 10) ||
    settings.value.category_id !== nullableString(data.category_id) ||
    settings.value.cover_image !== nullableString(data.cover_image) ||
    orderChanged
  );
});

watch(
  () => quizData.value?.data,
  (data) => {
    if (!data) return;
    settings.value = {
      points: Number(data.points ?? 10),
      duration_in_seconds: Number(data.duration_in_seconds ?? 10),
      category_id: nullableString(data.category_id),
      cover_image: nullableString(data.cover_image),
    };
    // The server value is authoritative after a refresh, and we no longer
    // have the filename that produced it.
    coverImageName.value = "";
    const serverQuestions = data.data || [];
    orderedQuestions.value = [...serverQuestions];
    originalOrderIds.value = serverQuestions.map((q) => q.question_id);
  }
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

const handleCoverImage = async (event) => {
  const picked = await pickCoverImage(event);
  if (!picked) return;
  settings.value.cover_image = picked.dataUrl;
  coverImageName.value = picked.name;
};

const removeCoverImage = () => {
  settings.value.cover_image = "";
  coverImageName.value = "";
};

const saveSettings = async () => {
  if (!canEditQuiz.value) return;

  const body = {
    points: Number(settings.value.points),
    duration_in_seconds: Number(settings.value.duration_in_seconds),
    question_ids: orderedQuestions.value.map((q) => q.question_id),
  };

  // Only send category/cover when they actually changed — omitting them leaves
  // the columns untouched, which keeps a points tweak from re-uploading the
  // whole base64 cover image (and from tripping the admin-only check).
  if (canEditPublicMeta.value) {
    const server = quizData.value?.data || {};
    if (settings.value.category_id !== nullableString(server.category_id)) {
      body.category_id = settings.value.category_id;
    }
    if (settings.value.cover_image !== nullableString(server.cover_image)) {
      body.cover_image = settings.value.cover_image;
    }
  }

  try {
    settingsPending.value = true;
    await $fetch(`${url.apiUrl}/quizzes/${quizId.value}/settings`, {
      method: "PUT",
      headers,
      credentials: "include",
      body,
    });
    toast.success("Quiz settings updated.");
    refresh();
  } catch (error) {
    toast.error(
      error?.data?.message ||
        error?.message ||
        "Failed to update quiz settings."
    );
  } finally {
    settingsPending.value = false;
  }
};

const saveNewQuestion = async ({ payload }) => {
  try {
    savingNewQuestion.value = true;
    const payloadWithDefaults = {
      ...payload,
      points: Number(settings.value.points),
      duration_in_seconds: Number(settings.value.duration_in_seconds),
    };
    await $fetch(`${url.apiUrl}/quizzes/${quizId.value}/questions`, {
      method: "POST",
      headers,
      credentials: "include",
      body: payloadWithDefaults,
    });

    toast.success("Question added successfully.");
    showNewQuestionForm.value = false;
    showAddAnother.value = true;
    refresh();
  } catch (error) {
    toast.error(
      error?.data?.message || error?.message || "Failed to add question."
    );
  } finally {
    savingNewQuestion.value = false;
  }
};

const saveExistingQuestion = async (question, { payload }) => {
  const questionId = question.question_id;

  try {
    savingQuestionId.value = questionId;
    await $fetch(
      `${url.apiUrl}/quizzes/${quizId.value}/questions/${questionId}`,
      {
        method: "PUT",
        headers,
        credentials: "include",
        body: {
          ...payload,
          points: Number(settings.value.points),
          duration_in_seconds: Number(settings.value.duration_in_seconds),
        },
      }
    );

    toast.success("Question updated successfully.");
    editingQuestionId.value = "";
    refresh();
  } catch (error) {
    toast.error(
      error?.data?.message || error?.message || "Failed to update question."
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
      }
    );
    toast.success("Question deleted successfully.");
    refresh();
  } catch (error) {
    toast.error(
      error?.data?.message || error?.message || "Failed to delete question."
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
      error?.data?.message || error?.message || "Failed to delete quiz."
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
        "Failed to import CSV."
    );
  } finally {
    importPending.value = false;
  }
};

const handleShareQuiz = () => {
  shareModalOpen.value = true;
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
      }
    );

    const activeQuizId = response?.data;
    if (!activeQuizId) {
      toast.error("Error while starting quiz.");
      return;
    }

    if (quizTitle.value) {
      sessionStore.setActiveQuizTitle(quizTitle.value);
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

<style scoped>
.flip-list-move {
  transition: transform 0.4s ease;
}

/* Drop indicator: faded preview of the card in its drop position. */
.question-ghost {
  position: relative;
  opacity: 0.4 !important;
  background-color: var(--jv-white) !important;
  border: 3px dashed var(--jv-coral) !important;
  box-shadow: none !important;
  transform: none !important;
  transition: opacity 0.15s ease, border-color 0.15s ease;
}

/* Disable the hover tilt on the ghost while it's acting as drop placeholder. */
.question-ghost:hover {
  transform: none !important;
}

/* Floating clone: the faded preview that follows the cursor. */
.question-drag {
  opacity: 0.92 !important;
  transform: rotate(1.5deg) !important;
  box-shadow: 8px 8px 0 0 var(--jv-ink) !important;
  cursor: grabbing !important;
}

/* Origin item: subtly dimmed so users see where it came from. */
.question-chosen {
  cursor: grabbing;
}
</style>
