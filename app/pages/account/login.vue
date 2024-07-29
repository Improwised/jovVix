<script setup>
import { useToast } from "vue-toastification";

const route = useRoute();
const router = useRouter();
const toast = useToast();
const csrfToken = ref();
const component = ref("waiting");
const loginURLWithFlowQuery = ref("")
const urls = useRuntimeConfig().public;
const errors = ref({
  email: "",
  password: "",
});
const kratos_url = urls.kratos_url;
console.log();

(async () => {
  if (process.client) {
    try {
      await $fetch(kratos_url + "/sessions/whoami", {
        method: "GET",
        credentials: "include",
        headers: {
          Accept: "application/json",
        },
        onResponse({ response }) {
          if (response.status >= 200 && response.status < 300) {
            toast.success("You are already logged in");
            navigateTo("/");
          }
        },
      });
    } catch (error) {
      console.log(); // this error will be produced when user is not logged in so not logging anything here or not doing anything here
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
              errors.value.password = "The provided credentials are invalid, check for spelling mistakes in your password or email"
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
</script>

<template>
  <Frame page-title="Login Page" page-message="Welcome To The Quizz World...">
    <form method="POST" :action="loginURLWithFlowQuery">
      <div class="mb-3">
        <label for="email" class="form-label">Email</label>
        <input
          id="email"
          v-model="email"
          type="email"
          name="identifier"
          class="form-control"
          required
        />
      </div>
      <div class="mb-3">
        <label for="password" class="form-label">Password</label>
        <input
          id="password"
          v-model="password"
          type="password"
          name="password"
          class="form-control"
          required
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
      <div class="p-2">
        <div class="text-end">
          Don't Have An
          <NuxtLink to="/account/register"><b>Account</b></NuxtLink>
          ?
        </div>
      </div>
      <button type="submit" class="btn text-white btn-primary">Submit</button>
      <button type="reset" class="btn text-white btn-primary ms-2">
        Clear
      </button>
    </form>
  </Frame>
</template>
