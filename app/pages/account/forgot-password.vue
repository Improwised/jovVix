<script setup>
import { usePush } from "notivue";
import { Mail, ArrowRight } from "lucide-vue-next";

definePageMeta({
  layout: "auth",
});

const seoCanonical = new URL(
  useRoute().path,
  useRuntimeConfig().public.baseUrl
).href;

useHead({
  link: [{ rel: "canonical", href: seoCanonical }],
});

useSeoMeta({
  title: "Reset Password - jovVix",
  description:
    "Reset your jovVix account password securely. Enter your email to receive a recovery link and regain access to your quiz dashboard.",
  ogTitle: "Reset Password - jovVix",
  ogDescription:
    "Reset your jovVix account password securely. Enter your email to receive a recovery link and regain access to your quiz dashboard.",
  ogUrl: seoCanonical,
});

const toast = usePush();
const urls = useRuntimeConfig().public;
const kratosUrl = urls.kratosUrl;
const email = ref("");
const isLoading = ref(false);

const handleForgotPassword = async () => {
  if (!email.value) {
    toast.error("Please enter your email first!");
    return;
  }

  isLoading.value = true;

  try {
    const recoveryResponse = await fetch(
      `${kratosUrl}/self-service/recovery/browser`,
      {
        credentials: "include",
        headers: {
          Accept: "application/json",
        },
      }
    );
    const recovery = await recoveryResponse.json();

    const recoverypage = recovery?.ui?.action;
    const csrfToken = recovery?.ui?.nodes?.find(
      (node) => node?.attributes?.name === "csrf_token"
    )?.attributes?.value;

    await fetch(recoverypage, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify({
        email: email.value,
        csrf_token: csrfToken,
        method: "code",
      }),
    });

    navigateTo(recoverypage, { external: true });
  } catch (error) {
    console.error(error);
    toast.error("Something went wrong. Please try again.");
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div class="relative z-10 mx-auto w-full max-w-[440px]">
    <div class="relative rotate-1">
      <span
        class="jv-card absolute -top-[8px] left-1/2 z-20 h-3 w-12 -translate-x-1/2 border-2 border-jv-ink bg-jv-slate shadow-brutal-sm"
        aria-hidden="true"
      ></span>

      <div
        class="jv-card border-2 border-jv-ink bg-jv-white px-6 py-7 shadow-brutal-lg sm:px-8 sm:py-9"
      >
        <header class="mb-6 flex flex-col items-center gap-1.5">
          <div class="relative inline-block">
            <h1
              class="m-0 font-headings text-[26px] leading-none text-jv-ink sm:text-[32px]"
            >
              Forgot Password
            </h1>
            <svg
              class="absolute -bottom-2 left-1/2 -translate-x-1/2 text-jv-mint"
              width="140"
              height="14"
              viewBox="0 0 140 14"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
              aria-hidden="true"
            >
              <path
                d="M3 9 Q 25 1, 50 7 T 95 6 T 137 4"
                stroke="currentColor"
                stroke-width="2.5"
                stroke-linecap="round"
                fill="none"
              />
            </svg>
          </div>
          <p
            class="mt-2 text-center font-body text-sm text-jv-ink/70 sm:text-[15px]"
          >
            Enter your email and we'll send you a reset code.
          </p>
        </header>

        <form
          class="flex flex-col gap-4"
          @submit.prevent="handleForgotPassword"
        >
          <div class="flex flex-col gap-1.5">
            <label
              for="email"
              class="px-0.5 font-body text-xs font-bold uppercase tracking-wide text-jv-ink sm:text-[13px]"
            >
              Email Address
            </label>
            <div
              class="jv-card flex items-center gap-2.5 border-2 border-jv-ink bg-jv-white px-3 py-2.5 shadow-brutal-sm transition-all focus-within:translate-x-[1px] focus-within:translate-y-[1px] focus-within:shadow-none"
            >
              <Mail
                class="size-[18px] shrink-0 text-jv-ink/70"
                :stroke-width="2.2"
              />
              <input
                id="email"
                v-model="email"
                type="email"
                placeholder="you@example.com"
                class="min-w-0 flex-1 border-0 bg-transparent font-body text-sm text-jv-ink outline-none placeholder:text-jv-ink/40 sm:text-base"
                required
              />
            </div>
          </div>

          <button
            type="submit"
            :disabled="isLoading"
            class="jv-card mt-2 inline-flex h-12 items-center justify-center gap-2 border-2 border-jv-ink bg-jv-coral font-headings text-base text-white shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none disabled:cursor-not-allowed disabled:opacity-70 disabled:hover:rotate-0 sm:text-lg"
          >
            {{ isLoading ? "Sending..." : "Send Reset Code" }}
            <ArrowRight v-if="!isLoading" class="size-5" :stroke-width="2.4" />
          </button>
        </form>

        <div class="mt-5 text-center font-body text-sm text-jv-ink/70">
          <NuxtLink
            to="/account/login"
            class="font-bold text-jv-coral hover:underline"
          >
            Back to Sign In
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>
