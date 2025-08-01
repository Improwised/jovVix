<template>
  <div class="header">
    <!-- navbar -->
    <nav class="navbar navbar-expand-lg">
      <!-- Logo -->
      <NuxtLink class="navbar-brand navbar-logo" style="color: black" to="/">
        <img class="logo" src="/jovvix-logo.png" alt="" />
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
              class="btn px-4 bg-primary btn-light nav-link nav-link-button btn-link mx-2"
              @click="navigate('/join')"
            >
              Join Quiz
            </button>
          </li>
          <!-- create quiz button -->
          <li v-if="isKratosUser" class="nav-item mb-1">
            <button
              class="btn px-4 bg-light-primary btn-light nav-link nav-link-button btn-link mx-2"
              @click="navigate('/admin/quiz')"
            >
              Quizzes
            </button>
          </li>
          <!-- Reports -->
          <li v-if="isKratosUser" class="nav-item mb-1">
            <button
              class="btn px-4 border bg-light-primary nav-link nav-link-button btn-link mx-2"
              @click="navigate('/admin/reports')"
            >
              Reports
            </button>
          </li>
          <!-- Login button -->
          <li class="nav-item mb-1">
            <div v-if="!isKratosUser" @click="navigate('/account/login')">
              <button
                class="btn px-4 border bg-light-primary nav-link nav-link-button btn-link mx-2"
              >
                Log in
              </button>
            </div>
            <button
              v-else
              class="btn px-4 border nav-link nav-link-button btn-link mx-2"
              @click="navigate('/admin')"
            >
              My Profile
            </button>
          </li>
          <!-- Register button -->
          <li class="nav-item mb-1">
            <button
              v-if="!isKratosUser"
              class="btn px-4 border bg-primary nav-link nav-link-button btn-link mx-2"
              @click="navigate('/account/register')"
            >
              Sign up
            </button>
          </li>
          <!-- Log out button -->
          <li v-if="isKratosUser" class="nav-item mb-1">
            <button
              class="btn px-4 border nav-link nav-link-button btn-link mx-2"
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
    v-if="activeSession && isKratosUser"
    class="alert bg-danger text-black d-flex justify-content-center align-items-center"
    role="alert"
  >
    <div class="doodle">&#128641;</div>
    <div class="row w-100">
      <div class="col-6 text-end">
        <span>
          <span>Your quiz is still runnig... check it out!</span>
          <font-awesome-icon
            type="button"
            class="ml-2 scale-110"
            :icon="['fas', 'arrow-up-right-from-square']"
            @click="navigateTo(`/admin/arrange/${activeSession}`)"
          />
        </span>
      </div>
      <div class="col-6">
        <span class="ml-5">
          <span>If you want to stop the running quiz please click here</span>

          <font-awesome-icon
            type="button"
            class="ml-2 scale-110"
            :icon="['fas', 'ban']"
            @click="stopQuiz"
          />
        </span>
      </div>
    </div>
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
  height: 48px;
}

@keyframes doodle-animation {
  0% {
    left: calc(100% + 50px);
  }

  100% {
    left: -50px;
  }
}

.scale-110 {
  transform: scale(1.2);
}

.doodle {
  position: absolute;
  animation: doodle-animation 5s linear infinite;
  font-size: 24px;
}

.nav-link-button {
  border-radius: calc(0.625rem + 4px);
  font-weight: bold;
}
</style>
