<template>
  <div class="container">
    <!-- Modal -->
    <div
      id="exampleModal"
      class="modal fade"
      tabindex="-1"
      aria-labelledby="exampleModalLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h1
              id="exampleModalLabel"
              class="modal-title fs-5 text-center text-white"
            >
              Change Password
            </h1>
            <button
              type="button"
              class="btn-close"
              data-bs-dismiss="modal"
              aria-label="Close"
            ></button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label for="name" class="form-label">Enter new password</label>
              <input
                id="name"
                v-model="password.new"
                type="password"
                class="form-control"
                placeholder="Enter new password"
                required
              />
            </div>
            <div class="form-group">
              <label for="name" class="form-label">Confirm New Password</label>
              <input
                id="name"
                v-model="password.confirm"
                type="password"
                class="form-control"
                placeholder="Confirm new password"
                required
              />
              <div v-if="passwordRequestError" class="form-text text-danger">
                {{ passwordRequestError }}
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button
              id="closeModalButton"
              type="button"
              class="btn btn-secondary text-white"
              data-bs-dismiss="modal"
            >
              Close
            </button>
            <button
              type="button"
              class="btn btn-primary text-white"
              @click="changePassword"
            >
              Change Password
            </button>
          </div>
        </div>
      </div>
    </div>
    <h3 class="d-flex align-item-center justify-content-center">
      Account Profile
    </h3>
    <div v-if="userError?.data?.code == 401">
      {{ navigateTo("/account/login") }}
    </div>
    <div v-if="userError" class="alert alert-danger" role="alert">
      {{ userError.data }}
    </div>
    <div v-else-if="userPending" class="text-center">Pending...</div>
    <div v-else class="row">
      <div class="col-md-5 mt-4">
        <div class="card d-flex justify-content-center align-items-center">
          <img
            :src="avatar"
            class="card-img-top mt-3"
            style="width: 14rem"
            alt="..."
          />
          <div class="card-body text-center">
            <h5 class="card-title">{{ profile.full_name }}</h5>
            <h5 class="card-title">{{ profile.email }}</h5>
            <div
              type="button"
              class="btn btn-primary btn-sm text-white"
              data-bs-toggle="modal"
              data-bs-target="#exampleModal"
            >
              Change Password
            </div>
          </div>
        </div>
      </div>
      <div class="col-md-7 mt-4">
        <div class="card">
          <div class="card-body">
            <form @submit.prevent="changeUserMetaData">
              <div class="form-group">
                <label for="name" class="form-label">First Name</label>
                <input
                  id="name"
                  v-model="userData.first_name"
                  type="text"
                  class="form-control"
                  placeholder="Pending..."
                  required
                  @focus="showCancleButton"
                />
              </div>
              <div class="form-group">
                <label for="name" class="form-label">Last Name</label>
                <input
                  id="name"
                  v-model="userData.last_name"
                  type="text"
                  class="form-control"
                  placeholder="Pending..."
                  required
                  @focus="showCancleButton"
                />
              </div>
              <div class="form-group">
                <label for="email" class="form-label">Email</label>
                <input
                  id="email"
                  v-model="userData.email"
                  type="email"
                  class="form-control"
                  placeholder="Pending..."
                  required
                  @focus="showCancleButton"
                />
              </div>
              <!-- Error message for password mismatch -->
              <div
                v-if="updateuserError"
                class="alert alert-danger"
                role="alert"
              >
                {{ updateuserError.data }}
              </div>
              <div v-else-if="updateuserPending">
                <button class="btn btn-primary">Pending...</button>
              </div>
              <div v-else class="form-group">
                <button type="submit" class="btn btn-primary text-white">
                  Save Changes
                </button>
                <div
                  v-if="cancleButton"
                  class="btn btn-secondary ml-2"
                  @click="hideCancleButton"
                >
                  Cancel
                </div>
              </div>
            </form>
          </div>
        </div>
      </div>
      <div class="d-flex flex-column justify-content-center">
        <!-- list loader -->
        <UtilsQuizListWaiting v-if="quizPending" />

        <div v-else-if="quizError" class="alert alert-danger" role="alert">
          {{ quizError.data }}
        </div>

        <!-- quiz details -->
        <div v-else>
          <!-- show quiz list -->
          <div class="mt-4">
            <div class="card">
              <div class="card-header">
                <nav class="navbar">
                  <div class="container-fluid p-0">
                    <h3 class="mb-0">Played Quiz List</h3>
                    <NuxtLink
                      class="btn text-white btn-primary"
                      to="/admin/played_quiz"
                    >
                      Played Quiz
                    </NuxtLink>
                  </div>
                </nav>
              </div>
              <div
                v-if="quizList == null || quizList.length < 1"
                class="no-quiz-list d-flex flex-column align-items-center my-4"
              >
                <h2>No Quiz Played By You !</h2>
              </div>
              <div class="row p-2">
                <div
                  v-for="(details, index) in quizList"
                  :key="index"
                  class="col-md-6 mb-5"
                >
                  <QuizListCard :details="details" :is-played-quiz="true" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useUsersStore } from "~~/store/users";
