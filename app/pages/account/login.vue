<script setup>
import { useToast } from "vue-toastification";

const toast = useToast();
const config = useRuntimeConfig();
const email = ref();
const password = ref();

async function login_user(e) {
  e.preventDefault();
  if (email.value.trim() == "") {
    alert("email and password are required");
  }

  const { error: error } = await useFetch(config.api_url + "/login", {
    method: "POST",
    credentials: "include",
    body: {
      email: email.value,
      password: password.value,
    },
    mode: 'cors'
  });

  if (error.value) {
    console.log(!error.value);
    toast.error("Error or password incorrect, try again");
    return;
  }

  toast.success("Login successfully!!!");
  setTimeout(() => navigateTo('/') , 3000)
}
</script>

<template>
  <AuthLayout
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
  </AuthLayout>
</template>
