<template>
  <div class="relative w-full max-w-[420px] mx-auto z-10">
    <!-- Card with subtle tilt -->
    <div class="relative rotate-1">
      <!-- Decorative "clip" tab on top of the card -->
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
              Sign In
            </h1>
            <!-- Hand-drawn underline doodle -->
            <svg
              class="absolute -bottom-2 left-1/2 -translate-x-1/2 text-jv-mint"
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
            Welcome to Jovvix
          </p>
        </header>

        <!-- Form -->
        <form
          method="POST"
          :action="loginURLWithFlowQuery"
          class="flex flex-col gap-3.5"
        >
          <!-- Email -->
          <div class="flex flex-col gap-1.5">
            <label
              for="email"
              class="font-body text-jv-ink text-xs sm:text-[13px] font-bold tracking-wide uppercase px-0.5"
            >
              Email Address
            </label>
            <div
              class="flex items-center gap-2.5 bg-jv-white border-2 border-jv-ink jv-card shadow-brutal-sm focus-within:translate-x-[1px] focus-within:translate-y-[1px] focus-within:shadow-none transition-all px-3 py-2.5"
            >
              <Mail
                class="size-[18px] text-jv-ink/70 shrink-0"
                :stroke-width="2.2"
              />
              <input
                id="email"
                v-model="email"
                type="email"
                name="identifier"
                class="flex-1 min-w-0 bg-transparent border-0 outline-none font-body text-jv-ink placeholder:text-jv-ink/40 text-sm sm:text-base"
                :readonly="isIdentityLocked"
                :class="{
                  'cursor-not-allowed bg-jv-slate/40': isIdentityLocked,
                }"
                placeholder="you@example.com"
                required
              />
            </div>
          </div>

          <!-- Password -->
          <div class="flex flex-col gap-1.5">
            <div class="flex items-end justify-between px-0.5">
              <label
                for="password"
                class="font-body text-jv-ink text-xs sm:text-[13px] font-bold tracking-wide uppercase"
              >
                Password
              </label>
            </div>
            <div
              class="flex items-center gap-2.5 bg-jv-white border-2 border-jv-ink jv-card shadow-brutal-sm focus-within:translate-x-[1px] focus-within:translate-y-[1px] focus-within:shadow-none transition-all px-3 py-2.5"
            >
              <Lock
                class="size-[18px] text-jv-ink/70 shrink-0"
                :stroke-width="2.2"
              />
              <input
                id="password"
                :type="showPassword ? 'text' : 'password'"
                name="password"
                class="flex-1 min-w-0 bg-transparent border-0 outline-none font-body text-jv-ink placeholder:text-jv-ink/40 text-sm sm:text-base"
                placeholder="••••••••"
                required
              />
              <button
                type="button"
                class="grid place-items-center size-6 text-jv-ink/70 hover:text-jv-ink transition-colors"
                :aria-label="showPassword ? 'Hide password' : 'Show password'"
                @click="showPassword = !showPassword"
              >
                <EyeOff
                  v-if="showPassword"
                  class="size-[18px]"
                  :stroke-width="2.2"
                />
                <Eye v-else class="size-[18px]" :stroke-width="2.2" />
              </button>
            </div>
            <NuxtLink
              to="/account/forgot-password"
              class="font-body text-jv-ink/60 hover:text-jv-coral text-xs sm:text-[13px] no-underline transition-colors self-end"
            >
              Forgot Password?
            </NuxtLink>
            <p
              v-if="errors.password"
              class="font-body text-jv-accent-red text-xs px-0.5 m-0"
            >
              {{ errors.password }}
            </p>
          </div>

          <input type="hidden" name="csrf_token" :value="csrfToken" />
          <input type="hidden" name="method" value="password" />

          <!-- Submit button -->
          <NavigationLink
            type="submit"
            class="mt-1 w-full bg-jv-coral text-white py-2.5 sm:py-3 text-sm sm:text-base"
            :disabled="!csrfToken || !loginURLWithFlowQuery"
          >
            <template v-if="!csrfToken || !loginURLWithFlowQuery">
              <Loader2 class="size-[18px] animate-spin" :stroke-width="2.4" />
              <span>Loading…</span>
            </template>
            <template v-else>
              <span>Sign In</span>
              <ArrowRight class="size-[18px]" :stroke-width="2.4" />
            </template>
          </NavigationLink>
        </form>

        <!-- Footer link -->
        <p class="text-center mt-6 font-body text-jv-ink/70 text-sm m-0">
          Don't have an account?
          <NuxtLink
            :to="createAccountLink"
            class="text-jv-coral underline underline-offset-4 decoration-2 hover:decoration-jv-coral font-semibold ml-1"
          >
            Create an account
          </NuxtLink>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { Mail, Lock, Eye, EyeOff, ArrowRight, Loader2 } from "lucide-vue-next";