import { getAvatarUrlByName } from "~~/composables/avatar";
const userStore = useUsersStore();
const { getUserData } = userStore;
const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const updateuserError = ref(false);
const updateuserPending = ref(false);
const passwordRequestError = ref(false);
const cancleButton = ref(false);

const avatar = computed(() => {
  const user = getUserData();
  return getAvatarUrlByName(user?.avatar);
});

const {
  data: user,
  pending: userPending,
  error: userError,
} = useFetch(url.api_url + "/kratos/whoami", {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
});

const {
  data: quizList,
  pending: quizPending,
  error: quizError,
} = useFetch(url.api_url + "/user_played_quizes", {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
  transform: (quizList) => {
    if (quizList.data) {
      return quizList?.data?.data.slice(0, 8);
    }
    return quizList?.data?.data;
  },
});

const userData = reactive({
  first_name: "",
  last_name: "",
  email: "",
});

const profile = reactive({
  full_name: "",
  email: "",
});

const password = reactive({
  new: "",
  confirm: "",
});

watch(
  [user, userError],
  () => {
    if (user.value) {
      profile.full_name =
        user.value.data.identity.traits.name.first +
        " " +
        user.value.data.identity.traits.name.last;
      userData.first_name = user.value.data.identity.traits.name.first;
      userData.last_name = user.value.data.identity.traits.name.last;
      userData.email = user.value.data.identity.traits.email;
      profile.email = user.value.data.identity.traits.email;
    }
    if (userError.value) {
      console.log("error");
    }
  },
  { immediate: true, deep: true }
);

const changeUserMetaData = async () => {
  updateuserPending.value = true;
  const { data, error } = await useFetch(url.api_url + "/kratos/user", {
    method: "PUT",
    headers: headers,
    mode: "cors",
    credentials: "include",
    body: {
      first_name: userData.first_name,
      last_name: userData.last_name,
      email: userData.email,
    },
  });

  updateuserPending.value = false;
  if (error.value) {
    updateuserError.value = error.value.data;
    setTimeout(() => {
      updateuserError.value = false;
    }, 2000);
  } else {
    profile.full_name =
      data.value.data.first_name + " " + data.value.data.last_name;
    profile.email = data.value.data.email;
    cancleButton.value = false;
  }
};

const showCancleButton = () => {
  cancleButton.value = true;
};

const hideCancleButton = () => {
  userData.first_name = user.value.data.identity.traits.name.first;
  userData.last_name = user.value.data.identity.traits.name.last;
  userData.email = user.value.data.identity.traits.email;
  cancleButton.value = false;
};

const flow = ref("");
const csrfToken = ref("");

const fetchFlowIdAndCsrfToken = async () => {
  try {
    const response = await fetch(
      `${url.kratos_url}/self-service/settings/browser?aal=&refresh=&return_to=`,
      {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
        credentials: "include",
      }
    );

    if (!response.ok) {
      throw new Error(`Error: ${response.statusText}`);
    }

    const data = await response.json();
    flow.value = data.id;
    csrfToken.value = data.ui.nodes.find(
      (node) => node.attributes.name === "csrf_token"
    ).attributes.value;
  } catch (error) {
    console.error("Error fetching flow ID and CSRF token:", error);
  }
};

const changePassword = async () => {
  try {
    if (password.new !== password.confirm) {
      throw new Error("Passwords do not match.");
    }

    await fetchFlowIdAndCsrfToken();

    const response = await fetch(
      `${url.kratos_url}/self-service/settings?flow=${flow.value}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          password: password.new,
          csrf_token: csrfToken.value,
          method: "password",
        }),
      }
    );

    if (response.status == 400) {
      throw new Error("the password is too similar to the user identifier");
    }

    if (!response.ok) {
      const errorData = await response.json();
      if (errorData.error.id === "session_refresh_required") {
        window.location.href = errorData.redirect_browser_to;
      } else {
        throw new Error(errorData.error.message);
      }
    }

    passwordRequestError.value = "";

    var closeModalButton = document.getElementById("closeModalButton");
    closeModalButton.click();
  } catch (error) {
    console.error("Error during password change:", error.message);
    passwordRequestError.value = error.message;
  }
};
</script>

<style scoped>
.profile-page {
  padding: 20px;
}

.user-info {
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 20px;
}
</style>
