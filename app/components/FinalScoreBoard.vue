<script setup>
import { useToast } from "vue-toastification";
import WinnerCard from "./WinnerCard.vue";
import ScoreBoardTable from "./ScoreBoardTable.vue";

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
const winnerUI = computed(() => route.query.winner_ui || false);

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

const changeUI = (value) => {
  navigateTo({ path: route.path, query: { ...route.query, winner_ui: value } });
};
</script>
<template>
  <div v-if="winnerUI == 'true' && props.isAdmin">
    <video id="myVideo" autoplay muted loop>
      <source src="@/assets/video/winner.mp4" type="video/mp4" />
      Your browser does not support HTML5 video.
    </video>
    <div v-if="requestPending" class="text-center">Loading...</div>
    <div
      v-else
      class="container-fluid justify-content-around row winners-container"
    >
      <div v-if="scoreboardData.length > 1" class="col-sm-12 col-lg-3 mt-4">
        <WinnerCard :winner="scoreboardData[1]" />
      </div>
      <div
        v-if="scoreboardData.length > 0"
        class="col-sm-12 col-lg-3 mt-4 first-rank"
      >
        <WinnerCard :winner="scoreboardData[0]" />
      </div>
      <div v-if="scoreboardData.length > 2" class="col-sm-12 col-lg-3 mt-4">
        <WinnerCard :winner="scoreboardData[2]" />
      </div>
      <div class="text-center change-ui-button" @click="changeUI(false)">
        <v-btn rounded color="light" dark x-large class="mt-3 px-7" flat
          >Next</v-btn
        >
      </div>
    </div>
  </div>

  <ClientOnly v-else>
    <div v-if="requestPending" class="text-center">Loading...</div>
    <div v-else>
      <div
        v-if="scoreboardData"
        class="table-responsive mt-5 w-100 container p-0 pb-2"
      >
        <ScoreBoardTable
          :scoreboard-data="scoreboardData"
          :is-admin="props.isAdmin"
          :user-name="props.userName"
        />
        <button
          v-if="props.isAdmin"
          class="btn btn-primary"
          @click="showAnalysis"
        >
          Show Analysis
        </button>
        <button
          v-if="props.isAdmin"
          class="btn m-2 btn-primary"
          @click="changeUI(true)"
        >
          Show Winners
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

<style scoped>
#myVideo {
  position: fixed;
  right: 0;
  bottom: 0;
  min-width: 100%;
  min-height: 100%;
}

.winners-container {
  margin-top: 10rem;
}

.winner-card {
  border-radius: 10px;
  box-shadow: 0px 4px 5px white;
  transition: transform 1s;
}

.winner-card:hover {
  transform: scale(1.1);
}

.first-rank {
  transform: scale(1.5);
}

.change-ui-button {
  margin-top: 8rem;
}

@media only screen and (max-width: 1079px) {
  .first-rank {
    transform: scale(1);
  }

  .change-ui-button {
    margin-top: 2rem;
  }

  .winners-container {
    margin-top: 0px;
  }
}
</style>
