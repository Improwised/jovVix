<script setup>
import { useToast } from "vue-toastification";

const url = useState("urls");
const dataForUsers = reactive([]);
const dataForAdmin = reactive([]);
const route = useRoute();
const router = useRouter();
const activeQuizId = ref("");
const toast = useToast();
const app = useNuxtApp();
useSystemEnv();

async function fetchScoreDataForUser() {
  const headers = useRequestHeaders(["cookie"]);
  const { data, error } = await useFetch(
    () => url.value.api_url + "/final_score/user",
    {
      method: "GET",
      headers: headers,
      credentials: "include",
      mode: "cors",
    }
  );

  watch(
    [data, error],
    () => {
      if (data.value) {
        dataForUsers.push(...data.value.data);
      }
      if (error.value) {
        toast.error(app.$Unauthorized);
        router.push("/");
      }
    },
    { immediate: true, deep: true }
  );
}

async function fetchScoreDataForAdmin() {
  const { data, error } = await useFetch(
    () =>
      url.value.api_url +
      "/final_score/admin?active_quiz_id=" +
      activeQuizId.value,
    {
      method: "GET",
      credentials: "include",
      mode: "cors",
    }
  );

  watch(
    [data, error],
    () => {
      if (data.value) {
        dataForAdmin.push(...data.value.data);
      }
      if (error.value) {
        toast.error(app.$Unauthorized);
        router.push("/");
      }
    },
    { immediate: true, deep: true }
  );
}

// on mounted check to whom show scoreboard data user or admin based on active quiz id
onBeforeMount(() => {
  activeQuizId.value = route.query.aqi ? route.query.aqi : "";
  if (activeQuizId.value != "") {
    fetchScoreDataForAdmin();
  } else {
    fetchScoreDataForUser();
  }
});
</script>

<template>
  <div class="table-responsive mt-5">
    <!-- table for showing to user -->
    <table
      v-if="dataForUsers.length"
      class="table table-striped table-hover table-borderless table-light align-middle"
    >
      <thead class="table-light">
        <caption>
          Rankings
        </caption>
        <tr>
          <th>Rank</th>
          <th>User</th>
          <th>Score</th>
          <th>Response Time (ms)</th>
        </tr>
      </thead>
      <tbody class="table-group-divider">
        <tr v-for="(user, index) in dataForUsers" :key="index">
          <td></td>
          <td>{{ user.username }}</td>
          <td>{{ user.score }}</td>
          <td>{{ user.response_time }}</td>
        </tr>
      </tbody>
      <tfoot></tfoot>
    </table>

    <!-- table for showing to Admin-->
    <table v-if="dataForAdmin.length" class="table table-dark align-middle">
      <thead>
        <caption>
          Rankings
        </caption>
        <tr>
          <th>Rank</th>
          <th>User</th>
          <th>Score</th>
          <th>Response Time (ms)</th>
        </tr>
      </thead>
      <tbody class="table-group-divider">
        <tr v-for="(user, index) in dataForAdmin" :key="index">
          <td></td>
          <td>{{ user.username }}</td>
          <td>{{ user.score }}</td>
          <td>{{ user.response_time }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
