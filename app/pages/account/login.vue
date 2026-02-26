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
  code: "",
});
const kratosUrl = urls.kratosUrl;
console.log();

(async () => {
  if (process.client) {
    const user = getUserData();
    const isReauth = !!route.query.returnTo;

    if (!isReauth && user) {
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
    // Build return_to URL
    const returnToUrl = route.query.returnTo
      ? `${window.location.origin}${route.query.returnTo}`
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
          return_to: returnToUrl
        },
        onResponseError({ response }) {
          console.error(
            "error while getting the flow id from the server",
            response
          );
        },
      }
    );
    const queryParams = route.query.returnTo
      ? `?flow=${kratosResponse?.id}&returnTo=${route.query.returnTo}`
      : `?flow=${kratosResponse?.id}`;

    router.push(queryParams);
    csrfToken.value = kratosResponse?.ui?.nodes.find(
      (node) => node.attributes.name === "csrf_token"
    )?.attributes?.value;
    loginURLWithFlowQuery.value = kratosResponse?.ui?.action;

    const identifierNode = kratosResponse?.ui?.nodes.find(
      (node) => node.attributes.name === "identifier"
    );
    if (identifierNode?.attributes?.value) {
      email.value = identifierNode.attributes.value;
    }
  } catch (error) {
    console.error(error);
  }
}

</script>

<template>
  <QuizLoadingSpace v-if="component === 'waiting'"></QuizLoadingSpace>
  <Frame v-else page-title="Sign in" page-message="Please enter your user information.">
    <!-- Form -->
    <form method="POST" :action="loginURLWithFlowQuery">
      <!-- Username -->
      <div class="mb-3">
        <label for="email" class="form-label">Email</label>
        <input id="email" v-model="email" type="email" name="identifier" class="form-control" :readonly="!!route.query.returnTo" :style="route.query.returnTo ? { backgroundColor: '#e9ecef', cursor: 'not-allowed' } : {}" required="" />
      </div>
      <!-- Password -->
      <div class="mb-3">
        <label for="password" class="form-label">Password</label>
        <input id="password" class="form-control" type="password" name="password" required="" />
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
            <NuxtLink to="/account/register" class="fs-5">Create An Account
            </NuxtLink>
          </div>
          <div>
            <NuxtLink to="/account/forgot-password" class="fs-5">Forgot Your Password?
            </NuxtLink>
          </div>
        </div>
      </div>
    </form>
  </Frame>
</template>
