<script setup>
import { useToast } from "vue-toastification";

const nuxtApp = useNuxtApp();
const route = useRoute();
const router = useRouter();
const toast = useToast();
const email = ref();
const password = ref();
let status = null;
useSystemEnv();

async function login_user(e) {
  e.preventDefault();
  const urls = useSystemEnv("urls");

  if (urls === undefined) {
    toast.error("some error occurred. please reload the page!!!");
  }

  const login_url = urls.value?.api_url + "/login";

  if (email.value.trim() == "" || password.value.trim() == "") {
    toast.error(nuxtApp.$IncorrectCredentials);
    return;
  }

  const { error: error } = await useFetch(login_url, {
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

  if (error?.value) {
    if (status >= 500) {
      toast.error(error.value);
    } else {
      toast.error(nuxtApp.$IncorrectCredentials);
    }
    return;
  }

  if (route.query.url) {
    router.push(route.query.url);
  } else {
    router.push("/");
  }
}
</script>

<template>
  <Frame page-title="Login page" page-message="Welcome to the quizz world...">
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
          don't have an
          <NuxtLink to="/account/register"><b>account</b></NuxtLink>
          ?
        </div>
      </div>
      <button type="submit" class="btn btn-primary">Submit</button>
      <button type="reset" class="btn btn-primary ms-2">Clear</button>
    </form>
  </Frame>
</template>
