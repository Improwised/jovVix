<template>
  <div class="join-page-container">
    <div class="background-container">
      <ul class="circles">
        <li></li>
        <li></li>
        <li></li>
        <li></li>
        <li></li>
        <li></li>
        <li></li>
        <li></li>
        <li></li>
        <li></li>
      </ul>
    </div>
    <div class="content-container">
      <div
        class="container h-100 d-flex align-items-center justify-content-center"
      >
        <div class="row w-100">
          <div class="col-12">
            <Frame
              page-title="Join Page"
              page-message="Let's Play Together"
              class="bg-white"
            >
              <form method="POST" @submit.prevent="join_quiz">
                <div class="mb-3 pe-3">
                  <label for="code" class="form-label purple-text"
                    >Invitation Code</label
                  >
                  <!-- Assuming v-otp-input is a custom component or external library -->
                  <v-otp-input
                    v-model="code"
                    max-width="500"
                    min-height="20"
                    type="number"
                    placeholder="0"
                    class="purple-text"
                  />
                </div>
                <div  class="mb-3">
                  <label for="username" class="form-label purple-text" v-if="!isUserLoggedIn"
                    >User Name</label
                  >
                  <input
                    id="username"
                    v-model.trim="username"
                    type="text"
                    name="username"
                    class="purple-text form-control"
                    v-if="!isUserLoggedIn"
                  />
                  <div v-if="isUserLoggedIn">Welcome <span class="font-weight-bold">{{ user?.name?.first }}</span></div>
                </div>
                <div class="p-2">
                  <div v-if="!isUserLoggedIn" class="text-center">
                    Want To Save Your Progress?
                    <NuxtLink to="/account/login"><b>Login</b></NuxtLink> Now.
                  </div>
                </div>
                <button
                  type="submit"
                  class="btn btn-primary btn-lg w-100 text-white join-button"
                >
                  Join
                </button>
              </form>
            </Frame>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useToast } from "vue-toastification";

// Assuming these are imported from external libraries or mixins
const code = ref("");
const username = ref("");
const user = ref({});
const isUserLoggedIn = ref(false)
const kratosURL = useRuntimeConfig().public.kratos_url
const router = useRouter();
const toast = useToast();

function join_quiz() {
  username.value = username.value.trim();

  if (!code.value || code.value.length !== 6) {
    toast.error("Please enter a valid quiz code (6 characters)");
    return;
  }

  if (!username.value && !isUserLoggedIn.value) {
    toast.error("Please add your username");
    return;
  }

  if (username.value.length > 12 && !isUserLoggedIn.value) {
    toast.error("Username must be a maximum of 12 characters long");
    return;
  }

  router.push(
    `/join/play/${code.value}?username=${encodeURIComponent(username.value)}`
  );
}

(async () => {
  try {
    await $fetch(kratosURL+"/sessions/whoami", {
      method: "GET",
      credentials: "include",
      headers: {
        Accept: "application/json",
      },
      onResponse({ response }) {
        if (response.status >= 200 && response.status < 300) {
          isUserLoggedIn.value = true;
          user.value = response?._data?.identity?.traits
          username.value = user?.value?.name?.first
        }
      },
    })
  } catch (error) {
    
  }
})();
</script>

<style scoped>
body {
  font-family: "Exo", sans-serif;
  margin: 0;
  padding: 0;
  overflow: hidden;
}

.purple-text {
  color: #663399;
  font-weight: bold;
}

.background-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: url("../../assets/images/join-page-bg.jpg");
  /* background: -webkit-linear-gradient(to left, #8f94fb, #4e54c8); */
  background-position: center;
  background-size: cover;
  z-index: -1;
}

.circles {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.circles li {
  position: absolute;
  display: block;
  list-style: none;
  width: 20px;
  height: 20px;
  background: #fff0f5;
  animation: animate 25s linear infinite;
  bottom: -150px;
}

.circles li:nth-child(1) {
  left: 25%;
  width: 80px;
  height: 80px;
  animation-delay: 0s;
}

.circles li:nth-child(2) {
  left: 10%;
  width: 20px;
  height: 20px;
  animation-delay: 2s;
  animation-duration: 12s;
}

.circles li:nth-child(3) {
  left: 70%;
  width: 20px;
  height: 20px;
  animation-delay: 4s;
}

.circles li:nth-child(4) {
  left: 40%;
  width: 60px;
  height: 60px;
  animation-delay: 0s;
  animation-duration: 18s;
}

.circles li:nth-child(5) {
  left: 65%;
  width: 20px;
  height: 20px;
  animation-delay: 0s;
}

.circles li:nth-child(6) {
  left: 75%;
  width: 110px;
  height: 110px;
  animation-delay: 3s;
}

.circles li:nth-child(7) {
  left: 35%;
  width: 150px;
  height: 150px;
  animation-delay: 7s;
}

.circles li:nth-child(8) {
  left: 50%;
  width: 25px;
  height: 25px;
  animation-delay: 15s;
  animation-duration: 45s;
}

.circles li:nth-child(9) {
  left: 20%;
  width: 15px;
  height: 15px;
  animation-delay: 2s;
  animation-duration: 35s;
}

.circles li:nth-child(10) {
  left: 85%;
  width: 150px;
  height: 150px;
  animation-delay: 0s;
  animation-duration: 11s;
}

@keyframes animate {
  0% {
    transform: translateY(0) rotate(0deg);
    opacity: 1;
    border-radius: 0;
  }
  100% {
    transform: translateY(-1000px) rotate(720deg);
    opacity: 0;
    border-radius: 50%;
  }
}

.join-page-container {
  position: relative;
  overflow-x: hidden; /* Prevent horizontal scrolling */
}

.content-container {
  position: relative;
  z-index: 2;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  padding: 20px;
}

.join-button {
  background: linear-gradient(270deg, #5a66ef 0, #8042e4);
}

.join-button:hover {
  background: #6f4eb8; /* Slightly darker shade */
}
</style>
