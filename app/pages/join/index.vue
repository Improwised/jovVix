<template>
  <div
    class="min-h-screen bg-jv-canvas relative flex items-center justify-center px-4 sm:px-6 py-6 sm:py-8 overflow-hidden"
  >
    <!-- Decorative background grid -->
    <div
      class="absolute inset-0 jv-grid pointer-events-none"
      aria-hidden="true"
    ></div>

    <!-- Decorative sparkles & squiggle -->
    <div
      class="absolute hidden sm:block left-[10%] top-[18%] rotate-12 text-jv-coral pointer-events-none"
      aria-hidden="true"
    >
      <Sparkle class="size-7 md:size-9 fill-current" :stroke-width="2" />
    </div>
    <div
      class="absolute hidden md:block right-[10%] top-[14%] -rotate-12 text-jv-yellow pointer-events-none"
      aria-hidden="true"
    >
      <svg
        width="120"
        height="48"
        viewBox="0 0 120 48"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          d="M2 28 Q 22 4, 42 22 T 82 26 T 118 14"
          stroke="currentColor"
          stroke-width="3.5"
          stroke-linecap="round"
          fill="none"
        />
      </svg>
    </div>
    <div
      class="absolute hidden md:block right-[14%] bottom-[18%] -rotate-12 text-jv-mint pointer-events-none"
      aria-hidden="true"
    >
      <Sparkle class="size-6 fill-current" :stroke-width="2" />
    </div>

    <div
      class="relative w-full max-w-[440px] mx-auto z-10 flex flex-col items-center"
    >
      <!-- Card -->
      <div class="relative rotate-1 w-full">
        <!-- "clip" tab -->
        <span
          class="absolute left-1/2 -top-[8px] z-20 h-3 w-12 -translate-x-1/2 jv-card border-2 border-jv-ink bg-jv-slate shadow-brutal-sm"
          aria-hidden="true"
        ></span>

        <div
          class="bg-jv-white border-2 border-jv-ink shadow-brutal-lg jv-card px-6 sm:px-8 py-7 sm:py-9"
        >
          <!-- Heading -->
          <header class="flex flex-col items-center gap-1.5 mb-6">
            <div class="relative inline-block">
              <h1
                class="font-headings text-jv-ink text-[32px] sm:text-[38px] leading-none m-0"
              >
                Join Quiz
              </h1>
              <svg
                class="absolute -bottom-2 left-1/2 -translate-x-1/2 text-jv-yellow"
                width="120"
                height="14"
                viewBox="0 0 120 14"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
                aria-hidden="true"
              >
                <path
                  d="M3 9 Q 20 1, 40 7 T 78 6 T 117 4"
                  stroke="currentColor"
                  stroke-width="2.5"
                  stroke-linecap="round"
                  fill="none"
                />
              </svg>
            </div>
            <p class="font-body text-jv-ink/70 text-sm sm:text-base m-0 mt-1">
              Enter the code to join
            </p>
          </header>

          <!-- Form -->
          <form
            method="POST"
            class="flex flex-col gap-4"
            @submit.prevent="join_quiz"
          >
            <!-- Invitation code -->
            <div class="flex flex-col gap-1.5">
              <label
                for="code"
                class="font-body text-jv-ink text-xs sm:text-[13px] font-bold tracking-wide uppercase px-0.5"
              >
                Invitation Code
              </label>
              <div
                class="flex items-center gap-2.5 bg-jv-white border-2 border-jv-ink jv-card shadow-brutal-sm focus-within:translate-x-[1px] focus-within:translate-y-[1px] focus-within:shadow-none transition-all px-3 py-2.5"
              >
                <Hash
                  class="size-[18px] text-jv-ink/70 shrink-0"
                  :stroke-width="2.2"
                />
                <input
                  id="code"
                  v-model="codeDisplay"
                  type="text"
                  inputmode="numeric"
                  autocomplete="off"
                  maxlength="7"
                  placeholder="000 000"
                  class="flex-1 min-w-0 bg-transparent border-0 outline-none font-body font-bold text-jv-ink placeholder:text-jv-ink/30 text-base sm:text-lg tracking-[0.15em]"
                />
                <span
                  v-if="isCodeValid"
                  class="grid place-items-center size-7 rounded-full bg-jv-mint text-jv-accent-green border-2 border-jv-ink/15 shrink-0"
                  aria-label="Valid code"
                >
                  <Check class="size-4" :stroke-width="3" />
                </span>
              </div>
            </div>

            <!-- Player profile (guest) -->
            <div v-if="!isUserLoggedIn" class="flex flex-col gap-1.5">
              <label
                for="username"
                class="font-body text-jv-ink text-xs sm:text-[13px] font-bold tracking-wide uppercase px-0.5"
              >
                Your Player Profile
              </label>
              <div class="flex items-stretch gap-2.5">
                <!-- Avatar with reroll -->
                <NavigationLink
                  type="button"
                  class="relative size-[46px] shrink-0 !p-0 bg-jv-mint !rounded-none !jv-card overflow-hidden shadow-none hover:rotate-0 hover:scale-110"
                  :aria-label="`Generate new avatar (current: ${avatarName})`"
                  @click="rerollAvatar"
                >
                  <img
                    :src="avatarUrl"
                    :alt="avatarName"
                    class="absolute inset-0 size-full object-cover"
                  />
                  <span
                    class="absolute -bottom-1 -right-1 grid place-items-center size-5 rounded-full bg-jv-ink text-white border-2 border-jv-canvas"
                    aria-hidden="true"
                  >
                    <RefreshCw class="size-2.5" :stroke-width="2.5" />
                  </span>
                </NavigationLink>

                <!-- Username -->
                <div
                  class="flex-1 flex items-center gap-2.5 bg-jv-white border-2 border-jv-ink jv-card shadow-brutal-sm focus-within:translate-x-[1px] focus-within:translate-y-[1px] focus-within:shadow-none transition-all px-3"
                >
                  <User
                    class="size-[18px] text-jv-ink/70 shrink-0"
                    :stroke-width="2.2"
                  />
                  <input
                    id="username"
                    v-model.trim="username"
                    type="text"
                    name="username"
                    maxlength="12"
                    placeholder="Pick a name"
                    class="flex-1 min-w-0 bg-transparent border-0 outline-none font-body text-jv-ink placeholder:text-jv-ink/40 text-sm sm:text-base"
                  />
                </div>
              </div>
            </div>

            <!-- Welcome (signed in) -->
            <div
              v-else
              class="flex items-center gap-2 bg-jv-yellow-soft border-2 border-jv-ink jv-card shadow-brutal-sm px-3 py-2.5"
            >
              <User
                class="size-[18px] text-jv-ink/70 shrink-0"
                :stroke-width="2.2"
              />
              <span class="font-body text-jv-ink text-sm sm:text-base">
                Welcome,
                <span class="font-bold">{{ firstname || username }}</span>
              </span>
            </div>

            <!-- Submit -->
            <NavigationLink
              type="submit"
              class="mt-1 w-full bg-jv-coral text-white py-2.5 sm:py-3 text-sm sm:text-base"
              :disabled="authChecking || quickUserPending"
            >
              <template v-if="authChecking">
                <Loader2 class="size-[18px] animate-spin" :stroke-width="2.4" />
                <span>Loading…</span>
              </template>
              <template v-else-if="quickUserPending">
                <Loader2 class="size-[18px] animate-spin" :stroke-width="2.4" />
                <span>Joining…</span>
              </template>
              <template v-else>
                <span>Enter Lobby</span>
                <ArrowRight class="size-[18px]" :stroke-width="2.4" />
              </template>
            </NavigationLink>
          </form>
        </div>
      </div>

      <!-- Below-card links -->
      <div
        class="mt-6 flex flex-col gap-1.5 text-center font-body text-sm text-jv-ink/70"
      >
        <p class="m-0">
          Don't have a code?
          <NuxtLink
            to="/"
            class="text-jv-coral underline underline-offset-4 decoration-2 font-semibold ml-1"
          >
            Browse Public Games
          </NuxtLink>
        </p>
        <p v-if="!isUserLoggedIn" class="m-0">
          Want to save your progress?
          <NuxtLink
            to="/account/login"
            class="text-jv-coral underline underline-offset-4 decoration-2 font-semibold ml-1"
          >
            Login Now
          </NuxtLink>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { usePush } from "notivue";
