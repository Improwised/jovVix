<script setup>
import { useToast } from "vue-toastification";
const url = useState("urls");
const scoreboardData = reactive([]);
const route = useRoute();
const router = useRouter();
const activeQuizId = ref("");
const toast = useToast();
const app = useNuxtApp();
const headers = useRequestHeaders(["cookie"]);
useSystemEnv();

const props = defineProps({
  userURL: {
    default: "",
    type: String,
    required: true,
  },
  isAdmin: {
    default: false,
    type: Boolean,
    required: false,
  },
});

async function getFinalScoreboardDetails(endpoint) {
  const { data, error } = await useFetch(() => url.value.api_url + endpoint, {
    method: "GET",
    headers: headers,
    credentials: "include",
    mode: "cors",
  });

  watch(
    [data, error],
    () => {
      if (data.value) {
        scoreboardData.push(...data.value.data);
      }
      if (error.value) {
        toast.error(app.$Unauthorized);
        router.push("/");
      }
    },
    { immediate: true, deep: true }
  );
}

if (props.isAdmin) {
  activeQuizId.value = props.isAdmin ? route.query.aqi : "";
  getFinalScoreboardDetails(
    props.userURL + "?active_quiz_id=" + activeQuizId.value
  );
} else {
  getFinalScoreboardDetails(props.userURL);
}
</script>
<template>
  <ClientOnly>
    <div>
      <div v-if="scoreboardData" class="table-responsive mt-5 w-100">
        <table
          class="table align-middle"
          :class="{
            'table-dark': props.isAdmin,
            'table-light': !props.isAdmin,
            'table-borderless': !props.isAdmin,
            'table-striped': !props.isAdmin,
            'table-hover': !props.isAdmin,
          }"
        >
          <thead>
            <caption>
              Rankings
            </caption>
            <tr>
              <th>Rank</th>
              <th>User</th>
              <th>Score</th>
            </tr>
          </thead>
          <tbody class="table-group-divider">
            <tr v-for="(user, index) in scoreboardData" :key="index">
              <td>{{ user.rank }}</td>
              <td>{{ user.username }}</td>
              <td>{{ user.score }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </ClientOnly>
</template>
