<script setup>
import { computed, ref, watchEffect } from "vue";
import { ChevronDown, Filter, Search } from "lucide-vue-next";
import { useToast } from "vue-toastification";
import AdminQuizListCard from "@/components/quiz-list/AdminQuizListCard.vue";
import usecopyToClipboard from "@/composables/copy_to_clipboard";
import { useListUserstore } from "~/store/userlist";
import { useSessionStore } from "~~/store/session";
import NavigationLink from "@/components/common/NavigationLink.vue";

definePageMeta({
  layout: "empty",
});

const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const toast = useToast();
const router = useRouter();
const sessionStore = useSessionStore();
const listUserStore = useListUserstore();

const searchQuery = ref("");
const startingQuizId = ref("");
const selectedFilter = ref("All Quiz");
const filterOpen = ref(false);
const filterOptions = ["All Quiz", "Shared By Me", "Shared With Me"];

const quizImages = [
  "/images/landing/homepage-public-quiz-1.png",
  "/images/landing/homepage-public-quiz-2.png",
  "/images/landing/homepage-public-quiz-3.png",
  "/images/landing/homepage-public-quiz-4.png",
];

const tiltClasses = [
  "rotate-[-0.8deg]",
  "rotate-[0.7deg]",
  "rotate-[-0.4deg]",
  "rotate-[0.5deg]",
];

const {
  data: quizList,
  pending: quizPending,
  error: quizError,
  refresh,
} = useFetch(url.apiUrl + "/quizzes", {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
});

watchEffect(() => {
  if (quizError.value?.data?.code === 401) {
    navigateTo("/account/login");
  }
});

const getTitle = (quiz) => {
  try {
    return decodeURI(quiz?.title || "Untitled Quiz");
  } catch {
    return quiz?.title || "Untitled Quiz";
  }
};

const getDescription = (quiz) => {
  const description = quiz?.description?.String || quiz?.description || "";
  return description || "No description added yet.";
};

const formatCreatedAt = (value) => {
  if (!value) return "Created recently";

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "Created recently";

  return `Created ${new Intl.DateTimeFormat("en", {
    month: "short",
    day: "2-digit",
    year: "numeric",
  }).format(date)}`;
};

const quizzes = computed(() =>
  (quizList.value?.data || []).map((quiz, index) => ({
    id: quiz.id,
    title: getTitle(quiz),
    description: getDescription(quiz),
    createdAt: formatCreatedAt(quiz.created_at),
    questionCount: quiz.total_questions || 0,
    image: quizImages[index % quizImages.length],
    tiltClass: tiltClasses[index % tiltClasses.length],
    viewUrl: `/admin/quiz/list-quiz/${quiz.id}`,
  }))
);

const filteredQuizzes = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  if (!query) return quizzes.value;

  return quizzes.value.filter((quiz) => {
    return (
      quiz.title.toLowerCase().includes(query) ||
      quiz.description.toLowerCase().includes(query)
    );
  });
});

