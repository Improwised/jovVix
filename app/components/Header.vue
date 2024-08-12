<template>
  <div class="header">
    <!-- navbar -->
    <nav class="navbar-classic navbar navbar-expand-lg">
      <a id="nav-toggle" href="#"
        ><i data-feather="menu" class="nav-icon me-2 icon-xs"></i
      ></a>
      <!-- Logo -->
      <NuxtLink class="navbar-brand navbar-logo" style="color: black" to="/">
        Quiz App
      </NuxtLink>
      <!--Navbar nav -->
      <ul class="navbar-nav navbar-right-wrap ms-auto d-flex nav-top-wrap">
        <button
          class="navbar-brand align-items-center d-md-none btn btn-link custom-btn rounded-pill"
          @click="navigate('/join')"
        >
          Join Quiz
        </button>
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
