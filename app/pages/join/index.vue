<template>
  <div class="join-page-container">
    <div class="content-container">
      <div class="h-100 d-flex align-items-center justify-content-center">
        <div class="row w-100">
          <div class="col-12">
            <QuizLoadingSpace v-if="pageLoading"></QuizLoadingSpace>
            <Frame
              v-else
              page-title="Join Quiz"
              page-message="Let's Play Together"
              class="bg-white"
            >
              <div v-if="userError">
                {{ userError }}
              </div>
              <form v-else method="POST" @submit.prevent="join_quiz">
                <div class="mb-3 pe-3">
                  <label
                    for="code"
                    class="form-label text-primary font-weight-bold"
                    >Invitation Code</label
                  >
                  <!-- Assuming v-otp-input is a custom component or external library -->
                  <v-otp-input
                    v-model="code"
                    max-width="500"
                    min-height="20"
                    type="number"
                    placeholder="0"
                    class="text-primary font-weight-bold"
                  />
                </div>
                <div class="mb-3">
                  <label
                    v-if="!isUserLoggedIn"
                    for="username"
                    class="form-label text-primary font-weight-bold"
                    >Username</label
                  >
                  <input
                    v-if="!isUserLoggedIn"
                    id="username"
                    v-model.trim="username"
                    type="text"
                    name="username"
                    class="text-primary font-weight-bold form-control"
                  />
                  <div v-if="isUserLoggedIn">
                    Welcome
                    <span class="font-weight-bold">{{ firstname }}</span>
                  </div>
                </div>
                <div class="p-2">
                  <div v-if="!isUserLoggedIn" class="text-center">
                    Want To Save Your Progress?
                    <NuxtLink to="/account/login"><b>Login</b></NuxtLink> Now.
                  </div>
                </div>
                <button
                  v-if="quickUserPending"
                  class="btn btn-primary btn-lg w-100 text-white join-button"
                >
                  Pending...
                </button>
                <button
                  v-else
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
import { useUsersStore } from "~~/store/users";
import { getRandomAvatarName } from "~~/composables/avatar";
const userData = useUsersStore();
const { setUserData } = userData;
const pageLoading = ref(true);
const route = useRoute();

// Assuming these are imported from external libraries or mixins
const codeparam = computed(() => route.query.code || "");
const code = ref(codeparam.value);
const username = ref("");
const firstname = ref({});
const isUserLoggedIn = ref(false);
const { apiUrl } = useRuntimeConfig().public;
const router = useRouter();
const toast = useToast();
const userError = ref(false);
const quickUserPending = ref(false);
const userPlayedQuiz = ref("");
const sessionId = ref("");

const join_quiz = async () => {
  username.value = username.value.trim().replace(/\s+/g, "");

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

  // create quick user
  if (!isUserLoggedIn.value) {
    quickUserPending.value = true;
    const avatarName = getRandomAvatarName();

    try {
      await $fetch(
        `${apiUrl}/user/${username.value}?avatar_name=${avatarName}`,
        {
          method: "POST",
          credentials: "include",
          headers: {
            Accept: "application/json",
          },
          onResponse({ response }) {
            if (response.status == 200) {
              isUserLoggedIn.value = true;
              quickUserPending.value = false;
              firstname.value = response._data?.data?.first_name;
              username.value = response._data?.data?.username;
              setUserData({
                role: "guest-user",
                avatar: response._data?.data?.img_key,
              });
            }
          },
        }
      );
    } catch (error) {
      userError.value = error.message;
      quickUserPending.value = false;
      return;
    }
  }

  try {
    quickUserPending.value = true;
    await $fetch(`${apiUrl}/user_played_quizes/${code.value}`, {
      method: "POST",
      credentials: "include",
      headers: {
        Accept: "application/json",
      },
      onResponse({ response }) {
        if (response.status == 200) {
          userPlayedQuiz.value = response._data?.data?.user_played_quiz;
          sessionId.value = response._data?.data?.session_id;
          quickUserPending.value = false;
        }
      },
    });
  } catch (error) {
    userError.value = error.message;
    quickUserPending.value = false;
    if (error?.status == 400) {
      userError.value = "invitation code not found";
    }
    return;
  }

  router.push(
    `/join/play/${code.value}?username=${encodeURIComponent(
      username.value
    )}&firstname=${firstname.value}&user_played_quiz=${
      userPlayedQuiz.value
    }&session_id=${sessionId.value}`
  );
};

// get user data
(async () => {
  try {
    await $fetch(apiUrl + "/user/who", {
      method: "GET",
      credentials: "include",
      headers: {
        Accept: "application/json",
      },
      onResponse({ response }) {
        if (response.status == 200) {
          isUserLoggedIn.value = true;
          firstname.value = response._data?.data?.firstname;
          username.value = response._data?.data?.username;
          pageLoading.value = false;
        }
      },
    });
  } catch (error) {
    if (error.status == 401) {
      isUserLoggedIn.value = false;
      pageLoading.value = false;
      return;
    }
    userError.value = error.message;
    pageLoading.value = false;
  }
})();
</script>

<style scoped>
.join-button {
  background: linear-gradient(270deg, #5a66ef 0, #8042e4);
}

.join-button:hover {
  background: #6f4eb8;
  /* Slightly darker shade */
}
</style>
