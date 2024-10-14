<script setup>
import AnswerSubmissionChart from "../AnswerSubmissionChart.vue";
const props = defineProps({
  data: {
    default: () => {
      return {};
    },
    type: Object,
    required: true,
  },
  isAdmin: {
    default: false,
    type: Boolean,
    required: false,
  },
  userName: {
    required: false,
    type: String,
    default: "",
  },
  selectedAnswer: {
    required: false,
    type: Number,
    default: 0,
  },
  analysisTab: {
    type: String,
    default: "",
  },
});

const emits = defineEmits(["askSkipTimer", "changeAnalysisTab"]);
const timer = ref(null);
const time = ref(0);
const isSkip = ref(false);

function handleTimer() {
  clearInterval(timer.value);
  timer.value = setInterval(() => {
    time.value += 0.1;
    if (time.value == props.data.data.duration + 1) {
      clearInterval(timer.value);
      time.value = -1;
      timer.value = null;
    }
  }, 100);
}

handleTimer();

function handleSkipTimer(e) {
  e.preventDefault();
  isSkip.value = true;
  emits("askSkipTimer");
}

const changeAnalysisTab = (tab) => emits("changeAnalysisTab", tab);
</script>

<template>
  <Frame
    page-title="Score Page"
    :music-component="true"
    page-message="Rank Board"
  >
    <v-progress-linear
      :striped="true"
      color="blue"
      :height="10"
      rounded="true"
      :model-value="(time * 100) / props.data.data.duration"
    ></v-progress-linear>

    <!-- Question -->
    <QuizQuestionAnalysis :question="props.data?.data" :is-for-quiz="true" />

    <!-- Options -->
    <div class="row d-flex align-items-stretch m-2">
      <div
        v-for="(answer, key) in props.data.data.options"
        :key="key"
        class="col-lg-6 col-md-12"
      >
        <div v-if="answer.isAnswer" class="bg-light-success option-box">
          <Option
            :order="Number(key)"
            :option="answer?.value"
            :is-correct="true"
            :options-media="props.data.data.options_media"
          />
        </div>
        <div
          v-else-if="!answer.isAnswer && key == props.selectedAnswer"
          class="bg-light-danger option-box wrong-option"
        >
          <Option
            :order="Number(key)"
            :option="answer?.value"
            :options-media="props.data.data.options_media"
          />
        </div>
        <div v-else class="option-box wrong-option">
          <Option
            :order="Number(key)"
            :option="answer?.value"
            :options-media="props.data.data.options_media"
          />
        </div>
      </div>
    </div>
    <button
      v-if="isAdmin"
      type="button"
      class="btn text-white btn-primary mt-3"
      :disabled="isSkip"
      @click="handleSkipTimer"
    >
      Skip
    </button>
    <ul
      v-if="isAdmin"
      id="pills-tab"
      class="nav nav-tabs mt-3 mb-3 nav-justified"
      role="tablist"
    >
      <!-- Ranking tab -->
      <li class="nav-item">
        <a
          id="pills-ranking-tab"
          class="nav-link"
          :class="{ active: props.analysisTab === 'ranking' }"
          data-bs-toggle="pill"
          href="#pills-ranking"
          role="tab"
          aria-controls="pills-ranking"
          aria-selected="true"
          @click="changeAnalysisTab('ranking')"
          >Rankings</a
        >
      </li>

      <!-- Charts tab -->
      <li class="nav-item">
        <a
          id="pills-chart-tab"
          class="nav-link"
          :class="{ active: props.analysisTab === 'chart' }"
          data-bs-toggle=""
          href="#pills-chart"
          role="tab"
          aria-controls="pills-chart"
          aria-selected="true"
          @click="changeAnalysisTab('chart')"
          >Chart</a
        >
      </li>
    </ul>
    <!-- Ranking Table -->
    <div id="pills-tabContent" class="tab-content p-4">
      <div
        id="pills-ranking"
        class="tab-pane fade"
        :class="{ 'active show': props.analysisTab === 'ranking' || !isAdmin }"
        role="tabpanel"
        aria-labelledby="pills-ranking-tab"
      >
        <div class="table-responsive mt-2">
          <table class="table table-borderless align-middle">
            <thead class="table-light"></thead>
            <tbody class="table-group-divider">
              <tr
                v-for="(user, index) in props.data.data.rankList"
                :key="index"
              >
                <td :class="{ 'bg-primary': user.username === props.userName }">
                  {{ user.rank }}
                </td>
                <td :class="{ 'bg-primary': user.username === props.userName }">
                  <img
                    :src="`${getAvatarUrlByName(user?.img_key)}&scale=75`"
                    alt="Avatar"
                    height="50px"
                  />
                  {{ user.firstname }}
                  <span v-if="user.username === props.userName"
                    >&nbsp; ({{ user.username }})</span
                  >
                </td>
                <td :class="{ 'bg-primary': user.username === props.userName }">
                  {{ user.score }}
                </td>
              </tr>
            </tbody>
            <tfoot></tfoot>
          </table>
        </div>
      </div>
      <!-- Chart -->
      <div
        id="pills-chart"
        class="tab-pane fade"
        :class="{ 'active show': props.analysisTab === 'chart' }"
        role="tabpanel"
        aria-labelledby="pills-chart-tab"
      >
        <AnswerSubmissionChart
          v-if="isAdmin"
          :options="props.data.data?.options"
          :options-media="props.data.data?.options_media"
          :responses="props.data.data?.userResponses"
        />
      </div>
    </div>
  </Frame>
</template>

<style scoped>
.table-custom {
  background-color: #f9f9f9; /* Example light color */
}

.table td,
.table th {
  background-color: #f9f9f9; /* Ensure each cell has a white background */
}

.option-box {
  min-height: 70px;
  padding-top: 3px;
  border-radius: 30px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 5px;
  margin-top: 3px;
}

.wrong-option {
  border: 1px solid var(--bs-light-primary);
}
</style>
