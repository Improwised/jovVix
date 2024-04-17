<script setup>
const url = useState("urls");
const dataForUsers = reactive([]);
const dataForAdmin = reactive([]);
const route = useRoute();
const activeQuizId = ref("");
useSystemEnv();

async function fetchScoreDataForUser() {
  const headers = useRequestHeaders(["cookie"]);
  const { data } = await useFetch(
    () => url.value.api_url + "/final_score/user",
    {
      method: "GET",
      headers: headers,
      credentials: "include",
      mode: "cors",
    }
  );
  watch(
    data,
    () => {
      if (data.value) {
        dataForUsers.push(...data.value.data);
        console.log("pending for user ", data.value);
      }
    },
    { immediate: true, deep: true }
  );
}

async function fetchScoreDataForAdmin() {
  const { data } = await useFetch(
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

  console.log("pending for admin ", data.value);
  watch(
    data,
    () => {
      if (data.value) {
        dataForAdmin.push(...data.value.data);
        console.log("pending for admin ", data.value);
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
          <th>Response Time</th>
        </tr>
      </thead>
      <tbody class="table-group-divider">
        <tr v-for="(user, index) in dataForUsers" :key="index">
          <td></td>
          <td>{{ user.userName }}</td>
          <td>{{ user.score }}</td>
          <td></td>
        </tr>
      </tbody>
      <tfoot></tfoot>
    </table>

    <!-- table for showing to Admin-->
    <table v-if="dataForAdmin.length" class="table table-danger align-middle">
      <thead>
        <caption>
          Rankings
        </caption>
        <tr>
          <th>Rank1</th>
          <th>User</th>
          <th>Score</th>
          <th>Response Time</th>
        </tr>
      </thead>
      <tbody class="table-group-divider">
        <tr v-for="(user, index) in dataForAdmin" :key="index">
          <td></td>
          <td>{{ user.userName }}</td>
          <td>{{ user.score }}</td>
          <td></td>
        </tr>
      </tbody>
      <tfoot></tfoot>
    </table>
  </div>
</template>
