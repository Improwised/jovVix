<script setup>
import { useToast } from "vue-toastification";

const config = useSystemEnv();
const nuxtApp = useNuxtApp();
const route = useRoute();
const toast = useToast();
const email = ref();
const password = ref();
let status = null;

async function login_user(e) {
  e.preventDefault();

  if (email.value.trim() == "" || password.value.trim() == "") {
    toast.error(nuxtApp.$IncorrectCredentials);
    return;
  }

  const { error: error } = await useFetch(config.value.api_url + "/login", {
    method: "POST",
    credentials: "include",
    body: {
      email: email.value,
      password: password.value,
    },
    mode: "cors",
    onResponseError: function ({ response }) {
      status = response.status;
    },
  });

  if (error.value) {
    if (status >= 500) {
      toast.error(error.value);
    } else {
      toast.error(nuxtApp.$IncorrectCredentials);
    }
    return;
  }

  if (route.query.url) {
    navigateTo(route.query.url);
  } else {
    navigateTo("/");
  }
}
</script>

<template>
  <Frame
    page-title="Login page"
    page-welcome-message="Welcome to the quizz world..."
  >
    <form method="POST" @submit="login_user">
      <div class="mb-3">
        <label for="email" class="form-label">Email</label>
        <input
          id="email"
          v-model="email"
          type="email"
          name="email"
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
      <div class="p-2">
        <div class="text-end">
          don't have an <a href="/account/register"><b>account</b></a
          >?
        </div>
      </div>
      <button type="submit" class="btn btn-primary">Submit</button>
      <button type="clear" class="btn btn-primary ms-2">Clear</button>
    </form>
  </Frame>
</template>