const handleStartQuiz = async (quizId) => {
  try {
    startingQuizId.value = quizId;
    const response = await $fetch(
      `${url.apiUrl}/quizzes/${quizId}/demo_session`,
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

    listUserStore.removeAllUsers();
    setSocketObject(null);
    router.push(`/admin/arrange/${activeQuizId}`);

    setTimeout(() => {
      sessionStore.setSession(activeQuizId);
    }, 1000);
  } catch (error) {
    toast.error(error?.message || "Error while starting quiz.");
  } finally {
    startingQuizId.value = "";
  }
};

const handleShareQuiz = (quiz) => {
  if (!import.meta.client) return;

  usecopyToClipboard(`${window.location.origin}${quiz.viewUrl}`);
};

const handleDeleteQuiz = async (quizId) => {
  try {
    await $fetch(`${url.apiUrl}/quizzes/${quizId}`, {
      method: "DELETE",
      headers: headers,
      credentials: "include",
    });
    toast.success("Quiz deleted successfully.");
    refresh();
  } catch (error) {
    console.error("Failed to delete quiz", error);
    toast.error("Failed to delete quiz.");
  }
};

const handleSelectFilter = (option) => {
  selectedFilter.value = option;
  filterOpen.value = false;
};
</script>

<template>
  <main
    class="min-h-screen bg-jv-canvas px-4 py-6 sm:px-6 md:px-8 lg:px-10 xl:px-12"
  >
    <div class="mx-auto max-w-[1180px]">
      <div
        class="mb-8 flex flex-col gap-5 sm:mb-10 md:flex-row md:items-start md:justify-between"
      >
        <div>
          <h1 class="font-headings text-jv-ink sm:text-[52px]">Quiz List</h1>
          <div
            class="ml-2 mt-1 h-3 w-28 rounded-full border-b-[3px] border-jv-yellow"
            aria-hidden="true"
          ></div>
        </div>
        <NavigationLink
          url="/admin/quiz/create-quiz"
          url-name="Create Quiz"
          class="bg-jv-coral text-white font-[500] py-2"
        />
      </div>

      <div
        class="mb-8 grid gap-4 md:grid-cols-[220px_minmax(280px,448px)] md:items-center md:justify-between"
      >
        <div class="relative w-fit">
          <button
            type="button"
            class="inline-flex h-12 w-fit rotate-[-1deg] items-center gap-2 jv-border-rough bg-jv-white px-4 text-[18px] font-semibold text-jv-ink shadow-brutal-sm"
            :aria-expanded="filterOpen"
            aria-haspopup="listbox"
            @click="filterOpen = !filterOpen"
          >
            <Filter class="size-5 text-jv-coral" :stroke-width="2.2" />
            <span>{{ selectedFilter }}</span>
            <ChevronDown
              class="size-4 text-jv-ink/60 transition-transform"
              :class="filterOpen ? 'rotate-180' : ''"
              :stroke-width="2.4"
            />
          </button>

          <div
            v-if="filterOpen"
            class="absolute left-0 top-14 z-30 w-48 rotate-[1deg] jv-border-rough bg-jv-white p-2 shadow-brutal-sm"
            role="listbox"
          >
            <button
              v-for="option in filterOptions"
              :key="option"
              type="button"
              class="block w-full rounded-[6px] px-3 py-2 text-left text-[15px] font-semibold text-jv-ink transition-colors hover:bg-jv-yellow/40"
              :class="option === selectedFilter ? 'bg-jv-yellow/60' : ''"
              role="option"
              :aria-selected="option === selectedFilter"
              @click="handleSelectFilter(option)"
            >
              {{ option }}
            </button>
          </div>
        </div>

        <label
          class="flex h-12 rotate-[1deg] items-center gap-3 jv-border-rough bg-jv-white px-4 shadow-brutal-sm"
        >
          <Search class="size-5 shrink-0 text-jv-ink/40" :stroke-width="2.4" />
          <input
            v-model="searchQuery"
            type="search"
            class="h-full min-w-0 flex-1 bg-transparent text-[17px] text-jv-ink outline-none placeholder:text-jv-ink/45"
            placeholder="Search by name or keyword..."
          />
        </label>
      </div>

      <UtilsQuizListWaiting v-if="quizPending" />

      <div
        v-else-if="quizError"
        class="jv-border-rough bg-jv-white p-5 text-[18px] font-semibold text-jv-coral shadow-brutal-sm"
      >
        {{ quizError.message }}
      </div>

      <section
        v-else-if="quizzes.length < 1"
        class="grid min-h-[420px] place-items-center text-center"
      >
        <div class="max-w-md">
          <h2 class="font-headings text-[34px] text-jv-ink">
            No Quiz Created By You!
          </h2>
          <p class="mt-2 text-[17px] text-jv-muted">Create your first quiz.</p>
          <NavigationLink
            url="/admin/quiz/create-quiz"
            url-name="Create Quiz"
            class="bg-jv-coral text-white font-[500] py-2"
          />
        </div>
      </section>

      <section
        v-else-if="filteredQuizzes.length < 1"
        class="jv-border-rough bg-jv-white p-5 text-[18px] font-semibold text-jv-muted shadow-brutal-sm"
      >
        No quizzes matched your search.
      </section>

      <section
        v-else
        class="grid gap-x-5 gap-y-7 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4"
      >
        <AdminQuizListCard
          v-for="quiz in filteredQuizzes"
          :key="quiz.id"
          :title="quiz.title"
          :description="quiz.description"
          :created-at="quiz.createdAt"
          :question-count="quiz.questionCount"
          :image="quiz.image"
          :tilt-class="quiz.tiltClass"
          :view-url="quiz.viewUrl"
          :starting="startingQuizId === quiz.id"
          @start-quiz="handleStartQuiz(quiz.id)"
          @share="handleShareQuiz(quiz)"
          @delete="handleDeleteQuiz(quiz.id)"
        />
      </section>
    </div>
  </main>
</template>
