<script setup>
import { computed, ref, watchEffect } from "vue";
import { ChevronDown, Filter, Search, X } from "lucide-vue-next";
import { usePush } from "notivue";
import AdminQuizListCard from "@/components/quiz-list/AdminQuizListCard.vue";
import ShareQuizModal from "@/components/Quiz/ShareQuizModal.vue";
import { useListUserstore } from "~/store/userlist";
import { useSessionStore } from "~~/store/session";
import { useUsersStore } from "~~/store/users";
import NavigationLink from "@/components/common/NavigationLink.vue";

definePageMeta({
  layout: "empty",
});

const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const toast = usePush();
const router = useRouter();
const route = useRoute();
const sessionStore = useSessionStore();
const listUserStore = useListUserstore();
const usersStore = useUsersStore();

const canCreatePublicQuiz = computed(
  () => !!usersStore.userData?.canCreatePublicQuiz
);

const searchQuery = ref("");
const startingQuizId = ref("");
const shareModalOpen = ref(false);
const shareQuizId = ref("");
const createQuizOpen = ref(false);
const createQuizPending = ref(false);
const createQuizForm = ref({
  title: "",
  description: "",
  is_public: false,
});
const selectedFilter = ref("All Quiz");
const filterOpen = ref(false);
const filterOptions = ["All Quiz", "Shared By Me", "Shared With Me"];
const filterApiPaths = {
  "All Quiz": "/quizzes",
  "Shared By Me": "/shared_quizzes?type=shared_by_me",
  "Shared With Me": "/shared_quizzes?type=shared_with_me",
};

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
const quizListUrl = computed(
  () => `${url.apiUrl}${filterApiPaths[selectedFilter.value]}`
);
const emptyState = computed(() => {
  if (selectedFilter.value === "Shared By Me") {
    return {
      title: "No Quiz Shared By You!",
      message: "Share your first quiz.",
      showCreateAction: false,
    };
  }

  if (selectedFilter.value === "Shared With Me") {
    return {
      title: "No Quiz Is Shared With You!",
      message: "Ask someone to share a quiz with you.",
      showCreateAction: false,
    };
  }

  return {
    title: "No Quiz Created By You!",
    message: "Create your first quiz.",
    showCreateAction: true,
  };
});

