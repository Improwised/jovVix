<template>
  <section
    v-if="publicQuizzes.length > 0"
    id="public-quiz"
    class="mt-0 sm:mt-1 scroll-mt-6"
  >
    <div class="mb-4 sm:mb-5 flex items-center justify-between gap-3 sm:gap-4">
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
      v-if="groupedQuizzes.length > 1"
      class="mb-5 flex flex-wrap gap-2 sm:mb-6 sm:gap-3"
    >
      <NavigationLink
        v-for="group in groupedQuizzes"
        :key="`nav-${group.name}`"
        type="button"
        :url-name="group.name"
        class="bg-jv-white rounded-full px-3 py-1.5 text-[14px] sm:px-4 sm:text-[15px] md:py-1.5 md:text-[15px]"
        @click="scrollToCategory(group.name)"
      />
    </div>

    <div class="flex flex-col gap-7 sm:gap-9">
      <div
        v-for="group in groupedQuizzes"
        :id="categoryAnchorId(group.name)"
        :key="group.name"
        class="scroll-mt-6"
      >
        <h4
          class="mb-1 font-headings font-bold text-[19px] leading-[1.2] text-jv-ink sm:text-[22px] md:text-[24px]"
        >
          {{ group.name }}
        </h4>
        <Carousel
          :opts="{ align: 'start', loop: false, slidesToScroll: 'auto' }"
          class="group"
        >
          <CarouselContent class="-ml-6 pb-3 pt-4 sm:-ml-8">
            <CarouselItem
              v-for="(quiz, index) in group.quizzes"
              :key="quiz.id"
              class="basis-1/2 pl-6 sm:basis-1/3 sm:pl-8 md:basis-1/4 lg:basis-1/5 xl:basis-1/6"
            >
              <PublicQuizCard
                :title="quiz.title"
                :image="coverOf(quiz, index)"
                :tilt-class="tiltFor(index)"
                :starting="startingQuizId === quiz.id"
                @start-quiz="handleStartQuiz(quiz.id)"
              />
            </CarouselItem>
          </CarouselContent>
          <CarouselPrevious class="jv-carousel-nav !left-2" />
          <CarouselNext class="jv-carousel-nav !right-2" />
        </Carousel>
      </div>
    </div>
  </section>
</template>

<script setup>
import { usePush } from "notivue";
import { computed, ref } from "vue";
import { Compass } from "lucide-vue-next";
import PublicQuizCard from "@/components/landing/PublicQuizCard.vue";
import NavigationLink from "@/components/common/NavigationLink.vue";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";
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

const tiltFor = (i) => tiltClasses[i % tiltClasses.length];

// The API serializes sql.NullString as { String, Valid } when populated; null otherwise.
const nullableString = (value) => {
  if (!value) return "";
  if (typeof value === "string") return value;
  return value.String || "";
};

const coverOf = (quiz, index) =>
  nullableString(quiz.cover_image) ||
  fallbackImages[index % fallbackImages.length];

const groupedQuizzes = computed(() => {
  const groups = new Map();
  for (const quiz of publicQuizzes.value) {
    const name = nullableString(quiz.category_name) || "Other";
    if (!groups.has(name)) groups.set(name, []);
    groups.get(name).push(quiz);
  }
  // Alphabetical, with the uncategorized "Other" group always last.
  return [...groups.entries()]
    .sort(([a], [b]) => {
      if (a === "Other") return 1;
      if (b === "Other") return -1;
      return a.localeCompare(b);
    })
    .map(([name, quizzes]) => ({ name, quizzes }));
});

const categoryAnchorId = (name) =>
  `public-quiz-category-${name.toLowerCase().replace(/[^a-z0-9]+/g, "-")}`;

const scrollToCategory = (name) => {
  document
    .getElementById(categoryAnchorId(name))
    ?.scrollIntoView({ behavior: "smooth", block: "start" });
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
