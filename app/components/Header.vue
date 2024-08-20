<template>
  <div class="header">
    <!-- navbar -->
    <nav class="navbar-classic navbar navbar-expand-lg">
      <!-- Logo -->
      <NuxtLink class="navbar-brand navbar-logo" style="color: black" to="/">
        Quiz App
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
          <!-- Reports -->
          <li v-if="isKratosUser" class="nav-item mb-1">
            <button
              class="btn custom-btn nav-link btn-link rounded-pill"
              @click="navigate('/admin/reports')"
            >
              Reports
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
    </nav>
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
.custom-btn {
  color: #fff;
  background: linear-gradient(270deg, #5a66ef 0, #8042e4);
  border-color: #007bff;
  margin-right: 10px; /* Space between buttons */
  padding: 8px 16px; /* Adjust padding as needed */
  border-radius: 20px; /* Increase border-radius for a pill shape */
  text-align: center; /* Center text */
}

.custom-btn:hover {
  color: #fff;
  background-color: #007bff;
  border-color: #007bff;
}
</style>
