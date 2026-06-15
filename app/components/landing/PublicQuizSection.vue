<template>
  <section v-if="publicQuizzes.length > 0" class="mt-2 sm:mt-4">
    <div class="mb-5 sm:mb-6 flex items-center justify-between gap-3 sm:gap-4">
      <h3
        class="flex items-center gap-2 font-headings text-[22px] min-[420px]:text-[26px] sm:text-[30px] md:text-[34px] leading-[1.2] text-jv-ink"
      >
        <span
          class="grid size-6 sm:size-7 shrink-0 place-items-center rounded-full border-[3px] border-jv-coral text-jv-coral"
        >
          <Compass class="size-3.5 sm:size-4" :stroke-width="2.6" />
        </span>
        <span>Explore Public Quizzes</span>
      </h3>
    </div>

    <div
      class="grid gap-5 sm:grid-cols-2 sm:gap-6 lg:grid-cols-3 2xl:grid-cols-4"
    >
      <AdminQuizListCard
        v-for="(quiz, index) in publicQuizzes"
        :key="quiz.id"
        :title="quiz.title"
        :description="descriptionOf(quiz)"
        :created-at="formatDate(quiz.created_at)"
        :question-count="quiz.total_questions"
        :image="imageFor(index)"
        :tilt-class="tiltFor(index)"
        :view-url="`/admin/quiz/list-quiz/${quiz.id}`"
        :is-public="true"
        :show-actions="false"
        :starting="startingQuizId === quiz.id"
        @start-quiz="handleStartQuiz(quiz.id)"
      />
    </div>
  </section>
</template>

<script setup>
import { usePush } from "notivue";
import { computed, ref } from "vue";
import { Compass } from "lucide-vue-next";
import AdminQuizListCard from "@/components/quiz-list/AdminQuizListCard.vue";
import { useListUserstore } from "~/store/userlist";
import { useSessionStore } from "~~/store/session";
import { useUsersStore } from "~~/store/users";
const toast = usePush();
const router = useRouter();
const { apiUrl } = useRuntimeConfig().public;
const listUserStore = useListUserstore();
const sessionStore = useSessionStore();
const usersStore = useUsersStore();
const startingQuizId = ref("");
const { data } = await useFetch(`${apiUrl}/quizzes/public`, {
  method: "GET",
  credentials: "include",
});

const publicQuizzes = computed(() => data.value?.data || []);

const fallbackImages = [
  "/images/landing/homepage-public-quiz-1.png",
  "/images/landing/homepage-public-quiz-2.png",
  "/images/landing/homepage-public-quiz-3.png",
  "/images/landing/homepage-public-quiz-4.png",
];

const tiltClasses = [
  "rotate-[-1deg]",
  "rotate-[1deg]",
  "rotate-[-0.4deg]",
  "rotate-[0.6deg]",
];

const imageFor = (i) => fallbackImages[i % fallbackImages.length];
const tiltFor = (i) => tiltClasses[i % tiltClasses.length];

const descriptionOf = (quiz) => {
  // The API serializes sql.NullString as { String, Valid } when populated; null otherwise.
  const d = quiz.description;
  if (!d) return "";
  if (typeof d === "string") return d;
  return d.String || "";
};

const formatDate = (iso) => {
  if (!iso) return "";
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) return "";
  return date.toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
};

// Logged-in visitors host immediately. Guests are routed to the host lobby with
// a sentinel session_id ("new") so the lobby can show the name-entry modal
// before creating their guest user and the public session.
const handleStartQuiz = async (quizId) => {
  if (startingQuizId.value) return;
  startingQuizId.value = quizId;
  const startedQuiz = publicQuizzes.value.find((q) => q.id === quizId);
  if (startedQuiz?.title) {
    sessionStore.setActiveQuizTitle(startedQuiz.title);
  }
  try {
    let isLoggedIn = false;
    try {
      const who = await $fetch(`${apiUrl}/user/who`, {
        method: "GET",
        credentials: "include",
      });
      if (who?.data) {
        usersStore.setUserData({
          role: who.data.role,
          avatar: who.data.avatar,
          firstname: who.data.firstname,
          username: who.data.username,
          email: who.data.email,
          canCreatePublicQuiz: !!who.data.can_create_public_quiz,
        });
        isLoggedIn = who.data.role && who.data.role !== "guest-user";
      }
    } catch (error) {
      // Not authenticated; treat as guest and let the host page collect a name.
    }

    if (!isLoggedIn) {
      await router.push(
        `/admin/arrange/new?public=1&quiz_id=${encodeURIComponent(quizId)}`
      );
      return;
    }

    const response = await $fetch(
      `${apiUrl}/quizzes/${quizId}/public_session`,
      {
        method: "POST",
        credentials: "include",
        headers: {
          Accept: "application/json",
        },
      }
    );

    const sessionId = response?.data;
    if (!sessionId) {
      toast.error("Error while starting quiz.");
      return;
    }

    listUserStore.removeAllUsers();
    setSocketObject(null);
    sessionStore.setSession(sessionId);
    // `public=1` tells the lobby this is a public session where the host may also play.
    await router.push(`/admin/arrange/${sessionId}?public=1`);
  } catch (error) {
    toast.error(
      error?.data?.message || error?.message || "Error while starting quiz."
    );
  } finally {
    startingQuizId.value = "";
  }
};
</script>
