<script setup>
// core dependencies
import { useToast } from "vue-toastification";

// define nuxt configs
const toast = useToast();
const router = useRouter();
useSystemEnv();

// define props and emits
const code = ref(0);
const username = ref();

function join_quiz(e) {
  e.preventDefault();

  if (!code.value || code.value.length != 6) {
    toast.error("please enter quiz code");
    return;
  } else if (!username.value) {
    toast.error("please add username");
    return;
  } else if (username.value.length > 12 || username.value.length < 6) {
    toast.error("username must be 6 to 12 char long");
  } else {
    router.push(`/join/play/${code.value}?username=${username.value}`);
  }
}
</script>

<template>
  <Frame page-title="Join page" page-message="Let's play together" max-width>
    <form method="POST" @submit="join_quiz">
      <div class="mb-3 pe-3">
        <label for="code" class="form-label">Code</label>
        <v-otp-input
          v-model="code"
          max-width="500"
          min-height="20"
          type="number"
        ></v-otp-input>
      </div>
      <div class="mb-3">
        <label for="username" class="form-label">User name</label>
        <input
          id="username"
          v-model="username"
          type="username"
          name="username"
          class="form-control"
          required
        />
      </div>
      <div class="p-2">
        <div class="text-center">
          Want to save your progress?
          <NuxtLink to="/account/login"><b>Login</b></NuxtLink> now.
        </div>
      </div>
      <button type="submit" class="btn btn-primary btn-lg bg-primary">
        Join
      </button>
    </form>
  </Frame>
</template>
