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
            <img src="/jovvix-logo.png" alt="Jovvix" class="h-9 w-auto" />
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
            class="jv-card mt-2 inline-flex h-12 items-center justify-center gap-2 border-2 border-jv-ink bg-jv-coral font-headings text-base text-white shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:text-lg"
          >
            Submit
            <ArrowRight class="size-5" :stroke-width="2.4" />
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

const otp = ref("");
const otpError = ref("");
const flow = ref("");
const csrfToken = ref("");

const route = useRoute();

onMounted(async () => {
  await fetchFlowIdAndCsrfToken();
});

const fetchFlowIdAndCsrfToken = async () => {
  try {
    flow.value = route.query.flow;

    const response = await fetch(
      `${kratosUrl}/self-service/recovery/browser?aal=&refresh=&return_to=`,
      {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
        credentials: "include",
      }
    );

    if (!response.ok) {
      throw new Error(`Failed to fetch CSRF token: ${response.statusText}`);
    }

    const data = await response.json();
    csrfToken.value = data.ui.nodes.find(
      (node) => node.attributes.name === "csrf_token"
    ).attributes.value;
  } catch (error) {
    console.error("Error fetching flow ID and CSRF token:", error.message);
  }
};

const verifyOTP = async () => {
  try {
    otpError.value = "";
    if (!otp.value) {
      otpError.value = "Please enter the OTP code";
      return;
    }

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
          code: otp.value,
          csrf_token: csrfToken.value,
          method: "code",
        }),
      }
    );

    if (!otpVerificationResponse.ok) {
      const errorData = await otpVerificationResponse.json();
      if (errorData.messages && errorData.messages.length > 0) {
        const errorMessage = errorData.messages[0].text;
        otpError.value = errorMessage;
      } else if (
        errorData.error &&
        errorData.error.id === "browser_location_change_required"
      ) {
        navigateTo("/account/change-password");
      } else {
        throw new Error(errorData.message || "Failed to verify OTP");
      }
      return;
    }
  } catch (error) {
    console.error("Error verifying OTP:", error.message);
    otpError.value = error.message;
  }
};
</script>
