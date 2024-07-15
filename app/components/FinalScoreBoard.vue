<script setup>
import { isCorrectAnswer } from "~/composables/check_is_correct.js/";
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

const analysisData = reactive([]);
const userAnalysisEndpoint = "/analytics_board/user";

const userAccuracy = ref();
const userAnswerAnalysis = ref([]);
const userCorrectAnswer = ref(0);
const userTotalScore = ref(0);

import { useUserScoreboardData } from "~/store/userScoreboardData";
const userScoreboardData = useUserScoreboardData();
const { addData, resetStore } = userScoreboardData;

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
        resetStore();
        addData(scoreboardData);
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

if (!props.isAdmin) {
  async function getAnalysisDetails() {
    const { data, error } = await useFetch(
      () => url.value.api_url + userAnalysisEndpoint,
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
          userAnalysis();
        }
        if (error.value) {
          toast.error(app.$$Unauthorized);
        }
      },
      { immediate: true, deep: true }
    );
  }

  getAnalysisDetails();

  const userAnalysis = () => {
    analysisData.filter((item) => {
      const correctAnswersCount = JSON.parse(item.correct_answer).length;

      // to not to consider the survey questions
      if (correctAnswersCount != Object.keys(item.options).length) {
        //for correct/incorrect answer (question-wise) and count of correct answer for accuracy
        userAnswerAnalysis.value.push(
          isCorrectAnswer(item.selected_answer.String, item.correct_answer)
        );

        //for counting total score
        userTotalScore.value += item.calculated_score;
      }
    });
    // get the count of correct answers
    userCorrectAnswer.value = userAnswerAnalysis.value.filter(Boolean).length;

    userAccuracy.value = (
      (userCorrectAnswer.value / userAnswerAnalysis.value.length) *
      100
    ).toFixed(2);
  };
}

const showAnalysis = () => {
  router.push({
    path: "/admin/analysis",
    query: { active_quiz_id: activeQuizId.value },
  });
};
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
              <td :class="{ 'user-row': user.username === props.userName && !props.isAdmin }">{{ user.rank }}</td>
              <td v-if="props.isAdmin">
                {{ user.firstname }} <span>({{ user.username }})</span>
              </td>
              <td v-else :class="{ 'user-row': user.username === props.userName }">
                {{ user.firstname }}
                <span v-if="props?.userName === user.username">
                  &nbsp;({{ user.username }})
                </span>
              </td>
              <td :class="{ 'user-row': user.username === props.userName && !props.isAdmin }">{{ user.score }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <hr v-if="!props.isAdmin" />
      <div v-if="!props.isAdmin">
        <h3 class="text-center">Accuracy: {{ userAccuracy }}%</h3>
        <h3 class="text-center">Total Score: {{ userTotalScore }}</h3>
        <Frame
          v-for="(item, index) in analysisData"
          :key="index"
          :page-title="'Q' + (index + 1) + '. ' + item.question"
        >
          <ul style="list-style-type: none; padding-left: 0">
            <li
              v-for="(option, key) in item.options"
              :key="key"
              style="display: flex; align-items: center; padding-left: 20px"
            >
              <span
                v-if="item.correct_answer.includes(key)"
                style="margin-right: 10px"
                >&#10004;</span
              >
              <span
                v-if="
                  item.selected_answer.String.includes(key) &&
                  !item.correct_answer.includes(key)
                "
                style="margin-right: 10px"
              >
                &#10006;
              </span>
              <span>{{ key }}: {{ option }}</span>
            </li>
          </ul>
          <div
            style="
              display: flex;
              flex: 1;
              margin-top: 10px;
              border-top: 1px solid #ccc;
            "
          >
            <div
              v-if="item.response_time > 0"
              style="flex: 1; padding: 10px; border-right: 1px solid #ccc"
            >
              Response Time:
              {{ (item.response_time / 1000).toFixed(2) }} seconds
            </div>
            <div
              v-else
              style="flex: 1; padding: 10px; border-right: 1px solid #ccc"
            >
              Response Time: -
            </div>
            <div style="flex: 1; padding: 10px">
              {{ item.is_attend ? "Attempted" : "Not Attempted" }}
            </div>
          </div>
        </Frame>
      </div>
    </div>
    <button v-if="props.isAdmin" class="btn btn-primary" @click="showAnalysis">
      Show Analysis
    </button>
  </ClientOnly>
</template>

<style scoped>
.user-row {
  background-color: #8968CD !important;
  box-shadow: none;
}
</style>