import { useUsersStore } from "~~/store/users";
import { useSessionStore } from "~~/store/session";
import { getRandomAvatarName, getAvatarUrlByName } from "~~/composables/avatar";
import {
  Hash,
  User,
  Check,
  RefreshCw,
  ArrowRight,
  Loader2,
  Sparkle,
} from "lucide-vue-next";
import NavigationLink from "@/components/common/NavigationLink.vue";

definePageMeta({
  layout: false,
});

const userData = useUsersStore();
const { setUserData } = userData;
const sessionStore = useSessionStore();
const { setActiveQuizTitle } = sessionStore;
const authChecking = ref(true);
const route = useRoute();

const codeparam = computed(() => route.query.code || "");
const code = ref(
  String(codeparam.value).replace(/\s+/g, "").replace(/\D/g, "").slice(0, 6)
);
const username = ref("");
const firstname = ref("");
const isUserLoggedIn = ref(false);
const { apiUrl } = useRuntimeConfig().public;
const router = useRouter();
const toast = usePush();
const userError = ref(false);
const quickUserPending = ref(false);
const userPlayedQuiz = ref("");
const sessionId = ref("");
const quizTitle = ref("");
const avatarName = ref("Sophia");
const avatarUrl = computed(() => getAvatarUrlByName(avatarName.value));

