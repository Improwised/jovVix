<script setup>
const url = useState("urls");
const dataForUsers = reactive([]);
const route = useRoute();
const activeQuizId = ref();

async function fetchScoreDataForUser() {
  const headers = useRequestHeaders(["cookie"]);

  const { data } = await useFetch(url.value.api_url + "/final_score/user", {
    method: "GET",
    headers: headers,
    credentials: "include",
    mode: "cors",
  });
  dataForUsers.push(...data.value.data);
}

onMounted(() => {
  activeQuizId.value = route.query.aqi ? route.query.aqi : "";
  if (activeQuizId.value != "") {
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
  </div>
</template>
