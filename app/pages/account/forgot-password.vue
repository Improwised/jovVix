<script setup>
import { useToast } from "vue-toastification";

const toast = useToast();
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
  <Frame
    page-title="Forgot Password"
    page-message="Enter your email and we'll send you a reset code."
  >
    <form @submit.prevent="handleForgotPassword">
      <!-- Email -->
      <div class="mb-3">
        <label for="email" class="form-label">Email</label>
        <input
          id="email"
          v-model="email"
          type="email"
          class="form-control"
          placeholder="Enter your email"
          required
        />
      </div>

      <!-- Submit Button -->
      <div class="d-grid mb-4">
        <button
          type="submit"
          class="btn btn-primary text-white"
          :disabled="isLoading"
        >
          {{ isLoading ? "Sending..." : "Send Reset Password Link" }}
        </button>
      </div>

      <!-- Back to login -->
      <div class="text-center">
        <NuxtLink to="/account/login" class="fs-5">
          Back to Sign In
        </NuxtLink>
      </div>
    </form>
  </Frame>
</template>