onMounted(() => {
  avatarName.value = getRandomAvatarName();
});

const codeDisplay = computed({
  get() {
    const c = code.value || "";
    return c.length > 3 ? c.slice(0, 3) + " " + c.slice(3) : c;
  },
  set(val) {
    code.value = (val || "").replace(/\s+/g, "").replace(/\D/g, "").slice(0, 6);
  },
});

const isCodeValid = computed(() => code.value.length === 6);

function rerollAvatar() {
  avatarName.value = getRandomAvatarName();
}

const join_quiz = async () => {
  username.value = username.value.trim().replace(/\s+/g, "");

  if (!code.value || code.value.length !== 6) {
    toast.error("Please enter a valid quiz code (6 characters)");
    return;
  }

  if (!username.value && !isUserLoggedIn.value) {
    toast.error("Please add your username");
    return;
  }

  if (username.value.length > 12 && !isUserLoggedIn.value) {
    toast.error("Username must be a maximum of 12 characters long");
    return;
  }

  // create quick user
  if (!isUserLoggedIn.value) {
    quickUserPending.value = true;

    try {
      await $fetch(
        `${apiUrl}/user/${username.value}?avatar_name=${avatarName.value}`,
        {
          method: "POST",
          credentials: "include",
          headers: {
            Accept: "application/json",
          },
          onResponse({ response }) {
            if (response.status == 200) {
              isUserLoggedIn.value = true;
              quickUserPending.value = false;
              firstname.value = response._data?.data?.first_name;
              username.value = response._data?.data?.username;
              setUserData({
                role: "guest-user",
                avatar: response._data?.data?.img_key,
              });
            }
          },
        }
      );
    } catch (error) {
      userError.value = error.message;
      quickUserPending.value = false;
      return;
    }
  }

  try {
    quickUserPending.value = true;
    await $fetch(`${apiUrl}/user_played_quizes/${code.value}`, {
      method: "POST",
      credentials: "include",
      headers: {
        Accept: "application/json",
      },
      onResponse({ response }) {
        if (response.status == 200) {
          userPlayedQuiz.value = response._data?.data?.user_played_quiz;
          sessionId.value = response._data?.data?.session_id;
          quizTitle.value = response._data?.data?.quiz_title || "";
          if (quizTitle.value) {
            setActiveQuizTitle(quizTitle.value);
          }
          quickUserPending.value = false;
        }
      },
    });
  } catch (error) {
    userError.value = error.message;
    quickUserPending.value = false;
    if (error?.status == 400) {
      toast.error("Please enter a valid quiz code (6 characters)");
      code.value = "";
    }
    return;
  }

  router.push(
    `/join/play/${code.value}?username=${encodeURIComponent(
      username.value
    )}&firstname=${firstname.value}&user_played_quiz=${
      userPlayedQuiz.value
    }&session_id=${sessionId.value}&quiz_title=${encodeURIComponent(
      quizTitle.value
    )}`
  );
};

// get user data — 401 is the expected guest case, so don't throw on it.
(async () => {
  try {
    const response = await $fetch.raw(apiUrl + "/user/who", {
      method: "GET",
      credentials: "include",
      headers: { Accept: "application/json" },
      ignoreResponseError: true,
    });
    if (response.status === 200) {
      isUserLoggedIn.value = true;
      firstname.value = response._data?.data?.firstname;
      username.value = response._data?.data?.username;
    } else if (response.status === 401) {
      isUserLoggedIn.value = false;
    } else {
      userError.value =
        response._data?.message || `Unexpected status ${response.status}`;
    }
  } catch (error) {
    userError.value = error?.message || "Failed to load user";
  } finally {
    authChecking.value = false;
  }
})();
</script>
