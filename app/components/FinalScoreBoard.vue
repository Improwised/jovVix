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
const winningSound = ref(null);

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

watch(
  winnerUI,
  (newValue) => {
    const music = newValue == "true";
    if (!music && winningSound.value) {
      winningSound.value.pause();
    } else if (music && winningSound.value) {
      winningSound.value.play();
    }
  },
  { deep: true, immediate: true }
);

onMounted(() => {
  winningSound.value = new Audio("/music/winning.mp3");

  if (winnerUI.value == "true") {
    winningSound.value.play();
  }
});
</script>
<template>
  <div v-if="winnerUI == 'true' && props.isAdmin">
    <img id="myVideo" src="@/assets/images/medal/bg.jpg" alt="" />
    <div v-if="requestPending" class="text-center">Loading...</div>
    <div
      v-else
      class="container-fluid justify-content-around row winners-container"
    >
      <div
        v-if="scoreboardData.length > 0"
        class="col-sm-12 col-lg-3 order-sm-1 order-lg-2 rank-one"
      >
        <WinnerCard :winner="scoreboardData[0]" />
      </div>
      <div
        v-if="scoreboardData.length > 1"
        class="col-sm-12 col-lg-3 order-sm-2 order-lg-1 rank-two"
      >
        <WinnerCard :winner="scoreboardData[1]" />
      </div>
      <div
        v-if="scoreboardData.length > 2"
        class="col-sm-12 col-lg-3 order-sm-3 order-lg-3"
      >
        <WinnerCard :winner="scoreboardData[2]" />
      </div>
      <div
        class="col-12 order-4 text-center change-ui-button"
        @click="changeUI(false)"
      >
        <v-btn rounded color="light" dark x-large class="px-7" flat>Next</v-btn>
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
        <QuizStatisticsBadges :user-statistics="userStatistics" />
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
  width: 100%;
  height: auto;
}

.rank-one {
  transform: scale(1.25);
}

.rank-two {
  transform: scale(1.1);
}

@media only screen and (max-width: 1079px) {
  .change-ui-button {
    margin-top: 2rem;
  }

  .winners-container {
    margin-top: 0px;
  }
}
</style>
