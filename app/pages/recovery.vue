<template>
  <div
    id="otp-recovery"
    class="d-flex justify-content-center align-items-center vh-100"
  >
    <div class="col-6">
      <div class="card shadow-lg">
        <div class="card-body text-center" style="min-height: 400px">
          <div>
            <NuxtLink to="/"> 
              <img class="logo" src="/jovvix-logo.png" alt="" />
            </NuxtLink>
            <h3 class="mt-3 welcome-text">Recover your account</h3>
          </div>

          <form @submit.prevent="verifyOTP">
            <div class="form-group mb-4 d-flex align-items-center flex-column">
              <input
                id="otp"
                v-model="otp"
                type="text"
                class="form-control otp-input"
                placeholder="Enter OTP code"
              />
              <span v-if="otpError" class="text-danger">{{ otpError }}</span>
            </div>

            <button
              class="btn btn-primary btn-block btn-lg shadow-sm mt-3 text-white"
              type="submit"
            >
              Submit
            </button>
          </form>

          <div class="mt-3">
            Back to
            <NuxtLink to="/account/login" class="text-primary">
              Sign in</NuxtLink
            >
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
const config = useRuntimeConfig();
const { kratosUrl } = config.public;
definePageMeta({
  layout: false,
});

const otp = ref("");
const otpError = ref("");
const flow = ref("");
const csrfToken = ref("");

const route = useRoute();

// Fetch flow ID and CSRF token on component mount
onMounted(async () => {
  await fetchFlowIdAndCsrfToken();
});

// Method to fetch flow ID and CSRF token
const fetchFlowIdAndCsrfToken = async () => {
  try {
    // Parse flow ID from route query parameters
    flow.value = route.query.flow;

    // Example endpoint to fetch CSRF token (replace with your actual endpoint)
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

    // Parse the response data
    const data = await response.json();
    csrfToken.value = data.ui.nodes.find(
      (node) => node.attributes.name === "csrf_token"
    ).attributes.value;
  } catch (error) {
    console.error("Error fetching flow ID and CSRF token:", error.message);
  }
};

// Method to verify OTP
const verifyOTP = async () => {
  try {
    otpError.value = "";
    if (!otp.value) {
      otpError.value = "Please enter the OTP code";
      return;
    }

    // Perform OTP verification request
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
        navigateTo("/admin");
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

<style scoped>
.otp-input {
  max-width: 350px;
}

.logo {
  height: 40px;
  transform: scale(1.8);
  margin-bottom: 20px;
}
</style>
