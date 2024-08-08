<template>
  <div
    class="container-fluid d-flex align-items-center min-vh-100 position-relative"
  >
    <header class="d-flex header bg-light position-fixed start-0 end-0 top-0">
      <nav class="navbar navbar-expand-md navbar-dark w-100">
        <div class="container-fluid p-2">
          <!-- Title on the left -->
          <NuxtLink
            class="navbar-brand navbar-logo"
            style="color: black"
            to="/"
          >
            Quiz App
          </NuxtLink>
          <div>
            <!-- Join Quiz button (visible on mobile) -->
            <button
              class="navbar-brand align-items-center d-md-none btn btn-link custom-btn rounded-pill"
              @click="navigate('/join')"
            >
              Join Quiz
            </button>
            <!-- Navbar toggler for mobile -->
            <button
              class="navbar-toggler bg-dark"
              type="button"
              data-bs-toggle="collapse"
              data-bs-target="#navbarNavAltMarkup"
              aria-controls="navbarNavAltMarkup"
              aria-expanded="false"
              aria-label="Toggle navigation"
            >
              <span class="navbar-toggler-icon text-dark"></span>
            </button>
          </div>

          <!-- Profile/Login and Admin buttons -->
          <div
            id="navbarNavAltMarkup"
            class="collapse navbar-collapse justify-content-end"
          >
            <ul class="navbar-nav">
              <!-- Register button -->
              <li class="nav-item mb-1">
                <button
                  v-if="!isKratosUser"
                  class="btn custom-btn nav-link btn-link rounded-pill"
                  @click="navigate('/account/register')"
                >
                  Register
                </button>
              </li>
              <!-- Login button -->
              <li class="nav-item mb-1">
                <button
                  v-if="!isKratosUser"
                  class="btn custom-btn nav-link btn-link rounded-pill"
                  @click="navigate('/account/login')"
                >
                  Login
                </button>
                <button
                  v-else
                  class="btn custom-btn nav-link btn-link rounded-pill"
                  @click="navigate('/admin')"
                >
                  My Profile
                </button>
              </li>
              <!-- create quiz button -->
              <li v-if="isKratosUser" class="nav-item mb-1">
                <button
                  class="btn custom-btn nav-link btn-link rounded-pill"
                  @click="navigate('/admin/quiz/list-quiz')"
                >
                  Create Quiz
                </button>
              </li>
              <!-- join quiz button -->
              <li class="nav-item mb-1">
                <button
                  class="btn custom-btn nav-link btn-link rounded-pill"
                  @click="navigate('/join')"
                >
                  Join Quiz
                </button>
              </li>
              <!-- Log out button -->
              <li v-if="isKratosUser" class="nav-item mb-1">
                <button
                  class="btn custom-btn nav-link btn-link rounded-pill"
                  @click="handleLogout()"
                >
                  Log Out
                </button>
              </li>
            </ul>
          </div>
        </div>
      </nav>
    </header>
    <main
      class="container-md d-flex justify-content-center align-items-center h-100 pageContainer px-0"
    >
      <slot />
    </main>
    <footer
      class="d-flex justify-content-center position-absolute start-0 end-0 bottom-0"
    >
      <div class="text-center text-capitalize p-3 text-white w-100">
        © {{ new Date().getFullYear() }} Copyright:
        <NuxtLink class="text-light" to="/">Quizz_app</NuxtLink>
      </div>
    </footer>
  </div>
</template>

<script setup>
const router = useRouter();
import { useUsersStore } from "~~/store/users";
const userData = useUsersStore();
const { getUserData } = userData;

const navigate = (url) => {
  router.push(url);
};

const isKratosUser = computed(() => {
  const user = getUserData();
  if (user && user == "admin-user") {
    return true;
  }
  return false;
});
</script>

<style scoped>
header,
footer {
  z-index: 0;
  background: #290f5a;
}

.pageContainer {
  padding-top: 100px;
  padding-bottom: 100px;
}

.custom-btn {
  color: #fff;
  background: linear-gradient(270deg, #5a66ef 0, #8042e4);
  border-color: #007bff;
  margin-right: 10px;
  /* Space between buttons */
  padding: 8px 16px;
  /* Adjust padding as needed */
  border-radius: 20px;
  /* Increase border-radius for a pill shape */
  text-align: center;
  /* Center text */
}

.custom-btn:hover {
  color: #fff;
  background-color: #007bff;
  border-color: #007bff;
}

/* Mobile view adjustments */
@media (max-width: 767px) {
  .navbar-brand {
    color: black;
  }

  .custom-btn {
    color: white !important;
    /* Ensure components are black on mobile */
  }

  .navbar-toggler-icon {
    color: black;
    /* Set navbar toggler icon color to black */
  }
}
</style>
