<script setup>
import { useToast } from "vue-toastification";
import { useUsersStore } from "~~/store/users";
const userData = useUsersStore();
const { getUserData } = userData;
const route = useRoute();
const router = useRouter();
const toast = useToast();
const csrfToken = ref();
const component = ref("waiting");
const loginURLWithFlowQuery = ref("");
const urls = useRuntimeConfig().public;
const email = ref();
const code = ref("");
const flowID = ref("");
const errors = ref({
  email: "",
  password: "",
  code: "",
});
const kratosUrl = urls.kratosUrl;
console.log();

(async () => {
  if (process.client) {
    const user = getUserData();
    if (user && user?.role == "admin-user") {
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
            if (response._data?.error?.code === 410) {
              navigateTo("/account/login");
            }

            if (response.status === 400) {
              toast.warning(
                "Please fill out the form correctly, password or email was incorrect"
              );
            }
            component.value = "form";
          },
          onResponse({ response }) {
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
        component.value = "form";
      } catch (error) {
        console.error(error);
        await setFlowIDAndCSRFToken();
        component.value = "form";
      }
    } else {
      await setFlowIDAndCSRFToken();
      component.value = "form";
    }
  } else {
    await setFlowIDAndCSRFToken();
    component.value = "form";
  }
})();

async function setFlowIDAndCSRFToken() {
  try {
    const kratosResponse = await $fetch(
      kratosUrl + "/self-service/login/browser?refresh=true",
      {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
        credentials: "include",
        onResponseError({ response }) {
          console.error(
            "error while getting the flow id from the server",
            response
          );
        },
      }
    );

    router.push("?flow=" + kratosResponse?.id);
    csrfToken.value = kratosResponse?.ui?.nodes.find(
      (node) => node.attributes.name === "csrf_token"
    )?.attributes?.value;
    loginURLWithFlowQuery.value = kratosResponse?.ui?.action;
  } catch (error) {
    console.error(error);
  }
}

const handleForgotPassword = async () => {
  if (!email.value) {
    toast.error("please enter email first!");
    return;
  }
  const recoveryResponse = await fetch(
    `${kratosUrl}/self-service/recovery/browser`,
    {
      headers: {
        Accept: "application/json",
      },
    }
  );
  const recovery = await recoveryResponse.json();

  const recoverypage = recovery?.ui?.action;

  await fetch(recoverypage, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    body: JSON.stringify({
      email: email.value,
      csrf_token: csrfToken.value,
      method: "code",
    }),
  });

  navigateTo(recoverypage, { external: true });
};

// Handle email verification
const handleEmailVerification = async () => {
  if (!email.value) {
    toast.error("Please enter email first!");
    return;
  }

  try {
    // Request email verification flow from Kratos
    const verificationResponse = await fetch(
      `${kratosUrl}/self-service/verification/browser`,
      {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
      }
    );
    const verification = await verificationResponse.json();
    const verificationPage = verification?.ui?.action;
    flowID.value = verification?.id;

    // Trigger email verification
    await fetch(verificationPage, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify({
        email: email.value,
        csrf_token: csrfToken.value,
        method: "code",
      }),
    });

    toast.success("Verification email has been sent!");
    component.value = "verifyCode";
  } catch (error) {
    toast.error("An error occurred while sending the verification email.");
    console.error(error);
  }
};

// Handle verification code input
const verifyCode = async () => {
  if (!code.value) {
    errors.value.code = "Verification code is required.";
    return;
  }

  try {
    const response = await fetch(
      `${kratosUrl}/self-service/verification?flow=${flowID.value}`,
      {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: JSON.stringify({
          csrf_token: csrfToken.value,
          method: "code",
          code: code.value,
        }),
      }
    );

    const result = await response.json();

    if (response.ok) {
      const messages = result?.ui?.messages;

      // code is invalid
      if (messages && messages[0]?.id == 4070006) {
        toast.error(
          "The verification code is invalid or has already been used. Please try again."
        );
        return;
      }

      toast.success("Email verified successfully!");
      component.value = "form";
    } else {
      toast.error(
        result?.error?.message || "Verification failed, please try again."
      );
      errors.value.code = "The verification code is invalid or expired.";
    }
  } catch (error) {
    toast.error("An error occurred during verification.");
    console.error(error);
  }
};
</script>

<template>
  <QuizLoadingSpace v-if="component === 'waiting'"></QuizLoadingSpace>
  <!-- Verification of email -->
  <Frame
    v-else-if="component === 'verifyCode'"
    page-title="Email Verification"
    page-message="Please enter the verification code sent to your email"
  >
    <form @submit.prevent="verifyCode">
      <div class="mb-3">
        <label for="code" class="form-label">Verification Code</label>
        <input
          id="code"
          v-model="code"
          type="text"
          class="form-control"
          maxlength="6"
          required
        />
      </div>
      <div v-if="errors.code" class="text-danger">{{ errors.code }}</div>
      <button type="submit" class="btn btn-primary text-light">Verify</button>
    </form>
  </Frame>
  <Frame
    v-else
    page-title="Login"
    page-message="Please enter your user information."
  >
    <!-- Form -->
    <form method="POST" :action="loginURLWithFlowQuery">
      <!-- Username -->
      <div class="mb-3">
        <label for="email" class="form-label">Email</label>
        <input
          id="email"
          v-model="email"
          type="email"
          name="identifier"
          class="form-control"
          required=""
        />
      </div>
      <!-- Password -->
      <div class="mb-3">
        <label for="password" class="form-label">Password</label>
        <input
          id="password"
          class="form-control"
          type="password"
          name="password"
          required=""
        />
      </div>
      <div v-if="errors.password" class="text-danger">
        {{ errors.password }}
      </div>
      <div>
        <input type="hidden" name="csrf_token" :value="csrfToken" />
      </div>
      <div>
        <input type="hidden" name="method" value="password" />
      </div>
      <div>
        <!-- Button -->
        <div class="d-grid">
          <button type="submit" class="btn btn-primary text-white">
            Sign in
          </button>
        </div>

        <div class="d-md-flex justify-content-between mt-4">
          <div class="mb-2 mb-md-0">
            <NuxtLink to="/account/register" class="fs-5"
              >Create An Account
            </NuxtLink>
          </div>
          <div>
            <button
              class="text-primary mb-2 fs-5"
              @click.prevent="handleForgotPassword"
            >
              Forgot your password?
            </button>
          </div>
          <div>
            <button
              class="text-primary fs-5"
              @click.prevent="handleEmailVerification"
            >
              Verify your email
            </button>
          </div>
        </div>
      </div>
    </form>
  </Frame>
</template>
