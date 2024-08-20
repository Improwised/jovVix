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
const errors = ref({
  email: "",
  password: "",
});
const kratos_url = urls.kratos_url;
console.log();

(async () => {
  if (process.client) {
    const user = getUserData();
    if (user && user == "admin-user") {
      navigateTo("/");
      return;
    }
    if (route.query.flow) {
      try {
        await $fetch(kratos_url + "/self-service/login/flows", {
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
            if (response?._data?.ui?.messages[0]?.type === "error") {
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
      kratos_url + "/self-service/login/browser",
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
    csrfToken.value = kratosResponse?.ui?.nodes[1]?.attributes?.value;
    loginURLWithFlowQuery.value = kratosResponse?.ui?.action;
  } catch (error) {
    console.error(error);
  }
}

const handleForgotPassword = async () => {
  if (!email.value) {
    toast.error("please enter email first!")
    return;
  }
  const recoveryResponse = await fetch(
    `${kratos_url}/self-service/recovery/browser`,
    {
      headers: {
        Accept: "application/json",
      },
    }
  );
  const recovery = await recoveryResponse.json();

  const recoverypage = recovery?.ui?.action;

  const recoveryCodeResponse = await fetch(recoverypage, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    body: JSON.stringify({
      email: email.value,
      csrf_token: csrfToken.value, // Include CSRF token in the request body
      method: "code", // Specify the method (e.g., code) for password recovery
    }),
  });

  const recoveryCode = await recoveryCodeResponse.json();

  navigateTo(recoverypage, { external: true });
};
</script>

<template>
  <QuizLoadingSpace v-if="component === 'waiting'"></QuizLoadingSpace>
  <div
    v-else
    class="row align-items-center justify-content-center g-0 min-vh-100"
  >
    <div class="col-12 col-md-8 col-lg-6 col-xxl-4 py-8 py-xl-0">
      <!-- Card -->
      <div class="card smooth-shadow-md">
        <!-- Card body -->
        <div class="card-body p-6">
          <div class="mb-4">
            <NuxtLink to="/" style="font-size: 30px; font-weight: bold"
              >Quiz App</NuxtLink
            >
            <p class="mb-6">Please enter your user information.</p>
          </div>
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
                class="form-control"
                id="password"
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
                    @click.prevent="handleForgotPassword"
                    class="text-inherit fs-5"
                  >
                    Forgot your password?
                  </button>
                </div>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>
