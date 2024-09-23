<script setup>
import { isCorrectAnswer } from "~/composables/check_is_correct.js/";
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

const userAccuracy = ref(0);
const userTotalScore = ref(0);
const requestPending = ref(false);

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
          let analysis = questionsAnalysis(data.value?.data);
          userAccuracy.value = analysis.accuracy;
          userTotalScore.value = analysis.totalScore;
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
    <div>
      <div v-if="scoreboardData" class="table-responsive mt-5 w-100">
        <table class="table align-middle table-light">
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
              <td
                :class="{
                  'user-row':
                    user.username === props.userName && !props.isAdmin,
                }"
              >
                {{ user.rank }}
              </td>
              <td v-if="props.isAdmin">
                {{ user.firstname }} <span>({{ user.username }})</span>
              </td>
              <td
                v-else
                :class="{ 'user-row': user.username === props.userName }"
              >
                {{ user.firstname }}
                <span v-if="props?.userName === user.username">
                  &nbsp;({{ user.username }})
                </span>
              </td>
              <td
                :class="{
                  'user-row':
                    user.username === props.userName && !props.isAdmin,
                }"
              >
                {{ user.score }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-if="!props.isAdmin">
        <h3 class="text-center">Accuracy: {{ userAccuracy }}%</h3>
        <h3 class="text-center">Total Score: {{ userTotalScore }}</h3>
        <QuizQuestionAnalysis :data="analysisData" />
      </div>
    </div>
    <button v-if="props.isAdmin" class="btn btn-primary" @click="showAnalysis">
      Show Analysis
    </button>
  </ClientOnly>
</template>

<style scoped>
.user-row {
  background-color: #8968cd !important;
  color: white;
  box-shadow: none;
}
</style>