import { usePush } from "notivue";
import { useUsersStore } from "~~/store/users";
import NavigationLink from "@/components/common/NavigationLink.vue";

definePageMeta({
  layout: "auth",
});

const seoCanonical = new URL(useRoute().path, useRuntimeConfig().public.baseUrl)
  .href;

useHead({
  link: [{ rel: "canonical", href: seoCanonical }],
});

useSeoMeta({
  title: "Login - jovVix",
  description:
    "Log in to jovVix to create, host, and manage your real-time quizzes, view reports, and track participant scores from your dashboard.",
  ogTitle: "Login - jovVix",
  ogDescription:
    "Log in to jovVix to create, host, and manage your real-time quizzes, view reports, and track participant scores from your dashboard.",
  ogUrl: seoCanonical,
});

const userData = useUsersStore();
const { getUserData } = userData;
const route = useRoute();
const router = useRouter();
const toast = usePush();
const csrfToken = ref();
const loginURLWithFlowQuery = ref("");
const urls = useRuntimeConfig().public;
const email = ref();
const isIdentityLocked = ref(false);
const showPassword = ref(false);
const returnTo = ref(route.query.returnTo || "");
const createAccountLink = computed(() =>
  returnTo.value
    ? `/account/register?returnTo=${encodeURIComponent(returnTo.value)}`
    : "/account/register"
);
const errors = ref({
  email: "",
  password: "",
  code: "",
});
const kratosUrl = urls.kratosUrl;

(async () => {
  if (process.client) {
    const user = getUserData();
    const isReauth = !!route.query.returnTo;

    if (!isReauth && user && user.role !== "guest-user") {
      navigateTo("/");
      return;
    }
    if (route.query.flow) {
      try {
        await $fetch(kratosUrl + "/self-service/login/flows", {
          method: "GET",
          credentials: "include",
          headers: {
            Accept: "application/json",
          },
          query: {
            id: route.query.flow,
          },
          onResponseError({ response }) {
            if (response.status === 400) {
              toast.warning(
                "Please fill out the form correctly, password or email was incorrect"
              );
            }
          },
          onResponse({ response }) {
            if (!returnTo.value) {
              returnTo.value = returnToPathFromUrl(response?._data?.return_to);
            }
            const messages = response?._data?.ui?.messages;
            if (messages && messages[0]?.type === "error") {
              // error indicating unverified email
              if (messages[0]?.id === 4000010) {
                toast.info("Please verify your email before logging in.");
                return;
              }
              errors.value.password =
                "The provided credentials are invalid, check for spelling mistakes in your password or email";
            }
          },
        });
        await setFlowIDAndCSRFToken();
      } catch (error) {
        console.error(error);
        await setFlowIDAndCSRFToken();
      }
    } else {
      await setFlowIDAndCSRFToken();
    }
  } else {
    await setFlowIDAndCSRFToken();
  }
})();

async function setFlowIDAndCSRFToken() {
  try {
    // Build return_to URL
    const returnToUrl = returnTo.value
      ? `${window.location.origin}${returnTo.value}`
      : `${window.location.origin}/`;

    const kratosResponse = await $fetch(
      kratosUrl + "/self-service/login/browser",
      {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
        credentials: "include",
        query: {
          refresh: true,
          return_to: returnToUrl,
        },
        onResponseError({ response }) {
          console.error(
            "error while getting the flow id from the server",
            response
          );
        },
      }
    );
    const queryParams = returnTo.value
      ? `?flow=${kratosResponse?.id}&returnTo=${encodeURIComponent(
          returnTo.value
        )}`
      : `?flow=${kratosResponse?.id}`;

    router.replace(queryParams);
    csrfToken.value = kratosResponse?.ui?.nodes.find(
      (node) => node.attributes.name === "csrf_token"
    )?.attributes?.value;
    loginURLWithFlowQuery.value = kratosResponse?.ui?.action;

    const identifierNode = kratosResponse?.ui?.nodes.find(
      (node) => node.attributes.name === "identifier"
    );
    isIdentityLocked.value = !!identifierNode?.attributes?.value;
    if (isIdentityLocked.value) {
      email.value = identifierNode.attributes.value;
    }
  } catch (error) {
    console.error(error);
  }
}
</script>
