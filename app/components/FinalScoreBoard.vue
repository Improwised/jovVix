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
    await $fetch(`${url.apiUrl}${endpoint}`, {
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
        `${url.apiUrl}${userAnalysisEndpoint}?user_played_quiz=${props.userPlayedQuiz}`,
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
  if (process.client) {
    winningSound.value = new Audio("/music/winning.mp3");

    if (winnerUI.value == "true") {
      winningSound.value.play();
    }
  }
});
</script>
<template>
  <div v-if="winnerUI == 'true' && props.isAdmin">
    <img id="myVideo" src="@/assets/images/medal/bg.jpg" alt="" />
    <div v-if="requestPending" class="text-center">Loading...</div>
    <div v-else class="container-fluid justify-content-around row winners-container">
      <div class="podium-row">

        <div v-if="scoreboardData.length > 3" class="runner-up">
          <div class="podium-position">{{ scoreboardData[3].rank }}</div>
          <WinnerCard :winner="scoreboardData[3]" class="rank-four" />
        </div>
        <!-- Second Place -->
        <div v-if="scoreboardData.length > 1" class="podium-item rank-two">
          <div class="podium-position">{{ scoreboardData[1].rank }}</div>
          <WinnerCard :winner="scoreboardData[1]" />
        </div>

        <!-- First Place -->
        <div v-if="scoreboardData.length > 0" class="podium-item rank-one">
          <div class="podium-position">{{ scoreboardData[0].rank }}</div>
          <WinnerCard :winner="scoreboardData[0]" />
        </div>

        <!-- Third Place -->
        <div v-if="scoreboardData.length > 2" class="podium-item rank-three">
          <div class="podium-position">{{ scoreboardData[2].rank }}</div>
          <WinnerCard :winner="scoreboardData[2]" />
        </div>


        <div v-if="scoreboardData.length > 4" class="runner-up">
          <div class="podium-position">{{ scoreboardData[4].rank }}</div>
          <WinnerCard :winner="scoreboardData[4]" class="rank-five" />
        </div>
      </div>

      <div class="col-12 order-4 text-center change-ui-button" @click="changeUI(false)">
        <v-btn rounded color="light" dark x-large class="px-7" flat>Next</v-btn>
      </div>
    </div>
  </div>

  <ClientOnly v-else>
    <div v-if="requestPending" class="text-center">Loading...</div>
    <div v-else>
      <div v-if="scoreboardData" class="table-responsive mt-5 w-100 container p-0 pb-2">
        <ScoreBoardTable :scoreboard-data="scoreboardData" :is-admin="props.isAdmin" :user-name="props.userName" />
        <button v-if="props.isAdmin" class="btn btn-primary" @click="showAnalysis">
          Show Analysis
        </button>
        <button v-if="props.isAdmin" class="btn m-2 btn-primary" @click="changeUI(true)">
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

.winners-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 2rem 1rem;
  position: relative;
}

.podium-row {
  display: flex;
  justify-content: space-around;
  align-items: flex-end;
  gap: 1.5rem;
  width: 100%;
  margin-bottom: 2rem;
  height: 50%;
}

.podium-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  transition: all 0.3s ease;
  margin: 0 2rem;
}

.rank-one {
  transform: scale(1.25);
  z-index: 5;
}

.rank-two {
  transform: scale(1.1);
  z-index: 4;

}

.rank-three {
  transform: scale(1.1);
  z-index: 4;

}

.rank-four {
  transform: scale(0.90);
  z-index: 3;
}

.rank-five {
  transform: scale(0.90);
  z-index: 3;
}

.podium-position {
  font-size: 3rem;
  font-weight: bold;
  color: gold;
  text-shadow: 0 0 8px rgba(0, 0, 0, 0.7);
  margin-bottom: -1.5rem;
  z-index: 10;
}


.runner-up {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.position-badge {
  position: absolute;
  top: -20px;
  background: #4a5568;
  color: white;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 1.2rem;
  z-index: 2;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.change-ui-button {
  margin-top: 7rem;
  z-index: 10;
}

/* Responsive adjustments */
@media (max-width: 992px) {
  .podium-row {
    gap: 1rem;
  }

 

  .rank-one {
    transform: scale(1.15);
  }

  .rank-two,
  .rank-three {
    transform: scale(1.05);
  }
}

@media (max-width: 768px) {
  .podium-row {
    flex-direction: column;
    align-items: center;
    gap: 3rem;
  }

  .podium-item {
    width: 80%;
  }



  .rank-one,
  .rank-two,
  .rank-three {
    transform: scale(1);
    order: 0;
  }

  .podium-position {
    font-size: 2.5rem;
  }
}

@media (max-width: 576px) {
  .podium-item {
    width: 95%;
  }

  .change-ui-button {
    margin-top: 2rem;
  }
}
</style>
