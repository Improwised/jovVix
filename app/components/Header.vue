<template>
  <div class="header">
    <!-- navbar -->
    <nav class="navbar navbar-expand-lg">
      <!-- Logo -->
      <NuxtLink class="navbar-brand navbar-logo" style="color: black" to="/">
        <img class="logo" src="@/assets/images/logo.png" alt="" />
        <!-- Quiz App -->
      </NuxtLink>
      <button
        class="navbar-toggler bg-light"
        type="button"
        data-bs-toggle="collapse"
        data-bs-target="#navbarNavAltMarkup"
        aria-controls="navbarNavAltMarkup"
        aria-expanded="false"
        aria-label="Toggle navigation"
      >
        <span class="navbar-toggler-icon text-dark"></span>
      </button>
      <div
        id="navbarNavAltMarkup"
        class="collapse navbar-collapse justify-content-end"
      >
        <ul class="navbar-nav">
          <!-- join quiz button -->
          <li class="nav-item mb-1">
            <button
              class="btn p-2 border border-danger btn-light nav-link btn-link mx-2"
              @click="navigate('/join')"
            >
              Join Quiz
            </button>
          </li>
          <!-- create quiz button -->
          <li v-if="isKratosUser" class="nav-item mb-1">
            <button
              class="btn p-2 border border-danger btn-light nav-link btn-link mx-2"
              @click="navigate('/admin/quiz/list-quiz')"
            >
              Create Quiz
            </button>
          </li>
          <!-- Reports -->
          <li v-if="isKratosUser" class="nav-item mb-1">
            <button
              class="btn p-2 border bg-light-primary nav-link btn-link mx-2"
              @click="navigate('/admin/reports')"
            >
              Reports
            </button>
          </li>
          <!-- Login button -->
          <li class="nav-item mb-1">
            <button
              v-if="!isKratosUser"
              class="btn p-2 border bg-light-primary nav-link btn-link mx-2"
              @click="navigate('/account/login')"
            >
              Log in
            </button>
            <button
              v-else
              class="btn p-2 border bg-light-primary nav-link btn-link mx-2"
              @click="navigate('/admin')"
            >
              My Profile
            </button>
          </li>
          <!-- Register button -->
          <li class="nav-item mb-1">
            <button
              v-if="!isKratosUser"
              class="btn p-2 border bg-primary nav-link btn-link mx-2"
              @click="navigate('/account/register')"
            >
              Sign up
            </button>
          </li>
          <!-- Log out button -->
          <li v-if="isKratosUser" class="nav-item mb-1">
            <button
              class="btn p-2 border bg-primary nav-link btn-link mx-2"
              @click="handleLogout()"
            >
              Log Out
            </button>
          </li>
        </ul>
      </div>
    </nav>
  </div>
  <div
    v-if="activeSession"
    class="alert bg-danger text-black d-flex justify-content-center align-items-center"
    role="alert"
  >
    <div class="doodle">&#128641;</div>
    <span> Your quiz is still runnig... check it out! </span>
    <span
      type="button"
      class="text-end ml-2"
      @click="navigateTo(`/admin/arrange/${activeSession}`)"
    >
      <font-awesome-icon :icon="['fas', 'arrow-up-right-from-square']" />
    </span>
    <span class="text-end ml-5">
      <span>If you want to stop the running quiz please click here</span>

      <font-awesome-icon
        type="button"
        class="ml-2"
        :icon="['fas', 'ban']"
        @click="stopQuiz"
      />
    </span>
  </div>
</template>

<script setup>
const router = useRouter();
import { setUserDataStore } from "~~/composables/auth";
import { useUsersStore } from "~~/store/users";
const userData = useUsersStore();
const { getUserData } = userData;
import { useSessionStore } from "~~/store/session";
const sessionStore = useSessionStore();
const { getSession, setSession } = sessionStore;
const activeSession = computed(() => getSession());

const navigate = (url) => {
  router.push(url);
};

const isKratosUser = computed(() => {
  const user = getUserData();
  if (user && user?.role == "admin-user") {
    return true;
  }
  return false;
});

onMounted(() => {
  setUserDataStore();
});

const stopQuiz = () => {
  setSession(null);
  console.log("stop called");
};
</script>

<style scoped>
.logo {
  height: 35px;
  transform: scale(2);
  margin-top: 4px;
  margin-left: 45px;
}

@keyframes doodle-animation {
  0% {
    left: calc(100% + 50px);
  }

  100% {
    left: -50px;
  }
}

.doodle {
  position: absolute;
  animation: doodle-animation 5s linear infinite;
  font-size: 24px;
}
</style>