const {
  data: quizList,
  pending: quizPending,
  error: quizError,
  refresh,
} = useFetch(quizListUrl, {
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

watchEffect(() => {
  if (route.query.create === "1") {
    createQuizOpen.value = true;
  }
});

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
    title: quiz.title,
    description: quiz.description.String,
    createdAt: formatCreatedAt(quiz.created_at),
    questionCount: quiz.total_questions || 0,
    image: quizImages[index % quizImages.length],
    tiltClass: tiltClasses[index % tiltClasses.length],
    viewUrl: `/admin/quiz/list-quiz/${quiz.id}`,
    isPublic: !!quiz.is_public,
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

    const startedQuiz = quizzes.value.find((q) => q.id === quizId);
    if (startedQuiz?.title) {
      sessionStore.setActiveQuizTitle(startedQuiz.title);
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
  shareQuizId.value = quiz.id;
  shareModalOpen.value = true;
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
    toast.error(
      error?.data?.message || error?.message || "Failed to delete quiz."
    );
  }
};

const handleSelectFilter = (option) => {
  selectedFilter.value = option;
  filterOpen.value = false;
};

const closeCreateQuizModal = () => {
  createQuizOpen.value = false;
  createQuizForm.value = {
    title: "",
    description: "",
    is_public: false,
  };
  if (route.query.create === "1") {
    router.replace({ path: route.path, query: {} });
  }
};

const handleCreateQuiz = async () => {
  try {
    createQuizPending.value = true;
    const response = await $fetch(`${url.apiUrl}/quizzes`, {
      method: "POST",
      headers: {
        Accept: "application/json",
      },
      body: createQuizForm.value,
      credentials: "include",
    });

    const quizId = response?.data;
    if (!quizId) {
      toast.error("Error while creating quiz.");
      return;
    }

    toast.success("Quiz created successfully.");
    closeCreateQuizModal();
    router.push(`/admin/quiz/list-quiz/${quizId}`);
  } catch (error) {
    toast.error(
      error?.data?.message || error?.message || "Error while creating quiz."
    );
  } finally {
    createQuizPending.value = false;
  }
};
</script>

<template>
  <main
    class="flex min-h-screen flex-col gap-8 bg-jv-canvas px-4 py-5 sm:gap-10 sm:px-6 md:px-8 md:py-6"
  >
    <div
      class="flex flex-col gap-5 md:flex-row md:items-start md:justify-between"
    >
      <div class="min-w-0">
        <h1
          class="font-headings text-[38px] leading-none text-jv-ink min-[420px]:text-[44px] sm:text-[52px] md:text-[56px]"
        >
          Quiz List
        </h1>
        <div
          class="ml-1 mt-1 h-3 w-24 rounded-full border-b-[3px] border-jv-yellow sm:ml-2 sm:w-28"
          aria-hidden="true"
        ></div>
      </div>
      <NavigationLink
        url-name="Create Quiz"
        class="w-full bg-jv-coral py-2 font-[500] text-white sm:w-fit md:mt-1"
        @click="createQuizOpen = true"
      />
    </div>

    <div
      class="grid gap-3 sm:gap-4 md:grid-cols-[minmax(0,220px)_minmax(280px,448px)] md:items-center md:justify-between"
    >
      <div class="relative w-full sm:w-fit">
        <button
          type="button"
          class="inline-flex h-11 w-full rotate-[-1deg] items-center justify-between gap-2 jv-border-rough bg-jv-white px-3 text-[16px] font-semibold text-jv-ink shadow-brutal-sm sm:h-12 sm:w-fit sm:justify-start sm:px-4 sm:text-[18px]"
          :aria-expanded="filterOpen"
          aria-haspopup="listbox"
          @click="filterOpen = !filterOpen"
        >
          <span class="flex min-w-0 items-center gap-2">
            <Filter class="size-5 shrink-0 text-jv-coral" :stroke-width="2.2" />
            <span class="truncate">{{ selectedFilter }}</span>
          </span>
          <ChevronDown
            class="size-4 shrink-0 text-jv-ink/60 transition-transform"
            :class="filterOpen ? 'rotate-180' : ''"
            :stroke-width="2.4"
          />
        </button>

        <div
          v-if="filterOpen"
          class="absolute left-0 top-[52px] z-30 w-full min-w-[190px] rotate-[1deg] jv-border-rough bg-jv-white p-2 shadow-brutal-sm sm:top-14 sm:w-52"
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
        class="flex h-11 min-w-0 rotate-[1deg] items-center gap-2.5 jv-border-rough bg-jv-white px-3 shadow-brutal-sm sm:h-12 sm:gap-3 sm:px-4"
      >
        <Search class="size-5 shrink-0 text-jv-ink/40" :stroke-width="2.4" />
        <input
          v-model="searchQuery"
          type="search"
          class="h-full min-w-0 flex-1 bg-transparent text-[15px] text-jv-ink outline-none placeholder:text-jv-ink/45 sm:text-[17px]"
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
      class="grid min-h-[320px] place-items-center px-2 text-center sm:min-h-[420px]"
    >
      <div class="max-w-md py-8">
        <h2
          class="font-headings text-[28px] leading-tight text-jv-ink sm:text-[34px]"
        >
          {{ emptyState.title }}
        </h2>
        <p class="mt-2 text-[17px] text-jv-muted">
          {{ emptyState.message }}
        </p>
        <NavigationLink
          v-if="emptyState.showCreateAction"
          url-name="Create Quiz"
          class="bg-jv-coral text-white font-[500] py-2"
          @click="createQuizOpen = true"
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
      class="grid gap-5 sm:grid-cols-2 sm:gap-6 lg:grid-cols-3 2xl:grid-cols-4"
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
        :is-public="quiz.isPublic"
        :starting="startingQuizId === quiz.id"
        @start-quiz="handleStartQuiz(quiz.id)"
        @share="handleShareQuiz(quiz)"
        @delete="handleDeleteQuiz(quiz.id)"
      />
    </section>

    <Teleport to="body">
      <div
        v-if="createQuizOpen"
        class="fixed inset-0 z-50 grid place-items-center bg-jv-ink/35 px-4 py-6 backdrop-blur-[2px]"
        @click.self="closeCreateQuizModal"
      >
        <form
          class="w-full max-w-[680px] rotate-[-0.4deg] border-[4px] border-jv-ink bg-jv-white shadow-brutal-lg"
          @submit.prevent="handleCreateQuiz"
        >
          <div
            class="flex items-center justify-between gap-4 border-b-[3px] border-jv-ink bg-jv-ink px-5 py-4 text-jv-white sm:px-6"
          >
            <h2
              class="font-body text-[24px] font-black leading-none text-jv-white sm:text-[28px]"
            >
              Create New Quiz
            </h2>
            <button
              type="button"
              class="grid size-9 place-items-center text-jv-white transition-transform hover:rotate-[6deg]"
              aria-label="Close create quiz modal"
              @click="closeCreateQuizModal"
            >
              <X class="size-6" :stroke-width="2.4" />
            </button>
          </div>

          <div class="grid gap-5 px-5 py-6 sm:px-8">
            <label class="grid gap-2">
              <span
                class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
              >
                Quiz Title <span class="text-jv-coral">*</span>
              </span>
              <input
                v-model.trim="createQuizForm.title"
                type="text"
                required
                maxlength="50"
                class="h-14 border-[3px] border-jv-ink bg-jv-canvas px-4 text-[17px] font-semibold text-jv-ink outline-none transition-shadow focus:shadow-brutal-sm"
              />
            </label>

            <label class="grid gap-2">
              <span
                class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
              >
                Quiz Description
              </span>
              <textarea
                v-model.trim="createQuizForm.description"
                maxlength="150"
                rows="5"
                class="resize-none border-[3px] border-jv-ink bg-jv-canvas px-4 py-3 text-[17px] font-semibold text-jv-ink outline-none transition-shadow focus:shadow-brutal-sm"
              ></textarea>
            </label>

            <label
              v-if="canCreatePublicQuiz"
              class="flex cursor-pointer items-start gap-3 jv-border-rough bg-jv-canvas p-4"
            >
              <input
                v-model="createQuizForm.is_public"
                type="checkbox"
                class="mt-1 size-5 shrink-0 cursor-pointer accent-jv-coral"
              />
              <span class="flex flex-col gap-1">
                <span
                  class="text-[14px] font-black uppercase tracking-[0.14em] text-jv-ink"
                >
                  Make this quiz public
                </span>
                <span class="text-[14px] font-medium text-jv-muted">
                  Public quizzes appear under "Explore Public Quizzes" on the
                  homepage. Anyone can discover them.
                </span>
              </span>
            </label>
          </div>

          <div
            class="flex flex-col-reverse gap-3 border-t-[3px] border-jv-ink bg-jv-canvas px-5 py-4 sm:flex-row sm:justify-end sm:px-8"
          >
            <NavigationLink
              url-name="Cancel"
              class="bg-jv-white font-[500]"
              @click="closeCreateQuizModal"
            />
            <NavigationLink
              url-name="Create Quiz"
              class="bg-jv-mint font-[500]"
            />
          </div>
        </form>
      </div>
    </Teleport>

    <ShareQuizModal v-model="shareModalOpen" :quiz-id="shareQuizId" />
  </main>
</template>
