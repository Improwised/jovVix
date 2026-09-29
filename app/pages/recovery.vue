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
        <header class="mb-6 flex flex-col items-center gap-2">
          <NuxtLink to="/" class="inline-block">
            <NuxtImg
              src="/images/jovvix-logo.png"
              alt="Jovvix"
              width="89"
              height="36"
              class="h-9 w-auto"
            />
          </NuxtLink>
          <div class="relative inline-block">
            <h1
              class="m-0 font-headings text-[26px] leading-none text-jv-ink sm:text-[32px]"
            >
              Recover Account
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
            If an account exists for this email, you will receive a recovery
            code shortly.
          </p>
        </header>

        <form class="flex flex-col gap-4" @submit.prevent="verifyOTP">
          <div class="flex flex-col gap-1.5">
            <label
              for="otp"
              class="px-0.5 font-body text-xs font-bold uppercase tracking-wide text-jv-ink sm:text-[13px]"
            >
              OTP Code
            </label>
            <div
              class="jv-card flex items-center gap-2.5 border-2 border-jv-ink bg-jv-white px-3 py-2.5 shadow-brutal-sm transition-all focus-within:translate-x-[1px] focus-within:translate-y-[1px] focus-within:shadow-none"
              :class="{ 'border-jv-coral': otpError }"
            >
              <KeyRound
                class="size-[18px] shrink-0 text-jv-ink/70"
                :stroke-width="2.2"
              />
              <input
                id="otp"
                v-model="otp"
                type="text"
                placeholder="Enter OTP code"
                class="min-w-0 flex-1 border-0 bg-transparent font-body text-sm tracking-[0.15em] text-jv-ink outline-none placeholder:text-jv-ink/40 sm:text-base"
              />
            </div>
            <p
              v-if="otpError"
              class="flex items-center gap-1.5 px-0.5 font-body text-xs font-bold text-jv-coral"
            >
              <AlertCircle class="size-3.5" :stroke-width="2.4" />
              {{ otpError }}
            </p>
          </div>

          <button
            type="submit"
            :disabled="isLoading"
            class="jv-card mt-2 inline-flex h-12 items-center justify-center gap-2 border-2 border-jv-ink bg-jv-coral font-headings text-base text-white shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none disabled:cursor-not-allowed disabled:opacity-70 disabled:hover:rotate-0 sm:text-lg"
          >
            {{ isLoading ? "Verifying..." : "Submit" }}
            <ArrowRight v-if="!isLoading" class="size-5" :stroke-width="2.4" />
          </button>
        </form>

        <div class="mt-5 text-center font-body text-sm text-jv-ink/70">
          Back to
          <NuxtLink
            to="/account/login"
            class="font-bold text-jv-coral hover:underline"
          >
            Sign in
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import { KeyRound, ArrowRight, AlertCircle } from "lucide-vue-next";

const config = useRuntimeConfig();
const { kratosUrl } = config.public;

definePageMeta({
  layout: "auth",
});

useSeoMeta({
  title: "Account Recovery - jovVix",
  description:
    "Recover access to your jovVix account using your secure recovery link.",
  robots: "noindex, nofollow",
});

const otp = ref("");
const otpError = ref("");
const flow = ref("");
const csrfToken = ref("");
const isLoading = ref(false);

const route = useRoute();

onMounted(async () => {
  await fetchFlowIdAndCsrfToken();
});

const updateFlowState = (data) => {
  const tokenNode = data?.ui?.nodes?.find(
    (node) => node?.attributes?.name === "csrf_token"
  );
  if (tokenNode?.attributes?.value) {
    csrfToken.value = tokenNode.attributes.value;
  }

  const messages = data?.ui?.messages || [];
  const errorMsg = messages.find((m) => m.type === "error");
  if (errorMsg) {
    otpError.value = errorMsg.text;
  }
};

const fetchFlowIdAndCsrfToken = async () => {
  try {
    flow.value = route.query.flow;

    if (!flow.value) {
      otpError.value = "Recovery session not found. Please request a new code.";
      return;
    }

    const response = await fetch(
      `${kratosUrl}/self-service/recovery/flows?id=${flow.value}`,
      {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
        credentials: "include",
      }
    );

    if (!response.ok) {
      throw new Error(`Failed to fetch recovery flow: ${response.statusText}`);
    }

    const data = await response.json();
    updateFlowState(data);
  } catch (error) {
    console.error("Error fetching flow ID and CSRF token:", error.message);
  }
};

const verifyOTP = async () => {
  try {
    otpError.value = "";
    const cleanOtp = otp.value ? otp.value.trim() : "";
    if (!cleanOtp) {
      otpError.value = "Please enter the OTP code";
      return;
    }

    if (!flow.value) {
      otpError.value = "Recovery session is missing. Please restart recovery.";
      return;
    }

    isLoading.value = true;

    const otpVerificationResponse = await fetch(
      `${kratosUrl}/self-service/recovery?flow=${flow.value}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          code: cleanOtp,
          csrf_token: csrfToken.value,
          method: "code",
        }),
      }
    );

    const responseData = await otpVerificationResponse.json();

    // 1. Successful verification in Kratos returns 422 with browser_location_change_required
    if (
      otpVerificationResponse.status === 422 &&
      responseData?.error?.id === "browser_location_change_required"
    ) {
      navigateTo("/account/change-password");
      return;
    }

    // 2. Kratos returned HTTP 200 (re-rendered flow with error messages or status)
    if (otpVerificationResponse.ok) {
      updateFlowState(responseData);

      const uiError = responseData?.ui?.messages?.find(
        (m) => m.type === "error"
      )?.text;

      const nodeError = responseData?.ui?.nodes
        ?.flatMap((n) => n.messages || [])
        ?.find((m) => m.type === "error")?.text;

      if (uiError || nodeError) {
        otpError.value = uiError || nodeError;
        return;
      }

      if (responseData.state === "passed_challenge") {
        navigateTo("/account/change-password");
        return;
      }
    }

    // 3. Any other non-OK response from Kratos
    const uiError = responseData?.ui?.messages?.find(
      (m) => m.type === "error"
    )?.text;
    const generalError =
      responseData?.error?.message ||
      responseData?.message ||
      responseData?.error?.reason;
    otpError.value =
      uiError || generalError || "The recovery code is invalid or has expired.";
  } catch (error) {
    console.error("Error verifying OTP:", error.message);
    otpError.value = error.message || "Failed to verify OTP";
  } finally {
    isLoading.value = false;
  }
};
</script>
