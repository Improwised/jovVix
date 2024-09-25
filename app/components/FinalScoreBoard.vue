<script setup>
import { useToast } from "vue-toastification";

const url = useRuntimeConfig().public;
const scoreboardData = ref([]);
const route = useRoute();
const router = useRouter();
const activeQuizId = ref("");
const toast = useToast();
const app = useNuxtApp();
const headers = useRequestHeaders(["cookie"]);
useSystemEnv();

const analysisData = reactive([]);
const userAnalysisEndpoint = "/analytics_board/user";
const requestPending = ref(false);
const userStatistics = ref({});

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
  userName: {
    type: String,
    required: false,
    default: "",
  },
  userPlayedQuiz: {
    type: String,
    required: false,
    default: "",
  },
});

const getFinalScoreboardDetails = async (endpoint) => {
  try {
    requestPending.value = true;
    await $fetch(`${url.api_url}${endpoint}`, {
      method: "GET",
      credentials: "include",
      headers: {
        Accept: "application/json",
      },
      onResponse({ response }) {
        if (response.status != 200) {
          requestPending.value = false;
          toast.error("error while get final scorcard");
          return;
        }
        if (response.status == 200) {
          scoreboardData.value = response._data?.data;
          requestPending.value = false;
        }
      },
    });
  } catch (error) {
    toast.error(error.message);
    requestPending.value = false;
    return;
  }
};

// for admin
if (props.isAdmin) {
  activeQuizId.value = props.isAdmin ? route.query.aqi : "";
  getFinalScoreboardDetails(
    props.userURL + "?active_quiz_id=" + activeQuizId.value
  );
} else {
  getFinalScoreboardDetails(
    `${props.userURL}?user_played_quiz=${props.userPlayedQuiz}`
  );
}

//for users
if (!props.isAdmin) {
  async function getAnalysisDetails() {
    const { data, error } = await useFetch(
      () =>
        `${url.api_url}${userAnalysisEndpoint}?user_played_quiz=${props.userPlayedQuiz}`,
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
          analysisData.push(...data.value.data);
          userStatistics.value = questionsAnalysis(data.value?.data);
        }
        if (error.value) {
          toast.error(app.$$Unauthorized);
        }
      },
      { immediate: true, deep: true }
    );
  }

  getAnalysisDetails();
}

const showAnalysis = () => {
  router.push({
    path: `/admin/reports/${activeQuizId.value}`,
  });
};
</script>
<template>
  <ClientOnly>
    <div v-if="requestPending" class="text-center">Loading...</div>
    <div v-else>
      <div v-if="scoreboardData" class="table-responsive mt-5 w-100 container">
        <table class="table align-middle table-bordered">
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
            <tr
              v-for="(user, index) in scoreboardData"
              :key="index"
              :class="{
                'table-primary':
                  user.username === props.userName && !props.isAdmin,
              }"
            >
              <td>
                {{ user.rank }}
              </td>
              <td v-if="props.isAdmin">
                {{ user.firstname }} <span>({{ user.username }})</span>
              </td>
              <td v-else>
                {{ user.firstname }}
                <span v-if="props?.userName === user.username">
                  &nbsp;({{ user.username }})
                </span>
              </td>
              <td>{{ user.score }}</td>
            </tr>
          </tbody>
        </table>
        <button
          v-if="props.isAdmin"
          class="btn btn-primary"
          @click="showAnalysis"
        >
          Show Analysis
        </button>
      </div>
      <div v-if="!props.isAdmin">
        <div
          class="d-flex flex-wrap align-items-center justify-content-around container bg-white rounded p-5"
        >
          <span
            class="badge rounded-pill bg-light-primary text-dark m-2 px-2 fs-5"
          >
            Accuracy: {{ userStatistics?.accuracy }}%
          </span>
          <span
            class="badge rounded-pill bg-light-primary text-dark m-2 px-2 fs-5"
          >
            Total Score: {{ userStatistics?.totalScore }}
          </span>
          <span
            class="badge rounded-pill bg-light-primary text-dark m-2 px-2 fs-5"
          >
            Total Correct: {{ userStatistics?.correctAnwers }}</span
          >
          <span
            class="badge rounded-pill bg-light-secondary text-dark m-2 px-2 fs-5"
          >
            Total Incorrect: {{ userStatistics?.wrongAnwers }}
          </span>
          <span
            class="badge rounded-pill bg-light-secondary text-dark m-2 px-2 fs-5"
          >
            Total Un-attmpted: {{ userStatistics?.unAttemptedQuestions }}
          </span>
        </div>
        <QuizAnalysis :data="analysisData" />
      </div>
    </div>
  </ClientOnly>
</template>
